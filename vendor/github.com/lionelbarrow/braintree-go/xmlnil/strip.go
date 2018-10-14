package xmlnil

import (
	"bytes"
	"encoding/xml"
	"io"
)

// StripNilElements parses the xml input, removing any elements that
// are decorated with the `nil="true"` attribute returning the XML
// without those elements.
func StripNilElements(x []byte) ([]byte, error) {
	d := xml.NewDecoder(bytes.NewReader(x))
	w := &bytes.Buffer{}
	e := xml.NewEncoder(w)

	for {
		t, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if isTokenNil(t) {
			err := d.Skip()
			if err != nil {
				return nil, err
			}
			continue
		}

		if err := e.EncodeToken(t); err != nil {
			return nil, err
		}
	}

	if err := e.Flush(); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func isTokenNil(t xml.Token) bool {
	if t, ok := t.(xml.StartElement); ok {
		for _, a := range t.Attr {
			if a.Name.Space == "" && a.Name.Local == "nil" && a.Value == "true" {
				return true
			}
		}
	}
	return false
}
