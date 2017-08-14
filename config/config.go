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

type AppConfig struct {
	TLS     TLS               `yaml:"TLS"`
	Clients []Client          `yaml:"clients"`
	Payment []PaymentProvider `yaml:"payment"`
	Storage Storage           `yaml:"storage"`
	Locker  Locker            `yaml:"locker"`
}

type Client struct {
	Serial string   `yaml:"serial"`
	Scopes []string `yaml:"scopes"`
}

type TLS struct {
	Key   string `yaml:"key"`
	Crt   string `yaml:"crt"`
	CACrt string `yaml:"ca"`
}

type Storage struct {
	Handler  string   `yaml:"handler"`
	Address  []string `yaml:"address"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	Database string   `yaml:"database"`
}

type Locker struct {
	Handler string   ` yaml:"handler"`
	Address []string ` yaml:"address"`
}

type PaymentProvider struct {
	Provider   string `yaml:"provider"`
	Secret     string `yaml:"secret"`
	Live       bool   `yaml:"live"`
	MerchId    string `yaml:"merchId"`
	PublicKey  string `yaml:"pubkey"`
	PrivateKey string `yaml:"priKey"`
}

func LoadConfig(configPath string) (conf *AppConfig, err error) {
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(b, &conf)
	return
}
