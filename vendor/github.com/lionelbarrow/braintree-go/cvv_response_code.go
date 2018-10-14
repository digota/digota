package braintree

// CVVResponseCode is the response
// to the card verification value (CVV) provided by the customer.
type CVVResponseCode string

const (
	// CVVResponseCodeMatches means the CVV provided matches the information on file with the cardholder's bank.
	CVVResponseCodeMatches CVVResponseCode = "M"

	// CVVResponseCodeDoesNotMatch means the The CVV provided does not match the information on file with the cardholder's bank.
	CVVResponseCodeDoesNotMatch CVVResponseCode = "N"

	// CVVResponseCodeNotVerified means the the card-issuing bank received the CVV but did not verify whether it was correct.
	// This typically happens if the processor declines an authorization before the bank evaluates the CVV.
	CVVResponseCodeNotVerified CVVResponseCode = "U"

	// CVVResponseCodeNotProvided means the no CVV was provided.
	CVVResponseCodeNotProvided CVVResponseCode = "I"

	// CVVResponseCodeIssuerDoesNotParticipate means the the CVV was provided
	// but the card-issuing bank does not participate in card verification.
	CVVResponseCodeIssuerDoesNotParticipate CVVResponseCode = "S"

	// CVVResponseCodeNotApplicable means the the CVV was provided
	// but this type of transaction does not support card verification.
	CVVResponseCodeNotApplicable CVVResponseCode = "A"
)
