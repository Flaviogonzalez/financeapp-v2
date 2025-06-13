package main

import (
	"auth-service/config"
	"log"
)

func main() {
	log.Println("Starting auth service on localhost:80")
	config.StartConfig().StartServer()
}
