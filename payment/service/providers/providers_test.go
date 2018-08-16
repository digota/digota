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

package providers

import (
	"github.com/digota/digota/config"
	"github.com/digota/digota/payment/paymentpb"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	New([]config.PaymentProvider{
		{
			Provider: "DigotaInternalTestOnly",
		},
	})
	if len(providers) != 1 {
		t.Fatal()
	}
	if providers[paymentpb.PaymentProviderId_Stripe].ProviderId() != paymentpb.PaymentProviderId_Stripe {
		t.Fatal()
	}
	func() {
		providers = make(map[paymentpb.PaymentProviderId]Interface)
		defer func() {
			if r := recover(); r == nil {
				t.Fatal()
			}
		}()
		New([]config.PaymentProvider{
			{
				Provider: "NA",
			},
		})
	}()
}

func TestProvider(t *testing.T) {
	providers = make(map[paymentpb.PaymentProviderId]Interface)
	New([]config.PaymentProvider{
		{
			Provider: "DigotaInternalTestOnly",
		},
	})
	if reflect.TypeOf(Provider(paymentpb.PaymentProviderId_Stripe)).String() != "*internalTestOnly.provider" {
		t.Fatal()
	}
	providers = make(map[paymentpb.PaymentProviderId]Interface)
	New([]config.PaymentProvider{
		{
			Provider: "Stripe",
		},
	})
	if reflect.TypeOf(Provider(paymentpb.PaymentProviderId_Stripe)).String() != "*stripe.provider" {
		t.Fatal()
	}
	func() {
		providers = make(map[paymentpb.PaymentProviderId]Interface)
		defer func() {
			if r := recover(); r == nil {
				t.Fatal()
			}
		}()
		Provider(paymentpb.PaymentProviderId_Stripe)
	}()

}
