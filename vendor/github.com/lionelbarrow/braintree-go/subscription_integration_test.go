package braintree

import (
	"fmt"
	"testing"
	"time"

	"github.com/lionelbarrow/braintree-go/testhelpers"
)

// This test will fail unless you set up your Braintree sandbox account correctly. See TESTING.md for details.
func TestSubscriptionSimple(t *testing.T) {
	t.Parallel()

	customer, err := testGateway.Customer().Create(&Customer{})
	if err != nil {
		t.Fatal(err)
	}
	paymentMethod, err := testGateway.PaymentMethod().Create(&PaymentMethodRequest{
		CustomerId:         customer.Id,
		PaymentMethodNonce: FakeNonceTransactable,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(customer)

	g := testGateway.Subscription()

	// Create
	sub, err := g.Create(&SubscriptionRequest{
		PaymentMethodToken: paymentMethod.GetToken(),
		PlanId:             "test_plan",
	})

	t.Log("sub1", sub)

	if err != nil {
		t.Fatal(err)
	}
	if sub.Id == "" {
		t.Fatal("invalid subscription id")
	}

	// Update
	sub2, err := g.Update(&SubscriptionRequest{
		Id:     sub.Id,
		PlanId: "test_plan_2",
		Options: &SubscriptionOptions{
			ProrateCharges:                       true,
			RevertSubscriptionOnProrationFailure: true,
			StartImmediately:                     true,
		},
	})

	t.Log("sub2", sub2)

	if err != nil {
		t.Fatal(err)
	}
	if sub2.Id != sub.Id {
		t.Fatal(sub2.Id)
	}
	if x := sub2.PlanId; x != "test_plan_2" {
		t.Fatal(x)
	}

	// Find
	sub3, err := g.Find(sub.Id)
	if err != nil {
		t.Fatal(err)
	}
	if sub3.Id != sub2.Id {
		t.Fatal(sub3.Id)
	}

	// Cancel
	_, err = g.Cancel(sub2.Id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSubscriptionAllFieldsWithBillingDayOfMonth(t *testing.T) {
	t.Parallel()

	customer, err := testGateway.Customer().Create(&Customer{})
	if err != nil {
		t.Fatal(err)
	}
	paymentMethod, err := testGateway.PaymentMethod().Create(&PaymentMethodRequest{
		CustomerId:         customer.Id,
		PaymentMethodNonce: FakeNonceTransactable,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(customer)

	g := testGateway.Subscription()

	// Create
	sub1, err := g.Create(&SubscriptionRequest{
		PaymentMethodToken:    paymentMethod.GetToken(),
		PlanId:                "test_plan",
		MerchantAccountId:     testMerchantAccountId,
		BillingDayOfMonth:     testhelpers.IntPtr(15),
		NumberOfBillingCycles: testhelpers.IntPtr(2),
		Price: NewDecimal(100, 2),
		Descriptor: &Descriptor{
			Name:  "Company Name*Product 1",
			Phone: "0000000000",
			URL:   "example.com",
		},
	})

	t.Log("sub1", sub1)

	if err != nil {
		t.Fatal(err)
	}
	if sub1.Id == "" {
		t.Fatal("invalid subscription id")
	}
	if sub1.BillingDayOfMonth != "15" {
		t.Fatalf("got billing day of month %#v, want %#v", sub1.BillingDayOfMonth, "15")
	}
	if x := sub1.NeverExpires; x {
		t.Fatalf("got never expires %#v, want false", x)
	}
	if x := sub1.NumberOfBillingCycles; x == nil || *x != 2 {
		t.Fatalf("got number billing cycles %#v, want 2", x)
	}
	if x := sub1.Price; x == nil || x.Scale != 2 || x.Unscaled != 100 {
		t.Fatalf("got price %#v, want 1.00", x)
	}
	if x := sub1.TrialPeriod; x {
		t.Fatalf("got trial period %#v, want false", x)
	}
	if x := sub1.Status; x != SubscriptionStatusPending && x != SubscriptionStatusActive {
		t.Fatalf("got status %#v, want Pending or Active (it will be active if todays date matches the billing day of month)", x)
	}
	if x := sub1.Descriptor.Name; x != "Company Name*Product 1" {
		t.Fatalf("got descriptor name %#v, want Company Name*Product 1", x)
	}
	if x := sub1.Descriptor.Phone; x != "0000000000" {
		t.Fatalf("got descriptor phone %#v, want 0000000000", x)
	}
	if x := sub1.Descriptor.URL; x != "example.com" {
		t.Fatalf("got descriptor url %#v, want example.com", x)
	}

	// Update
	sub2, err := g.Update(&SubscriptionRequest{
		Id:     sub1.Id,
		PlanId: "test_plan_2",
		Options: &SubscriptionOptions{
			ProrateCharges:                       true,
			RevertSubscriptionOnProrationFailure: true,
			StartImmediately:                     true,
		},
	})

	t.Log("sub2", sub2)

	if err != nil {
		t.Fatal(err)
	}
	if sub2.Id != sub1.Id {
		t.Fatal(sub2.Id)
	}
	if x := sub2.PlanId; x != "test_plan_2" {
		t.Fatal(x)
	}

	// Find
	sub3, err := g.Find(sub1.Id)
	if err != nil {
		t.Fatal(err)
	}
	if sub3.Id != sub1.Id {
		t.Fatal(sub3.Id)
	}

	// Cancel
	sub4, err := g.Cancel(sub1.Id)
	if err != nil {
		t.Fatal(err)
	}
	if x := sub4.Status; x != SubscriptionStatusCanceled {
		t.Fatalf("got status %#v, want Canceled", x)
	}
}

func TestSubscriptionAllFieldsWithBillingDayOfMonthNeverExpires(t *testing.T) {
	t.Parallel()

	customer, err := testGateway.Customer().Create(&Customer{})
	if err != nil {
		t.Fatal(err)
	}
	paymentMethod, err := testGateway.PaymentMethod().Create(&PaymentMethodRequest{
		CustomerId:         customer.Id,
		PaymentMethodNonce: FakeNonceTransactable,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(customer)

	g := testGateway.Subscription()

	// Create
	sub1, err := g.Create(&SubscriptionRequest{
		PaymentMethodToken: paymentMethod.GetToken(),
		PlanId:             "test_plan",
		MerchantAccountId:  testMerchantAccountId,
		BillingDayOfMonth:  testhelpers.IntPtr(15),
		NeverExpires:       testhelpers.BoolPtr(true),
		Price:              NewDecimal(100, 2),
		Descriptor: &Descriptor{
			Name:  "Company Name*Product 1",
			Phone: "0000000000",
			URL:   "example.com",
		},
	})

	t.Logf("sub1 %#v", sub1)

	if err != nil {
		t.Fatal(err)
	}
	if sub1.Id == "" {
		t.Fatal("invalid subscription id")
	}
	if sub1.BillingDayOfMonth != "15" {
		t.Fatalf("got billing day of month %#v, want %#v", sub1.BillingDayOfMonth, "15")
	}
	if x := sub1.NeverExpires; !x {
		t.Fatalf("got never expires %#v, want true", x)
	}
	if x := sub1.NumberOfBillingCycles; x != nil {
		t.Fatalf("got number billing cycles %#v, didn't want", x)
	}
	if x := sub1.Price; x == nil || x.Scale != 2 || x.Unscaled != 100 {
		t.Fatalf("got price %#v, want 1.00", x)
	}
	if x := sub1.TrialPeriod; x {
		t.Fatalf("got trial period %#v, want false", x)
	}
	if x := sub1.Status; x != SubscriptionStatusPending && x != SubscriptionStatusActive {
		t.Fatalf("got status %#v, want Pending or Active (it will be active if todays date matches the billing day of month)", x)
	}
	if x := sub1.Descriptor.Name; x != "Company Name*Product 1" {
		t.Fatalf("got descriptor name %#v, want Company Name*Product 1", x)
	}
	if x := sub1.Descriptor.Phone; x != "0000000000" {
		t.Fatalf("got descriptor phone %#v, want 0000000000", x)
	}
	if x := sub1.Descriptor.URL; x != "example.com" {
		t.Fatalf("got descriptor url %#v, want example.com", x)
	}

	// Update
	sub2, err := g.Update(&SubscriptionRequest{
		Id:     sub1.Id,
		PlanId: "test_plan_2",
		Options: &SubscriptionOptions{
			ProrateCharges:                       true,
			RevertSubscriptionOnProrationFailure: true,
			StartImmediately:                     true,
		},
	})

	t.Log("sub2", sub2)

	if err != nil {
		t.Fatal(err)
	}
	if sub2.Id != sub1.Id {
		t.Fatal(sub2.Id)
	}
	if x := sub2.PlanId; x != "test_plan_2" {
		t.Fatal(x)
	}

	// Find
	sub3, err := g.Find(sub1.Id)
	if err != nil {
		t.Fatal(err)
	}
	if sub3.Id != sub1.Id {
		t.Fatal(sub3.Id)
	}

	// Cancel
	sub4, err := g.Cancel(sub1.Id)
	if err != nil {
		t.Fatal(err)
	}
	if x := sub4.Status; x != SubscriptionStatusCanceled {
		t.Fatalf("got status %#v, want Canceled", x)
	}
}

func TestSubscriptionAllFieldsWithFirstBillingDate(t *testing.T) {
	t.Parallel()

	customer, err := testGateway.Customer().Create(&Customer{})
	if err != nil {
		t.Fatal(err)
	}
	paymentMethod, err := testGateway.PaymentMethod().Create(&PaymentMethodRequest{
		CustomerId:         customer.Id,
		PaymentMethodNonce: FakeNonceTransactable,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(customer)

	g := testGateway.Subscription()

	// Create
	firstBillingDate := fmt.Sprintf("%d-12-31", time.Now().Year())
	sub1, err := g.Create(&SubscriptionRequest{
		PaymentMethodToken:    paymentMethod.GetToken(),
		PlanId:                "test_plan",
		MerchantAccountId:     testMerchantAccountId,
		FirstBillingDate:      firstBillingDate,
		NumberOfBillingCycles: testhelpers.IntPtr(2),
		Price: NewDecimal(100, 2),
		Descriptor: &Descriptor{
			Name:  "Company Name*Product 1",
			Phone: "0000000000",
			URL:   "example.com",
		},
	})

	t.Log("sub1", sub1)

	if err != nil {
		t.Fatal(err)
	}
	if sub1.Id == "" {
		t.Fatal("invalid subscription id")
	}
	if sub1.BillingDayOfMonth != "31" {
		t.Fatalf("got billing day of month %#v, want %#v", sub1.BillingDayOfMonth, "31")
	}
	if sub1.FirstBillingDate != firstBillingDate {
		t.Fatalf("got first billing date %#v, want %#v", sub1.FirstBillingDate, firstBillingDate)
	}
	if x := sub1.NeverExpires; x {
		t.Fatalf("got never expires %#v, want false", x)
	}
	if x := sub1.NumberOfBillingCycles; x == nil {
		t.Fatalf("got number billing cycles nil, want 2")
	} else if *x != 2 {
		t.Fatalf("got number billing cycles %#v, want 2", *x)
	}
	if x := sub1.Price; x == nil || x.Scale != 2 || x.Unscaled != 100 {
		t.Fatalf("got price %#v, want 1.00", x)
	}
	if x := sub1.TrialPeriod; x {
		t.Fatalf("got trial period %#v, want false", x)
	}
	if x := sub1.Status; x != SubscriptionStatusPending {
		t.Fatalf("got status %#v, want Pending", x)
	}
	if x := sub1.Descriptor.Name; x != "Company Name*Product 1" {
		t.Fatalf("got descriptor name %#v, want Company Name*Product 1", x)
	}
	if x := sub1.Descriptor.Phone; x != "0000000000" {
		t.Fatalf("got descriptor phone %#v, want 0000000000", x)
	}
	if x := sub1.Descriptor.URL; x != "example.com" {
		t.Fatalf("got descriptor url %#v, want example.com", x)
	}

	// Update
	sub2, err := g.Update(&SubscriptionRequest{
		Id:     sub1.Id,
		PlanId: "test_plan_2",
		Options: &SubscriptionOptions{
			ProrateCharges:                       true,
			RevertSubscriptionOnProrationFailure: true,
			StartImmediately:                     true,
		},
	})

	t.Log("sub2", sub2)

	if err != nil {
		t.Fatal(err)
	}
	if sub2.Id != sub1.Id {
		t.Fatal(sub2.Id)
	}
	if x := sub2.PlanId; x != "test_plan_2" {
		t.Fatal(x)
	}

	// Find
	sub3, err := g.Find(sub1.Id)
	if err != nil {
		t.Fatal(err)
	}
	if sub3.Id != sub1.Id {
		t.Fatal(sub3.Id)
	}

	// Cancel
	sub4, err := g.Cancel(sub1.Id)
	if err != nil {
		t.Fatal(err)
	}
	if x := sub4.Status; x != SubscriptionStatusCanceled {
		t.Fatalf("got status %#v, want Canceled", x)
	}
}

func TestSubscriptionAllFieldsWithFirstBillingDateNeverExpires(t *testing.T) {
	t.Parallel()

	customer, err := testGateway.Customer().Create(&Customer{})
	if err != nil {
		t.Fatal(err)
	}
	paymentMethod, err := testGateway.PaymentMethod().Create(&PaymentMethodRequest{
		CustomerId:         customer.Id,
		PaymentMethodNonce: FakeNonceTransactable,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(customer)

	g := testGateway.Subscription()

	// Create
	firstBillingDate := fmt.Sprintf("%d-12-31", time.Now().Year())
	sub1, err := g.Create(&SubscriptionRequest{
		PaymentMethodToken: paymentMethod.GetToken(),
		PlanId:             "test_plan",
		MerchantAccountId:  testMerchantAccountId,
		FirstBillingDate:   firstBillingDate,
		NeverExpires:       testhelpers.BoolPtr(true),
		Price:              NewDecimal(100, 2),
		Descriptor: &Descriptor{
			Name:  "Company Name*Product 1",
			Phone: "0000000000",
			URL:   "example.com",
		},
	})

	t.Log("sub1", sub1)

	if err != nil {
		t.Fatal(err)
	}
	if sub1.Id == "" {
		t.Fatal("invalid subscription id")
	}
	if sub1.BillingDayOfMonth != "31" {
		t.Fatalf("got billing day of month %#v, want %#v", sub1.BillingDayOfMonth, "31")
	}
	if sub1.FirstBillingDate != firstBillingDate {
		t.Fatalf("got first billing date %#v, want %#v", sub1.FirstBillingDate, firstBillingDate)
	}
	if x := sub1.NeverExpires; !x {
		t.Fatalf("got never expires %#v, want true", x)
	}
	if x := sub1.NumberOfBillingCycles; x != nil {
		t.Fatalf("got number billing cycles %#v, didn't want", x)
	}
	if x := sub1.Price; x == nil || x.Scale != 2 || x.Unscaled != 100 {
		t.Fatalf("got price %#v, want 1.00", x)
	}
	if x := sub1.TrialPeriod; x {
		t.Fatalf("got trial period %#v, want false", x)
	}
	if x := sub1.Status; x != SubscriptionStatusPending {
		t.Fatalf("got status %#v, want Pending", x)
	}
	if x := sub1.Descriptor.Name; x != "Company Name*Product 1" {
		t.Fatalf("got descriptor name %#v, want Company Name*Product 1", x)
	}
	if x := sub1.Descriptor.Phone; x != "0000000000" {
		t.Fatalf("got descriptor phone %#v, want 0000000000", x)
	}
	if x := sub1.Descriptor.URL; x != "example.com" {
		t.Fatalf("got descriptor url %#v, want example.com", x)
	}

	// Update
	sub2, err := g.Update(&SubscriptionRequest{
		Id:     sub1.Id,
		PlanId: "test_plan_2",
		Options: &SubscriptionOptions{
			ProrateCharges:                       true,
			RevertSubscriptionOnProrationFailure: true,
			StartImmediately:                     true,
		},
	})

	t.Log("sub2", sub2)

	if err != nil {
		t.Fatal(err)
	}
	if sub2.Id != sub1.Id {
		t.Fatal(sub2.Id)
	}
	if x := sub2.PlanId; x != "test_plan_2" {
		t.Fatal(x)
	}

	// Find
	sub3, err := g.Find(sub1.Id)
	if err != nil {
		t.Fatal(err)
	}
	if sub3.Id != sub1.Id {
		t.Fatal(sub3.Id)
	}

	// Cancel
	sub4, err := g.Cancel(sub1.Id)
	if err != nil {
		t.Fatal(err)
	}
	if x := sub4.Status; x != SubscriptionStatusCanceled {
		t.Fatalf("got status %#v, want Canceled", x)
	}
}

func TestSubscriptionAllFieldsWithTrialPeriod(t *testing.T) {
	t.Parallel()

	customer, err := testGateway.Customer().Create(&Customer{})
	if err != nil {
		t.Fatal(err)
	}
	paymentMethod, err := testGateway.PaymentMethod().Create(&PaymentMethodRequest{
		CustomerId:         customer.Id,
		PaymentMethodNonce: FakeNonceTransactable,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(customer)

	g := testGateway.Subscription()

	// Create
	firstBillingDate := time.Now().In(testTimeZone).AddDate(0, 0, 7)
	sub1, err := g.Create(&SubscriptionRequest{
		PaymentMethodToken:    paymentMethod.GetToken(),
		PlanId:                "test_plan",
		MerchantAccountId:     testMerchantAccountId,
		TrialPeriod:           testhelpers.BoolPtr(true),
		TrialDuration:         "7",
		TrialDurationUnit:     SubscriptionTrialDurationUnitDay,
		NumberOfBillingCycles: testhelpers.IntPtr(2),
		Price: NewDecimal(100, 2),
		Descriptor: &Descriptor{
			Name:  "Company Name*Product 1",
			Phone: "0000000000",
			URL:   "example.com",
		},
	})

	t.Log("sub1", sub1)

	if err != nil {
		t.Fatal(err)
	}
	if sub1.Id == "" {
		t.Fatal("invalid subscription id")
	}
	if sub1.BillingDayOfMonth != fmt.Sprintf("%d", firstBillingDate.Day()) {
		t.Fatalf("got billing day of month %#v, want %#v", sub1.BillingDayOfMonth, firstBillingDate.Day())
	}
	if sub1.FirstBillingDate != firstBillingDate.Format("2006-01-02") {
		t.Fatalf("got first billing date %#v, want %#v", sub1.FirstBillingDate, firstBillingDate)
	}
	if x := sub1.NeverExpires; x {
		t.Fatalf("got never expires %#v, want false", x)
	}
	if x := sub1.NumberOfBillingCycles; x == nil || *x != 2 {
		t.Fatalf("got number billing cycles %#v, want 2", x)
	}
	if x := sub1.Price; x == nil || x.Scale != 2 || x.Unscaled != 100 {
		t.Fatalf("got price %#v, want 1.00", x)
	}
	if x := sub1.TrialPeriod; !x {
		t.Fatalf("got trial period %#v, want false", x)
	}
	if sub1.TrialDuration != "7" {
		t.Fatalf("got trial duration %#v, want 7", sub1.TrialDuration)
	}
	if sub1.TrialDurationUnit != SubscriptionTrialDurationUnitDay {
		t.Fatalf("got trial duration unit %#v, want day", sub1.TrialDurationUnit)
	}
	if x := sub1.Status; x != SubscriptionStatusActive {
		t.Fatalf("got status %#v, want Active", x)
	}
	if x := sub1.Descriptor.Name; x != "Company Name*Product 1" {
		t.Fatalf("got descriptor name %#v, want Company Name*Product 1", x)
	}
	if x := sub1.Descriptor.Phone; x != "0000000000" {
		t.Fatalf("got descriptor phone %#v, want 0000000000", x)
	}
	if x := sub1.Descriptor.URL; x != "example.com" {
		t.Fatalf("got descriptor url %#v, want example.com", x)
	}

	// Update
	sub2, err := g.Update(&SubscriptionRequest{
		Id:     sub1.Id,
		PlanId: "test_plan_2",
		Options: &SubscriptionOptions{
			ProrateCharges:                       true,
			RevertSubscriptionOnProrationFailure: true,
			StartImmediately:                     true,
		},
	})

	t.Log("sub2", sub2)

	if err != nil {
		t.Fatal(err)
	}
	if sub2.Id != sub1.Id {
		t.Fatal(sub2.Id)
	}
	if x := sub2.PlanId; x != "test_plan_2" {
		t.Fatal(x)
	}

	// Find
	sub3, err := g.Find(sub1.Id)
	if err != nil {
		t.Fatal(err)
	}
	if sub3.Id != sub1.Id {
		t.Fatal(sub3.Id)
	}

	// Cancel
	_, err = g.Cancel(sub1.Id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSubscriptionAllFieldsWithTrialPeriodNeverExpires(t *testing.T) {
	t.Parallel()

	customer, err := testGateway.Customer().Create(&Customer{})
	if err != nil {
		t.Fatal(err)
	}
	paymentMethod, err := testGateway.PaymentMethod().Create(&PaymentMethodRequest{
		CustomerId:         customer.Id,
		PaymentMethodNonce: FakeNonceTransactable,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(customer)

	g := testGateway.Subscription()

	// Create
	firstBillingDate := time.Now().In(testTimeZone).AddDate(0, 0, 7)
	sub1, err := g.Create(&SubscriptionRequest{
		PaymentMethodToken: paymentMethod.GetToken(),
		PlanId:             "test_plan",
		MerchantAccountId:  testMerchantAccountId,
		TrialPeriod:        testhelpers.BoolPtr(true),
		TrialDuration:      "7",
		TrialDurationUnit:  SubscriptionTrialDurationUnitDay,
		NeverExpires:       testhelpers.BoolPtr(true),
		Price:              NewDecimal(100, 2),
		Descriptor: &Descriptor{
			Name:  "Company Name*Product 1",
			Phone: "0000000000",
			URL:   "example.com",
		},
	})

	t.Log("sub1", sub1)

	if err != nil {
		t.Fatal(err)
	}
	if sub1.Id == "" {
		t.Fatal("invalid subscription id")
	}
	if sub1.BillingDayOfMonth != fmt.Sprintf("%d", firstBillingDate.Day()) {
		t.Fatalf("got billing day of month %#v, want %#v", sub1.BillingDayOfMonth, firstBillingDate.Day())
	}
	if sub1.FirstBillingDate != firstBillingDate.Format("2006-01-02") {
		t.Fatalf("got first billing date %#v, want %#v", sub1.FirstBillingDate, firstBillingDate)
	}
	if x := sub1.NeverExpires; !x {
		t.Fatalf("got never expires %#v, want true", x)
	}
	if x := sub1.NumberOfBillingCycles; x != nil {
		t.Fatalf("got number billing cycles %#v, didn't want", x)
	}
	if x := sub1.Price; x == nil || x.Scale != 2 || x.Unscaled != 100 {
		t.Fatalf("got price %#v, want 1.00", x)
	}
	if x := sub1.TrialPeriod; !x {
		t.Fatalf("got trial period %#v, want false", x)
	}
	if sub1.TrialDuration != "7" {
		t.Fatalf("got trial duration %#v, want 7", sub1.TrialDuration)
	}
	if sub1.TrialDurationUnit != SubscriptionTrialDurationUnitDay {
		t.Fatalf("got trial duration unit %#v, want day", sub1.TrialDurationUnit)
	}
	if x := sub1.Status; x != SubscriptionStatusActive {
		t.Fatalf("got status %#v, want Active", x)
	}
	if x := sub1.Descriptor.Name; x != "Company Name*Product 1" {
		t.Fatalf("got descriptor name %#v, want Company Name*Product 1", x)
	}
	if x := sub1.Descriptor.Phone; x != "0000000000" {
		t.Fatalf("got descriptor phone %#v, want 0000000000", x)
	}
	if x := sub1.Descriptor.URL; x != "example.com" {
		t.Fatalf("got descriptor url %#v, want example.com", x)
	}

	// Update
	sub2, err := g.Update(&SubscriptionRequest{
		Id:     sub1.Id,
		PlanId: "test_plan_2",
		Options: &SubscriptionOptions{
			ProrateCharges:                       true,
			RevertSubscriptionOnProrationFailure: true,
			StartImmediately:                     true,
		},
	})

	t.Log("sub2", sub2)

	if err != nil {
		t.Fatal(err)
	}
	if sub2.Id != sub1.Id {
		t.Fatal(sub2.Id)
	}
	if x := sub2.PlanId; x != "test_plan_2" {
		t.Fatal(x)
	}

	// Find
	sub3, err := g.Find(sub1.Id)
	if err != nil {
		t.Fatal(err)
	}
	if sub3.Id != sub1.Id {
		t.Fatal(sub3.Id)
	}

	// Cancel
	_, err = g.Cancel(sub1.Id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSubscriptionModifications(t *testing.T) {
	t.Parallel()

	customer, err := testGateway.Customer().Create(&Customer{})
	if err != nil {
		t.Fatal(err)
	}
	paymentMethod, err := testGateway.PaymentMethod().Create(&PaymentMethodRequest{
		CustomerId:         customer.Id,
		PaymentMethodNonce: FakeNonceTransactable,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(customer)

	g := testGateway.Subscription()

	// Create
	sub, err := g.Create(&SubscriptionRequest{
		PaymentMethodToken: paymentMethod.GetToken(),
		PlanId:             "test_plan_2",
	})

	t.Log("sub1", sub)

	if err != nil {
		t.Fatal(err)
	}
	if sub.Id == "" {
		t.Fatal("invalid subscription id")
	}

	// Add AddOn
	sub2, err := g.Update(&SubscriptionRequest{
		Id: sub.Id,
		AddOns: &ModificationsRequest{
			Add: []AddModificationRequest{
				{
					InheritedFromID: "test_add_on",
					ModificationRequest: ModificationRequest{
						Amount:       NewDecimal(300, 2),
						Quantity:     1,
						NeverExpires: true,
					},
				},
			},
		},
		Discounts: &ModificationsRequest{
			Add: []AddModificationRequest{
				{
					InheritedFromID: "test_discount",
					ModificationRequest: ModificationRequest{
						Amount:       NewDecimal(100, 2),
						Quantity:     1,
						NeverExpires: true,
					},
				},
			},
		},
	})

	t.Log("sub2", sub2)

	if err != nil {
		t.Fatal(err)
	}
	if sub2.Id != sub.Id {
		t.Fatal(sub2.Id)
	}
	if x := sub2.PlanId; x != "test_plan_2" {
		t.Fatal(x)
	}
	if x := sub2.AddOns.AddOns; len(x) != 1 {
		t.Fatalf("got %d add ons, want 1 add on", len(x))
	}
	if x := sub2.AddOns.AddOns[0].Amount; x.String() != NewDecimal(300, 2).String() {
		t.Fatalf("got %v add on, want 3.00 add on", x)
	}
	if x := sub2.Discounts.Discounts; len(x) != 1 {
		t.Fatalf("got %d discounts, want 1 discount", len(x))
	}
	if x := sub2.Discounts.Discounts[0].Amount; x.String() != NewDecimal(100, 2).String() {
		t.Fatalf("got %v discount, want 1.00 discount", x)
	}

	// Update AddOn
	sub3, err := g.Update(&SubscriptionRequest{
		Id: sub.Id,
		AddOns: &ModificationsRequest{
			Update: []UpdateModificationRequest{
				{
					ExistingID: "test_add_on",
					ModificationRequest: ModificationRequest{
						Amount: NewDecimal(150, 2),
					},
				},
			},
		},
		Discounts: &ModificationsRequest{
			RemoveExistingIDs: []string{
				"test_discount",
			},
		},
	})

	t.Log("sub3", sub3)

	if err != nil {
		t.Fatal(err)
	}
	if sub3.Id != sub.Id {
		t.Fatal(sub3.Id)
	}
	if x := sub3.PlanId; x != "test_plan_2" {
		t.Fatal(x)
	}
	if x := sub3.AddOns.AddOns; len(x) != 1 {
		t.Fatalf("got %d add ons, want 1 add on", len(x))
	}
	if x := sub3.AddOns.AddOns[0].Amount; x.String() != NewDecimal(150, 2).String() {
		t.Fatalf("got %v add on, want 1.50 add on", x)
	}
	if x := sub3.Discounts.Discounts; len(x) != 0 {
		t.Fatalf("got %d discounts, want 0 discounts", len(x))
	}

	// Cancel
	_, err = g.Cancel(sub3.Id)
	if err != nil {
		t.Fatal(err)
	}
}

// This test will fail unless you set up your Braintree sandbox account correctly. See TESTING.md for details.
func TestSubscriptionTransactions(t *testing.T) {
	t.Parallel()

	customer, err := testGateway.Customer().Create(&Customer{})
	if err != nil {
		t.Fatal(err)
	}
	paymentMethod, err := testGateway.PaymentMethod().Create(&PaymentMethodRequest{
		CustomerId:         customer.Id,
		PaymentMethodNonce: FakeNonceTransactable,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(customer)

	g := testGateway.Subscription()

	// Create
	sub, err := g.Create(&SubscriptionRequest{
		PaymentMethodToken: paymentMethod.GetToken(),
		PlanId:             "test_plan",
		Options: &SubscriptionOptions{
			StartImmediately: true,
		},
	})

	t.Log("sub1", sub)

	if err != nil {
		t.Fatal(err)
	}
	if sub.Id == "" {
		t.Fatal("invalid subscription id")
	}

	// Find
	sub2, err := g.Find(sub.Id)
	if err != nil {
		t.Fatal(err)
	}
	if sub2.Id != sub.Id {
		t.Fatal(sub2.Id)
	}
	if x := sub2.PlanId; x != "test_plan" {
		t.Fatal(x)
	}
	if len(sub2.Transactions.Transaction) < 1 {
		t.Fatalf("Expected transactions slice not to be empty")
	}
	if x := sub2.Transactions.Transaction[0].PlanId; x != "test_plan" {
		t.Fatal(x)
	}
	if x := sub2.Transactions.Transaction[0].SubscriptionId; x != sub.Id {
		t.Fatal(x)
	}

	// Cancel
	_, err = g.Cancel(sub2.Id)
	if err != nil {
		t.Fatal(err)
	}
}
