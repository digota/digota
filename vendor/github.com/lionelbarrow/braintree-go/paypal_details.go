package braintree

type PayPalDetails struct {
	PayerEmail                    string `xml:"payer-email,omitempty"`
	PaymentID                     string `xml:"payment-id,omitempty"`
	AuthorizationID               string `xml:"authorization-id,omitempty"`
	Token                         string `xml:"token,omitempty"`
	ImageURL                      string `xml:"image-url,omitempty"`
	DebugID                       string `xml:"debug-id,omitempty"`
	PayeeEmail                    string `xml:"payee-email,omitempty"`
	CustomField                   string `xml:"custom-field,omitempty"`
	PayerID                       string `xml:"payer-id,omitempty"`
	PayerFirstName                string `xml:"payer-first-name,omitempty"`
	PayerLastName                 string `xml:"payer-last-name,omitempty"`
	PayerStatus                   string `xml:"payer-status,omitempty"`
	SellerProtectionStatus        string `xml:"seller-protection-status,omitempty"`
	RefundID                      string `xml:"refund-id,omitempty"`
	CaptureID                     string `xml:"capture-id,omitempty"`
	TransactionFeeAmount          string `xml:"transaction-fee-amount,omitempty"`
	TransactionFeeCurrencyISOCode string `xml:"transaction-fee-currency-iso-code,omitempty"`
	Description                   string `xml:"description,omitempty"`
}
