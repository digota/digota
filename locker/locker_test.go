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

package locker

import (
	"github.com/digota/digota/config"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	err := New(config.Locker{
		Handler: "zookeeper",
		Address: []string{"localhost"},
	})
	if err != nil {
		t.Fatal(err)
	}
	handler = nil
	New(config.Locker{Handler: "not-valid"})
	if reflect.TypeOf(handler).String() != "*memlock.locker" {
		t.Fatal()
	}
}

func TestHandler(t *testing.T) {
	handler = nil
	err := New(config.Locker{
		Handler: "zookeeper",
		Address: []string{"localhost"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if reflect.TypeOf(Handler()).String() != "*zookeeper.lock" {
		t.Fatal()
	}
}
