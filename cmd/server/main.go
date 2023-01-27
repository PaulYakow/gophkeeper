package main

import (
	"log"

	"github.com/PaulYakow/gophkeeper/cmd/server/config"
	"github.com/PaulYakow/gophkeeper/internal/server/app"
)

func main() {
	// Configuration
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	srv := app.New(cfg)
	srv.Run()
}
