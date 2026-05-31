package server

import (
	"bytes"
	"log"
	"net/http"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/auth"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

type responseBuffer struct {
	http.ResponseWriter
	statusCode  int
	body        bytes.Buffer
	wroteHeader bool
}

func (rb *responseBuffer) WriteHeader(statusCode int) {
	if rb.wroteHeader {
		return
	}
	rb.statusCode = statusCode
	rb.wroteHeader = true
}

func (rb *responseBuffer) Write(b []byte) (int, error) {
	if !rb.wroteHeader {
		rb.WriteHeader(http.StatusOK)
	}
	return rb.body.Write(b)
}

func (rb *responseBuffer) flush() {
	if rb.wroteHeader {
		rb.ResponseWriter.WriteHeader(rb.statusCode)
	}
	if rb.body.Len() > 0 {
		rb.ResponseWriter.Write(rb.body.Bytes())
	}
}

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
			buf := &responseBuffer{ResponseWriter: w}
			next.ServeHTTP(buf, r.WithContext(ctx))

			if err := tx.Commit(r.Context()); err != nil {
				log.Printf("tx commit failed: %v", err)
				apperror.HandleError(w, apperror.Internal("Transaction commit failed", err))
				return
			}
			buf.flush()
		})
	}
}
