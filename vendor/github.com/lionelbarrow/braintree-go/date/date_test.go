package date

import (
	"encoding/xml"
	"testing"
	"time"
)

func TestDateUnmarshalXML(t *testing.T) {
	t.Parallel()

	date := &Date{}

	dateXML := []byte(`<?xml version="1.0" encoding="UTF-8"?><foo>2014-02-09</foo></xml>`)
	if err := xml.Unmarshal(dateXML, date); err != nil {
		t.Fatal(err)
	}

	if date.Format("2006-01-02") != "2014-02-09" {
		t.Fatalf("expected 2014-02-09 got %s", date)
	}
}

func TestDateMarshalXML(t *testing.T) {
	t.Parallel()

	date := &Date{Time: time.Date(2014, 2, 9, 0, 0, 0, 0, time.Local)}
	expected := `<Date>2014-02-09</Date>`

	b, err := xml.Marshal(date)
	if err != nil {
		t.Fatal(err)
	}

	if string(b) != expected {
		t.Fatalf("expected %s got %s", expected, string(b))
	}
}
