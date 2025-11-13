package servitor

import (
	"math/rand"
	"regexp"
	"strconv"
	"testing"
)

func TestNewServitorBeta(t *testing.T) {
	type wanted struct {
		passwordLength     int
		salsa20NonceLength int
	}

	testCases := []struct {
		name               string
		passwordLength     int
		salsa20NonceLength int
		want               wanted
	}{
		{
			name:               "VALID_LENGTHS",
			passwordLength:     32,
			salsa20NonceLength: 24,
			want: wanted{
				passwordLength:     32,
				salsa20NonceLength: 24,
			},
		},
		{
			name:               "INVALID_PASSWORD_LENGTH_LESS_THAN_16",
			passwordLength:     15,
			salsa20NonceLength: 24,
			want: wanted{
				passwordLength:     defaultPasswordLength,
				salsa20NonceLength: 24,
			},
		},
		{
			name:               "INVALID_PASSWORD_LENGTH_MORE_THAN_64",
			passwordLength:     65,
			salsa20NonceLength: 24,
			want: wanted{
				passwordLength:     defaultPasswordLength,
				salsa20NonceLength: 24,
			},
		},
		{
			name:               "INVALID_SALSA20_NONCE_LENGTH",
			passwordLength:     32,
			salsa20NonceLength: 17,
			want: wanted{
				passwordLength:     32,
				salsa20NonceLength: defaultSalsa20NonceLength,
			},
		},
		{
			name:               "BOTH_INVALID_LENGTHS",
			passwordLength:     10,
			salsa20NonceLength: 19,
			want: wanted{
				passwordLength:     defaultPasswordLength,
				salsa20NonceLength: defaultSalsa20NonceLength,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sb := NewServitorBeta(tc.passwordLength, tc.salsa20NonceLength)

			if sb.PasswordLength != tc.want.passwordLength {
				t.Errorf("expected password length %d, got %d", tc.want.passwordLength, sb.PasswordLength)
			}

			if sb.Salsa20NonceLength != tc.want.salsa20NonceLength {
				t.Errorf("expected salsa20 nonce length %d, got %d", tc.want.salsa20NonceLength, sb.Salsa20NonceLength)
			}
		})
	}
}

func TestBeta_DefaultPassword(t *testing.T) {
	testCases := []struct {
		name           string
		passwordLength int
		wantRegex      string
		wantErr        bool
	}{
		{
			name:           "PASSWORD_LENGTH_DEFAULT",
			passwordLength: defaultPasswordLength,
			wantRegex:      "^[a-zA-Z0-9~!@#$%^&*()_+-=]{" + strconv.Itoa(defaultPasswordLength) + "}$",
			wantErr:        false,
		},
		{
			name:           "PASSWORD_LENGTH_16",
			passwordLength: 16,
			wantRegex:      "^[a-zA-Z0-9~!@#$%^&*()_+-=]{16}$",
			wantErr:        false,
		},
		{
			name:           "PASSWORD_LENGTH_32",
			passwordLength: 32,
			wantRegex:      `^[a-zA-Z0-9~!@#$%^&*()_+-=]{32}$`,
			wantErr:        false,
		},
		{
			name:           "PASSWORD_LENGTH_64",
			passwordLength: 64,
			wantRegex:      `^[a-zA-Z0-9~!@#$%^&*()_+-=]{64}$`,
			wantErr:        false,
		},
		{
			name:           "PASSWORD_LENGTH_LESS_THAN_16",
			passwordLength: 16 - 1 - rand.Intn(10),
			wantRegex:      "^[a-zA-Z0-9~!@#$%^&*()_+-=]{" + strconv.Itoa(defaultPasswordLength) + "}$",
			wantErr:        false,
		},
		{
			name:           "PASSWORD_LENGTH_MORE_THAN_64",
			passwordLength: 64 + 1 + rand.Intn(10),
			wantRegex:      "^[a-zA-Z0-9~!@#$%^&*()_+-=]{" + strconv.Itoa(defaultPasswordLength) + "}$",
			wantErr:        false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sa := NewServitorBeta(tc.passwordLength, defaultSalsa20NonceLength)
			defaultPassword, err := sa.DefaultPassword()

			if err != nil && !tc.wantErr {
				t.Errorf("DefaultPassword() throws unexpected error %v", err)
			}

			if err == nil && tc.wantErr {
				t.Errorf("DefaultPassword() got: %v, want: error", defaultPassword)
			}

			passwordRegex := regexp.MustCompile(tc.wantRegex)

			if !passwordRegex.MatchString(defaultPassword) {
				t.Errorf("DefaultPassword() got: %v, want: valid password vs regex %v", defaultPassword, tc.wantRegex)
			}
		})
	}

}

func TestBeta_SymmetricEncryption(t *testing.T) {
	sb := NewServitorBeta(defaultPasswordLength, defaultSalsa20NonceLength)

	testCases := []struct {
		name      string
		key       []byte
		plaintext []byte
		wantErr   bool
	}{
		{
			name:      "VALID_KEY_32",
			key:       []byte("thisIs32BitKey121234567812345678"),
			plaintext: []byte("this is a secret message"),
			wantErr:   false,
		},
		{
			name:      "INVALID_KEY_LENGTH",
			key:       []byte("not32BitKey"),
			plaintext: []byte("this is a secret message"),
			wantErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := sb.SymmetricEncryption(tc.key, tc.plaintext)

			if err != nil && !tc.wantErr {
				t.Errorf("SymmetricEncryption() throws unexpected error %v", err)
			}

			if err == nil && tc.wantErr {
				t.Errorf("SymmetricEncryption() got: %v, want: error", err)
			}
		})
	}
}

func TestBeta_SymmetricDecryption(t *testing.T) {
	sb := NewServitorBeta(defaultPasswordLength, defaultSalsa20NonceLength)

	testCases := []struct {
		name      string
		key       []byte
		plaintext []byte
		wantErr   bool
	}{
		{
			name:      "VALID_KEY_32",
			key:       []byte("thisIs32BitKey121234567812345678"),
			plaintext: []byte("this is a secret message"),
			wantErr:   false,
		},
		{
			name:      "INVALID_KEY_LENGTH",
			key:       []byte("not32BitKey"),
			plaintext: []byte("this is a secret message"),
			wantErr:   true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt the plaintext first
			ciphertext, err := sb.SymmetricEncryption(tt.key, tt.plaintext)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("SymmetricEncryption() got unexpected error %v", err)
				}
				return
			}

			// Decrypt the ciphertext
			decryptedText, err := sb.SymmetricDecryption(tt.key, ciphertext)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("SymmetricDecryption() got unexpected error %v", err)
				}
				return
			}

			if tt.wantErr {
				t.Errorf("SymmetricDecryption() got no error, want error")
				return
			}

			if string(decryptedText) != string(tt.plaintext) {
				t.Errorf("SymmetricDecryption() got: %v, want: %v", string(decryptedText), string(tt.plaintext))
			}
		})
	}
}
