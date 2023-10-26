package main

import (
	"github.com/realPointer/EnrichInfo/internal/app"
)

// @title EnrichInfo service
// @version 1.0.0
// @description A service that will receive by api the full name, from open api enrich the response with the most probable age, gender and nationality and save the data in the db.

// @host localhost:8080
// @BasePath /v1

// @contact.name Andrew
// @contact.url https://t.me/realPointer

func main() {
	// Run application
	app.Run()
}
