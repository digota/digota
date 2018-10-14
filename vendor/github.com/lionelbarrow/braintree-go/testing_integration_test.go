package braintree

import "testing"

func TestSettleTransaction(t *testing.T) {
	t.Parallel()

	txn, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:   "sale",
		Amount: randomAmount(),
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	txn, err = testGateway.Transaction().SubmitForSettlement(txn.Id, txn.Amount)
	if err != nil {
		t.Fatal(err)
	}

	prodGateway := New(Production, "my_merchant_id", "my_public_key", "my_private_key")

	_, err = prodGateway.Testing().Settle(txn.Id)
	if err.Error() != "Operation not allowed in production environment" {
		t.Log(testGateway.Environment())
		t.Fatal(err)
	}

	txn, err = testGateway.Testing().Settle(txn.Id)
	if err != nil {
		t.Fatal(err)
	}

	if txn.Status != TransactionStatusSettled {
		t.Fatal(txn.Status)
	}
}

func TestSettlementConfirmTransaction(t *testing.T) {
	txn, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:   "sale",
		Amount: randomAmount(),
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	txn, err = testGateway.Transaction().SubmitForSettlement(txn.Id, txn.Amount)
	if err != nil {
		t.Fatal(err)
	}

	prodGateway := New(Production, "my_merchant_id", "my_public_key", "my_private_key")

	_, err = prodGateway.Testing().SettlementConfirm(txn.Id)
	if err.Error() != "Operation not allowed in production environment" {
		t.Log(prodGateway.Environment())
		t.Fatal(err)
	}

	txn, err = testGateway.Testing().SettlementConfirm(txn.Id)
	if err != nil {
		t.Fatal(err)
	}

	if txn.Status != TransactionStatusSettlementConfirmed {
		t.Fatal(txn.Status)
	}
}

func TestSettlementDeclinedTransaction(t *testing.T) {
	txn, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:   "sale",
		Amount: randomAmount(),
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	txn, err = testGateway.Transaction().SubmitForSettlement(txn.Id, txn.Amount)
	if err != nil {
		t.Fatal(err)
	}

	prodGateway := New(Production, "my_merchant_id", "my_public_key", "my_private_key")

	_, err = prodGateway.Testing().SettlementDecline(txn.Id)
	if err.Error() != "Operation not allowed in production environment" {
		t.Log(prodGateway.Environment())
		t.Fatal(err)
	}

	txn, err = testGateway.Testing().SettlementDecline(txn.Id)
	if err != nil {
		t.Fatal(err)
	}

	if txn.Status != TransactionStatusSettlementDeclined {
		t.Fatal(txn.Status)
	}
}

func TestSettlementPendingTransaction(t *testing.T) {
	txn, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:   "sale",
		Amount: randomAmount(),
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	txn, err = testGateway.Transaction().SubmitForSettlement(txn.Id, txn.Amount)
	if err != nil {
		t.Fatal(err)
	}

	prodGateway := New(Production, "my_merchant_id", "my_public_key", "my_private_key")

	_, err = prodGateway.Testing().SettlementPending(txn.Id)
	if err.Error() != "Operation not allowed in production environment" {
		t.Log(prodGateway.Environment())
		t.Fatal(err)
	}

	txn, err = testGateway.Testing().SettlementPending(txn.Id)
	if err != nil {
		t.Fatal(err)
	}

	if txn.Status != TransactionStatusSettlementPending {
		t.Fatal(txn.Status)
	}
}

func TestTransactionCreateSettleCheckCreditCardDetails(t *testing.T) {
	t.Parallel()

	amount := NewDecimal(10000, 2)
	txn, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:   "sale",
		Amount: amount,
		CreditCard: &CreditCard{
			Number:         testCreditCards["discover"].Number,
			ExpirationDate: "05/14",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if txn.PaymentInstrumentType != "credit_card" {
		t.Fatalf("Returned payment instrument doesn't match input, expected %q, got %q",
			"credit_card", txn.PaymentInstrumentType)
	}
	if txn.CreditCard.CardType != "Discover" {
		t.Fatalf("Returned credit card detail doesn't match input, expected %q, got %q",
			"Visa", txn.CreditCard.CardType)
	}

	txn, err = testGateway.Transaction().SubmitForSettlement(txn.Id, txn.Amount)
	if err != nil {
		t.Fatal(err)
	}

	txn, err = testGateway.Testing().Settle(txn.Id)
	if err != nil {
		t.Fatal(err)
	}

	if txn.Status != TransactionStatusSettled {
		t.Fatal(txn.Status)
	}
}
