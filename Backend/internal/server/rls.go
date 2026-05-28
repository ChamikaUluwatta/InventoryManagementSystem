package server

import (
	"net/http"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/auth"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TenantRLS(pool *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := auth.GetClaimsFromContext(r.Context())
			if claims == nil || claims.TenantID == "" {
				apperror.HandlerespondUnauthorized(w, "Missing tenant id")
				return
			}

			tx, err := pool.Begin(r.Context())
			if err != nil {
				apperror.HandleError(w, apperror.Internal("Database error", err))
				return
			}
			defer tx.Rollback(r.Context())

			_, err = tx.Exec(r.Context(),
				"SELECT set_config('app.current_tenant_id', $1, true)",
				claims.TenantID,
			)
			if err != nil {
				apperror.HandleError(w, apperror.Internal("Database error", err))
				return
			}

			ctx := database.SetTx(r.Context(), tx)
			next.ServeHTTP(w, r.WithContext(ctx))

			if err := tx.Commit(r.Context()); err != nil {
				// response already written
			}
		})
	}
}
