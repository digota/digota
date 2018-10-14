package braintree

import "strconv"

type ProcessorResponseCode int

func (rc ProcessorResponseCode) Int() int {
	return int(rc)
}

// UnmarshalText fills the response code with the integer value if the text contains one in string form. If the text is zero length, the response code's value is unchanged but unmarshaling is successful.
func (rc *ProcessorResponseCode) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return nil
	}

	n, err := strconv.Atoi(string(text))
	if err != nil {
		return err
	}

	*rc = ProcessorResponseCode(n)

	return nil
}

// MarshalText returns a string in bytes of the number, or nil in the case it is zero.
func (rc ProcessorResponseCode) MarshalText() ([]byte, error) {
	if rc == 0 {
		return nil, nil
	}
	return []byte(strconv.Itoa(int(rc))), nil
}
