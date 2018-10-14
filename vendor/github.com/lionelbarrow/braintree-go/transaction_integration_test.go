package braintree

import (
	"math/rand"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/lionelbarrow/braintree-go/testhelpers"
)

func randomAmount() *Decimal {
	return NewDecimal(rand.Int63n(10000), 2)
}

func TestTransactionCreateSubmitForSettlementAndVoid(t *testing.T) {
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
	if tx.Id == "" {
		t.Fatal("Received invalid ID on new transaction")
	}
	if tx.Status != TransactionStatusAuthorized {
		t.Fatal(tx.Status)
	}

	// Submit for settlement
	ten := NewDecimal(1000, 2)
	tx2, err := testGateway.Transaction().SubmitForSettlement(tx.Id, ten)

	t.Log(tx2)

	if err != nil {
		t.Fatal(err)
	}
	if x := tx2.Status; x != TransactionStatusSubmittedForSettlement {
		t.Fatal(x)
	}
	if amount := tx2.Amount; amount.Cmp(ten) != 0 {
		t.Fatalf("transaction settlement amount (%s) did not equal amount requested (%s)", amount, ten)
	}

	// Void
	tx3, err := testGateway.Transaction().Void(tx2.Id)

	t.Log(tx3)

	if err != nil {
		t.Fatal(err)
	}
	if x := tx3.Status; x != TransactionStatusVoided {
		t.Fatal(x)
	}
}

func TestTransactionSearch(t *testing.T) {
	t.Parallel()

	txg := testGateway.Transaction()
	createTx := func(amount *Decimal, customerName string) error {
		_, err := txg.Create(&TransactionRequest{
			Type:   "sale",
			Amount: amount,
			Customer: &Customer{
				FirstName: customerName,
			},
			CreditCard: &CreditCard{
				Number:         testCreditCards["visa"].Number,
				ExpirationDate: "05/14",
			},
		})
		return err
	}

	unique := testhelpers.RandomString()

	name0 := "Erik-" + unique
	if err := createTx(randomAmount(), name0); err != nil {
		t.Fatal(err)
	}

	name1 := "Lionel-" + unique
	if err := createTx(randomAmount(), name1); err != nil {
		t.Fatal(err)
	}

	query := new(SearchQuery)
	f := query.AddTextField("customer-first-name")
	f.Is = name0

	result, err := txg.Search(query)
	if err != nil {
		t.Fatal(err)
	}

	if result.TotalItems != 1 {
		t.Fatal(result.Transactions)
	}

	tx := result.Transactions[0]
	if x := tx.Customer.FirstName; x != name0 {
		t.Log(name0)
		t.Fatal(x)
	}
}

func TestTransactionSearchTime(t *testing.T) {
	txg := testGateway.Transaction()
	createTx := func(amount *Decimal, customerName string) error {
		_, err := txg.Create(&TransactionRequest{
			Type:   "sale",
			Amount: amount,
			Customer: &Customer{
				FirstName: customerName,
			},
			CreditCard: &CreditCard{
				Number:         testCreditCards["visa"].Number,
				ExpirationDate: "05/14",
			},
		})
		return err
	}

	unique := testhelpers.RandomString()

	name0 := "Erik-" + unique
	if err := createTx(randomAmount(), name0); err != nil {
		t.Fatal(err)
	}

	name1 := "Lionel-" + unique
	if err := createTx(randomAmount(), name1); err != nil {
		t.Fatal(err)
	}

	{ // test: txn is returned if querying for created at before now
		query := new(SearchQuery)
		f1 := query.AddTextField("customer-first-name")
		f1.Is = name0
		f2 := query.AddTimeField("created-at")
		f2.Max = time.Now()

		result, err := txg.Search(query)
		if err != nil {
			t.Fatal(err)
		}

		if result.TotalItems != 1 {
			t.Fatal(result.Transactions)
		}

		tx := result.Transactions[0]
		if x := tx.Customer.FirstName; x != name0 {
			t.Log(name0)
			t.Fatal(x)
		}
	}

	{ // test: txn is not returned if querying for created at before 1 hour ago
		query := new(SearchQuery)
		f1 := query.AddTextField("customer-first-name")
		f1.Is = name0
		f2 := query.AddTimeField("created-at")
		f2.Max = time.Now().Add(-time.Hour)

		result, err := txg.Search(query)
		if err != nil {
			t.Fatal(err)
		}

		if result.TotalItems != 0 {
			t.Fatal(result.Transactions)
		}
	}
}

// This test will fail unless you set up your Braintree sandbox account correctly. See TESTING.md for details.
func TestTransactionCreateWhenGatewayRejected(t *testing.T) {
	t.Parallel()

	_, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:   "sale",
		Amount: NewDecimal(201000, 2),
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
		},
	})
	if err == nil {
		t.Fatal("Did not receive error when creating invalid transaction")
	}
	if err.Error() != "Card Issuer Declined CVV" {
		t.Fatal(err)
	}
	if err.(*BraintreeError).Transaction.ProcessorResponseCode != 2010 {
		t.Fatalf("expected err.Transaction.ProcessorResponseCode to be 2010, but got %d", err.(*BraintreeError).Transaction.ProcessorResponseCode)
	}

	if err.(*BraintreeError).Transaction.AdditionalProcessorResponse != "2010 : Card Issuer Declined CVV" {
		t.Fatalf("expected err.Transaction.ProcessorResponseCode to be `2010 : Card Issuer Declined CVV`, but got %s", err.(*BraintreeError).Transaction.AdditionalProcessorResponse)
	}
}

