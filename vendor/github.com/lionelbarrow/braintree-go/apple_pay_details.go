package braintree

type ApplePayDetails struct {
	Token                 string `xml:"token"`
	CardType              string `xml:"card-type"`
	PaymentInstrumentName string `xml:"payment-instrument-name"`
	SourceDescription     string `xml:"source-description"`
	CardholderName        string `xml:"cardholder-name"`
	ExpirationMonth       string `xml:"expiration-month"`
	ExpirationYear        string `xml:"expiration-year"`
	Last4                 string `xml:"last-4"`
}
