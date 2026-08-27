package seed

import (
	"context"
	"testing"

	"github.com/matiasbinagora/inventory-dashboard/internal/domain"
)

type fakeInventory struct {
	projects []domain.Project
}

func (f *fakeInventory) ListProjects(context.Context) ([]domain.Project, error) {
	return append([]domain.Project(nil), f.projects...), nil
}

func (f *fakeInventory) CreateProject(_ context.Context, project domain.Project) (domain.Project, error) {
	project.ID = "generated-" + project.Name
	f.projects = append(f.projects, project)
	return project, nil
}

func TestApplySeedsOnlyCuratedMetadata(t *testing.T) {
	store := &fakeInventory{}

	if err := Apply(context.Background(), store); err != nil {
		t.Fatal(err)
	}
	if len(store.projects) != 2 {
		t.Fatalf("seeded projects = %d, want 2", len(store.projects))
	}

	for _, project := range store.projects {
		if project.ID != "generated-"+project.Name {
			t.Fatalf("seed supplied an identifier: %+v", project)
		}
		if len(project.Links) != 0 || len(project.Media) != 0 || len(project.Milestones) != 0 {
			t.Fatalf("unconfirmed children were seeded: %+v", project)
		}
		if project.Description == "" || len(project.Technologies) == 0 {
			t.Fatalf("incomplete curated project: %+v", project)
		}
	}
}

func TestApplyIsIdempotentByCuratedIdentity(t *testing.T) {
	store := &fakeInventory{}

	if err := Apply(context.Background(), store); err != nil {
		t.Fatal(err)
	}
	if err := Apply(context.Background(), store); err != nil {
		t.Fatal(err)
	}
	if len(store.projects) != 2 {
		t.Fatalf("reapplied seed projects = %d, want 2", len(store.projects))
	}
}

func TestFixturesContainNoSensitiveOrUnconfirmedValues(t *testing.T) {
	for _, project := range Fixtures() {
		if err := project.Validate(); err != nil {
			t.Fatalf("fixture %q is invalid: %v", project.Name, err)
		}
		if project.AgenticPlatform != "" || len(project.Links) != 0 || len(project.Media) != 0 || len(project.Milestones) != 0 {
			t.Fatalf("fixture %q contains unconfirmed fields: %+v", project.Name, project)
		}
	}
}