func TestTransactionCreateWhenGatewayRejectedFraud(t *testing.T) {
	t.Parallel()

	_, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:               "sale",
		Amount:             NewDecimal(201000, 2),
		PaymentMethodNonce: FakeNonceGatewayRejectedFraud,
	})
	if err == nil {
		t.Fatal("Did not receive error when creating invalid transaction")
	}

	if err.Error() != "Gateway Rejected: fraud" {
		t.Fatal(err)
	}

	txn := err.(*BraintreeError).Transaction
	if txn.Status != TransactionStatusGatewayRejected {
		t.Fatalf("Got status %q, want %q", txn.Status, TransactionStatusGatewayRejected)
	}

	if txn.GatewayRejectionReason != GatewayRejectionReasonFraud {
		t.Fatalf("Got gateway rejection reason %q, wanted %q", txn.GatewayRejectionReason, GatewayRejectionReasonFraud)
	}

	if txn.ProcessorResponseCode != 0 {
		t.Fatalf("Got processor response code %q, want %q", txn.ProcessorResponseCode, 0)
	}
}

func TestTransactionCreatedWhenCVVDoesNotMatch(t *testing.T) {
	t.Parallel()

	_, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:   "sale",
		Amount: randomAmount(),
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
			CVV:            "200", // Should cause CVV does not match response
		},
	})

	if err.Error() != "Gateway Rejected: cvv" {
		t.Fatal(err)
	}

	txn := err.(*BraintreeError).Transaction

	if txn.Status != TransactionStatusGatewayRejected {
		t.Fatalf("Got status %q, want %q", txn.Status, TransactionStatusGatewayRejected)
	}

	if txn.GatewayRejectionReason != GatewayRejectionReasonCVV {
		t.Fatalf("Got gateway rejection reason %q, wanted %q", txn.GatewayRejectionReason, GatewayRejectionReasonCVV)
	}

	if txn.CVVResponseCode != CVVResponseCodeDoesNotMatch {
		t.Fatalf("Got CVV Response Code %q, wanted %q", txn.CVVResponseCode, CVVResponseCodeDoesNotMatch)
	}
}

func TestTransactionCreatedWhenAVSBankDoesNotSupport(t *testing.T) {
	t.Parallel()

	_, err := testGateway.Transaction().Create(&TransactionRequest{
		MerchantAccountId: avsAndCVVTestMerchantAccountId,
		Type:              "sale",
		Amount:            randomAmount(),
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
			CVV:            "100",
		},
		BillingAddress: &Address{
			StreetAddress: "1 E Main St",
			Locality:      "Chicago",
			Region:        "IL",
			PostalCode:    "30001", // Should cause AVS bank does not support error response.
		},
	})

	if err == nil {
		t.Fatal("Did not receive error when creating invalid transaction")
	}

	if err.Error() != "Gateway Rejected: avs" {
		t.Fatal(err)
	}

	txn := err.(*BraintreeError).Transaction

	if txn.Status != TransactionStatusGatewayRejected {
		t.Fatalf("Got status %q, want %q", txn.Status, TransactionStatusGatewayRejected)
	}

	if txn.GatewayRejectionReason != GatewayRejectionReasonAVS {
		t.Fatalf("Got gateway rejection reason %q, wanted %q", txn.GatewayRejectionReason, GatewayRejectionReasonAVS)
	}

	if txn.AVSErrorResponseCode != AVSResponseCodeNotSupported {
		t.Fatalf("Got AVS Error Response Code %q, wanted %q", txn.AVSErrorResponseCode, AVSResponseCodeNotSupported)
	}
}

func TestTransactionCreatedWhenAVSPostalDoesNotMatch(t *testing.T) {
	t.Parallel()

	_, err := testGateway.Transaction().Create(&TransactionRequest{
		MerchantAccountId: avsAndCVVTestMerchantAccountId,
		Type:              "sale",
		Amount:            randomAmount(),
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
			CVV:            "100",
		},
		BillingAddress: &Address{
			StreetAddress: "1 E Main St",
			Locality:      "Chicago",
			Region:        "IL",
			PostalCode:    "20000", // Should cause AVS postal code does not match response.
		},
	})

	if err == nil {
		t.Fatal("Did not receive error when creating invalid transaction")
	}

	if err.Error() != "Gateway Rejected: avs" {
		t.Fatal(err)
	}

	txn := err.(*BraintreeError).Transaction

	if txn.Status != TransactionStatusGatewayRejected {
		t.Fatalf("Got status %q, want %q", txn.Status, TransactionStatusGatewayRejected)
	}

	if txn.GatewayRejectionReason != GatewayRejectionReasonAVS {
		t.Fatalf("Got gateway rejection reason %q, wanted %q", txn.GatewayRejectionReason, GatewayRejectionReasonAVS)
	}

	if txn.AVSPostalCodeResponseCode != AVSResponseCodeDoesNotMatch {
		t.Fatalf("Got AVS postal response code %q, wanted %q", txn.AVSPostalCodeResponseCode, AVSResponseCodeDoesNotMatch)
	}
}

