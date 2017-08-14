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
