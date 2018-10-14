package braintree

import "testing"

// This test will fail unless you set up your Braintree sandbox account correctly. See TESTING.md for details.
func TestClientToken(t *testing.T) {
	t.Parallel()

	g := testGateway.ClientToken()
	token, err := g.Generate()
	if err != nil {
		t.Fatalf("failed to generate client token: %s", err)
	}
	if len(token) == 0 {
		t.Fatalf("empty client token!")
	}
}

func TestClientTokenWithCustomer(t *testing.T) {
	t.Parallel()

	customerRequest := &Customer{FirstName: "Lionel"}

	customer, err := testGateway.Customer().Create(customerRequest)
	if err != nil {
		t.Error(err)
	}

	customerId := customer.Id

	token, err := testGateway.ClientToken().GenerateWithCustomer(customerId)
	if err != nil {
		t.Error(err)
	} else if len(token) == 0 {
		t.Fatalf("Received empty client token")
	}
}
