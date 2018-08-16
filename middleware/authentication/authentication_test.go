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
