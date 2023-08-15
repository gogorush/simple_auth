// auth/verify.go

package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	//"strings"
)

var Admin *User

func InitAdmin(admin User) {
	Admin = &User{}
	Admin.Username = admin.Username
	Admin.Password = admin.Password
	a := Role{
		Name:    "admin",
		Ability: ALL,
	}
	Admin.Roles = []Role{a}
	Roles.Set("admin", a)
	b, o := Roles.Get("admin")
	fmt.Println("init ", b, o)
}

func VerifyMethod(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var requestData UserRequest
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		path := r.URL.Path

		if requestData.Token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		roles, err := service.GetAllRoles(requestData.Token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var ability []string
		for _, role := range roles {
			ability = append(ability, role.Ability)
		}
		switch path {
		case "/create-user":
			if containsString(ability, ALL) || containsString(ability, CREATE_USER) {
				if requestData.Username == "" || requestData.Password == "" {
					http.Error(w, "Bad Request", http.StatusBadRequest)
					return
				}
				// Set requestData into the request's context
				ctx := context.WithValue(r.Context(), "requestData", requestData)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		case "/delete-user":
			if containsString(ability, ALL) || containsString(ability, DELETE_USER) {
				if requestData.Username == "" {
					http.Error(w, "Bad Request", http.StatusBadRequest)
					return
				}
				// Set requestData into the request's context
				ctx := context.WithValue(r.Context(), "requestData", requestData)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		case "/create-role":
			if containsString(ability, ALL) || containsString(ability, CREATE_ROLE) {
				if requestData.RoleName == "" || requestData.Ability == "" {
					http.Error(w, "Bad Request", http.StatusBadRequest)
					return
				}
				// Set requestData into the request's context
				ctx := context.WithValue(r.Context(), "requestData", requestData)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		case "/delete-role":
			if containsString(ability, ALL) || containsString(ability, DELETE_ROLE) {
				if requestData.RoleName == "" {
					http.Error(w, "Bad Request", http.StatusBadRequest)
					return
				}
				// Set requestData into the request's context
				ctx := context.WithValue(r.Context(), "requestData", requestData)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		case "/add-role-to-user":
			if containsString(ability, ALL) || containsString(ability, ADD_ROLE_TO_USER) {
				if requestData.RoleName == "" || requestData.Username == "" {
					http.Error(w, "Bad Request", http.StatusBadRequest)
					return
				}
				// Set requestData into the request's context
				ctx := context.WithValue(r.Context(), "requestData", requestData)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		case "/check-role":
			if containsString(ability, ALL) || containsString(ability, CHECK_ROLE) {
				if requestData.RoleName == "" || requestData.Token == "" {
					http.Error(w, "Bad Request", http.StatusBadRequest)
					return
				}
				// Set requestData into the request's context
				ctx := context.WithValue(r.Context(), "requestData", requestData)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		case "/get-all-roles":
			if containsString(ability, ALL) || containsString(ability, GET_ALL_ROLE) {
				if requestData.Token == "" {
					http.Error(w, "Bad Request", http.StatusBadRequest)
					return
				}
				// Set requestData into the request's context
				ctx := context.WithValue(r.Context(), "requestData", requestData)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		default:
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}
}

func containsString(arr []string, str string) bool {
    fmt.Println(arr, str)
	for _, a := range arr {
		if strings.Contains(a, str) {
			return true
		}
	}
	return false
}
