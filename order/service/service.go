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
	"errors"
	"fmt"
	"github.com/digota/digota/locker"
	orderInterface "github.com/digota/digota/order"
	"github.com/digota/digota/order/orderpb"
	"github.com/digota/digota/payment"
	"github.com/digota/digota/payment/paymentpb"
	"github.com/digota/digota/sku"
	"github.com/digota/digota/sku/skupb"
	"github.com/digota/digota/storage"
	"github.com/digota/digota/storage/object"
	"github.com/digota/digota/util"
	"github.com/digota/digota/validation"
	"github.com/rhymond/go-money"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
	"time"
)

const (
	ns                         = "order"
	orderTTL                   = time.Minute * 2
	defaultTaxDescription      = "Tax"
	defaultDiscountDescription = "Discount"
	defaultShippingDescription = "Shipping"
)

func init() {
	orderInterface.RegisterService(&orderService{})
}

type lockedOrderItem struct {
	OrderItem *orderpb.OrderItem
	Sku       *skupb.Sku
	Unlock    func() error
	Update    func() error
	Err       error
}

type orders []*orderpb.Order

func (o *orders) GetNamespace() string { return ns }

// Order wrapper
type order struct {
	orderpb.Order `bson:",inline"`
}

// implements object.Interface interface
func (o *order) GetNamespace() string { return ns }

// implements object.IdSetter interface
func (o *order) SetId(id string) { o.Id = id }

// implements object.TimeTracker interface
func (o *order) SetCreated(t int64) { o.Created = t }

// implements object.TimeTracker interface
func (o *order) SetUpdated(t int64) { o.Updated = t }

// IsReturnable checks che
func (o *order) IsReturnable(amount int64) error {
	if o.Status != orderpb.Order_Paid && o.Status != orderpb.Order_Fulfilled && o.Status != orderpb.Order_Canceled {
		return status.Error(codes.Internal, "Order is not paid or fulfilled.")
	}
	// if refund amount is bigger than the order amount return err
	if amount > o.GetAmount() {
		return status.Error(codes.Internal, "Refund amount is greater then order amount.")
	}
	return nil
}

func (o *order) IsPayable() error {
	if o.Status != orderpb.Order_Created {
		return status.Error(codes.Internal, "Order is not in created status.")
	}
	if time.Since(time.Unix(o.Created, 0)) > orderTTL {
		return status.Error(codes.Internal, "Order is too old for paying.")
	}
	if o.GetAmount() <= 0 {
		return status.Error(codes.Internal, "Order amount is Zero.")
	}
	return nil
}

type orderService struct{}

// New implements the orderpb.New interface.
// Creates new order from order Items and return Order or error, error will
// returned if something went wrong .. let say something such as inactive Items.
func (s *orderService) New(ctx context.Context, req *orderpb.NewRequest) (*orderpb.Order, error) {
	// validate input
	if err := validation.Validate(req); err != nil {
		return nil, err
	}
	// order wrapper
	o := &order{
		Order: orderpb.Order{
			Email:    req.GetEmail(),
			Status:   orderpb.Order_Created,
			Amount:   0,
			Currency: req.GetCurrency(),
			Shipping: req.GetShipping(),
		},
	}
	// get relevant order items
	orderItems, err := getUpdatedOrderItems(req.GetItems())
	if err != nil {
		return nil, err
	}
	// update order items
	o.Items = orderItems
	// calculate final amount
	amount, err := calculateTotal(o.GetCurrency(), orderItems)
	if err != nil {
		return nil, err
	}
	// update order amount
	o.Amount = amount
	// Insert and return order
	return &o.Order, storage.Handler().Insert(o)
}

// Get implements the orderpb.Get interface.
// Retrieve order or returns error.
func (s *orderService) Get(ctx context.Context, req *orderpb.GetRequest) (*orderpb.Order, error) {
	// validate input
	if err := validation.Validate(req); err != nil {
		return nil, err
	}
	// order wrapper
	o := &order{
		Order: orderpb.Order{
			Id: req.GetId(),
		},
	}
	// acquire order lock
	if unlock, err := locker.Handler().TryLock(o, locker.DefaultTimeout); err != nil {
		return nil, err
	} else {
		defer unlock()
	}
	// get and return order
	return &o.Order, storage.Handler().One(o)
}

// List implements the orderpb.List interface.
func (s *orderService) List(ctx context.Context, req *orderpb.ListRequest) (*orderpb.OrderList, error) {

	if err := validation.Validate(req); err != nil {
		return nil, err
	}

	slice := orders{}

	n, err := storage.Handler().List(&slice, object.ListOpt{
		Limit: req.GetLimit(),
		Page:  req.GetPage(),
		Sort:  object.SortNatural,
	})

	if err != nil {
		return nil, err
	}

	return &orderpb.OrderList{Orders: slice, Total: int32(n)}, nil

}

