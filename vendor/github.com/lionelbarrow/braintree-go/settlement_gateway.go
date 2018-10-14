package braintree

type SettlementGateway struct {
	*Braintree
}

func (sg *SettlementGateway) Generate(s *Settlement) (*SettlementBatchSummary, error) {
	resp, err := sg.execute("POST", "settlement_batch_summary", s)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.settlement()
	}
	return nil, &invalidResponseError{resp}
}
