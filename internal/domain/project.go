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
	return nil
}

func (l PublicLink) Validate() error {
	if l.Kind != GitHub && l.Kind != Trello {
		return fmt.Errorf("%w: unsupported kind", ErrInvalidLink)
	}
	u, err := url.Parse(l.URL)
	if err != nil || u.Scheme != "https" || u.User != nil || u.Hostname() == "" || u.Path == "" || u.Path == "/" {
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
	if path.IsAbs(m.Source) || clean == "." || clean == ".." || strings.HasPrefix(clean, "../") || strings.HasPrefix(clean, "media/../") || !strings.HasPrefix(clean, "media/") {
		return fmt.Errorf("%w: local source must be a managed relative media path", ErrInvalidMedia)
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
	if err != nil || u.Scheme != "https" || u.User != nil || u.Hostname() == "" {
		return errors.New("URL must be public HTTPS")
	}
	if privateHost(strings.ToLower(u.Hostname())) {
		return errors.New("URL host is private or local")
	}
	return nil
}

func privateHost(host string) bool {
	if host == "localhost" || strings.HasSuffix(host, ".localhost") || strings.HasSuffix(host, ".local") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && (ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsUnspecified())
}
