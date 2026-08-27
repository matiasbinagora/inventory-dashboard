# Inventory Dashboard

Local-first dashboard for cataloging current and historical development
projects and their repositories.

## Planned Stack

- Next.js and TypeScript for the web interface.
- Go for the HTTP API and application services.
- SQLite for local persistence.
- OpenSpec for requirements and implementation traceability.
- Trello for backlog and workflow state.

The stack and product behavior are subject to the approved OpenSpec changes.

## Reference Projects

The initial inventory may include metadata from:

- `~/Projects/aws-elemental-inference`
- `~/Projects/slack-bot-video-assistant`
- `~/Projects/acronis-manager`
- `~/Projects/candidates-collaborators-migration`

The importer must remain read-only and must never copy source code or secrets.
Paths should be configurable locally rather than embedded as required runtime
assumptions.

## Agent Workflow

```text
User request
    |
    v
inventory-orchestrator
    |
    +-- brainstorm, clarify, OpenSpec, Trello backlog
    +-- developer handoff
    |
    v
inventory-developer
    |
    +-- dedicated .worktrees/<task> worktree
    +-- Next.js, Go, SQLite, tests
    +-- focused branch and PR
    |
    v
inventory-orchestrator
    |
    +-- technical review gate
    +-- functional acceptance gate
```

The developer must use `.worktrees/` for every task. The folder is intentionally
ignored by Git and is never a location for permanent source changes.

## OpenSpec

Start discovery before defining implementation work:

```text
/opsx-explore <topic>
/opsx-propose <change-name>
/opsx-apply <change-name>
/opsx-validate <change-name>
```

Validate all artifacts with:

```bash
openspec validate --all --strict
```

Every implementation task must reference an OpenSpec change/spec or document
an allowed exception in its Trello card and pull request.

## Local Context Indexes

Graphify creates the repository knowledge graph at `graphify-out/graph.json`.
Codebase Memory maintains the structural symbol index under
`.codebase-memory/`. Both locations are generated and ignored by Git.

Initial local graph setup:

```bash
graphify extract . --code-only
graphify query "Where are the main application entry points?"
```

Refresh Graphify after code changes:

```bash
graphify update .
```

Agents must use narrow Graphify and Codebase Memory queries before broad source
reads. The active worktree must have its own graph and index.

## MCP Access

The project uses separate Docker MCP gateways for GitHub and Trello. Keep the
real `.env.mcp` local; only `.env.mcp.example` is tracked. Current configured
targets are:

- GitHub: `https://github.com/matiasbinagora/inventory-dashboard`
- Trello: `https://trello.com/b/0Seyljto/inventory-dashboard`

Read-only access should be tested before any card, branch, pull request, or
board mutation. Never paste credentials into this repository or chat.

## Repository Hygiene

Ignored local/generated content includes `.env.mcp`, `.worktrees/`,
`graphify-out/`, `.codebase-memory/`, SQLite files, `node_modules/`, and build
outputs.

## Local media runtime

The Go API serves curated local assets at `/media/<managed-relative-path>`;
the configured root contains the managed `media/` directory. By default the
root is the current runtime directory. Configure it with
`INVENTORY_MEDIA_ROOT` and the per-file limit with `INVENTORY_MEDIA_MAX_BYTES`
(default 50 MiB). The API remains loopback-only by default:

```bash
INVENTORY_MEDIA_ROOT=/path/to/runtime INVENTORY_DB=/path/to/inventory.db go run ./cmd/inventory-api
```

Only relative references under `media/` with supported image/video extensions
are accepted. Missing files return 404; remote storage, automatic importing,
and synchronization are intentionally out of scope.
