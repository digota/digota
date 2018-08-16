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

package sku

import (
	"github.com/digota/digota/sku/skupb"
	"github.com/digota/digota/util"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"regexp"
	"time"
)

const baseMethod = "^(.skupb.SkuService/)"

// Interface defines the functionality of the sku service
type Interface interface {
	skupb.SkuServiceServer
	GetWithInventoryLock(ctx context.Context, req *GetWithInventoryLockRequest) (*skupb.Sku, func() error, util.Fn, error)
	ProductData(ctx context.Context, req *ProductDataReq) ([]*skupb.Sku, error)
}

// GetWithInventoryLockRequest request for getting inventory lock
type GetWithInventoryLockRequest struct {
	Id       string `validate:"uuid4"`
	Duration time.Duration
}

// ProductDataReq get product data req
type ProductDataReq struct {
	Id string `validate:"uuid4"`
}

var service Interface

// RegisterSkuServer register service to the grpc server
func RegisterSkuServer(server *grpc.Server) {
	skupb.RegisterSkuServiceServer(server, Service())
}

// RegisterService register service as the service provider
func RegisterService(p Interface) {
	if service != nil {
		panic("SkuService is already registered")
	}
	service = p
}

// Service return the registered service
func Service() Interface {
	if service == nil {
		panic("SkuService is not registered")
	}
	return service
}

// WriteMethods returns regexp slice of writable methods, mostly used by the acl
func WriteMethods() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(baseMethod + "New"),
		regexp.MustCompile(baseMethod + "Update"),
		regexp.MustCompile(baseMethod + "Delete"),
	}
}

// ReadMethods returns regexp slice of readable methods, mostly used by the acl
func ReadMethods() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(baseMethod + "Get"),
		regexp.MustCompile(baseMethod + "List"),
	}
}
