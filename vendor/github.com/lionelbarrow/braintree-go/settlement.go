package braintree

type Record struct {
	XMLName           string   `xml:"record"`
	CardType          string   `xml:"card-type"`
	Count             int      `xml:"count"`
	MerchantAccountId string   `xml:"merchant-account-id"`
	Kind              string   `xml:"kind"`
	AmountSettled     *Decimal `xml:"amount-settled"`
}

type XMLRecords struct {
	XMLName string   `xml:"records"`
	Type    []Record `xml:"record"`
}
type SettlementBatchSummary struct {
	XMLName string     `xml:"settlement-batch-summary"`
	Records XMLRecords `xml:"records"`
}

type Settlement struct {
	XMLName string `xml:"settlement_batch_summary"`
	Date    string `xml:"settlement_date"`
}
