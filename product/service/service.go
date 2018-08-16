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

	unlock, err := locker.Handler().TryLock(p, time.Second)
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

	unlock, err := locker.Handler().TryLock(p, time.Second)
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

	unlock, err := locker.Handler().TryLock(p, time.Second)
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
