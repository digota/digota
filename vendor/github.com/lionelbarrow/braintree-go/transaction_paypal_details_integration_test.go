package braintree

import "testing"

func TestTransactionPayPalDetails(t *testing.T) {
	tx, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:               "sale",
		Amount:             NewDecimal(2000, 2),
		PaymentMethodNonce: FakeNoncePayPalOneTimePayment,
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

	if tx.PayPalDetails == nil {
		t.Fatal("Expected PayPalDetails for transaction created with PayPal nonce")
	}

	t.Log(tx.PayPalDetails)

	if tx.PayPalDetails.PayerEmail == "" {
		t.Fatal("Expected PayPalDetails to have PayerEmail set")
	}
	if tx.PayPalDetails.PaymentID == "" {
		t.Fatal("Expected PayPalDetails to have PaymentID set")
	}
	if tx.PayPalDetails.ImageURL == "" {
		t.Fatal("Expected PayPalDetails to have ImageURL set")
	}
	if tx.PayPalDetails.DebugID == "" {
		t.Fatal("Expected PayPalDetails to have DebugID set")
	}
	if tx.PayPalDetails.PayerID == "" {
		t.Fatal("Expected PayPalDetails to have DebugID set")
	}
	if tx.PayPalDetails.PayerFirstName == "" {
		t.Fatal("Expected PayPalDetails to have PayerFirstName set")
	}
	if tx.PayPalDetails.PayerLastName == "" {
		t.Fatal("Expected PayPalDetails to have PayerLastName set")
	}
	if tx.PayPalDetails.PayerStatus == "" {
		t.Fatal("Expected PayPalDetails to have PayerStatus set")
	}
	if tx.PayPalDetails.SellerProtectionStatus == "" {
		t.Fatal("Expected PayPalDetails to have SellerProtectionStatus set")
	}
}

func TestTransactionWithoutPayPalDetails(t *testing.T) {
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

	if tx.PayPalDetails != nil {
		t.Fatalf("Expected PayPalDetails to be nil for transaction created without PayPal, but was %#v", tx.PayPalDetails)
	}
}
