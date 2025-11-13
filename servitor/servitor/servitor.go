package servitor

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

type PasswordProvider interface {
	DefaultPassword() (string, error)
}

type CryptoProvider interface {
	SymmetricEncryption(key []byte, plaintext []byte) ([]byte, error)
	SymmetricDecryption(key []byte, ciphertext []byte) ([]byte, error)
}

func addPKCS7Padding(text []byte, blockSize int) ([]byte, error) {
	if blockSize < 1 || blockSize > 255 {
		return nil, fmt.Errorf("invalid block size: %d, must be 1-255", blockSize)
	}

	paddingLength := blockSize - (len(text) % blockSize)
	padding := bytes.Repeat([]byte{
		byte(paddingLength),
	}, paddingLength)

	return append(text, padding...), nil
}

func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func HexEncode(data []byte) string {
	return hex.EncodeToString(data)
}

func removePKCS7Padding(text []byte, blockSize int) ([]byte, error) {
	if blockSize < 1 || blockSize > 255 {
		return nil, fmt.Errorf("invalid block size: %d, must be 1-255", blockSize)
	}

	if len(text) == 0 {
		return nil, fmt.Errorf("invalid plaintext")
	}

	if (len(text) % blockSize) != 0 {
		return nil, fmt.Errorf("invalid plaintext length")
	}

	paddingLength := int(text[len(text)-1])
	paddingBytes := bytes.Repeat(
		[]byte{byte(paddingLength)}, paddingLength)

	if paddingLength == 0 || paddingLength > blockSize || !bytes.HasSuffix(text, paddingBytes) {
		return nil, fmt.Errorf("invalid padding")
	}

	return text[:len(text)-paddingLength], nil
}

func Base64Decode(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}

func HexDecode(data string) ([]byte, error) {
	return hex.DecodeString(data)
}