func TestTransactionCreatedWhenAVStreetAddressDoesNotMatch(t *testing.T) {
	t.Parallel()

	_, err := testGateway.Transaction().Create(&TransactionRequest{
		MerchantAccountId: avsAndCVVTestMerchantAccountId,
		Type:              "sale",
		Amount:            randomAmount(),
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
			CVV:            "100",
		},
		BillingAddress: &Address{
			StreetAddress: "201 E Main St", // Should cause AVS street address not verified response.
			Locality:      "Chicago",
			Region:        "IL",
			PostalCode:    "60637",
		},
	})

	if err == nil {
		t.Fatal("Did not receive error when creating invalid transaction")
	}

	if err.Error() != "Gateway Rejected: avs" {
		t.Fatal(err)
	}

	txn := err.(*BraintreeError).Transaction

	if txn.Status != TransactionStatusGatewayRejected {
		t.Fatalf("Got status %q, want %q", txn.Status, TransactionStatusGatewayRejected)
	}

	if txn.GatewayRejectionReason != GatewayRejectionReasonAVS {
		t.Fatalf("Got gateway rejection reason %q, wanted %q", txn.GatewayRejectionReason, GatewayRejectionReasonAVS)
	}

	if txn.AVSStreetAddressResponseCode != AVSResponseCodeNotVerified {
		t.Fatalf("Got AVS street address response code %q, wanted %q", txn.AVSStreetAddressResponseCode, AVSResponseCodeNotVerified)
	}
}

