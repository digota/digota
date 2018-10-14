package braintree

import "testing"

func TestTransactionAndroidPayDetails_AndroidPayProxyCardNonce(t *testing.T) {
	tx, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:               "sale",
		Amount:             NewDecimal(2000, 2),
		PaymentMethodNonce: FakeNonceAndroidPayDiscover,
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

	if tx.AndroidPayDetails == nil {
		t.Fatal("Expected AndroidPayDetails for transaction created with AndroidPay nonce")
	}

	t.Log(tx.AndroidPayDetails)

	if tx.AndroidPayDetails.CardType == "" {
		t.Fatal("Expected AndroidPayDetails to have CardType set")
	}
	if tx.AndroidPayDetails.Last4 == "" {
		t.Fatal("Expected AndroidPayDetails to have Last4 set")
	}
	if tx.AndroidPayDetails.SourceCardType == "" {
		t.Fatal("Expected AndroidPayDetails to have SourceCardType set")
	}
	if tx.AndroidPayDetails.SourceCardLast4 == "" {
		t.Fatal("Expected AndroidPayDetails to have SourceCardLast4 set")
	}
	if tx.AndroidPayDetails.SourceDescription == "" {
		t.Fatal("Expected AndroidPayDetails to have SourceDescription set")
	}
	if tx.AndroidPayDetails.VirtualCardType == "" {
		t.Fatal("Expected AndroidPayDetails to have VirtualCardType set")
	}
	if tx.AndroidPayDetails.VirtualCardLast4 == "" {
		t.Fatal("Expected AndroidPayDetails to have VirtualCardLast4 set")
	}
	if tx.AndroidPayDetails.ExpirationMonth == "" {
		t.Fatal("Expected AndroidPayDetails to have ExpirationMonth set")
	}
	if tx.AndroidPayDetails.ExpirationYear == "" {
		t.Fatal("Expected AndroidPayDetails to have ExpirationYear set")
	}
	if tx.AndroidPayDetails.GoogleTransactionID == "" {
		t.Fatal("Expected AndroidPayDetails to have GoogleTransactionID set")
	}
	if tx.AndroidPayDetails.BIN == "" {
		t.Fatal("Expected AndroidPayDetails to have BIN set")
	}
}

func TestTransactionAndroidPayDetails_AndroidPayNetworkTokenNonce(t *testing.T) {
	tx, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:               "sale",
		Amount:             NewDecimal(2000, 2),
		PaymentMethodNonce: FakeNonceAndroidPayMasterCard,
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

	if tx.AndroidPayDetails == nil {
		t.Fatal("Expected AndroidPayDetails for transaction created with AndroidPay nonce")
	}

	t.Log(tx.AndroidPayDetails)

	if tx.AndroidPayDetails.CardType == "" {
		t.Fatal("Expected AndroidPayDetails to have CardType set")
	}
	if tx.AndroidPayDetails.Last4 == "" {
		t.Fatal("Expected AndroidPayDetails to have Last4 set")
	}
	if tx.AndroidPayDetails.SourceCardType == "" {
		t.Fatal("Expected AndroidPayDetails to have SourceCardType set")
	}
	if tx.AndroidPayDetails.SourceCardLast4 == "" {
		t.Fatal("Expected AndroidPayDetails to have SourceCardLast4 set")
	}
	if tx.AndroidPayDetails.SourceDescription == "" {
		t.Fatal("Expected AndroidPayDetails to have SourceDescription set")
	}
	if tx.AndroidPayDetails.VirtualCardType == "" {
		t.Fatal("Expected AndroidPayDetails to have VirtualCardType set")
	}
	if tx.AndroidPayDetails.VirtualCardLast4 == "" {
		t.Fatal("Expected AndroidPayDetails to have VirtualCardLast4 set")
	}
	if tx.AndroidPayDetails.ExpirationMonth == "" {
		t.Fatal("Expected AndroidPayDetails to have ExpirationMonth set")
	}
	if tx.AndroidPayDetails.ExpirationYear == "" {
		t.Fatal("Expected AndroidPayDetails to have ExpirationYear set")
	}
	if tx.AndroidPayDetails.GoogleTransactionID == "" {
		t.Fatal("Expected AndroidPayDetails to have GoogleTransactionID set")
	}
	if tx.AndroidPayDetails.BIN == "" {
		t.Fatal("Expected AndroidPayDetails to have BIN set")
	}
}

func TestTransactionWithoutAndroidPayDetails(t *testing.T) {
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

	if tx.AndroidPayDetails != nil {
		t.Fatalf("Expected AndroidPayDetails to be nil for transaction created without AndroidPay, but was %#v", tx.AndroidPayDetails)
	}
}
