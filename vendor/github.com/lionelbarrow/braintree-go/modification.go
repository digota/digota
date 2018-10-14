package braintree

import "time"

const (
	ModificationKindDiscount = "discount"
	ModificationKindAddOn    = "add_on"
)

type Modification struct {
	Id                    string     `xml:"id"`
	Amount                *Decimal   `xml:"amount"`
	Description           string     `xml:"description"`
	Kind                  string     `xml:"kind"`
	Name                  string     `xml:"name"`
	NeverExpires          bool       `xml:"never-expires"`
	Quantity              int        `xml:"quantity"`
	NumberOfBillingCycles int        `xml:"number-of-billing-cycles"`
	CurrentBillingCycle   int        `xml:"current-billing-cycle"`
	UpdatedAt             *time.Time `xml:"updated_at"`
}
