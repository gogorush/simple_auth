package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWTGenerateToken(t *testing.T) {
	service := NewJWTTokenService(JwtKey, TokenDuration)
	username := "testuser"
	tokenDetails, err := service.GenerateToken(username)

	assert.Nil(t, err, "Error should be nil")
	assert.NotEmpty(t, tokenDetails.Token, "Token should not be empty")
	assert.True(t, tokenDetails.ExpiresAt > time.Now().Unix(), "Token expiration should be in the future")

	// Check if token exists in the in-memory store
	_, exists := Tokens.Get(tokenDetails.Token)
	assert.True(t, exists, "Token should exist in the store")
}

func TestJWTValidateToken(t *testing.T) {
	service := NewJWTTokenService(JwtKey, TokenDuration)
	username := "testuser2"
	tokenDetails, _ := service.GenerateToken(username)

	retrievedUsername, err := service.ValidateToken(tokenDetails.Token)

	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, username, retrievedUsername, "Usernames should match")

	// Check for invalid token
	_, err = service.ValidateToken("invalidtoken")
	assert.NotNil(t, err, "Error should not be nil for an invalid token")
}

func TestJWTInvalidateToken(t *testing.T) {
	service := NewJWTTokenService(JwtKey, TokenDuration)
	username := "testuser3"
	tokenDetails, _ := service.GenerateToken(username)

	// Invalidate the token
	service.InvalidateToken(tokenDetails.Token)

	_, exists := Tokens.Get(tokenDetails.Token)
	assert.False(t, exists, "Token should not exist in the store after invalidation")

	a, err := service.ValidateToken(tokenDetails.Token)
	assert.Equal(t, a, "", "should be empty")
	assert.NotNil(t, err, "Error should not be nil for a deleted token")
}

func TestInMemoryTokenService(t *testing.T) {
	service := NewInMemoryTokenService(2 * time.Second)

	username := "testMemoryUser"
	tokenDetails, err := service.GenerateToken(username)
	assert.Nil(t, err, "Error should be nil")
	assert.NotEmpty(t, tokenDetails.Token, "Token should not be empty")

	// Validate token
	retrievedUsername, err := service.ValidateToken(tokenDetails.Token)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, username, retrievedUsername, "Usernames should match")

	// Wait for token to expire
	time.Sleep(3 * time.Second)

	// Validate expired token
	_, err = service.ValidateToken(tokenDetails.Token)
	assert.NotNil(t, err, "Expected error for expired token")

	// Invalidate the token
	service.InvalidateToken(tokenDetails.Token)
	_, err = service.ValidateToken(tokenDetails.Token)
	assert.NotNil(t, err, "Expected error for invalidated token")
}

