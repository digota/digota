package braintree

import (
	"time"
)

// SubscriptionStatusEvent contains information about what and when
// something changed about the subscription.
type SubscriptionStatusEvent struct {
	Timestamp          time.Time          `xml:"timestamp"`
	Balance            *Decimal           `xml:"balance"`
	Price              *Decimal           `xml:"price"`
	Status             SubscriptionStatus `xml:"status"`
	CurrencyISOCode    string             `xml:"currency-iso-code"`
	User               string             `xml:"user"`
	PlanID             string             `xml:"plan-id"`
	SubscriptionSource string             `xml:"subscription-source"`
}