// Pay implements the orderpb.Pay interface.
// Pay will call payment service to charge the same order amount, charge id will
// be assigned to the order. This method locks the order till defers.
func (s *orderService) Pay(ctx context.Context, req *orderpb.PayRequest) (*orderpb.Order, error) {
	// validate input
	if err := validation.Validate(req); err != nil {
		return nil, err
	}
	// order wrapper
	o := &order{
		Order: orderpb.Order{
			Id: req.GetId(),
		},
	}
	// lock order for any change!
	if unlock, err := locker.Handler().TryLock(o, time.Second*5); err != nil {
		return nil, err
	} else {
		defer unlock()
	}
	// get the order
	if err := storage.Handler().One(o); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// check if order is payable
	if err := o.IsPayable(); err != nil {
		return nil, err
	}
	// lock all inventory order items (inventory objects)
	lockedItems := getLockedOrderItems(ctx, &o.Order)
	// Free all locks at func return
	defer func() {
		for _, v1 := range lockedItems {
			v1.Unlock()
		}
	}()
	// Check for errors and oversell
	for _, item := range lockedItems {
		if item.Err != nil {
			return nil, item.Err
		}
		// check for oversell
		if item.Sku.Inventory.Type == skupb.Inventory_Finite && item.Sku.Inventory.Quantity < item.OrderItem.Quantity {
			return nil, status.Error(codes.Canceled, "Oversell "+item.Sku.Id)
		}
	}
	// charge full amount for the order order
	c, err := payment.Service().Charge(ctx, &paymentpb.ChargeRequest{
		PaymentProviderId: req.GetPaymentProviderId(),
		Card:              req.GetCard(),
		Total:             uint64(o.GetAmount()),
		Currency:          o.GetCurrency(),
		Email:             o.GetEmail(),
		Statement:         fmt.Sprintf("Order %s", o.GetId()),
	})
	// return the charge error
	if err != nil {
		return nil, err
	}
	// Order has been Paid !
	// Update order object
	o.ChargeId = c.GetId()
	o.Status = orderpb.Order_Paid
	// update order with retries
	updateErr := util.Retry(func() error {
		return storage.Handler().Update(o)
	})
	// update has been failed after few times.. refund it to prevent data corruption
	// todo change to two-phase payment method ?
	if updateErr != nil {
		r, err := payment.Service().Refund(ctx, &paymentpb.RefundRequest{
			Id:     c.Id,
			Amount: uint64(o.GetAmount()),
			Reason: paymentpb.RefundReason_GeneralError,
		})
		if err != nil {
			return nil, status.Error(codes.DataLoss, fmt.Sprintf("could not update order {%s} object and could not refund the charge {%s}!", o.Id, o.ChargeId))
		}
		return nil, status.Error(codes.DataLoss, fmt.Sprintf("could not update order {%s} object, order has been refunded {%s}!", o.Id, r.Id))
	}
	// update all inventories
	for _, item := range lockedItems {
		if item.Sku.Inventory.Type == skupb.Inventory_Finite {
			// update inventory Quantity
			item.Sku.Inventory.Quantity -= item.OrderItem.Quantity
			item.Update()
		}
	}
	//
	return &o.Order, nil
}

// Pay implements the orderpb.Pay interface. Pay on specific orderId with CC based on payment provider supplied.
func (s *orderService) Return(ctx context.Context, req *orderpb.ReturnRequest) (*orderpb.Order, error) {
	// validate input
	if err := validation.Validate(req); err != nil {
		return nil, err
	}
	// order wrapper
	o := &order{
		Order: orderpb.Order{
			Id: req.GetId(),
		},
	}
	// lock order
	if unlock, err := locker.Handler().Lock(o); err != nil {
		return nil, err
	} else {
		defer unlock()
	}
	// get order
	if err := storage.Handler().One(o); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// calculate returns amount
	amount, err := calculateTotal(o.GetCurrency(), o.GetItems())
	if err != nil {
		return nil, err
	}
	// check if order can be refunded
	if err := o.IsReturnable(amount); err != nil {
		return nil, err
	}
	// lock all inventory order items (inventory objects)
	lockedItems := getLockedOrderItems(ctx, &o.Order)
	// Free all locks at func return
	defer func() {
		for _, item := range lockedItems {
			item.Unlock()
		}
	}()
	// refund the order
	if _, err := payment.Service().Refund(ctx, &paymentpb.RefundRequest{Id: o.GetChargeId(), Amount: uint64(amount)}); err != nil {
		return nil, err
	}
	// update order status
	switch o.Status {
	// if the order has been paid but never fulfilled
	// order status turns from paid into canceled
	// inventory item will get updated since the sku is in
	// stock again
	case orderpb.Order_Paid:
		o.Status = orderpb.Order_Canceled
		// notify listeners we want to return the items back in inventory
		// update inventories
		for _, item := range lockedItems {
			if item.Sku.Inventory.Type == skupb.Inventory_Finite {
				// update inventory Quantity
				item.Sku.Inventory.Quantity += item.OrderItem.Quantity
				item.Update()
			}
		}
		// if the order has been fulfilled
		// order status turns from fulfilled into returned
		// inventory will not get updated since we still got no
		// item to sell again
	case orderpb.Order_Fulfilled:
		o.Status = orderpb.Order_Returned
	}
	// update order with retries
	updateErr := util.Retry(func() error {
		return storage.Handler().Update(o)
	})
	// return err
	if updateErr != nil {
		return nil, status.Error(codes.DataLoss, fmt.Sprintf("could not update order {%s} object, order has been refunded {%s}!", o.Id, o.ChargeId))
	}
	// return order
	return &o.Order, nil
}

