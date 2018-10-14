package braintree

import (
	"encoding/xml"
	"time"
)

type SearchQuery struct {
	XMLName string `xml:"search"`
	Fields  []interface{}
}

type SearchResults struct {
	XMLName  string `xml:"search-results"`
	PageSize string `xml:"page-size"`
	Ids      struct {
		Item []string `xml:"item"`
	} `xml:"ids"`
}

type TextField struct {
	XMLName    xml.Name
	Is         string `xml:"is,omitempty"`
	IsNot      string `xml:"is-not,omitempty"`
	StartsWith string `xml:"starts-with,omitempty"`
	EndsWith   string `xml:"ends-with,omitempty"`
	Contains   string `xml:"contains,omitempty"`
}

type RangeField struct {
	XMLName xml.Name
	Is      float64 `xml:"is,omitempty"`
	Min     float64 `xml:"min,omitempty"`
	Max     float64 `xml:"max,omitempty"`
}

type TimeField struct {
	XMLName xml.Name
	Is      time.Time
	Min     time.Time
	Max     time.Time
}

func (d TimeField) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start = start.Copy()
	start.Name = d.XMLName

	var err error
	err = e.EncodeToken(start)
	if err != nil {
		return err
	}

	err = d.marshalXMLCriterion(e, "is", d.Is)
	if err != nil {
		return err
	}

	err = d.marshalXMLCriterion(e, "min", d.Min)
	if err != nil {
		return err
	}

	err = d.marshalXMLCriterion(e, "max", d.Max)
	if err != nil {
		return err
	}

	err = e.EncodeToken(start.End())
	return err
}

func (d TimeField) marshalXMLCriterion(e *xml.Encoder, name string, value time.Time) error {
	if value.IsZero() {
		return nil
	}
	const format = "2006-01-02T15:04:05Z"
	start := xml.StartElement{Name: xml.Name{Local: name}}
	start.Attr = []xml.Attr{{Name: xml.Name{Local: "type"}, Value: "datetime"}}
	return e.EncodeElement(value.UTC().Format(format), start)
}

type MultiField struct {
	XMLName xml.Name
	Type    string   `xml:"type,attr"` // type=array
	Items   []string `xml:"item"`
}

func (s *SearchQuery) AddTextField(field string) *TextField {
	f := &TextField{XMLName: xml.Name{Local: field}}
	s.Fields = append(s.Fields, f)
	return f
}

func (s *SearchQuery) AddRangeField(field string) *RangeField {
	f := &RangeField{XMLName: xml.Name{Local: field}}
	s.Fields = append(s.Fields, f)
	return f
}

func (s *SearchQuery) AddTimeField(field string) *TimeField {
	f := &TimeField{XMLName: xml.Name{Local: field}}
	s.Fields = append(s.Fields, f)
	return f
}

func (s *SearchQuery) AddMultiField(field string) *MultiField {
	f := &MultiField{
		XMLName: xml.Name{Local: field},
		Type:    "array",
	}
	s.Fields = append(s.Fields, f)
	return f
}
