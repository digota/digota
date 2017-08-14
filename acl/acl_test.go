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
	"github.com/digota/digota/config"
	"golang.org/x/net/context"
	"math/big"
	"regexp"
	"testing"
)

const testBase = "^(.test.Test/)"
const readMethod = "/test.Test/READ"
const writeMethod = "/test.Test/WRITE"
const publicMethod = "/test.Test/PUBLIC"

func TestCanAccessMethod(t *testing.T) {

	accessMap = map[client.Scope][][]*regexp.Regexp{
		// All methods
		client.WildcardScope: {
			[]*regexp.Regexp{
				regexp.MustCompile("(.*)"),
			},
		},
		// Public methods
		client.PublicScope: {
			[]*regexp.Regexp{
				regexp.MustCompile(testBase + "PUBLIC"),
			},
		},
		// Write only methods
		client.WriteScope: {
			[]*regexp.Regexp{
				regexp.MustCompile(testBase + "WRITE"),
			},
		},
		// Read only methods
		client.ReadScope: {
			[]*regexp.Regexp{
				regexp.MustCompile(testBase + "READ"),
			},
		},
	}

	readWriteClientSerial := new(big.Int)
	readWriteClientSerial.SetString("A2FF9503829A3A0DDE9CB87191A472D4", 16)

	readClientSerial := new(big.Int)
	readClientSerial.SetString("B2FF9503829A3A0DDE9CB87191A472D4", 16)

	writeClientSerial := new(big.Int)
	writeClientSerial.SetString("C2FF9503829A3A0DDE9CB87191A472D4", 16)

	client.New([]config.Client{
		{
			Serial: "A2FF9503829A3A0DDE9CB87191A472D4",
			Scopes: []string{"READ", "WRITE"},
		},
		{
			Serial: "B2FF9503829A3A0DDE9CB87191A472D4",
			Scopes: []string{"READ"},
		},
		{
			Serial: "C2FF9503829A3A0DDE9CB87191A472D4",
			Scopes: []string{"WRITE"},
		},
	})

	// check readWrite client
	ctx := client.NewContext(context.Background(), readWriteClientSerial)
	if !CanAccessMethod(ctx, readMethod) {
		t.Fatal()
	}
	if !CanAccessMethod(ctx, writeMethod) {
		t.Fatal()
	}
	if !CanAccessMethod(ctx, publicMethod) {
		t.Fatal()
	}

	// check read client
	ctx1 := client.NewContext(context.Background(), readClientSerial)
	if !CanAccessMethod(ctx1, readMethod) {
		t.Fatal()
	}
	if CanAccessMethod(ctx1, writeMethod) {
		t.Fatal()
	}
	if !CanAccessMethod(ctx1, publicMethod) {
		t.Fatal()
	}

	// check write client
	ctx2 := client.NewContext(context.Background(), writeClientSerial)
	if CanAccessMethod(ctx2, readMethod) {
		t.Fatal()
	}
	if !CanAccessMethod(ctx2, writeMethod) {
		t.Fatal()
	}
	if !CanAccessMethod(ctx2, publicMethod) {
		t.Fatal()
	}

	// check guest client
	ctx3 := client.NewContext(context.Background(), nil)
	if CanAccessMethod(ctx3, readMethod) {
		t.Fatal()
	}
	if CanAccessMethod(ctx3, writeMethod) {
		t.Fatal()
	}
	if !CanAccessMethod(ctx2, publicMethod) {
		t.Fatal()
	}

}
