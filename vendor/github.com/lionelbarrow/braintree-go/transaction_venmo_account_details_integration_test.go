package braintree

import "testing"

func TestTransactionVenmoAccountDetails(t *testing.T) {
	tx, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:               "sale",
		Amount:             NewDecimal(2000, 2),
		PaymentMethodNonce: FakeNonceVenmoAccount,
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

	if tx.VenmoAccountDetails == nil {
		t.Fatal("Expected VenmoAccountDetails for transaction created with VenmoAccount nonce")
	}

	t.Log(tx.VenmoAccountDetails)

	if tx.VenmoAccountDetails.Token != "" {
		t.Fatal("Expected VenmoAccountDetails to not have Token set")
	}
	if tx.VenmoAccountDetails.Username == "" {
		t.Fatal("Expected VenmoAccountDetails to have Username set")
	}
	if tx.VenmoAccountDetails.VenmoUserID == "" {
		t.Fatal("Expected VenmoAccountDetails to have VenmoUserID set")
	}
	if tx.VenmoAccountDetails.SourceDescription == "" {
		t.Fatal("Expected VenmoAccountDetails to have SourceDescription set")
	}
	if tx.VenmoAccountDetails.ImageURL == "" {
		t.Fatal("Expected VenmoAccountDetails to have ImageURL set")
	}
}

func TestTransactionWithoutVenmoAccountDetails(t *testing.T) {
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

	if tx.VenmoAccountDetails != nil {
		t.Fatalf("Expected VenmoAccountDetails to be nil for transaction created without a VenmoAccount, but was %#v", tx.VenmoAccountDetails)
	}
}
