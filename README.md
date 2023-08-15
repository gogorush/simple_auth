# Simple Authentication Service
This project provides a simple authentication service using Go. It supports functionalities such as user creation, role management, and token-based authentication.
**This is a simple version which covered all apis in the description, but not a real-world working solution, so I made another version is storage branch**

## Table of Contents

- [Project Structure](#-structure)
- [Getting Started](#-getting-started)
  - [Prerequisites](#prerequisites)
  - [Setup & Run](#setup--run)
  - [Testing](#-testing)
- [Features](#-features)
- [External Libraries](#-external-libs-used)
- [API Test Suite](#-api-test-suit)
- [Known Issues & Considerations](#-issues)
- [License](#license)

## 📂 Structure
The project is organized into the following structure:
```
├── Makefile
├── README.md
├── auth
│   ├── handler.go - HTTP handlers for the authentication endpoints.
│   ├── handler_test.go - Tests for the HTTP handlers.
│   ├── model.go - Data models used in the authentication service.
│   ├── service.go - Business logic for authentication and authorization.
│   ├── service_test.go - Tests for the business logic.
│   ├── tokens.go - JWT token generation, validation, and invalidation.
│   └── tokens_test.go - Tests for JWT token functionalities.
├── go.mod
├── go.sum
├── main.go - Entry point for the application.
├── simple_auth
└── utils
    ├── concurrent_map.go - A thread-safe concurrent map implementation.
    ├── concurrent_map_test.go - Tests for the concurrent map.
    ├── hasher.go - Utility for hashing passwords.
    └── hasher_test.go - Tests for the hashing utility.

```
## 🚀 Getting Started

### Prerequisites
Install Go[https://go.dev/].

### Setup & Run
Clone the repository:
```
git clone <repository-url> (or download the zip file from github)
cd <repository-directory>
```
Build and run the application:
```
make build
make run
```

### 🔍 Testing

Run the test suite with:
```
make test
```

### ✨ Features
- **User Management:** Register and authenticate users.
- **Role Management:** Create, delete, and assign roles to users.
- **Authentication:** Secure endpoints with JWT token-based authentication.
- **Storage:** Utilizes thread-safe in-memory storage.

### 📚 External libs used
- [golang-jwt](https://github.com/golang-jwt/jwt)

### 🧪 API Test Suit
Use the \`simple_auth_api.json\` Postman collection for testing all API endpoints. Just import it into Postman, and you're ready to go!


### 🚧 Issues
- **Role Functionality:** The actual use-case for roles (e.g., only certain roles can create users) isn't clear.
- **Role Deletion:** Removing roles can lead to inconsistencies, especially if users still possess the deleted role.
- **Server Type:** There's no clear distinction between HTTP and HTTPS. Consider using a proxy like nginx for HTTPS.
- **Token Storage:** While JWTs are efficient for authentication, storing them in memory isn't scalable. Although Redis can be a solution, it adds extra overhead.

### 📜 License
This project is licensed under the MIT License.
