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

package client

import (
	"github.com/digota/digota/config"
	"github.com/digota/digota/util"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"math/big"
)

const (
	// WildcardScope represents access to all methods
	WildcardScope Scope = "WILDCARD"
	// PublicScope represents access to all public methods
	PublicScope Scope = "PUBLIC"
	// WriteScope represents access to all write methods
	WriteScope Scope = "WRITE"
	// ReadScope represents access to all read methods
	ReadScope Scope = "READ"
)

type (
	// Client serial and scopes to determine if certain client can access certain method
	Client struct {
		Serial string
		Scopes []Scope
	}
	// Scope represents the level of access to various methods
	Scope string
)

type clientKey struct{}

var clients []Client

// New initiate client slice based on the []config.Client slice
func New(c []config.Client) {
	for _, v := range c {
		var scopes []Scope
		for _, scope := range v.Scopes {
			scopes = append(scopes, Scope(scope))
		}
		clients = append(clients, Client{
			Serial: v.Serial,
			Scopes: scopes,
		})
	}
}

// NewContext store user in ctx and return new ctx.
func NewContext(ctx context.Context, serialId *big.Int) context.Context {
	var c *Client
	var err error
	if c, err = GetClient(util.BigIntToHex(serialId)); err != nil {
		return ctx
	}
	return context.WithValue(ctx, clientKey{}, c)
}

// FromContext returns the User stored in ctx, if any.
func FromContext(ctx context.Context) (*Client, bool) {
	u, ok := ctx.Value(clientKey{}).(*Client)
	return u, ok
}

// GetClient return client based on provided serialId
func GetClient(serialId string) (*Client, error) {
	// search for user
	for _, c := range clients {
		if c.Serial == serialId {
			return &c, nil
		}
	}
	return nil, errors.New("Cant find client.")
}
