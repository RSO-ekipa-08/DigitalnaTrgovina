package main

import (
	"context"
	"log"
	"net/http"
	"fmt"
  	"io/ioutil"
	"encoding/json"
	"os"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	paymentv1 "github.com/RSO-ekipa-08/DigitalnaTrgovina/gen/payment/v1"
	paymentv1connect "github.com/RSO-ekipa-08/DigitalnaTrgovina/gen/payment/v1/paymentv1connect"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
)


type PayServer struct{}

func (s *PayServer) Pay(
	ctx context.Context,
	req *connect.Request[paymentv1.PayRequest],
) (*connect.Response[paymentv1.PayResponse], error) {
	log.Println("Request headers: ", req.Header())
	
	stripe.Key = "sk_test_51QKFbdByxGhH6PCLUgWMXEkNABlWw4eo9gWxo7eEEuy6ZV9wGAam0kLxUJESNhvNg8fM4qmi1zuh4qCb7J7LbEHv00Lod1NWXL"
	// Create the payment link
	params := &stripe.CheckoutSessionParams{
        PaymentMethodTypes: stripe.StringSlice([]string{"card"}), // or add other payment methods
        LineItems: []*stripe.CheckoutSessionLineItemParams{
            {
                PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
                    Currency: stripe.String("usd"),
                    ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
                        Name: stripe.String("app"),
                    },
                    UnitAmount: stripe.Int64(2000), // amount in cents, e.g., $20.00
                },
                Quantity: stripe.Int64(1),
            },
        },
        Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
        SuccessURL: stripe.String("https://example.com/success"), // Replace with your success URL
        CancelURL: stripe.String("https://example.com/cancel"),   // Replace with your cancel URL
		Metadata: map[string]string{
			"user_id": req.Msg.UserId,
			"app_id": req.Msg.AppId,
		},
    }
	
		log.Printf("Sending metadata", req.Msg.UserId)

	// Create the checkout session
    session, err := session.New(params)
    if err != nil {
        log.Fatalf("Unable to create checkout session: %v", err)
    }
	
	res := connect.NewResponse(&paymentv1.PayResponse{
		StripeLink: session.URL,
	})
	res.Header().Set("Payment-Version", "v1")
	return res, nil
}

func handleWebhook(w http.ResponseWriter, req *http.Request) {
	log.Println("handling webhook")
	const MaxBodyBytes = int64(65536)
	req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(req.Body)
	if err != nil {
	  fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
	  w.WriteHeader(http.StatusServiceUnavailable)
	  return
	}
  
	event := stripe.Event{}
  
	if err := json.Unmarshal(payload, &event); err != nil {
	  fmt.Fprintf(os.Stderr, "⚠️  Webhook error while parsing basic request. %v\n", err.Error())
	  w.WriteHeader(http.StatusBadRequest)
	  return
	}
  
	// Replace this endpoint secret with your endpoint's unique secret
	// If you are testing with the CLI, find the secret by running 'stripe listen'
	// If you are using an endpoint defined with the API or dashboard, look in your webhook settings
	// at https://dashboard.stripe.com/webhooks
	endpointSecret := "whsec_1c0a48ebf1cfffb46eaf2d993aaa0abd0c2975aa8bfa5236425a845d6f53e45f"
	signatureHeader := req.Header.Get("Stripe-Signature")
	event, err = webhook.ConstructEvent(payload, signatureHeader, endpointSecret)
	if err != nil {
	  fmt.Fprintf(os.Stderr, "⚠️  Webhook signature verification failed. %v\n", err)
	  w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
	  return
	}
	// Unmarshal the event data into an appropriate struct depending on its Type
	switch event.Type {
	case "payment_intent.succeeded":
	  var paymentIntent stripe.PaymentIntent
	  err := json.Unmarshal(event.Data.Raw, &paymentIntent)
	  if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
	  }
	  log.Printf("Successful payment for %d.", paymentIntent.Amount)
	  metadata := paymentIntent.Metadata
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			log.Printf("Error parsing checkout session object: %v", err)
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}

		metadata := session.Metadata
		for key, value := range metadata {
			log.Printf("Metadata key: %s, value: %s", key, value)
		}

		userId := metadata["user_id"]
		appId := metadata["app_id"]
		fmt.Printf("User ID: %s, App ID: %s\n", userId, appId)
	default:
	  fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	}
  
	w.WriteHeader(http.StatusOK)
  }

func main() {
	stripe.Key = "stripepk_test_51QKFbdByxGhH6PCLaKk13SRmwAlYySzPUkJXlm7USqcYya8Im6QrattPgJUEXFznnmRnAlgb3WFj4DOtXCBuWo6K004b1oTwIn"

	payServer := &PayServer{}
	mux := http.NewServeMux()
	path, handler := paymentv1connect.NewPaymentServiceHandler(payServer)
	mux.Handle(path, handler)
	mux.HandleFunc("/webhook", handleWebhook)
	http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
