package braintree

import (
	"time"
)

type Plan struct {
	XMLName               string       `xml:"plan"`
	Id                    string       `xml:"id"`
	MerchantId            string       `xml:"merchant-id"`
	BillingDayOfMonth     *int         `xml:"billing-day-of-month"`
	BillingFrequency      *int         `xml:"billing-frequency"`
	CurrencyISOCode       string       `xml:"currency-iso-code"`
	Description           string       `xml:"description"`
	Name                  string       `xml:"name"`
	NumberOfBillingCycles *int         `xml:"number-of-billing-cycles"`
	Price                 *Decimal     `xml:"price"`
	TrialDuration         *int         `xml:"trial-duration"`
	TrialDurationUnit     string       `xml:"trial-duration-unit"`
	TrialPeriod           bool         `xml:"trial-period"`
	CreatedAt             *time.Time   `xml:"created-at"`
	UpdatedAt             *time.Time   `xml:"updated-at"`
	AddOns                AddOnList    `xml:"add-ons"`
	Discounts             DiscountList `xml:"discounts"`
}

type Plans struct {
	XMLName string  `xml:"plans"`
	Plan    []*Plan `xml:"plan"`
}
