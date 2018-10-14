package braintree

import (
	"encoding/xml"
	"testing"
)

var errorXML = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<api-error-response>
  <errors>
    <errors type="array"/>
    <transaction>
      <errors type="array">
        <error>
          <code>81502</code>
          <attribute type="symbol">amount</attribute>
          <message>Amount is required.</message>
        </error>
        <error>
          <code>91526</code>
          <attribute type="symbol">custom_fields</attribute>
          <message>Custom field is invalid: store_me.</message>
        </error>
        <error>
          <code>91513</code>
          <attribute type="symbol">merchant_account_id</attribute>
          <message>Merchant account ID is invalid.</message>
        </error>
      </errors>
      <credit-card>
        <errors type="array">
          <error>
            <code>91708</code>
            <attribute type="symbol">base</attribute>
            <message>Cannot provide expiration_date if you are also providing expiration_month and expiration_year.</message>
          </error>
          <error>
            <code>81714</code>
            <attribute type="symbol">number</attribute>
            <message>Credit card number is required.</message>
          </error>
          <error>
            <code>81725</code>
            <attribute type="symbol">base</attribute>
            <message>Credit card must include either number or venmo_sdk_payment_method_code.</message>
          </error>
          <error>
            <code>81703</code>
            <attribute type="symbol">number</attribute>
            <message>Credit card type is not accepted by this merchant account.</message>
          </error>
        </errors>
      </credit-card>
      <customer>
        <errors type="array">
          <error>
            <code>81606</code>
            <attribute type="symbol">email</attribute>
            <message>Email is an invalid format.</message>
          </error>
        </errors>
      </customer>
    </transaction>
  </errors>
  <message>Everything is broken!</message>
</api-error-response>`)

func TestErrorsUnmarshalEverything(t *testing.T) {
	t.Parallel()

	apiErrors := &BraintreeError{}
	err := xml.Unmarshal(errorXML, apiErrors)
	if err != nil {
		t.Fatal("Error unmarshalling: " + err.Error())
	}

	allErrors := apiErrors.All()

	if len(allErrors) != 8 {
		t.Fatal("Did not get all errors")
	}
}

func TestAccessors(t *testing.T) {
	t.Parallel()

	apiErrors := &BraintreeError{}
	err := xml.Unmarshal(errorXML, apiErrors)
	if err != nil {
		t.Fatal("Error unmarshalling: " + err.Error())
	}

	ccErrors := apiErrors.For("Transaction").On("CreditCard")
	if len(ccErrors) != 4 {
		t.Fatal("Did not get the right credit card errors")
	}

	numberErrors := apiErrors.For("Transaction").For("CreditCard").On("Number")
	if len(numberErrors) != 2 {
		t.Fatal("Did not get the right number errors")
	}

	customerErrors := apiErrors.For("Transaction").On("Customer")
	if len(customerErrors) != 1 {
		t.Fatal("Did not get the right customer errors")
	}

	baseErrors := apiErrors.For("Transaction").On("Base")
	if len(baseErrors) != 3 {
		t.Fatal("Did not get the right base errors")
	}
}
