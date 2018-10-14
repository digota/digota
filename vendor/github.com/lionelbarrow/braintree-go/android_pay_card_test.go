package braintree

import (
	"encoding/xml"
	"testing"
)

func TestAndroidPayCard_UnmarshalXML(t *testing.T) {
	x := `
	<android-pay-card>
		<source-card-type>mastercard</source-card-type>
		<source-card-last-4>1111</source-card-last-4>
		<virtual-card-type>visa</virtual-card-type>
		<virtual-card-last-4>2222</virtual-card-last-4>
	</android-pay-card>
	`
	var a AndroidPayCard

	if err := xml.Unmarshal([]byte(x), &a); err != nil {
		t.Fatalf("%v", err)
	}

	if g, w := a.SourceCardType, "mastercard"; g != w {
		t.Errorf("SourceCardType got %v, want %v", g, w)
	}
	if g, w := a.SourceCardLast4, "1111"; g != w {
		t.Errorf("SourceCardLast4 got %v, want %v", g, w)
	}
	if g, w := a.VirtualCardType, "visa"; g != w {
		t.Errorf("VirtualCardType got %v, want %v", g, w)
	}
	if g, w := a.VirtualCardLast4, "2222"; g != w {
		t.Errorf("VirtualCardLast4 got %v, want %v", g, w)
	}
	if g, w := a.CardType, "visa"; g != w {
		t.Errorf("CardType got %v, want %v", g, w)
	}
	if g, w := a.Last4, "2222"; g != w {
		t.Errorf("Last4 got %v, want %v", g, w)
	}
}
