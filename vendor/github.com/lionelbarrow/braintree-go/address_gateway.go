package braintree

import (
	"encoding/xml"
)

type AddressGateway struct {
	*Braintree
}

// Create creates a new address for the specified customer id.
func (g *AddressGateway) Create(a *Address) (*Address, error) {
	// Copy address so that field sanitation won't affect original
	var cp Address = *a
	cp.CustomerId = ""
	cp.XMLName = xml.Name{Local: "address"}

	resp, err := g.execute("POST", "customers/"+a.CustomerId+"/addresses", &cp)
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
func (g *AddressGateway) Delete(customerId, addrId string) error {
	resp, err := g.execute("DELETE", "customers/"+customerId+"/addresses/"+addrId, nil)
	if err != nil {
		return err
	}
	switch resp.StatusCode {
	case 200:
		return nil
	}
	return &invalidResponseError{resp}
}
