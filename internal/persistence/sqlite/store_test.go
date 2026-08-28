//go:build integration

package sqlite

import (
	"context"
	"database/sql"
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

func TestDeleteProjectRemovesCompleteAggregate(t *testing.T) {
	ctx := context.Background()
	store, err := Open(ctx, "file:delete-aggregate?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	project, err := store.CreateProject(ctx, domain.Project{
		Name:         "Delete aggregate",
		Technologies: []string{"Go"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = store.AddLink(ctx, domain.PublicLink{
		ProjectID: project.ID,
		Kind:      domain.GitHub,
		URL:       "https://github.com/example/delete-aggregate",
	}); err != nil {
		t.Fatal(err)
	}
	original, err := store.AddMedia(ctx, domain.MediaAsset{
		ProjectID: project.ID,
		Role:      domain.Original,
		Source:    "media/delete/original.png",
		Curated:   true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = store.AddMedia(ctx, domain.MediaAsset{
		ProjectID:       project.ID,
		Role:            domain.Thumbnail,
		Source:          "media/delete/thumbnail.png",
		OriginalMediaID: original.ID,
		Curated:         true,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err = store.AddMedia(ctx, domain.MediaAsset{
		ProjectID: project.ID,
		Role:      domain.Screenshot,
		Source:    "media/delete/screenshot.png",
		Curated:   true,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err = store.AddMilestone(ctx, domain.Milestone{
		ProjectID:   project.ID,
		Date:        "2026-01-01",
		Title:       "Created",
		Description: "Aggregate deletion regression",
		MediaIDs:    []string{original.ID},
	}); err != nil {
		t.Fatal(err)
	}

	if err := store.DeleteProject(ctx, project.ID); err != nil {
		t.Fatalf("delete complete aggregate: %v", err)
	}
	if _, err := store.GetProject(ctx, project.ID); err != domain.ErrNotFound {
		t.Fatalf("GetProject after delete error=%v", err)
	}
	for _, table := range []string{"technologies", "public_links", "media", "milestones", "milestone_media"} {
		var count int
		if err := store.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+table).Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 0 {
			t.Fatalf("%s contains %d rows after delete", table, count)
		}
	}
}

func TestDeleteProjectFailureDoesNotPartiallyDelete(t *testing.T) {
	ctx := context.Background()
	store, err := Open(ctx, "file:delete-failure?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	project, err := store.CreateProject(ctx, domain.Project{Name: "Protected delete"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = store.db.ExecContext(ctx, `CREATE TRIGGER prevent_project_delete BEFORE DELETE ON projects BEGIN SELECT RAISE(ABORT, 'protected'); END`); err != nil {
		t.Fatal(err)
	}

	if err := store.DeleteProject(ctx, project.ID); err == nil {
		t.Fatal("delete succeeded despite trigger failure")
	}
	if _, err := store.GetProject(ctx, project.ID); err != nil {
		t.Fatalf("project was partially deleted: %v", err)
	}
}

func TestListProjectsReleasesReadCursorBeforeDelete(t *testing.T) {
	ctx := context.Background()
	store, err := Open(ctx, "file:list-delete?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	project, err := store.CreateProject(ctx, domain.Project{Name: "List then delete"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.ListProjects(ctx); err != nil {
		t.Fatal(err)
	}
	if err := store.DeleteProject(ctx, project.ID); err != nil {
		t.Fatalf("delete after list: %v", err)
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

func TestStorePersistsAndReplacesRepositoryReferencesAcrossRestart(t *testing.T) {
	ctx := context.Background()
	dsn := filepath.Join(t.TempDir(), "links.db")
	store, err := Open(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	project, err := store.CreateProject(ctx, domain.Project{Name: "Links", GitHubRepositoryURL: "https://github.com/acme/old", TrelloBacklogURL: "https://trello.com/b/old"})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.UpdateProject(ctx, domain.Project{ID: project.ID, Name: "Links", GitHubRepositoryURL: "https://github.com/acme/new", TrelloBacklogURL: "https://trello.com/b/new"}); err != nil {
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
	loaded, err := reopened.GetProject(ctx, project.ID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.GitHubRepositoryURL != "https://github.com/acme/new" || loaded.TrelloBacklogURL != "https://trello.com/b/new" {
		t.Fatalf("references not replaced after restart: %+v", loaded)
	}
	if err := reopened.UpdateProject(ctx, domain.Project{ID: project.ID, Name: "Links"}); err != nil {
		t.Fatal(err)
	}
	cleared, err := reopened.GetProject(ctx, project.ID)
	if err != nil {
		t.Fatal(err)
	}
	if cleared.GitHubRepositoryURL != "" || cleared.TrelloBacklogURL != "" {
		t.Fatalf("references not cleared: %+v", cleared)
	}
}

func TestOpenMigratesLegacyProjectsTableIdempotently(t *testing.T) {
	ctx := context.Background()
	dsn := filepath.Join(t.TempDir(), "legacy.db")
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, `CREATE TABLE projects (id TEXT PRIMARY KEY, name TEXT NOT NULL, description TEXT NOT NULL DEFAULT '', agentic_platform TEXT NOT NULL DEFAULT ''); INSERT INTO projects(id,name) VALUES('legacy','Legacy project')`); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	store, err := Open(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	loaded, err := store.GetProject(ctx, "legacy")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Name != "Legacy project" || loaded.GitHubRepositoryURL != "" || loaded.TrelloBacklogURL != "" {
		t.Fatalf("legacy project was not preserved: %+v", loaded)
	}
	if _, err := store.db.ExecContext(ctx, `ALTER TABLE projects ADD COLUMN github_repository_url TEXT`); err == nil {
		t.Fatal("migration was not idempotently applied")
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
