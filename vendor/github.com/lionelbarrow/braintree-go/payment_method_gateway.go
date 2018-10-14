package braintree

import "encoding/xml"

type PaymentMethodGateway struct {
	*Braintree
}

type PaymentMethodRequest struct {
	XMLName            xml.Name                     `xml:"payment-method"`
	CustomerId         string                       `xml:"customer-id,omitempty"`
	Token              string                       `xml:"token,omitempty"`
	PaymentMethodNonce string                       `xml:"payment-method-nonce,omitempty"`
	Options            *PaymentMethodRequestOptions `xml:"options,omitempty"`
}

type PaymentMethodRequestOptions struct {
	MakeDefault                   bool   `xml:"make-default,omitempty"`
	FailOnDuplicatePaymentMethod  bool   `xml:"fail-on-duplicate-payment-method,omitempty"`
	VerifyCard                    *bool  `xml:"verify-card,omitempty"`
	VerificationMerchantAccountId string `xml:"verification-merchant-account-id,omitempty"`
}

func (g *PaymentMethodGateway) Create(paymentMethodRequest *PaymentMethodRequest) (PaymentMethod, error) {
	resp, err := g.executeVersion("POST", "payment_methods", paymentMethodRequest, apiVersion4)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 201:
		return resp.paymentMethod()
	}
	return nil, &invalidResponseError{resp}
}

func (g *PaymentMethodGateway) Update(token string, paymentMethod *PaymentMethodRequest) (PaymentMethod, error) {
	resp, err := g.executeVersion("PUT", "payment_methods/any/"+token, paymentMethod, apiVersion4)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.paymentMethod()
	}
	return nil, &invalidResponseError{resp}
}

func (g *PaymentMethodGateway) Find(token string) (PaymentMethod, error) {
	resp, err := g.executeVersion("GET", "payment_methods/any/"+token, nil, apiVersion4)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.paymentMethod()
	}
	return nil, &invalidResponseError{resp}
}

func (g *PaymentMethodGateway) Delete(token string) error {
	resp, err := g.executeVersion("DELETE", "payment_methods/any/"+token, nil, apiVersion4)
	if err != nil {
		return err
	}
	switch resp.StatusCode {
	case 200:
		return nil
	}
	return &invalidResponseError{resp}
}
