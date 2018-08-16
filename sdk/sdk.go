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
