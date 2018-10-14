package braintree

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/lionelbarrow/braintree-go/testhelpers"
)

// This test will fail unless you set up your Braintree sandbox account correctly. See TESTING.md for details.
func TestCustomer(t *testing.T) {
	t.Parallel()

	oc := &Customer{
		FirstName: "Lionel",
		LastName:  "Barrow",
		Company:   "Braintree",
		Email:     "lionel.barrow@example.com",
		Phone:     "312.555.1234",
		Fax:       "614.555.5678",
		Website:   "http://www.example.com",
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
			CVV:            "200",
			Options: &CreditCardOptions{
				VerifyCard: testhelpers.BoolPtr(true),
			},
		},
	}

	// Create with errors
	_, err := testGateway.Customer().Create(oc)
	if err == nil {
		t.Fatal("Did not receive error when creating invalid customer")
	}

	// Create
	oc.CreditCard.CVV = ""
	oc.CreditCard.Options = nil
	customer, err := testGateway.Customer().Create(oc)

	t.Log(customer)

	if err != nil {
		t.Fatal(err)
	}
	if customer.Id == "" {
		t.Fatal("invalid customer id")
	}
	if card := customer.DefaultCreditCard(); card == nil {
		t.Fatal("invalid credit card")
	}
	if card := customer.DefaultCreditCard(); card.Token == "" {
		t.Fatal("invalid token")
	}

	// Update
	unique := testhelpers.RandomString()
	newFirstName := "John" + unique
	c2, err := testGateway.Customer().Update(&Customer{
		Id:        customer.Id,
		FirstName: newFirstName,
	})

	t.Log(c2)

	if err != nil {
		t.Fatal(err)
	}
	if c2.FirstName != newFirstName {
		t.Fatal("first name not changed")
	}

	// Find
	c3, err := testGateway.Customer().Find(customer.Id)

	t.Log(c3)

	if err != nil {
		t.Fatal(err)
	}
	if c3.Id != customer.Id {
		t.Fatal("ids do not match")
	}

	// Search
	query := new(SearchQuery)
	f := query.AddTextField("first-name")
	f.Is = newFirstName
	searchResult, err := testGateway.Customer().Search(query)
	if err != nil {
		t.Fatal(err)
	}
	if len(searchResult.Customers) == 0 {
		t.Fatal("could not search for a customer")
	}
	if id := searchResult.Customers[0].Id; id != customer.Id {
		t.Fatalf("id from search does not match: got %s, wanted %s", id, customer.Id)
	}

	// Delete
	err = testGateway.Customer().Delete(customer.Id)
	if err != nil {
		t.Fatal(err)
	}

	// Test customer 404
	c4, err := testGateway.Customer().Find(customer.Id)
	if err == nil {
		t.Fatal("should return 404")
	}
	if err.Error() != "Not Found (404)" {
		t.Fatal(err)
	}
	if apiErr, ok := err.(APIError); !(ok && apiErr.StatusCode() == http.StatusNotFound) {
		t.Fatal(err)
	}
	if c4 != nil {
		t.Fatal(c4)
	}
}

func TestCustomerWithCustomFields(t *testing.T) {
	t.Parallel()

	customFields := map[string]string{
		"custom_field_1": "custom value",
	}

	c := &Customer{
		CustomFields: customFields,
	}

	customer, err := testGateway.Customer().Create(c)
	if err != nil {
		t.Fatal(err)
	}

	if x := map[string]string(customer.CustomFields); !reflect.DeepEqual(x, customFields) {
		t.Fatalf("Returned custom fields doesn't match input, got %q, want %q", x, customFields)
	}

	customer, err = testGateway.Customer().Find(customer.Id)
	if err != nil {
		t.Fatal(err)
	}

	if x := map[string]string(customer.CustomFields); !reflect.DeepEqual(x, customFields) {
		t.Fatalf("Returned custom fields doesn't match input, got %q, want %q", x, customFields)
	}
}

