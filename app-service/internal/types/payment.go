package types

// PaymentRequest predstavlja zahtevek za plačilo
type PaymentRequest struct {
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
	ProductName string `json:"productName"`
	Quantity    int    `json:"quantity"`
}

// PaymentResponse predstavlja odgovor na zahtevek za plačilo
type PaymentResponse struct {
	CheckoutURL string `json:"checkoutUrl"`
}
