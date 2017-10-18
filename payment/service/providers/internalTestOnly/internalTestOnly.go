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
