// auth/tokens.go

package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var tokenDuration = 2 * time.Hour

var jwtKey = []byte("your-secret-key") // This should ideally be more secure and not hardcoded

// GenerateToken generates a JWT for the given user
func GenerateToken(username string) (TokenDetails, error) {
    expirationTime := time.Now().Add(tokenDuration).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": username,
		"exp":  expirationTime,
	})
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return TokenDetails{}, err
	}
	Tokens.Set(tokenString, username)
	return TokenDetails{Token: tokenString, ExpiresAt: expirationTime}, nil
}

// ValidateToken checks the given token's validity
func ValidateToken(tokenString string) (string, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
    //fmt.Println(token, " hello ", err)

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}
	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", errors.New("invalid token claims exp")
	}
	if int64(exp) < time.Now().Unix() {
        InvalidateToken(tokenString)
		return "", errors.New("token expired here")
	}

	username, ok := claims["user"].(string)
	if !ok {
		return "", errors.New("invalid token claims username")
	}

	return username, nil
}

// InvalidateToken removes a token, making it invalid
func InvalidateToken(tokenString string) {
	Tokens.Delete(tokenString)
}

// SetTokenDuration allows changing the token duration for testing purposes
func setTokenDuration(duration time.Duration) {
	tokenDuration = duration
}

