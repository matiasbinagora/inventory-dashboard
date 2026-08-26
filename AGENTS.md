# Inventory Dashboard Agent Instructions

## Project Purpose

This repository will contain a local-first inventory dashboard for projects and
repositories developed in the local Projects workspace. The initial product
target is a local Next.js website, a Go backend, and SQLite persistence.

The initial source examples are:

- `~/Projects/aws-elemental-inference`
- `~/Projects/slack-bot-video-assistant`
- `~/Projects/acronis-manager`
- `~/Projects/candidates-collaborators-migration`

These paths are read-only reference sources. The application must store project
metadata, not source code, credentials, transcripts, private URLs, or generated
artifacts.

## Source Of Truth

- Read this file before acting.
- OpenSpec is the source of truth for requirements, behavior, API contracts,
  data contracts, design decisions, and task traceability.
- Trello is the source of truth for task status, ownership, predecessors,
  acceptance criteria, and transition evidence.
- Source code, automated tests, and fresh runtime output are authoritative over
  memory indexes and generated graphs.
- Do not implement a Trello task without an OpenSpec Change, OpenSpec Spec, or
  documented OpenSpec Exception.

## Agent Roster

### `inventory-orchestrator`

Owns brainstorming, product clarification, OpenSpec governance, Trello backlog
quality, developer handoffs, and verification of delivered work against the
acceptance criteria.

The orchestrator:

- starts discovery with `/opsx-explore`;
- uses `/opsx-propose` for new behavior or contract changes;
- may edit governance files, OpenSpec files, and documentation;
- must not implement production frontend, backend, database, or test code;
- must not silently change scope or acceptance criteria;
- must not merge pull requests or move work to Done without human approval;
- performs separate technical and functional verification gates.

### `inventory-developer`

Owns approved implementation work for the Next.js frontend, Go backend, SQLite
database, tests, local runtime, and pull request preparation.

The developer:

- works only from a complete Trello handoff and OpenSpec reference;
- starts implementation with `/opsx-apply`;
- must use a dedicated worktree under `.worktrees/` for every task;
- must never edit implementation files in the repository root after a task
  worktree exists;
- must not create or rewrite backlog scope;
- must not skip required tests, OpenSpec validation, or context evidence;
- must not merge pull requests or move tasks to Done.

## Mandatory Worktree Policy

`.worktrees/` is the only approved location for task worktrees. It is ignored
by Git and must remain local.

Use this naming convention:

```text
.worktrees/inventory-dashboard-task-<number>-<kebab-case-name>
```

Use this branch convention:

```text
feature/task-<number>-<kebab-case-name>
```

Before editing implementation files, the developer must verify the active
worktree and branch with `git worktree list --porcelain` and report both in the
handoff. Each task gets one worktree and one focused branch. Do not delete or
reuse a worktree for another task without explicit approval.

## OpenSpec Workflow

OpenSpec uses the `spec-driven` schema. Required commands are:

1. `/opsx-explore <topic>` for discovery and brainstorming.
2. `/opsx-propose <change-name>` for a new formal change.
3. `/opsx-continue <change-name>` when artifacts are incomplete.
4. `/opsx-apply <change-name>` for implementation.
5. `/opsx-validate <change-name>` for acceptance validation.
6. `/opsx-sync <change-name>` only after review gates pass.
7. `/opsx-archive <change-name>` only after release readiness and approval.

Every behavior, API, database, UI, security, or workflow change requires a
new change or an explicit existing spec reference. Specs must use RFC 2119
keywords and verifiable Given/When/Then scenarios. Validate with:

```bash
openspec validate --all --strict
```

## Trello Workflow

Every implementation card must use `<TYPE>-<ID> <Title>` and include:

- ID, Type, Title, Description;
- Requirements;
- numbered Acceptance Criteria;
- Predecessors;
- OpenSpec Change, Spec, or Exception;
- Suggested agent;
- Validation expectations;
- Unit/integration test expectations;
- Risks and non-scope.

The workflow is:

```text
backlog -> ready -> in progress -> code review
        -> functional review -> ready to release -> done
```

Use `blocked` for missing information or failed validation. Every transition
requires a factual Trello comment. The developer may transition only `ready ->
in progress` and `in progress -> code review`. The orchestrator owns backlog
readiness and final closure, but must not bypass either review gate.

## Graph-First Context Policy

Graphify and Codebase Memory are mandatory, read-only context sources for
implementation and review work. They reduce broad reads but do not replace
tests, OpenSpec, runtime validation, or the handoff.

Before broad source search, after reading the task and OpenSpec artifacts:

1. Query Graphify with a narrow question using depth 1-2 and a bounded budget.
2. Use `graphify path` or `graphify explain` for focused relationships.
3. Call Codebase Memory `list_projects` and use the exact project identifier.
4. Use `search_graph`, `trace_path`, and `get_code_snippet` for exact symbols.
5. Open only the files selected by those queries unless more context is
   justified and reported.

The active worktree must have its own graph and index. Never point a worktree
at another checkout's `graphify-out` or Codebase Memory cache. If either is
missing or stale, stop and refresh the active worktree before implementation.

After code changes, refresh the local graph with:

```bash
graphify update .
```

Do not use `graphify global` or merge graphs from private repositories. Do not
put credentials, private URLs, source contents, or sensitive metadata in graph
questions, reports, or handoffs.

## Architecture Rules

Backend dependency direction:

```text
HTTP API -> Application -> Domain
HTTP API -> Infrastructure
Infrastructure -> Application / Domain
Domain -> no framework or infrastructure dependencies
```

Frontend direction:

```text
Routes -> Features / UI / API client
Features -> UI / API client / domain types
UI -> shared design tokens and utilities
```

The Next.js application must call the Go API rather than duplicating backend
business logic in Next.js route handlers. Keep the application local-only and
bind services to `127.0.0.1` by default.

## Security And Privacy

- Never print, commit, or place secrets in prompts, docs, logs, Trello, or PRs.
- Keep `.env.mcp` ignored and use `.env.mcp.example` only for placeholders.
- Treat all configured project paths as an explicit read-only allowlist.
- Do not follow symlinks outside approved roots during repository scans.
- Do not read `.env`, private keys, credential files, or secret-bearing Git
  configuration while importing metadata.
- Do not expose source code through the dashboard.
- Restrict CORS to the local frontend origin and bind local services to
  `127.0.0.1`.

## Validation And Definition Of Done

Completion reports must include the task ID, branch, worktree, changed files,
OpenSpec reference, tests, commands and results, Graphify evidence, Codebase
Memory evidence, PR URL, and remaining risks.

A task is complete only when its scope is respected, required tests pass,
OpenSpec validation passes, runtime acceptance criteria are verified, generated
artifacts remain ignored, and no blocker remains.
