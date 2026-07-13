package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type Claims struct {
	UserId   string `json:"user_id"`
	TenantId string `json:"tenant_id"`
	jwt.RegisteredClaims
}

type JwtManager struct {
	secretKey     []byte
	tokenDuration time.Duration
}

func NewJwtManager(secret string, duration time.Duration) *JwtManager {
	return &JwtManager{
		secretKey:     []byte(secret),
		tokenDuration: duration,
	}
}

func (m *JwtManager) Generate(userId, tenantId string) (string, error) {
	claims := Claims{
		UserId:   userId,
		TenantId: tenantId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now())},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(m.secretKey)
}

func (m *JwtManager) Verify(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return m.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
