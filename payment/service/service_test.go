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

package service

import (
	"github.com/digota/digota/config"
	"github.com/digota/digota/locker"
	"github.com/digota/digota/payment/paymentpb"
	"github.com/digota/digota/payment/service/providers"
	"github.com/digota/digota/storage"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"os"
	"testing"
	"time"
)

var db = "testing-payment-" + uuid.NewV4().String()

func TestMain(m *testing.M) {

	// setup
	if err := storage.New(config.Storage{
		Address:  []string{"localhost"},
		Handler:  "mongodb",
		Database: db,
	}); err != nil {
		panic(err)
	}

	// in-memory locker
	locker.New(config.Locker{})

	providers.New([]config.PaymentProvider{{Provider: "DigotaInternalTestOnly"}})

	retCode := m.Run()
	// teardown
	storage.Handler().DropDatabase(db)
	os.Exit(retCode)
}

func TestCharges_GetNamespace(t *testing.T) {
	o := charges{}
	if o.GetNamespace() != "charge" {
		t.FailNow()
	}
}

func TestCharge_GetNamespace(t *testing.T) {
	o := charge{}
	if o.GetNamespace() != "charge" {
		t.FailNow()
	}
}

func TestCharge_SetId(t *testing.T) {
	o := charge{}
	o.SetId("1234-1234-1234-1234")
	if o.Id != "1234-1234-1234-1234" {
		t.FailNow()
	}
}

func TestCharge_SetCreated(t *testing.T) {
	o := charge{}
	now := time.Now().Unix()
	o.SetCreated(now)
	if o.Created != now {
		t.FailNow()
	}
}

func TestCharge_SetUpdated(t *testing.T) {
	o := charge{}
	now := time.Now().Unix()
	o.SetUpdated(now)
	if o.Updated != now {
		t.FailNow()
	}
}

func TestService_Get(t *testing.T) {

}

func TestService_Charge(t *testing.T) {

	s := &paymentService{}

	chReq := &paymentpb.ChargeRequest{
		Total:     1000,
		Currency:  paymentpb.Currency_USD,
		Statement: "test statement",
		Metadata: map[string]string{
			"key": "val",
		},
		Email: "yaron@digota.com",
		Card: &paymentpb.Card{
			Type:        paymentpb.CardType_Visa,
			CVC:         "123",
			ExpireMonth: "12",
			ExpireYear:  "2022",
			LastName:    "Sumel",
			FirstName:   "Yaron",
			Number:      "4111111111111111",
		},
		// actually we are using the internal provider
		// the internal provider can't pass the input validation so
		// it is registered itself as stripe provider
		PaymentProviderId: paymentpb.PaymentProviderId_Stripe,
	}

	_, err := s.NewCharge(context.Background(), chReq)

	if err != nil {
		t.Fatal(err)
	}

	chReq.Email = "yarondigota.com"

	if _, err := s.NewCharge(context.Background(), chReq); err == nil {
		t.Fatal(err)
	}

	chReq.Email = "yaron@digota.com"
	chReq.Card.Type = paymentpb.CardType_Mastercard

	if _, err := s.NewCharge(context.Background(), chReq); err == nil {
		t.Fatal(err)
	}

	chReq.Card.Type = paymentpb.CardType_Visa
	// will cause charge error
	chReq.Email = "error@error.com"

	if _, err := s.NewCharge(context.Background(), chReq); err == nil {
		t.Fatal(err)
	}

}

func TestService_Refund(t *testing.T) {

	s := &paymentService{}

	ch, err := s.NewCharge(context.Background(), &paymentpb.ChargeRequest{
		Total:     1000,
		Currency:  paymentpb.Currency_USD,
		Statement: "test statement",
		Metadata: map[string]string{
			"key": "val",
		},
		Email: "yaron@digota.com",
		Card: &paymentpb.Card{
			Type:        paymentpb.CardType_Visa,
			CVC:         "123",
			ExpireMonth: "12",
			ExpireYear:  "2022",
			LastName:    "Sumel",
			FirstName:   "Yaron",
			Number:      "4111111111111111",
		},
		// actually we are using the internal provider
		// the internal provider can't pass the input validation so
		// it is registered itself as stripe provider
		PaymentProviderId: paymentpb.PaymentProviderId_Stripe,
	})

	if err != nil {
		t.Fatal(err)
	}

	if _, err := s.RefundCharge(context.Background(), &paymentpb.RefundRequest{
		Id:     ch.GetId(),
		Amount: ch.GetChargeAmount(),
		Reason: paymentpb.RefundReason_GeneralError,
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := s.RefundCharge(context.Background(), &paymentpb.RefundRequest{
		Id:     ch.GetId(),
		Amount: 990099,
		Reason: paymentpb.RefundReason_GeneralError,
	}); err == nil {
		t.Fatal(err)
	}

}

func TestPaymentService_Get(t *testing.T) {

	s := &paymentService{}

	ch, err := s.NewCharge(context.Background(), &paymentpb.ChargeRequest{
		Total:     1000,
		Currency:  paymentpb.Currency_USD,
		Statement: "test statement",
		Metadata: map[string]string{
			"key": "val",
		},
		Email: "yaron@digota.com",
		Card: &paymentpb.Card{
			Type:        paymentpb.CardType_Visa,
			CVC:         "123",
			ExpireMonth: "12",
			ExpireYear:  "2022",
			LastName:    "Sumel",
			FirstName:   "Yaron",
			Number:      "4111111111111111",
		},
		// actually we are using the internal provider
		// the internal provider can't pass the input validation so
		// it is registered itself as stripe provider
		PaymentProviderId: paymentpb.PaymentProviderId_Stripe,
	})

	if err != nil {
		t.Fatal(err)
	}

	if _, err := s.Get(context.Background(), &paymentpb.GetRequest{
		Id: ch.GetId(),
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := s.Get(context.Background(), &paymentpb.GetRequest{
		Id: uuid.NewV4().String(),
	}); err == nil {
		t.Fatal()
	}

}

func TestPaymentService_List(t *testing.T) {

	storage.Handler().DropDatabase(db)

	s := &paymentService{}

	for k := 0; k < 10; k++ {
		_, err := s.NewCharge(context.Background(), &paymentpb.ChargeRequest{
			Total:     1000,
			Currency:  paymentpb.Currency_USD,
			Statement: "test statement",
			Metadata: map[string]string{
				"key": "val",
			},
			Email: "yaron@digota.com",
			Card: &paymentpb.Card{
				Type:        paymentpb.CardType_Visa,
				CVC:         "123",
				ExpireMonth: "12",
				ExpireYear:  "2022",
				LastName:    "Sumel",
				FirstName:   "Yaron",
				Number:      "4111111111111111",
			},
			// actually we are using the internal provider
			// the internal provider can't pass the input validation so
			// it is registered itself as stripe provider
			PaymentProviderId: paymentpb.PaymentProviderId_Stripe,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	l, err := s.List(context.Background(), &paymentpb.ListRequest{
		Page:  0,
		Limit: 10,
	})

	if err != nil || l.Total != 10 {
		t.Fatal(err)
	}

}
