package braintree

import (
	"encoding/xml"
	"time"
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive       SubscriptionStatus = "Active"
	SubscriptionStatusCanceled     SubscriptionStatus = "Canceled"
	SubscriptionStatusExpired      SubscriptionStatus = "Expired"
	SubscriptionStatusPastDue      SubscriptionStatus = "Past Due"
	SubscriptionStatusPending      SubscriptionStatus = "Pending"
	SubscriptionStatusUnrecognized SubscriptionStatus = "Unrecognized"
)

const (
	SubscriptionTrialDurationUnitDay   = "day"
	SubscriptionTrialDurationUnitMonth = "month"
)

type Subscription struct {
	XMLName                 string                     `xml:"subscription"`
	Id                      string                     `xml:"id"`
	Balance                 *Decimal                   `xml:"balance"`
	BillingDayOfMonth       string                     `xml:"billing-day-of-month"`
	BillingPeriodEndDate    string                     `xml:"billing-period-end-date"`
	BillingPeriodStartDate  string                     `xml:"billing-period-start-date"`
	CurrentBillingCycle     string                     `xml:"current-billing-cycle"`
	DaysPastDue             string                     `xml:"days-past-due"`
	FailureCount            string                     `xml:"failure-count"`
	FirstBillingDate        string                     `xml:"first-billing-date"`
	MerchantAccountId       string                     `xml:"merchant-account-id"`
	NeverExpires            bool                       `xml:"never-expires"`
	NextBillAmount          *Decimal                   `xml:"next-bill-amount"`
	NextBillingPeriodAmount *Decimal                   `xml:"next-billing-period-amount"`
	NextBillingDate         string                     `xml:"next-billing-date"`
	NumberOfBillingCycles   *int                       `xml:"number-of-billing-cycles"`
	PaidThroughDate         string                     `xml:"paid-through-date"`
	PaymentMethodToken      string                     `xml:"payment-method-token"`
	PlanId                  string                     `xml:"plan-id"`
	Price                   *Decimal                   `xml:"price"`
	Status                  SubscriptionStatus         `xml:"status"`
	TrialDuration           string                     `xml:"trial-duration"`
	TrialDurationUnit       string                     `xml:"trial-duration-unit"`
	TrialPeriod             bool                       `xml:"trial-period"`
	Transactions            *Transactions              `xml:"transactions"`
	Options                 *SubscriptionOptions       `xml:"options"`
	StatusEvents            []*SubscriptionStatusEvent `xml:"status-history>status-event"`
	Descriptor              *Descriptor                `xml:"descriptor"`
	AddOns                  *AddOnList                 `xml:"add-ons"`
	Discounts               *DiscountList              `xml:"discounts"`
	CreatedAt               *time.Time                 `xml:"created-at,omitempty"`
	UpdatedAt               *time.Time                 `xml:"updated-at,omitempty"`
}

type SubscriptionRequest struct {
	XMLName               string                `xml:"subscription"`
	Id                    string                `xml:"id,omitempty"`
	BillingDayOfMonth     *int                  `xml:"billing-day-of-month,omitempty"`
	FailureCount          string                `xml:"failure-count,omitempty"`
	FirstBillingDate      string                `xml:"first-billing-date,omitempty"`
	MerchantAccountId     string                `xml:"merchant-account-id,omitempty"`
	NeverExpires          *bool                 `xml:"never-expires,omitempty"`
	NumberOfBillingCycles *int                  `xml:"number-of-billing-cycles,omitempty"`
	Options               *SubscriptionOptions  `xml:"options,omitempty"`
	PaymentMethodNonce    string                `xml:"paymentMethodNonce,omitempty"`
	PaymentMethodToken    string                `xml:"paymentMethodToken,omitempty"`
	PlanId                string                `xml:"planId,omitempty"`
	Price                 *Decimal              `xml:"price,omitempty"`
	TrialDuration         string                `xml:"trial-duration,omitempty"`
	TrialDurationUnit     string                `xml:"trial-duration-unit,omitempty"`
	TrialPeriod           *bool                 `xml:"trial-period,omitempty"`
	Descriptor            *Descriptor           `xml:"descriptor,omitempty"`
	AddOns                *ModificationsRequest `xml:"add-ons,omitempty"`
	Discounts             *ModificationsRequest `xml:"discounts,omitempty"`
}

type Subscriptions struct {
	Subscription []*Subscription `xml:"subscription"`
}

type SubscriptionOptions struct {
	DoNotInheritAddOnsOrDiscounts        bool `xml:"do-not-inherit-add-ons-or-discounts,omitempty"`
	ProrateCharges                       bool `xml:"prorate-charges,omitempty"`
	ReplaceAllAddOnsAndDiscounts         bool `xml:"replace-all-add-ons-and-discounts,omitempty"`
	RevertSubscriptionOnProrationFailure bool `xml:"revert-subscription-on-proration-failure,omitempty"`
	StartImmediately                     bool `xml:"start-immediately,omitempty"`
}

type SubscriptionTransactionRequest struct {
	Amount         *Decimal
	SubscriptionID string
	Options        *SubscriptionTransactionOptionsRequest
}

func (s *SubscriptionTransactionRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	x := struct {
		XMLName        xml.Name                               `xml:"transaction"`
		Type           string                                 `xml:"type"`
		SubscriptionID string                                 `xml:"subscription-id"`
		Amount         *Decimal                               `xml:"amount,omitempty"`
		Options        *SubscriptionTransactionOptionsRequest `xml:"options,omitempty"`
	}{
		Type:           "sale",
		SubscriptionID: s.SubscriptionID,
		Amount:         s.Amount,
		Options:        s.Options,
	}

	return e.Encode(x)
}

type SubscriptionTransactionOptionsRequest struct {
	SubmitForSettlement bool `xml:"submit-for-settlement"`
}

type SubscriptionSearchResult struct {
	TotalItems int
	TotalIDs   []string

	CurrentPageNumber int
	PageSize          int
	Subscriptions     []*Subscription
}
