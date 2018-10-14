package braintree

import "encoding/xml"

type AndroidPayDetails struct {
	Token               string `xml:"token,omitempty"`
	CardType            string `xml:"-"`
	Last4               string `xml:"-"`
	SourceCardType      string `xml:"source-card-type,omitempty"`
	SourceCardLast4     string `xml:"source-card-last-4,omitempty"`
	SourceDescription   string `xml:"source-description,omitempty"`
	VirtualCardType     string `xml:"virtual-card-type,omitempty"`
	VirtualCardLast4    string `xml:"virtual-card-last-4,omitempty"`
	ExpirationMonth     string `xml:"expiration-month,omitempty"`
	ExpirationYear      string `xml:"expiration-year,omitempty"`
	BIN                 string `xml:"bin,omitempty"`
	GoogleTransactionID string `xml:"google-transaction-id,omitempty"`
	ImageURL            string `xml:"image-url,omitempty"`
}

func (a *AndroidPayDetails) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type typeWithNoFunctions AndroidPayDetails
	if err := d.DecodeElement((*typeWithNoFunctions)(a), &start); err != nil {
		return err
	}
	a.CardType = a.VirtualCardType
	a.Last4 = a.VirtualCardLast4
	return nil
}
