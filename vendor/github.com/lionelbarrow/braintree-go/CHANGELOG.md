## 0.20.0 (August 8th, 2018)

BACKWARDS INCOMPATIBILITES:

* Add subscription ID parameter to SubscriptionGateway Update. See #243.

DEPRECATED:
* TransactionGateway Search and SearchNext functions are deprecated, use SearchIDs and SearchPage instead.
* SubscriptionGateway Search and SearchNext functions are deprecated, use SearchIDs and SearchPage instead.
* CustomerGateway Search function is deprecated, use SearchIDs and SearchPage instead.

IMPROVEMENTS:

* Add CreatedAt, UpdatedAt to Subscription.
* Add StoreInVaultOnSuccess to TransactionOptions.
* Add PaymentInstrumentType constants to Transaction.
* Add SearchPage to TransactionGateway, deprecating Search and SearchNext.
* Add SearchPage to SubscriptionGateway, deprecating Search and SearchNext.
* Add ExpiringBetweenIDs, ExpiringBetweenPage to CreditCardGateway.
* Add SearchIDs, SearchPage to CustomerGateway.
* Add PaymentMethodNonce, PaymentMethodNonceDetails, PaymentMethodNonceGateway.
* Add 3DS support:
  * Add ThreeDSecure to TransactionOptions.
  * Add ThreeDSecureInfo to Transaction.

## 0.19.0 (June 8th, 2018)

BACKWARDS INCOMPATIBILITES:

* Rewrite error handling and interfaces. See #222 for details.

IMPROVEMENTS:

* Add TransactionLineItemGateway, TransactionLineItem, TransactionLineItemRequest, and add LineItems to TransactionRequest.
* Add DisputeGateway, Dispute, DisputeEvidence, DisputeTextEvidenceRequest, related types, and add Disputes to Transaction.
* Add support for validation errors on all fields. See #222 for details.
* Add RetryCharge to SubscriptionGateway.
* Add TransactionSource to TransactionRequest.
* Add CurrentBillingCycle to Modification.
* Add SubscriptionDetails to Transaction.
* Add StatusEvents to Subscription.
* Add SearchIDs, Search, and SearchNext to SubscriptionGateway.
* Add CreatedAt and UpdatedAt to Customer.
* Add comments to some exported functions, constants, and types.

INTERNAL:

* Fix VenmoSDK tests.
* Fix errcheck checks.
* Other minor test fixes.

## 0.18.0 (March 2nd, 2018)

BACKWARDS INCOMPATIBILITES:

* Require TLS 1.2 on outgoing connections.
* Replace Address with AddressRequest for creating Addresses.
* Replace Customer with CustomerRequest for creating and updating Customers.

IMPROVEMENTS:

* Add SearchIDs to TransactionGateway.
* Add SearchNext to TransactionGateway.
* Add SkipAdvancedFraudChecking to TransactionOptions.
* Add Update to AddressGateway.
* Add Addresses to Customer.

## 0.17.0 (January 5th, 2018)

IMPROVEMENTS:

* Add NumerOfBillingCycles to Modifications.
* Add PurchaseOrderNumber to Transaction and TransactionRequest.

INTERNAL:

* Remove omitempty from structs only used for deserialization.
* Add build tags for unit and integration tests.
* Add analysis detailed warning messages.
* Remove slow unit tests.

## 0.16.0 (November 11th, 2017)

BACKWARDS INCOMPATIBILITES:

* Support for Go1.6 has been dropped. Support from this version onwards is Go1.7+.
* Add `context.Context` as the first parameter to call functions that make an outgoing remote call.

IMPROVEMENTS:

* Add `go test` `-log` flag that logs requests and responses to stderr.

BUG FIXES:

* Responses not GZIP encoded were still attempted to be GZIP decoded.

## 0.15.0 (October 3rd, 2017)

BACKWARDS INCOMPATIBILITES:

* Change `VerifyCard` on `CreditCardOptions` and `PaymentMethodRequestOptions` from `bool` to `*bool`.

BUG FIXES:

* Setting `false` for the `VerifyCard` field on `CreditCardOptions` and `PaymentMethodRequestOptions` did not send a value to Braintree and would not turn off verification if verification was enabled account wide.

## 0.14.0 (September 20th, 2017)

IMPROVEMENTS:

* Add support for transaction cloning:
  * `Clone` function added to transaction gateway.
  * `TransactionCloneRequest` added.
  * `TransactionCloneOptions` added.
* Add support for escrow:
  * `EscrowStatus` added to `Transaction`.
  * `HoldInEscrow` added to `TransactionOptions`.
  * `CancelRelease` function added to transaction gateway.
  * `ReleaseFromEscrow` function added to transaction gateway.
  * `HoldInEscrow` function added to transaction gateway.
