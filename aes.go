package libsgo

import (
	"io"
	"bytes"
	"errors"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

/**
 * @docs https://golang.org/src/crypto/cipher/example_test.go
 * @docs https://blog.csdn.net/whatday/article/details/98292648
 * @docs https://segmentfault.com/a/1190000021267253
 * @docs http://www.361way.com/golang-rsa-aes/5828.html
 */

func NewAesLib() *AesLib {
	return &AesLib{}
}

type AesLib struct {
	secKey []byte
}

func (lib *AesLib) SetSecretKey(secretKey []byte) {
	lib.secKey = secretKey
}

func (lib *AesLib) EncryptCFB(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(lib.secKey)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

func (lib *AesLib) DecryptCFB(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(lib.secKey)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

func (lib *AesLib) EncryptCBC(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(lib.secKey)
	if err != nil {
		return nil, err
	}

	plaintext = pkcs7Padding(plaintext, aes.BlockSize)

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

func (lib *AesLib) DecryptCBC(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(lib.secKey)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext) // CryptBlocks can work in-place if the two arguments are the same.

	return pkcs7UnPadding(ciphertext), nil
}

func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

func pkcs7UnPadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}
