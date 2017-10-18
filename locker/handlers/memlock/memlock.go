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
