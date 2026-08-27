//go:build integration

package sqlite

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/matiasbinagora/inventory-dashboard/internal/domain"
)

func TestStorePersistsProjectsAndIsolatesChildren(t *testing.T) {
	ctx := context.Background()
	store, err := Open(ctx, "file:store-test?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	one, err := store.CreateProject(ctx, domain.Project{Name: "One", Technologies: []string{"Go"}})
	if err != nil {
		t.Fatal(err)
	}
	two, err := store.CreateProject(ctx, domain.Project{Name: "Two"})
	if err != nil {
		t.Fatal(err)
	}
	media, err := store.AddMedia(ctx, domain.MediaAsset{ProjectID: one.ID, Role: domain.Original, Source: "media/one/original.png"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = store.AddMedia(ctx, domain.MediaAsset{ProjectID: two.ID, Role: domain.Thumbnail, Source: "media/two/thumb.png", OriginalMediaID: media.ID}); err == nil {
		t.Fatal("cross-project media association succeeded")
	}
	if _, err = store.AddMilestone(ctx, domain.Milestone{ProjectID: one.ID, Date: "2026-01-02", Title: "Started", Description: "First milestone", MediaIDs: []string{media.ID}}); err != nil {
		t.Fatal(err)
	}
	loaded, err := store.GetProject(ctx, one.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Media) != 1 || len(loaded.Milestones) != 1 || loaded.Technologies[0] != "Go" {
		t.Fatalf("loaded project lost children: %+v", loaded)
	}
	if err := store.DeleteProject(ctx, one.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := store.GetProject(ctx, one.ID); err != domain.ErrNotFound {
		t.Fatalf("GetProject after delete error=%v", err)
	}
	other, err := store.GetProject(ctx, two.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(other.Media) != 0 || len(other.Milestones) != 0 {
		t.Fatalf("project isolation broken: %+v", other)
	}
}

func TestStoreSurvivesReopen(t *testing.T) {
	ctx := context.Background()
	dsn := filepath.Join(t.TempDir(), "inventory.db")
	store, err := Open(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	created, err := store.CreateProject(ctx, domain.Project{Name: "Persistent project"})
	if err != nil {
		store.Close()
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}

	reopened, err := Open(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	loaded, err := reopened.GetProject(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Name != created.Name {
		t.Fatalf("reopened project name=%q, want %q", loaded.Name, created.Name)
	}
}
