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
	"github.com/digota/digota/payment/errors"
	"github.com/digota/digota/payment/paymentpb"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/refund"
	"strings"
	"time"
)

type provider struct {
}

// NewProvider creates and prepare the stripe provider
func NewProvider(paymentConfig *config.PaymentProvider) (*provider, error) {
	stripe.Key = paymentConfig.Secret
	stripe.Logger = logrus.New()
	return &provider{}, nil
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
	return paymentpb.PaymentProviderId(paymentpb.PaymentProviderId_Stripe)
}

func (p *provider) Charge(req *paymentpb.ChargeRequest) (*paymentpb.Charge, error) {
	// perform new charge
	ch, err := charge.New(&stripe.ChargeParams{
		Amount:   uint64(req.GetTotal()),
		Currency: stripe.Currency(strings.ToLower(req.GetCurrency().String())),
		//Desc:      charge.Description,
		Desc:  req.GetStatement(),
		Email: req.GetEmail(),
		Source: &stripe.SourceParams{
			Card: &stripe.CardParams{
				Number: req.GetCard().GetNumber(),
				Month:  req.GetCard().GetExpireMonth(),
				Year:   req.GetCard().GetExpireYear(),
				CVC:    req.GetCard().GetCVC(),
				Name:   req.GetCard().GetLastName() + " " + req.GetCard().GetFirstName(),
			},
		},
	})

	// convert err
	if err != nil || !ch.Paid {
		switch x := err.(type) {
		case *stripe.Error:
			return nil, convertStripeError(x.Code)
		default:
			return nil, err
		}
	}

	// return charge
	return &paymentpb.Charge{
		ProviderId:       p.ProviderId(),
		ProviderChargeId: ch.ID,
		Paid:             ch.Paid,
		Email:            req.GetEmail(),
		Currency:         req.GetCurrency(),
		ChargeAmount:     req.GetTotal(),
		Statement:        req.GetStatement(),
	}, nil

}

func (p *provider) Refund(chargeId string, amount uint64, currency paymentpb.Currency, reason paymentpb.RefundReason) (*paymentpb.Refund, error) {
	var stripeReason stripe.RefundReason
	// convert reason to stripe reason
	switch reason {
	case paymentpb.RefundReason_Fraud:
		stripeReason = "fraudulent"
	case paymentpb.RefundReason_Duplicate:
		stripeReason = "duplicate"
	case paymentpb.RefundReason_RequestedByCustomer:
		stripeReason = "requested_by_customer"
	}
	// perform refund
	rf, err := refund.New(&stripe.RefundParams{
		Amount: amount,
		Charge: chargeId,
		Reason: stripeReason,
	})
	// convert err
	if err != nil {
		return nil, err
	}
	return &paymentpb.Refund{
		ProviderRefundId: rf.ID,
		RefundAmount:     rf.Amount,
		Created:          time.Now().Unix(),
		Reason:           reason,
	}, nil
}

func convertStripeError(errorCode stripe.ErrorCode) error {
	switch errorCode {
	case stripe.IncorrectNum:
		return errors.IncorrectNum
	case stripe.InvalidNum:
		return errors.InvalidNum
	case stripe.InvalidExpM:
		return errors.InvalidExpM
	case stripe.InvalidExpY:
		return errors.InvalidExpY
	case stripe.InvalidCvc:
		return errors.InvalidCvc
	case stripe.ExpiredCard:
		return errors.ExpiredCard
	case stripe.IncorrectCvc:
		return errors.IncorrectCvc
	case stripe.IncorrectZip:
		return errors.IncorrectZip
	case stripe.CardDeclined:
		return errors.CardDeclined
	case stripe.Missing:
		return errors.Missing
	case stripe.ProcessingErr:
		return errors.ProcessingErr
	case stripe.RateLimit:
		return errors.RateLimit
	default:
		break
	}
	return nil
}
