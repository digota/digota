package braintree

import (
	"encoding/xml"
	"time"
)

type ApplePayCard struct {
	XMLName               xml.Name       `xml:"apple-pay-card"`
	Token                 string         `xml:"token"`
	ImageURL              string         `xml:"image-url"`
	CardType              string         `xml:"card-type"`
	PaymentInstrumentName string         `xml:"payment-instrument-name"`
	SourceDescription     string         `xml:"source-description"`
	Last4                 string         `xml:"last-4"`
	ExpirationMonth       string         `xml:"expiration-month"`
	ExpirationYear        string         `xml:"expiration-year"`
	Expired               bool           `xml:"expired"`
	Default               bool           `xml:"default"`
	CustomerId            string         `xml:"customer-id"`
	CreatedAt             *time.Time     `xml:"created-at"`
	UpdatedAt             *time.Time     `xml:"updated-at"`
	Subscriptions         *Subscriptions `xml:"subscriptions"`
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
