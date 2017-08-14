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

package mongo

import (
	"github.com/digota/digota/config"
	"github.com/digota/digota/storage/object"
	"github.com/satori/go.uuid"
	"log"
	"reflect"
	"testing"
)

type testParentObjs []*testParentObj

func (o *testParentObjs) GetNamespace() string {
	return "mongo_test"
}

type testParentObj struct {
	Id     string `bson:"_id"`
	Parent string
}

func (o *testParentObj) GetNamespace() string {
	return "mongo_test"
}

func (o *testParentObj) GetId() string {
	return o.Id
}

func (o *testParentObj) SetId(id string) {
	o.Id = id
}

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

type testObjWithTimeTracker struct {
	Id      string `bson:"_id"`
	Updated int64
	Created int64
	Data    string
}

func (o *testObjWithTimeTracker) GetNamespace() string {
	return "mongo_test"
}

func (o *testObjWithTimeTracker) GetId() string {
	return o.Id
}

func (o *testObjWithTimeTracker) SetId(id string) {
	o.Id = id
}

func (o *testObjWithTimeTracker) GetCreated() int64 {
	return o.Created
}

func (o *testObjWithTimeTracker) SetCreated(t int64) {
	o.Created = t
}

func (o *testObjWithTimeTracker) GetUpdated() int64 {
	return o.Updated
}

func (o *testObjWithTimeTracker) SetUpdated(t int64) {
	o.Updated = t
}

func TestNewHandler(t *testing.T) {

	iface := NewHandler(config.Storage{
		Address: []string{"localhost"},
	})

	if reflect.TypeOf(iface).String() != "*mongo.handler" {
		t.Fatal()
	}

	iface.Close()

}

func TestHandler_Prepare(t *testing.T) {

	db := uuid.NewV4().String()

	iface := NewHandler(config.Storage{
		Address:  []string{"localhost"},
		Database: db,
	})

	defer func() {
		iface.DropDatabase(db)
		iface.Close()
	}()

	if reflect.TypeOf(iface).String() != "*mongo.handler" {
		t.Fatal()
	}

	if err := iface.Prepare(); err != nil {
		t.Fatal(err)
	}

	iface2 := NewHandler(config.Storage{
		Address: []string{"notvalidhost"},
	})

	if err := iface2.Prepare(); err == nil {
		t.Fatal(err)
	}

}

func TestHandler_Close(t *testing.T) {

	iface := NewHandler(config.Storage{
		Address: []string{"localhost"},
	})

	// close before prepare
	if err := iface.Close(); err == nil {
		t.Fatal(err)
	}

	// clean close
	if err := iface.Prepare(); err != nil {
		t.Fatal(err)
	}

	// close before prepare
	if err := iface.Close(); err != nil {
		t.Fatal(err)
	}

}

func TestHandler_One(t *testing.T) {

	db := uuid.NewV4().String()

	iface := NewHandler(config.Storage{
		Address:  []string{"localhost"},
		Database: db,
	})

	defer func() {
		log.Println(iface.DropDatabase(db))
		iface.Close()
	}()

	if err := iface.Prepare(); err != nil {
		t.Fatal(err)
	}

	obj := &testObj{
		Id: "not_found",
	}

	if err := iface.One(obj); err == nil {
		t.Fatal(err)
	}

	if err := iface.Insert(obj); err != nil {
		t.Fatal(err)
	}

	if err := iface.One(obj); err != nil {
		t.Fatal(err)
	}

}

func TestHandler_Insert(t *testing.T) {

	db := uuid.NewV4().String()

	iface := NewHandler(config.Storage{
		Address:  []string{"localhost"},
		Database: db,
	})

	defer func() {
		iface.DropDatabase(db)
		iface.Close()
	}()

	if err := iface.Prepare(); err != nil {
		t.Fatal(err)
	}

	obj := &testObj{
		Id: "not_found",
	}

	if err := iface.One(obj); err == nil {
		t.Fatal(err)
	}

	if err := iface.Insert(obj); err != nil {
		t.Fatal(err)
	}

	if err := iface.One(obj); err != nil {
		t.Fatal(err)
	}

	obj1 := &testObjWithTimeTracker{
		Id: uuid.NewV4().String(),
	}

	if err := iface.Insert(obj1); err != nil {
		t.Fatal(err)
	}

}

