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

// Interface defines the functionality of the product service
type Interface interface {
	productpb.ProductServer
}

// RegisterService register p as the service provider
func RegisterService(p Interface) {
	if service != nil {
		panic("ProductService is already registered")
	}
	service = p
}

// Service return the registered service
func Service() Interface {
	if service == nil {
		panic("ProductService is not registered")
	}
	return service
}

// RegisterProductServer register service to the grpc server
func RegisterProductServer(server *grpc.Server) {
	productpb.RegisterProductServer(server, Service())
}

// ReadMethods returns regexp slice of readable methods, mostly used by the acl
func ReadMethods() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(baseMethod + "Get"),
		regexp.MustCompile(baseMethod + "List"),
	}
}

// WriteMethods returns regexp slice of writable methods, mostly used by the acl
func WriteMethods() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(baseMethod + "New"),
		regexp.MustCompile(baseMethod + "Update"),
		regexp.MustCompile(baseMethod + "Delete"),
	}
}
