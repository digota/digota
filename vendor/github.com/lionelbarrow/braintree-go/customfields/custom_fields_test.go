package customfields

import (
	"encoding/xml"
	"reflect"
	"testing"
)

func TestCustomFieldsMarshalXMLNil(t *testing.T) {
	type object struct {
		Field        string       `xml:"field,omitempty"`
		CustomFields CustomFields `xml:"custom-fields"`
	}
	o := object{}

	output, err := xml.Marshal(o)
	if err != nil {
		t.Fatalf("Error marshaling custom fields: %#v", err)
	}
	xml := string(output)

	if xml != "<object></object>" {
		t.Fatalf("Got %#v, wanted custom fields ommited", xml)
	}
}

func TestCustomFieldsMarshalXMLEmpty(t *testing.T) {
	type object struct {
		Field        string       `xml:"field,omitempty"`
		CustomFields CustomFields `xml:"custom-fields"`
	}
	o := object{CustomFields: CustomFields{}}

	output, err := xml.Marshal(o)
	if err != nil {
		t.Fatalf("Error marshaling custom fields: %#v", err)
	}
	xml := string(output)

	if xml != "<object></object>" {
		t.Fatalf("Got %#v, wanted custom fields ommited", xml)
	}
}

func TestCustomFieldsMarshalXML(t *testing.T) {
	type object struct {
		Field        string       `xml:"field"`
		CustomFields CustomFields `xml:"custom-fields"`
	}
	o := object{
		Field: "1.00",
		CustomFields: CustomFields{
			"custom_field_1": "Custom Value",
		},
	}
	wantedXML := `<object>
  <field>1.00</field>
  <custom-fields>
    <custom-field-1>Custom Value</custom-field-1>
  </custom-fields>
</object>`

	output, err := xml.MarshalIndent(o, "", "  ")
	if err != nil {
		t.Fatalf("Error marshaling custom fields: %#v", err)
	}
	xml := string(output)

	if xml != wantedXML {
		t.Fatalf("Got XML %#v, wanted %#v", xml, wantedXML)
	}
}

func TestCustomFieldsUnmarshalXMLNilEmpty(t *testing.T) {
	type object struct {
		Field        string       `xml:"field"`
		CustomFields CustomFields `xml:"custom-fields"`
	}

	s := `<object>
  <field>1.00</field>
</object>`
	wantedObject := object{
		Field:        "1.00",
		CustomFields: nil,
	}

	o := object{}
	err := xml.Unmarshal([]byte(s), &o)
	if err != nil {
		t.Fatalf("Error marshaling: %#v", err)
	}

	if !reflect.DeepEqual(o, wantedObject) {
		t.Fatalf("Got %#v, wanted %#v", o, wantedObject)
	}
}

func TestCustomFieldsUnmarshalXMLEmpty(t *testing.T) {
	type object struct {
		Field        string       `xml:"field"`
		CustomFields CustomFields `xml:"custom-fields"`
	}

	s := `<object>
  <field>1.00</field>
</object>`
	wantedObject := object{
		Field:        "1.00",
		CustomFields: nil,
	}

	o := object{}
	err := xml.Unmarshal([]byte(s), &o)
	if err != nil {
		t.Fatalf("Error marshaling: %#v", err)
	}

	if !reflect.DeepEqual(o, wantedObject) {
		t.Fatalf("Got %#v, wanted %#v", o, wantedObject)
	}
}

func TestCustomFieldsUnmarshalXMLNil(t *testing.T) {
	type object struct {
		Field        string       `xml:"field"`
		CustomFields CustomFields `xml:"custom-fields"`
	}

	s := `<object>
  <field>1.00</field>
  <custom-fields>
    <custom-field-1>Custom Value One</custom-field-1>
    <custom-field-2>Custom Value Two</custom-field-2>
  </custom-fields>
</object>`
	wantedObject := object{
		Field: "1.00",
		CustomFields: CustomFields{
			"custom_field_1": "Custom Value One",
			"custom_field_2": "Custom Value Two",
		},
	}

	o := object{}
	err := xml.Unmarshal([]byte(s), &o)
	if err != nil {
		t.Fatalf("Error marshaling: %#v", err)
	}

	if !reflect.DeepEqual(o, wantedObject) {
		t.Fatalf("Got %#v, wanted %#v", o, wantedObject)
	}
}

func TestCustomFieldsUnmarshalXML(t *testing.T) {
	type object struct {
		Field        string       `xml:"field"`
		CustomFields CustomFields `xml:"custom-fields"`
	}

	s := `<object>
  <field>1.00</field>
  <custom-fields>
    <custom-field-1>Custom Value One</custom-field-1>
    <custom-field-2>Custom Value Two</custom-field-2>
  </custom-fields>
</object>`
	wantedObject := object{
		Field: "1.00",
		CustomFields: CustomFields{
			"custom_field_1": "Custom Value One",
			"custom_field_2": "Custom Value Two",
		},
	}

	o := object{CustomFields: CustomFields{}}
	err := xml.Unmarshal([]byte(s), &o)
	if err != nil {
		t.Fatalf("Error marshaling: %#v", err)
	}

	if !reflect.DeepEqual(o, wantedObject) {
		t.Fatalf("Got %#v, wanted %#v", o, wantedObject)
	}
}
