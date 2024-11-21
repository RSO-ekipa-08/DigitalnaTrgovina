package main

import (
	"context"
	"log"
	"net/http"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	paymentv1 "github.com/RSO-ekipa-08/DigitalnaTrgovina/gen/payment/v1"
	paymentv1connect "github.com/RSO-ekipa-08/DigitalnaTrgovina/gen/payment/v1/paymentv1connect"
	"github.com/stripe/stripe-go/v74/checkout/session"

	"github.com/stripe/stripe-go/v74"
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
    }

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

func main() {
	stripe.Key = "stripepk_test_51QKFbdByxGhH6PCLaKk13SRmwAlYySzPUkJXlm7USqcYya8Im6QrattPgJUEXFznnmRnAlgb3WFj4DOtXCBuWo6K004b1oTwIn"

	payServer := &PayServer{}
	mux := http.NewServeMux()
	path, handler := paymentv1connect.NewPaymentServiceHandler(payServer)
	mux.Handle(path, handler)
	http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
