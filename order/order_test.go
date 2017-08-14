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

package order

import (
	"github.com/digota/digota/order/orderpb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"reflect"
	"regexp"
	"testing"
)

// dummy service
type dummyService struct{}

// dummy implementations
func (s *dummyService) New(context.Context, *orderpb.NewRequest) (*orderpb.Order, error) {
	return nil, nil
}
func (s *dummyService) Get(context.Context, *orderpb.GetRequest) (*orderpb.Order, error) {
	return nil, nil
}
func (s *dummyService) Pay(context.Context, *orderpb.PayRequest) (*orderpb.Order, error) {
	return nil, nil
}
func (s *dummyService) Return(context.Context, *orderpb.ReturnRequest) (*orderpb.Order, error) {
	return nil, nil
}
func (s *dummyService) List(context.Context, *orderpb.ListRequest) (*orderpb.OrderList, error) {
	return nil, nil
}

func TestRegisterOrderServer(t *testing.T) {
	service = &dummyService{}
	server := grpc.NewServer()
	RegisterOrderServer(server)
}

func TestRegisterService(t *testing.T) {
	service = nil
	service := &dummyService{}
	RegisterService(service)
	if !reflect.DeepEqual(service, service) {
		t.Fatal()
	}
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
		regexp.MustCompile(baseMethod + "Pay"),
		regexp.MustCompile(baseMethod + "Return"),
	}
	// check methods in same order
	for k, v := range WriteMethods() {
		if v.String() != methods[k].String() {
			t.FailNow()
		}
	}
}
