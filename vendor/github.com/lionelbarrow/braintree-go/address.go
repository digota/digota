package braintree

import (
	"encoding/xml"
	"time"
)

type Address struct {
	XMLName            xml.Name
	Id                 string     `xml:"id,omitempty"`
	CustomerId         string     `xml:"customer-id,omitempty"`
	FirstName          string     `xml:"first-name,omitempty"`
	LastName           string     `xml:"last-name,omitempty"`
	Company            string     `xml:"company,omitempty"`
	StreetAddress      string     `xml:"street-address,omitempty"`
	ExtendedAddress    string     `xml:"extended-address,omitempty"`
	Locality           string     `xml:"locality,omitempty"`
	Region             string     `xml:"region,omitempty"`
	PostalCode         string     `xml:"postal-code,omitempty"`
	CountryCodeAlpha2  string     `xml:"country-code-alpha2,omitempty"`
	CountryCodeAlpha3  string     `xml:"country-code-alpha3,omitempty"`
	CountryCodeNumeric string     `xml:"country-code-numeric,omitempty"`
	CountryName        string     `xml:"country-name,omitempty"`
	CreatedAt          *time.Time `xml:"created-at,omitempty"`
	UpdatedAt          *time.Time `xml:"updated-at,omitempty"`
}

type Addresses struct {
	XMLName string     `xml:"addresses"`
	Address []*Address `xml:"address"`
}
