package braintree

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"text/template"
)

const payloadTemplate = `
<notification>
	<timestamp type="datetime">%s</timestamp>
	<kind>%s</kind>
	<subject>%s</subject>
</notification>
`

// WebhookTestingGateway exports actions only available in the sandbox environment.
type WebhookTestingGateway struct {
	*Braintree
	apiKey apiKey
}

// Request simulates an incoming webhook notification request
func (g *WebhookTestingGateway) Request(kind, id string) (*http.Request, error) {
	payload := g.SamplePayload(kind, id)
	signature, err := g.SignPayload(payload)
	if err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Add("bt_signature", signature)
	form.Add("bt_payload", payload)

	body := form.Encode()
	return &http.Request{
		Method:        "POST",
		Header:        http.Header{"Content-Type": {"application/x-www-form-urlencoded"}},
		ContentLength: int64(len(body)),
		Body:          ioutil.NopCloser(strings.NewReader(body)),
	}, nil
}

// SamplePayload builds a sample payload of the given kind, different kind of
// notifications are defined as constants in webhook_notification.go
// id represents an id used in XML subject of the specified kind
func (g *WebhookTestingGateway) SamplePayload(kind, id string) string {
	datetime := time.Now().UTC().Format(time.RFC3339)
	payload := fmt.Sprintf(
		payloadTemplate,
		datetime,
		kind,
		g.subjectXML(kind, id),
	)
	encodedpayload := base64.StdEncoding.EncodeToString([]byte(payload))
	return strings.Replace(encodedpayload, "\r", "", -1)
}

// SignPayload signs payload and returns HMAC signature
func (g *WebhookTestingGateway) SignPayload(payload string) (string, error) {
	hmacer := newHmacer(g.apiKey.publicKey, g.apiKey.privateKey)
	hmacedPayload, err := hmacer.hmac(payload)
	if err != nil {
		return "", err
	}
	return hmacer.publicKey + "|" + hmacedPayload, nil
}

func (g *WebhookTestingGateway) subjectXML(kind, id string) string {
	type subjectXMLData struct {
		ID string
	}

	xmlTmpl := g.subjectXMLTemplate(kind)
	var b bytes.Buffer
	tmpl, err := template.New("").Parse(xmlTmpl)
	err = tmpl.Execute(&b, subjectXMLData{ID: id})
	if err != nil {
		panic(fmt.Errorf("creating xml template: " + err.Error()))
	}
	return b.String()
}

func (g *WebhookTestingGateway) subjectXMLTemplate(kind string) string {
	switch kind {
	case CheckWebhook:
		return g.checkXML()
	case SubMerchantAccountApprovedWebhook:
		return g.merchantAccountXMLApproved()
	case SubMerchantAccountDeclinedWebhook:
		return g.merchantAccountXMLDeclined()
	case TransactionDisbursedWebhook:
		return g.transactionDisbursedXML()
	case TransactionSettledWebhook:
		return g.transactionSettledXML()
	case TransactionSettlementDeclinedWebhook:
		return g.transactionSettlementDeclinedXML()
	case DisbursementWebhook:
		return g.disbursementXML()
	case DisputeOpenedWebhook:
		return g.disputeOpenedXML()
	case DisputeLostWebhook:
		return g.disputeLostXML()
	case DisputeWonWebhook:
		return g.disputeWonXML()
	case DisbursementExceptionWebhook:
		return g.disbursementExceptionXML()
	case PartnerMerchantConnectedWebhook:
		return g.partnerMerchantConnectedXML()
	case PartnerMerchantDisconnectedWebhook:
		return g.partnerMerchantDisconnectedXML()
	case PartnerMerchantDeclinedWebhook:
		return g.partnerMerchantDeclinedXML()
	case SubscriptionChargedSuccessfullyWebhook:
		return g.subscriptionChargedSuccessfullyXML()
	case AccountUpdaterDailyReportWebhook:
		return g.accountUpdaterDailyReportXML()
	default:
		return g.subscriptionXML()
	}
}

