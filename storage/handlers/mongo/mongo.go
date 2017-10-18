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
	"fmt"
	"github.com/digota/digota/config"
	"github.com/digota/digota/storage/object"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"sync"
	"time"
)

type handler struct {
	client   *mgo.Session
	dailInfo *mgo.DialInfo
	database string
}

// NewHandler create new mongo handler
func NewHandler(s config.Storage) *handler {
	if s.Database == "" {
		s.Database = object.DefaultDatabase
	}
	h := &handler{
		dailInfo: &mgo.DialInfo{
			PoolLimit: 4096,
			Timeout:   time.Second,
			FailFast:  true,
			Username:  s.Username,
			Password:  s.Password,
			Addrs:     s.Address,
			Database:  s.Database,
		},
		database: s.Database,
	}
	return h
}

// Prepare
func (h *handler) Prepare() error {
	mgoSession, err := mgo.DialWithInfo(h.dailInfo)
	if err != nil {
		return err
	}
	// Switch the session to a monotonic behavior.
	mgoSession.SetMode(mgo.Monotonic, true)
	h.client = mgoSession

	//

	s := mgoSession.Clone()
	defer s.Close()

	// EnsureIndexes
	//log.Println(s.DB(h.database).C("sku").EnsureIndex(mgo.Index{
	//	Key:        []string{"product"},
	//	Unique:     true,
	//	DropDups:   true,
	//	Background: true, // See notes.
	//	Sparse:     true,
	//}))

	return nil
}

// Close
func (h *handler) Close() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Close err %s", r)
		}
	}()
	h.client.Close()
	return
}

func (h *handler) ListParent(parent string, obj object.Interfaces) (err error) {
	s := h.client.Clone()
	defer s.Close()
	return s.DB(h.database).C(obj.GetNamespace()).Find(bson.M{"parent": parent}).All(obj)
}

func (h *handler) List(obj object.Interfaces, opt object.ListOpt) (n int, err error) {

	wg := sync.WaitGroup{}
	wg.Add(2)

	// list
	go func() {
		defer wg.Done()
		s := h.client.Clone()
		defer s.Close()
		var msort string
		switch opt.Sort {
		case object.SortNatural:
			msort = "$natural"
		case object.SortCreatedDesc:
			msort = "-created"
		case object.SortCreatedAsc:
			msort = "+created"
		case object.SortUpdatedDesc:
			msort = "-updated"
		case object.SortUpdatedAsc:
			msort = "+updated"
		}
		err = s.DB(h.database).C(obj.GetNamespace()).Find(nil).Skip(int(opt.Page * opt.Limit)).Limit(int(opt.Limit)).Sort(msort).All(obj)
	}()

	// count
	go func() {
		defer wg.Done()
		s := h.client.Clone()
		defer s.Close()
		n, err = s.DB(h.database).C(obj.GetNamespace()).Count()
	}()

	wg.Wait()
	return
}

// ListIds
func (h *handler) DropCollection(db string, obj object.Interface) error {
	s := h.client.Clone()
	defer s.Close()
	return s.DB(h.database).C(obj.GetNamespace()).DropCollection()
}

func (h *handler) DropDatabase(db string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			err = fmt.Errorf("Close err %s", r)
		}
	}()
	s := h.client.Clone()
	defer s.Close()
	return s.DB(h.database).DropDatabase()
}

func (h *handler) One(obj object.Interface) error {

	s := h.client.New()

	defer s.Close()

	if err := s.DB(h.database).C(obj.GetNamespace()).Find(bson.D{bson.DocElem{Name: "_id", Value: obj.GetId()}}).One(obj); err != nil {
		return status.Errorf(codes.NotFound, "`%s::%s::%s`", obj.GetNamespace(), obj.GetId(), err.Error())
	}

	return nil

}

// Insert
func (h *handler) Insert(obj object.Interface) error {

	s := h.client.New()

	defer s.Close()

	if v, ok := obj.(object.IdSetter); ok {
		v.SetId(uuid.NewV4().String())
	}

	if v, ok := obj.(object.TimeTracker); ok {
		v.SetCreated(time.Now().Unix())
		v.SetUpdated(time.Now().Unix())
	}

	if err := s.DB(h.database).C(obj.GetNamespace()).Insert(obj); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil

}

func (h *handler) Update(obj object.Interface) error {

	s := h.client.Clone()

	defer s.Close()

	if v, ok := obj.(object.TimeTracker); ok {
		v.SetUpdated(time.Now().Unix())
	}

	if err := s.DB(h.database).C(obj.GetNamespace()).UpdateId(obj.GetId(), obj); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil

}

func (h *handler) Remove(obj object.Interface) error {

	s := h.client.Clone()

	defer s.Close()

	if err := s.DB(h.database).C(obj.GetNamespace()).Remove(bson.D{bson.DocElem{Name: "_id", Value: obj.GetId()}}); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil

}
