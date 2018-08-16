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
	"github.com/digota/digota/acl"
	"github.com/digota/digota/client"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor returns a new unary server interceptors that performs per-request auth.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if acl.SkipAuth() {
			return handler(ctx, req)
		}
		//
		peer, ok := peer.FromContext(ctx)
		if !ok {
			return nil, status.Error(codes.PermissionDenied, "")
		}
		// get tlsInfo from peer
		tlsInfo, ok := peer.AuthInfo.(credentials.TLSInfo)
		if !ok || len(tlsInfo.State.VerifiedChains) == 0 {
			return nil, status.Error(codes.PermissionDenied, "")
		}
		ctx = client.NewContext(ctx, tlsInfo.State.VerifiedChains[0][0].SerialNumber)
		// the only Unauthenticated method is Auth.Auth
		// to access other methods JWT should have been set to context
		if !acl.CanAccessMethod(ctx, info.FullMethod) {
			//log.Infof("User %+v does not have permission to execute %s ", user.FromContext(ctx),info.FullMethod)
			return nil, status.Errorf(codes.PermissionDenied, "User does not have permission to execute %s ", info.FullMethod)
		}
		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a new unary server interceptors that performs per-request auth.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		wrapped := grpc_middleware.WrapServerStream(stream)
		if acl.SkipAuth() {
			return handler(srv, wrapped)
		}
		//
		peer, ok := peer.FromContext(ctx)
		if !ok {
			return status.Error(codes.PermissionDenied, "")
		}
		// get tlsInfo from peer
		tlsInfo, ok := peer.AuthInfo.(credentials.TLSInfo)
		if !ok || len(tlsInfo.State.VerifiedChains) == 0 {
			return status.Error(codes.PermissionDenied, "")
		}
		ctx = client.NewContext(ctx, tlsInfo.State.VerifiedChains[0][0].SerialNumber)
		// the only Unauthenticated method is Auth.Auth
		// to access other methods JWT should have been set to context
		if !acl.CanAccessMethod(ctx, info.FullMethod) {
			//log.Infof("User %+v does not have permission to execute %s ", user.FromContext(ctx),info.FullMethod)
			return status.Errorf(codes.PermissionDenied, "User does not have permission to execute %s ", info.FullMethod)
		}
		wrapped.WrappedContext = ctx
		return handler(srv, wrapped)
	}
}
