package beans
type StripePayment struct {
	Amount    int64  `json:"amount"`
	Currency  string `json:"currency"`
	CustomerID string `json:"customer_id"`
	PaymentMethodTypes []string `json:"payment_method_types"`
}