func TestFindTransaction(t *testing.T) {
	t.Parallel()

	createdTransaction, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:   "sale",
		Amount: randomAmount(),
		CreditCard: &CreditCard{
			Number:         testCreditCards["mastercard"].Number,
			ExpirationDate: "05/14",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	foundTransaction, err := testGateway.Transaction().Find(createdTransaction.Id)
	if err != nil {
		t.Fatal(err)
	}

	if createdTransaction.Id != foundTransaction.Id {
		t.Fatal("transaction ids do not match")
	}
}

func TestFindNonExistantTransaction(t *testing.T) {
	t.Parallel()

	_, err := testGateway.Transaction().Find("bad_transaction_id")
	if err == nil {
		t.Fatal("Did not receive error when finding an invalid tx ID")
	}
	if err.Error() != "Not Found (404)" {
		t.Fatal(err)
	}
	if apiErr, ok := err.(APIError); !(ok && apiErr.StatusCode() == http.StatusNotFound) {
		t.Fatal(err)
	}
}

// This test will fail unless you set up your Braintree sandbox account correctly. See TESTING.md for details.
func TestTransactionDescriptorFields(t *testing.T) {
	t.Parallel()

	tx := &TransactionRequest{
		Type:               "sale",
		Amount:             randomAmount(),
		PaymentMethodNonce: FakeNonceTransactable,
		Options: &TransactionOptions{
			SubmitForSettlement: true,
		},
		Descriptor: &Descriptor{
			Name:  "Company Name*Product 1",
			Phone: "0000000000",
			URL:   "example.com",
		},
	}

	tx2, err := testGateway.Transaction().Create(tx)
	if err != nil {
		t.Fatal(err)
	}

	if tx2.Type != tx.Type {
		t.Fatalf("expected Type to be equal, but %s was not %s", tx2.Type, tx.Type)
	}
	if tx2.Amount.Cmp(tx.Amount) != 0 {
		t.Fatalf("expected Amount to be equal, but %s was not %s", tx2.Amount, tx.Amount)
	}
	if tx2.Status != TransactionStatusSubmittedForSettlement {
		t.Fatalf("expected tx2.Status to be %s, but got %s", TransactionStatusSubmittedForSettlement, tx2.Status)
	}
	if tx2.Descriptor.Name != "Company Name*Product 1" {
		t.Fatalf("expected tx2.Descriptor.Name to be Company Name*Product 1, but got %s", tx2.Descriptor.Name)
	}
	if tx2.Descriptor.Phone != "0000000000" {
		t.Fatalf("expected tx2.Descriptor.Phone to be 0000000000, but got %s", tx2.Descriptor.Phone)
	}
	if tx2.Descriptor.URL != "example.com" {
		t.Fatalf("expected tx2.Descriptor.URL to be example.com, but got %s", tx2.Descriptor.URL)
	}
}

// This test will fail unless you set up your Braintree sandbox account correctly. See TESTING.md for details.
func TestTransactionPaypalFields(t *testing.T) {
	t.Parallel()

	const (
		PayeeEmail  = "payee@payal.com"
		Description = "One tasty sandwich"
		CustomField = "foo"
	)
	subData := make(map[string]string)
	subData["faz"] = "bar"

	customer, err := testGateway.Customer().Create(&Customer{})
	if err != nil {
		t.Fatal(err)
	}

	nonce := FakeNoncePayPalFuturePayment

	paymentMethod, err := testGateway.PaymentMethod().Create(&PaymentMethodRequest{
		CustomerId:         customer.Id,
		PaymentMethodNonce: nonce,
	})
	if err != nil {
		t.Fatal(err)
	}
	paypalAccount, ok := paymentMethod.(*PayPalAccount)
	if !ok {
		t.Fatal("Could not assert paypal account")
	}

	tx := &TransactionRequest{
		Type:               "sale",
		Amount:             randomAmount(),
		PaymentMethodToken: paypalAccount.GetToken(),
		OrderId:            "123456ABC",
		Options: &TransactionOptions{
			SubmitForSettlement: true,
			TransactionOptionsPaypalRequest: &TransactionOptionsPaypalRequest{
				PayeeEmail:        PayeeEmail,
				Description:       Description,
				CustomField:       CustomField,
				SupplementaryData: subData,
			},
		},
	}
	tx2, err := testGateway.Transaction().Create(tx)
	if err != nil {
		t.Fatal(err)
	}
	if tx2.Type != tx.Type {
		t.Fatalf("expected Type to be equal, but %s was not %s", tx2.Type, tx.Type)
	}
	if tx2.Amount.Cmp(tx.Amount) != 0 {
		t.Fatalf("expected Amount to be equal, but %s was not %s", tx2.Amount, tx.Amount)
	}
	if tx2.Status != TransactionStatusSettling {
		t.Fatalf("expected tx2.Status to be %s, but got %s", TransactionStatusSettling, tx2.Status)
	}
	if tx2.PayPalDetails.PayeeEmail != PayeeEmail {
		t.Fatalf("expected tx2.PaypalDetails.PayeeEmail to be %s, but got %s", PayeeEmail, tx2.PayPalDetails.PayeeEmail)
	}
	if tx2.PayPalDetails.Description != Description {
		t.Fatalf("expected tx2.PaypalDetails.Description to be %s, but got %s", Description, tx2.PayPalDetails.Description)
	}
	if tx2.PayPalDetails.CustomField != CustomField {
		t.Fatalf("expected tx2.PayPalDetails.CustomField to be %s, but got %s", CustomField, tx2.PayPalDetails.CustomField)
	}
}

func TestTransactionRiskDataFields(t *testing.T) {
	t.Parallel()

	tx := &TransactionRequest{
		Type:               "sale",
		Amount:             randomAmount(),
		PaymentMethodNonce: FakeNonceTransactable,
		RiskData: &RiskDataRequest{
			CustomerBrowser: "Mozilla/5.0 (X11; U; Linux x86_64; en-US) AppleWebKit/540.0 (KHTML,like Gecko) Chrome/9.1.0.0 Safari/540.0",
			CustomerIP:      "127.0.0.1",
		},
	}

	tx2, err := testGateway.Transaction().Create(tx)
	if err != nil {
		t.Fatal(err)
	}

	if tx2.Type != tx.Type {
		t.Fatalf("expected Type to be equal, but %s was not %s", tx2.Type, tx.Type)
	}
	if tx2.Amount.Cmp(tx.Amount) != 0 {
		t.Fatalf("expected Amount to be equal, but %s was not %s", tx2.Amount, tx.Amount)
	}
}

func TestAllTransactionFields(t *testing.T) {
	t.Parallel()

	amount := randomAmount()
	taxAmount := NewDecimal(amount.Unscaled/10, amount.Scale)

	tx := &TransactionRequest{
		Type:    "sale",
		Amount:  amount,
		OrderId: "my_custom_order",
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
			CVV:            "100",
		},
		Customer: &Customer{
			FirstName: "Lionel",
		},
		BillingAddress: &Address{
			StreetAddress: "1 E Main St",
			Locality:      "Chicago",
			Region:        "IL",
			PostalCode:    "60637",
		},
		ShippingAddress: &Address{
			StreetAddress: "1 E Main St",
			Locality:      "Chicago",
			Region:        "IL",
			PostalCode:    "60637",
		},
		TaxAmount:  taxAmount,
		DeviceData: `{"device_session_id": "dsi_1234", "fraud_merchant_id": "fmi_1234", "correlation_id": "ci_1234"}`,
		Channel:    "ChannelA",
		Options: &TransactionOptions{
			SubmitForSettlement:              true,
			StoreInVault:                     true,
			AddBillingAddressToPaymentMethod: true,
			StoreShippingAddressInVault:      true,
		},
	}

	tx2, err := testGateway.Transaction().Create(tx)
	if err != nil {
		t.Fatal(err)
	}

	if tx2.Type != tx.Type {
		t.Fatalf("expected Type to be equal, but %s was not %s", tx2.Type, tx.Type)
	}
	if tx2.CurrencyISOCode != "USD" {
		t.Fatalf("expected CurrencyISOCode to be %s but was %s", "USD", tx2.CurrencyISOCode)
	}
	if tx2.Amount.Cmp(tx.Amount) != 0 {
		t.Fatalf("expected Amount to be equal, but %s was not %s", tx2.Amount, tx.Amount)
	}
	if tx2.OrderId != tx.OrderId {
		t.Fatalf("expected OrderId to be equal, but %s was not %s", tx2.OrderId, tx.OrderId)
	}
	if tx2.Customer.FirstName != tx.Customer.FirstName {
		t.Fatalf("expected Customer.FirstName to be equal, but %s was not %s", tx2.Customer.FirstName, tx.Customer.FirstName)
	}
	if tx2.BillingAddress.StreetAddress != tx.BillingAddress.StreetAddress {
		t.Fatalf("expected BillingAddress.StreetAddress to be equal, but %s was not %s", tx2.BillingAddress.StreetAddress, tx.BillingAddress.StreetAddress)
	}
	if tx2.BillingAddress.Locality != tx.BillingAddress.Locality {
		t.Fatalf("expected BillingAddress.Locality to be equal, but %s was not %s", tx2.BillingAddress.Locality, tx.BillingAddress.Locality)
	}
	if tx2.BillingAddress.Region != tx.BillingAddress.Region {
		t.Fatalf("expected BillingAddress.Region to be equal, but %s was not %s", tx2.BillingAddress.Region, tx.BillingAddress.Region)
	}
	if tx2.BillingAddress.PostalCode != tx.BillingAddress.PostalCode {
		t.Fatalf("expected BillingAddress.PostalCode to be equal, but %s was not %s", tx2.BillingAddress.PostalCode, tx.BillingAddress.PostalCode)
	}
	if tx2.ShippingAddress.StreetAddress != tx.ShippingAddress.StreetAddress {
		t.Fatalf("expected ShippingAddress.StreetAddress to be equal, but %s was not %s", tx2.ShippingAddress.StreetAddress, tx.ShippingAddress.StreetAddress)
	}
	if tx2.ShippingAddress.Locality != tx.ShippingAddress.Locality {
		t.Fatalf("expected ShippingAddress.Locality to be equal, but %s was not %s", tx2.ShippingAddress.Locality, tx.ShippingAddress.Locality)
	}
	if tx2.ShippingAddress.Region != tx.ShippingAddress.Region {
		t.Fatalf("expected ShippingAddress.Region to be equal, but %s was not %s", tx2.ShippingAddress.Region, tx.ShippingAddress.Region)
	}
	if tx2.ShippingAddress.PostalCode != tx.ShippingAddress.PostalCode {
		t.Fatalf("expected ShippingAddress.PostalCode to be equal, but %s was not %s", tx2.ShippingAddress.PostalCode, tx.ShippingAddress.PostalCode)
	}
	if tx2.TaxAmount == nil {
		t.Fatalf("expected TaxAmount to be set, but was nil")
	}
	if tx2.TaxAmount.Cmp(tx.TaxAmount) != 0 {
		t.Fatalf("expected TaxAmount to be equal, but %s was not %s", tx2.TaxAmount, tx.TaxAmount)
	}
	if tx2.TaxExempt != tx.TaxExempt {
		t.Fatalf("expected TaxExempt to be equal, but %t was not %t", tx2.TaxExempt, tx.TaxExempt)
	}
	if tx2.CreditCard.Token == "" {
		t.Fatalf("expected CreditCard.Token to be equal, but %s was not %s", tx2.CreditCard.Token, tx.CreditCard.Token)
	}
	if tx2.Customer.Id == "" {
		t.Fatalf("expected Customer.Id to be equal, but %s was not %s", tx2.Customer.Id, tx.Customer.Id)
	}
	if tx2.Status != TransactionStatusSubmittedForSettlement {
		t.Fatalf("expected tx2.Status to be %s, but got %s", TransactionStatusSubmittedForSettlement, tx2.Status)
	}
	if tx2.PaymentInstrumentType != "credit_card" {
		t.Fatalf("expected tx2.PaymentInstrumentType to be %s, but got %s", "credit_card", tx2.PaymentInstrumentType)
	}
	if tx2.AdditionalProcessorResponse != "" {
		t.Fatalf("expected tx2.AdditionalProcessorResponse to be empty, but got %s", tx2.AdditionalProcessorResponse)
	}
	if tx2.RiskData == nil {
		t.Fatal("expected tx2.RiskData not to be empty")
	}
	if tx2.RiskData.ID == "" {
		t.Fatal("expected tx2.RiskData.ID not to be empty")
	}
	if tx2.RiskData.Decision != "Approve" {
		t.Fatalf("expected tx2.RiskData.Decision to be Approve, but got %s", tx2.RiskData.Decision)
	}
	if tx2.Channel != "ChannelA" {
		t.Fatalf("expected tx2.Channel to be ChannelA, but got %s", tx2.Channel)
	}
}

