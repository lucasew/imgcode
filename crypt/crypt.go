package crypt

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

// Crypter holds things that will be used to {de,}crypt things
type Crypter struct {
	cipherkey []byte
}

const contextCrypter = "crypter"

func ContextWithCrypter(from context.Context, crypter Crypter) context.Context {
	return context.WithValue(from, contextCrypter, crypter)
}

func GetCrypterFromContext(ctx context.Context) *Crypter {
	switch ret := ctx.Value(contextCrypter).(type) {
	case Crypter:
		return &ret
	default:
		return nil
	}
}

// GetCipher returns the cipher generated from crypter
func (c Crypter) GetCipher() (cipher.Block, error) {
	ciph, err := aes.NewCipher(c.cipherkey)
	if err != nil {
		return nil, err
	}
	return ciph, nil

}

// GetOFB returns the stream generated from the cipher that is generated from the crypter
func (c Crypter) GetOFB() (cipher.Stream, error) {
	ciph, err := c.GetCipher()
	if err != nil {
		return nil, err
	}
	iv := make([]byte, aes.BlockSize)
	_, err = rand.Reader.Read(iv)
	if err != nil {
		return nil, err
	}
	return cipher.NewOFB(ciph, iv), nil
}

// NewCrypterFromPassword Creates a crypter from a password
func NewCrypterFromPassword(passwd string) Crypter {
	cipherkey := make([]byte, aes.BlockSize)
	sha := sha256.New()
	sum := sha.Sum([]byte(passwd))
	for i := 0; i < aes.BlockSize; i++ {
		cipherkey[i] = sum[i]
	}
	fmt.Printf("%X\n", cipherkey)
	return Crypter{cipherkey}
}

func InplaceEncrypt(c Crypter, buf []byte) error {
	ciph, err := c.GetCipher()
	if err != nil {
		return err
	}
	var i int
	for i = 0; i < (len(buf) / aes.BlockSize); i++ {
		ciph.Encrypt(buf[aes.BlockSize*i:], buf[aes.BlockSize*i:])
	}
	return nil
}

func InplaceDecrypt(c Crypter, buf []byte) error {
	ciph, err := c.GetCipher()
	if err != nil {
		return err
	}
	var i int
	for i = 0; i < (len(buf) / aes.BlockSize); i++ {
		ciph.Decrypt(buf[aes.BlockSize*i:], buf[aes.BlockSize*i:])
	}
	return nil
}
