package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/matiasbinagora/inventory-dashboard/internal/application"
	"github.com/matiasbinagora/inventory-dashboard/internal/domain"
)

const DefaultFrontendOrigin = "http://127.0.0.1:3000"

type Handler struct{ inventory *application.Inventory }

func NewHandler(inventory *application.Inventory) http.Handler {
	return NewHandlerWithOrigin(inventory, DefaultFrontendOrigin)
}

func NewHandlerWithOrigin(inventory *application.Inventory, allowedOrigin string) http.Handler {
	h := &Handler{inventory: inventory}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/projects", h.listProjects)
	mux.HandleFunc("POST /api/projects", h.createProject)
	mux.HandleFunc("GET /api/projects/{id}", h.getProject)
	mux.HandleFunc("PUT /api/projects/{id}", h.updateProject)
	mux.HandleFunc("DELETE /api/projects/{id}", h.deleteProject)
	mux.HandleFunc("POST /api/projects/{id}/links", h.addLink)
	mux.HandleFunc("DELETE /api/projects/{id}/links/{linkID}", h.deleteLink)
	mux.HandleFunc("POST /api/projects/{id}/technologies", h.addTechnology)
	mux.HandleFunc("POST /api/projects/{id}/media", h.addMedia)
	mux.HandleFunc("DELETE /api/projects/{id}/media/{mediaID}", h.deleteMedia)
	mux.HandleFunc("POST /api/projects/{id}/milestones", h.addMilestone)
	mux.HandleFunc("PUT /api/projects/{id}/milestones/{milestoneID}", h.updateMilestone)
	mux.HandleFunc("DELETE /api/projects/{id}/milestones/{milestoneID}", h.deleteMilestone)
	return withCORS(withJSON(mux), allowedOrigin)
}

type projectRequest struct {
	Name            string              `json:"name"`
	Description     string              `json:"description"`
	AgenticPlatform string              `json:"agentic_platform"`
	Technologies    []string            `json:"technologies"`
	Links           []domain.PublicLink `json:"links"`
	Media           []domain.MediaAsset `json:"media"`
	Milestones      []domain.Milestone  `json:"milestones"`
}

type technologyRequest struct {
	Name string `json:"name"`
}

func (r projectRequest) project() domain.Project {
	return domain.Project{Name: r.Name, Description: r.Description, AgenticPlatform: r.AgenticPlatform, Technologies: r.Technologies, Links: r.Links, Media: r.Media, Milestones: r.Milestones}
}

func (h *Handler) listProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.inventory.ListProjects(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, projects)
}
func (h *Handler) createProject(w http.ResponseWriter, r *http.Request) {
	var request projectRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	project, err := h.inventory.CreateProject(r.Context(), request.project())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, project)
}
func (h *Handler) getProject(w http.ResponseWriter, r *http.Request) {
	project, err := h.inventory.GetProject(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, project)
}
func (h *Handler) updateProject(w http.ResponseWriter, r *http.Request) {
	var request projectRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	project := request.project()
	project.ID = r.PathValue("id")
	if err := h.inventory.UpdateProject(r.Context(), project); err != nil {
		writeError(w, err)
		return
	}
	h.getProject(w, r)
}
func (h *Handler) deleteProject(w http.ResponseWriter, r *http.Request) {
	if err := h.inventory.DeleteProject(r.Context(), r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *Handler) addLink(w http.ResponseWriter, r *http.Request) {
	var link domain.PublicLink
	if !decodeJSON(w, r, &link) {
		return
	}
	link.ProjectID = r.PathValue("id")
	created, err := h.inventory.AddLink(r.Context(), link)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}
func (h *Handler) deleteLink(w http.ResponseWriter, r *http.Request) {
	if err := h.inventory.DeleteLink(r.Context(), r.PathValue("id"), r.PathValue("linkID")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *Handler) addTechnology(w http.ResponseWriter, r *http.Request) {
	var request technologyRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	if err := h.inventory.AddTechnology(r.Context(), r.PathValue("id"), request.Name); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
func (h *Handler) addMedia(w http.ResponseWriter, r *http.Request) {
	var media domain.MediaAsset
	if !decodeJSON(w, r, &media) {
		return
	}
	media.ProjectID = r.PathValue("id")
	created, err := h.inventory.AddMedia(r.Context(), media)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}
func (h *Handler) deleteMedia(w http.ResponseWriter, r *http.Request) {
	if err := h.inventory.DeleteMedia(r.Context(), r.PathValue("id"), r.PathValue("mediaID")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *Handler) addMilestone(w http.ResponseWriter, r *http.Request) {
	var milestone domain.Milestone
	if !decodeJSON(w, r, &milestone) {
		return
	}
	milestone.ProjectID = r.PathValue("id")
	created, err := h.inventory.AddMilestone(r.Context(), milestone)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}
func (h *Handler) updateMilestone(w http.ResponseWriter, r *http.Request) {
	var milestone domain.Milestone
	if !decodeJSON(w, r, &milestone) {
		return
	}
	milestone.ID, milestone.ProjectID = r.PathValue("milestoneID"), r.PathValue("id")
	if err := h.inventory.UpdateMilestone(r.Context(), milestone); err != nil {
		writeError(w, err)
		return
	}
	h.getProject(w, r)
}
func (h *Handler) deleteMilestone(w http.ResponseWriter, r *http.Request) {
	if err := h.inventory.DeleteMilestone(r.Context(), r.PathValue("id"), r.PathValue("milestoneID")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, domain.ErrInvalidProject), errors.Is(err, domain.ErrInvalidLink), errors.Is(err, domain.ErrInvalidMedia), errors.Is(err, domain.ErrInvalidMilestone):
		status = http.StatusBadRequest
	case errors.Is(err, domain.ErrNotFound):
		status = http.StatusNotFound
	case strings.Contains(err.Error(), "UNIQUE"):
		status = http.StatusConflict
	}
	writeJSON(w, status, errorResponse{Error: err.Error()})
}
func decodeJSON(w http.ResponseWriter, r *http.Request, destination any) bool {
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON body"})
		return false
	}
	return true
}
func withJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func withCORS(next http.Handler, allowedOrigin string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}
		if origin != allowedOrigin {
			http.Error(w, "origin is not allowed", http.StatusForbidden)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Add("Vary", "Origin")
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
