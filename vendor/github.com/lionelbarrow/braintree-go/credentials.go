package braintree

type credentials interface {
	Environment() Environment
	MerchantID() string
	AuthorizationHeader() string
}
