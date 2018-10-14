package braintree

import (
	"reflect"
	"testing"
)

func TestCustomerApplePayCard(t *testing.T) {
	t.Parallel()

	customer, err := testGateway.Customer().Create(&Customer{})
	if err != nil {
		t.Fatal(err)
	}

	nonce := FakeNonceApplePayVisa

	paymentMethod, err := testGateway.PaymentMethod().Create(&PaymentMethodRequest{
		CustomerId:         customer.Id,
		PaymentMethodNonce: nonce,
	})
	if err != nil {
		t.Fatal(err)
	}
	applePayCard := paymentMethod.(*ApplePayCard)

	customerFound, err := testGateway.Customer().Find(customer.Id)
	if err != nil {
		t.Fatal(err)
	}

	if customerFound.ApplePayCards == nil || len(customerFound.ApplePayCards.ApplePayCard) != 1 {
		t.Fatalf("Customer %#v expected to have one ApplePayCard", customerFound)
	}
	if !reflect.DeepEqual(customerFound.ApplePayCards.ApplePayCard[0], applePayCard) {
		t.Fatalf("Got Customer %#v ApplePayCard %#v, want %#v", customerFound, customerFound.ApplePayCards.ApplePayCard[0], applePayCard)
	}
}
