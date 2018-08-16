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
