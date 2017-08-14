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

	if err := locker.New(config.Locker{
		Handler: "zookeeper",
		Address: []string{"localhost"},
	}); err != nil {
		panic(err)
	}

	providers.New([]config.PaymentProvider{
		{
			Provider: "DigotaInternalTestOnly",
		},
	})

	retCode := m.Run()
	// teardown
	storage.Handler().DropDatabase(db)
	os.Exit(retCode)
}

func TestCharges_GetNamespace(t *testing.T) {
	o := Charges{}
	if o.GetNamespace() != "charge" {
		t.FailNow()
	}
}

func TestCharge_GetNamespace(t *testing.T) {
	o := Charge{}
	if o.GetNamespace() != "charge" {
		t.FailNow()
	}
}

func TestCharge_SetId(t *testing.T) {
	o := Charge{}
	o.SetId("1234-1234-1234-1234")
	if o.Id != "1234-1234-1234-1234" {
		t.FailNow()
	}
}

func TestCharge_SetCreated(t *testing.T) {
	o := Charge{}
	now := time.Now().Unix()
	o.SetCreated(now)
	if o.Created != now {
		t.FailNow()
	}
}

func TestCharge_SetUpdated(t *testing.T) {
	o := Charge{}
	now := time.Now().Unix()
	o.SetUpdated(now)
	if o.Updated != now {
		t.FailNow()
	}
}

func TestService_Get(t *testing.T) {

}

func TestService_Charge(t *testing.T) {

	s := &PaymentService{}

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

	_, err := s.Charge(context.Background(), chReq)

	if err != nil {
		t.Fatal(err)
	}

	chReq.Email = "yarondigota.com"

	if _, err := s.Charge(context.Background(), chReq); err == nil {
		t.Fatal(err)
	}

	chReq.Email = "yaron@digota.com"
	chReq.Card.Type = paymentpb.CardType_Mastercard

	if _, err := s.Charge(context.Background(), chReq); err == nil {
		t.Fatal(err)
	}

	chReq.Card.Type = paymentpb.CardType_Visa
	// will cause charge error
	chReq.Email = "error@error.com"

	if _, err := s.Charge(context.Background(), chReq); err == nil {
		t.Fatal(err)
	}

}

func TestService_Refund(t *testing.T) {

	s := &PaymentService{}

	ch, err := s.Charge(context.Background(), &paymentpb.ChargeRequest{
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

	if _, err := s.Refund(context.Background(), &paymentpb.RefundRequest{
		Id:     ch.GetId(),
		Amount: ch.GetChargeAmount(),
		Reason: paymentpb.RefundReason_GeneralError,
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := s.Refund(context.Background(), &paymentpb.RefundRequest{
		Id:     ch.GetId(),
		Amount: 990099,
		Reason: paymentpb.RefundReason_GeneralError,
	}); err == nil {
		t.Fatal(err)
	}

}

func TestPaymentService_Get(t *testing.T) {

	s := &PaymentService{}

	ch, err := s.Charge(context.Background(), &paymentpb.ChargeRequest{
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

	s := &PaymentService{}

	for k := 0; k < 10; k++ {
		_, err := s.Charge(context.Background(), &paymentpb.ChargeRequest{
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
