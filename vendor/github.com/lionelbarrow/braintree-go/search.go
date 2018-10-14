package braintree

import (
	"encoding/xml"
	"time"
)

type SearchQuery struct {
	fields     []interface{}
	fieldIndex map[string]int
}

func (s *SearchQuery) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "search"
	x := struct {
		Fields []interface{}
	}{
		Fields: s.fields,
	}
	return e.EncodeElement(&x, start)
}

type SearchResult struct {
	PageSize  int
	PageCount int
	IDs       []string
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

func (s *SearchQuery) addField(fieldName string, field interface{}) {
	if i, ok := s.fieldIndex[fieldName]; !ok {
		s.fields = append(s.fields, field)
		if s.fieldIndex == nil {
			s.fieldIndex = map[string]int{}
		}
		s.fieldIndex[fieldName] = len(s.fields) - 1
	} else {
		s.fields[i] = field
	}
}

func (s *SearchQuery) AddTextField(field string) *TextField {
	f := &TextField{XMLName: xml.Name{Local: field}}
	s.addField(field, f)
	return f
}

func (s *SearchQuery) AddRangeField(field string) *RangeField {
	f := &RangeField{XMLName: xml.Name{Local: field}}
	s.addField(field, f)
	return f
}

func (s *SearchQuery) AddTimeField(field string) *TimeField {
	f := &TimeField{XMLName: xml.Name{Local: field}}
	s.addField(field, f)
	return f
}

func (s *SearchQuery) AddMultiField(field string) *MultiField {
	f := &MultiField{
		XMLName: xml.Name{Local: field},
		Type:    "array",
	}
	s.addField(field, f)
	return f
}

func (s *SearchQuery) shallowCopy() *SearchQuery {
	return &SearchQuery{
		fields: func() []interface{} {
			a := make([]interface{}, len(s.fields))
			copy(a, s.fields)
			return a
		}(),
		fieldIndex: func() map[string]int {
			m := map[string]int{}
			for f, i := range s.fieldIndex {
				m[f] = i
			}
			return m
		}(),
	}
}
