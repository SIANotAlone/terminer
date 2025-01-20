package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"terminer/internal/models"
	"terminer/internal/repository"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

const (
	salt      = "s4564tadaskldnmlad13mnlak23mdlamd"
	signinkey = "kshabdnksndaskdnaksdnaksdan"
	tokenTTL  = 744 * time.Hour // Строк дія авторизації 744 годин = 31 день
)

type AuthService struct {
	repo repository.Authorization
}
type tokenClaims struct {
	jwt.StandardClaims
	UserID uuid.UUID `json:"id"`
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user models.UserRegistration) (uuid.UUID, error) {
	user.Password = s.generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(email string, password string) (string, error) {
	user, err := s.repo.GetUser(email, s.generatePasswordHash(password))
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user})
	return token.SignedString([]byte(signinkey))

}

func (s *AuthService) ParseToken(accessToken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signinkey), nil
	})
	if err != nil {
		return uuid.UUID{}, err

	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return uuid.UUID{}, errors.New("token claims are not of type *tokenClaims")
	}
	return claims.UserID, nil
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
