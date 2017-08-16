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
	"github.com/digota/digota/payment/service/providers/internalTestOnly"
	"github.com/digota/digota/payment/service/providers/stripe"
	log "github.com/sirupsen/logrus"
	"sync"
)

var mtx = sync.Mutex{}
var providers = make(map[paymentpb.PaymentProviderId]Interface)

// Interface defines the base functionality which any payment
// provider should implement to become valid payment provider
type Interface interface {
	ProviderId() paymentpb.PaymentProviderId
	Charge(req *paymentpb.ChargeRequest) (*paymentpb.Charge, error)
	Refund(chargeId string, amount uint64, currency paymentpb.Currency, reason paymentpb.RefundReason) (*paymentpb.Refund, error)
	SupportedCards() []paymentpb.CardType
}

// New creates several payment providers from the provided
// []config.PaymentProvider and saves them in provider map
// for further use. panics on the first error.
func New(paymentConfig []config.PaymentProvider) {

	// init the payment providers
	for _, v := range paymentConfig {
		var p Interface
		var err error
		switch v.Provider {
		case paymentpb.PaymentProviderId_Stripe.String():
			p, err = stripe.NewProvider(&v)
			if err != nil {
				log.Panic(err)
			}
			//case paymentpb.PaymentProviderId_Braintree.String():
			//	p, err = braintree.NewProvider(&v)
			//	if err == nil {
			//		break
			//	}
			//	fallthrough
			// just for testing
		case "DigotaInternalTestOnly":
			p, err = internalTestOnly.NewProvider()
			if err != nil {
				log.Panic(err)
			}
		default:
			log.Panicf("Payment provider %s is not valid", v.Provider)
		}
		mtx.Lock()
		providers[p.ProviderId()] = p
		mtx.Unlock()
	}
}

// Provider returns payment provider from the map using the provided id
// panics if the provider could not be found
func Provider(p paymentpb.PaymentProviderId) Interface {
	mtx.Lock()
	defer mtx.Unlock()
	if v, ok := providers[p]; ok {
		return v
	}
	log.Panicf("payment provider %s isn't registered", p.String())
	return nil
}