// calculateTotal will calculate the new amount of the cart. using
// go-money library, which helps us to do money calculations of
// the `Fowler's Money pattern`. will return error if something went
// wrong, like wrong currencies.
func calculateTotal(currency paymentpb.Currency, orderItems []*orderpb.OrderItem) (int64, error) {
	var err error
	m := money.New(0, currency.String())
	for _, v := range orderItems {
		if v.Quantity <= 0 {
			v.Quantity = 1
		}
		m, err = m.Add(money.New(v.Quantity*v.Amount, v.Currency.String()))
		if err != nil {
			return 0, status.Error(codes.Internal, err.Error())
		}
	}
	return m.Amount(), nil
}

func getUpdatedOrderItems(reqItems []*orderpb.OrderItem) (orderItems []*orderpb.OrderItem, err error) {

	var skuMap = make(map[string]*orderpb.OrderItem)
	var mtx = sync.Mutex{}
	var errs []error
	var wg = sync.WaitGroup{}

	// merge duplicated items
	for _, v := range reqItems {
		switch v.GetType() {
		case orderpb.OrderItem_sku:
			if skuItem, ok := skuMap[v.Parent]; ok {
				skuItem.Quantity += v.Quantity
				continue
			} else {
				skuMap[v.Parent] = v
			}
			fallthrough
		case orderpb.OrderItem_discount:
			if v.Quantity <= 0 {
				v.Quantity = 1
			}
			fallthrough
		case orderpb.OrderItem_shipping:
			if v.Quantity <= 0 {
				v.Quantity = 1
			}
			fallthrough
		case orderpb.OrderItem_tax:
			if v.Quantity <= 0 {
				v.Quantity = 1
			}
			orderItems = append(orderItems, v)
		}
	}

	// update order item data
	for _, v := range orderItems {
		switch v.GetType() {
		case orderpb.OrderItem_sku:
			wg.Add(1)
			// get the sku object
			go func(orderItem *orderpb.OrderItem, wg *sync.WaitGroup) {
				defer wg.Done()
				if item, err := sku.Service().Get(context.Background(), &skupb.GetRequest{Id: orderItem.GetParent()}); err != nil {
					mtx.Lock()
					errs = append(errs, err)
					mtx.Unlock()
				} else {
					orderItem.Amount = int64(item.GetPrice())
					orderItem.Currency = item.GetCurrency()
					orderItem.Description = item.GetName()
				}
			}(v, &wg)
		case orderpb.OrderItem_discount:
			// nothing to fetch yet
			if v.Description == "" {
				v.Description = defaultDiscountDescription
			}
			fallthrough
		case orderpb.OrderItem_shipping:
			// nothing to fetch yet
			if v.Description == "" {
				v.Description = defaultShippingDescription
			}
			fallthrough
		case orderpb.OrderItem_tax:
			// nothing to fetch yet
			if v.Description == "" {
				v.Description = defaultTaxDescription
			}
		}
	}

	// wait for goroutines
	wg.Wait()

	// return the first error
	if errs != nil {
		return nil, errs[0]
	}

	return
}

//
func getLockedOrderItems(ctx context.Context, order *orderpb.Order) (items []*lockedOrderItem) {

	var wg = sync.WaitGroup{}
	var mtx = sync.Mutex{}

	for _, orderItem := range order.GetItems() {
		switch orderItem.Type {
		case orderpb.OrderItem_sku:
			wg.Add(1)
			go func(orderItem *orderpb.OrderItem, wg *sync.WaitGroup) {
				defer wg.Done()
				item, unlock, set, err := sku.Service().GetWithInventoryLock(ctx, &sku.GetWithInventoryLockRequest{
					Id:       orderItem.Parent,
					Duration: time.Second,
				})
				mtx.Lock()
				items = append(items, &lockedOrderItem{
					OrderItem: orderItem,
					Sku:       item,
					Err:       err,
					Unlock: func() error {
						if unlock != nil {
							return unlock()
						}
						return errors.New("Unlock is nil")
					},
					Update: func() error {
						return set()
					},
				})
				mtx.Unlock()
			}(orderItem, &wg)
		case orderpb.OrderItem_discount:
			fallthrough
		case orderpb.OrderItem_shipping:
			fallthrough
		case orderpb.OrderItem_tax:
			// nothing to get or lock
			continue
		}

	}
	wg.Wait()
	return
}
