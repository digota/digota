package braintree

type VenmoAccountDetails struct {
	Token             string `xml:"token,omitempty"`
	Username          string `xml:"username,omitempty"`
	VenmoUserID       string `xml:"venmo-user-id,omitempty"`
	SourceDescription string `xml:"source-description,omitempty"`
	ImageURL          string `xml:"image-url,omitempty"`
}
