package braintree

import "context"

type PayPalAccountGateway struct {
	*Braintree
}

func (g *PayPalAccountGateway) Update(ctx context.Context, paypalAccount *PayPalAccount) (*PayPalAccount, error) {
	resp, err := g.executeVersion(ctx, "PUT", "payment_methods/paypal_account/"+paypalAccount.Token, paypalAccount, apiVersion4)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.paypalAccount()
	}
	return nil, &invalidResponseError{resp}
}

func (g *PayPalAccountGateway) Find(ctx context.Context, token string) (*PayPalAccount, error) {
	resp, err := g.executeVersion(ctx, "GET", "payment_methods/paypal_account/"+token, nil, apiVersion4)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.paypalAccount()
	}
	return nil, &invalidResponseError{resp}
}

func (g *PayPalAccountGateway) Delete(ctx context.Context, paypalAccount *PayPalAccount) error {
	resp, err := g.executeVersion(ctx, "DELETE", "payment_methods/paypal_account/"+paypalAccount.Token, nil, apiVersion4)
	if err != nil {
		return err
	}
	switch resp.StatusCode {
	case 200:
		return nil
	}
	return &invalidResponseError{resp}
}
