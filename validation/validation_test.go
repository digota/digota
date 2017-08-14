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
