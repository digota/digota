package braintree

import (
	"context"
)

type AddressGateway struct {
	*Braintree
}

// Create creates a new address for the specified customer id.
func (g *AddressGateway) Create(ctx context.Context, customerID string, a *AddressRequest) (*Address, error) {
	resp, err := g.execute(ctx, "POST", "customers/"+customerID+"/addresses", &a)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 201:
		return resp.address()
	}
	return nil, &invalidResponseError{resp}
}

// Delete deletes the address for the specified id and customer id.
func (g *AddressGateway) Delete(ctx context.Context, customerId, addrId string) error {
	resp, err := g.execute(ctx, "DELETE", "customers/"+customerId+"/addresses/"+addrId, nil)
	if err != nil {
		return err
	}
	switch resp.StatusCode {
	case 200:
		return nil
	}
	return &invalidResponseError{resp}
}

// Update updates an address for the address id and customer id.
func (g *AddressGateway) Update(ctx context.Context, customerID, addrID string, a *AddressRequest) (*Address, error) {
	resp, err := g.execute(ctx, "PUT", "customers/"+customerID+"/addresses/"+addrID, a)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.address()
	}
	return nil, &invalidResponseError{resp}
}
