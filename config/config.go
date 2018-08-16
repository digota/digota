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
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// AppConfig is the main config structure
type AppConfig struct {
	TLS     TLS               `yaml:"TLS"`
	Clients []Client          `yaml:"clients"`
	Payment []PaymentProvider `yaml:"payment"`
	Storage Storage           `yaml:"storage"`
	Locker  Locker            `yaml:"locker"`
}

// Client is the client config structure
// accept serial and scopse
//	clients:
//	- serial: "A2FF9503829A3A0DDE9CB87191A472D4"
//	scopes:
//	- READ
//	- WRITE
type Client struct {
	Serial string   `yaml:"serial"`
	Scopes []string `yaml:"scopes"`
}

// TLS is the tls config for running the server
//	TLS:
//	crt: out/server.com.crt
//	key: out/server.com.key
//	ca: out/ca.crt
type TLS struct {
	Key   string `yaml:"key"`
	Crt   string `yaml:"crt"`
	CACrt string `yaml:"ca"`
}

// Storage is the storage handler config
//	storage:
//	handler: mongodb
//	address:
//	- localhost
type Storage struct {
	Handler  string   `yaml:"handler"`
	Address  []string `yaml:"address"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	Database string   `yaml:"database"`
}

// Locker is the lock server handler config
//	locker:
//	handler: zookeeper
//	address:
//	- localhost
type Locker struct {
	Handler string   ` yaml:"handler"`
	Address []string ` yaml:"address"`
}

// PaymentProvider is the payment provider config
//	payment:
//	- provider: Stripe
//	secret: sk_test_0000000000000000000
type PaymentProvider struct {
	Provider   string `yaml:"provider"`
	Secret     string `yaml:"secret"`
	Live       bool   `yaml:"live"`
	MerchId    string `yaml:"merchId"`
	PublicKey  string `yaml:"pubkey"`
	PrivateKey string `yaml:"priKey"`
}

// LoadConfig read file from provided path and returns *AppConfig or error
func LoadConfig(configPath string) (conf *AppConfig, err error) {
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(b, &conf)
	return
}
