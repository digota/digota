package braintree

type ApplePayDetails struct {
	Token                 string `xml:"token,omitempty"`
	CardType              string `xml:"card-type,omitempty"`
	PaymentInstrumentName string `xml:"payment-instrument-name,omitempty"`
	SourceDescription     string `xml:"source-description,omitempty"`
	CardholderName        string `xml:"cardholder-name,omitempty"`
	ExpirationMonth       string `xml:"expiration-month,omitempty"`
	ExpirationYear        string `xml:"expiration-year,omitempty"`
	Last4                 string `xml:"last-4,omitempty"`
}
