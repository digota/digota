//   Copyright 2015 Digota Inc.
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
	productInterface "github.com/digota/digota/product"
	"github.com/digota/digota/product/productpb"
	"github.com/digota/digota/sku"
	"github.com/digota/digota/storage"
	"github.com/digota/digota/storage/object"
	"github.com/digota/digota/validation"
	"golang.org/x/net/context"
	"time"
)

const ns = "product"

func init() {
	productInterface.RegisterService(&productService{})
}

type products []*productpb.Product

func (p *products) GetNamespace() string { return ns }

type product struct {
	productpb.Product `bson:",inline"`
}

func (p *product) GetNamespace() string { return ns }

func (p *product) SetId(id string) { p.Id = id }

func (p *product) SetCreated(t int64) { p.Created = t }

func (p *product) SetUpdated(t int64) { p.Updated = t }

type productService struct{}

// New
func (s *productService) New(ctx context.Context, req *productpb.NewRequest) (*productpb.Product, error) {

	if err := validation.Validate(req); err != nil {
		return nil, err
	}

	p := &product{
		Product: productpb.Product{
			Name:        req.GetName(),
			Description: req.GetDescription(),
			Shippable:   req.GetShippable(),
			Images:      req.GetImages(),
			Attributes:  req.GetAttributes(),
			Metadata:    req.GetMetadata(),
			Active:      req.GetActive(),
			Url:         req.GetUrl(),
		},
	}

	return &p.Product, storage.Handler().Insert(p)

}

// Get
func (s *productService) Get(ctx context.Context, req *productpb.GetRequest) (*productpb.Product, error) {

	if err := validation.Validate(req); err != nil {
		return nil, err
	}

	p := &product{
		Product: productpb.Product{
			Id: req.Id,
		},
	}

	unlock, err := locker.Handler().TryLock(p,time.Second)
	if err != nil {
		return nil, err
	}
	defer unlock()

	if err := storage.Handler().One(p); err != nil {
		return nil, err
	}

	// get products skus

	skus, err := sku.Service().ProductData(ctx, &sku.ProductDataReq{Id: p.GetId()})

	if err != nil {
		return nil, err
	}

	p.Skus = skus

	return &p.Product, nil

}

// UpdateProduct
func (s *productService) Update(ctx context.Context, req *productpb.UpdateRequest) (*productpb.Product, error) {

	if err := validation.Validate(req); err != nil {
		return nil, err
	}

	p := &product{
		Product: productpb.Product{
			Id: req.Id,
		},
	}

	unlock, err := locker.Handler().TryLock(p,time.Second)
	if err != nil {
		return nil, err
	}
	defer unlock()

	if err := storage.Handler().One(p); err != nil {
		return nil, err
	}

	// update fields and keep the rest the same

	p.Shippable = req.GetShippable()
	p.Active = req.GetActive()

	if x := req.GetName(); x != "" {
		p.Name = x
	}

	if x := req.GetDescription(); x != "" {
		p.Description = x
	}

	if x := req.GetImages(); x != nil {
		p.Images = x
	}

	if x := req.GetAttributes(); x != nil {
		p.Attributes = x
	}

	if x := req.GetMetadata(); x != nil {
		p.Metadata = x
	}

	if x := req.GetUrl(); x != "" {
		p.Url = x
	}

	return &p.Product, storage.Handler().Update(p)

}

// DeleteProduct
func (s *productService) Delete(ctx context.Context, req *productpb.DeleteRequest) (*productpb.Empty, error) {

	if err := validation.Validate(req); err != nil {
		return nil, err
	}

	p := &product{
		Product: productpb.Product{
			Id: req.Id,
		},
	}

	unlock, err := locker.Handler().TryLock(p,time.Second)
	if err != nil {
		return nil, err
	}
	defer unlock()

	// remove product skus
	//
	//skus, err := sku.Service().GetProductSkus(ctx, &sku.GetProductSkusRequest{Id: product.GetId()})
	//
	//if err != nil {
	//	return nil, err
	//}
	//
	//// delete skus with best effort
	//for _, v := range skus {
	//	sku.Service().Delete(ctx, &skupb.DeleteRequest{
	//		Id: v.GetId(),
	//	})
	//}

	return &productpb.Empty{}, storage.Handler().Remove(p)

}

// List
func (s *productService) List(ctx context.Context, req *productpb.ListRequest) (*productpb.ProductList, error) {

	if err := validation.Validate(req); err != nil {
		return nil, err
	}

	slice := &products{}

	n, err := storage.Handler().List(slice, object.ListOpt{
		Limit: req.GetLimit(),
		Page:  req.GetPage(),
		Sort:  object.SortNatural,
	})

	if err != nil {
		return nil, err
	}

	return &productpb.ProductList{Products: *slice, Total: int32(n)}, nil

}
