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

// NewProvider creates and prepare the braintree provider
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
