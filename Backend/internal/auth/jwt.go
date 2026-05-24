package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func ParseAndValidate(tokenString string, jwks *JWKSCache) (*Claims, error) {
	claims, err := parseToken(tokenString, jwks)
	if err != nil {
		if refreshErr := jwks.Refresh(); refreshErr != nil {
			return nil, fmt.Errorf("failed to refresh JWKS: %w", refreshErr)
		}
		claims, err = parseToken(tokenString, jwks)
		if err != nil {
			return nil, err
		}
	}
	return claims, nil
}

func parseToken(tokenString string, jwks *JWKSCache) (*Claims, error) {
	publicKey, err := jwks.GetPublicKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %w", err)
	}

	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	if err := ValidateClaims(&claims); err != nil {
		return nil, fmt.Errorf("invalid claims: %w", err)
	}

	return &claims, nil
}
