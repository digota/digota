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

package sku

import (
	"github.com/digota/digota/sku/skupb"
	"github.com/digota/digota/util"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"regexp"
	"time"
)

const baseMethod = "^(.skupb.Sku/)"

type Interface interface {
	skupb.SkuServer
	GetWithInventoryLock(ctx context.Context, req *GetWithInventoryLockRequest) (*skupb.Sku, func() error, util.Fn, error)
	ProductData(ctx context.Context, req *ProductDataReq) ([]*skupb.Sku, error)
}

type GetWithInventoryLockRequest struct {
	Id       string `validate:"uuid4"`
	Duration time.Duration
}

type ProductDataReq struct {
	Id string `validate:"uuid4"`
}

var service Interface

func RegisterSkuServer(server *grpc.Server) {
	skupb.RegisterSkuServer(server, Service())
}

func RegisterService(p Interface) {
	if service != nil {
		panic("SkuService is already registered")
	}
	service = p
}

func Service() Interface {
	if service == nil {
		panic("SkuService is not registered")
	}
	return service
}

func WriteMethods() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(baseMethod + "New"),
		regexp.MustCompile(baseMethod + "Update"),
		regexp.MustCompile(baseMethod + "Delete"),
	}
}

func ReadMethods() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(baseMethod + "Get"),
		regexp.MustCompile(baseMethod + "List"),
	}
}
