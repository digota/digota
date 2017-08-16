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

package authentication

import (
	"crypto/tls"
	"github.com/grpc-ecosystem/go-grpc-middleware/testing"
	pb_testproto "github.com/grpc-ecosystem/go-grpc-middleware/testing/testproto"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"testing"
)

var (
	goodPing = &pb_testproto.PingRequest{Value: "something", SleepTimeMs: 9999}
)

type assertingPingService struct {
	pb_testproto.TestServiceServer
	T *testing.T
}

type authOverrideTestService struct {
	pb_testproto.TestServiceServer
	T *testing.T
}

func TestAuthOverrideTestSuite(t *testing.T) {
	t.Skip()
	s := &AuthOverrideTestSuite{
		InterceptorTestSuite: &grpc_testing.InterceptorTestSuite{
			TestService: &authOverrideTestService{&assertingPingService{&grpc_testing.TestPingService{T: t}, t}, t},
			ServerOpts: []grpc.ServerOption{
				grpc.Creds(credentials.NewTLS(&tls.Config{
					MinVersion:         tls.VersionTLS12,
					InsecureSkipVerify: true,
				})),
				grpc.StreamInterceptor(StreamServerInterceptor()),
				grpc.UnaryInterceptor(UnaryServerInterceptor()),
			},
		},
	}
	suite.Run(t, s)
}

type AuthOverrideTestSuite struct {
	*grpc_testing.InterceptorTestSuite
}

func (s *AuthOverrideTestSuite) TestUnary_PassesAuth() {
	_, err := s.Client.Ping(context.Background(), goodPing)
	require.Error(s.T(), err, "error must occur")
}

func (s *AuthOverrideTestSuite) TestStream_PassesAuth() {
	stream, err := s.Client.PingList(context.Background(), goodPing)
	require.NoError(s.T(), err, "should not fail on establishing the stream")
	pong, err := stream.Recv()
	require.Error(s.T(), err, "error must occur")
	require.Nil(s.T(), pong, "pong must be nil")
}
