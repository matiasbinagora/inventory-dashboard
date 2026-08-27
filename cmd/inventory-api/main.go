package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/matiasbinagora/inventory-dashboard/internal/api"
	"github.com/matiasbinagora/inventory-dashboard/internal/application"
	"github.com/matiasbinagora/inventory-dashboard/internal/persistence/sqlite"
)

const defaultAPIAddress = "127.0.0.1:8080"

func main() {
	address := envOr("INVENTORY_API_ADDR", defaultAPIAddress)
	frontendOrigin := envOr("INVENTORY_FRONTEND_ORIGIN", api.DefaultFrontendOrigin)
	dsn := envOr("INVENTORY_DB", "inventory.db")
	store, err := sqlite.Open(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	log.Printf("inventory API listening on %s", address)
	if err := http.ListenAndServe(address, api.NewHandlerWithOrigin(application.NewInventory(store), frontendOrigin)); err != nil {
		log.Fatal(err)
	}
}

func envOr(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
