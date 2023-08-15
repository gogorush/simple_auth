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
    //Ability []string // TODO: give authorization to different apis?
    //Status  bool // TODO: not a good idea to remove a role, marked should be better
}

type TokenDetails struct {
	UserName  string
	Token     string
	ExpiresAt int64
}

var (
	Users  = utils.NewConcurrentMap()
	Roles  = utils.NewConcurrentMap()
	Tokens = utils.NewConcurrentMap()
)
