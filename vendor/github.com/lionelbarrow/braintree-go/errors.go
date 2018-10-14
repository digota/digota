package braintree

import "strings"

type errorGroup interface {
	For(string) errorGroup
	On(string) []FieldError
}

type BraintreeError struct {
	statusCode      int
	XMLName         string           `xml:"api-error-response"`
	Errors          responseErrors   `xml:"errors"`
	ErrorMessage    string           `xml:"message"`
	MerchantAccount *MerchantAccount `xml:",omitempty"`
	Transaction     Transaction      `xml:"transaction"`
}

func (e *BraintreeError) Error() string {
	return e.ErrorMessage
}

func (e *BraintreeError) StatusCode() int {
	return e.statusCode
}

func (e *BraintreeError) All() []FieldError {
	baseErrors := e.Errors.TransactionErrors.ErrorList.Errors
	creditCardErrors := e.Errors.TransactionErrors.CreditCardErrors.ErrorList.Errors
	customerErrors := e.Errors.TransactionErrors.CustomerErrors.ErrorList.Errors
	allErrors := append(baseErrors, creditCardErrors...)
	allErrors = append(allErrors, customerErrors...)
	return allErrors
}

func (e *BraintreeError) For(item string) errorGroup {
	switch item {
	default:
		return nil
	case "Transaction":
		return e.Errors.TransactionErrors
	}
}

func (e *BraintreeError) On(item string) []FieldError {
	return []FieldError{}
}

type responseErrors struct {
	TransactionErrors responseError `xml:"transaction"`
}

type responseError struct {
	ErrorList        errorList  `xml:"errors"`
	CreditCardErrors errorBlock `xml:"credit-card"`
	CustomerErrors   errorBlock `xml:"customer"`
}

func (r responseError) For(item string) errorGroup {
	switch item {
	default:
		return nil
	case "Base":
		return r.ErrorList.Errors
	case "Customer":
		return r.CustomerErrors.ErrorList.Errors
	case "CreditCard":
		return r.CreditCardErrors.ErrorList.Errors
	}
}

func (r responseError) On(item string) []FieldError {
	switch item {
	default:
		return []FieldError{}
	case "Base":
		return r.ErrorList.Errors
	case "Customer":
		return r.CustomerErrors.ErrorList.Errors
	case "CreditCard":
		return r.CreditCardErrors.ErrorList.Errors
	}
}

type errorBlock struct {
	ErrorList errorList `xml:"errors"`
}

type errorList struct {
	Errors FieldErrorList `xml:"error"`
}

type FieldErrorList []FieldError

func (f FieldErrorList) For(item string) errorGroup {
	return nil
}

func (f FieldErrorList) On(item string) []FieldError {
	errors := make([]FieldError, 0)
	for _, e := range f {
		if strings.ToLower(item) == e.Attribute {
			errors = append(errors, e)
		}
	}
	return errors
}

type FieldError struct {
	Code      string `xml:"code"`
	Attribute string `xml:"attribute"`
	Message   string `xml:"message"`
}
