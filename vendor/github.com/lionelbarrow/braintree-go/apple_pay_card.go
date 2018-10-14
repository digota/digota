package braintree

import (
	"encoding/xml"
	"time"
)

type ApplePayCard struct {
	XMLName               xml.Name       `xml:"apple-pay-card"`
	Token                 string         `xml:"token,omitempty"`
	ImageURL              string         `xml:"image-url,omitempty"`
	CardType              string         `xml:"card-type,omitempty"`
	PaymentInstrumentName string         `xml:"payment-instrument-name,omitempty"`
	SourceDescription     string         `xml:"source-description,omitempty"`
	Last4                 string         `xml:"last-4,omitempty"`
	ExpirationMonth       string         `xml:"expiration-month,omitempty"`
	ExpirationYear        string         `xml:"expiration-year,omitempty"`
	Expired               bool           `xml:"expired,omitempty"`
	Default               bool           `xml:"default,omitempty"`
	CustomerId            string         `xml:"customer-id,omitempty"`
	CreatedAt             *time.Time     `xml:"created-at,omitempty"`
	UpdatedAt             *time.Time     `xml:"updated-at,omitempty"`
	Subscriptions         *Subscriptions `xml:"subscriptions,omitempty"`
}

type ApplePayCards struct {
	ApplePayCard []*ApplePayCard `xml:"apple-pay-card"`
}

func (a *ApplePayCards) PaymentMethods() []PaymentMethod {
	if a == nil {
		return nil
	}
	var paymentMethods []PaymentMethod
	for _, a := range a.ApplePayCard {
		paymentMethods = append(paymentMethods, a)
	}
	return paymentMethods
}

func (a *ApplePayCard) GetCustomerId() string {
	return a.CustomerId
}

func (a *ApplePayCard) GetToken() string {
	return a.Token
}

func (a *ApplePayCard) IsDefault() bool {
	return a.Default
}

func (a *ApplePayCard) GetImageURL() string {
	return a.ImageURL
}

// AllSubscriptions returns all subscriptions for this paypal account, or nil if none present.
func (a *ApplePayCard) AllSubscriptions() []*Subscription {
	if a.Subscriptions != nil {
		subs := a.Subscriptions.Subscription
		if len(subs) > 0 {
			a := make([]*Subscription, 0, len(subs))
			for _, s := range subs {
				a = append(a, s)
			}
			return a
		}
	}
	return nil
}
