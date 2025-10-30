package auth

import (
	"context"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/warmhouse/warmhouse_authproxy/internal/config"
	"github.com/warmhouse/warmhouse_authproxy/internal/models"
)

type Service struct {
	secrets *config.Secrets
	conf    *config.Config
}

func NewService(secrets *config.Secrets, conf *config.Config) *Service {
	return &Service{secrets: secrets, conf: conf}
}

func (s *Service) GenerateToken(ctx context.Context, user models.User) (string, time.Time, error) {
	log.Println("config:", s.conf, "secrets:", s.secrets)

	claims := &JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.conf.JwtDurationMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := jwtToken.SignedString([]byte(s.secrets.JWTSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return token, claims.ExpiresAt.Time, nil
}

func (s *Service) ValidateToken(ctx context.Context, token string) (string, error) {
	claims := &JWTClaims{}

	jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return []byte(s.secrets.JWTSecret), nil
	})

	if err != nil {
		return "", err
	}

	log.Println("jwtToken:", jwtToken)
	log.Println("claims:", jwtToken.Claims)

	if time.Now().After(claims.ExpiresAt.Time) {
		return "", ErrTokenExpired
	}

	return claims.UserID.String(), nil
}
