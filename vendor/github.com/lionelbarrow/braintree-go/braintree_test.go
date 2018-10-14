package braintree

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHTTPClientTimeout_New(t *testing.T) {
	t.Parallel()
	testHTTPClientTimeout(
		t,
		func(env Environment) *Braintree {
			return New(env, "mid", "pubkey", "privkey")
		},
		time.Second*60,
	)
}

func TestHTTPClientTimeout_NewWithHttpClient(t *testing.T) {
	t.Parallel()
	testHTTPClientTimeout(
		t,
		func(env Environment) *Braintree {
			return NewWithHttpClient(env, "mid", "pubkey", "privkey", &http.Client{Timeout: time.Second * 10})
		},
		time.Second*10,
	)
}

func TestHTTPClientTimeout_NewWithHttpClient_Nil(t *testing.T) {
	t.Parallel()
	testHTTPClientTimeout(
		t,
		func(env Environment) *Braintree {
			return NewWithHttpClient(env, "mid", "pubkey", "privkey", nil)
		},
		time.Second*60,
	)
}

func testHTTPClientTimeout(t *testing.T, braintreeFactory func(env Environment) *Braintree, expectedTimeout time.Duration) {
	const gracePeriod = time.Second * 10

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second*60 + gracePeriod)
	}))
	env := NewEnvironment(server.URL)

	b := braintreeFactory(env)

	finished := make(chan bool)
	go func() {
		_, err := b.Transaction().Create(&TransactionRequest{})
		if err == nil {
			t.Fatal("Expected timeout error, received no error")
		}
		if !strings.Contains(err.Error(), "Timeout") && !strings.Contains(err.Error(), "read tcp") {
			t.Fatalf("Expected timeout error, received: %s", err)
		}
		finished <- true
	}()

	select {
	case <-finished:
		t.Logf("Timeout received as expected")
	case <-time.After(expectedTimeout + gracePeriod):
		t.Fatalf("Timeout did not occur around %s, has been at least %s", expectedTimeout, expectedTimeout+gracePeriod)
	}
}
