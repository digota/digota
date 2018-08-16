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

package memlock

import (
	"errors"
	"github.com/digota/digota/storage/object"
	"sync"
	"time"
)

const separator = "-"

var mtx sync.Mutex

type locker struct {
	smap map[string]*semaphore
}

// NewLocker return new lock
func NewLocker() *locker {
	return &locker{
		smap: make(map[string]*semaphore),
	}
}

func (m *locker) Close() error {
	m.smap = make(map[string]*semaphore)
	return nil
}

func (m *locker) getSemaphore(key string) *semaphore {
	mtx.Lock()
	defer mtx.Unlock()
	v, ok := m.smap[key]
	if !ok {
		m.smap[key] = newSemaphore()
		v = m.smap[key]
	}
	return v
}

func (m *locker) Lock(doc object.Interface) (func() error, error) {

	key, err := getKey(doc)

	if err != nil {
		return nil, err
	}

	s := m.getSemaphore(key)

	s.lock()

	return func() error {
		s.unlock()
		return nil
	}, nil

}

func (m *locker) TryLock(doc object.Interface, timeout time.Duration) (func() error, error) {

	key, err := getKey(doc)

	if err != nil {
		return nil, err
	}

	s := m.getSemaphore(key)

	if !s.tryLock(timeout) {
		return nil, errors.New("tryLock timeout")
	}

	return func() error {
		s.unlock()
		return nil
	}, nil

}

func getKey(doc object.Interface) (string, error) {
	if doc.GetId() == "" || doc.GetNamespace() == "" {
		return "", errors.New("Obj is missing information to make that lock")
	}
	return doc.GetNamespace() + separator + doc.GetId(), nil
}

func newSemaphore() *semaphore {
	s := make(semaphore, 1)
	return &s
}

type semaphore chan struct{}

func (s semaphore) lock() {
	s <- struct{}{}
}

func (s semaphore) unlock() {
	<-s
}

func (s semaphore) tryLock(timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case s <- struct{}{}:
		return true
	case <-time.After(timeout):
	}
	return false
}
