package braintree

import (
	"encoding/xml"
)

type DisbursementDetails struct {
	XMLName                        xml.Name `xml:"disbursement-details"`
	DisbursementDate               string   `xml:"disbursement-date"`
	SettlementAmount               *Decimal `xml:"settlement-amount"`
	SettlementCurrencyIsoCode      string   `xml:"settlement-currency-iso-code"`
	SettlementCurrencyExchangeRate *Decimal `xml:"settlement-currency-exchange-rate"`
	FundsHeld                      bool     `xml:"funds-held"`
	Success                        bool     `xml:"success"`
}
