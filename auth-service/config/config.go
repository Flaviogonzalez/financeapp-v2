package config

import (
	data "auth-service/models"
	"auth-service/routes"
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	db   *sql.DB
	data data.Models
}

func StartConfig() *Config {
	db := connectToDB()

	return &Config{
		db:   db,
		data: data.New(db), // initialize the data models
	}
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DATABASE_URL")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Panic("Possibly incorrect DSN: ", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Could not connect to the database: ", err)
	}

	return db
}

func (app *Config) StartServer() {
	srv := &http.Server{
		Addr:    ":80",
		Handler: routes.Routes(app.db),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Panic("jiji se me cayó el server")
	}
}
