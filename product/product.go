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

package product

import (
	"github.com/digota/digota/product/productpb"
	"google.golang.org/grpc"
	"regexp"
)

const baseMethod = "^(.productpb.Product/)"

var service Interface

type Interface interface {
	productpb.ProductServer
}

func RegisterService(p Interface) {
	if service != nil {
		panic("ProductService is already registered")
	}
	service = p
}

func Service() Interface {
	if service == nil {
		panic("ProductService is not registered")
	}
	return service
}

func RegisterProductServer(server *grpc.Server) {
	productpb.RegisterProductServer(server, Service())
}

func ReadMethods() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(baseMethod + "Get"),
		regexp.MustCompile(baseMethod + "List"),
	}
}

func WriteMethods() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(baseMethod + "New"),
		regexp.MustCompile(baseMethod + "Update"),
		regexp.MustCompile(baseMethod + "Delete"),
	}
}
