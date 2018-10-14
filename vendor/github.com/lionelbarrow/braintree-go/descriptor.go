package braintree

type Descriptor struct {
	Name  string `xml:"name,omitempty"`
	Phone string `xml:"phone,omitempty"`
	URL   string `xml:"url,omitempty"`
}
