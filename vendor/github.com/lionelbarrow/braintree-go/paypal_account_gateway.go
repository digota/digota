package braintree

type PayPalAccountGateway struct {
	*Braintree
}

func (g *PayPalAccountGateway) Update(paypalAccount *PayPalAccount) (*PayPalAccount, error) {
	resp, err := g.executeVersion("PUT", "payment_methods/paypal_account/"+paypalAccount.Token, paypalAccount, apiVersion4)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.paypalAccount()
	}
	return nil, &invalidResponseError{resp}
}

func (g *PayPalAccountGateway) Find(token string) (*PayPalAccount, error) {
	resp, err := g.executeVersion("GET", "payment_methods/paypal_account/"+token, nil, apiVersion4)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.paypalAccount()
	}
	return nil, &invalidResponseError{resp}
}

func (g *PayPalAccountGateway) Delete(paypalAccount *PayPalAccount) error {
	resp, err := g.executeVersion("DELETE", "payment_methods/paypal_account/"+paypalAccount.Token, nil, apiVersion4)
	if err != nil {
		return err
	}
	switch resp.StatusCode {
	case 200:
		return nil
	}
	return &invalidResponseError{resp}
}
