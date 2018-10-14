package braintree

type DiscountGateway struct {
	*Braintree
}

func (g *DiscountGateway) All() ([]Discount, error) {
	resp, err := g.execute("GET", "discounts", nil)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.discounts()
	}
	return nil, &invalidResponseError{resp}
}
