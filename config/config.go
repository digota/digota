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

package config

import (
	"github.com/kelseyhightower/envconfig"
)

// AppConfig is the main config structure
// export DIGOTA_LOCKER...=val
type AppConfig struct {
	TLS      TLS
	Clients  []Client
	Payment  []PaymentProvider
	Storage  Storage
	Locker   Locker
	Insecure bool
	Address  string
}

// Client is the client config structure
type Client struct {
	Serial string
	Scopes []string
}

// TLS is the tls config for running the server
//	TLS:
// export DIGOTA_CLIENT_SERIAL=val
// export DIGOTA_CLIENT_SCOPES=READ,WRITE
// export DIGOTA_CLIENT_SCOPES=READ,WRITE

type TLS struct {
	Key   string
	Crt   string
	CACrt string
}

// Storage is the storage handler config
type Storage struct {
	Handler  string
	Address  []string
	Username string
	Password string
	Database string
}

// Locker is the lock server handler config
type Locker struct {
	Handler string
	Address []string
}

// PaymentProvider is the payment provider config
type PaymentProvider struct {
	Provider   string
	Secret     string
	Live       bool
	MerchId    string
	PublicKey  string
	PrivateKey string
}

// LoadConfig read env vars and *AppConfig or error
func LoadConfig() (*AppConfig, error) {
	var (
		conf AppConfig
		err  error
	)
	err = envconfig.Process("digota", &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
