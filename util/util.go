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