// This test will only pass on Travis. See TESTING.md for more details.
func TestTransactionDisbursementDetails(t *testing.T) {
	t.Parallel()

	txn, err := testGateway.Transaction().Find("dskdmb")
	if err != nil {
		t.Fatal(err)
	}

	if txn.DisbursementDetails.DisbursementDate != "2013-06-27" {
		t.Fatalf("expected disbursement date to be %s, was %s", "2013-06-27", txn.DisbursementDetails.DisbursementDate)
	}
	if txn.DisbursementDetails.SettlementAmount.Cmp(NewDecimal(10000, 2)) != 0 {
		t.Fatalf("expected settlement amount to be %s, was %s", NewDecimal(10000, 2), txn.DisbursementDetails.SettlementAmount)
	}
	if txn.DisbursementDetails.SettlementCurrencyIsoCode != "USD" {
		t.Fatalf("expected settlement currency code to be %s, was %s", "USD", txn.DisbursementDetails.SettlementCurrencyIsoCode)
	}
	if txn.DisbursementDetails.SettlementCurrencyExchangeRate.Cmp(NewDecimal(100, 2)) != 0 {
		t.Fatalf("expected settlement currency exchange rate to be %s, was %s", NewDecimal(100, 2), txn.DisbursementDetails.SettlementCurrencyExchangeRate)
	}
	if txn.DisbursementDetails.FundsHeld {
		t.Error("funds held doesn't match")
	}
	if !txn.DisbursementDetails.Success {
		t.Error("success doesn't match")
	}
}

