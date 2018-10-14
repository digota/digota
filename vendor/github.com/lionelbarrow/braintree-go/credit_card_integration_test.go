package braintree

import (
	"testing"

	"github.com/lionelbarrow/braintree-go/testhelpers"
)

func TestCreditCard(t *testing.T) {
	t.Parallel()

	cust, err := testGateway.Customer().Create(&Customer{})
	if err != nil {
		t.Fatal(err)
	}

	g := testGateway.CreditCard()
	card, err := g.Create(&CreditCard{
		CustomerId:     cust.Id,
		Number:         testCreditCards["visa"].Number,
		ExpirationDate: "05/14",
		CVV:            "100",
		Options: &CreditCardOptions{
			VerifyCard: testhelpers.BoolPtr(true),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(card)

	if card.Token == "" {
		t.Fatal("invalid token")
	}

	// Update
	card2, err := g.Update(&CreditCard{
		Token:          card.Token,
		Number:         testCreditCards["mastercard"].Number,
		ExpirationDate: "05/14",
		CVV:            "100",
		Options: &CreditCardOptions{
			VerifyCard: testhelpers.BoolPtr(true),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(card2)

	if card2.Token != card.Token {
		t.Fatal("tokens do not match")
	}
	if card2.CardType != "MasterCard" {
		t.Fatal("card type does not match")
	}

	// Delete
	err = g.Delete(card2)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreditCardFailedAutoVerification(t *testing.T) {
	t.Parallel()

	cust, err := testGateway.Customer().Create(&Customer{})
	if err != nil {
		t.Fatal(err)
	}

	g := testGateway.CreditCard()
	card, err := g.Create(&CreditCard{
		CustomerId:         cust.Id,
		PaymentMethodNonce: FakeNonceProcessorDeclinedVisa,
	})
	if err == nil {
		t.Fatal("Got no error, want error")
	}
	if g, w := err.(*BraintreeError).ErrorMessage, "Do Not Honor"; g != w {
		t.Fatalf("Got error %q, want error %q", g, w)
	}

	t.Logf("%#v\n", err)
	t.Logf("%#v\n", card)
}

func TestCreditCardForceNotVerified(t *testing.T) {
	t.Parallel()

	cust, err := testGateway.Customer().Create(&Customer{})
	if err != nil {
		t.Fatal(err)
	}

	g := testGateway.CreditCard()
	card, err := g.Create(&CreditCard{
		CustomerId:         cust.Id,
		PaymentMethodNonce: FakeNonceProcessorDeclinedVisa,
		Options: &CreditCardOptions{
			VerifyCard: testhelpers.BoolPtr(false),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%#v\n", card)
}

func TestCreateCreditCardWithExpirationMonthAndYear(t *testing.T) {
	t.Parallel()

	customer, err := testGateway.Customer().Create(&Customer{})
	if err != nil {
		t.Fatal(err)
	}
	card, err := testGateway.CreditCard().Create(&CreditCard{
		CustomerId:      customer.Id,
		Number:          testCreditCards["visa"].Number,
		ExpirationMonth: "05",
		ExpirationYear:  "2014",
		CVV:             "100",
	})

	if err != nil {
		t.Fatal(err)
	}
	if card.Token == "" {
		t.Fatal("invalid token")
	}
}

func TestCreateCreditCardInvalidInput(t *testing.T) {
	t.Parallel()

	card, err := testGateway.CreditCard().Create(&CreditCard{
		Number:         testCreditCards["visa"].Number,
		ExpirationDate: "05/14",
	})

	t.Log(card)

	// This test should fail because customer id is required
	if err == nil {
		t.Fatal("expected to get error creating card because of required fields, but did not")
	}

	// TODO: validate fields
}

func TestFindCreditCard(t *testing.T) {
	t.Parallel()

	customer, err := testGateway.Customer().Create(&Customer{})
	if err != nil {
		t.Fatal(err)
	}
	card, err := testGateway.CreditCard().Create(&CreditCard{
		CustomerId:     customer.Id,
		Number:         testCreditCards["visa"].Number,
		ExpirationDate: "05/14",
		CVV:            "100",
		Options: &CreditCardOptions{
			VerifyCard: testhelpers.BoolPtr(true),
		},
	})

	t.Log(card)

	if err != nil {
		t.Fatal(err)
	}
	if card.Token == "" {
		t.Fatal("invalid token")
	}

	card2, err := testGateway.CreditCard().Find(card.Token)

	t.Log(card2)

	if err != nil {
		t.Fatal(err)
	}
	if card2.Token != card.Token {
		t.Fatal("tokens do not match")
	}
}

func TestFindCreditCardBadData(t *testing.T) {
	t.Parallel()

	card, err := testGateway.CreditCard().Find("invalid_token")

	t.Log(card)

	if err == nil {
		t.Fatal("expected to get error because the token is invalid")
	}
}

func TestSaveCreditCardWithVenmoSDKPaymentMethodCode(t *testing.T) {
	t.Parallel()

	customer, err := testGateway.Customer().Create(&Customer{})
	if err != nil {
		t.Fatal(err)
	}
	card, err := testGateway.CreditCard().Create(&CreditCard{
		CustomerId:                customer.Id,
		VenmoSDKPaymentMethodCode: "stub-" + testCreditCards["visa"].Number,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !card.VenmoSDK {
		t.Fatal("venmo card not marked")
	}
}

func TestSaveCreditCardWithVenmoSDKSession(t *testing.T) {
	t.Parallel()

	customer, err := testGateway.Customer().Create(&Customer{})
	if err != nil {
		t.Fatal(err)
	}
	card, err := testGateway.CreditCard().Create(&CreditCard{
		CustomerId:     customer.Id,
		Number:         testCreditCards["visa"].Number,
		ExpirationDate: "05/14",
		Options: &CreditCardOptions{
			VenmoSDKSession: "stub-session",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !card.VenmoSDK {
		t.Fatal("venmo card not marked")
	}
}
