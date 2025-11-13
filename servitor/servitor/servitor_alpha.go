package servitor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"

	"github.com/google/uuid"
)

var (
	aesValidKeyLengths = []int{16, 24, 32}
)

type ServitorAlpha struct {
}

func NewServitorAlpha() *ServitorAlpha {
	return &ServitorAlpha{}
}

func (sa *ServitorAlpha) DefaultPassword() (string, error) {
	return uuid.NewString(), nil
}

func (sa *ServitorAlpha) SymmetricEncryption(key []byte, plaintext []byte) ([]byte, error) {
	// Generate a random initialization vector
	iv := make([]byte, aes.BlockSize)

	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	// check key length
	validKeyLength := false

	for _, length := range aesValidKeyLengths {
		if len(key) == length {
			validKeyLength = true
			break
		}
	}

	if !validKeyLength {
		return nil, fmt.Errorf("invalid key length: %d, must be 16, 24, or 32", len(key))
	}

	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	paddedPlaintext, err := addPKCS7Padding(plaintext, aes.BlockSize)

	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, len(paddedPlaintext))
	mode.CryptBlocks(ciphertext, paddedPlaintext)

	encrypted := append(iv, ciphertext...)

	return encrypted, nil
}

func (sa *ServitorAlpha) SymmetricDecryption(key []byte, ciphertext []byte) ([]byte, error) {
	// extract the IV
	iv := ciphertext[:aes.BlockSize]
	ciphertextWithoutIV := ciphertext[aes.BlockSize:]

	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	plaintext := make([]byte, len(ciphertextWithoutIV))
	mode.CryptBlocks(plaintext, ciphertextWithoutIV)

	plaintextWithoutPadding, err := removePKCS7Padding(plaintext, aes.BlockSize)

	if err != nil {
		return nil, err
	}

	return plaintextWithoutPadding, nil
}
