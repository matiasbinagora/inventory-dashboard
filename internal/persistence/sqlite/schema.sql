PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS schema_migrations (version INTEGER PRIMARY KEY, applied_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP);
CREATE TABLE IF NOT EXISTS projects (id TEXT PRIMARY KEY, name TEXT NOT NULL CHECK (length(trim(name)) > 0), description TEXT NOT NULL DEFAULT '', agentic_platform TEXT NOT NULL DEFAULT '', github_repository_url TEXT, trello_backlog_url TEXT);
CREATE TABLE IF NOT EXISTS technologies (project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE, name TEXT NOT NULL CHECK (length(trim(name)) > 0), PRIMARY KEY (project_id, name));
CREATE TABLE IF NOT EXISTS public_links (id TEXT PRIMARY KEY, project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE, kind TEXT NOT NULL CHECK (kind IN ('github', 'trello')), url TEXT NOT NULL, label TEXT NOT NULL DEFAULT '', UNIQUE(project_id, kind));
CREATE TABLE IF NOT EXISTS media (id TEXT PRIMARY KEY, project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE, role TEXT NOT NULL CHECK (role IN ('thumbnail', 'original', 'screenshot', 'video')), source TEXT NOT NULL, original_media_id TEXT REFERENCES media(id) ON DELETE RESTRICT, alt_text TEXT NOT NULL DEFAULT '', caption TEXT NOT NULL DEFAULT '', curated INTEGER NOT NULL DEFAULT 0 CHECK (curated IN (0, 1)));
CREATE TABLE IF NOT EXISTS milestones (id TEXT PRIMARY KEY, project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE, date TEXT NOT NULL, title TEXT NOT NULL, description TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS milestone_media (milestone_id TEXT NOT NULL REFERENCES milestones(id) ON DELETE CASCADE, media_id TEXT NOT NULL REFERENCES media(id) ON DELETE CASCADE, PRIMARY KEY (milestone_id, media_id));
CREATE INDEX IF NOT EXISTS idx_media_project ON media(project_id);
CREATE INDEX IF NOT EXISTS idx_milestones_project_date ON milestones(project_id, date);
