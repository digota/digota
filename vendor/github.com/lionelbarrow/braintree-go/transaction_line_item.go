package braintree

import "encoding/xml"

type TransactionLineItemKind string

const (
	TransactionLineItemKindDebit  TransactionLineItemKind = "debit"
	TransactionLineItemKindCredit TransactionLineItemKind = "credit"
)

type TransactionLineItem struct {
	Quantity       *Decimal                `xml:"quantity"`
	Name           string                  `xml:"name"`
	Description    string                  `xml:"description"`
	Kind           TransactionLineItemKind `xml:"kind"`
	UnitAmount     *Decimal                `xml:"unit-amount"`
	UnitTaxAmount  *Decimal                `xml:"unit-tax-amount"`
	TotalAmount    *Decimal                `xml:"total-amount"`
	TaxAmount      *Decimal                `xml:"tax-amount"`
	DiscountAmount *Decimal                `xml:"discount-amount"`
	UnitOfMeasure  string                  `xml:"unit-of-measure"`
	ProductCode    string                  `xml:"product-code"`
	CommodityCode  string                  `xml:"commodity-code"`
	URL            string                  `xml:"url"`
}

type TransactionLineItemRequest struct {
	Name           string                  `xml:"name"`
	Description    string                  `xml:"description,omitempty"`
	Kind           TransactionLineItemKind `xml:"kind"`
	Quantity       *Decimal                `xml:"quantity"`
	UnitAmount     *Decimal                `xml:"unit-amount"`
	UnitTaxAmount  *Decimal                `xml:"unit-tax-amount,omitempty"`
	TotalAmount    *Decimal                `xml:"total-amount"`
	TaxAmount      *Decimal                `xml:"tax-amount,omitempty"`
	DiscountAmount *Decimal                `xml:"discount-amount,omitempty"`
	UnitOfMeasure  string                  `xml:"unit-of-measure,omitempty"`
	ProductCode    string                  `xml:"product-code,omitempty"`
	CommodityCode  string                  `xml:"commodity-code,omitempty"`
	URL            string                  `xml:"url,omitempty"`
}

type TransactionLineItems []*TransactionLineItem

func (ls *TransactionLineItems) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	s := struct {
		LineItems []*TransactionLineItem `xml:"line-item"`
	}{}

	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	*ls = TransactionLineItems(s.LineItems)

	return nil
}

type TransactionLineItemRequests []*TransactionLineItemRequest

func (rs TransactionLineItemRequests) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(rs) == 0 {
		return nil
	}

	s := struct {
		XMLName   string                        `xml:"line-items"`
		Type      string                        `xml:"type,attr"`
		LineItems []*TransactionLineItemRequest `xml:"item"`
	}{
		Type:      "array",
		LineItems: []*TransactionLineItemRequest(rs),
	}

	return e.EncodeElement(s, start)
}