func (g *WebhookTestingGateway) checkXML() string {
	return `<check type="boolean">true</check>`
}

func (g *WebhookTestingGateway) merchantAccountXMLApproved() string {
	return `
		<merchant_account>
			<id>{{ $.ID }}</id>
			<master_merchant_account>
				<id>master_ma_for_{{ $.ID }}</id>
				<status>active</status>
			</master_merchant_account>
			<status>active</status>
		</merchant_account>
		`
}

func (g *WebhookTestingGateway) merchantAccountXMLDeclined() string {
	return `
		<api-error-response>
			<message>Credit score is too low</message>
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
				<id>{{ $.ID }}</id>
				<status>suspended</status>
				<master-merchant-account>
					<id>master_ma_for_{{ $.ID }}</id>
					<status>suspended</status>
				</master-merchant-account>
			</merchant-account>
		</api-error-response>
		`
}

func (g *WebhookTestingGateway) subscriptionXML() string {
	return `
		<subscription>
			<id>{{ $.ID }}</id>
			<transactions type="array">
			</transactions>
			<add_ons type="array">
			</add_ons>
			<discounts type="array">
			</discounts>
		</subscription>
		`
}

func (g *WebhookTestingGateway) subscriptionChargedSuccessfullyXML() string {
	return `
		<subscription>
			<id>{{ $.ID }}</id>
			<transactions type="array">
				<transaction>
					<id>{{ $.ID }}</id>
					<status>submitted_for_settlement</status>
					<amount>49.99</amount>
				</transaction>
			</transactions>
			<add_ons type="array">
			</add_ons>
			<discounts type="array">
			</discounts>
		</subscription>
		`
}

func (g *WebhookTestingGateway) transactionDisbursedXML() string {
	return `
		<transaction>
			<id>{{ $.ID }}</id>
			<amount>100</amount>
			<disbursement-details>
				<disbursement-date type="date">2013-07-09</disbursement-date>
			</disbursement-details>
		</transaction>
		`
}

func (g *WebhookTestingGateway) transactionSettledXML() string {
	return `
		<transaction>
			<id>{{ $.ID}}</id>
			<status>settled</status>
			<type>sale</type>
			<currency-iso-code>USD</currency-iso-code>
			<amount>100.00</amount>
			<merchant-account-id>ogaotkivejpfayqfeaimuktty</merchant-account-id>
			<payment-instrument-type>us_bank_account</payment-instrument-type>
			<us-bank-account>
				<routing-number>123456789</routing-number>
				<last-4>1234</last-4>
				<account-type>checking</account-type>
				<account-holder-name>Dan Schulman</account-holder-name>
			</us-bank-account>
		</transaction>
		`
}

func (g *WebhookTestingGateway) transactionSettlementDeclinedXML() string {
	return `
		<transaction>
			<id>{{ $.ID }}</id>
			<status>settlement_declined</status>
			<type>sale</type>
			<currency-iso-code>USD</currency-iso-code>
			<amount>100.00</amount>
			<merchant-account-id>ogaotkivejpfayqfeaimuktty</merchant-account-id>
			<payment-instrument-type>us_bank_account</payment-instrument-type>
			<us-bank-account>
				<routing-number>123456789</routing-number>
				<last-4>1234</last-4>
				<account-type>checking</account-type>
				<account-holder-name>Dan Schulman</account-holder-name>
			</us-bank-account>
		</transaction>
		`
}

func (g *WebhookTestingGateway) disputeOpenedXML() string {
	return `
		<dispute>
			<amount>250.00</amount>
			<currency-iso-code>USD</currency-iso-code>
			<received-date type="date">2014-03-01</received-date>
			<reply-by-date type="date">2014-03-21</reply-by-date>
			<kind>chargeback</kind>
			<status>open</status>
			<reason>fraud</reason>
			<id>{{ $.ID }}</id>
			<transaction>
				<id>{{ $.ID }}</id>
				<amount>250.00</amount>
			</transaction>
			<date-opened type="date">2014-03-21</date-opened>
		</dispute>
		`
}

