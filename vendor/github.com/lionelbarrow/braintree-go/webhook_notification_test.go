package braintree

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestWebhookParseRequest(t *testing.T) {
	t.Parallel()

	gateway := New(Sandbox, "mid", "sz9g7zhxz8838v7h", "0c809a2d2e8f4e4c817900ff441c9554")
	webhookGateway := gateway.WebhookNotification()

	body := strings.NewReader("bt_signature=sz9g7zhxz8838v7h%7C4b532339b3107eae876d7637d59217858f320098&bt_payload=PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPG5vdGlm%0AaWNhdGlvbj4KICA8a2luZD5jaGVjazwva2luZD4KICA8dGltZXN0YW1wIHR5%0AcGU9ImRhdGV0aW1lIj4yMDE3LTA0LTI2VDA3OjEyOjI0WjwvdGltZXN0YW1w%0APgogIDxzdWJqZWN0PgogICAgPGNoZWNrIHR5cGU9ImJvb2xlYW4iPnRydWU8%0AL2NoZWNrPgogIDwvc3ViamVjdD4KPC9ub3RpZmljYXRpb24%2BCg%3D%3D%0A")
	r := &http.Request{
		Method:        "POST",
		Header:        http.Header{"Content-Type": {"application/x-www-form-urlencoded"}},
		ContentLength: int64(body.Len()),
		Body:          ioutil.NopCloser(body),
	}

	notification, err := webhookGateway.ParseRequest(r)

	if err != nil {
		t.Fatal(err)
	} else if notification.Kind != CheckWebhook {
		t.Fatal("Incorrect Notification kind, expected check got", notification.Kind)
	}
}

func TestWebhookParseMerchantAccountAccepted(t *testing.T) {
	t.Parallel()

	webhookGateway := testGateway.WebhookNotification()
	apiKey := testGateway.credentials.(apiKey)
	hmacer := newHmacer(apiKey.publicKey, apiKey.privateKey)

	payload := base64.StdEncoding.EncodeToString([]byte(`
<notification>
  <timestamp type="datetime">2014-01-26T10:32:28+00:00</timestamp>
  <kind>sub_merchant_account_approved</kind>
  <subject>
    <merchant-account>
      <id>123</id>
      <master-merchant-account>
        <id>master_ma_for_123</id>
        <status>active</status>
      </master-merchant-account>
      <status>active</status>
    </merchant-account>
  </subject>
</notification>`))
	hmacedPayload, err := hmacer.hmac(payload)
	if err != nil {
		t.Fatal(err)
	}
	signature := hmacer.publicKey + "|" + hmacedPayload

	notification, err := webhookGateway.Parse(signature, payload)

	if err != nil {
		t.Fatal(err)
	} else if notification.Kind != SubMerchantAccountApprovedWebhook {
		t.Fatal("Incorrect Notification kind, expected sub_merchant_account_approved got", notification.Kind)
	} else if notification.MerchantAccount() == nil {
		t.Log(notification.Subject)
		t.Fatal("Notification should have a merchant account")
	} else if notification.MerchantAccount().Id != "123" {
		t.Fatal("Incorrect Merchant Id, expected '123' got", notification.Subject.MerchantAccount.Id)
	} else if notification.MerchantAccount().Status != "active" {
		t.Fatal("Incorrect Merchant Status, expected 'active' got", notification.Subject.MerchantAccount.Status)
	}
}

func TestWebhookParseMerchantAccountDeclined(t *testing.T) {
	t.Parallel()

	webhookGateway := testGateway.WebhookNotification()
	apiKey := testGateway.credentials.(apiKey)
	hmacer := newHmacer(apiKey.publicKey, apiKey.privateKey)

	payload := base64.StdEncoding.EncodeToString([]byte(`
<notification>
   <timestamp type="datetime">2014-01-26T10:32:28+00:00</timestamp>
   <kind>sub_merchant_account_declined</kind>
   <subject>
     <api-error-response>
       <message>Credit score is too low</message>
       <errors>
         <errors type="array"/>
           <merchant-account>
             <errors type="array">
               <error>
                 <code>82621</code>
                 <message>Credit score is too low</message>
                 <attribute type="symbol">base</attribute>
               </error>
             </errors>
           </merchant-account>
         </errors>
         <merchant-account>
           <id>123</id>
           <status>suspended</status>
           <master-merchant-account>
             <id>master_ma_for_123</id>
             <status>suspended</status>
           </master-merchant-account>
         </merchant-account>
     </api-error-response>
  </subject>
</notification>`))
	hmacedPayload, err := hmacer.hmac(payload)
	if err != nil {
		t.Fatal(err)
	}
	signature := hmacer.publicKey + "|" + hmacedPayload

	notification, err := webhookGateway.Parse(signature, payload)

	if err != nil {
		t.Fatal(err)
	} else if notification.Kind != SubMerchantAccountDeclinedWebhook {
		t.Fatal("Incorrect Notification kind, expected sub_merchant_account_declined got", notification.Kind)
	} else if notification.Subject.APIErrorResponse == nil {
		t.Fatal("Notification should have an error response")
	} else if notification.Subject.APIErrorResponse.ErrorMessage != "Credit score is too low" {
		t.Fatal("Incorrect Error Message, expected 'Credit score is too low' got", notification.Subject.APIErrorResponse.ErrorMessage)
	} else if notification.MerchantAccount() == nil {
		t.Fatal("Notification should have a merchant account")
	} else if notification.MerchantAccount().Id != "123" {
		t.Fatal("Incorrect Merchant Id, expected '123' got", notification.Subject.MerchantAccount.Id)
	} else if notification.MerchantAccount().Status != "suspended" {
		t.Fatal("Incorrect Merchant Status, expected 'suspended' got", notification.Subject.MerchantAccount.Status)
	}
}