func TestTransactionCreateFromPaymentMethodCode(t *testing.T) {
	t.Parallel()

	customer, err := testGateway.Customer().Create(&Customer{
		CreditCard: &CreditCard{
			Number:         testCreditCards["discover"].Number,
			ExpirationDate: "05/14",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if customer.CreditCards.CreditCard[0].Token == "" {
		t.Fatal("invalid token")
	}

	tx, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:               "sale",
		CustomerID:         customer.Id,
		Amount:             randomAmount(),
		PaymentMethodToken: customer.CreditCards.CreditCard[0].Token,
	})

	if err != nil {
		t.Fatal(err)
	}
	if tx.Id == "" {
		t.Fatal("invalid tx id")
	}
}

func TestTrxPaymentMethodNonce(t *testing.T) {
	t.Parallel()

	txn, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:               "sale",
		Amount:             randomAmount(),
		PaymentMethodNonce: "fake-apple-pay-mastercard-nonce",
	})
	if err != nil {
		t.Fatal(err)
	}

	txn, err = testGateway.Transaction().SubmitForSettlement(txn.Id, txn.Amount)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTransactionCreateSettleAndFullRefund(t *testing.T) {
	t.Parallel()

	amount := NewDecimal(20000, 2)
	txn, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:   "sale",
		Amount: amount,
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

	txn, err = testGateway.Testing().Settle(txn.Id)
	if err != nil {
		t.Fatal(err)
	}

	if txn.Status != TransactionStatusSettled {
		t.Fatal(txn.Status)
	}

	// Refund
	refundTxn, err := testGateway.Transaction().Refund(txn.Id)

	t.Log(refundTxn)

	if err != nil {
		t.Fatal(err)
	}
	if x := refundTxn.Status; x != TransactionStatusSubmittedForSettlement {
		t.Fatal(x)
	}

	refundTxn, err = testGateway.Testing().Settle(refundTxn.Id)
	if err != nil {
		t.Fatal(err)
	}

	if refundTxn.Status != TransactionStatusSettled {
		t.Fatal(txn.Status)
	}

	if *refundTxn.RefundedTransactionId != txn.Id {
		t.Fatal(*refundTxn.RefundedTransactionId)
	}

	// Check that the refund shows up in the original transaction
	txn, err = testGateway.Transaction().Find(txn.Id)
	if err != nil {
		t.Fatal(err)
	}

	if txn.RefundIds != nil && (*txn.RefundIds)[0] != refundTxn.Id {
		t.Fatal(*txn.RefundIds)
	}

	// Second refund should fail
	refundTxn, err = testGateway.Transaction().Refund(txn.Id)
	t.Log(refundTxn)

	if err.Error() != "Transaction has already been completely refunded." {
		t.Fatal(err)
	}
}

func TestTransactionCreateSettleAndPartialRefund(t *testing.T) {
	t.Parallel()

	amount := NewDecimal(10000, 2)
	refundAmt1 := NewDecimal(5000, 2)
	refundAmt2 := NewDecimal(5001, 2)
	txn, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:   "sale",
		Amount: amount,
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

	txn, err = testGateway.Testing().Settle(txn.Id)
	if err != nil {
		t.Fatal(err)
	}

	if txn.Status != TransactionStatusSettled {
		t.Fatal(txn.Status)
	}

	// Refund
	refundTxn, err := testGateway.Transaction().Refund(txn.Id, refundAmt1)

	t.Log(refundTxn)

	if err != nil {
		t.Fatal(err)
	}
	if x := refundTxn.Status; x != TransactionStatusSubmittedForSettlement {
		t.Fatal(x)
	}

	refundTxn, err = testGateway.Testing().Settle(refundTxn.Id)
	if err != nil {
		t.Fatal(err)
	}

	if refundTxn.Status != TransactionStatusSettled {
		t.Fatal(txn.Status)
	}

	// Refund amount too large
	refundTxn, err = testGateway.Transaction().Refund(txn.Id, refundAmt2)

	t.Log(refundTxn)

	if err.Error() != "Refund amount is too large." {
		t.Fatal(err)
	}
}

