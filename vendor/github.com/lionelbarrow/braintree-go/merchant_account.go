package braintree

type MerchantAccount struct {
	XMLName                 string                         `xml:"merchant-account,omitempty"`
	Id                      string                         `xml:"id,omitempty"`
	MasterMerchantAccountId string                         `xml:"master-merchant-account-id,omitempty"`
	TOSAccepted             bool                           `xml:"tos_accepted,omitempty"`
	Individual              *MerchantAccountPerson         `xml:"individual,omitempty"`
	Business                *MerchantAccountBusiness       `xml:"business,omitempty"`
	FundingOptions          *MerchantAccountFundingOptions `xml:"funding,omitempty"`
	Status                  string                         `xml:"status,omitempty"`
}

type MerchantAccountPerson struct {
	FirstName   string   `xml:"first-name,omitempty"`
	LastName    string   `xml:"last-name,omitempty"`
	Email       string   `xml:"email,omitempty"`
	Phone       string   `xml:"phone,omitempty"`
	DateOfBirth string   `xml:"date-of-birth,omitempty"`
	SSN         string   `xml:"ssn,omitempty"`
	Address     *Address `xml:"address,omitempty"`
}

type MerchantAccountBusiness struct {
	LegalName string   `xml:"legal-name,omitempty"`
	DbaName   string   `xml:"dba-name,omitempty"`
	TaxId     string   `xml:"tax-id,omitempty"`
	Address   *Address `xml:"address,omitempty"`
}

type MerchantAccountFundingOptions struct {
	Destination   string `xml:"destination,omitempty"`
	Email         string `xml:"email,omitempty"`
	MobilePhone   string `xml:"mobile-phone,omitempty"`
	AccountNumber string `xml:"account-number,omitempty"`
	RoutingNumber string `xml:"routing-number,omitempty"`
}

const (
	FUNDING_DEST_BANK         = "bank"
	FUNDING_DEST_MOBILE_PHONE = "mobile_phone"
	FUNDING_DEST_EMAIL        = "email"
)
