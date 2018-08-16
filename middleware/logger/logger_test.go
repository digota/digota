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
