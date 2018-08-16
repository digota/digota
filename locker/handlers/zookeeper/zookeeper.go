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

// watch dump every second:
// 				watch -n 1 -d '{ echo "dump"; sleep 1; } | telnet localhost 2181'

package zookeeper

import (
	"errors"
	"github.com/digota/digota/config"
	"github.com/digota/digota/storage/object"
	"github.com/yaronsumel/go-zookeeper/zk"
	"time"
)

const separator = "/"

type lock struct {
	*zk.Conn
}

// NewLocker return new lock
func NewLocker(lockerConfig config.Locker) (*lock, error) {
	c, _, err := zk.Connect(lockerConfig.Address, time.Millisecond*100)
	if err != nil {
		return nil, err
	}
	return &lock{c}, nil
}

func (l *lock) newLock(obj object.Interface) (*zk.Lock, error) {
	if obj.GetNamespace() == "" || obj.GetId() == "" {
		return nil, errors.New("Obj is missing information to make that lock")
	}
	return zk.NewLock(l.Conn, separator+obj.GetNamespace()+separator+obj.GetId(), zk.WorldACL(zk.PermAll)), nil
}

func (l *lock) Close() error {
	l.Conn.Close()
	return nil
}

func (l *lock) Lock(obj object.Interface) (func() error, error) {
	z, err := l.newLock(obj)
	if err != nil {
		return nil, err
	}
	if err := z.Lock(); err != nil {
		return nil, err
	}
	return func() error { return z.Unlock() }, nil
}

func (l *lock) TryLock(obj object.Interface, t time.Duration) (func() error, error) {
	z, err := l.newLock(obj)
	if err != nil {
		return nil, err
	}
	if err := z.TryLock(t); err != nil {
		return nil, err
	}
	return func() error { return z.Unlock() }, nil
}
