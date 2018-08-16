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
	"reflect"
	"regexp"
	"testing"
)

// dummy service
type dummyService struct{}

func (s *dummyService) New(context.Context, *skupb.NewRequest) (*skupb.Sku, error) {
	return nil, nil
}
func (s *dummyService) Get(context.Context, *skupb.GetRequest) (*skupb.Sku, error) {
	return nil, nil
}
func (s *dummyService) Update(context.Context, *skupb.UpdateRequest) (*skupb.Sku, error) {
	return nil, nil
}
func (s *dummyService) Delete(context.Context, *skupb.DeleteRequest) (*skupb.Empty, error) {
	return nil, nil
}
func (s *dummyService) List(context.Context, *skupb.ListRequest) (*skupb.SkuList, error) {
	return nil, nil
}
func (s *dummyService) GetWithInventoryLock(ctx context.Context, req *GetWithInventoryLockRequest) (*skupb.Sku, func() error, util.Fn, error) {
	return nil, nil, nil, nil
}
func (s *dummyService) ProductData(ctx context.Context, req *ProductDataReq) ([]*skupb.Sku, error) {
	return nil, nil
}

func TestRegisterService(t *testing.T) {

}

func TestRegisterSkuServer(t *testing.T) {
	service = &dummyService{}
	server := grpc.NewServer()
	RegisterSkuServer(server)
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
