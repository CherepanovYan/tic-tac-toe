package domain

import (
	"time"

	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtService struct {
	secret []byte
}

func NewJwtService(secret string) *JwtService {
	return &JwtService{
		secret: []byte(secret),
	}
}

func (s *JwtService) GenerateAccessToken(userUUID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"uuid": userUUID.String(),
		"type": "access",
		"exp":  time.Now().Add(15 * time.Minute).Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.secret)
}

func (s *JwtService) GenerateRefreshToken(userUUID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"uuid": userUUID.String(),
		"type": "refresh",
		"exp":  time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.secret)
}

func (s *JwtService) ValidateAccessToken(
	tokenString string,
) error {
	return s.validateToken(
		tokenString,
		"access",
	)
}

func (s *JwtService) ValidateRefreshToken(
	tokenString string,
) error {
	return s.validateToken(
		tokenString,
		"refresh",
	)
}

func (s *JwtService) GetUUID(
	tokenString string,
) (uuid.UUID, error) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New(
					"unexpected signing method",
				)
			}

			return s.secret, nil
		},
	)

	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("invalid claims")
	}

	value, ok := claims["uuid"].(string)
	if !ok {
		return uuid.Nil, errors.New("uuid not found")
	}

	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, errors.New("invalid uuid")
	}

	return id, nil
}

func (s *JwtService) validateToken(
	tokenString string,
	tokenType string,
) error {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New(
					"unexpected signing method",
				)
			}

			return s.secret, nil
		},
	)

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid claims")
	}

	actualType, ok := claims["type"].(string)
	if !ok {
		return errors.New("token type not found")
	}

	if actualType != tokenType {
		return errors.New("invalid token type")
	}

	return nil
}
