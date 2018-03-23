package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
)

func NewAESWriter(key []byte, w io.Writer) *cipher.StreamWriter {
	kl := len(key)
	if kl < 16 {
		panic("res key 长度不能小于16")
	}
	if kl >= 32 {
		key = key[:32]
	} else if kl >= 24 {
		key = key[:24]
	} else {
		key = key[:16]
	}
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return &cipher.StreamWriter{
		S: cipher.NewOFB(aesCipher, key[:aes.BlockSize]),
		W: w,
	}
}

func NewAESReader(key []byte, r io.Reader) *cipher.StreamReader {
	kl := len(key)
	if kl < 16 {
		panic("res key 长度不能小于16")
	}
	if kl >= 32 {
		key = key[:32]
	} else if kl >= 24 {
		key = key[:24]
	} else {
		key = key[:16]
	}
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return &cipher.StreamReader{
		S: cipher.NewOFB(aesCipher, key[:aes.BlockSize]),
		R: r,
	}
}
