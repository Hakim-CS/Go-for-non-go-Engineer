package servitor

const (
	AlgorithmUnknown = "Unknown"
	AlgorithmAES     = "AES"
	AlgorithmSalsa20 = "Salsa20"
)

type ServitorOmega struct {
	passwordProvider PasswordProvider
	cryptoProvider   CryptoProvider
}

func NewServitorOmega(passwordProvider PasswordProvider, cryptoProvider CryptoProvider) *ServitorOmega {
	return &ServitorOmega{
		passwordProvider: passwordProvider,
		cryptoProvider:   cryptoProvider,
	}
}

func (so *ServitorOmega) GeneratePassword() (string, error) {
	return so.passwordProvider.DefaultPassword()
}

func (so *ServitorOmega) Encrypt(key, plaintext []byte) ([]byte, string, error) {
	algorithm := AlgorithmUnknown
	ciphertext, err := so.cryptoProvider.SymmetricEncryption(key, plaintext)

	if err != nil {
		return nil, algorithm, err
	}

	switch so.cryptoProvider.(type) {
	case *ServitorAlpha:
		algorithm = AlgorithmAES
	case *ServitorBeta:
		algorithm = AlgorithmSalsa20
	}

	return ciphertext, algorithm, nil
}

func (so *ServitorOmega) Decrypt(key, ciphertext []byte) ([]byte, string, error) {
	algorithm := AlgorithmUnknown
	plaintext, err := so.cryptoProvider.SymmetricDecryption(key, ciphertext)

	if err != nil {
		return nil, algorithm, err
	}

	switch so.cryptoProvider.(type) {
	case *ServitorAlpha:
		algorithm = AlgorithmAES
	case *ServitorBeta:
		algorithm = AlgorithmSalsa20
	}

	return plaintext, algorithm, nil
}
