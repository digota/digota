package braintree

type AVSResponseCode string

const (
	// The postal code or street address provided matches the information on file with the cardholder's bank.
	AVSResponseCodeMatches AVSResponseCode = "M"

	// The postal code or street address provided does not match the information on file with the cardholder's bank.
	AVSResponseCodeDoesNotMatch AVSResponseCode = "N"

	// The card-issuing bank received the postal code or street address but did not verify whether it was correct.
	// This typically happens if the processor declines an authorization before the bank evaluates the postal code.
	AVSResponseCodeNotVerified AVSResponseCode = "U"

	// No postal code or street address was provided.
	AVSResponseCodeNotProvided AVSResponseCode = "I"

	// AVS information was provided but the card-issuing bank does not participate in address verification.
	// This typically indicates a card-issuing bank outside of the US, Canada, and the UK.
	AVSResponseCodeNotSupported AVSResponseCode = "S"

	// A system error prevented any verification of street address or postal code.
	AVSResponseCodeSystemError AVSResponseCode = "E"

	// AVS information was provided but this type of transaction does not support address verification.
	AVSResponseCodeNotApplicable AVSResponseCode = "A"
)
