package braintree

import "testing"

func TestTransactionClone(t *testing.T) {
	t.Parallel()

	tx, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:   "sale",
		Amount: NewDecimal(2000, 2),
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
		},
	})
	t.Log(tx)
	if err != nil {
		t.Fatal(err)
	}

	// Clone
	tx2, err := testGateway.Transaction().Clone(tx.Id, &TransactionCloneRequest{
		Amount:  NewDecimal(1000, 2),
		Channel: "ChannelA",
		Options: &TransactionCloneOptions{
			SubmitForSettlement: false,
		},
	})
	t.Log(tx2)
	if err != nil {
		t.Fatal(err)
	}
	if g, w := tx2.Status, TransactionStatusAuthorized; g != w {
		t.Errorf("Transaction status got %v, want %v", g, w)
	}
	if g, w := tx2.Amount, NewDecimal(1000, 2); g.Cmp(w) != 0 {
		t.Errorf("Transaction amount got %v, want %v", g, w)
	}
	if g, w := tx2.Channel, "ChannelA"; g != w {
		t.Errorf("Transaction channel got %v, want %v", g, w)
	}
	if g, w := tx2.CreditCard.ExpirationMonth, "05"; g != w {
		t.Errorf("Transaction credit card expiration month got %v, want %v", g, w)
	}
	if g, w := tx2.CreditCard.ExpirationYear, "2014"; g != w {
		t.Errorf("Transaction credit card expiration year got %v, want %v", g, w)
	}
	if g, w := tx2.CreditCard.Last4, "1111"; g != w {
		t.Errorf("Transaction credit card last 4 got %v, want %v", g, w)
	}
}

func TestTransactionCloneSubmittedForSettlement(t *testing.T) {
	t.Parallel()

	tx, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:   "sale",
		Amount: NewDecimal(2000, 2),
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
		},
	})
	t.Log(tx)
	if err != nil {
		t.Fatal(err)
	}

	// Clone
	tx2, err := testGateway.Transaction().Clone(tx.Id, &TransactionCloneRequest{
		Amount:  NewDecimal(1000, 2),
		Channel: "ChannelA",
		Options: &TransactionCloneOptions{
			SubmitForSettlement: true,
		},
	})
	t.Log(tx2)
	if err != nil {
		t.Fatal(err)
	}
	if g, w := tx2.Status, TransactionStatusSubmittedForSettlement; g != w {
		t.Errorf("Transaction status got %v, want %v", g, w)
	}
	if g, w := tx2.Amount, NewDecimal(1000, 2); g.Cmp(w) != 0 {
		t.Errorf("Transaction amount got %v, want %v", g, w)
	}
	if g, w := tx2.Channel, "ChannelA"; g != w {
		t.Errorf("Transaction channel got %v, want %v", g, w)
	}
	if g, w := tx2.CreditCard.ExpirationMonth, "05"; g != w {
		t.Errorf("Transaction credit card expiration month got %v, want %v", g, w)
	}
	if g, w := tx2.CreditCard.ExpirationYear, "2014"; g != w {
		t.Errorf("Transaction credit card expiration year got %v, want %v", g, w)
	}
	if g, w := tx2.CreditCard.Last4, "1111"; g != w {
		t.Errorf("Transaction credit card last 4 got %v, want %v", g, w)
	}
}
