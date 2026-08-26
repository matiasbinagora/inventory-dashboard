## ADDED Requirements

### Requirement: Sensitive local data is excluded
The repository MUST keep MCP credentials, generated indexes, SQLite data, source contents, private URLs, and sensitive project metadata out of tracked artifacts and agent output.

#### Scenario: Repository artifacts are prepared for version control
- **WHEN** files are staged or reviewed for a commit
- **THEN** `.env.mcp`, `.worktrees/`, `graphify-out/`, `.codebase-memory/`, local databases, and generated artifacts MUST be excluded

#### Scenario: Reference repositories are inspected
- **WHEN** project metadata is collected from an approved local path
- **THEN** the process MUST use the explicit read-only allowlist and MUST NOT read secrets, private keys, credential files, or source contents

### Requirement: Context indexes are isolated and read-only
Graphify and Codebase Memory MUST target the active worktree, provide navigation context only, and never be treated as authoritative over source code, tests, or fresh runtime output.

#### Scenario: Implementation context is prepared
- **WHEN** a developer begins broad source investigation
- **THEN** the developer MUST query the active worktree's Graphify and Codebase Memory indexes first and MUST stop if either is missing or stale

#### Scenario: Implementation is validated
- **WHEN** code changes are ready for handoff
- **THEN** the developer MUST refresh Graphify, verify Codebase Memory status, and report both as supporting evidence alongside tests and runtime results

### Requirement: Local services default to local binding
Future dashboard services MUST bind to `127.0.0.1` by default and MUST restrict browser access to the configured local frontend origin.

#### Scenario: Local runtime starts with defaults
- **WHEN** the dashboard is started without an explicit network override
- **THEN** the frontend and backend MUST listen only on loopback interfaces
