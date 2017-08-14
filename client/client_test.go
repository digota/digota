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
