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

package payment

import (
	"github.com/digota/digota/payment/paymentpb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"reflect"
	"regexp"
	"testing"
)

// dummy service
type dummyService struct{}

// dummy implementations
func (s *dummyService) NewCharge(context.Context, *paymentpb.ChargeRequest) (*paymentpb.Charge, error) {
	return nil, nil
}
func (s *dummyService) Get(context.Context, *paymentpb.GetRequest) (*paymentpb.Charge, error) {
	return nil, nil
}
func (s *dummyService) RefundCharge(context.Context, *paymentpb.RefundRequest) (*paymentpb.Charge, error) {
	return nil, nil
}
func (s *dummyService) List(context.Context, *paymentpb.ListRequest) (*paymentpb.ChargeList, error) {
	return nil, nil
}

func TestRegisterService(t *testing.T) {
	service := &dummyService{}
	RegisterService(service)
	if !reflect.DeepEqual(s, service) {
		t.FailNow()
	}
}

func TestRegisterOrderServer(t *testing.T) {
	s = &dummyService{}
	server := grpc.NewServer()
	RegisterPaymentServer(server)
}

func TestService(t *testing.T) {
	s = &dummyService{}
	if !reflect.DeepEqual(Service(), s) {
		t.FailNow()
	}
	s = nil
	defer func() {
		if r := recover(); r == nil {
			t.Fatal(r)
		}
	}()
	Service()
}

func TestWriteMethods(t *testing.T) {
	methods := []*regexp.Regexp{
		regexp.MustCompile(baseMethod + "NewCharge"),
		regexp.MustCompile(baseMethod + "RefundCharge"),
	}
	// check methods in same order
	for k, v := range WriteMethods() {
		if v.String() != methods[k].String() {
			t.FailNow()
		}
	}
}

func TestReadMethods(t *testing.T) {
	methods := []*regexp.Regexp{
		regexp.MustCompile(baseMethod + "List"),
		regexp.MustCompile(baseMethod + "Get"),
	}
	// check methods in same order
	for k, v := range ReadMethods() {
		if v.String() != methods[k].String() {
			t.FailNow()
		}
	}
}
