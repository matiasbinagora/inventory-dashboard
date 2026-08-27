package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matiasbinagora/inventory-dashboard/internal/application"
	"github.com/matiasbinagora/inventory-dashboard/internal/domain"
	"github.com/matiasbinagora/inventory-dashboard/internal/persistence/sqlite"
)

func TestProjectCRUDAndValidation(t *testing.T) {
	store, err := sqlite.Open(context.Background(), "file:api-test?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	handler := NewHandler(application.NewInventory(store))

	request := func(method, path string, body any) *httptest.ResponseRecorder {
		var payload bytes.Buffer
		if body != nil && json.NewEncoder(&payload).Encode(body) != nil {
			t.Fatal("encode request")
		}
		req := httptest.NewRequest(method, path, &payload)
		req.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, req)
		return response
	}

	invalid := request(http.MethodPost, "/api/projects", map[string]string{"name": "  "})
	if invalid.Code != http.StatusBadRequest {
		t.Fatalf("invalid project status = %d", invalid.Code)
	}

	created := request(http.MethodPost, "/api/projects", map[string]any{"name": "Editorial API", "technologies": []string{" Go "}})
	if created.Code != http.StatusCreated {
		t.Fatalf("create status = %d, body=%s", created.Code, created.Body.String())
	}
	var project domain.Project
	if err := json.NewDecoder(created.Body).Decode(&project); err != nil {
		t.Fatal(err)
	}
	if project.ID == "" || project.Technologies[0] != "Go" {
		t.Fatalf("unexpected project: %+v", project)
	}

	link := request(http.MethodPost, "/api/projects/"+project.ID+"/links", domain.PublicLink{Kind: domain.GitHub, URL: "http://github.com/example/repo"})
	if link.Code != http.StatusBadRequest {
		t.Fatalf("invalid link status = %d", link.Code)
	}
	milestone := request(http.MethodPost, "/api/projects/"+project.ID+"/milestones", domain.Milestone{Date: "2026-01-01", Title: "First", Description: "Curated"})
	if milestone.Code != http.StatusCreated {
		t.Fatalf("milestone status = %d, body=%s", milestone.Code, milestone.Body.String())
	}

	got := request(http.MethodGet, "/api/projects/"+project.ID, nil)
	if got.Code != http.StatusOK {
		t.Fatalf("get status = %d", got.Code)
	}
	var loaded domain.Project
	if err := json.NewDecoder(got.Body).Decode(&loaded); err != nil {
		t.Fatal(err)
	}
	if len(loaded.Milestones) != 1 {
		t.Fatalf("milestones = %+v", loaded.Milestones)
	}
}

func TestCORS(t *testing.T) {
	handler := NewHandler(nil)

	tests := []struct {
		name               string
		method             string
		origin             string
		requestMethod      string
		requestHeaders     string
		wantStatus         int
		wantAllowedOrigin  string
		wantAllowedMethods string
		wantAllowedHeaders string
	}{
		{
			name:              "allows configured frontend origin",
			method:            http.MethodGet,
			origin:            DefaultFrontendOrigin,
			wantStatus:        http.StatusNotFound,
			wantAllowedOrigin: DefaultFrontendOrigin,
		},
		{
			name:       "rejects external origin",
			method:     http.MethodGet,
			origin:     "https://external.example",
			wantStatus: http.StatusForbidden,
		},
		{
			name:               "answers allowed preflight without invoking route",
			method:             http.MethodOptions,
			origin:             DefaultFrontendOrigin,
			requestMethod:      http.MethodPost,
			requestHeaders:     "Content-Type",
			wantStatus:         http.StatusNoContent,
			wantAllowedOrigin:  DefaultFrontendOrigin,
			wantAllowedMethods: "GET, POST, PUT, DELETE, OPTIONS",
			wantAllowedHeaders: "Content-Type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/api/unknown", nil)
			req.Header.Set("Origin", tt.origin)
			if tt.requestMethod != "" {
				req.Header.Set("Access-Control-Request-Method", tt.requestMethod)
			}
			if tt.requestHeaders != "" {
				req.Header.Set("Access-Control-Request-Headers", tt.requestHeaders)
			}
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, req)

			if response.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, tt.wantStatus)
			}
			if got := response.Header().Get("Access-Control-Allow-Origin"); got != tt.wantAllowedOrigin {
				t.Fatalf("allow origin = %q, want %q", got, tt.wantAllowedOrigin)
			}
			if got := response.Header().Get("Access-Control-Allow-Methods"); got != tt.wantAllowedMethods {
				t.Fatalf("allow methods = %q, want %q", got, tt.wantAllowedMethods)
			}
			if got := response.Header().Get("Access-Control-Allow-Headers"); got != tt.wantAllowedHeaders {
				t.Fatalf("allow headers = %q, want %q", got, tt.wantAllowedHeaders)
			}
		})
	}
}

func TestPrivacyBoundariesRejectUncuratedAndSensitiveMedia(t *testing.T) {
	store, err := sqlite.Open(context.Background(), "file:privacy-api-test?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	handler := NewHandler(application.NewInventory(store))
	request := func(body domain.MediaAsset) *httptest.ResponseRecorder {
		var payload bytes.Buffer
		if err := json.NewEncoder(&payload).Encode(body); err != nil {
			t.Fatal(err)
		}
		project, err := store.CreateProject(context.Background(), domain.Project{Name: "Privacy"})
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(http.MethodPost, "/api/projects/"+project.ID+"/media", &payload)
		req.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, req)
		return response
	}

	for _, media := range []domain.MediaAsset{
		{Role: domain.Screenshot, Source: "media/project/screenshot.png"},
		{Role: domain.Screenshot, Source: "media/project/.env", Curated: true},
		{Role: domain.Screenshot, Source: "media/project/source.go", Curated: true},
	} {
		if response := request(media); response.Code != http.StatusBadRequest {
			t.Fatalf("media %+v status = %d, want %d", media, response.Code, http.StatusBadRequest)
		}
	}
}
