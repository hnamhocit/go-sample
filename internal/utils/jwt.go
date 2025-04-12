package utils

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var (
	key string
	t   *jwt.Token
	s   string
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func GenerateToken(sub int32, secretKey string) (string, error) {
	key = os.Getenv(secretKey)
	t = jwt.NewWithClaims(jwt.SigningMethodES256,
		jwt.MapClaims{
			"sub": sub,
		})

	s, err := t.SignedString(key)

	return s, err
}

func VerifyToken(tokenString string, secretKey string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv(secretKey)), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func GenerateTokens(sub int32) (*Tokens, error) {
	accessToken, err := GenerateToken(sub, "JWT_ACCESS_SECRET")
	if err != nil {
		return nil, err
	}

	refreshToken, err := GenerateToken(sub, "JWT_REFRESH_SECRET")
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
