package config

import (
	data "auth-service/models"
	"auth-service/routes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

const (
	port = 80
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
		Addr:    fmt.Sprintf(":%v", port),
		Handler: routes.Routes(app.db),
	}

	log.Println("Starting auth service on http://localhost:80")
	if err := srv.ListenAndServe(); err != nil {
		log.Panic("jiji se me cayó el server")
	}
}
