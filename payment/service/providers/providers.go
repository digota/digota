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
	"sync"

	"github.com/digota/digota/config"
	"github.com/digota/digota/payment/paymentpb"
	"github.com/digota/digota/payment/service/providers/internalTestOnly"
	"github.com/digota/digota/payment/service/providers/stripe"
	log "github.com/sirupsen/logrus"
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
