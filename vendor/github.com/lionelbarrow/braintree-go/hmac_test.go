package braintree

import (
	"testing"
)

func TestHmacerParseSignature(t *testing.T) {
	t.Parallel()

	hmacer := newHmacer("pubkey", "privkey")

	testCases := []struct {
		Description       string
		SignatureKeyPairs string
		ExpectedSignature string
		ExpectedError     error
	}{
		{
			Description:       "Single signature",
			SignatureKeyPairs: "pubkey|the_signature",
			ExpectedSignature: "the_signature",
		},
		{
			Description:       "Single signature (not matching)",
			SignatureKeyPairs: "pubkey2|the_signature",
			ExpectedError:     SignatureError{message: "Signature-key pair contains the wrong public key!"},
		},
		{
			Description:       "Multiple signatures (one matching)",
			SignatureKeyPairs: "pubkey|the_signature&pubkey2|the_signature",
			ExpectedSignature: "the_signature",
		},
		{
			Description:       "Invalid signature (no pipe)",
			SignatureKeyPairs: "pubkeythe_signature",
			ExpectedError:     SignatureError{message: "Signature-key pair does not contain |"},
		},
	}

	for _, tc := range testCases {
		signature, err := hmacer.parseSignature(tc.SignatureKeyPairs)
		if err != nil && (tc.ExpectedError == nil || err.(SignatureError) != tc.ExpectedError) {
			t.Errorf("Test Case %q with Signature %q encountered error %#v want %#v", tc.Description, tc.SignatureKeyPairs, err, tc.ExpectedError)
		} else if signature != tc.ExpectedSignature {
			t.Errorf("Test Case %q with Signature %q returned %#v want %#v", tc.Description, tc.SignatureKeyPairs, signature, tc.ExpectedSignature)
		}
	}
}

func TestHmacerVerifySignature(t *testing.T) {
	t.Parallel()

	hmacer := newHmacer("jkq28pcxj4r85dwr", "66062a3876e2dc298f2195f0bf173f5a")

	testCases := []struct {
		Description   string
		Signature     string
		Payload       string
		ExpectedError error
		ExpectedValid bool
	}{
		{
			Description:   "Single signature (valid)",
			Signature:     "jkq28pcxj4r85dwr|4af78bab15cc58195871c636c786716f34cd9711",
			Payload:       "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPG5vdGlm\naWNhdGlvbj4KICA8a2luZD5jaGVjazwva2luZD4KICA8dGltZXN0YW1wIHR5\ncGU9ImRhdGV0aW1lIj4yMDE3LTA0LTI0VDA0OjI1OjEwWjwvdGltZXN0YW1w\nPgogIDxzdWJqZWN0PgogICAgPGNoZWNrIHR5cGU9ImJvb2xlYW4iPnRydWU8\nL2NoZWNrPgogIDwvc3ViamVjdD4KPC9ub3RpZmljYXRpb24+Cg==\n",
			ExpectedValid: true,
		},
		{
			Description:   "Multiple signature (valid)",
			Signature:     "4zn8jg4gdmzyvcyd|dd6390bc9d75985f0cc986d5d5f55fcdb52531cb&cd7jwvrw8jytyfm3|d7fdd777e30a1fd93b58770d7682b577b461cf6f&jkq28pcxj4r85dwr|96dd50905f51f6de1c24790c4d77aa460cb55a3d",
			Payload:       "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPG5vdGlm\naWNhdGlvbj4KICA8a2luZD5jaGVjazwva2luZD4KICA8dGltZXN0YW1wIHR5\ncGU9ImRhdGV0aW1lIj4yMDE3LTA0LTI0VDA0OjUwOjA0WjwvdGltZXN0YW1w\nPgogIDxzdWJqZWN0PgogICAgPGNoZWNrIHR5cGU9ImJvb2xlYW4iPnRydWU8\nL2NoZWNrPgogIDwvc3ViamVjdD4KPC9ub3RpZmljYXRpb24+Cg==\n",
			ExpectedValid: true,
		},
		{
			Description:   "Single signature (invalid)",
			Signature:     "jkq28pcxj4r85dwr|4af78bab15cc58195871c636c786716f34cd9711",
			Payload:       "payloadthatdoesntmatchsignature",
			ExpectedValid: false,
		},
		{
			Description:   "Single signature (unknown public key)",
			Signature:     "cd7jwvrw8jytyfm3|d7fdd777e30a1fd93b58770d7682b577b461cf6f&jkq28pcxj4r85dwr",
			Payload:       "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPG5vdGlm\naWNhdGlvbj4KICA8a2luZD5jaGVjazwva2luZD4KICA8dGltZXN0YW1wIHR5\ncGU9ImRhdGV0aW1lIj4yMDE3LTA0LTI0VDA0OjI1OjEwWjwvdGltZXN0YW1w\nPgogIDxzdWJqZWN0PgogICAgPGNoZWNrIHR5cGU9ImJvb2xlYW4iPnRydWU8\nL2NoZWNrPgogIDwvc3ViamVjdD4KPC9ub3RpZmljYXRpb24+Cg==\n",
			ExpectedError: SignatureError{message: "Signature-key pair contains the wrong public key!"},
		},
	}

	for _, tc := range testCases {
		valid, err := hmacer.verifySignature(tc.Signature, tc.Payload)
		if err != nil && (tc.ExpectedError == nil || err.(SignatureError) != tc.ExpectedError) {
			t.Errorf("Test Case %q with Signature %q and Payload %q encountered error %#v want %#v", tc.Description, tc.Signature, tc.Payload, err, tc.ExpectedError)
		} else if valid != tc.ExpectedValid {
			t.Errorf("Test Case %q with Signature %q and Payload %q was valid %#v want %#v", tc.Description, tc.Signature, tc.Payload, valid, tc.ExpectedValid)
		}
	}
}
