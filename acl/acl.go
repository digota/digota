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

// SkipAuth is
var skipAuth = false

// SkipAuth set skipAuth flag to true
func SetSkipAuth() {
	skipAuth = true
}

// GetSkipAuth get flag
func SkipAuth() bool {
	return skipAuth
}

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
