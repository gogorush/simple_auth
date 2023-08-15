// auth/InMemoryTokenService.go

package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/gogorush/simple_auth/utils"
	"time"
)

type InMemoryTokenService struct {
	tokens        *utils.ConcurrentMap
	tokenQueue    *utils.PriorityQueue
	tokenDuration time.Duration
}

func NewInMemoryTokenService(tokenDuration time.Duration) *InMemoryTokenService {
	s := &InMemoryTokenService{
		tokens:        utils.NewConcurrentMap(),
		tokenQueue:    utils.NewPriorityQueue(),
		tokenDuration: tokenDuration,
	}
	go s.expireTokens()
	return s
}

func (s *InMemoryTokenService) GenerateToken(username string) (TokenDetails, error) {
	tokenString, err := generateRandomTokenString(32) // Generate a random token string
	if err != nil {
		return TokenDetails{}, errors.New("error generate token")
	}
	expirationTime := time.Now().Add(s.tokenDuration).Unix()
	tokenDetails := TokenDetails{UserName: username, Token: tokenString, ExpiresAt: expirationTime}
	s.tokens.Set(tokenString, tokenDetails)
	s.tokenQueue.Push(&utils.Item{Token: tokenString, ExpiresAt: expirationTime})
	return tokenDetails, nil
}

func (s *InMemoryTokenService) ValidateToken(tokenString string) (string, error) {
	val, exists := s.tokens.Get(tokenString)
	if !exists {
		return "", errors.New("invalid token")
	}
	tokenDetails := val.(TokenDetails)
	if tokenDetails.ExpiresAt < time.Now().Unix() {
		s.InvalidateToken(tokenString)
		return "", errors.New("token expired")
	}
	return tokenDetails.UserName, nil
}

func (s *InMemoryTokenService) InvalidateToken(tokenString string) {
	s.tokens.Delete(tokenString)
}

func (s *InMemoryTokenService) expireTokens() {
	ticker := time.NewTicker(10 * time.Minute)
	for {
		<-ticker.C
		now := time.Now().Unix()
		for {
			item, exists := s.tokenQueue.Peek()
			if !exists {
				break
			}
			if item.ExpiresAt < now {
				s.InvalidateToken(item.Token)
				s.tokenQueue.Pop()
			} else {
				break
			}
		}
	}
}

func generateRandomTokenString(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(randomBytes), nil
}
