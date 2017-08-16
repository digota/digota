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

package sdk

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"time"
)

// ClientOpt additional connection information
type ClientOpt struct {
	Crt, Key, ServerName, CaCrt   string
	InsecureSkipVerify, WithBlock bool
}

// NewClient creates new grpc connection to the addr using the ClientOpt
func NewClient(addr string, opt *ClientOpt) (*grpc.ClientConn, error) {

	if opt == nil {
		return nil, fmt.Errorf("Empty ClientOpt")
	}

	// Load key pair
	certificate, err := tls.LoadX509KeyPair(opt.Crt, opt.Key)
	if err != nil {
		return nil, err
	}

	tlsConf := &tls.Config{
		ServerName:         opt.ServerName,
		Certificates:       []tls.Certificate{certificate},
		InsecureSkipVerify: opt.InsecureSkipVerify,
	}

	// Mutual handshake with private ca
	if opt.CaCrt != "" {
		certPool := x509.NewCertPool()
		bs, err := ioutil.ReadFile(opt.CaCrt)
		if err != nil {
			return nil, err
		}
		ok := certPool.AppendCertsFromPEM(bs)
		if !ok {
			return nil, fmt.Errorf("failed to append certs")
		}
		tlsConf.InsecureSkipVerify = false
		tlsConf.RootCAs = certPool
	}

	dailOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConf)),
	}

	if opt.WithBlock {
		dailOpts = append(dailOpts, grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	}

	// Dail and return connection
	conn, err := grpc.Dial(addr, dailOpts...)
	if err != nil {
		return nil, err
	}

	return conn, nil

}
