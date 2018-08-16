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

package zookeeper

import (
	"github.com/digota/digota/config"
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
	if _, err := NewLocker(config.Locker{Address: []string{"notvalid"}}); err == nil {
		t.Fatal(err)
	}
	if l, err := NewLocker(config.Locker{Address: []string{"localhost"}}); err != nil {
		t.Fatal(err)
	} else {
		l.Close()
	}
}

func TestLock_Lock(t *testing.T) {

	l, err := NewLocker(config.Locker{Address: []string{"localhost"}})
	if err != nil {
		t.Fatal(err)
	}
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

	l, err := NewLocker(config.Locker{Address: []string{"localhost"}})

	if err != nil {
		t.Fatal(err)
	}

	if err := l.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestLock_TryLock(t *testing.T) {

	l, err := NewLocker(config.Locker{Address: []string{"localhost"}})
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	id := uuid.NewV4().String()
	unlock, err := l.Lock(&testObj{Id: id})
	if err != nil {
		t.Fatal(err)
	}
	defer unlock()

	if unlock, err := l.TryLock(&testObj{Id: id}, time.Second); err == nil {
		t.Fatal(err)
		unlock()
	}

	if unlock, err := l.TryLock(&testObj{Id: uuid.NewV4().String()}, time.Second); err != nil {
		t.Fatal(err)
	} else {
		unlock()
	}
}