func TestTransactionCreateWithCustomFields(t *testing.T) {
	t.Parallel()

	customFields := map[string]string{
		"custom_field_1": "custom value",
	}

	amount := NewDecimal(10000, 2)
	txn, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:               "sale",
		Amount:             amount,
		PaymentMethodNonce: FakeNonceTransactable,
		CustomFields:       customFields,
	})
	if err != nil {
		t.Fatal(err)
	}

	if x := map[string]string(txn.CustomFields); !reflect.DeepEqual(x, customFields) {
		t.Fatalf("Returned custom fields doesn't match input, got %q, want %q", x, customFields)
	}

	txn, err = testGateway.Transaction().Find(txn.Id)
	if err != nil {
		t.Fatal(err)
	}

	if x := map[string]string(txn.CustomFields); !reflect.DeepEqual(x, customFields) {
		t.Fatalf("Returned custom fields doesn't match input, got %q, want %q", x, customFields)
	}
}

func TestTransactionTaxExempt(t *testing.T) {
	t.Parallel()

	amount := NewDecimal(10000, 2)
	txn, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:               "sale",
		Amount:             amount,
		TaxExempt:          true,
		PaymentMethodNonce: FakeNonceTransactable,
	})
	if err != nil {
		t.Fatal(err)
	}

	txn, err = testGateway.Transaction().Find(txn.Id)
	if err != nil {
		t.Fatal(err)
	}

	if !txn.TaxExempt {
		t.Fatalf("Transaction did not return tax exempt")
	}
	if txn.TaxAmount != nil {
		t.Fatalf("Transaction TaxAmount got %v, want nil", txn.TaxAmount)
	}
}

func TestTransactionTaxFieldsNotProvided(t *testing.T) {
	t.Parallel()

	amount := NewDecimal(10000, 2)
	txn, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:               "sale",
		Amount:             amount,
		PaymentMethodNonce: FakeNonceTransactable,
	})
	if err != nil {
		t.Fatal(err)
	}

	txn, err = testGateway.Transaction().Find(txn.Id)
	if err != nil {
		t.Fatal(err)
	}

	if txn.TaxExempt {
		t.Fatalf("Transaction returned tax exempt, expected not to")
	}
	if txn.TaxAmount != nil {
		t.Fatalf("Transaction tax amount got %v, want nil", *txn.TaxAmount)
	}
}

func TestEscrowHoldOnCreate(t *testing.T) {
	t.Parallel()

	txn, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:   "sale",
		Amount: NewDecimal(6200, 2),
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
		},
		MerchantAccountId: testSubMerchantAccount(),
		ServiceFeeAmount:  NewDecimal(1000, 2),
		Options: &TransactionOptions{
			HoldInEscrow: true,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if txn.EscrowStatus != EscrowStatusHoldPending {
		t.Fatalf("Transaction EscrowStatus got %s, want %s", txn.EscrowStatus, EscrowStatusHoldPending)
	}
}

func TestEscrowHoldOnCreateOnMasterMerchant(t *testing.T) {
	t.Parallel()

	_, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:   "sale",
		Amount: NewDecimal(6301, 2),
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
		},
		Options: &TransactionOptions{
			HoldInEscrow: true,
		},
	})
	if err == nil {
		t.Fatal("Transaction Sale got no error, want error")
	}
	errors := err.(*BraintreeError).Errors.TransactionErrors.For("Base").On("base")
	if len(errors) != 1 {
		t.Fatalf("Transaction Sale got %d errors, want 1 error", len(errors))
	}
	if g, w := errors[0].Code, "91560"; g != w {
		t.Errorf("Transaction Sale got error code %s, want %s", g, w)
	}
	if g, w := errors[0].Message, "Transaction could not be held in escrow."; g != w {
		t.Errorf("Transaction Sale got error message %s, want %s", g, w)
	}
}

func TestEscrowHoldAfterSale(t *testing.T) {
	t.Parallel()

	txn, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:   "sale",
		Amount: NewDecimal(6300, 2),
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
		},
		MerchantAccountId: testSubMerchantAccount(),
		ServiceFeeAmount:  NewDecimal(1000, 2),
	})
	if err != nil {
		t.Fatal(err)
	}
	txn, err = testGateway.Transaction().HoldInEscrow(txn.Id)
	if err != nil {
		t.Fatal(err)
	}
	if txn.EscrowStatus != EscrowStatusHoldPending {
		t.Fatalf("Transaction EscrowStatus got %s, want %s", txn.EscrowStatus, EscrowStatusHoldPending)
	}
}

