package braintree

import (
	"encoding/xml"
)

type ModificationRequest struct {
	Amount                *Decimal `xml:"amount,omitempty"`
	NumberOfBillingCycles int      `xml:"number-of-billing-cycles,omitempty"`
	Quantity              int      `xml:"quantity,omitempty"`
	NeverExpires          bool     `xml:"never-expires,omitempty"`
}

type AddModificationRequest struct {
	ModificationRequest
	InheritedFromID string `xml:"inherited-from-id,omitempty"`
}

type UpdateModificationRequest struct {
	ModificationRequest
	ExistingID string `xml:"existing-id,omitempty"`
}

type ModificationsRequest struct {
	Add               []AddModificationRequest
	Update            []UpdateModificationRequest
	RemoveExistingIDs []string
}

func (m ModificationsRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type addSchema struct {
		Type          string                   `xml:"type,attr"`
		Modifications []AddModificationRequest `xml:"modification"`
	}
	type updateSchema struct {
		Type          string                      `xml:"type,attr"`
		Modifications []UpdateModificationRequest `xml:"modification"`
	}
	type removeSchema struct {
		Type        string   `xml:"type,attr"`
		ExistingIDs []string `xml:"modification"`
	}
	type schema struct {
		Add    *addSchema    `xml:"add,omitempty"`
		Update *updateSchema `xml:"update,omitempty"`
		Remove *removeSchema `xml:"remove,omitempty"`
	}

	x := schema{}
	if len(m.Add) > 0 {
		x.Add = &addSchema{
			Type:          "array",
			Modifications: m.Add,
		}
	}
	if len(m.Update) > 0 {
		x.Update = &updateSchema{
			Type:          "array",
			Modifications: m.Update,
		}
	}
	if len(m.RemoveExistingIDs) > 0 {
		x.Remove = &removeSchema{
			Type:        "array",
			ExistingIDs: m.RemoveExistingIDs,
		}
	}
	return e.EncodeElement(x, start)
}