func TestHandler_Update(t *testing.T) {

	db := uuid.NewV4().String()

	iface := NewHandler(config.Storage{
		Address:  []string{"localhost"},
		Database: db,
	})

	defer func() {
		iface.DropDatabase(db)
		iface.Close()
	}()

	if err := iface.Prepare(); err != nil {
		t.Fatal(err)
	}

	obj := &testObjWithTimeTracker{
		Id:   uuid.NewV4().String(),
		Data: "beforeUpdate",
	}

	if err := iface.Insert(obj); err != nil {
		t.Fatal(err)
	}

	obj.Data = "afterUpdate"

	if err := iface.Update(obj); err != nil {
		t.Fatal(err)
	}

	if err := iface.One(obj); err != nil || obj.Data != "afterUpdate" || obj.Updated == 0 {
		t.Fatal(err)
	}

	obj.Id = ""

	if err := iface.Update(obj); err == nil {
		t.Fatal(err)
	}

}

func TestHandler_Remove(t *testing.T) {

	db := uuid.NewV4().String()

	iface := NewHandler(config.Storage{
		Address:  []string{"localhost"},
		Database: db,
	})

	defer func() {
		iface.DropDatabase(db)
		iface.Close()
	}()

	if err := iface.Prepare(); err != nil {
		t.Fatal(err)
	}

	obj := &testObj{
		Id: uuid.NewV4().String(),
	}

	if err := iface.Insert(obj); err != nil {
		t.Fatal(err)
	}

	if err := iface.Remove(obj); err != nil {
		t.Fatal(err)
	}

	if err := iface.Remove(obj); err == nil {
		t.Fatal(err)
	}

}

func TestHandler_ListParent(t *testing.T) {

	db := uuid.NewV4().String()

	iface := NewHandler(config.Storage{
		Address:  []string{"localhost"},
		Database: db,
	})

	defer func() {
		iface.DropDatabase(db)
		iface.Close()
	}()

	if err := iface.Prepare(); err != nil {
		t.Fatal(err)
	}

	parent := uuid.NewV4().String()

	for k := 0; k < 10; k++ {
		err := iface.Insert(&testParentObj{
			Id:     uuid.NewV4().String(),
			Parent: parent,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	slice := &testParentObjs{}

	err := iface.ListParent(parent, slice)

	if err != nil {
		t.Fatal(err)
	}

	if len(*slice) != 10 {
		t.Fatal()
	}

}

func TestHandler_List(t *testing.T) {

	db := uuid.NewV4().String()

	iface := NewHandler(config.Storage{
		Address:  []string{"localhost"},
		Database: db,
	})

	defer func() {
		iface.DropDatabase(db)
		iface.Close()
	}()

	if err := iface.Prepare(); err != nil {
		t.Fatal(err)
	}

	parent := uuid.NewV4().String()

	for k := 0; k < 10; k++ {
		err := iface.Insert(&testParentObj{
			Id:     uuid.NewV4().String(),
			Parent: parent,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	slice := &testParentObjs{}

	for _, v := range []object.Sort{
		object.SortNatural,
		object.SortCreatedDesc,
		object.SortCreatedAsc,
		object.SortUpdatedDesc,
		object.SortUpdatedAsc,
	} {

		n, err := iface.List(slice, object.ListOpt{
			Limit: 10,
			Page:  0,
			Sort:  v,
		})

		if err != nil || n != 10 {
			t.Fatal(err)
		}
	}

}

func TestHandler_DropCollection(t *testing.T) {

	db := uuid.NewV4().String()

	iface := NewHandler(config.Storage{
		Address:  []string{"localhost"},
		Database: db,
	})

	defer func() {
		iface.DropDatabase(db)
		iface.Close()
	}()

	if err := iface.Prepare(); err != nil {
		t.Fatal(err)
	}

	obj := &testObj{
		Id: uuid.NewV4().String(),
	}

	if err := iface.Insert(obj); err != nil {
		t.Fatal(err)
	}

	if err := iface.DropCollection(db, obj); err != nil {
		t.Fatal(err)
	}

}
