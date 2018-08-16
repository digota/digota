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

package main

import (
	"github.com/digota/digota/order/orderpb"
	"github.com/digota/digota/payment/paymentpb"
	"github.com/digota/digota/sdk"
	"golang.org/x/net/context"
	"log"
)

func main() {

	c, err := sdk.NewClient("localhost:3051", &sdk.ClientOpt{
		InsecureSkipVerify: false,
		ServerName:         "server.com",
		CaCrt:              "out/ca.crt",
		Crt:                "out/client.com.crt",
		Key:                "out/client.com.key",
	})

	if err != nil {
		panic(err)
	}

	defer c.Close()

	// Create new order
	o, err := orderpb.NewOrderServiceClient(c).New(context.Background(), &orderpb.NewRequest{
		Currency: paymentpb.Currency_USD,
		Items: []*orderpb.OrderItem{
			{
				Parent:   "af350ecc-56c8-485f-8858-74d4faffa9cb",
				Quantity: 2,
				Type:     orderpb.OrderItem_sku,
			},
			{
				Parent:   "af350ecc-56c8-485f-8858-74d4faffa9cb",
				Quantity: 2,
				Type:     orderpb.OrderItem_sku,
			},
			//{
			//	Parent:   "480e53bf-b409-4a34-8c74-13786b35ae11",
			//	Quantity: 1,
			//	Type:     orderpb.OrderItem_sku,
			//},
			//{
			//	Parent:   "480e53bf-b409-4a34-8c74-13786b35ae11",
			//	Quantity: 1,
			//	Type:     orderpb.OrderItem_sku,
			//},
			{
				Amount:      -1000,
				Description: "on the fly discount without parent",
				Currency:    paymentpb.Currency_USD,
				Type:        orderpb.OrderItem_discount,
			},
			{
				Amount:      1000,
				Description: "Tax (Included)",
				Currency:    paymentpb.Currency_USD,
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

	if err != nil {
		panic(err)
	}

	log.Println(o.GetId())

}
