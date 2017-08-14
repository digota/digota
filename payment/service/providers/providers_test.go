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
