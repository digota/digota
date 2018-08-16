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
	"errors"
	"github.com/digota/digota/payment/paymentpb"
	"github.com/satori/go.uuid"
	"time"
)

type provider struct {
}

// NewProvider create new provider for testings
func NewProvider() (*provider, error) {
	return &provider{}, nil
}

func (p *provider) SupportedCards() []paymentpb.CardType {
	return []paymentpb.CardType{
		paymentpb.CardType_Visa,
	}
}

func (p *provider) ProviderId() paymentpb.PaymentProviderId {
	return paymentpb.PaymentProviderId(paymentpb.PaymentProviderId_Stripe)
}

func (p *provider) Charge(req *paymentpb.ChargeRequest) (*paymentpb.Charge, error) {

	if req.GetEmail() == "error@error.com" {
		return nil, errors.New("expected charge error")
	}

	return &paymentpb.Charge{
		ProviderId:       p.ProviderId(),
		ProviderChargeId: uuid.NewV4().String(),
		Paid:             true,
		Email:            req.GetEmail(),
		Currency:         req.GetCurrency(),
		ChargeAmount:     req.GetTotal(),
		Statement:        req.GetStatement(),
	}, nil

}

func (p *provider) Refund(ch string, amount uint64, currency paymentpb.Currency, reason paymentpb.RefundReason) (*paymentpb.Refund, error) {

	if amount == 990099 {
		return nil, errors.New("expected refund error")
	}

	return &paymentpb.Refund{
		ProviderRefundId: uuid.NewV4().String(),
		RefundAmount:     amount,
		Created:          time.Now().Unix(),
		Reason:           reason,
	}, nil
}
