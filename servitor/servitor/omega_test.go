package servitor

import (
	"testing"
)

func TestNewServitorOmega(t *testing.T) {
	testCases := []struct {
		name             string
		passwordProvider PasswordProvider
		cryptoProvider   CryptoProvider
	}{
		{
			name:             "IMPLEMENTATION_SERVITOR_ALPHA_ALPHA",
			passwordProvider: NewServitorAlpha(),
			cryptoProvider:   NewServitorAlpha(),
		},
		{
			name:             "IMPLEMENTATION_SERVITOR_BETA_BETA",
			passwordProvider: NewServitorBeta(defaultPasswordLength, defaultSalsa20NonceLength),
			cryptoProvider:   NewServitorBeta(defaultPasswordLength, defaultSalsa20NonceLength),
		},
		{
			name:             "IMPLEMENTATION_SERVITOR_ALPHA_BETA",
			passwordProvider: NewServitorAlpha(),
			cryptoProvider:   NewServitorBeta(defaultPasswordLength, defaultSalsa20NonceLength),
		},
		{
			name:             "IMPLEMENTATION_SERVITOR_BETA_ALPHA",
			passwordProvider: NewServitorBeta(defaultPasswordLength, defaultSalsa20NonceLength),
			cryptoProvider:   NewServitorAlpha(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			so := NewServitorOmega(tc.passwordProvider, tc.cryptoProvider)

			if so == nil {
				t.Errorf("NewServitorOmega() returned nil")
			}

			if so.passwordProvider != tc.passwordProvider {
				t.Errorf("NewServitorOmega() returned passwordProvider %v, want %v",
					so.passwordProvider, tc.passwordProvider)
			}

			if so.cryptoProvider != tc.cryptoProvider {
				t.Errorf("NewServitorOmega() returned cryptoProvider %v, want %v",
					so.cryptoProvider, tc.cryptoProvider)
			}
		})
	}
}

func TestGeneratePassword(t *testing.T) {
	type wanted struct {
		passwordLength int
		wantErr        bool
	}

	testCases := []struct {
		name             string
		passwordProvider PasswordProvider
		cryptoProvider   CryptoProvider
		want             wanted
	}{
		{
			name:             "IMPLEMENTATION_SERVITOR_ALPHA_ALPHA",
			passwordProvider: NewServitorAlpha(),
			cryptoProvider:   NewServitorAlpha(),
			want: wanted{
				passwordLength: 36,
				wantErr:        false,
			},
		},
		{
			name:             "IMPLEMENTATION_SERVITOR_ALPHA_BETA",
			passwordProvider: NewServitorAlpha(),
			cryptoProvider:   NewServitorBeta(defaultPasswordLength, defaultSalsa20NonceLength),
			want: wanted{
				passwordLength: 36,
				wantErr:        false,
			},
		},
		{
			name:             "IMPLEMENTATION_SERVITOR_BETA_ALPHA",
			passwordProvider: NewServitorBeta(defaultPasswordLength, defaultSalsa20NonceLength),
			cryptoProvider:   NewServitorAlpha(),
			want: wanted{
				passwordLength: defaultPasswordLength,
				wantErr:        false,
			},
		},
		{
			name:             "IMPLEMENTATION_SERVITOR_BETA_BETA",
			passwordProvider: NewServitorBeta(defaultPasswordLength, defaultSalsa20NonceLength),
			cryptoProvider:   NewServitorBeta(defaultPasswordLength, defaultSalsa20NonceLength),
			want: wanted{
				passwordLength: defaultPasswordLength,
				wantErr:        false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			so := NewServitorOmega(tc.passwordProvider, tc.cryptoProvider)
			password, err := so.GeneratePassword()

			if err != nil && !tc.want.wantErr {
				t.Errorf("GeneratePassword() returned error: %v", err)
			}

			if err == nil && tc.want.wantErr {
				t.Errorf("GeneratePassword() expected error, got nil")
			}

			if len(password) != tc.want.passwordLength {
				t.Errorf("GeneratePassword() returned password length %d, want %d",
					len(password), tc.want.passwordLength)
			}
		})
	}
}

