package braintree

import "context"

type SettlementGateway struct {
	*Braintree
}

func (sg *SettlementGateway) Generate(ctx context.Context, s *Settlement) (*SettlementBatchSummary, error) {
	resp, err := sg.execute(ctx, "POST", "settlement_batch_summary", s)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.settlement()
	}
	return nil, &invalidResponseError{resp}
}
