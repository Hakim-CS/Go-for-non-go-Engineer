package servitor

import (
	cryptorand "crypto/rand"
	"fmt"
	"math/rand"
	"strings"

	"golang.org/x/crypto/salsa20"
)

const (
	allowablePasswordCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~!@#$%^&*()_+-="
	minPasswordLength           = 16
	maxPasswordLength           = 64
	defaultPasswordLength       = 32
	defaultSalsa20NonceLength   = 24
	defaultSalsa20KeyLength     = 32
)

var (
	salsa20ValidNonceLengths = []int{8, 24}
)

type ServitorBeta struct {
	PasswordLength     int
	Salsa20NonceLength int
}

func NewServitorBeta(passwordLength, Salsa20NonceLength int) *ServitorBeta {
	result := ServitorBeta{
		PasswordLength: passwordLength,
	}

	if passwordLength < minPasswordLength || passwordLength > maxPasswordLength {
		result.PasswordLength = defaultPasswordLength
	}

	result.Salsa20NonceLength = defaultSalsa20NonceLength

	for _, validNonceLength := range salsa20ValidNonceLengths {
		if Salsa20NonceLength == validNonceLength {
			result.Salsa20NonceLength = Salsa20NonceLength
			break
		}
	}

	return &result
}

func (sb *ServitorBeta) DefaultPassword() (string, error) {
	var password strings.Builder

	for range sb.PasswordLength {
		randomChar := allowablePasswordCharacters[rand.Intn(len(allowablePasswordCharacters))]
		password.WriteByte(randomChar)
	}

	return password.String(), nil
}

func (sb *ServitorBeta) SymmetricEncryption(key []byte, plaintext []byte) ([]byte, error) {
	// Generate a nonce
	nonce := make([]byte, sb.Salsa20NonceLength)
	if _, err := cryptorand.Read(nonce); err != nil {
		return nil, err
	}

	// Check key length
	if len(key) != defaultSalsa20KeyLength {
		return nil, fmt.Errorf("invalid key length: %d, must be %d", len(key), defaultSalsa20KeyLength)
	}

	// Encrypt the plaintext
	var keyArray [defaultSalsa20KeyLength]byte
	copy(keyArray[:], key)
	ciphertext := make([]byte, len(plaintext))
	salsa20.XORKeyStream(ciphertext, plaintext, nonce, &keyArray)

	// Prepend the nonce to the ciphertext
	encrypted := append(nonce, ciphertext...)

	// Return the encrypted data
	return encrypted, nil
}

func (sb *ServitorBeta) SymmetricDecryption(key []byte, ciphertext []byte) ([]byte, error) {
	nonce := ciphertext[:sb.Salsa20NonceLength]
	ciphertextWithoutNonce := ciphertext[24:]

	if len(key) != defaultSalsa20KeyLength {
		return nil, fmt.Errorf("invalid key length: %d, must be %d", len(key), defaultSalsa20KeyLength)
	}

	// Decrypt the ciphertext
	var keyArray [defaultSalsa20KeyLength]byte
	copy(keyArray[:], key)
	plaintext := make([]byte, len(ciphertextWithoutNonce))
	salsa20.XORKeyStream(plaintext, ciphertextWithoutNonce, nonce, &keyArray)

	return plaintext, nil
}
