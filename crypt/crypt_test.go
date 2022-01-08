package crypt

import (
	"testing"
)

func TestIfGoAndReturn(t *testing.T) {
	from := "The quick brown fox jumps over the lazy dog"
	bytes := make([]byte, 48)
	copy(bytes, from)
	passwd := "1234"
	c := NewCrypterFromPassword(passwd)
	err := InplaceEncrypt(c, bytes)
	if err != nil {
		t.Error(err)
	}
	println(string(bytes))
	err = InplaceDecrypt(c, bytes)
	if err != nil {
		t.Error(err)
	}
	if from != string(bytes)[:len(from)] {
		t.Errorf("expected %s got %s", from, string(bytes))
	}
}
