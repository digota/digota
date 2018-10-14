package braintree

import (
	"encoding/xml"
	"time"
)

type VenmoAccount struct {
	XMLName           xml.Name       `xml:"venmo-account"`
	CustomerId        string         `xml:"customer-id,omitempty"`
	Token             string         `xml:"token,omitempty"`
	Username          string         `xml:"username,omitempty"`
	VenmoUserID       string         `xml:"venmo-user-id,omitempty"`
	SourceDescription string         `xml:"source-description,omitempty"`
	ImageURL          string         `xml:"image-url,omitempty"`
	CreatedAt         *time.Time     `xml:"created-at,omitempty"`
	UpdatedAt         *time.Time     `xml:"updated-at,omitempty"`
	Subscriptions     *Subscriptions `xml:"subscriptions,omitempty"`
	Default           bool           `xml:"default,omitempty"`
}

type VenmoAccounts struct {
	VenmoAccount []*VenmoAccount `xml:"venmo-account"`
}

func (v *VenmoAccounts) PaymentMethods() []PaymentMethod {
	if v == nil {
		return nil
	}
	var paymentMethods []PaymentMethod
	for _, va := range v.VenmoAccount {
		paymentMethods = append(paymentMethods, va)
	}
	return paymentMethods
}

func (v *VenmoAccount) GetCustomerId() string {
	return v.CustomerId
}

func (v *VenmoAccount) GetToken() string {
	return v.Token
}

func (v *VenmoAccount) IsDefault() bool {
	return v.Default
}

func (v *VenmoAccount) GetImageURL() string {
	return v.ImageURL
}

// AllSubscriptions returns all subscriptions for this venmo account, or nil if none present.
func (v *VenmoAccount) AllSubscriptions() []*Subscription {
	if v.Subscriptions != nil {
		subs := v.Subscriptions.Subscription
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
