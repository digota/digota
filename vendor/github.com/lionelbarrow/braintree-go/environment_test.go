package braintree

import "testing"

func TestEnvironmentBaseURL(t *testing.T) {
	tests := []struct {
		Environment   Environment
		WantedBaseURL string
	}{
		{Development, "http://localhost:3000"},
		{Sandbox, "https://api.sandbox.braintreegateway.com:443"},
		{Production, "https://api.braintreegateway.com:443"},
	}

	for _, test := range tests {
		actual := test.Environment.BaseURL()
		if actual != test.WantedBaseURL {
			t.Fatalf("%#v.BaseURL() = %#v, want %#v", test.Environment, actual, test.WantedBaseURL)
		}
	}
}
