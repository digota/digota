/*
Package braintree is a client library for Braintree.

Initializing

Initialize it with API Keys:

	braintree.New(Sandbox, merchantId, publicKey, privateKey)

Initialize it with an Access Token:

	braintree.NewWithAccessToken(accessToken)

Loggers and HTTP Clients

Optionally configure a logger and HTTP client:

	bt := braintree.New(...)
	bt.Logger = log.New(...)
	bt.HttpClient = ...

Creating Transactions

Create transactions:

	ctx := context.Background()
	t, err := bt.Transaction().Create(ctx, &braintree.TransactionRequest{
		Type:   "sale",
		Amount: braintree.NewDecimal(100, 2), // $1.00
		PaymentMethodNonce: braintree.FakeNonceTransactable,
	})

API Errors

API errors are intended to be consumed in two ways. One, they can be dealt with as a single unit:

	t, err := bt.Transaction().Create(...)
	err.Error() => "A top level error message"

Second, you can drill down to see specific error messages on a field-by-field basis:

	err.For("Transaction").On("Base")[0].Message => "A more specific error message"
*/
package braintree
