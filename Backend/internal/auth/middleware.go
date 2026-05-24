package auth

import (
	"net/http"
	"slices"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
)

func JWTAuth(jwks *JWKSCache) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractBearerToken(r)
			if token == "" {
				apperror.HandlerespondUnauthorized(w, "Missing or invalid Authorization header")
				return
			}

			claims, err := ParseAndValidate(token, jwks)
			if err != nil {
				apperror.HandlerespondUnauthorized(w, "Invalid or expired token")
				return
			}

			ctx := ContextWithClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func VersionCheck(redis *RedisClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaimsFromContext(r.Context())
			if claims == nil {
				apperror.HandlerespondUnauthorized(w, "Unauthorized")
				return
			}

			redisVersion, err := redis.GetJwtVersion(claims.UserID)
			if err != nil {
				apperror.HandlerespondUnauthorized(w, "Unable to verify claims")
				return
			}

			if redisVersion > 0 && redisVersion > claims.JwtVersion {
				apperror.HandlerespondUnauthorized(w, "Token version is outdated, please refresh your token")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaimsFromContext(r.Context())
			if claims == nil {
				apperror.HandlerespondUnauthorized(w, "Unauthorized")
				return
			}

			if !slices.Contains(claims.Permissions, permission) {
				apperror.HandlerespondForbidden(w, "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
