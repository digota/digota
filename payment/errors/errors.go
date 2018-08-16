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
