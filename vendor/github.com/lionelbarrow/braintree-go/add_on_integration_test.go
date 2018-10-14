package braintree

import (
	"testing"
)

func TestAddOn(t *testing.T) {
	t.Parallel()

	addOns, err := testGateway.AddOn().All()

	if err != nil {
		t.Fatal(err)
	} else if len(addOns) != 1 {
		t.Fatalf("expected to retrieve one add-on, but retrieved %d", len(addOns))
	}

	addOn := addOns[0]

	t.Log(addOn)

	if addOn.Id != "test_add_on" {
		t.Fatalf("expected Id to be %s, was %s", "test_add_on", addOn.Id)
	} else if addOn.Amount.Cmp(NewDecimal(1000, 2)) != 0 {
		t.Fatalf("expected Amount to be %s, was %s", NewDecimal(1000, 2), addOn.Amount)
	} else if addOn.Kind != ModificationKindAddOn {
		t.Fatalf("expected Kind to be %s, was %s", ModificationKindAddOn, addOn.Kind)
	} else if addOn.Name != "test_add_on_name" {
		t.Fatalf("expected Name to be %s, was %s", "test_add_on_name", addOn.Name)
	} else if addOn.NeverExpires != true {
		t.Fatalf("expected NeverExpires to be %v, was %v", true, addOn.NeverExpires)
	} else if addOn.Description != "A test add-on" {
		t.Fatalf("expected Description to be %s, was %s", "A test add-on", addOn.Description)
	}
}
