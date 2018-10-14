package braintree

import (
	"context"
	"encoding/xml"
	"fmt"
)

type TransactionGateway struct {
	*Braintree
}

// Create initiates a transaction.
func (g *TransactionGateway) Create(ctx context.Context, tx *TransactionRequest) (*Transaction, error) {
	resp, err := g.execute(ctx, "POST", "transactions", tx)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 201:
		return resp.transaction()
	}
	return nil, &invalidResponseError{resp}
}

// Clone clones a transaction.
func (g *TransactionGateway) Clone(ctx context.Context, id string, tx *TransactionCloneRequest) (*Transaction, error) {
	resp, err := g.execute(ctx, "POST", "transactions/"+id+"/clone", tx)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 201:
		return resp.transaction()
	}
	return nil, &invalidResponseError{resp}
}

// SubmitForSettlement submits the transaction with the specified id for settlement.
// If the amount is omitted, the full amount is settled.
func (g *TransactionGateway) SubmitForSettlement(ctx context.Context, id string, amount ...*Decimal) (*Transaction, error) {
	var tx *TransactionRequest
	if len(amount) > 0 {
		tx = &TransactionRequest{
			Amount: amount[0],
		}
	}
	resp, err := g.execute(ctx, "PUT", "transactions/"+id+"/submit_for_settlement", tx)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.transaction()
	}
	return nil, &invalidResponseError{resp}
}

// Settle settles a transaction.
// This action is only available in the sandbox environment.
// Deprecated: use the Settle function on the TestingGateway instead. e.g. g.Testing().Settle(id).
func (g *TransactionGateway) Settle(ctx context.Context, id string) (*Transaction, error) {
	return g.Testing().Settle(ctx, id)
}

// Void voids the transaction with the specified id if it has a status of authorized or
// submitted_for_settlement. When the transaction is voided Braintree will do an authorization
// reversal if possible so that the customer won't have a pending charge on their card
func (g *TransactionGateway) Void(ctx context.Context, id string) (*Transaction, error) {
	resp, err := g.execute(ctx, "PUT", "transactions/"+id+"/void", nil)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.transaction()
	}
	return nil, &invalidResponseError{resp}
}

// CancelRelease cancels a pending release of a transaction with the given id from escrow.
func (g *TransactionGateway) CancelRelease(ctx context.Context, id string) (*Transaction, error) {
	resp, err := g.execute(ctx, "PUT", "transactions/"+id+"/cancel_release", nil)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.transaction()
	}
	return nil, &invalidResponseError{resp}
}

// ReleaseFromEscrow submits the transaction with the given id for release from escrow.
func (g *TransactionGateway) ReleaseFromEscrow(ctx context.Context, id string) (*Transaction, error) {
	resp, err := g.execute(ctx, "PUT", "transactions/"+id+"/release_from_escrow", nil)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.transaction()
	}
	return nil, &invalidResponseError{resp}
}

// HoldInEscrow holds the transaction with the given id for escrow.
func (g *TransactionGateway) HoldInEscrow(ctx context.Context, id string) (*Transaction, error) {
	resp, err := g.execute(ctx, "PUT", "transactions/"+id+"/hold_in_escrow", nil)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.transaction()
	}
	return nil, &invalidResponseError{resp}
}

// A transaction can be refunded if it is settled or settling.
// If the transaction has not yet begun settlement, use Void() instead.
// If you do not specify an amount to refund, the entire transaction amount will be refunded.
func (g *TransactionGateway) Refund(ctx context.Context, id string, amount ...*Decimal) (*Transaction, error) {
	var tx *TransactionRequest
	if len(amount) > 0 {
		tx = &TransactionRequest{
			Amount: amount[0],
		}
	}
	resp, err := g.execute(ctx, "POST", "transactions/"+id+"/refund", tx)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.transaction()
	case 201:
		return resp.transaction()
	}
	return nil, &invalidResponseError{resp}
}

// Find finds the transaction with the specified id.
func (g *TransactionGateway) Find(ctx context.Context, id string) (*Transaction, error) {
	resp, err := g.execute(ctx, "GET", "transactions/"+id, nil)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.transaction()
	}
	return nil, &invalidResponseError{resp}
}

