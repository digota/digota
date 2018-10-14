package braintree

type TransactionCloneRequest struct {
	XMLName string                   `xml:"transaction-clone"`
	Amount  *Decimal                 `xml:"amount"`
	Channel string                   `xml:"channel"`
	Options *TransactionCloneOptions `xml:"options"`
}

type TransactionCloneOptions struct {
	SubmitForSettlement bool `xml:"submit-for-settlement"`
}
