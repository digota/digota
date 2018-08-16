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

package validation

import (
	"testing"
)

type (
	ok struct {
		A string `validate:"required,eq=1"`
	}
	ok1 struct {
		A string `validate:"required,eq=1"`
		B string `validate:"required,eq=1"`
	}
	fail struct {
		A string `validate:"get"`
	}
)

func TestValidate(t *testing.T) {
	// should pass
	if err := Validate(&ok{A: "1"}); err != nil {
		t.FailNow()
	}
	// should fail
	if err := Validate(&ok1{A: "1", B: "y"}); err == nil {
		t.FailNow()
	}
	// should fail
	if err := Validate(&ok{A: "x"}); err == nil {
		t.FailNow()
	}
	// should panic
	defer func() {
		if r := recover(); r == nil {
			t.FailNow()
		}
	}()
	Validate(&fail{A: "x"})
}
