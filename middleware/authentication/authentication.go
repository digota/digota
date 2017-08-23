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
