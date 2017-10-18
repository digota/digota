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

package memlock

import (
	"github.com/satori/go.uuid"
	"testing"
	"time"
)

type testObj struct {
	Id string `bson:"_id"`
}

func (o *testObj) GetNamespace() string {
	return "mongo_test"
}

func (o *testObj) GetId() string {
	return o.Id
}

func (o *testObj) SetId(id string) {
	o.Id = id
}

func TestNewLocker(t *testing.T) {
	l := NewLocker()
	l.Close()
}

func TestLock_Lock(t *testing.T) {
	l := NewLocker()
	defer l.Close()
	uuid := uuid.NewV4().String()
	unlock, err := l.Lock(&testObj{Id: uuid})
	if err != nil {
		t.Fatal(err)
	} else {
		unlock()
	}
	if _, err := l.Lock(&testObj{Id: ""}); err == nil {
		t.Fatal(err)
	}
}

func TestLock_Close(t *testing.T) {
	l := NewLocker()
	if err := l.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestLock_TryLock(t *testing.T) {
	l := NewLocker()
	defer l.Close()
	id := uuid.NewV4().String()
	unlock, err := l.Lock(&testObj{Id: id})
	if err != nil {
		t.Fatal(err)
	}
	defer unlock()
	if unlock, err := l.TryLock(&testObj{Id: id}, 1*time.Second); err == nil {
		t.Fatal(err)
		unlock()
	}
	if unlock, err := l.TryLock(&testObj{Id: uuid.NewV4().String()}, time.Second); err != nil {
		t.Fatal(err)
	} else {
		unlock()
	}
}
