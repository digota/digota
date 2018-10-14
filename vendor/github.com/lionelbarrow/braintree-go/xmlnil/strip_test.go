package xmlnil

import (
	"encoding/xml"
	"testing"
)

func TestStripNilElements(t *testing.T) {
	testCases := []struct {
		In  string
		Out string
	}{
		{``, ``},
		{`<element/>`, `<element></element>`},
		{`<element></element>`, `<element></element>`},
		{`<element nil="hi"/>`, `<element nil="hi"></element>`},
		{`<element nil="hi"></element>`, `<element nil="hi"></element>`},
		{`<element nil="true"/>`, ``},
		{`<element nil="true"></element>`, ``},
		{`<root><element><sub>name</sub></element></root>`, `<root><element><sub>name</sub></element></root>`},
		{`<root><element nil="true"><sub>name</sub></element></root>`, `<root></root>`},
		{`<root><element><sub nil="true">name</sub></element></root>`, `<root><element></element></root>`},
	}

	for _, tc := range testCases {
		out, err := StripNilElements([]byte(tc.In))
		if err != nil {
			t.Errorf("StripNilElements(%+v) got err %+v", tc.In, err)
		} else if string(out) != tc.Out {
			t.Errorf("StripNilElements(%+v) got %+v, want %+v", tc.In, string(out), tc.Out)
		} else {
			t.Logf("StripNilElements(%+v) got %q", tc.In, string(out))
		}
	}
}

func TestStripNilElementsErrors(t *testing.T) {
	testCases := []struct {
		In string
	}{
		{`<element>`},
		{`<element><element>`},
		{`<element></sub>`},
		{`<element nil="hi">`},
		{`<element nil="true/>`},
	}

	for _, tc := range testCases {
		_, err := StripNilElements([]byte(tc.In))
		if err == nil {
			t.Errorf("StripNilElements(%+v) got no err", tc.In)
		} else {
			t.Logf("StripNilElements(%+v) got err %+v", tc.In, err)
		}
	}
}

func TestIsTokenNil(t *testing.T) {
	testCases := []struct {
		Token xml.Token
		IsNil bool
	}{
		{xml.Comment{}, false},
		{xml.CharData{}, false},
		{xml.EndElement{}, false},
		{xml.StartElement{}, false},
		{xml.StartElement{Name: xml.Name{Local: "element"}}, false},
	}

	for _, tc := range testCases {
		isNil := isTokenNil(tc.Token)
		if isNil != tc.IsNil {
			t.Errorf("isTokenNil(%+v) got %+v, want %+v", tc.Token, isNil, tc.IsNil)
		}
	}
}
