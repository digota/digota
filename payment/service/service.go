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
	"github.com/digota/digota/locker"
	"github.com/digota/digota/payment"
	"github.com/digota/digota/payment/paymentpb"
	"github.com/digota/digota/payment/service/providers"
	"github.com/digota/digota/storage"
	"github.com/digota/digota/storage/object"
	"github.com/digota/digota/util"
	"github.com/digota/digota/validation"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const ns = "charge"

func init() {
	payment.RegisterService(&PaymentService{})
}

type Charges []*paymentpb.Charge

func (c *Charges) GetNamespace() string { return ns }

type Charge struct {
	paymentpb.Charge `bson:",inline"`
}

func (c *Charge) GetNamespace() string { return ns }

func (c *Charge) SetId(id string) { c.Id = id }

func (c *Charge) SetUpdated(t int64) { c.Updated = t }

func (c *Charge) SetCreated(t int64) { c.Created = t }

type PaymentService struct{}

// Get implements the payment.pb/Get method.
// Get charge object by charge id or error.
func (p *PaymentService) Get(ctx context.Context, req *paymentpb.GetRequest) (*paymentpb.Charge, error) {

	if err := validation.Validate(req); err != nil {
		return nil, err
	}

	c := &Charge{
		Charge: paymentpb.Charge{
			Id: req.GetId(),
		},
	}

	if unlock, err := locker.Handler().TryLock(c, locker.DefaultTimeout); err != nil {
		return nil, err
	} else {
		defer unlock()
	}

	return &c.Charge, storage.Handler().One(c)

}

// Get implements the payment.pb/Get method.
// Get charge object by charge id or error.
func (p *PaymentService) List(ctx context.Context, req *paymentpb.ListRequest) (*paymentpb.ChargeList, error) {

	if err := validation.Validate(req); err != nil {
		return nil, err
	}

	slice := &Charges{}

	n, err := storage.Handler().List(slice, object.ListOpt{
		Limit: req.GetLimit(),
		Page:  req.GetPage(),
		Sort:  object.SortNatural,
	})

	if err != nil {
		return nil, err
	}

	return &paymentpb.ChargeList{Charges: *slice, Total: int32(n)}, nil

}

// Charge
func (p *PaymentService) Charge(ctx context.Context, req *paymentpb.ChargeRequest) (*paymentpb.Charge, error) {

	if err := validation.Validate(req); err != nil {
		return nil, err
	}

	provider := providers.Provider(req.GetPaymentProviderId())

	// check if card type is supported with payment provider
	if err := func() error {
		for _, v := range provider.SupportedCards() {
			if req.GetCard().GetType() == v {
				return nil
			}
		}
		return status.Error(codes.Internal, "Card type is not supported with payment provider.")
	}(); err != nil {
		return nil, err
	}

	ch, err := provider.Charge(req)
	if err != nil {
		return nil, err
	}

	//if ch == nil {
	//	return nil, status.Error(codes.Internal, "Something went wrong with the charge. 0")
	//}

	charge := &Charge{
		Charge: *ch,
	}

	// critical operations wrapped util.Retry to keep trying when failing
	if err := util.Retry(func() (err error) { return storage.Handler().Insert(charge) }); err != nil {
		// if Insert failed => refund that amount instantly with the provider
		if _, err := provider.Refund(ch.ProviderChargeId, uint64(req.GetTotal()), req.GetCurrency(), paymentpb.RefundReason_GeneralError); err != nil {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "Something went wrong with the charge.")
	}

	return &charge.Charge, nil

}

// Refund
//
//
//
//
func (p *PaymentService) Refund(ctx context.Context, req *paymentpb.RefundRequest) (*paymentpb.Charge, error) {

	if err := validation.Validate(req); err != nil {
		return nil, err
	}

	c := &Charge{
		Charge: paymentpb.Charge{
			Id: req.GetId(),
		},
	}

	if unlock, err := locker.Handler().Lock(c); err != nil {
		return nil, err
	} else {
		defer unlock()
	}

	if err := storage.Handler().One(c); err != nil {
		return nil, err
	}

	if !c.Paid || c.GetChargeAmount() <= 0 || req.GetAmount() > c.GetChargeAmount() || c.GetRefundAmount()+req.GetAmount() > c.GetChargeAmount() {
		return nil, status.Error(codes.Canceled, "Refund is unavailable for this charge.")
	}

	refund, err := providers.Provider(c.ProviderId).Refund(c.ProviderChargeId, uint64(req.GetAmount()), c.GetCurrency(), req.GetReason())
	if err != nil {
		return nil, err
	}

	// append refund to refunds
	c.Refunds = append(c.Refunds, refund)
	// sum amount of refunds
	c.RefundAmount += refund.RefundAmount
	// mark as refunded
	c.Refunded = true
	// update charge
	if err := util.Retry(func() (err error) { return storage.Handler().Update(c) }); err != nil {
		return nil, status.Error(codes.DataLoss, "Storage could not update object.")
	}

	// return response
	return &c.Charge, nil

}
