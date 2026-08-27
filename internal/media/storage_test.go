package media

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/matiasbinagora/inventory-dashboard/internal/domain"
)

func TestStoreOpenRejectsUnsafeAndUnsupportedSources(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "media", "project"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "media", "project", "image.png"), []byte("image"), 0o644); err != nil {
		t.Fatal(err)
	}
	store, err := NewStore(root, 100)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		source string
		want   error
	}{
		{name: "absolute", source: filepath.Join(root, "media", "project", "image.png"), want: ErrInvalidAsset},
		{name: "traversal", source: "media/project/../.env", want: ErrInvalidAsset},
		{name: "unsupported extension", source: "media/project/source.go", want: ErrInvalidAsset},
		{name: "missing", source: "media/project/missing.png", want: domain.ErrNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, _, gotErr := store.Open(tt.source)
			if file != nil {
				file.Close()
			}
			if !errors.Is(gotErr, tt.want) {
				t.Fatalf("Open(%q) error = %v, want %v", tt.source, gotErr, tt.want)
			}
		})
	}
}

func TestStoreSaveUsesManagedRelativeReferenceAndSizeLimit(t *testing.T) {
	store, err := NewStore(t.TempDir(), 4)
	if err != nil {
		t.Fatal(err)
	}
	reference, err := store.Save("project", "original.png", bytes.NewReader([]byte("1234")))
	if err != nil {
		t.Fatal(err)
	}
	if reference == "" || filepath.IsAbs(reference) || reference[:len("media/")] != "media/" {
		t.Fatalf("unexpected managed reference %q", reference)
	}
	if _, err := store.Save("project", "too-large.png", bytes.NewReader([]byte("12345"))); !errors.Is(err, ErrInvalidAsset) {
		t.Fatalf("oversized save error = %v", err)
	}
}

func TestStoreRejectsEscapingSymlink(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "secret.png"), []byte("private"), 0o644); err != nil {
		t.Fatal(err)
	}
	managed := filepath.Join(root, "media", "project")
	if err := os.MkdirAll(managed, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(outside, "secret.png"), filepath.Join(managed, "linked.png")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	store, err := NewStore(root, DefaultMaxBytes)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := store.Open("media/project/linked.png"); !errors.Is(err, ErrInvalidAsset) {
		t.Fatalf("escaping symlink error = %v", err)
	}
}
