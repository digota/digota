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
