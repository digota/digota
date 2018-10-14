package braintree

import "encoding/base64"

type apiKey struct {
	env        Environment
	merchantID string
	publicKey  string
	privateKey string
}

func newAPIKey(env Environment, merchantID, publicKey, privateKey string) credentials {
	return apiKey{
		env:        env,
		merchantID: merchantID,
		publicKey:  publicKey,
		privateKey: privateKey,
	}
}

func (k apiKey) Environment() Environment {
	return k.env
}

func (k apiKey) MerchantID() string {
	return k.merchantID
}

func (k apiKey) AuthorizationHeader() string {
	auth := k.publicKey + ":" + k.privateKey
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
