package braintree

import "time"

type DisputeKind string

const (
	DisputeKindChargeback     DisputeKind = "chargeback"
	DisputeKindPreArbitration DisputeKind = "pre_arbitration"
	DisputeKindRetrieval      DisputeKind = "retrieval"
)

type DisputeReason string

const (
	DisputeReasonCancelledRecurringTransaction DisputeReason = "cancelled_recurring_transaction"
	DisputeReasonCreditNotProcessed            DisputeReason = "credit_not_processed"
	DisputeReasonDuplicate                     DisputeReason = "duplicate"
	DisputeReasonFraud                         DisputeReason = "fraud"
	DisputeReasonGeneral                       DisputeReason = "general"
	DisputeReasonInvalidAccount                DisputeReason = "invalid_account"
	DisputeReasonNotRecognized                 DisputeReason = "not_recognized"
	DisputeReasonProductNotReceived            DisputeReason = "product_not_received"
	DisputeReasonProductUnsatisfactory         DisputeReason = "product_unsatisfactory"
	DisputeReasonTransactionAmountDiffers      DisputeReason = "transaction_amount_differs"
)

type DisputeStatus string

const (
	DisputeStatusAccepted DisputeStatus = "accepted"
	DisputeStatusDisputed DisputeStatus = "disputed"
	DisputeStatusExpired  DisputeStatus = "expired"
	DisputeStatusOpen     DisputeStatus = "open"
	DisputeStatusLost     DisputeStatus = "lost"
	DisputeStatusWon      DisputeStatus = "won"
)

type Dispute struct {
	XMLName           string                       `xml:"dispute"`
	AmountDisputed    *Decimal                     `xml:"amount-disputed"`
	AmountWon         *Decimal                     `xml:"amount-won"`
	CaseNumber        string                       `xml:"case-number"`
	CreatedAt         *time.Time                   `xml:"created-at"`
	CurrencyISOCode   string                       `xml:"currency-iso-code"`
	Evidence          []*DisputeEvidence           `xml:"evidence>evidence"`
	ID                string                       `xml:"id"`
	Kind              DisputeKind                  `xml:"kind"`
	MerchantAccountID string                       `xml:"merchant-account-id"`
	OriginalDisputeID string                       `xml:"original-dispute-id"`
	ProcessorComments string                       `xml:"processor-comments"`
	Reason            DisputeReason                `xml:"reason"`
	ReturnCode        string                       `xml:"return-code"`
	ReceivedDate      string                       `xml:"received-date"`
	ReferenceNumber   string                       `xml:"reference-number"`
	ReplyByDate       string                       `xml:"reply-by-date"`
	Status            DisputeStatus                `xml:"status"`
	StatusHistory     []*DisputeStatusHistoryEvent `xml:"status-history>status-history"`
	Transaction       *DisputeTransaction          `xml:"transaction"`
	UpdatedAt         *time.Time                   `xml:"updated-at"`
}

type DisputeStatusHistoryEvent struct {
	XMLName          string     `xml:"status-history"`
	DisbursementDate string     `xml:"disbursement-date"`
	EffectiveDate    string     `xml:"effective-date"`
	Status           string     `xml:"status"`
	Timestamp        *time.Time `xml:"timestamp"`
}

type DisputeTransaction struct {
	XMLName                  string     `xml:"transaction"`
	Amount                   *Decimal   `xml:"amount"`
	CreatedAt                *time.Time `xml:"created-at"`
	ID                       string     `xml:"id"`
	OrderID                  string     `xml:"order-id"`
	PaymentInstrumentSubtype string     `xml:"payment-instrument-subtype"`
	PurchaseOrderNumber      string     `xml:"purchase-order-number"`
}