func TestCustomerPaymentMethods(t *testing.T) {
	t.Parallel()

	customer, err := testGateway.Customer().Create(&Customer{})
	if err != nil {
		t.Fatal(err)
	}

	paymentMethod1, err := testGateway.PaymentMethod().Create(&PaymentMethodRequest{
		CustomerId:         customer.Id,
		PaymentMethodNonce: FakeNoncePayPalBillingAgreement,
	})
	if err != nil {
		t.Fatal(err)
	}
	paymentMethod2, err := testGateway.PaymentMethod().Create(&PaymentMethodRequest{
		CustomerId:         customer.Id,
		PaymentMethodNonce: FakeNonceTransactable,
	})
	if err != nil {
		t.Fatal(err)
	}

	expectedPaymentMethods := []PaymentMethod{
		paymentMethod2,
		paymentMethod1,
	}

	customerFound, err := testGateway.Customer().Find(customer.Id)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(customerFound.PaymentMethods(), expectedPaymentMethods) {
		t.Fatalf("Got Customer %#v PaymentMethods %#v, want %#v", customerFound, customerFound.PaymentMethods(), expectedPaymentMethods)
	}
}

func TestCustomerDefaultPaymentMethod(t *testing.T) {
	t.Parallel()

	customer, err := testGateway.Customer().Create(&Customer{})
	if err != nil {
		t.Fatal(err)
	}

	defaultPaymentMethod, err := testGateway.PaymentMethod().Create(&PaymentMethodRequest{
		CustomerId:         customer.Id,
		PaymentMethodNonce: FakeNonceTransactable,
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = testGateway.PaymentMethod().Create(&PaymentMethodRequest{
		CustomerId:         customer.Id,
		PaymentMethodNonce: FakeNoncePayPalBillingAgreement,
	})
	if err != nil {
		t.Fatal(err)
	}

	customerFound, err := testGateway.Customer().Find(customer.Id)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(customerFound.DefaultPaymentMethod(), defaultPaymentMethod) {
		t.Fatalf("Got Customer %#v DefaultPaymentMethod %#v, want %#v", customerFound, customerFound.DefaultPaymentMethod(), defaultPaymentMethod)
	}
}

func TestCustomerDefaultPaymentMethodManuallySet(t *testing.T) {
	t.Parallel()

	customer, err := testGateway.Customer().Create(&Customer{})
	if err != nil {
		t.Fatal(err)
	}

	_, err = testGateway.PaymentMethod().Create(&PaymentMethodRequest{
		CustomerId:         customer.Id,
		PaymentMethodNonce: FakeNonceTransactable,
	})
	if err != nil {
		t.Fatal(err)
	}
	paymentMethod2, err := testGateway.PaymentMethod().Create(&PaymentMethodRequest{
		CustomerId:         customer.Id,
		PaymentMethodNonce: FakeNoncePayPalBillingAgreement,
	})
	if err != nil {
		t.Fatal(err)
	}
	paypalAccount, err := testGateway.PayPalAccount().Update(&PayPalAccount{
		Token: paymentMethod2.GetToken(),
		Options: &PayPalAccountOptions{
			MakeDefault: true,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	customerFound, err := testGateway.Customer().Find(customer.Id)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(customerFound.DefaultPaymentMethod(), paypalAccount) {
		t.Fatalf("Got Customer %#v DefaultPaymentMethod %#v, want %#v", customerFound, customerFound.DefaultPaymentMethod(), paypalAccount)
	}
}

func TestCustomerPaymentMethodNonce(t *testing.T) {
	t.Parallel()

	customer, err := testGateway.Customer().Create(&Customer{PaymentMethodNonce: FakeNonceTransactable})
	if err != nil {
		t.Fatal(err)
	}

	customerFound, err := testGateway.Customer().Find(customer.Id)
	if err != nil {
		t.Fatal(err)
	}

	if len(customer.PaymentMethods()) != 1 {
		t.Fatalf("Customer %#v has %#v payment method(s), want 1 payment method", customerFound, len(customer.PaymentMethods()))
	}
}
