package braintree

import "context"

type AddOnGateway struct {
	*Braintree
}

// All gets all addons.
func (g *AddOnGateway) All(ctx context.Context) ([]AddOn, error) {
	resp, err := g.execute(ctx, "GET", "add_ons", nil)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.addOns()
	}
	return nil, &invalidResponseError{resp}
}
