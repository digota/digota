package braintree

import (
	"encoding/xml"
)

type PlanGateway struct {
	*Braintree
}

// All returns all available plans
func (g *PlanGateway) All() ([]*Plan, error) {
	resp, err := g.execute("GET", "plans", nil)
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
func (g *PlanGateway) Find(id string) (*Plan, error) {
	plans, err := g.All()
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
