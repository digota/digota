//    Copyright 2015 Digota Inc.
//
//    This program is free software: you can redistribute it and/or  modify
//    it under the terms of the GNU Affero General Public License, version 3,
//    as published by the Free Software Foundation.
//
//    This program is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU Affero General Public License for more details.
//
//    You should have received a copy of the GNU Affero General Public License
//    along with this program.  If not, see <https://www.gnu.org/licenses/agpl-3.0.en.html>.

package service

import (
	"github.com/digota/digota/locker"
	"github.com/digota/digota/payment/paymentpb"
	"github.com/digota/digota/product"
	"github.com/digota/digota/product/productpb"
	skuInterface "github.com/digota/digota/sku"
	"github.com/digota/digota/sku/skupb"
	"github.com/digota/digota/storage"
	"github.com/digota/digota/storage/object"
	"github.com/digota/digota/util"
	"github.com/digota/digota/validation"
	"golang.org/x/net/context"
	"time"
)

const ns = "sku"

func init() {
	skuInterface.RegisterService(&skuService{})
}

type skus []*skupb.Sku

func (s *skus) GetNamespace() string { return ns }

type sku struct {
	skupb.Sku `bson:",inline"`
}

func (s *sku) GetNamespace() string { return ns }

func (s *sku) SetId(id string) { s.Id = id }

func (s *sku) SetCreated(t int64) { s.Created = t }

func (s *sku) SetUpdated(t int64) { s.Updated = t }

// service implementations

type skuService struct{}

func (s *skuService) New(ctx context.Context, req *skupb.NewRequest) (*skupb.Sku, error) {

	if err := validation.Validate(req); err != nil {
		return nil, err
	}

	p, err := product.Service().Get(ctx, &productpb.GetRequest{Id: req.GetParent()})
	if err != nil {
		return nil, err
	}

	var validAttr = make(map[string]string)
	for k, v := range req.GetAttributes() {
		for _, pv := range p.GetAttributes() {
			if k == pv {
				validAttr[k] = v
			}
		}
	}

	item := &sku{
		Sku: skupb.Sku{
			Price:             req.GetPrice(),
			Currency:          req.GetCurrency(),
			Active:            req.GetActive(),
			Metadata:          req.GetMetadata(),
			Name:              req.GetName(),
			Parent:            req.GetParent(),
			Image:             req.GetImage(),
			Attributes:        validAttr,
			Inventory:         req.GetInventory(),
			PackageDimensions: req.GetPackageDimensions(),
		},
	}

	return &item.Sku, storage.Handler().Insert(item)

}

func (s *skuService) Get(ctx context.Context, req *skupb.GetRequest) (*skupb.Sku, error) {

	// input validation
	if err := validation.Validate(req); err != nil {
		return nil, err
	}

	// item wrapper
	item := &sku{
		Sku: skupb.Sku{
			Id: req.GetId(),
		},
	}

	// acquire lock
	unlock, err := locker.Handler().TryLock(item,time.Second)
	if err != nil {
		return nil, err
	}
	defer unlock()

	// return item or error
	return &item.Sku, storage.Handler().One(item)

}

func (s *skuService) Update(ctx context.Context, req *skupb.UpdateRequest) (*skupb.Sku, error) {

	// input validation
	if err := validation.Validate(req); err != nil {
		return nil, err
	}

	// item wrapper
	item := &sku{
		Sku: skupb.Sku{
			Id: req.GetId(),
		},
	}

	// acquire lock
	unlock, err := locker.Handler().TryLock(item,time.Second)
	if err != nil {
		return nil, err
	}
	defer unlock()

	if err := storage.Handler().One(item); err != nil {
		return nil, err
	}

	if parent := req.GetParent(); parent != "" {
		item.Parent = parent
	}

	p, err := product.Service().Get(ctx, &productpb.GetRequest{
		Id: item.GetParent(),
	})

	if err != nil {
		return nil, err
	}

	var attrs = make(map[string]string)

	if x := req.GetAttributes(); x != nil {
		attrs = x
	} else {
		attrs = item.GetAttributes()
	}

	var validAttr = make(map[string]string)
	// save only the valid attr
	for k, v := range attrs {
		for _, pv := range p.GetAttributes() {
			if k == pv {
				validAttr[k] = v
			}
		}
	}
	item.Attributes = validAttr

	// update fields

	item.Active = req.GetActive()

	if x := req.GetName(); x != "" {
		item.Name = x
	}

	if x := req.GetPrice(); x != 0 {
		item.Price = x
	}

	if x := req.GetCurrency(); x != paymentpb.Currency_CUR_RESERVED {
		item.Currency = x
	}

	if x := req.GetMetadata(); x != nil {
		item.Metadata = x
	}

	if x := req.GetImage(); x != "" {
		item.Image = x
	}

	if x := req.GetPackageDimensions(); x != nil {
		item.PackageDimensions = x
	}

	if x := req.GetInventory(); x != nil {
		item.Inventory = x
	}

	return &item.Sku, storage.Handler().Update(item)

}

func (s *skuService) Delete(ctx context.Context, req *skupb.DeleteRequest) (*skupb.Empty, error) {

	if err := validation.Validate(req); err != nil {
		return nil, err
	}

	item := &sku{
		Sku: skupb.Sku{
			Id: req.GetId(),
		},
	}

	unlock, err := locker.Handler().TryLock(item,time.Second)
	if err != nil {
		return nil, err
	}
	defer unlock()

	return &skupb.Empty{}, storage.Handler().Remove(item)

}

func (s *skuService) GetWithInventoryLock(ctx context.Context, req *skuInterface.GetWithInventoryLockRequest) (*skupb.Sku, func() error, util.Fn, error) {

	if err := validation.Validate(req); err != nil {
		return nil, nil, nil, err
	}

	item := &sku{
		Sku: skupb.Sku{
			Id: req.Id,
		},
	}

	//skuInventory := &internal.SKUInventory{
	//	Sku: item,
	//}

	unlock, err := locker.Handler().TryLock(item, req.Duration)
	if err != nil {
		return nil, nil, nil, err
	}

	if err := storage.Handler().One(item); err != nil {
		return nil, nil, nil, err
	}

	// release lock if inventory type is not finite
	//if sku.Inventory.Type != skupb.Inventory_Finite {
	//	unlock()
	//}

	setFn := func() error {
		return storage.Handler().Update(item)
	}

	return &item.Sku, unlock, setFn, nil

}

// List implements the orderpb.List interface.
func (s *skuService) List(ctx context.Context, req *skupb.ListRequest) (*skupb.SkuList, error) {

	if err := validation.Validate(req); err != nil {
		return nil, err
	}

	slice := skus{}

	n, err := storage.Handler().List(&slice, object.ListOpt{
		Limit: req.GetLimit(),
		Page:  req.GetPage(),
		Sort:  object.SortNatural,
	})

	if err != nil {
		return nil, err
	}

	return &skupb.SkuList{Orders: slice, Total: int32(n)}, nil

}

//
func (s *skuService) ProductData(ctx context.Context, req *skuInterface.ProductDataReq) ([]*skupb.Sku, error) {

	if err := validation.Validate(req); err != nil {
		return nil, err
	}

	slice := skus{}

	return slice, storage.Handler().ListParent(req.Id, &slice)

}
