package routes

import (
	"auth-service/handlers"
	"auth-service/middleware"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Routes(db *sql.DB) http.Handler {
	mux := chi.NewRouter()
	// middleware
	mux.Use(chi_middleware.Logger)
	mux.Use(chi_middleware.Recoverer)
	mux.Use(middleware.WithDB(db)) // Add database to context

	// cors
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	// handlers
	mux.Post("/login", handlers.Login)
	mux.Post("/register", handlers.Register)
	return mux
}
