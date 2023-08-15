// auth/service_test.go

package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/gogorush/simple_auth/utils"
)

var authService AuthService = &InMemoryAuthService{}

func setup() {

	// Reset the global maps to ensure a fresh state for each test
	Users = utils.NewConcurrentMap()
	Roles = utils.NewConcurrentMap()
	Tokens = utils.NewConcurrentMap()
	authService = &InMemoryAuthService{} // Reset to mock service for each test
}

func TestCreateUser(t *testing.T) {
    setup()
	err := authService.CreateUser("testuser1", "password123")

	assert.Nil(t, err, "Error should be nil")

	// Attempting to create the same user should return an error
	err = authService.CreateUser("testuser1", "password123")
	assert.NotNil(t, err, "Error should not be nil")
}

func TestDeleteUser(t *testing.T) {
    setup()
	authService.CreateUser("userToDelete", "password123")

	err := authService.DeleteUser("userToDelete")
	assert.Nil(t, err, "Error should be nil")

	// Attempting to delete the user again should return an error
	err = authService.DeleteUser("userToDelete")
	assert.NotNil(t, err, "Error should not be nil")
}

func TestCreateRole(t *testing.T) {
    setup()
	err := authService.CreateRole("testRole")
	assert.Nil(t, err, "Error should be nil")

	// Attempting to create the same role should return an error
	err = authService.CreateRole("testRole")
	assert.NotNil(t, err, "Error should not be nil")
}

func TestDeleteRole(t *testing.T) {
    setup()
	authService.CreateRole("roleToDelete")

	err := authService.DeleteRole("roleToDelete")
	assert.Nil(t, err, "Error should be nil")

	// Attempting to delete the role again should return an error
	err = authService.DeleteRole("roleToDelete")
	assert.NotNil(t, err, "Error should not be nil")
}

func TestAddRoleToUser(t *testing.T) {
    setup()
	authService.CreateUser("userForRole", "password123")
	authService.CreateRole("roleToAdd")

	err := authService.AddRoleToUser("userForRole", "roleToAdd")
	assert.Nil(t, err, "Error should be nil")

	// Attempting to add the role again should not return an error (idempotent operation)
	err = authService.AddRoleToUser("userForRole", "roleToAdd")
	assert.Nil(t, err, "Error should be nil")

	err = authService.AddRoleToUser("userForRole", "notExistRole")
	assert.NotNil(t, err, "Error should not be nil")
}

func TestAuthenticate(t *testing.T) {
    setup()
	authService.CreateUser("userToAuth", "password123")

	_, err := authService.Authenticate("userToAuth", "password123")
	assert.Nil(t, err, "Error should be nil")

	_, err = authService.Authenticate("userToAuth", "wrongpassword")
	assert.NotNil(t, err, "Error should not be nil")
}

func TestCheckUserRole(t *testing.T) {
    setup()
	authService.CreateUser("userForRoleCheck", "password123")
	authService.CreateRole("roleToCheck")
    //authService.AddRoleToUser("userForRoleCheck", "roleToCheck")
	tokenDetails, _ := authService.Authenticate("userForRoleCheck", "password123")

	hasRole, err := authService.CheckUserRole(tokenDetails.Token, "roleToCheck")
	assert.Nil(t, err, "Error should be nil")
	assert.False(t, hasRole, "User should not have the role")

	authService.AddRoleToUser("userForRoleCheck", "roleToCheck")
	hasRole, err = authService.CheckUserRole(tokenDetails.Token, "roleToCheck")
	assert.Nil(t, err, "Error should be nil")
	assert.True(t, hasRole, "User should have the role")

	hasRole, err = authService.CheckUserRole(tokenDetails.Token, "notExistRole")
	assert.NotNil(t, err, "Error should be nil")
	assert.False(t, hasRole, "User should not have the role")
}

func TestGetAllRoles(t *testing.T) {
    setup()
	authService.CreateUser("userForGetAllRoles", "password123")
	authService.CreateRole("role1")
	authService.CreateRole("role2")
	tokenDetails, _ := authService.Authenticate("userForGetAllRoles", "password123")

	roles, err := authService.GetAllRoles(tokenDetails.Token)
	assert.Nil(t, err, "Error should be nil")
	assert.Len(t, roles, 0, "User should have no roles")

	authService.AddRoleToUser("userForGetAllRoles", "role1")
	authService.AddRoleToUser("userForGetAllRoles", "role2")

	roles, err = authService.GetAllRoles(tokenDetails.Token)
	assert.Nil(t, err, "Error should be nil")
	assert.Len(t, roles, 2, "User should have 2 roles")
}
