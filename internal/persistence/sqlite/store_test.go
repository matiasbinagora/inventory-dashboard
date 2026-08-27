//go:build integration

package sqlite

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/matiasbinagora/inventory-dashboard/internal/application"
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
	media, err := store.AddMedia(ctx, domain.MediaAsset{ProjectID: one.ID, Role: domain.Original, Source: "media/one/original.png", Curated: true})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = store.AddMedia(ctx, domain.MediaAsset{ProjectID: two.ID, Role: domain.Thumbnail, Source: "media/two/thumb.png", OriginalMediaID: media.ID, Curated: true}); err == nil {
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

func TestCreateProjectWithChildrenIsAtomic(t *testing.T) {
	tests := []struct {
		name       string
		trigger    string
		project    domain.Project
		childTable string
	}{
		{
			name:       "link failure rolls back project and earlier links",
			trigger:    `CREATE TRIGGER fail_link BEFORE INSERT ON public_links WHEN NEW.url = 'fail' BEGIN SELECT RAISE(ABORT, 'link failure'); END`,
			project:    domain.Project{Name: "Atomic links", Links: []domain.PublicLink{{Kind: domain.GitHub, URL: "https://github.com/example/repo"}, {Kind: domain.Trello, URL: "fail"}}},
			childTable: "public_links",
		},
		{
			name:       "media failure rolls back project and earlier media",
			trigger:    `CREATE TRIGGER fail_media BEFORE INSERT ON media WHEN NEW.source = 'fail' BEGIN SELECT RAISE(ABORT, 'media failure'); END`,
			project:    domain.Project{Name: "Atomic media", Media: []domain.MediaAsset{{Role: domain.Original, Source: "media/original.png", Curated: true}, {Role: domain.Screenshot, Source: "fail", Curated: true}}},
			childTable: "media",
		},
		{
			name:       "milestone failure rolls back project and earlier milestones",
			trigger:    `CREATE TRIGGER fail_milestone BEFORE INSERT ON milestones WHEN NEW.title = 'fail' BEGIN SELECT RAISE(ABORT, 'milestone failure'); END`,
			project:    domain.Project{Name: "Atomic milestones", Milestones: []domain.Milestone{{Date: "2026-01-01", Title: "Started", Description: "First"}, {Date: "2026-01-02", Title: "fail", Description: "Second"}}},
			childTable: "milestones",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			store, err := Open(ctx, fmt.Sprintf("file:atomic-%s?mode=memory&cache=shared", tt.childTable))
			if err != nil {
				t.Fatal(err)
			}
			defer store.Close()
			if _, err := store.db.ExecContext(ctx, tt.trigger); err != nil {
				t.Fatal(err)
			}

			_, err = application.NewInventory(store).CreateProject(ctx, tt.project)
			if err == nil {
				t.Fatal("CreateProject succeeded despite child write failure")
			}
			for _, table := range []string{"projects", "technologies", "public_links", "media", "milestones", "milestone_media"} {
				var count int
				if err := store.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+table).Scan(&count); err != nil {
					t.Fatal(err)
				}
				if count != 0 {
					t.Fatalf("%s contains %d rows after rollback", table, count)
				}
			}
		})
	}
}

func TestCreateProjectWithChildrenPersistsCompleteAggregate(t *testing.T) {
	ctx := context.Background()
	store, err := Open(ctx, "file:atomic-success?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	created, err := application.NewInventory(store).CreateProject(ctx, domain.Project{
		Name:         "Complete aggregate",
		Technologies: []string{"Go"},
		Links:        []domain.PublicLink{{Kind: domain.GitHub, URL: "https://github.com/example/repo"}},
		Media:        []domain.MediaAsset{{Role: domain.Original, Source: "media/original.png", Curated: true}},
		Milestones:   []domain.Milestone{{Date: "2026-01-01", Title: "Started", Description: "First"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(created.Technologies) != 1 || len(created.Links) != 1 || len(created.Media) != 1 || len(created.Milestones) != 1 {
		t.Fatalf("created aggregate missing children: %+v", created)
	}
}
