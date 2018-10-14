package braintree

import (
	"testing"
)

// This test will fail unless you set up your Braintree sandbox account correctly. See TESTING.md for details.
func TestPlan(t *testing.T) {
	t.Parallel()

	g := testGateway.Plan()
	plans, err := g.All()
	if err != nil {
		t.Fatal(err)
	}
	if len(plans) == 0 {
		t.Fatal(plans)
	}

	var plan *Plan
	for _, p := range plans {
		if p.Id == "test_plan" {
			plan = p
			break
		}
	}

	t.Log(plan)

	if plan == nil {
		t.Fatal("plan not found")
	}
	if x := plan.Id; x != "test_plan" {
		t.Fatal(x)
	}
	if x := plan.MerchantId; x == "" {
		t.Fatal(x)
	}
	if x := plan.BillingFrequency; x == nil || *x != 1 {
		t.Fatal(x)
	}
	if x := plan.CurrencyISOCode; x != "USD" {
		t.Fatal(x)
	}
	if x := plan.Description; x != "test_plan_desc" {
		t.Fatal(x)
	}
	if x := plan.Name; x != "test_plan_name" {
		t.Fatal(x)
	}
	if x := plan.NumberOfBillingCycles; x == nil || *x != 2 {
		t.Fatal(x)
	}
	if x := plan.Price; x.Cmp(NewDecimal(1000, 2)) != 0 {
		t.Fatal(x)
	}
	if x := plan.TrialDuration; x == nil || *x != 14 {
		t.Fatal(x)
	}
	if x := plan.TrialDurationUnit; x != "day" {
		t.Fatal(x)
	}
	if x := plan.TrialPeriod; !x {
		t.Fatal(x)
	}
	if x := plan.CreatedAt; x == nil {
		t.Fatal(x)
	}
	if x := plan.UpdatedAt; x == nil {
		t.Fatal(x)
	}

	// Add Ons
	if len(plan.AddOns.AddOns) == 0 {
		t.Fatal(plan.AddOns)
	}
	addOn := plan.AddOns.AddOns[0]
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

	// Discounts
	if len(plan.Discounts.Discounts) == 0 {
		t.Fatal(plan.Discounts)
	}
	discount := plan.Discounts.Discounts[0]
	t.Log(discount)
	if discount.Id != "test_discount" {
		t.Fatalf("expected Id to be %s, was %s", "test_discount", discount.Id)
	} else if discount.Amount.Cmp(NewDecimal(1000, 2)) != 0 {
		t.Fatalf("expected Amount to be %s, was %s", NewDecimal(1000, 2), discount.Amount)
	} else if discount.Kind != ModificationKindDiscount {
		t.Fatalf("expected Kind to be %s, was %s", ModificationKindDiscount, discount.Kind)
	} else if discount.Name != "test_discount_name" {
		t.Fatalf("expected Name to be %s, was %s", "test_discount_name", discount.Name)
	} else if discount.NeverExpires != true {
		t.Fatalf("expected NeverExpires to be %v, was %v", true, discount.NeverExpires)
	} else if discount.Description != "A test discount" {
		t.Fatalf("expected Description to be %s, was %s", "A test discount", discount.Description)
	}

	// Find
	plan2, err := g.Find("test_plan_2")
	if err != nil {
		t.Fatal(err)
	}
	if plan2.Id != "test_plan_2" {
		t.Fatal(plan2)
	}
	if len(plan2.AddOns.AddOns) != 0 {
		t.Fatal(plan2.AddOns)
	}
	if len(plan2.Discounts.Discounts) != 0 {
		t.Fatal(plan2.Discounts)
	}
}
