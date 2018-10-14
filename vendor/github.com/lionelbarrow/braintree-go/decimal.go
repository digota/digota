package braintree

import (
	"bytes"
	"strconv"
	"strings"
)

const precision = 16

// Decimal represents fixed precision numbers
type Decimal struct {
	Unscaled int64
	Scale    int
}

// NewDecimal creates a new decimal number equal to
// unscaled ** 10 ^ (-scale)
func NewDecimal(unscaled int64, scale int) *Decimal {
	return &Decimal{Unscaled: unscaled, Scale: scale}
}

// MarshalText outputs a decimal representation of the scaled number
func (d *Decimal) MarshalText() (text []byte, err error) {
	b := new(bytes.Buffer)
	if d.Scale <= 0 {
		b.WriteString(strconv.FormatInt(d.Unscaled, 10))
		b.WriteString(strings.Repeat("0", -d.Scale))
	} else {
		str := strconv.FormatInt(d.Unscaled, 10)
		if len(str) < d.Scale {
			str = strings.Repeat("0", d.Scale) + str
		}
		b.WriteString(str[:len(str)-d.Scale])
		b.WriteString(".")
		b.WriteString(str[len(str)-d.Scale:])
	}
	return b.Bytes(), nil
}

// UnmarshalText creates a Decimal from a string representation (e.g. 5.20)
// Currently only supports decimal strings
func (d *Decimal) UnmarshalText(text []byte) (err error) {
	var (
		str            = string(text)
		unscaled int64 = 0
		scale    int   = 0
	)

	if str == "" {
		return nil
	}

	if i := strings.Index(str, "."); i != -1 {
		scale = len(str) - i - 1
		str = strings.Replace(str, ".", "", 1)
	}

	if unscaled, err = strconv.ParseInt(str, 10, 64); err != nil {
		return err
	}

	d.Unscaled = unscaled
	d.Scale = scale

	return nil
}

// Cmp compares x and y and returns:
//
//   -1 if x <  y
//    0 if x == y
//   +1 if x >  y
//
func (x *Decimal) Cmp(y *Decimal) int {
	xUnscaled, yUnscaled := x.Unscaled, y.Unscaled
	xScale, yScale := x.Scale, y.Scale

	for ; xScale > yScale; xScale-- {
		yUnscaled = yUnscaled * 10
	}

	for ; yScale > xScale; yScale-- {
		xUnscaled = xUnscaled * 10
	}

	switch {
	case xUnscaled < yUnscaled:
		return -1
	case xUnscaled > yUnscaled:
		return 1
	default:
		return 0
	}
}

// String returns string representation of Decimal
func (d *Decimal) String() string {
	b, err := d.MarshalText()

	if err != nil {
		panic(err) //should never happen (see: MarshalText)
	}

	return string(b)
}
