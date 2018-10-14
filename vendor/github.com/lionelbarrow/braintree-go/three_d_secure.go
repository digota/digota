package braintree

type ThreeDSecureInfo struct {
	Status                 ThreeDSecureStatus   `xml:"status"`
	Enrolled               ThreeDSecureEnrolled `xml:"enrolled"`
	LiabilityShiftPossible bool                 `xml:"liability-shift-possible"`
	LiabilityShifted       bool                 `xml:"liability-shifted"`
}

type ThreeDSecureStatus string

const (
	ThreeDSecureStatusUnsupportedCard                              ThreeDSecureStatus = "unsupported_card"
	ThreeDSecureStatusLookupError                                  ThreeDSecureStatus = "lookup_error"
	ThreeDSecureStatusLookupEnrolled                               ThreeDSecureStatus = "lookup_enrolled"
	ThreeDSecureStatusLookupNotEnrolled                            ThreeDSecureStatus = "lookup_not_enrolled"
	ThreeDSecureStatusAuthenticateSuccessfulIssuerNotParticipating ThreeDSecureStatus = "authenticate_successful_issuer_not_participating"
	ThreeDSecureStatusAuthenticationUnavailable                    ThreeDSecureStatus = "authentication_unavailable"
	ThreeDSecureStatusAuthenticateSignatureVerificationFailed      ThreeDSecureStatus = "authenticate_signature_verification_failed"
	ThreeDSecureStatusAuthenticateSuccessful                       ThreeDSecureStatus = "authenticate_successful"
	ThreeDSecureStatusAuthenticateAttemptSuccessful                ThreeDSecureStatus = "authenticate_attempt_successful"
	ThreeDSecureStatusAuthenticateFailed                           ThreeDSecureStatus = "authenticate_failed"
	ThreeDSecureStatusAuthenticateUnableToAuthenticate             ThreeDSecureStatus = "authenticate_unable_to_authenticate"
	ThreeDSecureStatusAuthenticateError                            ThreeDSecureStatus = "authenticate_error"
)

type ThreeDSecureEnrolled string

const (
	ThreeDSecureEnrolledYes            ThreeDSecureEnrolled = "Y"
	ThreeDSecureEnrolledNo             ThreeDSecureEnrolled = "N"
	ThreeDSecureEnrolledUnavailable    ThreeDSecureEnrolled = "U"
	ThreeDSecureEnrolledBypass         ThreeDSecureEnrolled = "B"
	ThreeDSecureEnrolledRequestFailure ThreeDSecureEnrolled = "E"
)
