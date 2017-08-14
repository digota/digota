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
	"testing"
)

type testCase struct {
	code  codes.Code
	level log.Level
}

var testCases = []testCase{
	{
		codes.OK,
		log.InfoLevel,
	},
	{
		codes.Canceled,
		log.ErrorLevel,
	},
	{
		codes.Unknown,
		log.ErrorLevel,
	},
	{
		codes.InvalidArgument,
		log.ErrorLevel,
	},
	{
		codes.DeadlineExceeded,
		log.ErrorLevel,
	},
	{
		codes.NotFound,
		log.ErrorLevel,
	},
	{
		codes.AlreadyExists,
		log.ErrorLevel,
	},
	{
		codes.PermissionDenied,
		log.ErrorLevel,
	},
	{
		codes.Unauthenticated,
		log.ErrorLevel,
	},
	{
		codes.ResourceExhausted,
		log.ErrorLevel,
	},
	{
		codes.FailedPrecondition,
		log.ErrorLevel,
	},
	{
		codes.Aborted,
		log.ErrorLevel,
	},
	{
		codes.OutOfRange,
		log.ErrorLevel,
	},
	{
		codes.Unimplemented,
		log.ErrorLevel,
	},
	{
		codes.Internal,
		log.ErrorLevel,
	},
	{
		codes.Unavailable,
		log.ErrorLevel,
	},
	{
		codes.Unimplemented,
		log.ErrorLevel,
	},
	{
		codes.DataLoss,
		log.ErrorLevel,
	},
}

func TestCodeToLevel(t *testing.T) {
	for _, v := range testCases {
		if CodeToLevel(v.code) != v.level {
			t.Fatal()
		}
	}
	// cover default
	if CodeToLevel(150) != log.InfoLevel {
		t.Fatal()
	}
}
