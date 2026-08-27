package seed

import (
	"context"
	"fmt"

	"github.com/matiasbinagora/inventory-dashboard/internal/domain"
)

type inventory interface {
	CreateProject(context.Context, domain.Project) (domain.Project, error)
	ListProjects(context.Context) ([]domain.Project, error)
}

// Fixtures returns the manually curated, privacy-safe catalog entries. It
// intentionally contains no repository paths, links, lifecycle fields, dates,
// owners, metrics, or media because those values were not confirmed.
func Fixtures() []domain.Project {
	return []domain.Project{
		{
			Name:        "Slack Video Assistant",
			Description: "A local-first Slack bot for understanding and exporting short videos.",
			Technologies: []string{
				"Python",
				"Slack Bolt",
				"Socket Mode",
				"Claude Agent SDK",
				"FFmpeg",
				"FFprobe",
			},
		},
		{
			Name:        "AWS Elemental Inference Smart Crop Demo",
			Description: "A proof of concept for converting landscape video into portrait video with AWS Elemental MediaConvert Smart Cropping using AWS Elemental Inference.",
			Technologies: []string{
				"AWS Elemental MediaConvert",
				"AWS Elemental Inference",
				"Amazon S3",
				"AWS CloudFormation",
				"H.264",
				"AAC",
				"FFmpeg",
				"FFprobe",
			},
		},
	}
}

// Apply adds missing fixtures by their curated display identity. Existing
// projects are left untouched so administration remains the owner of edits.
func Apply(ctx context.Context, target inventory) error {
	projects, err := target.ListProjects(ctx)
	if err != nil {
		return fmt.Errorf("list projects before seed: %w", err)
	}
	existing := make(map[string]struct{}, len(projects))
	for _, project := range projects {
		existing[project.Name] = struct{}{}
	}

	for _, fixture := range Fixtures() {
		if _, ok := existing[fixture.Name]; ok {
			continue
		}
		if _, err := target.CreateProject(ctx, fixture); err != nil {
			return fmt.Errorf("seed project %q: %w", fixture.Name, err)
		}
	}
	return nil
}
