package braintree

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
	XMLName                 string               `xml:"subscription"`
	Id                      string               `xml:"id,omitempty"`
	Balance                 *Decimal             `xml:"balance,omitempty"`
	BillingDayOfMonth       string               `xml:"billing-day-of-month,omitempty"`
	BillingPeriodEndDate    string               `xml:"billing-period-end-date,omitempty"`
	BillingPeriodStartDate  string               `xml:"billing-period-start-date,omitempty"`
	CurrentBillingCycle     string               `xml:"current-billing-cycle,omitempty"`
	DaysPastDue             string               `xml:"days-past-due,omitempty"`
	FailureCount            string               `xml:"failure-count,omitempty"`
	FirstBillingDate        string               `xml:"first-billing-date,omitempty"`
	MerchantAccountId       string               `xml:"merchant-account-id,omitempty"`
	NeverExpires            bool                 `xml:"never-expires,omitempty"`
	NextBillAmount          *Decimal             `xml:"next-bill-amount,omitempty"`
	NextBillingPeriodAmount *Decimal             `xml:"next-billing-period-amount,omitempty"`
	NextBillingDate         string               `xml:"next-billing-date,omitempty"`
	NumberOfBillingCycles   *int                 `xml:"number-of-billing-cycles,omitempty"`
	PaidThroughDate         string               `xml:"paid-through-date,omitempty"`
	PaymentMethodToken      string               `xml:"payment-method-token,omitempty"`
	PlanId                  string               `xml:"plan-id,omitempty"`
	Price                   *Decimal             `xml:"price,omitempty"`
	Status                  SubscriptionStatus   `xml:"status,omitempty"`
	TrialDuration           string               `xml:"trial-duration,omitempty"`
	TrialDurationUnit       string               `xml:"trial-duration-unit,omitempty"`
	TrialPeriod             bool                 `xml:"trial-period,omitempty"`
	Transactions            *Transactions        `xml:"transactions,omitempty"`
	Options                 *SubscriptionOptions `xml:"options,omitempty"`
	Descriptor              *Descriptor          `xml:"descriptor,omitempty"`
	AddOns                  *AddOnList           `xml:"add-ons,omitempty"`
	Discounts               *DiscountList        `xml:"discounts,omitempty"`
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
