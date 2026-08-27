package domain

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"path"
	"strings"
	"time"
)

var (
	ErrNotFound         = errors.New("record not found")
	ErrInvalidProject   = errors.New("invalid project")
	ErrInvalidLink      = errors.New("invalid public link")
	ErrInvalidMedia     = errors.New("invalid media asset")
	ErrInvalidMilestone = errors.New("invalid milestone")
)

type Project struct {
	ID              string       `json:"id"`
	Name            string       `json:"name"`
	Description     string       `json:"description,omitempty"`
	AgenticPlatform string       `json:"agentic_platform,omitempty"`
	Technologies    []string     `json:"technologies"`
	Links           []PublicLink `json:"links"`
	Media           []MediaAsset `json:"media"`
	Milestones      []Milestone  `json:"milestones"`
}

type LinkKind string

const (
	GitHub LinkKind = "github"
	Trello LinkKind = "trello"
)

type PublicLink struct {
	ID        string   `json:"id,omitempty"`
	ProjectID string   `json:"project_id,omitempty"`
	Kind      LinkKind `json:"kind"`
	URL       string   `json:"url"`
	Label     string   `json:"label,omitempty"`
}

type MediaRole string

const (
	Thumbnail  MediaRole = "thumbnail"
	Original   MediaRole = "original"
	Screenshot MediaRole = "screenshot"
	Video      MediaRole = "video"
)

type MediaAsset struct {
	ID              string    `json:"id,omitempty"`
	ProjectID       string    `json:"project_id,omitempty"`
	Role            MediaRole `json:"role"`
	Source          string    `json:"source"`
	OriginalMediaID string    `json:"original_media_id,omitempty"`
	AltText         string    `json:"alt_text,omitempty"`
	Caption         string    `json:"caption,omitempty"`
	Curated         bool      `json:"curated"`
}

type Milestone struct {
	ID          string   `json:"id,omitempty"`
	ProjectID   string   `json:"project_id,omitempty"`
	Date        string   `json:"date"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	MediaIDs    []string `json:"media_ids"`
}

func (p Project) Validate() error {
	if strings.TrimSpace(p.Name) == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidProject)
	}
	if err := validateCuratedText(p.Description, "description"); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidProject, err)
	}
	if err := validateCuratedText(p.AgenticPlatform, "agentic platform"); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidProject, err)
	}
	return nil
}

func (l PublicLink) Validate() error {
	if l.Kind != GitHub && l.Kind != Trello {
		return fmt.Errorf("%w: unsupported kind", ErrInvalidLink)
	}
	u, err := url.Parse(l.URL)
	if err != nil || u.Scheme != "https" || u.User != nil || u.Hostname() == "" || u.Path == "" || u.Path == "/" || u.RawQuery != "" || u.Fragment != "" {
		return fmt.Errorf("%w: URL must be an HTTPS URL with a path", ErrInvalidLink)
	}
	host := strings.ToLower(u.Hostname())
	allowed := map[LinkKind]string{GitHub: "github.com", Trello: "trello.com"}
	if host != allowed[l.Kind] && host != "www."+allowed[l.Kind] {
		return fmt.Errorf("%w: unsupported host", ErrInvalidLink)
	}
	if privateHost(host) {
		return fmt.Errorf("%w: private host", ErrInvalidLink)
	}
	return nil
}

func (m MediaAsset) Validate() error {
	if !m.Curated {
		return fmt.Errorf("%w: media must be explicitly curated", ErrInvalidMedia)
	}
	validRole := m.Role == Thumbnail || m.Role == Original || m.Role == Screenshot || m.Role == Video
	if !validRole || strings.TrimSpace(m.Source) == "" {
		return fmt.Errorf("%w: role and source are required", ErrInvalidMedia)
	}
	if m.Role == Video && strings.HasPrefix(strings.ToLower(m.Source), "https://") {
		if err := validatePublicHTTPS(m.Source); err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidMedia, err)
		}
		return nil
	}
	if m.Role == Video && strings.Contains(m.Source, "://") {
		return fmt.Errorf("%w: video URL must use HTTPS", ErrInvalidMedia)
	}
	clean := path.Clean(m.Source)
	if path.IsAbs(m.Source) || clean != m.Source || clean == "." || clean == ".." || strings.HasPrefix(clean, "../") || !strings.HasPrefix(clean, "media/") || !safeMediaPath(clean) {
		return fmt.Errorf("%w: local source must be a managed relative media path", ErrInvalidMedia)
	}
	if err := validateCuratedText(m.AltText, "alt text"); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidMedia, err)
	}
	if err := validateCuratedText(m.Caption, "caption"); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidMedia, err)
	}
	return nil
}

func (m Milestone) Validate() error {
	if _, err := time.Parse("2006-01-02", m.Date); err != nil {
		return fmt.Errorf("%w: date must be YYYY-MM-DD", ErrInvalidMilestone)
	}
	if strings.TrimSpace(m.Title) == "" || strings.TrimSpace(m.Description) == "" {
		return fmt.Errorf("%w: title and description are required", ErrInvalidMilestone)
	}
	return nil
}

func validatePublicHTTPS(raw string) error {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.User != nil || u.Hostname() == "" || u.Path == "" || u.Path == "/" || u.RawQuery != "" || u.Fragment != "" {
		return errors.New("URL must be public HTTPS")
	}
	if privateHost(strings.ToLower(u.Hostname())) {
		return errors.New("URL host is private or local")
	}
	return nil
}

func validateCuratedText(value, field string) error {
	lower := strings.ToLower(value)
	for _, marker := range []string{"transcript", "transcripción", "password=", "passwd=", "api_key", "apikey", "secret=", "token=", "authorization:", "package main", "-----begin"} {
		if strings.Contains(lower, marker) {
			return fmt.Errorf("%s contains excluded private content", field)
		}
	}
	return nil
}

func safeMediaPath(source string) bool {
	parts := strings.Split(source, "/")
	if len(parts) < 2 || parts[0] != "media" {
		return false
	}
	for _, part := range parts[1:] {
		lower := strings.ToLower(part)
		if part == "" || strings.HasPrefix(part, ".") || lower == "node_modules" || lower == ".git" || lower == "transcripts" || lower == "artifacts" {
			return false
		}
	}
	extension := strings.ToLower(path.Ext(parts[len(parts)-1]))
	for _, allowed := range []string{".png", ".jpg", ".jpeg", ".webp", ".gif", ".svg", ".mp4", ".webm", ".mov"} {
		if extension == allowed {
			return true
		}
	}
	return false
}

func privateHost(host string) bool {
	if host == "localhost" || strings.HasSuffix(host, ".localhost") || strings.HasSuffix(host, ".local") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && (ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsUnspecified())
}
