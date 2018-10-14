package braintree

import (
	"encoding/xml"
	"time"
)

type AndroidPayCard struct {
	XMLName             xml.Name       `xml:"android-pay-card"`
	Token               string         `xml:"token"`
	CardType            string         `xml:"-"`
	Last4               string         `xml:"-"`
	SourceCardType      string         `xml:"source-card-type"`
	SourceCardLast4     string         `xml:"source-card-last-4"`
	SourceDescription   string         `xml:"source-description"`
	VirtualCardType     string         `xml:"virtual-card-type"`
	VirtualCardLast4    string         `xml:"virtual-card-last-4"`
	ExpirationMonth     string         `xml:"expiration-month"`
	ExpirationYear      string         `xml:"expiration-year"`
	BIN                 string         `xml:"bin"`
	GoogleTransactionID string         `xml:"google-transaction-id"`
	ImageURL            string         `xml:"image-url"`
	Default             bool           `xml:"default"`
	CustomerId          string         `xml:"customer-id"`
	CreatedAt           *time.Time     `xml:"created-at"`
	UpdatedAt           *time.Time     `xml:"updated-at"`
	Subscriptions       *Subscriptions `xml:"subscriptions"`
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
