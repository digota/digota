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
