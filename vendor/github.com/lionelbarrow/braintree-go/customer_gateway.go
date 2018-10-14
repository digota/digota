package braintree

import (
	"context"
	"encoding/xml"
	"fmt"
)

type CustomerGateway struct {
	*Braintree
}

// Create creates a new customer from the passed in customer object.
// If no ID is set, Braintree will assign one.
func (g *CustomerGateway) Create(ctx context.Context, c *CustomerRequest) (*Customer, error) {
	resp, err := g.execute(ctx, "POST", "customers", c)
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
// The ID field is mandatory.
func (g *CustomerGateway) Update(ctx context.Context, c *CustomerRequest) (*Customer, error) {
	resp, err := g.execute(ctx, "PUT", "customers/"+c.ID, c)
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
func (g *CustomerGateway) Find(ctx context.Context, id string) (*Customer, error) {
	resp, err := g.execute(ctx, "GET", "customers/"+id, nil)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.customer()
	}
	return nil, &invalidResponseError{resp}
}

// SearchIDs finds customers matching the search query, returning the IDs
// only. Use SearchPage to get pages of customers.
func (g *CustomerGateway) SearchIDs(ctx context.Context, query *SearchQuery) (*SearchResult, error) {
	resp, err := g.execute(ctx, "POST", "customers/advanced_search_ids", query)
	if err != nil {
		return nil, err
	}

	var searchResult struct {
		PageSize int `xml:"page-size"`
		Ids      struct {
			Item []string `xml:"item"`
		} `xml:"ids"`
	}
	err = xml.Unmarshal(resp.Body, &searchResult)
	if err != nil {
		return nil, err
	}

	return &SearchResult{
		PageSize:  searchResult.PageSize,
		PageCount: (len(searchResult.Ids.Item) + searchResult.PageSize - 1) / searchResult.PageSize,
		IDs:       searchResult.Ids.Item,
	}, nil
}

// SearchPage gets the page of customers matching the search
// query.
// Use SearchIDs to start a search and get a list of IDs, use its
// result object to get pages.
// Page numbers start at 1.
// Returns a nil result and nil error when no more results are available.
func (g *CustomerGateway) SearchPage(ctx context.Context, query *SearchQuery, searchResult *SearchResult, page int) (*CustomerSearchResult, error) {
	if page < 1 || page > searchResult.PageCount {
		return nil, fmt.Errorf("page %d out of bounds, page numbers start at 1 and page count is %d", page, searchResult.PageCount)
	}
	startOffset := (page - 1) * searchResult.PageSize
	endOffset := startOffset + searchResult.PageSize
	if endOffset > len(searchResult.IDs) {
		endOffset = len(searchResult.IDs)
	}

	pageQuery := query.shallowCopy()
	pageQuery.AddMultiField("ids").Items = searchResult.IDs[startOffset:endOffset]
	customers, err := g.fetchCustomers(ctx, pageQuery)

	pageResult := &CustomerSearchResult{
		TotalItems:        len(searchResult.IDs),
		TotalIDs:          searchResult.IDs,
		CurrentPageNumber: page,
		PageSize:          searchResult.PageSize,
		Customers:         customers,
	}

	return pageResult, err
}

// Search finds customers matching the search query, returning the first page
// of results.
//
// Deprecated: Use SearchIDs and SearchPage.
func (g *CustomerGateway) Search(ctx context.Context, query *SearchQuery) (*CustomerSearchResult, error) {
	searchResult, err := g.SearchIDs(ctx, query)
	if err != nil {
		return nil, err
	}

	return g.SearchPage(ctx, query, searchResult, 1)
}

func (g *CustomerGateway) fetchCustomers(ctx context.Context, query *SearchQuery) ([]*Customer, error) {
	resp, err := g.execute(ctx, "POST", "customers/advanced_search", query)
	if err != nil {
		return nil, err
	}
	var v struct {
		XMLName   string      `xml:"customers"`
		Customers []*Customer `xml:"customer"`
	}
	err = xml.Unmarshal(resp.Body, &v)
	if err != nil {
		return nil, err
	}
	return v.Customers, err
}

// Delete deletes the customer with the given id.
func (g *CustomerGateway) Delete(ctx context.Context, id string) error {
	resp, err := g.execute(ctx, "DELETE", "customers/"+id, nil)
	if err != nil {
		return err
	}
	switch resp.StatusCode {
	case 200:
		return nil
	}
	return &invalidResponseError{resp}
}
