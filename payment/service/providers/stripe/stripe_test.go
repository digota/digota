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

package stripe

import (
	"github.com/digota/digota/config"
	"github.com/digota/digota/payment/paymentpb"
	"os"
	"reflect"
	"testing"
)

func GetTestKey() string {
	key := os.Getenv("STRIPE_KEY")
	if len(key) == 0 {
		panic("STRIPE_KEY environment variable is not set, but is needed to run tests!\n")
	}
	return key
}

func TestNewProvider(t *testing.T) {
	p, err := NewProvider(&config.PaymentProvider{
		Secret: GetTestKey(),
	})
	if err != nil {
		t.Fatal()
	}
	if reflect.TypeOf(p).String() != "*stripe.provider" {
		t.Fatal()
	}
}

func TestProvider_Charge(t *testing.T) {
	p, err := NewProvider(&config.PaymentProvider{
		Secret: GetTestKey(),
	})
	if err != nil {
		t.Fatal()
	}

	if _, err := p.Charge(&paymentpb.ChargeRequest{
		Total:             10 * 1000,
		Currency:          paymentpb.Currency_USD,
		Email:             "yaron@digota.com",
		Statement:         "Charge statement",
		PaymentProviderId: paymentpb.PaymentProviderId_Stripe,
		Card: &paymentpb.Card{
			Type:        paymentpb.CardType_Visa,
			CVC:         "123",
			ExpireMonth: "12",
			ExpireYear:  "2022",
			LastName:    "Sumel",
			FirstName:   "Yaron",
			Number:      "4111111111111111",
		},
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := p.Charge(&paymentpb.ChargeRequest{
		Total:             10 * 1000,
		Currency:          paymentpb.Currency_USD,
		Email:             "yaron@digota.com",
		Statement:         "Charge statement",
		PaymentProviderId: paymentpb.PaymentProviderId_Stripe,
		Card: &paymentpb.Card{
			Type:        paymentpb.CardType_Visa,
			CVC:         "123",
			ExpireMonth: "12",
			ExpireYear:  "2022",
			LastName:    "Sumel",
			FirstName:   "Yaron",
			Number:      "4000111111111111",
		},
	}); err == nil {
		t.Fatal(err)
	}

	if _, err := p.Charge(&paymentpb.ChargeRequest{
		Total:             10 * 1000,
		Currency:          paymentpb.Currency_USD,
		Email:             "yaron@digota.com",
		Statement:         "Charge statement",
		PaymentProviderId: paymentpb.PaymentProviderId_Stripe,
		Card: &paymentpb.Card{
			Type:        paymentpb.CardType_Visa,
			CVC:         "123",
			ExpireMonth: "12",
			ExpireYear:  "2010",
			Number:      "4111111111111111",
		},
	}); err == nil {
		t.Fatal(err)
	}

}

func TestProvider_ProviderId(t *testing.T) {
	p := provider{}
	if p.ProviderId() != paymentpb.PaymentProviderId(paymentpb.PaymentProviderId_Stripe) {
		t.Fatal()
	}
}

func TestProvider_Refund(t *testing.T) {
	p, err := NewProvider(&config.PaymentProvider{
		Secret: GetTestKey(),
	})
	if err != nil {
		t.Fatal()
	}

	ch, err := p.Charge(&paymentpb.ChargeRequest{
		Total:             10 * 1000,
		Currency:          paymentpb.Currency_USD,
		Email:             "yaron@digota.com",
		Statement:         "Charge statement",
		PaymentProviderId: paymentpb.PaymentProviderId_Stripe,
		Card: &paymentpb.Card{
			Type:        paymentpb.CardType_Visa,
			CVC:         "123",
			ExpireMonth: "12",
			ExpireYear:  "2022",
			LastName:    "Sumel",
			FirstName:   "Yaron",
			Number:      "4111111111111111",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := p.Refund(ch.ProviderChargeId, 10*1000, paymentpb.Currency_USD, paymentpb.RefundReason_GeneralError); err != nil {
		t.Fatal(err)
	}

	if ch, err := p.Charge(&paymentpb.ChargeRequest{
		Total:             10 * 1000,
		Currency:          paymentpb.Currency_USD,
		Email:             "yaron@digota.com",
		Statement:         "Charge statement",
		PaymentProviderId: paymentpb.PaymentProviderId_Stripe,
		Card: &paymentpb.Card{
			Type:        paymentpb.CardType_Visa,
			CVC:         "123",
			ExpireMonth: "12",
			ExpireYear:  "2022",
			LastName:    "Sumel",
			FirstName:   "Yaron",
			Number:      "4111111111111111",
		},
	}); err != nil {
		t.Fatal(err)
	} else {
		if _, err := p.Refund(ch.ProviderChargeId, 10*1000, paymentpb.Currency_USD, paymentpb.RefundReason_Duplicate); err != nil {
			t.Fatal(err)
		}
	}

	if ch, err := p.Charge(&paymentpb.ChargeRequest{
		Total:             10 * 1000,
		Currency:          paymentpb.Currency_USD,
		Email:             "yaron@digota.com",
		Statement:         "Charge statement",
		PaymentProviderId: paymentpb.PaymentProviderId_Stripe,
		Card: &paymentpb.Card{
			Type:        paymentpb.CardType_Visa,
			CVC:         "123",
			ExpireMonth: "12",
			ExpireYear:  "2022",
			LastName:    "Sumel",
			FirstName:   "Yaron",
			Number:      "4111111111111111",
		},
	}); err != nil {
		t.Fatal(err)
	} else {
		if _, err := p.Refund(ch.ProviderChargeId, 10*1000, paymentpb.Currency_USD, paymentpb.RefundReason_RequestedByCustomer); err != nil {
			t.Fatal(err)
		}
	}

	if ch, err := p.Charge(&paymentpb.ChargeRequest{
		Total:             10 * 1000,
		Currency:          paymentpb.Currency_USD,
		Email:             "yaron@digota.com",
		Statement:         "Charge statement",
		PaymentProviderId: paymentpb.PaymentProviderId_Stripe,
		Card: &paymentpb.Card{
			Type:        paymentpb.CardType_Visa,
			CVC:         "123",
			ExpireMonth: "12",
			ExpireYear:  "2022",
			LastName:    "Sumel",
			FirstName:   "Yaron",
			Number:      "4111111111111111",
		},
	}); err != nil {
		t.Fatal(err)
	} else {
		if _, err := p.Refund(ch.ProviderChargeId, 10*1000, paymentpb.Currency_USD, paymentpb.RefundReason_Fraud); err != nil {
			t.Fatal(err)
		}
	}

}

func TestProvider_SupportedCards(t *testing.T) {
	found := 0
	p := provider{}
	types := []paymentpb.CardType{
		paymentpb.CardType_Mastercard,
		paymentpb.CardType_Visa,
		paymentpb.CardType_AmericanExpress,
		paymentpb.CardType_JCB,
		paymentpb.CardType_Discover,
		paymentpb.CardType_DinersClub,
	}
	for _, v := range types {
		for _, v1 := range p.SupportedCards() {
			if v == v1 {
				found++
				break
			}
		}
	}

	if found != 6 {
		t.Fatal()
	}

}
