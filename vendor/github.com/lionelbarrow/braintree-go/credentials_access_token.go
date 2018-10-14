package braintree

import (
	"errors"
	"strings"
)

type accessToken struct {
	raw        string
	env        Environment
	merchantID string
}

func newAccessToken(accessTokenStr string) (credentials, error) {
	parts := strings.Split(accessTokenStr, "$")
	if len(parts) < 3 || parts[0] != "access_token" {
		return nil, errors.New("access token is not of expected format")
	}
	env, err := EnvironmentFromName(parts[1])
	if err != nil {
		return nil, errors.New("access token is for unsupported environment, " + err.Error())
	}
	t := accessToken{
		raw:        accessTokenStr,
		env:        env,
		merchantID: parts[2],
	}
	return t, nil
}

func (t accessToken) Environment() Environment {
	return t.env
}

func (t accessToken) MerchantID() string {
	return t.merchantID
}

func (t accessToken) AuthorizationHeader() string {
	return "Bearer " + t.raw
}
