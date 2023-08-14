// utils/hasher_test.go

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "testPassword123"
	hash, err := HashPassword(password)

	assert.Nil(t, err, "Error should be nil when hashing a password")
	assert.NotEqual(t, password, hash, "Hash should be different from the original password")
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testPassword123"
	hash, _ := HashPassword(password)

	isCorrect := CheckPasswordHash(password, hash)
	assert.True(t, isCorrect, "Password should match the hash")

	isCorrect = CheckPasswordHash("wrongPassword", hash)
	assert.False(t, isCorrect, "Incorrect password should not match the hash")
}