func TestEscrowHoldAfterSaleOnMasterMerchant(t *testing.T) {
	t.Parallel()

	txn, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:   "sale",
		Amount: NewDecimal(6301, 2),
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = testGateway.Transaction().HoldInEscrow(txn.Id)
	if err == nil {
		t.Fatal("Transaction HoldInEscrow got no error, want error")
	}
	errors := err.(*BraintreeError).Errors.TransactionErrors.For("Base").On("base")
	if len(errors) != 1 {
		t.Fatalf("Transaction HoldInEscrow got %d errors, want 1 error", len(errors))
	}
	if g, w := errors[0].Code, "91560"; g != w {
		t.Errorf("Transaction HoldInEscrow got error code %s, want %s", g, w)
	}
	if g, w := errors[0].Message, "Transaction could not be held in escrow."; g != w {
		t.Errorf("Transaction HoldInEscrow got error message %s, want %s", g, w)
	}
}

func TestEscrowRelease(t *testing.T) {
	t.Parallel()

	txn, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:   "sale",
		Amount: NewDecimal(6400, 2),
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
		},
		MerchantAccountId: testSubMerchantAccount(),
		ServiceFeeAmount:  NewDecimal(1000, 2),
		Options: &TransactionOptions{
			SubmitForSettlement: true,
			HoldInEscrow:        true,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	txn, err = testGateway.Transaction().Settle(txn.Id)
	if err != nil {
		t.Fatal(err)
	}
	txn, err = testGateway.Transaction().ReleaseFromEscrow(txn.Id)
	if err != nil {
		t.Fatal(err)
	}
	if txn.EscrowStatus != EscrowStatusReleasePending {
		t.Fatalf("Transaction EscrowStatus got %s, want %s", txn.EscrowStatus, EscrowStatusReleasePending)
	}
}

func TestEscrowReleaseNotEscrowed(t *testing.T) {
	t.Parallel()

	txn, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:   "sale",
		Amount: NewDecimal(6401, 2),
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
		},
		MerchantAccountId: testSubMerchantAccount(),
		ServiceFeeAmount:  NewDecimal(1000, 2),
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = testGateway.Transaction().ReleaseFromEscrow(txn.Id)
	if err == nil {
		t.Fatal("Transaction ReleaseFromEscrow got no error, want error")
	}
	errors := err.(*BraintreeError).Errors.TransactionErrors.For("Base").On("base")
	if len(errors) != 1 {
		t.Fatalf("Transaction ReleaseFromEscrow got %d errors, want 1 error", len(errors))
	}
	if g, w := errors[0].Code, "91561"; g != w {
		t.Errorf("Transaction ReleaseFromEscrow got error code %s, want %s", g, w)
	}
	if g, w := errors[0].Message, "Cannot release a transaction that is not escrowed."; g != w {
		t.Errorf("Transaction ReleaseFromEscrow got error message %s, want %s", g, w)
	}
}

func TestEscrowCancelRelease(t *testing.T) {
	t.Parallel()

	txn, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:   "sale",
		Amount: NewDecimal(6500, 2),
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
		},
		MerchantAccountId: testSubMerchantAccount(),
		ServiceFeeAmount:  NewDecimal(1000, 2),
		Options: &TransactionOptions{
			SubmitForSettlement: true,
			HoldInEscrow:        true,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	txn, err = testGateway.Transaction().Settle(txn.Id)
	if err != nil {
		t.Fatal(err)
	}
	txn, err = testGateway.Transaction().ReleaseFromEscrow(txn.Id)
	if err != nil {
		t.Fatal(err)
	}
	if txn.EscrowStatus != EscrowStatusReleasePending {
		t.Fatalf("Transaction EscrowStatus got %s, want %s", txn.EscrowStatus, EscrowStatusReleasePending)
	}
	txn, err = testGateway.Transaction().CancelRelease(txn.Id)
	if err != nil {
		t.Fatal(err)
	}
	if txn.EscrowStatus != EscrowStatusHeld {
		t.Fatalf("Transaction EscrowStatus got %s, want %s", txn.EscrowStatus, EscrowStatusHeld)
	}
}

func TestEscrowCancelReleaseNotPending(t *testing.T) {
	t.Parallel()

	txn, err := testGateway.Transaction().Create(&TransactionRequest{
		Type:   "sale",
		Amount: NewDecimal(6501, 2),
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
		},
		MerchantAccountId: testSubMerchantAccount(),
		ServiceFeeAmount:  NewDecimal(1000, 2),
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = testGateway.Transaction().CancelRelease(txn.Id)
	if err == nil {
		t.Fatal("Transaction Cancel Release got no error, want error")
	}
	errors := err.(*BraintreeError).Errors.TransactionErrors.For("Base").On("base")
	if len(errors) != 1 {
		t.Fatalf("Transaction Cancel Release got %d errors, want 1 error", len(errors))
	}
	if g, w := errors[0].Code, "91562"; g != w {
		t.Errorf("Transaction Cancel Release got error code %s, want %s", g, w)
	}
	if g, w := errors[0].Message, "Release can only be cancelled if the transaction is submitted for release."; g != w {
		t.Errorf("Transaction Cancel Release got error message %s, want %s", g, w)
	}
}
