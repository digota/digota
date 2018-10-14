package braintree

import (
	"github.com/lionelbarrow/braintree-go/customfields"
)

type Customer struct {
	XMLName            string                    `xml:"customer"`
	Id                 string                    `xml:"id,omitempty"`
	FirstName          string                    `xml:"first-name,omitempty"`
	LastName           string                    `xml:"last-name,omitempty"`
	Company            string                    `xml:"company,omitempty"`
	Email              string                    `xml:"email,omitempty"`
	Phone              string                    `xml:"phone,omitempty"`
	Fax                string                    `xml:"fax,omitempty"`
	Website            string                    `xml:"website,omitempty"`
	CustomFields       customfields.CustomFields `xml:"custom-fields,omitempty"`
	CreditCard         *CreditCard               `xml:"credit-card,omitempty"`
	CreditCards        *CreditCards              `xml:"credit-cards,omitempty"`
	PayPalAccounts     *PayPalAccounts           `xml:"paypal-accounts,omitempty"`
	VenmoAccounts      *VenmoAccounts            `xml:"venmo-accounts,omitempty"`
	AndroidPayCards    *AndroidPayCards          `xml:"android-pay-cards,omitempty"`
	ApplePayCards      *ApplePayCards            `xml:"apple-pay-cards,omitempty"`
	PaymentMethodNonce string                    `xml:"payment-method-nonce,omitempty"`
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

type CustomerSearchResult struct {
	XMLName           string      `xml:"customers"`
	CurrentPageNumber int         `xml:"current-page-number"`
	PageSize          int         `xml:"page-size"`
	TotalItems        int         `xml:"total-items"`
	Customers         []*Customer `xml:"customer"`
}
