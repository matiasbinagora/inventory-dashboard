package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/matiasbinagora/inventory-dashboard/internal/api"
	"github.com/matiasbinagora/inventory-dashboard/internal/application"
	"github.com/matiasbinagora/inventory-dashboard/internal/media"
	"github.com/matiasbinagora/inventory-dashboard/internal/persistence/sqlite"
	"github.com/matiasbinagora/inventory-dashboard/internal/seed"
)

const defaultAPIAddress = "127.0.0.1:8080"

func main() {
	address := envOr("INVENTORY_API_ADDR", defaultAPIAddress)
	frontendOrigin := envOr("INVENTORY_FRONTEND_ORIGIN", api.DefaultFrontendOrigin)
	dsn := envOr("INVENTORY_DB", "inventory.db")
	mediaRoot := envOr("INVENTORY_MEDIA_ROOT", ".")
	maxMediaBytes := envInt64Or("INVENTORY_MEDIA_MAX_BYTES", media.DefaultMaxBytes)
	store, err := sqlite.Open(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	if err := seed.Apply(context.Background(), application.NewInventory(store)); err != nil {
		log.Fatal(err)
	}
	log.Printf("inventory API listening on %s", address)
	if err := http.ListenAndServe(address, api.NewHandlerWithOriginAndMediaRoot(application.NewInventory(store), frontendOrigin, mediaRoot, maxMediaBytes)); err != nil {
		log.Fatal(err)
	}
}

func envInt64Or(name string, fallback int64) int64 {
	value, err := strconv.ParseInt(os.Getenv(name), 10, 64)
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}

func envOr(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
