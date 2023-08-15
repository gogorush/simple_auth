// auth/handlers_test.go

package auth

import (
	"bytes"
	"encoding/json"
	//"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gogorush/simple_auth/utils"
)

func setupService() {

	// Reset the global maps to ensure a fresh state for each test
	Users = utils.NewConcurrentMap()
	Roles = utils.NewConcurrentMap()
	Tokens = utils.NewConcurrentMap()
	service = &InMemoryAuthService{
		//TokenSvc: NewJWTTokenService(JwtKey, TokenDuration),
		TokenSvc: NewInMemoryTokenService(TokenDuration),
	} // Reset to mock service for each test
}

func TestHandleCreateUser(t *testing.T) {
	setupService()

	reqBody := bytes.NewBufferString(`{"username":"testuser", "password":"testpass"}`)
	req, err := http.NewRequest("POST", "/create-user", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	HandleCreateUser(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestHandleDeleteUser(t *testing.T) {
	setupService()

	// First, create the user
	createReqBody := bytes.NewBufferString(`{"username":"testuser", "password":"testpass"}`)
	createReq, err := http.NewRequest("POST", "/create-user", createReqBody)
	if err != nil {
		t.Fatal(err)
	}
	createRR := httptest.NewRecorder()
	HandleCreateUser(createRR, createReq)
	if status := createRR.Code; status != http.StatusCreated {
		t.Fatalf("Failed to create mock user: got %v want %v", status, http.StatusCreated)
	}

	// Now, delete the user
	deleteReqBody := bytes.NewBufferString(`{"username":"testuser"}`)
	deleteReq, err := http.NewRequest("POST", "/delete-user", deleteReqBody)
	if err != nil {
		t.Fatal(err)
	}
	deleteRR := httptest.NewRecorder()
	HandleDeleteUser(deleteRR, deleteReq)
	if status := deleteRR.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestHandleCreateRole(t *testing.T) {
	setupService()

	reqBody := bytes.NewBufferString(`{"roleName":"testrole"}`)
	req, err := http.NewRequest("POST", "/create-role", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	HandleCreateRole(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestHandleDeleteRole(t *testing.T) {
	setupService()

	// First, create the role
	createReqBody := bytes.NewBufferString(`{"roleName":"testrole"}`)
	createReq, err := http.NewRequest("POST", "/create-role", createReqBody)
	if err != nil {
		t.Fatal(err)
	}
	createRR := httptest.NewRecorder()
	HandleCreateRole(createRR, createReq)
	if status := createRR.Code; status != http.StatusCreated {
		t.Fatalf("Failed to create mock role: got %v want %v", status, http.StatusCreated)
	}

	// Now, delete the role
	deleteReqBody := bytes.NewBufferString(`{"roleName":"testrole"}`)
	deleteReq, err := http.NewRequest("POST", "/delete-role", deleteReqBody)
	if err != nil {
		t.Fatal(err)
	}
	deleteRR := httptest.NewRecorder()
	HandleDeleteRole(deleteRR, deleteReq)
	if status := deleteRR.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestHandleAddRoleToUser(t *testing.T) {
	setupService()

	// Create a mock user
	userReqBody := bytes.NewBufferString(`{"username":"testuser", "password":"testpass"}`)
	userReq, err := http.NewRequest("POST", "/create-user", userReqBody)
	if err != nil {
		t.Fatal(err)
	}

	userRR := httptest.NewRecorder()
	HandleCreateUser(userRR, userReq)

	if status := userRR.Code; status != http.StatusCreated {
		t.Fatalf("Failed to create mock user: got %v want %v", status, http.StatusCreated)
	}

	// Create a mock role
	roleReqBody := bytes.NewBufferString(`{"roleName":"testrole"}`)
	roleReq, err := http.NewRequest("POST", "/create-role", roleReqBody)
	if err != nil {
		t.Fatal(err)
	}

	roleRR := httptest.NewRecorder()
	HandleCreateRole(roleRR, roleReq)

	if status := roleRR.Code; status != http.StatusCreated {
		t.Fatalf("Failed to create mock role: got %v want %v", status, http.StatusCreated)
	}

	// Test adding role to user
	reqBody := bytes.NewBufferString(`{"username":"testuser", "roleName":"testrole"}`)
	req, err := http.NewRequest("POST", "/add-role", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	HandleAddRoleToUser(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code for adding role to user: got %v want %v", status, http.StatusCreated)
	}

	// Test adding the same role again to the user
	//dupRoleRR := httptest.NewRecorder()
	HandleAddRoleToUser(rr, req)

	// Assuming the handler doesn't throw an error for duplicate roles, it should return StatusCreated again
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code for adding duplicate role to user: got %v want %v", status, http.StatusCreated)
	}
}

func TestHandleAuthenticate(t *testing.T) {
	setupService()

	// 1. Authenticate a non-existent user
	nonExistentUserReq := bytes.NewBufferString(`{"username":"nonexistentuser", "password":"anyrandompass"}`)
	req1, err := http.NewRequest("POST", "/authenticate", nonExistentUserReq)
	if err != nil {
		t.Fatal(err)
	}

	rr1 := httptest.NewRecorder()
	HandleAuthenticate(rr1, req1)

	if status := rr1.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code for non-existent user: got %v want %v", status, http.StatusBadRequest)
	}

	// 2. Create a user and authenticate with the correct password
	userReqBody := bytes.NewBufferString(`{"username":"testuser", "password":"testpass"}`)
	userReq, err := http.NewRequest("POST", "/create-user", userReqBody)
	if err != nil {
		t.Fatal(err)
	}

	userRR := httptest.NewRecorder()
	HandleCreateUser(userRR, userReq)

	if status := userRR.Code; status != http.StatusCreated {
		t.Fatalf("Failed to create mock user: got %v want %v", status, http.StatusCreated)
	}

	authReqBody := bytes.NewBufferString(`{"username":"testuser", "password":"testpass"}`)
	req2, err := http.NewRequest("POST", "/authenticate", authReqBody)
	if err != nil {
		t.Fatal(err)
	}

	rr2 := httptest.NewRecorder()
	HandleAuthenticate(rr2, req2)

	if status := rr2.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code for correct password: got %v want %v", status, http.StatusOK)
	}
	var responseToken TokenDetails
	err = json.NewDecoder(rr2.Body).Decode(&responseToken)
	if err != nil {
		t.Fatalf("Failed to decode the response token: %v", err)
	}
	if responseToken.Token == "" {
		t.Errorf("Expected a token in the response, but got none")
	}
	// 3. Authenticate with the wrong password
	authWrongPassReqBody := bytes.NewBufferString(`{"username":"testuser", "password":"wrongpass"}`)
	req3, err := http.NewRequest("POST", "/authenticate", authWrongPassReqBody)
	if err != nil {
		t.Fatal(err)
	}

	rr3 := httptest.NewRecorder()
	HandleAuthenticate(rr3, req3)

	if status := rr3.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code for wrong password: got %v want %v", status, http.StatusBadRequest)
	}
}
func TestHandleGenerateToken(t *testing.T) {

	setupService()

	// Create a mock user
	userReqBody := bytes.NewBufferString(`{"username":"testuser", "password":"testpass"}`)
	userReq, err := http.NewRequest("POST", "/create-user", userReqBody)
	if err != nil {
		t.Fatal(err)
	}
	userRR := httptest.NewRecorder()
	HandleCreateUser(userRR, userReq)
	if status := userRR.Code; status != http.StatusCreated {
		t.Fatalf("Failed to create mock user: got %v want %v", status, http.StatusCreated)
	}

	// Now test token generation for the created user
	reqBody := bytes.NewBufferString(`{"username":"testuser", "password":"testpass"}`)
	req, err := http.NewRequest("POST", "/generate-token", reqBody)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	HandleGenerateToken(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var tokenDetails TokenDetails
	err = json.NewDecoder(rr.Body).Decode(&tokenDetails)
	if err != nil {
		t.Fatal("Failed decoding response body")
	}
	if tokenDetails.Token == "" {
		t.Errorf("Expected a token but got none")
	}
}

func TestHandleInvalidateToken(t *testing.T) {
	setupService()

	tokenDetails, _ := service.Authenticate("testuser", "testpass")
	reqBody := bytes.NewBufferString(`{"token":"` + tokenDetails.Token + `"}`)
	req, err := http.NewRequest("POST", "/invalidate-token", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	HandleInvalidateToken(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestHandleCheckRole(t *testing.T) {
	setupService()

	// Create a mock user
	userReqBody := bytes.NewBufferString(`{"username":"testuser", "password":"testpass"}`)
	userReq, err := http.NewRequest("POST", "/create-user", userReqBody)
	if err != nil {
		t.Fatal(err)
	}
	userRR := httptest.NewRecorder()
	HandleCreateUser(userRR, userReq)
	if status := userRR.Code; status != http.StatusCreated {
		t.Fatalf("Failed to create mock user: got %v want %v", status, http.StatusCreated)
	}

	// Authenticate the user to get a token
	authReqBody := bytes.NewBufferString(`{"username":"testuser", "password":"testpass"}`)
	authReq, err := http.NewRequest("POST", "/authenticate", authReqBody)
	if err != nil {
		t.Fatal(err)
	}
	authRR := httptest.NewRecorder()
	HandleAuthenticate(authRR, authReq)
	if status := authRR.Code; status != http.StatusOK {
		t.Fatalf("Failed to authenticate mock user: got %v want %v", status, http.StatusOK)
	}
	var tokenDetails TokenDetails
	err = json.NewDecoder(authRR.Body).Decode(&tokenDetails)
	if err != nil {
		t.Fatalf("Failed to decode authentication response: %v", err)
	}

	// Create a mock role
	roleReqBody := bytes.NewBufferString(`{"roleName":"testrole"}`)
	roleReq, err := http.NewRequest("POST", "/create-role", roleReqBody)
	if err != nil {
		t.Fatal(err)
	}
	roleRR := httptest.NewRecorder()
	HandleCreateRole(roleRR, roleReq)
	if status := roleRR.Code; status != http.StatusCreated {
		t.Fatalf("Failed to create mock role: got %v want %v", status, http.StatusCreated)
	}

	// Assign role to user
	assignReqBody := bytes.NewBufferString(`{"username":"testuser", "roleName":"testrole"}`)
	assignReq, err := http.NewRequest("POST", "/add-role", assignReqBody)
	if err != nil {
		t.Fatal(err)
	}
	assignRR := httptest.NewRecorder()
	HandleAddRoleToUser(assignRR, assignReq)
	if status := assignRR.Code; status != http.StatusCreated {
		t.Fatalf("Failed to assign role to user: got %v want %v", status, http.StatusCreated)
	}

	//fmt.Println(tokenDetails)
	// Test checking role for user using the token
	checkReqBody := bytes.NewBufferString(`{"token":"` + tokenDetails.Token + `", "roleName":"testrole"}`)
	checkReq, err := http.NewRequest("POST", "/check-role", checkReqBody)
	if err != nil {
		t.Fatal(err)
	}
	checkRR := httptest.NewRecorder()
	HandleCheckRole(checkRR, checkReq)
	if status := checkRR.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var roleCheck map[string]bool
	err = json.NewDecoder(checkRR.Body).Decode(&roleCheck)
	if err != nil {
		t.Fatal("Failed decoding response body")
	}
	if !roleCheck["hasRole"] {
		t.Errorf("Expected user to have role but they did not")
	}
}

func TestHandleGetAllRoles(t *testing.T) {
	setupService()

	// Create a mock user
	userReqBody := bytes.NewBufferString(`{"username":"testuser", "password":"testpass"}`)
	userReq, err := http.NewRequest("POST", "/create-user", userReqBody)
	if err != nil {
		t.Fatal(err)
	}
	userRR := httptest.NewRecorder()
	HandleCreateUser(userRR, userReq)
	if status := userRR.Code; status != http.StatusCreated {
		t.Fatalf("Failed to create mock user: got %v want %v", status, http.StatusCreated)
	}

	// Authenticate the user to get a token
	authReqBody := bytes.NewBufferString(`{"username":"testuser", "password":"testpass"}`)
	authReq, err := http.NewRequest("POST", "/authenticate", authReqBody)
	if err != nil {
		t.Fatal(err)
	}
	authRR := httptest.NewRecorder()
	HandleAuthenticate(authRR, authReq)
	if status := authRR.Code; status != http.StatusOK {
		t.Fatalf("Failed to authenticate mock user: got %v want %v", status, http.StatusOK)
	}
	var tokenDetails TokenDetails
	err = json.NewDecoder(authRR.Body).Decode(&tokenDetails)
	if err != nil {
		t.Fatalf("Failed to decode authentication response: %v", err)
	}

	// Create some roles for testing and assign them to the user
	roleNames := []string{"role1", "role2", "role3"}
	for _, roleName := range roleNames {
		roleReqBody := bytes.NewBufferString(`{"roleName":"` + roleName + `"}`)
		roleReq, err := http.NewRequest("POST", "/create-role", roleReqBody)
		if err != nil {
			t.Fatal(err)
		}
		roleRR := httptest.NewRecorder()
		HandleCreateRole(roleRR, roleReq)
		if status := roleRR.Code; status != http.StatusCreated {
			t.Fatalf("Failed to create role %s: got %v want %v", roleName, status, http.StatusCreated)
		}

		// Assign role to user
		assignReqBody := bytes.NewBufferString(`{"username":"testuser", "roleName":"` + roleName + `"}`)
		assignReq, err := http.NewRequest("POST", "/add-role", assignReqBody)
		if err != nil {
			t.Fatal(err)
		}
		assignRR := httptest.NewRecorder()
		HandleAddRoleToUser(assignRR, assignReq)
		if status := assignRR.Code; status != http.StatusCreated {
			t.Fatalf("Failed to assign role %s to user: got %v want %v", roleName, status, http.StatusCreated)
		}
	}

	// Test getting all roles using the obtained token
	reqBody := bytes.NewBufferString(`{"token":"` + tokenDetails.Token + `"}`)
	req, err := http.NewRequest("POST", "/get-all-roles", reqBody)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	HandleGetAllRoles(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var roles []Role
	err = json.NewDecoder(rr.Body).Decode(&roles)
	if err != nil {
		t.Fatal("Failed decoding response body")
	}
	if len(roles) != len(roleNames) {
		t.Errorf("Expected %d roles but got %d", len(roleNames), len(roles))
	}
}
