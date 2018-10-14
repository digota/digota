package braintree

import (
	"encoding/xml"
	"time"
)

const (
	CheckWebhook                             = "check"
	DisbursementWebhook                      = "disbursement"
	DisbursementExceptionWebhook             = "disbursement_exception"
	SubscriptionCanceledWebhook              = "subscription_canceled"
	SubscriptionChargedSuccessfullyWebhook   = "subscription_charged_successfully"
	SubscriptionChargedUnsuccessfullyWebhook = "subscription_charged_unsuccessfully"
	SubscriptionExpiredWebhook               = "subscription_expired"
	SubscriptionTrialEndedWebhook            = "subscription_trial_ended"
	SubscriptionWentActiveWebhook            = "subscription_went_active"
	SubscriptionWentPastDueWebhook           = "subscription_went_past_due"
	SubMerchantAccountApprovedWebhook        = "sub_merchant_account_approved"
	SubMerchantAccountDeclinedWebhook        = "sub_merchant_account_declined"
	PartnerMerchantConnectedWebhook          = "partner_merchant_connected"
	PartnerMerchantDisconnectedWebhook       = "partner_merchant_disconnected"
	PartnerMerchantDeclinedWebhook           = "partner_merchant_declined"
	TransactionSettledWebhook                = "transaction_settled"
	TransactionSettlementDeclinedWebhook     = "transaction_settlement_declined"
	TransactionDisbursedWebhook              = "transaction_disbursed"
	DisputeOpenedWebhook                     = "dispute_opened"
	DisputeLostWebhook                       = "dispute_lost"
	DisputeWonWebhook                        = "dispute_won"
	AccountUpdaterDailyReportWebhook         = "account_updater_daily_report"
)

type WebhookNotification struct {
	XMLName   xml.Name        `xml:"notification"`
	Timestamp time.Time       `xml:"timestamp"`
	Kind      string          `xml:"kind"`
	Subject   *webhookSubject `xml:"subject"`
}

func (n *WebhookNotification) MerchantAccount() *MerchantAccount {
	if n.Subject.APIErrorResponse != nil && n.Subject.APIErrorResponse.MerchantAccount != nil {
		return n.Subject.APIErrorResponse.MerchantAccount
	} else if n.Subject.MerchantAccount != nil {
		return n.Subject.MerchantAccount
	}
	return nil
}

func (n *WebhookNotification) Disbursement() *Disbursement {
	if n.Subject.Disbursement != nil {
		return n.Subject.Disbursement
	} else {
		return nil
	}
}

type webhookSubject struct {
	XMLName          xml.Name         `xml:"subject"`
	APIErrorResponse *BraintreeError  `xml:",omitempty"`
	Disbursement     *Disbursement    `xml:"disbursement,omitempty"`
	Subscription     *Subscription    `xml:",omitempty"`
	MerchantAccount  *MerchantAccount `xml:"merchant-account,omitempty"`
	Transaction      *Transaction     `xml:",omitempty"`

	// Remaining Fields:
	// partner_merchant
}
