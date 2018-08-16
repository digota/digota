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

package client

import (
	"github.com/digota/digota/config"
	"golang.org/x/net/context"
	"math/big"
	"reflect"
	"testing"
)

func TestInit(t *testing.T) {
	New([]config.Client{
		{
			Serial: "12FF9503829A3A0DDE9CB87191A472D4",
			Scopes: []string{"READ", "WRITE"},
		},
		{
			Serial: "22FF9503829A3A0DDE9CB87191A472D4",
			Scopes: []string{"READ"},
		},
	})

	for k, v := range clients {
		switch k {
		case 0:
			if v.Serial != "12FF9503829A3A0DDE9CB87191A472D4" {
				t.Fatal()
			}
		case 1:
			if v.Serial != "22FF9503829A3A0DDE9CB87191A472D4" {
				t.Fatal()
			}

		}
	}
}

func TestGetClient(t *testing.T) {
	New([]config.Client{
		{
			Serial: "12FF9503829A3A0DDE9CB87191A472D4",
			Scopes: []string{"READ", "WRITE"},
		},
	})
	if _, err := GetClient("12FF9503829A3A0DDE9CB87191A472D4"); err != nil {
		t.Fatal(err)
	}
	if _, err := GetClient("00FF9503829A3A0DDE9CB87191A472D4"); err == nil {
		t.Fatal(err)
	}
}

func TestFromContext(t *testing.T) {
	if _, ok := FromContext(context.Background()); ok {
		t.Fatal()
	}
}

func TestNewContext(t *testing.T) {
	New([]config.Client{
		{
			Serial: "12FF9503829A3A0DDE9CB87191A472D4",
			Scopes: []string{"READ", "WRITE"},
		},
	})
	c := new(big.Int)
	c.SetString("12FF9503829A3A0DDE9CB87191A472D4", 16)
	ctx1 := NewContext(context.Background(), c)
	if _, ok := FromContext(ctx1); !ok {
		t.Fatal()
	}
	ctxbg := context.Background()
	ctx2 := NewContext(ctxbg, nil)
	if !reflect.DeepEqual(ctx2, ctxbg) {
		t.Fatal()
	}
}
