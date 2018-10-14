package braintree

// TestingGateway exports actions only available in the sandbox environment.
type TestingGateway struct {
	*Braintree
}

// Settle changes the transaction's status to settle.
func (g *TestingGateway) Settle(transactionID string) (*Transaction, error) {
	return g.setStatus(transactionID, "settle")
}

// SettlementConfirm changes the transaction's status to settlement_confirm.
func (g *TestingGateway) SettlementConfirm(transactionID string) (*Transaction, error) {
	return g.setStatus(transactionID, "settlement_confirm")
}

// SettlementDecline changes the transaction's status to settlement_decline.
func (g *TestingGateway) SettlementDecline(transactionID string) (*Transaction, error) {
	return g.setStatus(transactionID, "settlement_decline")
}

// SettlementPending changes the transaction's status to settlement_pending.
func (g *TestingGateway) SettlementPending(transactionID string) (*Transaction, error) {
	return g.setStatus(transactionID, TransactionStatusSettlementPending)
}

func (g *TestingGateway) setStatus(transactionID string, status TransactionStatus) (*Transaction, error) {
	if g.Environment() != Production {
		resp, err := g.execute("PUT", "transactions/"+transactionID+"/"+string(status), nil)
		if err != nil {
			return nil, err
		}
		switch resp.StatusCode {
		case 200:
			return resp.transaction()
		}
		return nil, &invalidResponseError{resp}
	} else {
		return nil, &testOperationPerformedInProductionError{}
	}
}
