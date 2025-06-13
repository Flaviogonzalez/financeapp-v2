package routes

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes(db *sql.DB) http.Handler {
	mux := chi.NewRouter()
	// middleware

	// cors

	// handlers

	return mux
}
