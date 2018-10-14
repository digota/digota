package braintree

type CVVResponseCode string

const (
	// The CVV provided matches the information on file with the cardholder's bank.
	CVVResponseCodeMatches CVVResponseCode = "M"

	// The CVV provided does not match the information on file with the cardholder's bank.
	CVVResponseCodeDoesNotMatch CVVResponseCode = "N"

	// The card-issuing bank received the CVV but did not verify whether it was correct.
	// This typically happens if the processor declines an authorization before the bank evaluates the CVV.
	CVVResponseCodeNotVerified CVVResponseCode = "U"

	// No CVV was provided.
	CVVResponseCodeNotProvided CVVResponseCode = "I"

	// The CVV was provided but the card-issuing bank does not participate in card verification.
	CVVResponseCodeIssuerDoesNotParticipate CVVResponseCode = "S"

	// The CVV was provided but this type of transaction does not support card verification.
	CVVResponseCodeNotApplicable CVVResponseCode = "A"
)
