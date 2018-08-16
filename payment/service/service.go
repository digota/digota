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
	"github.com/digota/digota/locker"
	paymentInterface "github.com/digota/digota/payment"
	"github.com/digota/digota/payment/paymentpb"
	"github.com/digota/digota/payment/service/providers"
	"github.com/digota/digota/storage"
	"github.com/digota/digota/storage/object"
	"github.com/digota/digota/util"
	"github.com/digota/digota/validation"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

const ns = "charge"

func init() {
	paymentInterface.RegisterService(&paymentService{})
}

type charges []*paymentpb.Charge

func (c *charges) GetNamespace() string { return ns }

type charge struct {
	paymentpb.Charge `bson:",inline"`
}

func (c *charge) GetNamespace() string { return ns }

func (c *charge) SetId(id string) { c.Id = id }

func (c *charge) SetUpdated(t int64) { c.Updated = t }

func (c *charge) SetCreated(t int64) { c.Created = t }

type paymentService struct{}

// Get implements the payment.pb/Get method.
// Get charge object by charge id or error.
func (p *paymentService) Get(ctx context.Context, req *paymentpb.GetRequest) (*paymentpb.Charge, error) {

	if err := validation.Validate(req); err != nil {
		return nil, err
	}

	c := &charge{
		Charge: paymentpb.Charge{
			Id: req.GetId(),
		},
	}

	unlock, err := locker.Handler().TryLock(c, locker.DefaultTimeout)
	if err != nil {
		return nil, err
	}
	defer unlock()

	return &c.Charge, storage.Handler().One(c)

}

// Get implements the payment.pb/Get method.
// Get charge object by charge id or error.
func (p *paymentService) List(ctx context.Context, req *paymentpb.ListRequest) (*paymentpb.ChargeList, error) {

	if err := validation.Validate(req); err != nil {
		return nil, err
	}

	slice := &charges{}

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
func (p *paymentService) NewCharge(ctx context.Context, req *paymentpb.ChargeRequest) (*paymentpb.Charge, error) {

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

	charge := &charge{
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
func (p *paymentService) RefundCharge(ctx context.Context, req *paymentpb.RefundRequest) (*paymentpb.Charge, error) {

	if err := validation.Validate(req); err != nil {
		return nil, err
	}

	c := &charge{
		Charge: paymentpb.Charge{
			Id: req.GetId(),
		},
	}

	unlock, err := locker.Handler().TryLock(c, time.Second)
	if err != nil {
		return nil, err
	}
	defer unlock()

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
