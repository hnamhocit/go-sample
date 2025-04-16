package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	Sub          int32 `json:"sub"`
	TokenVersion int32 `json:"token_version"`
	jwt.RegisteredClaims
}

func GenerateToken(sub int32, token_version int32, secretKey string) (string, error) {
	key := os.Getenv(secretKey)
	if key == "" {
		return "", fmt.Errorf("missing secret key: %s", secretKey)
	}

	expiredAt := 1 // default expiration time
	if secretKey == "JWT_REFRESH_SECRET" {
		expiredAt = 24 * 7 // 7 days
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Sub:          sub,
		TokenVersion: token_version,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiredAt) * time.Hour)),
		},
	})

	signedString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedString, nil
}

func VerifyToken(tokenString string, secretKey string) (*Claims, error) {
	key := os.Getenv(secretKey)
	if key == "" {
		return nil, fmt.Errorf("missing secret key: %s", secretKey)
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func GenerateTokens(sub, token_version int32) (*Tokens, error) {
	accessToken, err := GenerateToken(sub, token_version, "JWT_ACCESS_SECRET")
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := GenerateToken(sub, token_version, "JWT_REFRESH_SECRET")
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
