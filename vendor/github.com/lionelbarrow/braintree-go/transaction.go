package braintree

import (
	"encoding/xml"
	"time"

	"github.com/lionelbarrow/braintree-go/customfields"
)

type TransactionStatus string

const (
	TransactionStatusAuthorizationExpired   TransactionStatus = "authorization_expired"
	TransactionStatusAuthorizing            TransactionStatus = "authorizing"
	TransactionStatusAuthorized             TransactionStatus = "authorized"
	TransactionStatusGatewayRejected        TransactionStatus = "gateway_rejected"
	TransactionStatusFailed                 TransactionStatus = "failed"
	TransactionStatusProcessorDeclined      TransactionStatus = "processor_declined"
	TransactionStatusSettled                TransactionStatus = "settled"
	TransactionStatusSettlementConfirmed    TransactionStatus = "settlement_confirmed"
	TransactionStatusSettlementDeclined     TransactionStatus = "settlement_declined"
	TransactionStatusSettlementPending      TransactionStatus = "settlement_pending"
	TransactionStatusSettling               TransactionStatus = "settling"
	TransactionStatusSubmittedForSettlement TransactionStatus = "submitted_for_settlement"
	TransactionStatusVoided                 TransactionStatus = "voided"
	TransactionStatusUnrecognized           TransactionStatus = "unrecognized"
)

type Transaction struct {
	XMLName                      string                    `xml:"transaction"`
	Id                           string                    `xml:"id,omitempty"`
	Status                       TransactionStatus         `xml:"status,omitempty"`
	Type                         string                    `xml:"type,omitempty"`
	CurrencyISOCode              string                    `xml:"currency-iso-code,omitempty"`
	Amount                       *Decimal                  `xml:"amount"`
	OrderId                      string                    `xml:"order-id,omitempty"`
	PaymentMethodToken           string                    `xml:"payment-method-token,omitempty"`
	PaymentMethodNonce           string                    `xml:"payment-method-nonce,omitempty"`
	MerchantAccountId            string                    `xml:"merchant-account-id,omitempty"`
	PlanId                       string                    `xml:"plan-id,omitempty"`
	SubscriptionId               string                    `xml:"subscription-id,omitempty"`
	CreditCard                   *CreditCard               `xml:"credit-card,omitempty"`
	Customer                     *Customer                 `xml:"customer,omitempty"`
	BillingAddress               *Address                  `xml:"billing,omitempty"`
	ShippingAddress              *Address                  `xml:"shipping,omitempty"`
	TaxAmount                    *Decimal                  `xml:"tax-amount,omitempty"`
	TaxExempt                    bool                      `xml:"tax-exempt,omitempty"`
	DeviceData                   string                    `xml:"device-data,omitempty"`
	ServiceFeeAmount             *Decimal                  `xml:"service-fee-amount,attr,omitempty"`
	CreatedAt                    *time.Time                `xml:"created-at,omitempty"`
	UpdatedAt                    *time.Time                `xml:"updated-at,omitempty"`
	DisbursementDetails          *DisbursementDetails      `xml:"disbursement-details,omitempty"`
	RefundId                     string                    `xml:"refund-id,omitempty"`
	RefundIds                    *[]string                 `xml:"refund-ids>item,omitempty"`
	RefundedTransactionId        *string                   `xml:"refunded-transaction-id,omitempty"`
	ProcessorResponseCode        ProcessorResponseCode     `xml:"processor-response-code,omitempty"`
	ProcessorResponseText        string                    `xml:"processor-response-text,omitempty"`
	ProcessorAuthorizationCode   string                    `xml:"processor-authorization-code,omitempty"`
	SettlementBatchId            string                    `xml:"settlement-batch-id,omitempty"`
	EscrowStatus                 EscrowStatus              `xml:"escrow-status,omitempty"`
	PaymentInstrumentType        string                    `xml:"payment-instrument-type,omitempty"`
	PayPalDetails                *PayPalDetails            `xml:"paypal,omitempty"`
	VenmoAccountDetails          *VenmoAccountDetails      `xml:"venmo-account,omitempty"`
	AndroidPayDetails            *AndroidPayDetails        `xml:"android-pay-card,omitempty"`
	ApplePayDetails              *ApplePayDetails          `xml:"apple-pay,omitempty"`
	AdditionalProcessorResponse  string                    `xml:"additional-processor-response,omitempty"`
	RiskData                     *RiskData                 `xml:"risk-data,omitempty"`
	Descriptor                   *Descriptor               `xml:"descriptor,omitempty"`
	Channel                      string                    `xml:"channel,omitempty"`
	CustomFields                 customfields.CustomFields `xml:"custom-fields,omitempty"`
	AVSErrorResponseCode         AVSResponseCode           `xml:"avs-error-response-code,omitempty"`
	AVSPostalCodeResponseCode    AVSResponseCode           `xml:"avs-postal-code-response-code,omitempty"`
	AVSStreetAddressResponseCode AVSResponseCode           `xml:"avs-street-address-response-code,omitempty"`
	CVVResponseCode              CVVResponseCode           `xml:"cvv-response-code,omitempty"`
	GatewayRejectionReason       GatewayRejectionReason    `xml:"gateway-rejection-reason,omitempty"`
}

