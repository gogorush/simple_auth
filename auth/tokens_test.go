// auth/tokens_test.go

package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	username := "testuser"
	tokenDetails, err := GenerateToken(username)

	assert.Nil(t, err, "Error should be nil")
	assert.NotEmpty(t, tokenDetails.Token, "Token should not be empty")
	assert.True(t, tokenDetails.ExpiresAt > time.Now().Unix(), "Token expiration should be in the future")

	// Check if token exists in the in-memory store
	_, exists := Tokens.Get(tokenDetails.Token)
	assert.True(t, exists, "Token should exist in the store")
}

func TestValidateToken(t *testing.T) {
	username := "testuser2"
	tokenDetails, _ := GenerateToken(username)

	retrievedUsername, err := ValidateToken(tokenDetails.Token)

	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, username, retrievedUsername, "Usernames should match")

	// Check for invalid token
	_, err = ValidateToken("invalidtoken")
	assert.NotNil(t, err, "Error should not be nil for an invalid token")
}

func TestInvalidateToken(t *testing.T) {
	username := "testuser3"
	tokenDetails, _ := GenerateToken(username)

	// Invalidate the token
	InvalidateToken(tokenDetails.Token)

	_, exists := Tokens.Get(tokenDetails.Token)
	assert.False(t, exists, "Token should not exist in the store after invalidation")

	a, err := ValidateToken(tokenDetails.Token)
	assert.Equal(t, a, username, "should be user name")
	assert.Nil(t, err, "Error should not be nil for a deleted token")
}

func TestTokenExpiry(t *testing.T) {
	// Set token duration to 2 seconds for this test
	setTokenDuration(2 * time.Second)

	username := "testuserExpiry"
	tokenDetails, _ := GenerateToken(username)

	// Wait for 3 seconds to ensure the token expires
	time.Sleep(3 * time.Second)

	_, err := ValidateToken(tokenDetails.Token)

	assert.NotNil(t, err, "Error should not be nil for an expired token")
	assert.Equal(t, "token has invalid claims: token is expired", err.Error(), "Expected token expired error")
}
