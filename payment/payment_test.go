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
