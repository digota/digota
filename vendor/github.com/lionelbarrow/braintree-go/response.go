package braintree

import (
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/lionelbarrow/braintree-go/xmlnil"
)

type Response struct {
	*http.Response
	Body []byte
}

// TODO: remove dedicated unmarshal methods (redundant)

func (r *Response) entityName() (string, error) {
	var b struct {
		XMLName xml.Name
	}
	if err := xml.Unmarshal(r.Body, &b); err != nil {
		return "", err
	}
	return b.XMLName.Local, nil
}

func (r *Response) merchantAccount() (*MerchantAccount, error) {
	var b MerchantAccount
	if err := xml.Unmarshal(r.Body, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Response) transaction() (*Transaction, error) {
	var b Transaction
	if err := xml.Unmarshal(r.Body, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Response) paymentMethod() (PaymentMethod, error) {
	entityName, err := r.entityName()
	if err != nil {
		return nil, err
	}

	switch entityName {
	case "credit-card":
		return r.creditCard()
	case "paypal-account":
		return r.paypalAccount()
	case "venmo-account":
		return r.venmoAccount()
	case "android-pay-card":
		return r.androidPayCard()
	case "apple-pay-card":
		return r.applePayCard()
	}

	return nil, fmt.Errorf("Unrecognized payment method %#v", entityName)
}

func (r *Response) creditCard() (*CreditCard, error) {
	var b CreditCard
	if err := xml.Unmarshal(r.Body, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Response) paypalAccount() (*PayPalAccount, error) {
	var b PayPalAccount
	if err := xml.Unmarshal(r.Body, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Response) venmoAccount() (*VenmoAccount, error) {
	var b VenmoAccount
	if err := xml.Unmarshal(r.Body, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Response) androidPayCard() (*AndroidPayCard, error) {
	var b AndroidPayCard
	if err := xml.Unmarshal(r.Body, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Response) applePayCard() (*ApplePayCard, error) {
	var b ApplePayCard
	if err := xml.Unmarshal(r.Body, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Response) customer() (*Customer, error) {
	var b Customer
	if err := xml.Unmarshal(r.Body, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Response) subscription() (*Subscription, error) {
	var b Subscription
	if err := xml.Unmarshal(r.Body, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Response) settlement() (*SettlementBatchSummary, error) {
	var b SettlementBatchSummary
	if err := xml.Unmarshal(r.Body, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Response) address() (*Address, error) {
	var b Address
	if err := xml.Unmarshal(r.Body, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Response) addOns() ([]AddOn, error) {
	var b AddOnList
	if err := xml.Unmarshal(r.Body, &b); err != nil {
		return nil, err
	}
	return b.AddOns, nil
}

func (r *Response) discounts() ([]Discount, error) {
	var b DiscountList
	if err := xml.Unmarshal(r.Body, &b); err != nil {
		return nil, err
	}
	return b.Discounts, nil
}

func (r *Response) unpackBody() error {
	if len(r.Body) == 0 {
		b, err := gzip.NewReader(r.Response.Body)
		if err != nil {
			return err
		}
		defer func() { _ = r.Response.Body.Close() }()

		buf, err := ioutil.ReadAll(b)
		if err != nil {
			return err
		}
		strippedBuf, err := xmlnil.StripNilElements(buf)
		if err == nil {
			r.Body = strippedBuf
		} else {
			r.Body = buf
		}
	}
	return nil
}

func (r *Response) apiError() error {
	var b BraintreeError
	err := xml.Unmarshal(r.Body, &b)
	if err == nil && b.ErrorMessage != "" {
		b.statusCode = r.StatusCode
		return &b
	}
	if r.StatusCode > 299 {
		return httpError(r.StatusCode)
	}
	return nil
}

type APIError interface {
	error
	StatusCode() int
}

type httpError int

func (e httpError) StatusCode() int {
	return int(e)
}

func (e httpError) Error() string {
	return fmt.Sprintf("%s (%d)", http.StatusText(e.StatusCode()), e.StatusCode())
}

type invalidResponseError struct {
	resp *Response
}

type InvalidResponseError interface {
	error
	Response() *Response
}

func (e *invalidResponseError) Error() string {
	return fmt.Sprintf("braintree returned invalid response (%d)", e.resp.StatusCode)
}

func (e *invalidResponseError) Response() *Response {
	return e.resp
}
