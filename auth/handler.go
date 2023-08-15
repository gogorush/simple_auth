// auth/handlers.go

package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var service AuthService = &InMemoryAuthService{
	//TokenSvc: NewJWTTokenService(JwtKey, TokenDuration),
	TokenSvc: NewInMemoryTokenService(TokenDuration),
} // Create an instance of the AuthService

type UserRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	RoleName string `json:"roleName,omitempty"`
	Token    string `json:"token,omitempty"`
	Ability  string `json:"ability,omitempty"`
}

func HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	ctxData := r.Context().Value("requestData")
	if ctxData == nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	requestData, ok := ctxData.(UserRequest)
	if !ok {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	err := service.CreateUser(requestData.Username, requestData.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	ctxData := r.Context().Value("requestData")
	if ctxData == nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	requestData, ok := ctxData.(UserRequest)
	if !ok {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err := service.DeleteUser(requestData.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HandleCreateRole(w http.ResponseWriter, r *http.Request) {
	ctxData := r.Context().Value("requestData")
	if ctxData == nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	requestData, ok := ctxData.(UserRequest)
	if !ok {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err := service.CreateRole(requestData.RoleName, requestData.Ability)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func HandleDeleteRole(w http.ResponseWriter, r *http.Request) {
	ctxData := r.Context().Value("requestData")
	if ctxData == nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	requestData, ok := ctxData.(UserRequest)
	if !ok {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err := service.DeleteRole(requestData.RoleName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HandleAddRoleToUser(w http.ResponseWriter, r *http.Request) {
	ctxData := r.Context().Value("requestData")
	if ctxData == nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	requestData, ok := ctxData.(UserRequest)
	if !ok {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err := service.AddRoleToUser(requestData.Username, requestData.RoleName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func HandleAuthenticate(w http.ResponseWriter, r *http.Request) {
	var requestData UserRequest
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	tokenDetails, err := service.Authenticate(requestData.Username, requestData.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(tokenDetails)
}

func HandleGenerateToken(w http.ResponseWriter, r *http.Request) {
	var requestData UserRequest
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	tokenDetails, err := service.Authenticate(requestData.Username, requestData.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(tokenDetails)
}

func HandleInvalidateToken(w http.ResponseWriter, r *http.Request) {
	var requestData UserRequest
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	service.(*InMemoryAuthService).InvalidateToken(requestData.Token)
	w.WriteHeader(http.StatusOK)
}

func HandleCheckRole(w http.ResponseWriter, r *http.Request) {
	ctxData := r.Context().Value("requestData")
	if ctxData == nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	requestData, ok := ctxData.(UserRequest)
	if !ok {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	hasRole, err := service.CheckUserRole(requestData.Token, requestData.RoleName)
	if err != nil {
		if err.Error() == "token has expired" {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"hasRole": hasRole})
}

func HandleGetAllRoles(w http.ResponseWriter, r *http.Request) {
	ctxData := r.Context().Value("requestData")
	if ctxData == nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	requestData, ok := ctxData.(UserRequest)
	if !ok {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

    fmt.Println("requestData", requestData)
	roles, err := service.GetAllRoles(requestData.Token)
    fmt.Println("HandleGetAllRoles: ", roles, err)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(roles) == 0 {
		w.WriteHeader(http.StatusOK)
		return
	}

	json.NewEncoder(w).Encode(roles)
}
