package main

import (
	"auth-service/config"
)

func main() {
	config.StartConfig().StartServer()
}
