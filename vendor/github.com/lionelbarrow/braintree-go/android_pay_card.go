package braintree

import (
	"encoding/xml"
	"time"
)

type AndroidPayCard struct {
	XMLName             xml.Name       `xml:"android-pay-card"`
	Token               string         `xml:"token,omitempty"`
	CardType            string         `xml:"-"`
	Last4               string         `xml:"-"`
	SourceCardType      string         `xml:"source-card-type,omitempty"`
	SourceCardLast4     string         `xml:"source-card-last-4,omitempty"`
	SourceDescription   string         `xml:"source-description,omitempty"`
	VirtualCardType     string         `xml:"virtual-card-type,omitempty"`
	VirtualCardLast4    string         `xml:"virtual-card-last-4,omitempty"`
	ExpirationMonth     string         `xml:"expiration-month,omitempty"`
	ExpirationYear      string         `xml:"expiration-year,omitempty"`
	BIN                 string         `xml:"bin,omitempty"`
	GoogleTransactionID string         `xml:"google-transaction-id,omitempty"`
	ImageURL            string         `xml:"image-url,omitempty"`
	Default             bool           `xml:"default,omitempty"`
	CustomerId          string         `xml:"customer-id,omitempty"`
	CreatedAt           *time.Time     `xml:"created-at,omitempty"`
	UpdatedAt           *time.Time     `xml:"updated-at,omitempty"`
	Subscriptions       *Subscriptions `xml:"subscriptions,omitempty"`
}

type AndroidPayCards struct {
	AndroidPayCard []*AndroidPayCard `xml:"android-pay-card"`
}

func (a *AndroidPayCards) PaymentMethods() []PaymentMethod {
	if a == nil {
		return nil
	}
	var paymentMethods []PaymentMethod
	for _, ac := range a.AndroidPayCard {
		paymentMethods = append(paymentMethods, ac)
	}
	return paymentMethods
}

func (a *AndroidPayCard) GetCustomerId() string {
	return a.CustomerId
}

func (a *AndroidPayCard) GetToken() string {
	return a.Token
}

func (a *AndroidPayCard) IsDefault() bool {
	return a.Default
}

func (a *AndroidPayCard) GetImageURL() string {
	return a.ImageURL
}

// AllSubscriptions returns all subscriptions for this paypal account, or nil if none present.
func (a *AndroidPayCard) AllSubscriptions() []*Subscription {
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

func (a *AndroidPayCard) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type typeWithNoFunctions AndroidPayCard
	if err := d.DecodeElement((*typeWithNoFunctions)(a), &start); err != nil {
		return err
	}
	a.CardType = a.VirtualCardType
	a.Last4 = a.VirtualCardLast4
	return nil
}
