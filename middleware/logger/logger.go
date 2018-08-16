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
)

// CodeToLevel translate grpc response code into log level (info/debug/warning/error)
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
