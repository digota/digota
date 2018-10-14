package braintree

import "testing"

func TestWebhookTestingGatewayRequest(t *testing.T) {
	t.Parallel()

	testingWebhookGateway := testGateway.WebhookTesting()
	webhookGateway := testGateway.WebhookNotification()

	kind := SubscriptionChargedSuccessfullyWebhook
	id := "123"

	r, err := testingWebhookGateway.Request(kind, id)
	if err != nil {
		t.Fatal(err)
	}

	err = r.ParseForm()
	if err != nil {
		t.Fatal(err)
	}
	payload := r.FormValue("bt_payload")
	signature := r.FormValue("bt_signature")

	notification, err := webhookGateway.Parse(signature, payload)

	if err != nil {
		t.Fatal(err)
	} else if notification.Kind != SubscriptionChargedSuccessfullyWebhook {
		t.Fatal("Incorrect Notification kind, expected subscription_charged_successfully got", notification.Kind)
	} else if notification.Subject.Subscription == nil {
		t.Log(notification.Subject)
		t.Fatal("Notification should have a subscription")
	} else if len(notification.Subject.Subscription.Transactions.Transaction) != 1 {
		t.Fatal("Incorrect number of subscription transactions, expected 1, got", len(notification.Subject.Subscription.Transactions.Transaction))
	} else if notification.Subject.Subscription.Transactions.Transaction[0].Id != "123" {
		t.Fatal("Incorrect transaction ID, expected '123' got", notification.Subject.Subscription.Transactions.Transaction[0].Id)
	}
}
