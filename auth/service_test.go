// auth/service_test.go

package auth

import (
	"testing"

	"github.com/gogorush/simple_auth/utils"
	"github.com/stretchr/testify/assert"
)

var authService AuthService = &InMemoryAuthService{}

func setup() {
	// Reset the global maps to ensure a fresh state for each test
	Users = utils.NewConcurrentMap()
	Roles = utils.NewConcurrentMap()
	Tokens = utils.NewConcurrentMap()

	// Initialize the AuthService
	authService = &InMemoryAuthService{
		TokenSvc: NewInMemoryTokenService(TokenDuration),
	}

	// Initialize admin user
	InitAdmin(User{
		Username: "admin",
		Password: "pass",
	})
}

func TestCreateUser(t *testing.T) {
	setup()

	// Create a new user
	err := authService.CreateUser("testuser1", "password123")
	assert.Nil(t, err, "Error should be nil while creating a new user")

	// Attempt to create the same user again
	err = authService.CreateUser("testuser1", "password123")
	assert.NotNil(t, err, "Error should not be nil when trying to create an existing user")
}

func TestDeleteUser(t *testing.T) {
	setup()

	// Create a new user
	authService.CreateUser("userToDelete", "password123")

	// Delete the user
	err := authService.DeleteUser("userToDelete")
	assert.Nil(t, err, "Error should be nil while deleting the user")

	// Attempt to delete the user again
	err = authService.DeleteUser("userToDelete")
	assert.NotNil(t, err, "Error should not be nil when trying to delete a non-existent user")
}

func TestAddRoleToUser(t *testing.T) {
	setup()

	// Create a new user and role
	authService.CreateUser("testUser", "password123")
	authService.CreateRole("testRole", "read")

	// Add role to user
	err := authService.AddRoleToUser("testUser", "testRole")
	assert.Nil(t, err, "Error should be nil while adding role to user")

	// Check if user has the role
	hasRole, _ := authService.CheckUserRole("testUser", "testRole")
	assert.True(t, hasRole, "User should have the role added")
}

func TestAuthenticate(t *testing.T) {
	setup()

	// Create a new user
	authService.CreateUser("testUser", "password123")

	// Authenticate with correct credentials
	tokenDetails, err := authService.Authenticate("testUser", "password123")
	assert.Nil(t, err, "Error should be nil for valid credentials")
	assert.NotEqual(t, "", tokenDetails.Token, "Token should not be empty for valid credentials")

	// Authenticate with incorrect credentials
	_, err = authService.Authenticate("testUser", "wrongpassword")
	assert.NotNil(t, err, "Error should not be nil for invalid credentials")
}

func TestGetAllRoles(t *testing.T) {
	setup()

	// Create a new user and roles
	authService.CreateUser("testUser", "password123")
	authService.CreateRole("testRole1", "read")
	authService.CreateRole("testRole2", "write")

	// Add roles to user
	authService.AddRoleToUser("testUser", "testRole1")
	authService.AddRoleToUser("testUser", "testRole2")

	// Get all roles for the user
	roles, err := authService.GetAllRoles("testUser")
	assert.Nil(t, err, "Error should be nil while getting all roles")
	assert.Equal(t, 2, len(roles), "User should have two roles")
}

func TestInvalidateToken(t *testing.T) {
	setup()

	// Create a new user
	authService.CreateUser("testUser", "password123")

	// Authenticate to get a token
	tokenDetails, _ := authService.Authenticate("testUser", "password123")

	// Invalidate the token
	authService.(*InMemoryAuthService).InvalidateToken(tokenDetails.Token)

	// Check if token is still valid (this functionality is not provided in the given code, so this is a placeholder)
	// isValid := authService.IsTokenValid(tokenDetails.Token)
	// assert.False(t, isValid, "Token should be invalid after invalidation")
}
