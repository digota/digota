# Braintree Go

[![GoDoc](https://godoc.org/github.com/lionelbarrow/braintree-go?status.svg)](http://godoc.org/github.com/lionelbarrow/braintree-go) [![Build Status](https://travis-ci.org/lionelbarrow/braintree-go.png?branch=master)](https://travis-ci.org/lionelbarrow/braintree-go)

A Go client library for [Braintree](https://www.braintreepayments.com), the payments company behind awesome companies like GitHub, Heroku, and 37signals.

This is *not* an official client library. Braintree maintains server-side libraries for [Ruby](https://www.github.com/braintree/braintree_ruby), [Python](https://www.github.com/braintree/braintree_python), [PHP](https://www.github.com/braintree/braintree_php), [Perl](https://www.github.com/braintree/braintree_perl), [Node](https://www.github.com/braintree/braintree_node), [C#](https://www.github.com/braintree/braintree_dotnet) and [Java](https://www.github.com/braintree/braintree_java), but not Go. This package implements the core functionality of the other client libraries, but it's missing a few advanced features.

With that said, this package contains more than enough to get you started accepting payments using Braintree. If there's a feature the other client libraries implement that you really need, open an issue (or better yet, a pull request).

### Usage

Setting up your credentials is easy.

```go
import "github.com/lionelbarrow/braintree-go"

bt := braintree.New(
  braintree.Sandbox,
  "YOUR_BRAINTREE_MERCH_ID",
  "YOUR_BRAINTREE_PUB_KEY",
  "YOUR_BRAINTREE_PRIV_KEY",
)
```

So is creating your first transaction.

```go
tx, err := bt.Transaction().Create(&braintree.TransactionRequest{
  Type: "sale",
  Amount: braintree.NewDecimal(100, 2), // 100 cents
  CreditCard: &braintree.CreditCard{
    Number:         "4111111111111111",
    ExpirationDate: "05/14",
  },
})
```

The error returned by these calls is typed. The package returns a generic error when something mechanical goes wrong, such as receiving malformed XML or being unable to connect to the Braintree gateway. However, if Braintree was able to process the request correctly, but was unable to fulfill it due to a semantic failure (such as the credit card being declined) then a `BraintreeError` type is returned.

In addition to creating transactions, you can also tokenize credit card information for repeat or subscription billing using the `CreditCard`, `Customer`, and `Subscription` types. This package is completely compatible with [Braintree.js](https://www.braintreepayments.com/braintrust/braintree-js), so if you encrypt your customers' credit cards in the browser, you can pass them on to Braintree without ever seeing them yourself. This decreases your PCI regulatory exposure and helps to secure your users' data. See the examples folder for a working implementation.

### Installation

The usual. `go get github.com/lionelbarrow/braintree-go`

### Supported Go Versions

* 1.6
* 1.7
* 1.8

### Documentation

Braintree provides a [ton of documentation](https://www.braintreepayments.com/docs/ruby/guide/overview) on how to use their API. I recommend you use the Ruby documentation when following along, as the Ruby client library is broadly similar to this one.

For details on this package, see [GoDoc](http://godoc.org/github.com/lionelbarrow/braintree-go).

### Testing

The integration tests run against a sandbox account created in the [Braintree Sandbox](https://sandbox.braintreegateway.com/).
See [TESTING.md](TESTING.md) for further instructions on how to set up your sandbox for integration testing.

You can run tests locally using the same credentials that Travis CI uses by using the credentials in `.default.env`. Simply `cp .default.env .env` if you use a tool that autoloads `.env` files, or `source .default.env` to load the credentials into your shell. Then run tests with `go test ./...`.

```
source .default.env
go test ./...
```

### Webhook Integration Testing

You can use the `WebhookTestingGateway` to write your own integration tests to verify your application is processing incoming webhook notifications correctly.

A simple example:

```go
package integration_test

import (
  "testing"
  "net/http/httptest"

  "github.com/lionelbarrow/braintree-go"
)

func TestMyWebhook(t *testing.T) {
  bt := braintree.New(
    braintree.Sandbox,
    "merchaint_id",
    "public_key",
    "private_key",
  )

  r, err := bt.WebhookTesting().Request(
    braintree.SubscriptionChargedSuccessfullyWebhook,
    "123", // id
  )
  if err != nil {
    t.Fatal(err)
  }

  // You can now send the payload and signature to your webhook handler
  // and test your application's busines logic

  w := httptest.NewRecorder()
  router.ServeHTTP(w, r)

  // assertions
}
```

### License

The MIT License (MIT)

Copyright (c) 2013 Lionel Barrow

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

### Contributors

- [Erik Aigner](http://github.com/eaigner)
- [Kayle Gishen](https://github.com/kayleg)
- [Jesse Szwedko](https://github.com/jszwedko)
- [Leigh McCulloch](https://github.com/leighmcculloch)
