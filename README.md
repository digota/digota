<a href="http://digota.com/">![Logo](http://i.imgur.com/hqEKC51.png)</a>
## Digota - ecommerce microservice [![Go Report Card](https://goreportcard.com/badge/github.com/digota/digota)](https://goreportcard.com/report/github.com/digota/digota) [![Build Status](https://travis-ci.org/digota/digota.svg?branch=master)](https://travis-ci.org/digota/digota)

[Digota](http://digota.com) is ecommerce microservice built to be the modern standard for ecommerce systems.It is based on grpc,protocol-buffers and http2 provides clean, powerful and secured RPC interface. 

Our Goal is to provide the best technology that covers most of the ecommerce flows, just focus of your business logic and not on the ecommerce logic.

___TLDR; scalable ecommerce microservice.___

## Getting started

#### Prerequisites

* Go > 1.8
* Database 
  * mongodb > 3.2
* Lock server 
  * zookeeper 

#### Installation

```bash
$ go get -u github.com/digota/digota
```

#### Run

```bash
$ digota --port=8080 --config=/etc/digota/digota.yml
```
Check out this [example](https://github.com/digota/digota/blob/master/_example/config/digota.yaml) to understand how to set up your config.

## Cross languages

Key benefit of using grpc is the native support of major languages (`C++`,`Java`,`Python`,`Go`,`Ruby`,`Node.js`,`C#`,`Objective-C`,`Android Java` and `PHP`). 
Learn How to compile your client right [here](https://grpc.io/docs/quickstart/), You can use you `Makefile` as well.

```bash
$ make php
```

## Flexible payment gateways

It does not matter which payment gateway you are using, it is just matter of config to register it.

Supported gateways for now: 
1. Stripe
2. Braintree

> Are you payment provider ? 
> Just implement the following [interface](https://github.com/digota/digota/blob/master/payment/service/providers/providers.go#L32) and PR you changes.

```yaml
...
payment:
- provider: Stripe
  secret: sk_test_00000000000000000
```


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

```yml
...
clients:
- serial: "A2FF9503829A3A0DDE9CB87191A472D4"
  scopes:
  - READ
  - WRITE
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
Digota <http://digota.com> - eCommerce microservice
Copyright (C) 2017  Yaron Sumel <yaron@digota.com>. All Rights Reserved.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
```
You can find the complete license file [here](https://github.com/digota/digota/blob/master/LICENSE), for any questions regarding the license please [contact](https://github.com/digota/digota#contact) us.

## Contact
For any questions or inquiries please contact ___yaron@digota.com___
