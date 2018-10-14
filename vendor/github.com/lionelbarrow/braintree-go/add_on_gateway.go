package braintree

type AddOnGateway struct {
	*Braintree
}

func (g *AddOnGateway) All() ([]AddOn, error) {
	resp, err := g.execute("GET", "add_ons", nil)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.addOns()
	}
	return nil, &invalidResponseError{resp}
}
