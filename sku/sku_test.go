//     Digota <http://digota.com> - eCommerce microservice
//     Copyright (C) 2017  Yaron Sumel <yaron@digota.com>. All Rights Reserved.
//
//     This program is free software: you can redistribute it and/or modify
//     it under the terms of the GNU Affero General Public License as published
//     by the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.
//
//     This program is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU Affero General Public License for more details.
//
//     You should have received a copy of the GNU Affero General Public License
//     along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
