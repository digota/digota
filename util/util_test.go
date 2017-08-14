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

package util

import (
	"errors"
	"math/big"
	"testing"
)

func TestRetry(t *testing.T) {
	if err := Retry(func() error { return errors.New("err") }); err == nil {
		t.FailNow()
	}
	if err := Retry(func() error { return nil }); err != nil {
		t.FailNow()
	}
}

func TestBigIntToHex(t *testing.T) {
	if BigIntToHex(big.NewInt(10)) != "A" {
		t.Fatal()
	}
}
