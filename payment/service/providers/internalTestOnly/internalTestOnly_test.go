//     Digota <http://digota.com> - eCommerce microservice
//     Copyright (C) 2017  Yaron Sumel <yaron@digota.com>. All Rights Reserved.
//
//     This program is free software: you can redistribute it and/or modify
//     it under the terms of the GNU Affero General Public License as published
//     by the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.
//
//     This program is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU Affero General Public License for more details.
//
//     You should have received a copy of the GNU Affero General Public License
//     along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
