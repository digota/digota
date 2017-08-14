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

package braintree

import (
	"github.com/digota/digota/config"
	"github.com/digota/digota/payment/paymentpb"
	"github.com/lionelbarrow/braintree-go"
	"strings"
	"time"
)

type provider struct {
	bt *braintree.Braintree
}

func NewProvider(paymentConfig *config.PaymentProvider) (*provider, error) {
	env := braintree.Sandbox
	if paymentConfig.Live {
		env = braintree.Production
	}
	return &provider{braintree.New(
		env,
		paymentConfig.MerchId,
		paymentConfig.PublicKey,
		paymentConfig.PrivateKey,
	)}, nil
}

func (p *provider) SupportedCards() []paymentpb.CardType {
	return []paymentpb.CardType{
		paymentpb.CardType_Mastercard,
		paymentpb.CardType_Visa,
		paymentpb.CardType_AmericanExpress,
		paymentpb.CardType_JCB,
		paymentpb.CardType_Discover,
		paymentpb.CardType_DinersClub,
	}
}

func (p *provider) ProviderId() paymentpb.PaymentProviderId {
	return paymentpb.PaymentProviderId(paymentpb.PaymentProviderId_Braintree)
}

func (p *provider) Charge(req *paymentpb.ChargeRequest) (*paymentpb.Charge, error) {

	tx, err := p.bt.Transaction().Create(&braintree.TransactionRequest{
		Type:   "sale",
		Amount: braintree.NewDecimal(int64(req.GetTotal()), 2),
		CreditCard: &braintree.CreditCard{
			Number:          req.GetCard().GetNumber(),
			CVV:             req.GetCard().GetCVC(),
			ExpirationMonth: req.GetCard().GetExpireMonth(),
			ExpirationYear:  req.GetCard().GetExpireYear(),
			CardType:        strings.ToLower(req.GetCard().GetType().String()),
		},
	})

	// convert err
	if err != nil {
		return nil, err
	}

	// return charge
	return &paymentpb.Charge{
		ProviderId:       p.ProviderId(),
		ProviderChargeId: tx.Id,
		Paid:             true,
		Email:            req.GetEmail(),
		Currency:         req.GetCurrency(),
		ChargeAmount:     req.GetTotal(),
		Statement:        req.GetStatement(),
	}, nil

}

func (p *provider) Refund(chargeId string, amount uint64, currency paymentpb.Currency, reason paymentpb.RefundReason) (*paymentpb.Refund, error) {

	tx, err := p.bt.Transaction().Refund(chargeId, braintree.NewDecimal(int64(amount), 2))

	// convert err
	if err != nil {
		return nil, err
	}

	return &paymentpb.Refund{
		ProviderRefundId: tx.Id,
		RefundAmount:     uint64(tx.Amount.Unscaled),
		Created:          time.Now().Unix(),
		Reason:           reason,
	}, nil

}
