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

package storage

import (
	"github.com/digota/digota/config"
	"github.com/digota/digota/storage/object"
	"testing"
)

type dummyStorage struct{}

func (d *dummyStorage) Prepare() error                                               { return nil }
func (d *dummyStorage) Close() error                                                 { return nil }
func (d *dummyStorage) DropCollection(db string, doc object.Interface) error         { return nil }
func (d *dummyStorage) DropDatabase(db string) error                                 { return nil }
func (d *dummyStorage) One(doc object.Interface) error                               { return nil }
func (d *dummyStorage) List(docs object.Interfaces, opt object.ListOpt) (int, error) { return 0, nil }
func (d *dummyStorage) ListParent(parent string, docs object.Interfaces) error       { return nil }
func (d *dummyStorage) Insert(doc object.Interface) error                            { return nil }
func (d *dummyStorage) Update(doc object.Interface) error                            { return nil }
func (d *dummyStorage) Remove(doc object.Interface) error                            { return nil }

func TestNew(t *testing.T) {

	if err := New(config.Storage{
		Address: []string{"localhost"},
		Handler: "mongodb",
	}); err != nil {
		t.Fatal(err)
	}

	if err := New(config.Storage{
		Address: []string{"localhost"},
		Handler: "mongodbnotvalid",
	}); err == nil {
		t.Fatal(err)
	}

}

func TestHandler(t *testing.T) {
	handler = &dummyStorage{}
	if _, ok := Handler().(Interface); !ok {
		t.Fatal()
	}
}
