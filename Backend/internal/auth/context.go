package auth

import (
	"context"
)

type contextKey string

const claimsKey contextKey = "claims"

type Claims struct {
	UserID      string   `json:"sub"`
	Email       string   `json:"email"`
	TenantID    string   `json:"tenant_id"`
	Permissions []string `json:"permissions"`
	JwtVersion  int      `json:"jwt_version"`
	IssuedAt    int64    `json:"iat"`
	ExpiresAt   int64    `json:"exp"`
}

func ContextWithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

func GetClaimsFromContext(ctx context.Context) *Claims {
	claims, ok := ctx.Value(claimsKey).(*Claims)
	if !ok {
		return nil
	}
	return claims
}
