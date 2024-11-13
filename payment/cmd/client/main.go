package main

import (
	"context"
	"log"
	"net/http"
	"fmt"

	paymentv1 "github.com/RSO-ekipa-08/DigitalnaTrgovina/gen/payment/v1"
	paymentv1connect "github.com/RSO-ekipa-08/DigitalnaTrgovina/gen/payment/v1/paymentv1connect"

	"connectrpc.com/connect"
)

func main() {
	client := paymentv1connect.NewPaymentServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
	)
	res, err := client.Pay(
		context.Background(),
		connect.NewRequest(&paymentv1.PayRequest{UserId: "Jane", AppId: "app"}),
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(fmt.Sprintf("link: %s", res.Msg.StripeLink))
}
