// auth/tokens.go

package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	GenerateToken(username string) (TokenDetails, error)
	ValidateToken(tokenString string) (string, error)
	InvalidateToken(tokenString string)
}

type JWTTokenService struct {
	jwtKey        []byte
	tokenDuration time.Duration
}

func NewJWTTokenService(jwtKey []byte, tokenDuration time.Duration) *JWTTokenService {
	return &JWTTokenService{jwtKey: jwtKey, tokenDuration: tokenDuration}
}

var TokenDuration time.Duration = 2 * time.Hour

var JwtKey []byte // This should ideally be more secure and not hardcoded

// GenerateToken generates a JWT for the given user
func (s *JWTTokenService) GenerateToken(username string) (TokenDetails, error) {

	expirationTime := time.Now().Add(TokenDuration).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": username,
		"exp":  expirationTime,
	})
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return TokenDetails{}, err
	}
	Tokens.Set(tokenString, username)
	return TokenDetails{Token: tokenString, ExpiresAt: expirationTime}, nil
}

// ValidateToken checks the given token's validity
func (s *JWTTokenService) ValidateToken(tokenString string) (string, error) {

	_, ok := Tokens.Get(tokenString)
	if !ok {
		return "", errors.New("invalid token")
	}
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
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
		s.InvalidateToken(tokenString)
		return "", errors.New("token expired here")
	}

	username, ok := claims["user"].(string)
	if !ok {
		return "", errors.New("invalid token claims username")
	}

	return username, nil
}

// InvalidateToken removes a token, making it invalid
func (s *JWTTokenService) InvalidateToken(tokenString string) {
	Tokens.Delete(tokenString)
}

// SetTokenDuration allows changing the token duration for testing purposes
// NOTE: for test only
func setTokenDuration(duration time.Duration) {
	TokenDuration = duration
}
