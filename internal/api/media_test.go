package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matiasbinagora/inventory-dashboard/internal/application"
	"github.com/matiasbinagora/inventory-dashboard/internal/domain"
	"github.com/matiasbinagora/inventory-dashboard/internal/persistence/sqlite"
)

func TestManagedMediaServingAndSQLiteRestart(t *testing.T) {
	root := t.TempDir()
	assetDir := filepath.Join(root, "media", "atlas")
	if err := os.MkdirAll(assetDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(assetDir, "original.png"), []byte("real temporary image"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(assetDir, "demo.webm"), []byte{0x1a, 0x45, 0xdf, 0xa3}, 0o644); err != nil {
		t.Fatal(err)
	}
	database := filepath.Join(t.TempDir(), "inventory.db")
	ctx := context.Background()
	store, err := sqlite.Open(ctx, database)
	if err != nil {
		t.Fatal(err)
	}
	project, err := store.CreateProject(ctx, domain.Project{Name: "Atlas"})
	if err != nil {
		t.Fatal(err)
	}
	handler := NewHandlerWithOriginAndMediaRoot(application.NewInventory(store), DefaultFrontendOrigin, root, 1024)
	add := func(asset domain.MediaAsset) {
		record := httptest.NewRecorder()
		body, _ := json.Marshal(asset)
		req := httptest.NewRequest(http.MethodPost, "/api/projects/"+project.ID+"/media", strings.NewReader(string(body)))
		handler.ServeHTTP(record, req)
		if record.Code != http.StatusCreated {
			t.Fatalf("register %q status = %d, body = %s", asset.Source, record.Code, record.Body.String())
		}
	}
	add(domain.MediaAsset{Role: domain.Original, Source: "media/atlas/original.png", Curated: true})
	add(domain.MediaAsset{Role: domain.Video, Source: "media/atlas/demo.webm", Curated: true})
	store.Close()

	store, err = sqlite.Open(ctx, database)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	handler = NewHandlerWithOriginAndMediaRoot(application.NewInventory(store), DefaultFrontendOrigin, root, 1024)
	loaded, err := store.GetProject(ctx, project.ID)
	if err != nil || len(loaded.Media) != 2 {
		t.Fatalf("reloaded media = %+v, error = %v", loaded.Media, err)
	}
	for _, test := range []struct {
		path        string
		contentType string
	}{
		{path: "/media/atlas/original.png", contentType: "image/png"},
		{path: "/media/atlas/demo.webm", contentType: "video/webm"},
	} {
		record := httptest.NewRecorder()
		handler.ServeHTTP(record, httptest.NewRequest(http.MethodGet, test.path, nil))
		if record.Code != http.StatusOK || record.Header().Get("Content-Type") != test.contentType {
			t.Errorf("serve %s = status %d type %q", test.path, record.Code, record.Header().Get("Content-Type"))
		}
	}
	missing := httptest.NewRecorder()
	handler.ServeHTTP(missing, httptest.NewRequest(http.MethodGet, "/media/atlas/missing.png", nil))
	if missing.Code != http.StatusNotFound || strings.Contains(missing.Body.String(), root) {
		t.Fatalf("missing response = status %d body %q", missing.Code, missing.Body.String())
	}
}
