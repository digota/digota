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
