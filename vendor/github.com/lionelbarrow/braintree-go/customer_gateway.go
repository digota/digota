package braintree

import "encoding/xml"

type CustomerGateway struct {
	*Braintree
}

// Create creates a new customer from the passed in customer object.
// If no Id is set, Braintree will assign one.
func (g *CustomerGateway) Create(c *Customer) (*Customer, error) {
	resp, err := g.execute("POST", "customers", c)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 201:
		return resp.customer()
	}
	return nil, &invalidResponseError{resp}
}

// Update updates any field that is set in the passed customer object.
// The Id field is mandatory.
func (g *CustomerGateway) Update(c *Customer) (*Customer, error) {
	resp, err := g.execute("PUT", "customers/"+c.Id, c)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.customer()
	}
	return nil, &invalidResponseError{resp}
}

// Find finds the customer with the given id.
func (g *CustomerGateway) Find(id string) (*Customer, error) {
	resp, err := g.execute("GET", "customers/"+id, nil)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.customer()
	}
	return nil, &invalidResponseError{resp}
}

func (g *CustomerGateway) Search(query *SearchQuery) (*CustomerSearchResult, error) {
	resp, err := g.execute("POST", "customers/advanced_search", query)
	if err != nil {
		return nil, err
	}
	var v CustomerSearchResult
	err = xml.Unmarshal(resp.Body, &v)
	if err != nil {
		return nil, err
	}
	return &v, err
}

// Delete deletes the customer with the given id.
func (g *CustomerGateway) Delete(id string) error {
	resp, err := g.execute("DELETE", "customers/"+id, nil)
	if err != nil {
		return err
	}
	switch resp.StatusCode {
	case 200:
		return nil
	}
	return &invalidResponseError{resp}
}