// SearchIDs finds transactions matching the search query, returning the IDs only.
// Use SearchPage to get pages of transactions.
func (g *TransactionGateway) SearchIDs(ctx context.Context, query *SearchQuery) (*SearchResult, error) {
	resp, err := g.execute(ctx, "POST", "transactions/advanced_search_ids", query)
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

// SearchPage gets the page of transactions matching the search
// query.
// Use SearchIDs to start a search and get a list of IDs, use its
// result object to get pages.
// Page numbers start at 1.
// Returns a nil result and nil error when no more results are available.
func (g *TransactionGateway) SearchPage(ctx context.Context, query *SearchQuery, searchResult *SearchResult, page int) (*TransactionSearchResult, error) {
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
	transactions, err := g.fetchTransactions(ctx, pageQuery)

	pageResult := &TransactionSearchResult{
		TotalItems:        len(searchResult.IDs),
		TotalIDs:          searchResult.IDs,
		CurrentPageNumber: page,
		PageSize:          searchResult.PageSize,
		Transactions:      transactions,
	}

	return pageResult, err
}

// Search finds transactions matching the search query, returning the first
// page of results. Use SearchNext to get subsequent pages.
//
// Deprecated: Use SearchIDs and SearchPage.
func (g *TransactionGateway) Search(ctx context.Context, query *SearchQuery) (*TransactionSearchResult, error) {
	searchResult, err := g.SearchIDs(ctx, query)
	if err != nil {
		return nil, err
	}

	pageSize := searchResult.PageSize
	ids := searchResult.IDs

	endOffset := pageSize
	if endOffset > len(ids) {
		endOffset = len(ids)
	}

	firstPageQuery := query.shallowCopy()
	firstPageQuery.AddMultiField("ids").Items = ids[:endOffset]
	firstPageTransactions, err := g.fetchTransactions(ctx, firstPageQuery)

	firstPageResult := &TransactionSearchResult{
		TotalItems:        len(ids),
		TotalIDs:          ids,
		CurrentPageNumber: 1,
		PageSize:          pageSize,
		Transactions:      firstPageTransactions,
	}

	return firstPageResult, err
}

// SearchNext finds the next page of transactions matching the search query.
// Use Search to start a search and get the first page of results.
// Returns a nil result and nil error when no more results are available.
//
// Deprecated: Use SearchIDs and SearchPage.
func (g *TransactionGateway) SearchNext(ctx context.Context, query *SearchQuery, prevResult *TransactionSearchResult) (*TransactionSearchResult, error) {
	startOffset := prevResult.CurrentPageNumber * prevResult.PageSize
	endOffset := startOffset + prevResult.PageSize
	if endOffset > len(prevResult.TotalIDs) {
		endOffset = len(prevResult.TotalIDs)
	}
	if startOffset >= endOffset {
		return nil, nil
	}

	nextPageQuery := query.shallowCopy()
	nextPageQuery.AddMultiField("ids").Items = prevResult.TotalIDs[startOffset:endOffset]
	nextPageTransactions, err := g.fetchTransactions(ctx, nextPageQuery)

	nextPageResult := &TransactionSearchResult{
		TotalItems:        prevResult.TotalItems,
		TotalIDs:          prevResult.TotalIDs,
		CurrentPageNumber: prevResult.CurrentPageNumber + 1,
		PageSize:          prevResult.PageSize,
		Transactions:      nextPageTransactions,
	}

	return nextPageResult, err
}

func (g *TransactionGateway) fetchTransactions(ctx context.Context, query *SearchQuery) ([]*Transaction, error) {
	resp, err := g.execute(ctx, "POST", "transactions/advanced_search", query)
	if err != nil {
		return nil, err
	}
	var v struct {
		XMLName      string         `xml:"credit-card-transactions"`
		Transactions []*Transaction `xml:"transaction"`
	}
	err = xml.Unmarshal(resp.Body, &v)
	if err != nil {
		return nil, err
	}
	return v.Transactions, err
}

type testOperationPerformedInProductionError struct {
	error
}

func (e *testOperationPerformedInProductionError) Error() string {
	return fmt.Sprint("Operation not allowed in production environment")
}
