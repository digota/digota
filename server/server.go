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

package server

// load services first
import (
	// register order service
	_ "github.com/digota/digota/order/service"
	// register payment service
	_ "github.com/digota/digota/payment/service"
	// register product service
	_ "github.com/digota/digota/product/service"
	// register sku service
	_ "github.com/digota/digota/sku/service"
)

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/digota/digota/client"
	"github.com/digota/digota/config"
	"github.com/digota/digota/locker"
	"github.com/digota/digota/middleware/authentication"
	"github.com/digota/digota/middleware/logger"
	"github.com/digota/digota/middleware/recovery"
	"github.com/digota/digota/order"
	"github.com/digota/digota/payment"
	"github.com/digota/digota/payment/service/providers"
	"github.com/digota/digota/product"
	"github.com/digota/digota/sku"
	"github.com/digota/digota/storage"
	"github.com/digota/digota/util"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type server struct {
	listener   net.Listener
	grpcServer *grpc.Server
}

// New create new digota server
func New(addr string, conf *config.AppConfig) *server {

	// create new storage handler
	if err := storage.New(conf.Storage); err != nil {
		log.Fatalf("Could not create storage handler => %s", err.Error())
	}

	// create new locker handler
	if err := locker.New(conf.Locker); err != nil {
		log.Fatalf("Could not create locker handler => %s", err.Error())
	}

	// load ca clients
	client.New(conf.Clients)
	providers.New(conf.Payment)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Infof("Listening on %s", addr)
	return &server{
		listener:   lis,
		grpcServer: newGRPCServer(conf),
	}
}

func getTlsOption(appConfig *config.AppConfig) grpc.ServerOption {

	// Load the certificates from disk
	certificate, err := tls.LoadX509KeyPair(appConfig.TLS.Crt, appConfig.TLS.Key)
	if err != nil {
		log.Panicf("could not load server key pair: %s", err)
	}

	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile(appConfig.TLS.CACrt)
	if err != nil {
		log.Fatalf("failed to read client ca cert: %s", err)
	}

	if ok := certPool.AppendCertsFromPEM(bs); !ok {
		log.Fatal("failed to append client certs")
	}

	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			for _, v := range appConfig.Clients {
				if v.Serial == util.BigIntToHex(verifiedChains[0][0].SerialNumber) {
					return nil
				}
			}
			return status.Error(codes.Unauthenticated, "")
		},
	}
	return grpc.Creds(credentials.NewTLS(tlsConfig))
}

func newGRPCServer(appConfig *config.AppConfig) *grpc.Server {
	// create new server with StreamServerInterceptors
	s := grpc.NewServer(
		// TLS with your certs
		getTlsOption(appConfig),
		// StreamInterceptor
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			// recover from panics
			grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandler(recovery.RecoveryHandlerFunc)),
			// will block Unauthenticated requests
			authentication.StreamServerInterceptor(),
			// logger interceptor
			grpc_logrus.StreamServerInterceptor(log.NewEntry(log.New()), grpc_logrus.WithLevels(logger.CodeToLevel)),
		)),
		// UnaryInterceptor
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			// recover from panics
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(recovery.RecoveryHandlerFunc)),
			// will block Unauthenticated requests
			authentication.UnaryServerInterceptor(),
			// logger interceptor
			grpc_logrus.UnaryServerInterceptor(log.NewEntry(log.New()), grpc_logrus.WithLevels(logger.CodeToLevel)),
		)),
	)
	registerServices(s)
	return s
}

func registerServices(s *grpc.Server) {
	product.RegisterProductServer(s)
	order.RegisterOrderServer(s)
	payment.RegisterPaymentServer(s)
	sku.RegisterSkuServer(s)
	reflection.Register(s)
}

func (s *server) Run() {
	// graceful stop on Interrupt
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range c {
			// sig is a ^C, handle it
			log.Warnf("Sig %s => GracefulStop triggered", sig.String())
			s.grpcServer.GracefulStop()
		}
	}()
	if err := s.grpcServer.Serve(s.listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
