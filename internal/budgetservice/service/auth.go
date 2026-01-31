package service

import (
	"errors"
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
	
}
type tokenClaims struct {
	jwt.StandardClaims
	UserID uuid.UUID `json:"id"`
}

func NewAuthService() *AuthService {
	return &AuthService{}
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