package main

import (
	"fmt"

	"github.com/timpamungkas/servitor/servitor"
)

func main() {
	servitorAlpha := servitor.NewServitorAlpha()
	servitorBeta := servitor.NewServitorBeta(32, 34)
	servitorOmega := servitor.NewServitorOmega(servitorBeta, servitorAlpha)

	// password
	password, err := servitorOmega.GeneratePassword()

	if err != nil {
		fmt.Println("Error generating password:", err)
		return
	}

	fmt.Println("Generated password:", password)

	// crypto
	key := []byte("f6SrJBymPB9eDyy1NmBu1RfnM5x1YTcF")
	plaintext := []byte("Hello, Servitor Omega!")

	encrypted, algorithm, err := servitorOmega.Encrypt(key, plaintext)

	if err != nil {
		fmt.Println("Error encrypting:", err)
		return
	}

	base64Encoded := servitor.Base64Encode(encrypted)
	hexEncoded := servitor.HexEncode(encrypted)

	fmt.Printf("[Encryption] Algorithm: %s, encrypted base64: %s\n", algorithm, base64Encoded)
	fmt.Printf("[Encryption] Algorithm: %s, encrypted hex: %s\n", algorithm, hexEncoded)

	// decryption
	ciphertext, err := servitor.Base64Decode(base64Encoded)

	if err != nil {
		fmt.Println("Error decoding base64:", err)
		return
	}

	decrypted, algorithm, err := servitorOmega.Decrypt(key, ciphertext)

	if err != nil {
		fmt.Println("Error decrypting:", err)
		return
	}

	fmt.Printf("[Decryption] Algorithm: %s, decrypted: %s\n", algorithm, decrypted)
}
