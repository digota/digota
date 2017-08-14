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

package acl

import (
	"github.com/digota/digota/client"
	"github.com/digota/digota/order"
	"github.com/digota/digota/payment"
	"github.com/digota/digota/product"
	"github.com/digota/digota/sku"
	"golang.org/x/net/context"
	"regexp"
)

var accessMap = map[client.Scope][][]*regexp.Regexp{
	// All methods
	client.WildcardScope: {
		[]*regexp.Regexp{
			regexp.MustCompile("(.*)"),
		},
	},
	// Public methods
	client.PublicScope: {
	//auth.PublicMethods(),
	},
	// Write only methods
	client.WriteScope: {
		payment.WriteMethods(),
		sku.WriteMethods(),
		order.WriteMethods(),
		product.WriteMethods(),
	},
	// Read only methods
	client.ReadScope: {
		payment.ReadMethods(),
		sku.ReadMethods(),
		order.ReadMethods(),
		product.ReadMethods(),
	},
}

// getAccessMap return access map for specific client
func getAccessMap(c *client.Client) [][]*regexp.Regexp {
	var m [][]*regexp.Regexp
	// append public scope
	m = append(m, accessMap[client.PublicScope]...)
	if c == nil {
		return m
	}
	for _, v := range c.Scopes {
		if methods, ok := accessMap[v]; ok {
			m = append(m, methods...)
		}
	}
	return m
}

// CanAccessMethod check if user can access fullMethod by getting its accessMap
func CanAccessMethod(ctx context.Context, fullMethod string) bool {
	u, _ := client.FromContext(ctx)
	for _, v := range getAccessMap(u) {
		for _, r := range v {
			if r.MatchString(fullMethod) {
				return true
			}
		}
	}
	return false
}
