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
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"reflect"
	"regexp"
	"testing"
)

// dummy service
type dummyService struct{}

func (s *dummyService) New(context.Context, *productpb.NewRequest) (*productpb.Product, error) {
	return nil, nil
}
func (s *dummyService) Get(context.Context, *productpb.GetRequest) (*productpb.Product, error) {
	return nil, nil
}
func (s *dummyService) Update(context.Context, *productpb.UpdateRequest) (*productpb.Product, error) {
	return nil, nil
}
func (s *dummyService) List(context.Context, *productpb.ListRequest) (*productpb.ProductList, error) {
	return nil, nil
}
func (s *dummyService) Delete(context.Context, *productpb.DeleteRequest) (*productpb.Empty, error) {
	return nil, nil
}

func TestRegisterService(t *testing.T) {
	service := &dummyService{}
	RegisterService(service)
	if !reflect.DeepEqual(Service(), service) {
		t.FailNow()
	}
}

func TestRegisterProductServer(t *testing.T) {
	service = &dummyService{}
	server := grpc.NewServer()
	RegisterProductServer(server)
}

func TestService(t *testing.T) {
	service = &dummyService{}
	if !reflect.DeepEqual(Service(), service) {
		t.FailNow()
	}
	service = nil
	defer func() {
		if r := recover(); r == nil {
			t.Fatal(r)
		}
	}()
	Service()
}

func TestReadMethods(t *testing.T) {
	methods := []*regexp.Regexp{
		regexp.MustCompile(baseMethod + "Get"),
		regexp.MustCompile(baseMethod + "List"),
	}
	// check methods in same order
	for k, v := range ReadMethods() {
		if v.String() != methods[k].String() {
			t.FailNow()
		}
	}
}

func TestWriteMethods(t *testing.T) {
	methods := []*regexp.Regexp{
		regexp.MustCompile(baseMethod + "New"),
		regexp.MustCompile(baseMethod + "Update"),
		regexp.MustCompile(baseMethod + "Delete"),
	}
	// check methods in same order
	for k, v := range WriteMethods() {
		if v.String() != methods[k].String() {
			t.FailNow()
		}
	}
}
