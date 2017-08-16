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
	"google.golang.org/grpc"
	"regexp"
)

const baseMethod = "^(.orderpb.Order/)"

var service Interface

// Interface defines the functionality of the order service
type Interface interface {
	orderpb.OrderServer
}

// RegisterService register p as the service provider
func RegisterService(p Interface) {
	if service != nil {
		panic("OrderService is already registered")
	}
	service = p
}

// Register register cartServer in-front of the grpc server
func RegisterOrderServer(server *grpc.Server) {
	orderpb.RegisterOrderServer(server, Service())
}

// Service returns the registered service
func Service() Interface {
	if service == nil {
		panic("OrderService is not registered")
	}
	return service
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
		regexp.MustCompile(baseMethod + "Pay"),
		regexp.MustCompile(baseMethod + "Return"),
	}
}
