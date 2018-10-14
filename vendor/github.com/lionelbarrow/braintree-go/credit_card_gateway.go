package braintree

type CreditCardGateway struct {
	*Braintree
}

func (g *CreditCardGateway) Create(card *CreditCard) (*CreditCard, error) {
	resp, err := g.execute("POST", "payment_methods", card)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 201:
		return resp.creditCard()
	}
	return nil, &invalidResponseError{resp}
}

func (g *CreditCardGateway) Update(card *CreditCard) (*CreditCard, error) {
	resp, err := g.execute("PUT", "payment_methods/"+card.Token, card)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.creditCard()
	}
	return nil, &invalidResponseError{resp}
}

func (g *CreditCardGateway) Find(token string) (*CreditCard, error) {
	resp, err := g.execute("GET", "payment_methods/"+token, nil)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.creditCard()
	}
	return nil, &invalidResponseError{resp}
}

func (g *CreditCardGateway) Delete(card *CreditCard) error {
	resp, err := g.execute("DELETE", "payment_methods/"+card.Token, nil)
	if err != nil {
		return err
	}
	switch resp.StatusCode {
	case 200:
		return nil
	}
	return &invalidResponseError{resp}
}
