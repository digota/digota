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

type TransactionSource string

const (
	TransactionSourceRecurringFirst TransactionSource = "recurring_first"
	TransactionSourceRecurring      TransactionSource = "recurring"
	TransactionSourceMOTO           TransactionSource = "moto"
	TransactionSourceMerchant       TransactionSource = "merchant"
)

type PaymentInstrumentType string

const (
	PaymentInstrumentTypeAndroidPayCard   PaymentInstrumentType = "android_pay_card"
	PaymentInstrumentTypeApplePayCard     PaymentInstrumentType = "apple_pay_card"
	PaymentInstrumentTypeCreditCard       PaymentInstrumentType = "credit_card"
	PaymentInstrumentTypeMasterpassCard   PaymentInstrumentType = "masterpass_card"
	PaymentInstrumentTypePaypalAccount    PaymentInstrumentType = "paypal_account"
	PaymentInstrumentTypeVenmoAccount     PaymentInstrumentType = "venmo_account"
	PaymentInstrumentTypeVisaCheckoutCard PaymentInstrumentType = "visa_checkout_card"
)

type Transaction struct {
	XMLName                      string                    `xml:"transaction"`
	Id                           string                    `xml:"id"`
	Status                       TransactionStatus         `xml:"status"`
	Type                         string                    `xml:"type"`
	CurrencyISOCode              string                    `xml:"currency-iso-code"`
	Amount                       *Decimal                  `xml:"amount"`
	OrderId                      string                    `xml:"order-id"`
	PaymentMethodToken           string                    `xml:"payment-method-token"`
	PaymentMethodNonce           string                    `xml:"payment-method-nonce"`
	MerchantAccountId            string                    `xml:"merchant-account-id"`
	PlanId                       string                    `xml:"plan-id"`
	SubscriptionId               string                    `xml:"subscription-id"`
	SubscriptionDetails          *SubscriptionDetails      `xml:"subscription"`
	CreditCard                   *CreditCard               `xml:"credit-card"`
	Customer                     *Customer                 `xml:"customer"`
	BillingAddress               *Address                  `xml:"billing"`
	ShippingAddress              *Address                  `xml:"shipping"`
	TaxAmount                    *Decimal                  `xml:"tax-amount"`
	TaxExempt                    bool                      `xml:"tax-exempt"`
	DeviceData                   string                    `xml:"device-data"`
	ServiceFeeAmount             *Decimal                  `xml:"service-fee-amount,attr"`
	CreatedAt                    *time.Time                `xml:"created-at"`
	UpdatedAt                    *time.Time                `xml:"updated-at"`
	DisbursementDetails          *DisbursementDetails      `xml:"disbursement-details"`
	RefundId                     string                    `xml:"refund-id"`
	RefundIds                    *[]string                 `xml:"refund-ids>item"`
	RefundedTransactionId        *string                   `xml:"refunded-transaction-id"`
	ProcessorResponseCode        ProcessorResponseCode     `xml:"processor-response-code"`
	ProcessorResponseText        string                    `xml:"processor-response-text"`
	ProcessorAuthorizationCode   string                    `xml:"processor-authorization-code"`
	SettlementBatchId            string                    `xml:"settlement-batch-id"`
	EscrowStatus                 EscrowStatus              `xml:"escrow-status"`
	PaymentInstrumentType        PaymentInstrumentType     `xml:"payment-instrument-type"`
	ThreeDSecureInfo             *ThreeDSecureInfo         `xml:"three-d-secure-info,omitempty"`
	PayPalDetails                *PayPalDetails            `xml:"paypal"`
	VenmoAccountDetails          *VenmoAccountDetails      `xml:"venmo-account"`
	AndroidPayDetails            *AndroidPayDetails        `xml:"android-pay-card"`
	ApplePayDetails              *ApplePayDetails          `xml:"apple-pay"`
	AdditionalProcessorResponse  string                    `xml:"additional-processor-response"`
	RiskData                     *RiskData                 `xml:"risk-data"`
	Descriptor                   *Descriptor               `xml:"descriptor"`
	Channel                      string                    `xml:"channel"`
	CustomFields                 customfields.CustomFields `xml:"custom-fields"`
	AVSErrorResponseCode         AVSResponseCode           `xml:"avs-error-response-code"`
	AVSPostalCodeResponseCode    AVSResponseCode           `xml:"avs-postal-code-response-code"`
	AVSStreetAddressResponseCode AVSResponseCode           `xml:"avs-street-address-response-code"`
	CVVResponseCode              CVVResponseCode           `xml:"cvv-response-code"`
	GatewayRejectionReason       GatewayRejectionReason    `xml:"gateway-rejection-reason"`
	PurchaseOrderNumber          string                    `xml:"purchase-order-number"`
	Disputes                     []*Dispute                `xml:"disputes>dispute"`
}

