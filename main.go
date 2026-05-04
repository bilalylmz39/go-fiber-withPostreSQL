package main

import (
	"log"

	productapp "ybilaly/internal/application/product"
	"ybilaly/internal/config"
	"ybilaly/internal/infrastructure/postgres"
	"ybilaly/internal/interfaces/httpapi"
)

func main() {
	cfg := config.Load()

	db, err := postgres.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("database close failed: %v", err)
		}
	}()

	productRepository := postgres.NewProductRepository(db)
	productService := productapp.NewService(productRepository)
	server := httpapi.NewServer(productService, cfg.RequestTimeout)

	log.Fatal(server.App().Listen(":" + cfg.Port))
}