func TestWebhookParseDisbursement(t *testing.T) {
	t.Parallel()

	webhookGateway := testGateway.WebhookNotification()
	apiKey := testGateway.credentials.(apiKey)
	hmacer := newHmacer(apiKey.publicKey, apiKey.privateKey)

	payload := base64.StdEncoding.EncodeToString([]byte(`
<notification>
	<timestamp type="datetime">2014-04-06T10:32:28+00:00</timestamp>
	<kind>disbursement</kind>
	<subject>
		<disbursement>
			<id>456</id>
			<transaction-ids type="array">
				<item>afv56j</item>
				<item>kj8hjk</item>
			</transaction-ids>
			<success type="boolean">true</success>
			<retry type="boolean">false</retry>
			<merchant-account>
				<id>123</id>
				<currency-iso-code>USD</currency-iso-code>
				<sub-merchant-account type="boolean">false</sub-merchant-account>
				<status>active</status>
			</merchant-account>
			<amount>100.00</amount>
			<disbursement-date type="date">2014-02-09</disbursement-date>
			<exception-message nil="true"/>
			<follow-up-action nil="true"/>
		</disbursement>
	</subject>
</notification>`))
	hmacedPayload, err := hmacer.hmac(payload)
	if err != nil {
		t.Fatal(err)
	}
	signature := hmacer.publicKey + "|" + hmacedPayload

	notification, err := webhookGateway.Parse(signature, payload)

	if err != nil {
		t.Fatal(err)
	} else if notification.Kind != DisbursementWebhook {
		t.Fatal("Incorrect Notification kind, expected disbursement got", notification.Kind)
	} else if notification.Disbursement() == nil {
		t.Fatal("Notification should have a disbursement")
	} else if notification.Disbursement().Id != "456" {
		t.Fatal("Incorrect disbursement id, expected 456 got", notification.Subject.MerchantAccount.Status)
	} else if len(notification.Disbursement().TransactionIds) != 2 {
		t.Fatal("Disbursement should have two txns")
	} else if notification.Disbursement().TransactionIds[1] != "kj8hjk" {
		t.Fatal("Incorrect txn id on disbursement, expected kj8hjk got", notification.Disbursement().TransactionIds[1])
	} else if notification.Disbursement().MerchantAccount.Id != "123" {
		t.Fatal("Disbursement not associated with correct merchant account")
	} else if notification.Disbursement().ExceptionMessage != "" {
		t.Fatal("Disbursement should not have an exception message")
	}
}

func TestWebhookParseDisbursementException(t *testing.T) {
	t.Parallel()

	webhookGateway := testGateway.WebhookNotification()
	apiKey := testGateway.credentials.(apiKey)
	hmacer := newHmacer(apiKey.publicKey, apiKey.privateKey)

	payload := base64.StdEncoding.EncodeToString([]byte(`
<notification>
	<timestamp type="datetime">2014-04-06T10:32:28+00:00</timestamp>
	<kind>disbursement_exception</kind>
	<subject>
		<disbursement>
			<id>456</id>
			<transaction-ids type="array">
				<item>afv56j</item>
				<item>kj8hjk</item>
			</transaction-ids>
			<success type="boolean">false</success>
			<retry type="boolean">false</retry>
			<merchant-account>
				<id>123</id>
				<currency-iso-code>USD</currency-iso-code>
				<sub-merchant-account type="boolean">false</sub-merchant-account>
				<status>active</status>
			</merchant-account>
			<amount>100.00</amount>
			<disbursement-date type="date">2014-02-09</disbursement-date>
      <exception-message>bank_rejected</exception-message>
      <follow-up-action>update_funding_information</follow-up-action>
		</disbursement>
	</subject>
</notification>`))
	hmacedPayload, err := hmacer.hmac(payload)
	if err != nil {
		t.Fatal(err)
	}
	signature := hmacer.publicKey + "|" + hmacedPayload

	notification, err := webhookGateway.Parse(signature, payload)

	if err != nil {
		t.Fatal(err)
	} else if notification.Disbursement().ExceptionMessage != BankRejected {
		t.Fatal("Disbursement should have a BankRejected exception message")
	} else if notification.Disbursement().FollowUpAction != UpdateFundingInformation {
		t.Fatal("Disbursement followup action should be UpdateFundingInformation")
	}

}
