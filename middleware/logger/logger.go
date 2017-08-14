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

package logger

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
)

func CodeToLevel(code codes.Code) log.Level {
	switch code {
	case codes.OK:
		return log.InfoLevel
	case codes.Canceled:
		fallthrough
	case codes.Unknown:
		fallthrough
	case codes.InvalidArgument:
		fallthrough
	case codes.DeadlineExceeded:
		fallthrough
	case codes.NotFound:
		fallthrough
	case codes.AlreadyExists:
		fallthrough
	case codes.PermissionDenied:
		fallthrough
	case codes.Unauthenticated:
		fallthrough
	case codes.ResourceExhausted:
		fallthrough
	case codes.FailedPrecondition:
		fallthrough
	case codes.Aborted:
		fallthrough
	case codes.OutOfRange:
		fallthrough
	case codes.Unimplemented:
		fallthrough
	case codes.Internal:
		fallthrough
	case codes.Unavailable:
		fallthrough
	case codes.DataLoss:
		return log.ErrorLevel
	default:
		return log.InfoLevel
	}
}
