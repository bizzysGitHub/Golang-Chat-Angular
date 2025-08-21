package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token")

func secret() (string, error) {
	s := os.Getenv("JWT_HS256_SECRET")
	if s == "" {
		return "", errors.New("server misconfigured: JWT_HS256_SECRET not set")
	}
	return s, nil
}

// IssueJWT creates a short-lived HS256 token with sub, iat, exp.
func IssueJWT(userID string, ttl time.Duration) (string, error) {
	sec, err := secret()
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(ttl).Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString([]byte(sec))
}

// VerifyJWT as before…
func VerifyJWT(tokenString string) (string, error) {
	sec, err := secret()
	if err != nil {
		return "", err
	}

	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithIssuedAt(),
		jwt.WithExpirationRequired(),
	)
	claims := jwt.MapClaims{}
	token, err := parser.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return []byte(sec), nil
	})
	if err != nil || !token.Valid {
		return "", ErrInvalidToken
	}
	sub, _ := claims["sub"].(string)
	if sub == "" {
		return "", ErrInvalidToken
	}
	return sub, nil
}
