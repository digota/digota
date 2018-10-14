package braintree

type SubscriptionGateway struct {
	*Braintree
}

func (g *SubscriptionGateway) Create(sub *SubscriptionRequest) (*Subscription, error) {
	resp, err := g.execute("POST", "subscriptions", sub)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 201:
		return resp.subscription()
	}
	return nil, &invalidResponseError{resp}
}

func (g *SubscriptionGateway) Update(sub *SubscriptionRequest) (*Subscription, error) {
	resp, err := g.execute("PUT", "subscriptions/"+sub.Id, sub)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.subscription()
	}
	return nil, &invalidResponseError{resp}
}

func (g *SubscriptionGateway) Find(subId string) (*Subscription, error) {
	resp, err := g.execute("GET", "subscriptions/"+subId, nil)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.subscription()
	}
	return nil, &invalidResponseError{resp}
}

func (g *SubscriptionGateway) Cancel(subId string) (*Subscription, error) {
	resp, err := g.execute("PUT", "subscriptions/"+subId+"/cancel", nil)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return resp.subscription()
	}
	return nil, &invalidResponseError{resp}
}
