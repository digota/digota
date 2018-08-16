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
