package braintree

type GatewayRejectionReason string

const (
	GatewayRejectionReasonApplicationIncomplete GatewayRejectionReason = "application_incomplete"
	GatewayRejectionReasonAVS                   GatewayRejectionReason = "avs"
	GatewayRejectionReasonAVSAndCVV             GatewayRejectionReason = "avs_and_cvv"
	GatewayRejectionReasonCVV                   GatewayRejectionReason = "cvv"
	GatewayRejectionReasonDuplicate             GatewayRejectionReason = "duplicate"
	GatewayRejectionReasonFraud                 GatewayRejectionReason = "fraud"
	GatewayRejectionReasonThreeDSecure          GatewayRejectionReason = "three_d_secure"
	GatewayRejectionReasonUnrecognized          GatewayRejectionReason = "unrecognized"
)
