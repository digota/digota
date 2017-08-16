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

package recovery

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime/debug"
)

// RecoveryHandlerFunc outputs the stack and the error
// returns valid grpc error
func RecoveryHandlerFunc(p interface{}) (err error) {
	// print stack to stderr
	debug.PrintStack()
	// return grpc error

	switch x := p.(type) {
	case *logrus.Entry:
		return status.Errorf(codes.Internal, "%s", x.Message)
	}

	return status.Errorf(codes.Internal, "%s", p)
}
