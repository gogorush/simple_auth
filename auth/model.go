// auth/model.go

package auth

import "github.com/gogorush/simple_auth/utils"

const (
	CREATE_USER = "create_user"
	DELETE_USER = "delete_user"

	CREATE_ROLE = "create_role"
	DELETE_ROLE = "delete_role"

	ADD_ROLE_TO_USER = "add_role_to_user"
	CHECK_ROLE       = "check_role"
	GET_ALL_ROLE     = "get_all_roles"

    ALL = "all"
)

type User struct {
	Username string
	Password string
	Roles    []Role
}

type Role struct {
	Name    string
	Ability string // TODO: give authorization to different apis?
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
