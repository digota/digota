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

package internalTestOnly

import (
	"github.com/digota/digota/payment/paymentpb"
	"reflect"
	"testing"
)

func TestNewProvider(t *testing.T) {
	p, _ := NewProvider()
	if reflect.TypeOf(p).String() != "*internalTestOnly.provider" {
		t.Fatal()
	}
}

func TestProvider_SupportedCards(t *testing.T) {
	found := 0
	p := provider{}
	types := []paymentpb.CardType{
		paymentpb.CardType_Visa,
	}
	for _, v := range types {
		for _, v1 := range p.SupportedCards() {
			if v == v1 {
				found++
				break
			}
		}
	}
	if found != 1 {
		t.Fatal()
	}
}

func TestProvider_Charge(t *testing.T) {
	p, _ := NewProvider()
	if _, err := p.Charge(&paymentpb.ChargeRequest{
		Currency:          paymentpb.Currency_USD,
		Total:             120,
		Email:             "aa@aa.com",
		PaymentProviderId: paymentpb.PaymentProviderId(10),
		Statement:         "statement",
		Card:              &paymentpb.Card{},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := p.Charge(&paymentpb.ChargeRequest{
		Currency:          paymentpb.Currency_USD,
		Total:             120,
		Email:             "error@error.com",
		PaymentProviderId: paymentpb.PaymentProviderId(10),
		Statement:         "statement",
		Card:              &paymentpb.Card{},
	}); err == nil {
		t.Fatal(err)
	}
}

func TestProvider_Refund(t *testing.T) {
	p, _ := NewProvider()
	if _, err := p.Refund("chargeid", 123, paymentpb.Currency_USD, paymentpb.RefundReason_RequestedByCustomer); err != nil {
		t.Fatal(err)
	}
	if _, err := p.Refund("chargeid", 990099, paymentpb.Currency_USD, paymentpb.RefundReason_RequestedByCustomer); err == nil {
		t.Fatal(err)
	}
}
