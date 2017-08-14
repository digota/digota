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

package object

const DefaultDatabase = "digota"

const (
	SortNatural Sort = iota
	SortCreatedDesc
	SortCreatedAsc
	SortUpdatedDesc
	SortUpdatedAsc
)

type (
	Sort int

	Interfaces interface {
		GetNamespace() string
	}

	// Interface very basic object interface
	Interface interface {
		GetNamespace() string
		GetId() string
	}

	// TimeTracker help storage handlers set created and updated time when needed.
	TimeTracker interface {
		SetCreated(t int64)
		GetCreated() int64
		SetUpdated(t int64)
		GetUpdated() int64
	}

	// IdSetter helps the storage handler creating new object with fresh uuid
	IdSetter interface {
		SetId(string)
	}

	ListOpt struct {
		Page  int64
		Limit int64
		Sort  Sort
	}
)
