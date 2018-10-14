package braintree

import "testing"

func TestCredentials(t *testing.T) {
	tests := []struct {
		Credentials                 credentials
		ExpectedEnvironment         Environment
		ExpectedMerchantID          string
		ExpectedAuthorizationHeader string
	}{
		{
			newAPIKey(Development, "development_merchant_id", "integration_public_key", "integration_private_key"),
			Development,
			"development_merchant_id",
			"Basic aW50ZWdyYXRpb25fcHVibGljX2tleTppbnRlZ3JhdGlvbl9wcml2YXRlX2tleQ==",
		},
		{
			func() credentials {
				c, err := newAccessToken("access_token$development$integration_merchant_id$4bff9793ed")
				if err != nil {
					t.Fatalf("Failed to create access token for test: %s", err.Error())
				}
				return c
			}(),
			Development,
			"integration_merchant_id",
			"Bearer access_token$development$integration_merchant_id$4bff9793ed",
		},
	}

	for _, test := range tests {
		if test.Credentials.Environment() != test.ExpectedEnvironment {
			t.Errorf("For %q, got environment %#v, want %#v", test.Credentials, test.Credentials.Environment(), test.ExpectedEnvironment)
		}
		if test.Credentials.MerchantID() != test.ExpectedMerchantID {
			t.Errorf("For %q, got merchant ID %#v, want %#v", test.Credentials, test.Credentials.MerchantID(), test.ExpectedMerchantID)
		}
		if test.Credentials.AuthorizationHeader() != test.ExpectedAuthorizationHeader {
			t.Errorf("For %q, got authorization header %#v, want %#v", test.Credentials, test.Credentials.AuthorizationHeader(), test.ExpectedAuthorizationHeader)
		}
	}
}
