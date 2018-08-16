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
	"errors"
	"github.com/digota/digota/config"
	"github.com/digota/digota/storage/handlers/mongo"
	"github.com/digota/digota/storage/object"
)

const (
	mongodbHandler handlerName = "mongodb"
)

type (
	handlerName string
	// Interface defines the base functionality which any storage handler
	// should implement to become valid storage handler
	Interface interface {
		Prepare() error
		Close() error
		DropCollection(db string, doc object.Interface) error
		DropDatabase(db string) error
		One(doc object.Interface) error
		List(docs object.Interfaces, opt object.ListOpt) (int, error)
		ListParent(parent string, docs object.Interfaces) error
		Insert(doc object.Interface) error
		Update(doc object.Interface) error
		Remove(doc object.Interface) error
	}
)

var handler Interface

// New creates storage handler from config.Storage and prepare it for use
// returns error if something went wrong during the preparations
func New(storageConfig config.Storage) error {
	// create handler based on the storage config
	switch handlerName(storageConfig.Handler) {
	case mongodbHandler:
		handler = mongo.NewHandler(storageConfig)
	default:
		return errors.New("Invalid storage handler `" + storageConfig.Handler + "`")
	}
	// prepare handler
	return handler.Prepare()
}

// Handler returns the registered storage handler
func Handler() Interface {
	return handler
}
