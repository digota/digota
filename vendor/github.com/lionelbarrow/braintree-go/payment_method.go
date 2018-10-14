package braintree

type PaymentMethod interface {
	GetCustomerId() string
	GetToken() string
	IsDefault() bool
	GetImageURL() string
}
