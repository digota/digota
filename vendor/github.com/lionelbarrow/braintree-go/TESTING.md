### Setup

In order to run the integration tests, you have to configure your sandbox account accordingly, or some tests might fail.

#### Environment

Before running the integration tests with `go test -tags='unit integration' ./...`, make sure the following environment variables are set

```
export BRAINTREE_MERCH_ID={your-merchant-id}
export BRAINTREE_MERCH_ACCT_ID={your-merchant-account-id}
export BRAINTREE_PUB_KEY={your-public-key}
export BRAINTREE_PRIV_KEY={your-private-key}
```

You can also optional set an environment variable with the timezone of your sandbox account. The default is below but should match your account's timezone to ensure tests pass.

```
export BRAINTREE_TIMEZONE=UTC
```

When using Braintree Go in a production environment, we recommend that you continue to store these credentials in environment variables. See [the 12 Factor App](http://www.12factor.net/config) for details.

#### Sandbox settings

In your sandbox account go to `Settings > Processing` and enable the following

  1. `Card Verification`

In your sandbox account go to `Settings > Processing > CVV` and enable the following

  1. `CVV does not match (when provided) (N)` to `For Any Transaction`
  2. `CVV is not verified (when provided) (U)` to `For Any Transaction`

Your sandbox account will also need to be enabled for dynamic descriptors for the transaction dynamic descriptor tests to pass. See the [Braintree docs](https://articles.braintreepayments.com/control-panel/transactions/descriptors).

Finally you will also need to create a transaction with a specific id, two plans, an add-on, and a discount, with the the add on and discount associated with the first plan. Once you do all of these things, the integration tests should all pass, with the one exception listed below.

**Transactions**

The details don't matter, as long as its id is `dskdmb`.

**Test Plan 1 Setup**

```
Plan ID:                test_plan
Plan Name:              test_plan_name
Description:            test_plan_desc
Price:                  10
Currency:               USD

Include Trial Period:   YES
Duration:               14 days

Billing Cycle:          Every 1 Month
End Date:               After 2 billing cycles
```

**Test Plan 2 Setup**

```
Plan ID:                test_plan_2
Plan Name:              test_plan_2_name
Price:                  20
Currency:               USD

Billing Cycle:          Every 1 Month
First Bill Date:        Specific Day - Last Day of the Month
End Date:               Never
```

**Add-on setup**

```
Add-on ID:              test_add_on
Name:                   test_add_on_name
Description:            "A test add-on"
Amount:                 10
Number of cycles:       For the duration of the subscription.
```

**Discount setup**

```
Discount ID:            test_discount
Name:                   test_discount_name
Description             "A test discount"
Amount:                 10
Number of cycles:       For the duration of the subscription.
```

#### Testing disbursement details

In Braintree's sandbox environment, transactions do not disburse immediately. The implementation of disbursement details is tested in `disbursement_integration_test.go` using an already-disbursed transaction on the sandbox account associated with the Travis CI build. This test does not pass in other environments, and serves as a proof-of-concept. To make all of the tests pass, you will either need to comment out this test, or replace the values inside of it with the values from a disbursed transaction on your account. 
