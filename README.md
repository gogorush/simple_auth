# Simple Authentication Service
This project provides a simple authentication service using Go. It supports functionalities such as user creation, role management, and token-based authentication.

## Structure
The project is organized into the following structure:
```

├── Makefile
├── README.md
├── auth
│   ├── InMemoryTokenService.go - self made token generation, validation, and invalidation (using priority queue)
│   ├── handler.go - HTTP handlers for the authentication endpoints.
│   ├── handler_test.go - Tests for the HTTP handlers.
│   ├── model.go - Data models used in the authentication service.
│   ├── service.go - Business logic for authentication and authorization.
│   ├── service_test.go - Tests for the business logic.
│   ├── jwt_tokens.go - JWT token generation, validation, and invalidation.
│   └── tokens_test.go - Tests for JWT token functionalities.
├── go.mod
├── go.sum
├── main.go - Entry point for the application.
└── utils
    ├── PriorityQueue.go - A priority queue
    ├── PriorityQueue_test.go - Tests for priority queue
    ├── concurrent_map.go - A thread-safe concurrent map implementation.
    ├── concurrent_map_test.go - Tests for the concurrent map.
    ├── hasher.go - Utility for hashing passwords.
    └── hasher_test.go - Tests for the hashing utility.

```

## Getting Started
### Prerequisites
Install Go[https://go.dev/].

### Running the application
Clone the repository:
```
git clone <repository-url>
cd <repository-directory>
```
Build and run the application:
```
make build
make run
```
Running the Tests
```
make test
```

### Features
- User registration and authentication.
- Role management (create, delete, assign?).
- Token-based authentication using JWT.
- Thread-safe in-memory storage.

### External libs used
- golang-jwt (optional)

### Issues
- role with actual meanful use case? (only some user could create-user,  some could delte roles)
- http or https server? Not defined, should be use nginx like proxy for https? (implemented and comment our)
- While using JWTs for authentication was my original idea, storing these tokens in memory presents challenges, especially when dealing with multiple services. While an in-memory database like Redis can be used for this purpose, it introduces additional costs. (implement a in-memory storage with auto expire)

### License
This project is licensed under the MIT License.
