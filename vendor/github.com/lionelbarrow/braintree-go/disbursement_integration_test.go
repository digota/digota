package braintree

import (
	"testing"
)

// This test will fail unless you have a transaction with this ID on your sandbox.
func TestDisbursementTransactions(t *testing.T) {
	t.Parallel()

	d := Disbursement{
		TransactionIds: []string{"dskdmb"},
	}

	result, err := d.Transactions(testGateway.Transaction())

	if err != nil {
		t.Fatal(err)
	}

	if result.TotalItems != 1 {
		t.Fatal(result)
	}

	txn := result.Transactions[0]
	if txn.Id != "dskdmb" {
		t.Fatal(txn.Id)
	}

}
