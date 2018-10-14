package braintree

import "encoding/xml"

type AndroidPayDetails struct {
	Token               string `xml:"token"`
	CardType            string `xml:"-"`
	Last4               string `xml:"-"`
	SourceCardType      string `xml:"source-card-type"`
	SourceCardLast4     string `xml:"source-card-last-4"`
	SourceDescription   string `xml:"source-description"`
	VirtualCardType     string `xml:"virtual-card-type"`
	VirtualCardLast4    string `xml:"virtual-card-last-4"`
	ExpirationMonth     string `xml:"expiration-month"`
	ExpirationYear      string `xml:"expiration-year"`
	BIN                 string `xml:"bin"`
	GoogleTransactionID string `xml:"google-transaction-id"`
	ImageURL            string `xml:"image-url"`
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
