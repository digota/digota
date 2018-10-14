package braintree

import (
	"encoding/base64"
	"encoding/xml"
	"github.com/lionelbarrow/braintree-go/xmlnil"
	"net/http"
)

type WebhookNotificationGateway struct {
	*Braintree
	apiKey apiKey
}

func (w *WebhookNotificationGateway) ParseRequest(r *http.Request) (*WebhookNotification, error) {
	signature := r.PostFormValue("bt_signature")
	payload := r.PostFormValue("bt_payload")
	return w.Parse(signature, payload)
}

func (w *WebhookNotificationGateway) Parse(signature, payload string) (*WebhookNotification, error) {
	hmacer := newHmacer(w.apiKey.publicKey, w.apiKey.privateKey)
	if verified, err := hmacer.verifySignature(signature, payload); err != nil {
		return nil, err
	} else if !verified {
		return nil, SignatureError{}
	}

	xmlNotification, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, err
	}

	strippedBuf, err := xmlnil.StripNilElements(xmlNotification)
	if err != nil {
		return nil, err
	}

	var n WebhookNotification
	err = xml.Unmarshal(strippedBuf, &n)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (w *WebhookNotificationGateway) Verify(challenge string) (string, error) {
	hmacer := newHmacer(w.apiKey.publicKey, w.apiKey.privateKey)
	digest, err := hmacer.hmac(challenge)
	if err != nil {
		return ``, err
	}
	return hmacer.publicKey + `|` + digest, nil
}
