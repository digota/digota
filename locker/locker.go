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

package locker

import (
	"errors"
	"github.com/digota/digota/config"
	"github.com/digota/digota/locker/handlers/zookeeper"
	"github.com/digota/digota/storage/object"
	"time"
)

const (
	ZookeeperHandler handlerName = "zookeeper"
	DefaultTimeout               = time.Millisecond * 100
)

type (
	handlerName string
	Interface   interface {
		Close() error
		Lock(doc object.Interface) (func() error, error)
		TryLock(doc object.Interface, t time.Duration) (func() error, error)
	}
)

var handler Interface

func New(lockerConfig config.Locker) error {
	var err error
	switch handlerName(lockerConfig.Handler) {
	case ZookeeperHandler:
		handler, err = zookeeper.NewLocker(lockerConfig)
		return err
	default:
		return errors.New("Locker is not valid")
	}
}

func Handler() Interface {
	return handler
}
