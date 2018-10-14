package braintree

import (
	"context"
)

type PaymentMethodNonceGateway struct {
	*Braintree
}

func (g *PaymentMethodNonceGateway) Find(ctx context.Context, nonce string) (*PaymentMethodNonce, error) {
	resp, err := g.executeVersion(ctx, "GET", "/payment_method_nonces/"+nonce, nil, apiVersion4)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.paymentMethodNonce()
	}
	return nil, &invalidResponseError{resp}
}

func (g *PaymentMethodNonceGateway) Create(ctx context.Context, token string) (*PaymentMethodNonce, error) {
	resp, err := g.executeVersion(ctx, "POST", "/payment_methods/"+token+"/nonces", nil, apiVersion4)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 201:
		return resp.paymentMethodNonce()
	}
	return nil, &invalidResponseError{resp}
}
