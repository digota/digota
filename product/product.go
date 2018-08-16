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

package product

import (
	"github.com/digota/digota/product/productpb"
	"google.golang.org/grpc"
	"regexp"
)

const baseMethod = "^(.productpb.ProductService/)"

var service Interface

// Interface defines the functionality of the product service
type Interface interface {
	productpb.ProductServiceServer
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
	productpb.RegisterProductServiceServer(server, Service())
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
