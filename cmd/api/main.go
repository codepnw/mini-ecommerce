package main

import (
	"log"

	"github.com/codepnw/mini-ecommerce/pkg/config"
	"github.com/codepnw/mini-ecommerce/pkg/routes"
)

const envPath = "dev.env"

func main() {
	cfg, err := config.LoadConfig(envPath)
	if err != nil {
		log.Fatal(err)
	}

	if err = routes.RegisterRoutes(cfg); err != nil {
		log.Fatal(err)
	}
}
