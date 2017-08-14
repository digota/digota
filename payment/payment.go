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

type Interface interface {
	paymentpb.PaymentServer
}

func RegisterService(p Interface) {
	s = p
}

func RegisterPaymentServer(server *grpc.Server) {
	paymentpb.RegisterPaymentServer(server, Service())
}

func Service() Interface {
	if s == nil {
		panic("PaymentService is not registered")
	}
	return s
}

func ReadMethods() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(baseMethod + "List"),
		regexp.MustCompile(baseMethod + "Get"),
	}
}

func WriteMethods() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(baseMethod + "Charge"),
		regexp.MustCompile(baseMethod + "Refund"),
	}
}
