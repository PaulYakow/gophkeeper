package main

import (
	"log"

	"github.com/PaulYakow/gophkeeper/cmd/client/config"
	"github.com/PaulYakow/gophkeeper/internal/client/app"
)

func main() {
	// Configuration
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("client config error: %s", err)
	}

	client := app.New(cfg)
	client.Run()
}
