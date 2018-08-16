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

package main

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/digota/digota/sdk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
)

func main() {

	// You can use sdk client to make things easier
	log.Println(sdk.NewClient("localhost:3051", &sdk.ClientOpt{
		InsecureSkipVerify: false,
		ServerName:         "server.com",
		CaCrt:              "out/ca.crt",
		Crt:                "out/client.com.crt",
		Key:                "out/client.com.key",
	}))

	// Or.. you can do things alone

	// Without mutual handshake validation (skipInsecure)
	func() {
		certificate, err := tls.LoadX509KeyPair("out/client.com.crt", "out/client.com.key")
		if err != nil {
			panic(err)
		}
		tlsConf := &tls.Config{
			ServerName:         "",
			Certificates:       []tls.Certificate{certificate},
			InsecureSkipVerify: true,
		}
		log.Println(grpc.Dial("localhost:1234", grpc.WithTransportCredentials(credentials.NewTLS(tlsConf))))
	}()

	// With mutual handshake validation using the ca certificate
	func() {
		certificate, err := tls.LoadX509KeyPair("out/client.com.crt", "out/client.com.key")
		if err != nil {
			panic(err)
		}
		tlsConf := &tls.Config{
			ServerName:         "",
			Certificates:       []tls.Certificate{certificate},
			InsecureSkipVerify: false,
		}
		certPool := x509.NewCertPool()
		bs, err := ioutil.ReadFile("out/ca.crt")
		if err != nil {
			panic(err)
		}
		ok := certPool.AppendCertsFromPEM(bs)
		if !ok {
			panic("failed to append certs")
		}
		tlsConf.RootCAs = certPool
		log.Println(grpc.Dial("localhost:1234", grpc.WithTransportCredentials(credentials.NewTLS(tlsConf))))
	}()

}
