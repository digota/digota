package braintree

import "testing"

func TestTransactionApplePayDetails(t *testing.T) {
	tx, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:               "sale",
		Amount:             NewDecimal(2000, 2),
		PaymentMethodNonce: FakeNonceApplePayVisa,
	})

	t.Log(tx)

	if err != nil {
		t.Fatal(err)
	}
	if tx.Id == "" {
		t.Fatal("Received invalid ID on new transaction")
	}
	if tx.Status != TransactionStatusAuthorized {
		t.Fatal(tx.Status)
	}

	if tx.ApplePayDetails == nil {
		t.Fatal("Expected ApplePayDetails for transaction created with ApplePay nonce")
	}

	t.Log(tx.ApplePayDetails)

	if tx.ApplePayDetails.CardType == "" {
		t.Fatal("Expected ApplePayDetails to have CardType set")
	}
	if tx.ApplePayDetails.PaymentInstrumentName == "" {
		t.Fatal("Expected ApplePayDetails to have PaymentInstrumentName set")
	}
	if tx.ApplePayDetails.SourceDescription == "" {
		t.Fatal("Expected ApplePayDetails to have SourceDescription set")
	}
	if tx.ApplePayDetails.CardholderName == "" {
		t.Fatal("Expected ApplePayDetails to have CardholderName set")
	}
	if tx.ApplePayDetails.ExpirationMonth == "" {
		t.Fatal("Expected ApplePayDetails to have ExpirationMonth set")
	}
	if tx.ApplePayDetails.ExpirationYear == "" {
		t.Fatal("Expected ApplePayDetails to have ExpirationYear set")
	}
	if tx.ApplePayDetails.Last4 == "" {
		t.Fatal("Expected ApplePayDetails to have Last3 set")
	}
}

func TestTransactionWithoutApplePayDetails(t *testing.T) {
	tx, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:               "sale",
		Amount:             NewDecimal(2000, 2),
		PaymentMethodNonce: FakeNonceTransactable,
	})

	t.Log(tx)

	if err != nil {
		t.Fatal(err)
	}
	if tx.Id == "" {
		t.Fatal("Received invalid ID on new transaction")
	}
	if tx.Status != TransactionStatusAuthorized {
		t.Fatal(tx.Status)
	}

	if tx.ApplePayDetails != nil {
		t.Fatalf("Expected ApplePayDetails to be nil for transaction created without ApplePay, but was %#v", tx.ApplePayDetails)
	}
}
