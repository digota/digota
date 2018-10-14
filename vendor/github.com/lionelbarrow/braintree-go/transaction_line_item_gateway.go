package braintree

import (
	"context"
)

type TransactionLineItemGateway struct {
	*Braintree
}

// Find finds the transaction line items with the specified transaction id.
func (g *TransactionLineItemGateway) Find(ctx context.Context, transactionID string) (TransactionLineItems, error) {
	resp, err := g.execute(ctx, "GET", "transactions/"+transactionID+"/line_items", nil)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.transactionLineItems()
	}
	return nil, &invalidResponseError{resp}
}
