package braintree

type AddOnList struct {
	XMLName string  `xml:"add-ons"`
	AddOns  []AddOn `xml:"add-on"`
}

type AddOn struct {
	XMLName string `xml:"add-on"`
	Modification
}
