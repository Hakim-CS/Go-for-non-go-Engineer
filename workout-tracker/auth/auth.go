package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWTSecret is the secret key used to sign JWT tokens
// In production, this should come from environment variables
var JWTSecret = []byte("your-secret-key-change-this-in-production")

// TokenExpirationHours defines how long a token is valid
var TokenExpirationHours = 24

// Claims represents the data stored in a JWT token
// This is what gets encoded into the token
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// HashPassword takes a plain text password and returns a bcrypt hash
// This is used when creating a new user or changing password
func HashPassword(password string) (string, error) {
	// bcrypt automatically adds salt and uses a secure hashing algorithm
	// The cost factor (10) determines how slow the hash is (higher = more secure but slower)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// CheckPassword compares a plain text password with a hashed password
// Returns true if they match, false otherwise
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken creates a new JWT token for a user
// The token contains the user's ID and username, and is valid for 24 hours
func GenerateToken(userID int, username string) (string, error) {
	// Set when the token expires
	expirationTime := time.Now().Add(time.Duration(TokenExpirationHours) * time.Hour)

	// Create the claims (the data that goes in the token)
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create the token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret key
	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken checks if a token is valid and returns the claims
// If the token is invalid or expired, it returns an error
func ValidateToken(tokenString string) (*Claims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Make sure the token uses the correct signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}
		return JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Extract and return the claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
