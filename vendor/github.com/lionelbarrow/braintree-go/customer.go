package braintree

import (
	"github.com/lionelbarrow/braintree-go/customfields"
	"time"
)

type Customer struct {
	XMLName            string                    `xml:"customer"`
	Id                 string                    `xml:"id"`
	FirstName          string                    `xml:"first-name"`
	LastName           string                    `xml:"last-name"`
	Company            string                    `xml:"company"`
	Email              string                    `xml:"email"`
	Phone              string                    `xml:"phone"`
	Fax                string                    `xml:"fax"`
	Website            string                    `xml:"website"`
	CustomFields       customfields.CustomFields `xml:"custom-fields"`
	CreditCard         *CreditCard               `xml:"credit-card"`
	CreditCards        *CreditCards              `xml:"credit-cards"`
	PayPalAccounts     *PayPalAccounts           `xml:"paypal-accounts"`
	VenmoAccounts      *VenmoAccounts            `xml:"venmo-accounts"`
	AndroidPayCards    *AndroidPayCards          `xml:"android-pay-cards"`
	ApplePayCards      *ApplePayCards            `xml:"apple-pay-cards"`
	PaymentMethodNonce string                    `xml:"payment-method-nonce"`
	Addresses          *Addresses                `xml:"addresses"`
	CreatedAt          *time.Time                `xml:"created-at"`
	UpdatedAt          *time.Time                `xml:"updated-at"`
}

// PaymentMethods returns a slice of all PaymentMethods this customer has
func (c *Customer) PaymentMethods() []PaymentMethod {
	var paymentMethods []PaymentMethod
	paymentMethods = append(paymentMethods, c.CreditCards.PaymentMethods()...)
	paymentMethods = append(paymentMethods, c.PayPalAccounts.PaymentMethods()...)
	paymentMethods = append(paymentMethods, c.VenmoAccounts.PaymentMethods()...)
	paymentMethods = append(paymentMethods, c.AndroidPayCards.PaymentMethods()...)
	paymentMethods = append(paymentMethods, c.ApplePayCards.PaymentMethods()...)
	return paymentMethods
}

// DefaultCreditCard returns the default credit card, or nil
func (c *Customer) DefaultCreditCard() *CreditCard {
	for _, card := range c.CreditCards.CreditCard {
		if card.Default {
			return card
		}
	}
	return nil
}

// DefaultPaymentMethod returns the default payment method, or nil
func (c *Customer) DefaultPaymentMethod() PaymentMethod {
	for _, pm := range c.PaymentMethods() {
		if pm.IsDefault() {
			return pm
		}
	}
	return nil
}

type CustomerRequest struct {
	XMLName            string                    `xml:"customer"`
	ID                 string                    `xml:"id,omitempty"`
	FirstName          string                    `xml:"first-name,omitempty"`
	LastName           string                    `xml:"last-name,omitempty"`
	Company            string                    `xml:"company,omitempty"`
	Email              string                    `xml:"email,omitempty"`
	Phone              string                    `xml:"phone,omitempty"`
	Fax                string                    `xml:"fax,omitempty"`
	Website            string                    `xml:"website,omitempty"`
	CustomFields       customfields.CustomFields `xml:"custom-fields,omitempty"`
	CreditCard         *CreditCard               `xml:"credit-card,omitempty"`
	PaymentMethodNonce string                    `xml:"payment-method-nonce,omitempty"`
}

type CustomerSearchResult struct {
	TotalItems int
	TotalIDs   []string

	CurrentPageNumber int
	PageSize          int
	Customers         []*Customer
}
