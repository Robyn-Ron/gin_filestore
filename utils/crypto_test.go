package utils

import (
	"crypto/sha1"
	"fmt"
	"io"
	"testing"
)

func TestCryptoHash(t *testing.T) {
	h := sha1.New()
	io.WriteString(h, "hello my boy.")
	io.WriteString(h, "Here is your father.")
	fmt.Println(len(h.Sum(nil)))
	io.WriteString(h, "Here is your father.")
	fmt.Println(len(h.Sum([]byte{12,42,42})))
	fmt.Println(string(h.Sum([]byte{12,42,42})))
	fmt.Println(string(h.Sum(nil)))
	fmt.Println([]byte("")==nil)
}

