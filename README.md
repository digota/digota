<a href="http://digota.com/">![Logo](http://i.imgur.com/hqEKC51.png)</a>
## Digota - ecommerce microservice [![Go Report Card](https://goreportcard.com/badge/github.com/digota/digota)](https://goreportcard.com/report/github.com/digota/digota) [![Build Status](https://travis-ci.org/digota/digota.svg?branch=master)](https://travis-ci.org/digota/digota) [![Coverage Status](https://coveralls.io/repos/github/digota/digota/badge.svg?branch=master)](https://coveralls.io/github/digota/digota?branch=master) [![Gitter chat](https://badges.gitter.im/gitterHQ/gitter.png)](https://gitter.im/digota/Lobby)

[Digota](http://digota.com) is ecommerce microservice built to be the modern standard for ecommerce systems.It is based on grpc,protocol-buffers and http2 provides clean, powerful and secured RPC interface. 

Our Goal is to provide the best technology that covers most of the ecommerce flows, just focus of your business logic and not on the ecommerce logic.

___TLDR; scalable ecommerce microservice.___

### Join our community at [gitter](https://gitter.im/digota/Lobby) ! 

## Getting started

#### Prerequisites

* Go > 1.8
* Database 
  * mongodb > 3.2
  * redis (TBD)
  * postgresql (TBD - [#2](https://github.com/digota/digota/issues/2))
* Lock server (default is in-memory locker)
  * zookeeper 
  * redis !! (thanks @Gerifield)
  * etcd (TBD - [#3](https://github.com/digota/digota/issues/3))

#### Installation

From source  

```bash
$ go get -u github.com/digota/digota
```

From docker hub

```bash
$ docker pull digota/digota:0.1
```

#### Run

```bash
$ docker run digota/digota:0.1
```

Or with flags

```bash
$ docker run digota/digota:0.1 bash -c "digota --version"
```

Check out this [docker-compose](https://github.com/digota/digota/blob/master/docker/docker-compose.yml) for more details.

Flags:
       --info                  Set log level to info
       --debug                 Set log level to debug
       --help, -h              show help
       --version, -v           print the version

## Cross languages

Key benefit of using grpc is the native support of major languages (`C++`,`Java`,`Python`,`Go`,`Ruby`,`Node.js`,`C#`,`Objective-C`,`Android Java` and `PHP`). 
Learn How to compile your client right [here](https://grpc.io/docs/quickstart/), You can use you `Makefile` as well.

Complied clients:
1. [php](https://github.com/digota/digota-php)

## Flexible payment gateways

It does not matter which payment gateway you are using, it is just matter of config to register it.

Supported gateways for now: 
1. Stripe
2. Braintree

> Are you payment provider ? 
> Just implement the following [interface](https://github.com/digota/digota/blob/master/payment/service/providers/providers.go#L32) and PR your changes.

## Auth & Security

##### We take security very seriously, don't hesitate to report a security issue.

Digota is fully Encrypted (end-to-end) using TLS, That fact is leveraged also to Authenticate Clients based on their Certificate in front of the Local Certificate Authority.
Basically we are creating CA and signing any certificate we want to approve with same CA. 

> How about revoking certificate? The CRL approch here is whitelist instead of blacklist, just remove client serial from your config.

##### Create CA

```bash
$ certstrap init --common-name "ca.company.com"
```

##### Create Client Certificate

```bash
$ certstrap request-cert --domain client.company.com
```

##### Sign Certificate

```bash
$ certstrap sign --CA "ca.company.com" client.company.com
```

##### Approve Certificate

Take the certificate serial and Append the serial and scopes(`WRITE`,`READ`,`WILDCARD`) to your config

```bash
$ openssl x509 -in out/client.com.crt -serial | grep -Po '(?<=serial=)\w+'
output: A2FF9503829A3A0DDE9CB87191A472D4
```

Follow [these](https://github.com/digota/digota/tree/master/_example/auth) steps to create your CA and Certificates.

## Money & Currencies

Floats are tricky when it comes to money, we don't want to lose money so the chosen money representation here is 
based on the [smallest currency unit](https://martinfowler.com/eaaCatalog/money.html). For example: `4726` is `$47.26`.

## Distributed lock

All the important data usage is `Exclusively Guaranteed`, means that you don't need to worry about any concurrent data-race across different nodes.
Typical data access is as following:
```
Client #1 GetSomething -> TryLock -> [lock accuired] ->  DoSomething -> ReleaseLock -> Return Something 
                                                                                 \ 
Client #2 GetSomething -> TryLock -> --------- [wait for lock] -------------------*-----> [lock accuired] -> ...
                                         
Client #3 GetSomething -> TryLock -> -------------------- [wait for lock] ---> [accuire error] -> Return Error
```

## Core Services 

### Payment

```proto
service Payment {
    rpc Charge  (chargeRequest) returns (charge)        {}
    rpc Refund  (refundRequest) returns (charge)        {}
    rpc Get     (getRequest)    returns (charge)        {}
    rpc List    (listRequest)   returns (chargeList)    {}
}
```

___Full service [definition](https://github.com/digota/digota/blob/master/payment/paymentpb/payment.proto).___

Payment service is used for credit/debit card charge and refund, it is provides support of multiple 
payment providers as well. Usually there is no use in this service externally if you are using `order` functionality.

### Order

```proto
service Order {
    rpc New     (newRequest)    returns (order)         {}
    rpc Get     (getRequest)    returns (order)         {}
    rpc Pay     (payRequest)    returns (order)         {}
    rpc Return  (returnRequest) returns (order)         {}
    rpc List    (listRequest)   returns (listResponse)  {}
}
```

___Full service [definition](https://github.com/digota/digota/blob/master/order/orderpb/order.proto).___

Order service helps you deal with structured purchases ie `order`. Naturally order is a collection of purchasable
products,discounts,invoices and basic customer information.

### Product

```proto
service Product {
    rpc New     (newRequest)    returns (product)       {}
    rpc Get     (getRequest)    returns (product)       {}
    rpc Update  (updateRequest) returns (product)       {}
    rpc Delete  (deleteRequest) returns (empty)         {}
    rpc List    (listRequest)   returns (productList)   {}
}
```

___Full service [definition](https://github.com/digota/digota/blob/master/product/productpb/product.proto).___

Product service helps you manage your products, product represent collection of purchasable items(sku), physical or digital.

### Sku

```proto
service Sku {
    rpc New     (newRequest)    returns (sku)           {}
    rpc Get     (getRequest)    returns (sku)           {}
    rpc Update  (updateRequest) returns (sku)           {}
    rpc Delete  (deleteRequest) returns (empty)         {}
    rpc List    (listRequest)   returns (skuList)       {}
}
```

___Full service [definition](https://github.com/digota/digota/blob/master/sku/skupb/sku.proto).___

Sku service helps you manage your product Stock Keeping Units(SKU), sku represent specific product configuration such as attributes, currency and price.

For example, a product may be a `football ticket`, whereas a specific SKU represents the stadium section. 

Sku is also used to manage its inventory and 
prevent oversell in case that the inventory type is `Finite`. 

## Usage example

Eventually the goal is to make life easier at the client-side, 
here's golang example of creating order and paying for it.. easy as that.

### Create new order

```go
order.New(context.Background(), &orderpb.NewRequest{
    Currency: paymentpb.Currency_EUR,
    Items: []*orderpb.OrderItem{
    	{
    		Parent:   "af350ecc-56c8-485f-8858-74d4faffa9cb",
    		Quantity: 2,
    		Type:     orderpb.OrderItem_sku,
    	},
    	{
    		Amount:      -1000,
    		Description: "Discount for being loyal customer",
    		Currency:    paymentpb.Currency_EUR,
    		Type:        orderpb.OrderItem_discount,
    	},
    	{
    		Amount:      1000,
    		Description: "Tax",
    		Currency:    paymentpb.Currency_EUR,
    		Type:        orderpb.OrderItem_tax,
    	},
    },
    Email: "yaron@digota.com",
    Shipping: &orderpb.Shipping{
    	Name:  "Yaron Sumel",
    	Phone: "+972 000 000 000",
    	Address: &orderpb.Shipping_Address{
    		Line1:      "Loren ipsum",
    		City:       "San Jose",
    		Country:    "USA",
    		Line2:      "",
    		PostalCode: "12345",
    		State:      "CA",
    	},
    },
})
```

### Pay the order

```go
			
order.Pay(context.Background(), &orderpb.PayRequest{
    Id:                "bf350ecc-56c8-485f-8858-74d4faffa9cb",
    PaymentProviderId: paymentpb.PaymentProviderId_Stripe,
    Card: &paymentpb.Card{
        Type:        paymentpb.CardType_Visa,
    	CVC:         "123",
    	ExpireMonth: "12",
    	ExpireYear:  "2022",
    	LastName:    "Sumel",
    	FirstName:   "Yaron",
    	Number:      "4242424242424242",
    },
})			
			
```

## Contribution

### Development

### Donations 

## License
```
// Digota <http://digota.com> - eCommerce microservice
// Copyright (c) 2018 Yaron Sumel <yaron@digota.com>
//
// MIT License
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

```
You can find the complete license file [here](https://github.com/digota/digota/blob/master/LICENSE), for any questions regarding the license please [contact](https://github.com/digota/digota#contact) us.

## Contact
For any questions or inquiries please contact ___yaron@digota.com___
