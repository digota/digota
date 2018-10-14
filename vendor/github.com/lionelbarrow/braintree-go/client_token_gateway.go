package braintree

import "encoding/xml"

const clientTokenVersion = 2

type ClientTokenGateway struct {
	*Braintree
}

func (g *ClientTokenGateway) Generate() (string, error) {
	return g.generate(&ClientTokenRequest{
		Version: clientTokenVersion,
	})
}

func (g *ClientTokenGateway) GenerateWithCustomer(customerId string) (string, error) {
	return g.generate(&ClientTokenRequest{
		Version:    clientTokenVersion,
		CustomerID: customerId,
	})
}

func (g *ClientTokenGateway) generate(req *ClientTokenRequest) (string, error) {
	resp, err := g.execute("POST", "client_token", req)
	if err != nil {
		return "", err
	}
	switch resp.StatusCode {
	case 201:
		var b clientToken
		if err := xml.Unmarshal(resp.Body, &b); err != nil {
			return "", err
		}
		return b.ClientToken, nil
	}
	return "", &invalidResponseError{resp}
}
