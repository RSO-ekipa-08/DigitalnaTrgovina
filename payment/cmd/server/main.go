package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	paymentv1 "github.com/RSO-ekipa-08/DigitalnaTrgovina/gen/payment/v1"
	paymentv1connect "github.com/RSO-ekipa-08/DigitalnaTrgovina/gen/payment/v1/paymentv1connect"
)

type PayServer struct{}

func (s *PayServer) Pay(
	ctx context.Context,
	req *connect.Request[paymentv1.PayRequest],
) (*connect.Response[paymentv1.PayResponse], error) {
	log.Println("Request headers: ", req.Header())
	res := connect.NewResponse(&paymentv1.PayResponse{
		StripeLink: fmt.Sprintf("SomeStripeLink for %s!", req.Msg.UserId),
	})
	res.Header().Set("Payment-Version", "v1")
	return res, nil
}

func main() {
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
