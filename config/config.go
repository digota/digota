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
