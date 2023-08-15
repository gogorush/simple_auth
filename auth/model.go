// auth/model.go

package auth

import "github.com/gogorush/simple_auth/utils"

type User struct {
	Username string
	Password string
	Roles    []Role
}

type Role struct {
	Name    string
	//Ability []string
	//Status  bool
}

type TokenDetails struct {
    UserName string
	Token     string
	ExpiresAt int64
}

type UserTokenDetails struct {
	Token     string
	ExpiresAt int64
    UserName string
}

var (
	Users  = utils.NewConcurrentMap()
	Roles  = utils.NewConcurrentMap()
	Tokens = utils.NewConcurrentMap()
)
