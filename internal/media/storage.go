package media

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/matiasbinagora/inventory-dashboard/internal/domain"
)

const DefaultMaxBytes int64 = 50 << 20

var ErrInvalidAsset = errors.New("invalid local media asset")

var contentTypes = map[string]string{
	".gif":  "image/gif",
	".jpeg": "image/jpeg",
	".jpg":  "image/jpeg",
	".mp4":  "video/mp4",
	".mov":  "video/quicktime",
	".png":  "image/png",
	".svg":  "image/svg+xml",
	".webm": "video/webm",
	".webp": "image/webp",
}

type Store struct {
	root     string
	maxBytes int64
}

func NewStore(root string, maxBytes int64) (*Store, error) {
	if strings.TrimSpace(root) == "" {
		return nil, fmt.Errorf("%w: media root is required", ErrInvalidAsset)
	}
	if maxBytes <= 0 {
		maxBytes = DefaultMaxBytes
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid media root", ErrInvalidAsset)
	}
	if err := os.MkdirAll(abs, 0o755); err != nil {
		return nil, fmt.Errorf("create media root: %w", err)
	}
	if resolved, err := filepath.EvalSymlinks(abs); err == nil {
		abs = resolved
	}
	return &Store{root: abs, maxBytes: maxBytes}, nil
}

func (s *Store) Open(source string) (*os.File, string, error) {
	path, err := s.safePath(source)
	if err != nil {
		return nil, "", err
	}
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, "", domain.ErrNotFound
	}
	if err != nil {
		return nil, "", fmt.Errorf("inspect media asset: %w", err)
	}
	if !info.Mode().IsRegular() || info.Size() > s.maxBytes {
		return nil, "", fmt.Errorf("%w: file is not a permitted regular asset", ErrInvalidAsset)
	}
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return nil, "", domain.ErrNotFound
	}
	if !s.insideRoot(resolved) {
		return nil, "", fmt.Errorf("%w: asset leaves media root", ErrInvalidAsset)
	}
	file, err := os.Open(resolved)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, "", domain.ErrNotFound
		}
		return nil, "", fmt.Errorf("open media asset: %w", err)
	}
	return file, contentType(source), nil
}

func (s *Store) Register(source string) error {
	asset := domain.MediaAsset{Curated: true, Role: roleForExtension(source), Source: source}
	if err := asset.Validate(); err != nil {
		return err
	}
	file, _, err := s.Open(source)
	if err != nil {
		return err
	}
	return file.Close()
}

func (s *Store) Save(projectID, filename string, src io.Reader) (string, error) {
	base := filepath.Base(filename)
	ext := strings.ToLower(filepath.Ext(base))
	if base == "." || base == ".." || strings.Contains(base, "\x00") || contentTypes[ext] == "" {
		return "", fmt.Errorf("%w: unsupported filename", ErrInvalidAsset)
	}
	if strings.TrimSpace(projectID) == "" {
		return "", fmt.Errorf("%w: project is required", ErrInvalidAsset)
	}
	relative := filepath.ToSlash(filepath.Join("media", projectID, randomName(base)))
	destination, err := s.safePath(relative)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return "", fmt.Errorf("create media directory: %w", err)
	}
	resolvedDir, err := filepath.EvalSymlinks(filepath.Dir(destination))
	if err != nil || !s.insideRoot(resolvedDir) {
		return "", fmt.Errorf("%w: upload directory leaves media root", ErrInvalidAsset)
	}
	temporary, err := os.CreateTemp(filepath.Dir(destination), ".upload-*")
	if err != nil {
		return "", fmt.Errorf("create media temporary file: %w", err)
	}
	temporaryName := temporary.Name()
	defer os.Remove(temporaryName)
	limited := io.LimitReader(src, s.maxBytes+1)
	count, copyErr := io.Copy(temporary, limited)
	if copyErr != nil {
		temporary.Close()
		return "", fmt.Errorf("write media asset: %w", copyErr)
	}
	if count > s.maxBytes {
		temporary.Close()
		return "", fmt.Errorf("%w: file exceeds configured size limit", ErrInvalidAsset)
	}
	if err := temporary.Close(); err != nil {
		return "", fmt.Errorf("close media temporary file: %w", err)
	}
	if err := os.Rename(temporaryName, destination); err != nil {
		return "", fmt.Errorf("store media asset: %w", err)
	}
	return relative, nil
}

func (s *Store) safePath(source string) (string, error) {
	if source == "" || filepath.IsAbs(source) || strings.Contains(source, "\\") || strings.Contains(source, "\x00") {
		return "", fmt.Errorf("%w: source must be a managed relative path", ErrInvalidAsset)
	}
	clean := filepath.Clean(filepath.FromSlash(source))
	if clean != filepath.FromSlash(source) || clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) || !strings.HasPrefix(clean, "media"+string(filepath.Separator)) {
		return "", fmt.Errorf("%w: unsafe source path", ErrInvalidAsset)
	}
	ext := strings.ToLower(filepath.Ext(clean))
	if contentTypes[ext] == "" {
		return "", fmt.Errorf("%w: unsupported media type", ErrInvalidAsset)
	}
	candidate := filepath.Join(s.root, clean)
	if !s.insideRoot(candidate) {
		return "", fmt.Errorf("%w: source leaves media root", ErrInvalidAsset)
	}
	return candidate, nil
}

func (s *Store) insideRoot(candidate string) bool {
	candidate, err := filepath.Abs(candidate)
	if err != nil {
		return false
	}
	root := filepath.Clean(s.root)
	return candidate == root || strings.HasPrefix(candidate, root+string(filepath.Separator))
}

func contentType(source string) string { return contentTypes[strings.ToLower(filepath.Ext(source))] }

func roleForExtension(source string) domain.MediaRole {
	if strings.HasPrefix(contentType(source), "video/") {
		return domain.Video
	}
	return domain.Original
}

func randomName(filename string) string {
	var token [8]byte
	if _, err := rand.Read(token[:]); err != nil {
		return filename
	}
	return fmt.Sprintf("%x-%s", token, filename)
}

func ValidateContentType(source, declared string) error {
	wanted := contentType(source)
	if wanted == "" {
		return fmt.Errorf("%w: unsupported content type", ErrInvalidAsset)
	}
	if declared != "" {
		parsed, _, err := mime.ParseMediaType(declared)
		if err != nil {
			return fmt.Errorf("%w: unsupported content type", ErrInvalidAsset)
		}
		if parsed != wanted {
			return fmt.Errorf("%w: content type does not match extension", ErrInvalidAsset)
		}
	}
	return nil
}

func DetectContentType(file *os.File) (string, error) {
	var header [512]byte
	n, err := file.Read(header[:])
	if _, seekErr := file.Seek(0, io.SeekStart); err != nil && !errors.Is(err, io.EOF) {
		return "", err
	} else if seekErr != nil {
		return "", seekErr
	}
	return http.DetectContentType(header[:n]), nil
}
