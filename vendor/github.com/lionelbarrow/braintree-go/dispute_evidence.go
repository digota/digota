package braintree

import (
	"encoding/xml"
	"time"
)

type DisputeEvidenceCategory string

const (
	DisputeEvidenceCategoryDeviceId                                   DisputeEvidenceCategory = "DEVICE_ID"
	DisputeEvidenceCategoryDeviceName                                 DisputeEvidenceCategory = "DEVICE_NAME"
	DisputeEvidenceCategoryPriorDigitalGoodsTransactionArn            DisputeEvidenceCategory = "PRIOR_DIGITAL_GOODS_TRANSACTION_ARN"
	DisputeEvidenceCategoryPriorDigitalGoodsTransactionDateTime       DisputeEvidenceCategory = "PRIOR_DIGITAL_GOODS_TRANSACTION_DATE_TIME"
	DisputeEvidenceCategoryDownloadDateTime                           DisputeEvidenceCategory = "DOWNLOAD_DATE_TIME"
	DisputeEvidenceCategoryGeographicalLocation                       DisputeEvidenceCategory = "GEOGRAPHICAL_LOCATION"
	DisputeEvidenceCategoryLegitPaymentsForSameMerchandise            DisputeEvidenceCategory = "LEGIT_PAYMENTS_FOR_SAME_MERCHANDISE"
	DisputeEvidenceCategoryMerchantWebsiteOrAppAccess                 DisputeEvidenceCategory = "MERCHANT_WEBSITE_OR_APP_ACCESS"
	DisputeEvidenceCategoryPriorNonDisputedTransactionArn             DisputeEvidenceCategory = "PRIOR_NON_DISPUTED_TRANSACTION_ARN"
	DisputeEvidenceCategoryPriorNonDisputedTransactionDateTime        DisputeEvidenceCategory = "PRIOR_NON_DISPUTED_TRANSACTION_DATE_TIME"
	DisputeEvidenceCategoryPriorNonDisputedTransactionEmailAddress    DisputeEvidenceCategory = "PRIOR_NON_DISPUTED_TRANSACTION_EMAIL_ADDRESS"
	DisputeEvidenceCategoryPriorNonDisputedTransactionIpAddress       DisputeEvidenceCategory = "PRIOR_NON_DISPUTED_TRANSACTION_IP_ADDRESS"
	DisputeEvidenceCategoryPriorNonDisputedTransactionPhoneNumber     DisputeEvidenceCategory = "PRIOR_NON_DISPUTED_TRANSACTION_PHONE_NUMBER"
	DisputeEvidenceCategoryPriorNonDisputedTransactionPhysicalAddress DisputeEvidenceCategory = "PRIOR_NON_DISPUTED_TRANSACTION_PHYSICAL_ADDRESS"
	DisputeEvidenceCategoryProfileSetupOrAppAccess                    DisputeEvidenceCategory = "PROFILE_SETUP_OR_APP_ACCESS"
	DisputeEvidenceCategoryProofOfAuthorizedSigner                    DisputeEvidenceCategory = "PROOF_OF_AUTHORIZED_SIGNER"
	DisputeEvidenceCategoryProofOfDeliveryEmpAddress                  DisputeEvidenceCategory = "PROOF_OF_DELIVERY_EMP_ADDRESS"
	DisputeEvidenceCategoryProofOfDelivery                            DisputeEvidenceCategory = "PROOF_OF_DELIVERY"
	DisputeEvidenceCategoryProofOfPossessionOrUsage                   DisputeEvidenceCategory = "PROOF_OF_POSSESSION_OR_USAGE"
	DisputeEvidenceCategoryPurchaserEmailAddress                      DisputeEvidenceCategory = "PURCHASER_EMAIL_ADDRESS"
	DisputeEvidenceCategoryPurchaserIpAddress                         DisputeEvidenceCategory = "PURCHASER_IP_ADDRESS"
	DisputeEvidenceCategoryPurchaserName                              DisputeEvidenceCategory = "PURCHASER_NAME"
	DisputeEvidenceCategoryRecurringTransactionArn                    DisputeEvidenceCategory = "RECURRING_TRANSACTION_ARN"
	DisputeEvidenceCategoryRecurringTransactionDateTime               DisputeEvidenceCategory = "RECURRING_TRANSACTION_DATE_TIME"
	DisputeEvidenceCategorySignedDeliveryForm                         DisputeEvidenceCategory = "SIGNED_DELIVERY_FORM"
	DisputeEvidenceCategorySignedOrderForm                            DisputeEvidenceCategory = "SIGNED_ORDER_FORM"
	DisputeEvidenceCategoryTicketProof                                DisputeEvidenceCategory = "TICKET_PROOF"
)

type DisputeEvidence struct {
	XMLName           string                  `xml:"evidence"`
	Comment           string                  `xml:"comment"`
	CreatedAt         *time.Time              `xml:"created-at"`
	ID                string                  `xml:"id"`
	SentToProcessorAt string                  `xml:"sent-to-processor-at"`
	URL               string                  `xml:"url"`
	Category          DisputeEvidenceCategory `xml:"category"`
	SequenceNumber    string                  `xml:"sequence-number"`
}

type DisputeTextEvidenceRequest struct {
	XMLName        xml.Name                `xml:"evidence"`
	Content        string                  `xml:"comments"`
	Category       DisputeEvidenceCategory `xml:"category,omitempty"`
	SequenceNumber string                  `xml:"sequence-number,omitempty"`
}
