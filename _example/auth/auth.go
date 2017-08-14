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
