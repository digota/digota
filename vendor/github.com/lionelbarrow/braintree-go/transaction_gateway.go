package braintree

import (
	"encoding/xml"
	"fmt"
)

type TransactionGateway struct {
	*Braintree
}

// Create initiates a transaction.
func (g *TransactionGateway) Create(tx *TransactionRequest) (*Transaction, error) {
	resp, err := g.execute("POST", "transactions", tx)
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
func (g *TransactionGateway) Clone(id string, tx *TransactionCloneRequest) (*Transaction, error) {
	resp, err := g.execute("POST", "transactions/"+id+"/clone", tx)
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
func (g *TransactionGateway) SubmitForSettlement(id string, amount ...*Decimal) (*Transaction, error) {
	var tx *TransactionRequest
	if len(amount) > 0 {
		tx = &TransactionRequest{
			Amount: amount[0],
		}
	}
	resp, err := g.execute("PUT", "transactions/"+id+"/submit_for_settlement", tx)
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
func (g *TransactionGateway) Settle(id string) (*Transaction, error) {
	return g.Testing().Settle(id)
}

// Void voids the transaction with the specified id if it has a status of authorized or
// submitted_for_settlement. When the transaction is voided Braintree will do an authorization
// reversal if possible so that the customer won’t have a pending charge on their card
func (g *TransactionGateway) Void(id string) (*Transaction, error) {
	resp, err := g.execute("PUT", "transactions/"+id+"/void", nil)
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
func (g *TransactionGateway) CancelRelease(id string) (*Transaction, error) {
	resp, err := g.execute("PUT", "transactions/"+id+"/cancel_release", nil)
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
func (g *TransactionGateway) ReleaseFromEscrow(id string) (*Transaction, error) {
	resp, err := g.execute("PUT", "transactions/"+id+"/release_from_escrow", nil)
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
func (g *TransactionGateway) HoldInEscrow(id string) (*Transaction, error) {
	resp, err := g.execute("PUT", "transactions/"+id+"/hold_in_escrow", nil)
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
func (g *TransactionGateway) Refund(id string, amount ...*Decimal) (*Transaction, error) {
	var tx *TransactionRequest
	if len(amount) > 0 {
		tx = &TransactionRequest{
			Amount: amount[0],
		}
	}
	resp, err := g.execute("POST", "transactions/"+id+"/refund", tx)
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
func (g *TransactionGateway) Find(id string) (*Transaction, error) {
	resp, err := g.execute("GET", "transactions/"+id, nil)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.transaction()
	}
	return nil, &invalidResponseError{resp}
}

// Search finds all transactions matching the search query.
func (g *TransactionGateway) Search(query *SearchQuery) (*TransactionSearchResult, error) {
	resp, err := g.execute("POST", "transactions/advanced_search", query)
	if err != nil {
		return nil, err
	}
	var v TransactionSearchResult
	err = xml.Unmarshal(resp.Body, &v)
	if err != nil {
		return nil, err
	}
	return &v, err
}

type testOperationPerformedInProductionError struct {
	error
}

func (e *testOperationPerformedInProductionError) Error() string {
	return fmt.Sprint("Operation not allowed in production environment")
}
