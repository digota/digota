package customfields

import (
	"encoding/xml"
	"io"
	"strings"
)

// CustomFields is a string to string map of custom field names to values.
// Ref: https://articles.braintreepayments.com/control-panel/custom-fields
type CustomFields map[string]string

type xmlField struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

var nameMarshalReplacer = strings.NewReplacer("_", "-")
var nameUnmarshalReplacer = strings.NewReplacer("-", "_")

// MarshalXML encodes the map of custom fields names to values with the name as
// an xml tag and the value as it's contents. Tag names have underscores
// replaced with hyphens, as required by the Braintree API.
func (c CustomFields) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(c) == 0 {
		return nil
	}

	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	for k, v := range c {
		tag := nameMarshalReplacer.Replace(k)
		err := e.Encode(xmlField{
			XMLName: xml.Name{Local: tag},
			Value:   v,
		})
		if err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

// UnmarshalXML decodes the tags within into a map of custom fields names to
// values with the xml tag is the name and the contents is it's value. Tag
// names have hyphens replaced with underscores.
func (c *CustomFields) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*c = CustomFields{}
	for {
		var cf xmlField

		err := d.Decode(&cf)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		name := nameUnmarshalReplacer.Replace(cf.XMLName.Local)
		(*c)[name] = cf.Value
	}
	return nil
}
