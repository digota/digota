package braintree

import (
	"bytes"
	"encoding/xml"
	"io"
	"strconv"
	"unicode"
)

type BraintreeError struct {
	statusCode int
	errors     ValidationErrors

	ErrorMessage    string
	MerchantAccount *MerchantAccount
	Transaction     *Transaction
}

func (e *BraintreeError) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var x struct {
		Errors          ValidationErrors `xml:"errors"`
		ErrorMessage    string           `xml:"message"`
		MerchantAccount *MerchantAccount `xml:"merchant-account"`
		Transaction     *Transaction     `xml:"transaction"`
	}
	err := d.DecodeElement(&x, &start)
	if err != nil {
		return err
	}
	e.errors = x.Errors
	e.ErrorMessage = x.ErrorMessage
	e.MerchantAccount = x.MerchantAccount
	e.Transaction = x.Transaction
	return nil
}

func (e *BraintreeError) Error() string {
	return e.ErrorMessage
}

func (e *BraintreeError) StatusCode() int {
	return e.statusCode
}

func (e *BraintreeError) All() []ValidationError {
	return e.errors.AllDeep()
}

func (e *BraintreeError) For(name string) *ValidationErrors {
	return e.errors.For(name)
}

type ValidationErrors struct {
	Object           string
	ValidationErrors []ValidationError
	Children         map[string]*ValidationErrors
}

func (r *ValidationErrors) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		t, err := d.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		r.Object = errorNameKebabToCamel(start.Name.Local)

		if subStart, ok := t.(xml.StartElement); ok {
			subName := subStart.Name.Local
			if subName == "errors" {
				errorList := struct {
					ErrorList []ValidationError `xml:"error"`
				}{}
				err := d.DecodeElement(&errorList, &subStart)
				if err != nil {
					return err
				}
				r.ValidationErrors = errorList.ErrorList
			} else {
				subSectionName := errorNameKebabToCamel(subName)
				subSection := &ValidationErrors{}
				err := d.DecodeElement(subSection, &subStart)
				if err != nil {
					return err
				}
				if r.Children == nil {
					r.Children = map[string]*ValidationErrors{}
				}
				r.Children[subSectionName] = subSection
			}
		}
	}
	return nil
}

func (r *ValidationErrors) All() []ValidationError {
	if r == nil {
		return nil
	}
	return r.ValidationErrors
}

func (r *ValidationErrors) AllDeep() []ValidationError {
	if r == nil {
		return nil
	}
	errorList := append([]ValidationError{}, r.All()...)
	for _, sub := range r.Children {
		errorList = append(errorList, sub.AllDeep()...)
	}
	return errorList
}

func (r *ValidationErrors) For(name string) *ValidationErrors {
	if r == nil || r.Children == nil {
		return (*ValidationErrors)(nil)
	}
	return r.Children[name]
}

func (r *ValidationErrors) ForIndex(i int) *ValidationErrors {
	if r == nil || r.Children == nil {
		return (*ValidationErrors)(nil)
	}
	return r.Children["Index"+strconv.Itoa(i)]
}

func (r *ValidationErrors) On(name string) []ValidationError {
	if r == nil {
		return nil
	}
	errors := make([]ValidationError, 0)
	for _, e := range r.ValidationErrors {
		if name == e.Attribute {
			errors = append(errors, e)
		}
	}
	return errors
}

type ValidationError struct {
	Code      string
	Attribute string
	Message   string
}

func (e *ValidationError) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var x struct {
		Code      string `xml:"code"`
		Attribute string `xml:"attribute"`
		Message   string `xml:"message"`
	}
	err := d.DecodeElement(&x, &start)
	if err != nil {
		return err
	}
	e.Code = x.Code
	e.Attribute = errorNameSnakeToCamel(x.Attribute)
	e.Message = x.Message
	return nil
}

func errorNameSnakeToCamel(snake string) string {
	if len(snake) == 0 {
		return ""
	}
	camel := bytes.Buffer{}
	capitalizeNext := true
	for i := 0; i < len(snake); i++ {
		s := snake[i]
		if s == '_' {
			capitalizeNext = true
			continue
		} else if capitalizeNext {
			camel.WriteByte(byte(unicode.ToUpper(rune(s))))
		} else {
			camel.WriteByte(s)
		}
		capitalizeNext = false
	}
	return camel.String()
}

func errorNameKebabToCamel(kebab string) string {
	if len(kebab) == 0 {
		return ""
	}
	camel := bytes.Buffer{}
	capitalizeNext := true
	for i := 0; i < len(kebab); i++ {
		k := kebab[i]
		if k == '-' {
			capitalizeNext = true
			continue
		} else if capitalizeNext {
			camel.WriteByte(byte(unicode.ToUpper(rune(k))))
		} else {
			camel.WriteByte(k)
		}
		capitalizeNext = false
	}
	return camel.String()
}
