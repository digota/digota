package braintree

import (
	"encoding/xml"
	"time"
)

type PayPalAccount struct {
	XMLName       xml.Name              `xml:"paypal-account"`
	CustomerId    string                `xml:"customer-id,omitempty"`
	Token         string                `xml:"token,omitempty"`
	Email         string                `xml:"email,omitempty"`
	ImageURL      string                `xml:"image-url,omitempty"`
	CreatedAt     *time.Time            `xml:"created-at,omitempty"`
	UpdatedAt     *time.Time            `xml:"updated-at,omitempty"`
	Subscriptions *Subscriptions        `xml:"subscriptions,omitempty"`
	Default       bool                  `xml:"default,omitempty"`
	Options       *PayPalAccountOptions `xml:"options,omitempty"`
}

type PayPalAccounts struct {
	PayPalAccount []*PayPalAccount `xml:"paypal-account"`
}

func (pp *PayPalAccounts) PaymentMethods() []PaymentMethod {
	if pp == nil {
		return nil
	}
	var paymentMethods []PaymentMethod
	for _, pp := range pp.PayPalAccount {
		paymentMethods = append(paymentMethods, pp)
	}
	return paymentMethods
}

type PayPalAccountOptions struct {
	MakeDefault bool `xml:"make-default,omitempty"`
}

func (paypalAccount *PayPalAccount) GetCustomerId() string {
	return paypalAccount.CustomerId
}

func (paypalAccount *PayPalAccount) GetToken() string {
	return paypalAccount.Token
}

func (paypalAccount *PayPalAccount) IsDefault() bool {
	return paypalAccount.Default
}

func (paypalAccount *PayPalAccount) GetImageURL() string {
	return paypalAccount.ImageURL
}

// AllSubscriptions returns all subscriptions for this paypal account, or nil if none present.
func (paypalAccount *PayPalAccount) AllSubscriptions() []*Subscription {
	if paypalAccount.Subscriptions != nil {
		subs := paypalAccount.Subscriptions.Subscription
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
