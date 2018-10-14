package braintree

import (
	"encoding/xml"
	"time"
)

type CreditCard struct {
	XMLName                   xml.Name           `xml:"credit-card"`
	CustomerId                string             `xml:"customer-id,omitempty"`
	Token                     string             `xml:"token,omitempty"`
	PaymentMethodNonce        string             `xml:"payment-method-nonce,omitempty"`
	Number                    string             `xml:"number,omitempty"`
	ExpirationDate            string             `xml:"expiration-date,omitempty"`
	ExpirationMonth           string             `xml:"expiration-month,omitempty"`
	ExpirationYear            string             `xml:"expiration-year,omitempty"`
	CVV                       string             `xml:"cvv,omitempty"`
	VenmoSDKPaymentMethodCode string             `xml:"venmo-sdk-payment-method-code,omitempty"`
	VenmoSDK                  bool               `xml:"venmo-sdk,omitempty"`
	Options                   *CreditCardOptions `xml:"options,omitempty"`
	CreatedAt                 *time.Time         `xml:"created-at,omitempty"`
	UpdatedAt                 *time.Time         `xml:"updated-at,omitempty"`
	Bin                       string             `xml:"bin,omitempty"`
	CardType                  string             `xml:"card-type,omitempty"`
	CardholderName            string             `xml:"cardholder-name,omitempty"`
	CustomerLocation          string             `xml:"customer-location,omitempty"`
	ImageURL                  string             `xml:"image-url,omitempty"`
	Default                   bool               `xml:"default,omitempty"`
	Expired                   bool               `xml:"expired,omitempty"`
	Last4                     string             `xml:"last-4,omitempty"`
	Commercial                string             `xml:"commercial,omitempty"`
	Debit                     string             `xml:"debit,omitempty"`
	DurbinRegulated           string             `xml:"durbin-regulated,omitempty"`
	Healthcare                string             `xml:"healthcare,omitempty"`
	Payroll                   string             `xml:"payroll,omitempty"`
	Prepaid                   string             `xml:"prepaid,omitempty"`
	CountryOfIssuance         string             `xml:"country-of-issuance,omitempty"`
	IssuingBank               string             `xml:"issuing-bank,omitempty"`
	UniqueNumberIdentifier    string             `xml:"unique-number-identifier,omitempty"`
	BillingAddress            *Address           `xml:"billing-address,omitempty"`
	Subscriptions             *Subscriptions     `xml:"subscriptions,omitempty"`
}

type CreditCards struct {
	CreditCard []*CreditCard `xml:"credit-card"`
}

func (cards *CreditCards) PaymentMethods() []PaymentMethod {
	if cards == nil {
		return nil
	}
	var paymentMethods []PaymentMethod
	for _, cc := range cards.CreditCard {
		paymentMethods = append(paymentMethods, cc)
	}
	return paymentMethods
}

type CreditCardOptions struct {
	VerifyCard                    *bool  `xml:"verify-card,omitempty"`
	VenmoSDKSession               string `xml:"venmo-sdk-session,omitempty"`
	MakeDefault                   bool   `xml:"make-default,omitempty"`
	FailOnDuplicatePaymentMethod  bool   `xml:"fail-on-duplicate-payment-method,omitempty"`
	VerificationMerchantAccountId string `xml:"verification-merchant-account-id,omitempty"`
	UpdateExistingToken           string `xml:"update-existing-token,omitempty"`
}

// GetCustomerId gets the customer id of the credit card.
func (card *CreditCard) GetCustomerId() string {
	return card.CustomerId
}

// GetToken gets the payment method token of the credit card.
func (card *CreditCard) GetToken() string {
	return card.Token
}

// IsDefault returns whether the credit card is the default for the customer.
func (card *CreditCard) IsDefault() bool {
	return card.Default
}

// GetImageURL gets a URL that points to a payment method
// image file hosted by Braintree.
func (card *CreditCard) GetImageURL() string {
	return card.ImageURL
}

// AllSubscriptions returns all subscriptions for this card, or nil if none present.
func (card *CreditCard) AllSubscriptions() []*Subscription {
	if card.Subscriptions != nil {
		subs := card.Subscriptions.Subscription
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

type CreditCardSearchResult struct {
	TotalItems int
	TotalIDs   []string

	CurrentPageNumber int
	PageSize          int
	CreditCards       []*CreditCard
}
