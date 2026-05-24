package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func extractBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return ""
	}

	return parts[1]
}

func ValidateClaims(claims *Claims) error {
	if claims.UserID == "" {
		return fmt.Errorf("missing subject claim")
	}
	if claims.Email == "" {
		return fmt.Errorf("missing email claim")
	}
	if claims.JwtVersion < 1 {
		return fmt.Errorf("invalid permissions version")
	}
	return nil
}
