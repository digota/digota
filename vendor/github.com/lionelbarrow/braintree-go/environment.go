package braintree

import "fmt"

type Environment struct {
	baseURL string
}

func NewEnvironment(baseURL string) Environment {
	return Environment{baseURL: baseURL}
}

func (e Environment) BaseURL() string {
	return e.baseURL
}

var (
	Development = NewEnvironment("http://localhost:3000")
	Sandbox     = NewEnvironment("https://api.sandbox.braintreegateway.com:443")
	Production  = NewEnvironment("https://api.braintreegateway.com:443")
)

func EnvironmentFromName(name string) (Environment, error) {
	switch name {
	case "development":
		return Development, nil
	case "sandbox":
		return Sandbox, nil
	case "production":
		return Production, nil
	}
	return Environment{}, fmt.Errorf("unknown environment %q", name)
}
