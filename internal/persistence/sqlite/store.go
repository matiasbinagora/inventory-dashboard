package sqlite

import (
	"context"
	"crypto/rand"
	"database/sql"
	"embed"
	"fmt"
	"strings"

	"github.com/matiasbinagora/inventory-dashboard/internal/domain"
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schema embed.FS

type Store struct{ db *sql.DB }

func Open(ctx context.Context, dsn string) (*Store, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	store := &Store{db: db}
	if err := store.initialize(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) initialize(ctx context.Context) error {
	contents, err := schema.ReadFile("schema.sql")
	if err != nil {
		return fmt.Errorf("read schema: %w", err)
	}
	if _, err := s.db.ExecContext(ctx, string(contents)); err != nil {
		return fmt.Errorf("initialize schema: %w", err)
	}
	if _, err := s.db.ExecContext(ctx, `INSERT OR IGNORE INTO schema_migrations(version) VALUES(1)`); err != nil {
		return fmt.Errorf("record schema migration: %w", err)
	}
	return nil
}

func (s *Store) CreateProject(ctx context.Context, project domain.Project) (domain.Project, error) {
	if err := project.Validate(); err != nil {
		return domain.Project{}, err
	}
	project.ID = newID()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Project{}, fmt.Errorf("begin project transaction: %w", err)
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `INSERT INTO projects(id,name,description,agentic_platform) VALUES(?,?,?,?)`, project.ID, strings.TrimSpace(project.Name), project.Description, project.AgenticPlatform)
	if err != nil {
		return domain.Project{}, fmt.Errorf("insert project: %w", err)
	}
	for _, technology := range project.Technologies {
		technology = strings.TrimSpace(technology)
		if technology == "" {
			return domain.Project{}, fmt.Errorf("%w: empty technology", domain.ErrInvalidProject)
		}
		if _, err = tx.ExecContext(ctx, `INSERT INTO technologies(project_id,name) VALUES(?,?)`, project.ID, technology); err != nil {
			return domain.Project{}, fmt.Errorf("insert technology: %w", err)
		}
	}
	if err = tx.Commit(); err != nil {
		return domain.Project{}, fmt.Errorf("commit project: %w", err)
	}
	return s.GetProject(ctx, project.ID)
}

func (s *Store) GetProject(ctx context.Context, id string) (domain.Project, error) {
	var p domain.Project
	err := s.db.QueryRowContext(ctx, `SELECT id,name,description,agentic_platform FROM projects WHERE id=?`, id).Scan(&p.ID, &p.Name, &p.Description, &p.AgenticPlatform)
	if err == sql.ErrNoRows {
		return domain.Project{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Project{}, fmt.Errorf("get project: %w", err)
	}
	p.Technologies = []string{}
	p.Links = []domain.PublicLink{}
	p.Media = []domain.MediaAsset{}
	p.Milestones = []domain.Milestone{}
	if err := s.loadChildren(ctx, &p); err != nil {
		return domain.Project{}, err
	}
	return p, nil
}

func (s *Store) ListProjects(ctx context.Context) ([]domain.Project, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id FROM projects ORDER BY name, id`)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	defer rows.Close()
	projects := []domain.Project{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		project, err := s.GetProject(ctx, id)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return projects, nil
}

func (s *Store) UpdateProject(ctx context.Context, project domain.Project) error {
	if project.ID == "" {
		return domain.ErrNotFound
	}
	if err := project.Validate(); err != nil {
		return err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin project update: %w", err)
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `UPDATE projects SET name=?,description=?,agentic_platform=? WHERE id=?`, strings.TrimSpace(project.Name), project.Description, project.AgenticPlatform, project.ID)
	if err != nil {
		return fmt.Errorf("update project: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check project update: %w", err)
	}
	if count == 0 {
		return domain.ErrNotFound
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM technologies WHERE project_id=?`, project.ID); err != nil {
		return err
	}
	for _, technology := range project.Technologies {
		technology = strings.TrimSpace(technology)
		if technology == "" {
			return fmt.Errorf("%w: empty technology", domain.ErrInvalidProject)
		}
		if _, err = tx.ExecContext(ctx, `INSERT INTO technologies(project_id,name) VALUES(?,?)`, project.ID, technology); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) DeleteProject(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM projects WHERE id=?`, id)
	if err != nil {
		return fmt.Errorf("delete project: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check project deletion: %w", err)
	}
	if count == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (s *Store) AddLink(ctx context.Context, link domain.PublicLink) (domain.PublicLink, error) {
	if err := link.Validate(); err != nil {
		return domain.PublicLink{}, err
	}
	if err := s.requireProject(ctx, link.ProjectID); err != nil {
		return domain.PublicLink{}, err
	}
	link.ID = newID()
	_, err := s.db.ExecContext(ctx, `INSERT INTO public_links(id,project_id,kind,url,label) VALUES(?,?,?,?,?)`, link.ID, link.ProjectID, link.Kind, link.URL, link.Label)
	if err != nil {
		return domain.PublicLink{}, fmt.Errorf("insert public link: %w", err)
	}
	return link, nil
}

func (s *Store) DeleteLink(ctx context.Context, projectID, id string) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM public_links WHERE id=? AND project_id=?`, id, projectID)
	if err != nil {
		return fmt.Errorf("delete public link: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (s *Store) AddTechnology(ctx context.Context, projectID, technology string) error {
	technology = strings.TrimSpace(technology)
	if technology == "" {
		return fmt.Errorf("%w: empty technology", domain.ErrInvalidProject)
	}
	if err := s.requireProject(ctx, projectID); err != nil {
		return err
	}
	if _, err := s.db.ExecContext(ctx, `INSERT INTO technologies(project_id,name) VALUES(?,?)`, projectID, technology); err != nil {
		return fmt.Errorf("insert technology: %w", err)
	}
	return nil
}

func (s *Store) AddMedia(ctx context.Context, media domain.MediaAsset) (domain.MediaAsset, error) {
	if err := media.Validate(); err != nil {
		return domain.MediaAsset{}, err
	}
	if err := s.requireProject(ctx, media.ProjectID); err != nil {
		return domain.MediaAsset{}, err
	}
	media.ID = newID()
	if media.OriginalMediaID != "" {
		var role, projectID string
		err := s.db.QueryRowContext(ctx, `SELECT role,project_id FROM media WHERE id=?`, media.OriginalMediaID).Scan(&role, &projectID)
		if err == sql.ErrNoRows || role != string(domain.Original) || projectID != media.ProjectID {
			return domain.MediaAsset{}, fmt.Errorf("%w: original must belong to the same project", domain.ErrInvalidMedia)
		}
		if media.Role != domain.Thumbnail {
			return domain.MediaAsset{}, fmt.Errorf("%w: only thumbnails can reference originals", domain.ErrInvalidMedia)
		}
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO media(id,project_id,role,source,original_media_id,alt_text,caption) VALUES(?,?,?,?,?,?,?)`, media.ID, media.ProjectID, media.Role, media.Source, nullable(media.OriginalMediaID), media.AltText, media.Caption)
	if err != nil {
		return domain.MediaAsset{}, fmt.Errorf("insert media: %w", err)
	}
	return media, nil
}

func (s *Store) AddMilestone(ctx context.Context, milestone domain.Milestone) (domain.Milestone, error) {
	if err := milestone.Validate(); err != nil {
		return domain.Milestone{}, err
	}
	if err := s.requireProject(ctx, milestone.ProjectID); err != nil {
		return domain.Milestone{}, err
	}
	milestone.ID = newID()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Milestone{}, err
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, `INSERT INTO milestones(id,project_id,date,title,description) VALUES(?,?,?,?,?)`, milestone.ID, milestone.ProjectID, milestone.Date, milestone.Title, milestone.Description); err != nil {
		return domain.Milestone{}, fmt.Errorf("insert milestone: %w", err)
	}
	for _, mediaID := range milestone.MediaIDs {
		var projectID string
		if err = tx.QueryRowContext(ctx, `SELECT project_id FROM media WHERE id=?`, mediaID).Scan(&projectID); err != nil || projectID != milestone.ProjectID {
			return domain.Milestone{}, fmt.Errorf("%w: media belongs to another project", domain.ErrInvalidMilestone)
		}
		if _, err = tx.ExecContext(ctx, `INSERT INTO milestone_media(milestone_id,media_id) VALUES(?,?)`, milestone.ID, mediaID); err != nil {
			return domain.Milestone{}, err
		}
	}
	if err = tx.Commit(); err != nil {
		return domain.Milestone{}, err
	}
	return milestone, nil
}

func (s *Store) DeleteMedia(ctx context.Context, projectID, id string) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM media WHERE id=? AND project_id=?`, id, projectID)
	if err != nil {
		return fmt.Errorf("delete media: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (s *Store) UpdateMilestone(ctx context.Context, milestone domain.Milestone) error {
	if err := milestone.Validate(); err != nil {
		return err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `UPDATE milestones SET date=?,title=?,description=? WHERE id=? AND project_id=?`, milestone.Date, milestone.Title, milestone.Description, milestone.ID, milestone.ProjectID)
	if err != nil {
		return fmt.Errorf("update milestone: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return domain.ErrNotFound
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM milestone_media WHERE milestone_id=?`, milestone.ID); err != nil {
		return err
	}
	for _, mediaID := range milestone.MediaIDs {
		var projectID string
		if err = tx.QueryRowContext(ctx, `SELECT project_id FROM media WHERE id=?`, mediaID).Scan(&projectID); err != nil || projectID != milestone.ProjectID {
			return fmt.Errorf("%w: media belongs to another project", domain.ErrInvalidMilestone)
		}
		if _, err = tx.ExecContext(ctx, `INSERT INTO milestone_media(milestone_id,media_id) VALUES(?,?)`, milestone.ID, mediaID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) DeleteMilestone(ctx context.Context, projectID, id string) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM milestones WHERE id=? AND project_id=?`, id, projectID)
	if err != nil {
		return fmt.Errorf("delete milestone: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (s *Store) requireProject(ctx context.Context, id string) error {
	var found string
	err := s.db.QueryRowContext(ctx, `SELECT id FROM projects WHERE id=?`, id).Scan(&found)
	if err == sql.ErrNoRows {
		return domain.ErrNotFound
	}
	return err
}

func (s *Store) loadChildren(ctx context.Context, p *domain.Project) error {
	rows, err := s.db.QueryContext(ctx, `SELECT name FROM technologies WHERE project_id=? ORDER BY name`, p.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return err
		}
		p.Technologies = append(p.Technologies, name)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	linkRows, err := s.db.QueryContext(ctx, `SELECT id,kind,url,label FROM public_links WHERE project_id=? ORDER BY id`, p.ID)
	if err != nil {
		return err
	}
	defer linkRows.Close()
	for linkRows.Next() {
		var link domain.PublicLink
		if err := linkRows.Scan(&link.ID, &link.Kind, &link.URL, &link.Label); err != nil {
			return err
		}
		link.ProjectID = p.ID
		p.Links = append(p.Links, link)
	}
	if err := linkRows.Err(); err != nil {
		return err
	}
	mediaRows, err := s.db.QueryContext(ctx, `SELECT id,role,source,COALESCE(original_media_id,''),alt_text,caption FROM media WHERE project_id=? ORDER BY id`, p.ID)
	if err != nil {
		return err
	}
	defer mediaRows.Close()
	for mediaRows.Next() {
		var media domain.MediaAsset
		if err := mediaRows.Scan(&media.ID, &media.Role, &media.Source, &media.OriginalMediaID, &media.AltText, &media.Caption); err != nil {
			return err
		}
		media.ProjectID = p.ID
		p.Media = append(p.Media, media)
	}
	if err := mediaRows.Err(); err != nil {
		return err
	}
	milestoneRows, err := s.db.QueryContext(ctx, `SELECT id,date,title,description FROM milestones WHERE project_id=? ORDER BY date,id`, p.ID)
	if err != nil {
		return err
	}
	defer milestoneRows.Close()
	for milestoneRows.Next() {
		var milestone domain.Milestone
		if err := milestoneRows.Scan(&milestone.ID, &milestone.Date, &milestone.Title, &milestone.Description); err != nil {
			return err
		}
		milestone.ProjectID = p.ID
		milestone.MediaIDs = []string{}
		mediaIDRows, err := s.db.QueryContext(ctx, `SELECT media_id FROM milestone_media WHERE milestone_id=? ORDER BY media_id`, milestone.ID)
		if err != nil {
			return err
		}
		for mediaIDRows.Next() {
			var mediaID string
			if err := mediaIDRows.Scan(&mediaID); err != nil {
				mediaIDRows.Close()
				return err
			}
			milestone.MediaIDs = append(milestone.MediaIDs, mediaID)
		}
		if err := mediaIDRows.Err(); err != nil {
			mediaIDRows.Close()
			return err
		}
		mediaIDRows.Close()
		p.Milestones = append(p.Milestones, milestone)
	}
	return milestoneRows.Err()
}

func nullable(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func newID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", b)
}