func (g *WebhookTestingGateway) disputeLostXML() string {
	return `
		<dispute>
			<amount>250.00</amount>
			<currency-iso-code>USD</currency-iso-code>
			<received-date type="date">2014-03-01</received-date>
			<reply-by-date type="date">2014-03-21</reply-by-date>
			<kind>chargeback</kind>
			<status>lost</status>
			<reason>fraud</reason>
			<id>{{ $.ID }}</id>
			<transaction>
				<id>{{ $.ID }}</id>
				<amount>250.00</amount>
			</transaction>
			<date-opened type="date">2014-03-21</date-opened>
		</dispute>
		`
}

func (g *WebhookTestingGateway) disputeWonXML() string {
	return `
		<dispute>
			<amount>250.00</amount>
			<currency-iso-code>USD</currency-iso-code>
			<received-date type="date">2014-03-01</received-date>
			<reply-by-date type="date">2014-03-21</reply-by-date>
			<kind>chargeback</kind>
			<status>won</status>
			<reason>fraud</reason>
			<id>{{ $.ID }}</id>
			<transaction>
				<id>{{ $.ID }}</id>
				<amount>250.00</amount>
			</transaction>
			<date-opened type="date">2014-03-21</date-opened>
			<date-won type="date">2014-03-22</date-won>
		</dispute>
		`
}

func (g *WebhookTestingGateway) disbursementXML() string {
	return `
		<disbursement>
			<id>{{ $.ID }}</id>
			<transaction-ids type="array">
				<item>afv56j</item>
				<item>kj8hjk</item>
			</transaction-ids>
			<success type="boolean">true</success>
			<retry type="boolean">false</retry>
			<merchant-account>
				<id>merchant_account_token</id>
				<currency-iso-code>USD</currency-iso-code>
				<sub-merchant-account type="boolean">false</sub-merchant-account>
				<status>active</status>
			</merchant-account>
			<amount>100.00</amount>
			<disbursement-date type="date">2014-02-10</disbursement-date>
			<exception-message nil="true"/>
			<follow-up-action nil="true"/>
		</disbursement>
		`
}

func (g *WebhookTestingGateway) disbursementExceptionXML() string {
	return `
		<disbursement>
			<id>{{ $.ID }}</id>
			<transaction-ids type="array">
				<item>afv56j</item>
				<item>kj8hjk</item>
			</transaction-ids>
			<success type="boolean">false</success>
			<retry type="boolean">false</retry>
			<merchant-account>
				<id>merchant_account_token</id>
				<currency-iso-code>USD</currency-iso-code>
				<sub-merchant-account type="boolean">false</sub-merchant-account>
				<status>active</status>
			</merchant-account>
			<amount>100.00</amount>
			<disbursement-date type="date">2014-02-10</disbursement-date>
			<exception-message>bank_rejected</exception-message>
			<follow-up-action>update_funding_information</follow-up-action>
		</disbursement>
	`
}

func (g *WebhookTestingGateway) partnerMerchantConnectedXML() string {
	return `
	<partner-merchant>
		<merchant-public-id>public_id</merchant-public-id>
		<public-key>public_key</public-key>
		<private-key>private_key</private-key>
		<partner-merchant-id>abc123</partner-merchant-id>
		<client-side-encryption-key>cse_key</client-side-encryption-key>
	</partner-merchant>
	`
}

func (g *WebhookTestingGateway) partnerMerchantDisconnectedXML() string {
	return `
	<partner-merchant>
		<partner-merchant-id>abc123</partner-merchant-id>
	</partner-merchant>
	`
}

func (g *WebhookTestingGateway) partnerMerchantDeclinedXML() string {
	return `
	<partner-merchant>
		<partner-merchant-id>abc123</partner-merchant-id>
	</partner-merchant>
	`
}

func (g *WebhookTestingGateway) accountUpdaterDailyReportXML() string {
	return `
	<account-updater-daily-report>
		<report-date type="date">2016-01-14</report-date>
		<report-url>link-to-csv-report</report-url>
	</account-updater-daily-report>
	`
}