type TransactionRequest struct {
	XMLName             string                      `xml:"transaction"`
	CustomerID          string                      `xml:"customer-id,omitempty"`
	Type                string                      `xml:"type,omitempty"`
	Amount              *Decimal                    `xml:"amount"`
	OrderId             string                      `xml:"order-id,omitempty"`
	PaymentMethodToken  string                      `xml:"payment-method-token,omitempty"`
	PaymentMethodNonce  string                      `xml:"payment-method-nonce,omitempty"`
	MerchantAccountId   string                      `xml:"merchant-account-id,omitempty"`
	PlanId              string                      `xml:"plan-id,omitempty"`
	CreditCard          *CreditCard                 `xml:"credit-card,omitempty"`
	Customer            *CustomerRequest            `xml:"customer,omitempty"`
	BillingAddress      *Address                    `xml:"billing,omitempty"`
	ShippingAddress     *Address                    `xml:"shipping,omitempty"`
	TaxAmount           *Decimal                    `xml:"tax-amount,omitempty"`
	TaxExempt           bool                        `xml:"tax-exempt,omitempty"`
	DeviceData          string                      `xml:"device-data,omitempty"`
	Options             *TransactionOptions         `xml:"options,omitempty"`
	ServiceFeeAmount    *Decimal                    `xml:"service-fee-amount,attr,omitempty"`
	RiskData            *RiskDataRequest            `xml:"risk-data,omitempty"`
	Descriptor          *Descriptor                 `xml:"descriptor,omitempty"`
	Channel             string                      `xml:"channel,omitempty"`
	CustomFields        customfields.CustomFields   `xml:"custom-fields,omitempty"`
	PurchaseOrderNumber string                      `xml:"purchase-order-number,omitempty"`
	TransactionSource   TransactionSource           `xml:"transaction-source,omitempty"`
	LineItems           TransactionLineItemRequests `xml:"line-items,omitempty"`
}

func (t *Transaction) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type typeWithNoFunctions Transaction
	if err := d.DecodeElement((*typeWithNoFunctions)(t), &start); err != nil {
		return err
	}
	if t.SubscriptionDetails != nil &&
		t.SubscriptionDetails.BillingPeriodStartDate == "" &&
		t.SubscriptionDetails.BillingPeriodEndDate == "" {
		t.SubscriptionDetails = nil
	}
	return nil
}

// TODO: not all transaction fields are implemented yet, here are the missing fields (add on demand)
//
// <transaction>
//   <voice-referral-number nil="true"></voice-referral-number>
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
	SubmitForSettlement              bool                                   `xml:"submit-for-settlement,omitempty"`
	StoreInVault                     bool                                   `xml:"store-in-vault,omitempty"`
	StoreInVaultOnSuccess            bool                                   `xml:"store-in-vault-on-success,omitempty"`
	AddBillingAddressToPaymentMethod bool                                   `xml:"add-billing-address-to-payment-method,omitempty"`
	StoreShippingAddressInVault      bool                                   `xml:"store-shipping-address-in-vault,omitempty"`
	HoldInEscrow                     bool                                   `xml:"hold-in-escrow,omitempty"`
	TransactionOptionsPaypalRequest  *TransactionOptionsPaypalRequest       `xml:"paypal,omitempty"`
	SkipAdvancedFraudChecking        bool                                   `xml:"skip_advanced_fraud_checking,omitempty"`
	ThreeDSecure                     *TransactionOptionsThreeDSecureRequest `xml:"three-d-secure,omitempty"`
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

type TransactionOptionsThreeDSecureRequest struct {
	Required bool `xml:"required"`
}

type TransactionSearchResult struct {
	TotalItems int
	TotalIDs   []string

	CurrentPageNumber int
	PageSize          int
	Transactions      []*Transaction
}

type RiskData struct {
	ID       string `xml:"id"`
	Decision string `xml:"decision"`
}

type RiskDataRequest struct {
	CustomerBrowser string `xml:"customer-browser"`
	CustomerIP      string `xml:"customer-ip"`
}

type SubscriptionDetails struct {
	BillingPeriodStartDate string `xml:"billing-period-start-date"`
	BillingPeriodEndDate   string `xml:"billing-period-end-date"`
}