* `Channel` added to `Transaction`, `TransactionRequest`, and `TransactionCloneRequest`.

BUG FIXES:

* Test and CI reliability improvements to reduce flakiness of integration tests.

## 0.13.0 (September 1st, 2017)

BACKWARDS INCOMPATIBILITES:

* Nullable struct types were removed, with *NullBool/*NullInt64 replaced with *bool/*int and bool/int.

IMPROVEMENTS:

* `TransactionOptionsPaypalRequest` added to `TransactionOptions`.
* `TaxAmount` added to `Transaction` and `TransactionRequest`.
* `TaxExempt` added to `Transaction` and `TransactionRequest`.

BUG FIXES:
* Minor example typo fix.

## 0.12.0 (June 27th, 2017)

BACKWARDS INCOMPATIBILITES:

* Support for Go1.2-1.5 has been dropped. Support from this version onwards is Go1.6+.

IMPROVEMENTS:

* API endpoints updated to use the `api.` prefix: `sandbox.braintreepayments.com` => `api.sandbox.braintreegateway.com`, `www.braintreegateway.com` => `api.braintreegateway.com`.
* Add support for `VenmoAccount`s.
* Add support for `AndroidPayCard`s.
* Add support for `ApplePayCard`s.
* Added a default timeout of 60s to all connections.
* `PaymentMethodNonce` added to `Customer`.
* `Transaction` `Status` constants added.

BUG FIXES:
* Numerous test and CI reliability improvements to reduce flakiness of integration tests.

## 0.11.0 (May 12th, 2017)

BACKWARDS INCOMPATIBILITES:

* `TransactionGateway` `Create` now takes a `TransactionRequest`. Fields were removed from `Transaction` that are not included in the response.

IMPROVEMENTS:

* `CurrencyISOCode` added to `Transaction`.
* `GatewayRejectionReason` added to `Transaction`.
* `CVVResponseCode` added to `Transaction`.
* `AVSErrorResponseCode`, `AVSPostalCodeResponseCode` and `AVSStreetAddressResponseCode` added to `Transaction`.
* `SubscriptionId` added to `Transaction`.
* `RiskDataRequest` added and can be set on `TransactionRequest` when calling `TransactionGateway` `Create`.
* `WebhookNotificationGateway` now supports validating signatures on webhook payloads for accounts that have multiple API Keys.
* `ParseRequest` added to `WebhookNotificationGateway`.
* `WebhookTestingGateway` added for generating sample webhook notifications.

## 0.10.0 (April 14th, 2017)

BACKWARDS INCOMPATIBILITES:

* `Environment`, `MerchantId` wer
* `PublicKey`, `PrivateKey` were removed from the `Braintree` struct.
* `Subscription` `Discounts` was changed from an `[]interface{}` to a concrete type.
* `SubscriptionRequest` replaced `Subscription` as the parameter when creating subscriptions.

DEPRECATED:
* `TransactionGateway` `Settle` function was marked as deprecated and the `TestingGateway` `Settle` function should be used instead.

IMPROVEMENTS:

* `NewWithAccessToken` added to allow initialization with an access token
  instead of an API key.
* `NewWithHttpClient` was added.
* `PaymentMethodGateway` `Update`, `Find`, and `Delete` were added, and support for non-credit-card payment methods were supported.
* `PaymentMethod` was added as an interface across all payment methods.
* `Customer` `DefaultPaymentMethod` was added.
* `Transaction` `PayPalDetails` was added.
* `PayPalAccount` and `PayPalAccountGateway` were added.
* `Transaction` `DeviceData` was added.
* `Transaction` `RefundedTransactionId` and `RefundIds` were added.
* `CustomFields` were added to `Transaction` and `Customer`.
* `Transaction` and `Subscription` `Descriptor` was added to support setting dynamic descriptors.
* `Transaction` `RiskData` in responses was added.
* `Transaction` `PaymentInstrumentType` was added.
* `Subscription` modifications for `AddOn`s and `Discounts` were added as concrete objects.
* `Plan` `AddOns` and `Discounts` were added.
* `SettlementGateway` and `Settlement` were added.
* `SearchQuery` `AddTimeField` was added to support querying on time/date ranges.
* `Transaction` was added and is accessible on `BraintreeError`.
* `APIError` is now returned for 404 responses to allow developers to get the
  status code and check for a 404.
* `TestingGateway` added with functions for changing the status of transactions in sandbox to `settle`, `settlement_confirm`, `settlement_decline`, and `settlement_pending` to aid with integration tests.
* `Environment` was changed from `string` to `struct` that can be constructed with a custom base URL.

BUG FIXES:

* `Transaction` `ProcessorResponseCode` will now be zero if no processor response code is returned from the Braintree API.
* `Decimal` now marshals correctly if they are under $1.00.
* `Decimal` marshal and unmarshals correctly if it's value is zero.
* Numerous test reliability improvements to reduce flakiness of integration tests.
* Documentation improvements.

## 0.9.0 (August 5th, 2015)

BACKWARDS INCOMPATIBILITES:

* Time fields such as `CreatedAt`, `UpdatedAt`, and `DisbursementDetails`
  changed to be `time.Time` or `date.Date` types to make their use simpler
  (rather than forcing the user to deserialize)
* Additional currency fields migrated from `float64` to `Decimal` to avoid
  precision loss
* `Null(Int64|Bool)` types added to support empty XML elements. Many of
  the nonstring, `string`, fields were updated to be nullable of their actual
  type.
* `ClientToken` was renamed to `ClientTokenRequest`
* `ClientToken` API changed to allow generation of client tokens with or
  without an associated customer. `NewClientTokenRequest` was removed.

IMPROVEMENTS:

* `CustomerGateway.Search` added to permit advanced searching for customers by
  metadata
* `BraintreeError` type was exposed to make it easier to inspect whether the
  errors returned by the library are network on Braintree Gatway errors
* `ClientTokenGateway.GenerateWithCustomer` added to generate a customer
  specific client token

## 0.8.0 (April 3, 2015)

BACKWARDS INCOMPATIBILITES:

* Webhook constants made more uniform via `Webhook` suffix
* All currency amounts changed from `float` to `Decimal` to remove loss of
  precision

IMPROVEMENTS:

* Specification of a custom `http.Client` to use via `Braintree.HttpClient`.
  This enables `AppEngine` support which required a being able to use a custom
  `http.Client`.
* `DisbursementDetails` added to `Transaction`
* Support for querying disbursement webhooks added via `WebhookNotification.Disbursement`
* `TransactionGateway.Settle` added to automatically settle transactions in
  sandbox (`SubmitForSettlement` should be used in production)
* `PaymentMethodNonce` added to `CreditCard`
* `PaymentMethodNonce` added to `Transaction`
* `Decimal` arbitrary precision numeric type added to be used for currency
  amounts
* `ClientToken` support added via `ClientTokenGateway` to generate new client
  tokens

BUG FIXES:

* Typo in path for merchant account updates (`MerchantAccountGateway.Update`)
  was fixed.

## 0.7.0 (April 3, 2014)

BACKWARDS INCOMPATIBILITES:

* `InvalidResponseError` was unexported to encourage use of the new
  `BraintreeError` type
* `CreditCard.Default` changed from string to bool
* `CreditCard.Expired` changed from string to bool

IMPROVEMENTS:

* `CustomerGateway.Update` added to update metadata about the customer
* `CustomerGateway.Delete` added to allow customers to be deleted
* `Customer.DefaultCreditCard` added to return the default credit card
  associated with the customer
* `BraintreeError` type added to expose metadata about gateway errors in
  a structured manner
* `TransactionGateway.SubmitForSettlement` added to allow transactions to be
  submitted to be settled
* `TransactionGateway.Void` added to allow transactions to be voided
* Additional fields added to `Plan` (all except `Addons` and `Discounts`)
* Additional fields added to `Subscription` (all except `Addons` and `Descriptor`)
* `Subscription.Update` added to allow subscription data to be updated
* Remaining fields added to `CreditCard` and `CreditCardOptions`
* `CreditCardGateway.Update` added to update credit card information
* `CreditCardGateway.Delete` added to allow credit cards to be deleted
* `CreditCard.AllSubscriptions` added to allow subscriptions for a credit card
  to be queried
* `PlanGateway.Find` added to lookup plan by id
* `SubscriptionStatus*` constants were added to make comparisons easier
* `TransactionGateway.Search` added to permit searching for transactions by
  metadata
* `CreatedAt`, `UpdatedAt`, `PlanId` added to `Transaction`
* `ParseDate` added to facilitate parsing the date format returned by Braintree
* Adedd `AddOn` support via `AddOnGateway`
* Adedd `Discount` support via `DiscountGateway`
* Adedd `MerchantAccount` support via `MerchantAccountGateway` for submerchant
  support. Includes addition of `ServiceFeeAmount` to `Transaction`

BUG FIXES:

* `AddressGateway.Create` now copies address for sanitization to avoid
  modifying passed struct
* Errors during failed HTTP requests no longer cause a nil pointer dereference
  (when a `nil` body was `Close`d)

## 0.6.0 (June 30, 2015)

BACKWARDS INCOMPATIBILITES:

* Large scale refactoring from `0.5.0`

IMPROVEMENTS:

* Start of `Subscription` and `Plan` support
* `Address` `Create` and `Delete` support added via `AddressGateway`
* `ExpirationMonth` and `ExpirationYear` added to `CreditCard`

## 0.5.0 (May 27, 2013)

Initial release
