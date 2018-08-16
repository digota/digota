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

package locker

import (
	"time"

	"github.com/digota/digota/config"
	"github.com/digota/digota/locker/handlers/memlock"
	"github.com/digota/digota/locker/handlers/redis"
	"github.com/digota/digota/locker/handlers/zookeeper"
	"github.com/digota/digota/storage/object"
)

const (
	zookeeperHandler handlerName = "zookeeper"
	redisHandler     handlerName = "redis"
	// lock acquire timeout
	// DefaultTimeout lock acquire timeout
	DefaultTimeout = time.Millisecond * 100
)

type (
	handlerName string
	// Interface is the base functionality that any locker handler
	// should implement in order to become valid handler
	Interface interface {
		Close() error
		Lock(doc object.Interface) (func() error, error)
		TryLock(doc object.Interface, t time.Duration) (func() error, error)
	}
)

var handler Interface

// New creates a locker handler from the provided
// config.Locker and save it in handler for further use
func New(lockerConfig config.Locker) error {
	var err error
	switch handlerName(lockerConfig.Handler) {
	case zookeeperHandler:
		handler, err = zookeeper.NewLocker(lockerConfig)
		return err
	case redisHandler:
		handler, err = redis.NewLocker(lockerConfig)
		return err
	default:
		handler = memlock.NewLocker()
	}
	return nil
}

// Handler returns the registered locker handler
func Handler() Interface {
	return handler
}
