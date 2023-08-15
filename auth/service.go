// auth/service.go

package auth

import (
	"errors"
	//"fmt"

	"github.com/gogorush/simple_auth/utils"
)

type AuthService interface {
	CreateUser(username, password string) error
	DeleteUser(username string) error
	CreateRole(roleName string) error
	DeleteRole(roleName string) error
	AddRoleToUser(username, roleName string) error
	Authenticate(username, password string) (TokenDetails, error)
	CheckUserRole(tokenString, roleName string) (bool, error)
	GetAllRoles(tokenString string) ([]Role, error)
}

type InMemoryAuthService struct{}

func (s *InMemoryAuthService) CreateUser(username, password string) error {

	if _, exists := Users.Get(username); exists {
		return errors.New("user already exists")
	}
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	newUser := User{
		Username: username,
		Password: hashedPassword,
	}
	Users.Set(username, newUser)
	return nil
}

// DeleteUser deletes an existing user
func (s *InMemoryAuthService) DeleteUser(username string) error {
	if _, exists := Users.Get(username); !exists {
		return errors.New("user does not exist")
	}
	Users.Delete(username)
	return nil
}

// CreateRole creates a new role
func (s *InMemoryAuthService) CreateRole(roleName string) error {

	if _, exists := Roles.Get(roleName); exists {
		return errors.New("role already exists")
	}
	newRole := Role{Name: roleName}
	Roles.Set(roleName, newRole)
	return nil
}

// DeleteRole deletes an existing role
func (s *InMemoryAuthService) DeleteRole(roleName string) error {
	if _, exists := Roles.Get(roleName); !exists {
		return errors.New("role does not exist")
	}
	Roles.Delete(roleName)
	return nil
}

// AddRoleToUser associates a role with a user
func (s *InMemoryAuthService) AddRoleToUser(username string, roleName string) error {

	userInterface, userExists := Users.Get(username)
	if !userExists {
		return errors.New("user does not exist")
	}
	user := userInterface.(User) // type assertion

	// check if role exits
	roleInterface, roleExists := Roles.Get(roleName)
	if !roleExists {
		return errors.New("role does not exist")
	}
	role, ok := roleInterface.(Role) // type assertion
	if !ok {
		return errors.New("type assertion failed: not a Role")
	}

	for _, r := range user.Roles {
		if r.Name == roleName {
			return nil // If the role is already associated with the user, nothing happens
		}
	}
	user.Roles = append(user.Roles, role)
	Users.Set(username, user)
	return nil
}

// Authenticate validates user credentials
func (s *InMemoryAuthService) Authenticate(username, password string) (TokenDetails, error) {

	userInterface, userExists := Users.Get(username)
	if !userExists {
		return TokenDetails{}, errors.New("invalid credentials")
	}
	user := userInterface.(User) // type assertion

	if !utils.CheckPasswordHash(password, user.Password) {
		return TokenDetails{}, errors.New("invalid credentials")
	}
	return GenerateToken(username)
}

// CheckUserRole checks if a user has a specific role
func (s *InMemoryAuthService) CheckUserRole(tokenString, roleName string) (bool, error) {

	username, err := ValidateToken(tokenString)
	if err != nil {
		return false, err
	}

	userInterface, exists := Users.Get(username)
	if !exists {
		return false, errors.New("user does not exist")
	}

	_, exists = Roles.Get(roleName)
	// check if role exists
	if !exists {
		return false, errors.New("role does not exist")
	}

	user := userInterface.(User) // type assertion

	for _, role := range user.Roles {
		if role.Name == roleName {
			return true, nil
		}
	}
	return false, nil
}

// GetAllRoles retrieves all roles for a user
func (s *InMemoryAuthService) GetAllRoles(tokenString string) ([]Role, error) {

	username, err := ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	userInterface, exists := Users.Get(username)
	if !exists {
		return nil, errors.New("user does not exist")
	}
	user := userInterface.(User) // type assertion

	// check if role exists
	var roles []Role
	for _, role := range user.Roles {
		_, exists := Roles.Get(role.Name)
		if exists {
			roles = append(roles, role)
		}
	}
	return roles, nil
}
