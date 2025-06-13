package middleware

import (
	"context"
	"database/sql"
	"net/http"
)

type contextKey string

const DBContextKey contextKey = "db"

func WithDB(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), DBContextKey, db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
