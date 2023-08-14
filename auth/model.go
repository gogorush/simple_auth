// auth/model.go

package auth

import "github.com/gogorush/simple_auth/utils"

type User struct {
	Username string
	Password string
	Roles    []Role
}

type Role struct {
	Name string
}

type TokenDetails struct {
	Token     string
	ExpiresAt int64
}

var (
	Users  = utils.NewConcurrentMap()
	Roles  = utils.NewConcurrentMap()
	Tokens = utils.NewConcurrentMap()
)

