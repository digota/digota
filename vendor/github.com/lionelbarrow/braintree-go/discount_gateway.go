package braintree

import "context"

type DiscountGateway struct {
	*Braintree
}

func (g *DiscountGateway) All(ctx context.Context) ([]Discount, error) {
	resp, err := g.execute(ctx, "GET", "discounts", nil)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.discounts()
	}
	return nil, &invalidResponseError{resp}
}
