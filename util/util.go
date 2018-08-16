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

package util

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/big"
	"time"
)

const (
	defaultRetryAttempts = 5
	defaultRetrySleep    = time.Millisecond * 100
)

// Fn generic closure function
type Fn func() error

// Retry function for default retry attempts
func Retry(callback func() error) (err error) {
	for i := 0; ; i++ {
		if err = callback(); err == nil {
			return
		}
		if i >= (defaultRetryAttempts - 1) {
			break
		}
		time.Sleep(defaultRetrySleep)
		log.Error("retrying after error:", err)
	}
	return fmt.Errorf("after %d attempts, last error: %s", defaultRetryAttempts, err)
}

// BigIntToHex convert big.Int to string
func BigIntToHex(b *big.Int) string {
	return fmt.Sprintf("%X", b)
}