type TransactionRequest struct {
	XMLName            string                    `xml:"transaction"`
	CustomerID         string                    `xml:"customer-id,omitempty"`
	Type               string                    `xml:"type,omitempty"`
	Amount             *Decimal                  `xml:"amount"`
	OrderId            string                    `xml:"order-id,omitempty"`
	PaymentMethodToken string                    `xml:"payment-method-token,omitempty"`
	PaymentMethodNonce string                    `xml:"payment-method-nonce,omitempty"`
	MerchantAccountId  string                    `xml:"merchant-account-id,omitempty"`
	PlanId             string                    `xml:"plan-id,omitempty"`
	CreditCard         *CreditCard               `xml:"credit-card,omitempty"`
	Customer           *Customer                 `xml:"customer,omitempty"`
	BillingAddress     *Address                  `xml:"billing,omitempty"`
	ShippingAddress    *Address                  `xml:"shipping,omitempty"`
	TaxAmount          *Decimal                  `xml:"tax-amount,omitempty"`
	TaxExempt          bool                      `xml:"tax-exempt,omitempty"`
	DeviceData         string                    `xml:"device-data,omitempty"`
	Options            *TransactionOptions       `xml:"options,omitempty"`
	ServiceFeeAmount   *Decimal                  `xml:"service-fee-amount,attr,omitempty"`
	RiskData           *RiskDataRequest          `xml:"risk-data,omitempty"`
	Descriptor         *Descriptor               `xml:"descriptor,omitempty"`
	Channel            string                    `xml:"channel,omitempty"`
	CustomFields       customfields.CustomFields `xml:"custom-fields,omitempty"`
}

// TODO: not all transaction fields are implemented yet, here are the missing fields (add on demand)
//
// <transaction>
//   <voice-referral-number nil="true"></voice-referral-number>
//   <purchase-order-number nil="true"></purchase-order-number>
//   <status-history type="array">
//     <status-event>
//       <timestamp type="datetime">2013-10-07T17:26:14Z</timestamp>
//       <status>authorized</status>
//       <amount>7.00</amount>
//       <user>eaigner</user>
//       <transaction-source>Recurring</transaction-source>
//     </status-event>
//     <status-event>
//       <timestamp type="datetime">2013-10-07T17:26:14Z</timestamp>
//       <status>submitted_for_settlement</status>
//       <amount>7.00</amount>
//       <user>eaigner</user>
//       <transaction-source>Recurring</transaction-source>
//     </status-event>
//     <status-event>
//       <timestamp type="datetime">2013-10-08T07:06:38Z</timestamp>
//       <status>settled</status>
//       <amount>7.00</amount>
//       <user nil="true"></user>
//       <transaction-source></transaction-source>
//     </status-event>
//   </status-history>
//   <plan-id>bronze</plan-id>
//   <subscription-id>jqsydb</subscription-id>
//   <subscription>
//     <billing-period-end-date type="date">2013-11-06</billing-period-end-date>
//     <billing-period-start-date type="date">2013-10-07</billing-period-start-date>
//   </subscription>
//   <add-ons type="array"/>
//   <discounts type="array"/>
//   <descriptor>
//     <name nil="true"></name>
//     <phone nil="true"></phone>
//   </descriptor>
//   <recurring type="boolean">true</recurring>
// </transaction>

type Transactions struct {
	Transaction []*Transaction `xml:"transaction"`
}

type TransactionOptions struct {
	SubmitForSettlement              bool                             `xml:"submit-for-settlement,omitempty"`
	StoreInVault                     bool                             `xml:"store-in-vault,omitempty"`
	AddBillingAddressToPaymentMethod bool                             `xml:"add-billing-address-to-payment-method,omitempty"`
	StoreShippingAddressInVault      bool                             `xml:"store-shipping-address-in-vault,omitempty"`
	HoldInEscrow                     bool                             `xml:"hold-in-escrow,omitempty"`
	TransactionOptionsPaypalRequest  *TransactionOptionsPaypalRequest `xml:"paypal,omitempty"`
}

type TransactionOptionsPaypalRequest struct {
	CustomField       string
	PayeeEmail        string
	Description       string
	SupplementaryData map[string]string
}

func (r TransactionOptionsPaypalRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type transactionOptionsPaypalRequest TransactionOptionsPaypalRequest

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if r.CustomField != "" {
		if err := e.EncodeElement(r.CustomField, xml.StartElement{Name: xml.Name{Local: "custom-field"}}); err != nil {
			return err
		}
	}
	if r.PayeeEmail != "" {
		if err := e.EncodeElement(r.PayeeEmail, xml.StartElement{Name: xml.Name{Local: "payee-email"}}); err != nil {
			return err
		}
	}
	if r.Description != "" {
		if err := e.EncodeElement(r.Description, xml.StartElement{Name: xml.Name{Local: "description"}}); err != nil {
			return err
		}
	}
	if len(r.SupplementaryData) > 0 {
		start := xml.StartElement{Name: xml.Name{Local: "supplementary-data"}}
		if err := e.EncodeToken(start); err != nil {
			return err
		}
		for k, v := range r.SupplementaryData {
			if err := e.EncodeElement(v, xml.StartElement{Name: xml.Name{Local: k}}); err != nil {
				return err
			}
		}
		if err := e.EncodeToken(start.End()); err != nil {
			return err
		}
	}

	if err := e.EncodeToken(start.End()); err != nil {
		return err
	}

	if err := e.Flush(); err != nil {
		return err
	}

	return nil
}

type TransactionSearchResult struct {
	XMLName           string         `xml:"credit-card-transactions"`
	CurrentPageNumber int            `xml:"current-page-number"`
	PageSize          int            `xml:"page-size"`
	TotalItems        int            `xml:"total-items"`
	Transactions      []*Transaction `xml:"transaction"`
}

type RiskData struct {
	ID       string `xml:"id"`
	Decision string `xml:"decision"`
}

type RiskDataRequest struct {
	CustomerBrowser string `xml:"customer-browser"`
	CustomerIP      string `xml:"customer-ip"`
}
