package braintree

import (
	"context"
	"encoding/xml"
)

type PlanGateway struct {
	*Braintree
}

// All returns all available plans
func (g *PlanGateway) All(ctx context.Context) ([]*Plan, error) {
	resp, err := g.execute(ctx, "GET", "plans", nil)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		var b Plans
		if err := xml.Unmarshal(resp.Body, &b); err != nil {
			return nil, err
		}
		return b.Plan, nil
	}
	return nil, &invalidResponseError{resp}
}

// Find returns the plan with the specified id, or nil
func (g *PlanGateway) Find(ctx context.Context, id string) (*Plan, error) {
	plans, err := g.All(ctx)
	if err != nil {
		return nil, err
	}
	for _, p := range plans {
		if p.Id == id {
			return p, nil
		}
	}
	return nil, nil
}
