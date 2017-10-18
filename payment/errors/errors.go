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

package errors

// ErrorCode represents payment provider specific error
type ErrorCode string

var (
	// IncorrectNum incorrect cc num
	IncorrectNum ErrorCode = "incorrect_number"
	// InvalidNum invalid cc num
	InvalidNum ErrorCode = "invalid_number"
	// InvalidExpM invalid exp month
	InvalidExpM ErrorCode = "invalid_expiry_month"
	// InvalidExpY invalid exp year
	InvalidExpY ErrorCode = "invalid_expiry_year"
	// InvalidCvc invalid cvc number
	InvalidCvc ErrorCode = "invalid_cvc"
	// ExpiredCard card is expired
	ExpiredCard ErrorCode = "expired_card"
	// IncorrectCvc incorrect cvc
	IncorrectCvc ErrorCode = "incorrect_cvc"
	// IncorrectZip incorrect zip code
	IncorrectZip ErrorCode = "incorrect_zip"
	// CardDeclined card declined
	CardDeclined ErrorCode = "card_declined"
	// Missing missing information
	Missing ErrorCode = "missing"
	// ProcessingErr processing error
	ProcessingErr ErrorCode = "processing_error"
	// RateLimit reached call rate limit
	RateLimit ErrorCode = "rate_limit"
)

func (e ErrorCode) Error() string {
	return string(e)
}
