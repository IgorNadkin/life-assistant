package main

import (
	"log"

	"life-assistant/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("server started :8080")

	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
