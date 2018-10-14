package testhelpers

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"strconv"
)

func RandomInt() int64 {
	var n int64
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(b)
	err = binary.Read(buf, binary.LittleEndian, &n)
	if err != nil {
		panic(err)
	}
	return n
}

func RandomString() string {
	return strconv.FormatInt(RandomInt(), 10)
}
