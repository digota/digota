package date

import (
	"encoding/xml"
	"time"
)

// Date wraps a time object but handles deserializing dates returned from the Braintree API
// e.g. "2014-02-09"
type Date struct {
	time.Time
}

// UnmarshalXML interprets Braintree's date format from XML to initialize the Date
// e.g. "2014-02-09"
func (d *Date) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var v string
	err := dec.DecodeElement(&v, &start)
	if err != nil {
		return err
	}

	parse, err := time.Parse("2006-01-02", v)
	if err != nil {
		return err
	}

	*d = Date{Time: parse}
	return nil
}

// MarshalXML Outputs the date in a format the Braintree expects to see date
// e.g. "2014-02-09"
func (d *Date) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(d.Format("2006-01-02"), start)
}
