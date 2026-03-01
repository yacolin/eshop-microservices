package main

import (
	"log"
	"os"

	"eshop-microservices/internal/inventory-service/app"
)

func main() {
	app, err := app.New(os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Fatalf("app: %v", err)
	}
	if err := app.Run(); err != nil {
		log.Fatalf("run: %v", err)
	}
}