func TestEncrypt(t *testing.T) {
	type args struct {
		key       []byte
		plaintext []byte
	}

	type wanted struct {
		algorithm string
		wantErr   bool
	}

	testCases := []struct {
		name             string
		passwordProvider PasswordProvider
		cryptoProvider   CryptoProvider
		args             args
		want             wanted
	}{
		{
			name:             "ENCRYPT_SERVITOR_ALPHA",
			passwordProvider: NewServitorAlpha(),
			cryptoProvider:   NewServitorAlpha(),
			args: args{
				key:       []byte("testkey123456789"),
				plaintext: []byte("plaintext"),
			},
			want: wanted{
				algorithm: AlgorithmAES,
				wantErr:   false,
			},
		},
		{
			name:             "ENCRYPT_SERVITOR_BETA",
			passwordProvider: NewServitorBeta(defaultPasswordLength, defaultSalsa20NonceLength),
			cryptoProvider:   NewServitorBeta(defaultPasswordLength, defaultSalsa20NonceLength),
			args: args{
				key:       []byte("testkey123456789testkey123456789"),
				plaintext: []byte("plaintext"),
			},
			want: wanted{
				algorithm: AlgorithmSalsa20,
				wantErr:   false,
			},
		},
		{
			name:             "ENCRYPT_SERVITOR_ALPHA_INVALID_KEY_LENGTH",
			passwordProvider: NewServitorAlpha(),
			cryptoProvider:   NewServitorAlpha(),
			args: args{
				key:       []byte("TheDummyKeyWithInvalidLength"),
				plaintext: []byte("plaintext"),
			},
			want: wanted{
				algorithm: AlgorithmUnknown,
				wantErr:   true,
			},
		},
		{
			name:             "ENCRYPT_SERVITOR_BETA_INVALID_KEY_LENGTH",
			passwordProvider: NewServitorBeta(defaultPasswordLength, defaultSalsa20NonceLength),
			cryptoProvider:   NewServitorBeta(defaultPasswordLength, defaultSalsa20NonceLength),
			args: args{
				key:       []byte("TheDummyKeyWithInvalidLength"),
				plaintext: []byte("plaintext"),
			},
			want: wanted{
				algorithm: AlgorithmUnknown,
				wantErr:   true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			so := NewServitorOmega(tc.passwordProvider, tc.cryptoProvider)
			ciphertext, algorithm, err := so.Encrypt(tc.args.key, tc.args.plaintext)

			if err != nil && !tc.want.wantErr {
				t.Errorf("Encrypt() returned error: %v", err)
			}

			if err == nil && tc.want.wantErr {
				t.Errorf("Encrypt() expected error, got nil")
			}

			if algorithm != tc.want.algorithm {
				t.Errorf("Encrypt() returned algorithm %s, want %s", algorithm, tc.want.algorithm)
			}

			if ciphertext == nil && !tc.want.wantErr {
				t.Errorf("Encrypt() returned nil ciphertext, want non-nil")
			}
		})
	}
}

func TestDecrypt(t *testing.T) {
	type args struct {
		key []byte
	}
	type wanted struct {
		algorithm string
		wantErr   bool
	}

	testCases := []struct {
		name             string
		passwordProvider PasswordProvider
		cryptoProvider   CryptoProvider
		args             args
		want             wanted
	}{
		{
			name:             "DECRYPT_SERVITOR_ALPHA",
			passwordProvider: NewServitorAlpha(),
			cryptoProvider:   NewServitorAlpha(),
			args: args{
				key: []byte("testkey123456789"),
			},
			want: wanted{
				algorithm: AlgorithmAES,
				wantErr:   false,
			},
		},
		{
			name:             "DECRYPT_SERVITOR_BETA",
			passwordProvider: NewServitorBeta(defaultPasswordLength, defaultSalsa20NonceLength),
			cryptoProvider:   NewServitorBeta(defaultPasswordLength, defaultSalsa20NonceLength),
			args: args{
				key: []byte("testkey123456789testkey123456789"),
			},
			want: wanted{
				algorithm: AlgorithmSalsa20,
				wantErr:   false,
			},
		},
		{
			name:             "DECRYPT_SERVITOR_ALPHA_INVALID_KEY_LENGTH",
			passwordProvider: NewServitorAlpha(),
			cryptoProvider:   NewServitorAlpha(),
			args: args{
				key: []byte("TheDummyKeyWithInvalidLength"),
			},
			want: wanted{
				algorithm: AlgorithmUnknown,
				wantErr:   true,
			},
		},
		{
			name:             "DECRYPT_SERVITOR_BETA_INVALID_KEY_LENGTH",
			passwordProvider: NewServitorBeta(defaultPasswordLength, defaultSalsa20NonceLength),
			cryptoProvider:   NewServitorBeta(defaultPasswordLength, defaultSalsa20NonceLength),
			args: args{
				key: []byte("TheDummyKeyWithInvalidLength"),
			},
			want: wanted{
				algorithm: AlgorithmUnknown,
				wantErr:   true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			so := NewServitorOmega(tc.passwordProvider, tc.cryptoProvider)
			plaintext := []byte("this is a secret message")

			// Encrypt the plaintext first
			ciphertext, _, err := so.Encrypt(tc.args.key, plaintext)

			if err != nil {
				if !tc.want.wantErr {
					t.Errorf("Encrypt() got unexpected error = %v", err)
				}

				return
			}

			// Decrypt the ciphertext
			decryptedText, algorithm, err := so.Decrypt(tc.args.key, ciphertext)
			if err != nil && !tc.want.wantErr {
				t.Errorf("Decrypt() got unexpected error = %v", err)
			}
			if err == nil && tc.want.wantErr {
				t.Errorf("Decrypt() got no error, want error")
			}
			if algorithm != tc.want.algorithm {
				t.Errorf("Decrypt() algorithm = %v, want %v", algorithm, tc.want.algorithm)
			}
			if string(decryptedText) != string(plaintext) && !tc.want.wantErr {
				t.Errorf("Decrypt() got = %v, want %v", string(decryptedText), string(plaintext))
			}
		})
	}
}
