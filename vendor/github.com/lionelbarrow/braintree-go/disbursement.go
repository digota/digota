package braintree

import (
	"encoding/xml"
	"github.com/lionelbarrow/braintree-go/date"
)

type Disbursement struct {
	XMLName            xml.Name         `xml:"disbursement"`
	Id                 string           `xml:"id"`
	ExceptionMessage   string           `xml:"exception-message"`
	DisbursementDate   *date.Date       `xml:"disbursement-date"`
	FollowUpAction     string           `xml:"follow-up-action"`
	Success            bool             `xml:"success"`
	Retry              bool             `xml:"retry"`
	Amount             *Decimal         `xml:"amount"`
	MerchantAccount    *MerchantAccount `xml:"merchant-account"`
	CurrencyIsoCode    string           `xml:"currency-iso-code"`
	SubmerchantAccount bool             `xml:"sub-merchant-account"`
	Status             string           `xml:"status"`
	TransactionIds     []string         `xml:"transaction-ids>item"`
}

const (
	// Exception messages
	BankRejected         = "bank_rejected"
	InsufficientFunds    = "insuffient_funds"
	AccountNotAuthorized = "account_not_authorized"

	// Followup actions
	ContactUs                = "contact_us"
	UpdateFundingInformation = "update_funding_information"
	None                     = "none"
)

func (d *Disbursement) Transactions(g *TransactionGateway) (*TransactionSearchResult, error) {
	query := new(SearchQuery)
	f := query.AddMultiField("ids")
	f.Items = d.TransactionIds

	result, err := g.Search(query)
	if err != nil {
		return nil, err
	}

	return result, err
}
