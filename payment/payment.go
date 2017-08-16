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
	"google.golang.org/grpc"
	"regexp"
)

const baseMethod = "^(.paymentpb.Payment/)"

var s Interface

// Interface defines the functionality of the payment service
type Interface interface {
	paymentpb.PaymentServer
}

// RegisterService register p as the service provider
func RegisterService(p Interface) {
	s = p
}

// RegisterPaymentServer register service to the grpc server
func RegisterPaymentServer(server *grpc.Server) {
	paymentpb.RegisterPaymentServer(server, Service())
}

// Service return the registered service
func Service() Interface {
	if s == nil {
		panic("PaymentService is not registered")
	}
	return s
}

// ReadMethods returns regexp slice of readable methods, mostly used by the acl
func ReadMethods() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(baseMethod + "List"),
		regexp.MustCompile(baseMethod + "Get"),
	}
}

// WriteMethods returns regexp slice of writable methods, mostly used by the acl
func WriteMethods() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(baseMethod + "Charge"),
		regexp.MustCompile(baseMethod + "Refund"),
	}
}
