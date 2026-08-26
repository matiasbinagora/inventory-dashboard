## Context

The repository is a local-first project that will eventually combine a Next.js frontend, a Go API, and SQLite. The bootstrap currently contains agent definitions, OpenSpec configuration, local MCP gateways, and navigation indexes, but no product runtime. The main stakeholders are the orchestrator, the implementation developer, and the repository owner.

## Goals / Non-Goals

**Goals:**

- Make ownership and permissions explicit between planning and implementation.
- Ensure implementation changes are isolated from the root checkout.
- Keep local project metadata and generated context artifacts out of tracked files and agent output.
- Preserve a clear path from OpenSpec requirements to Trello work and validation evidence.

**Non-Goals:**

- Building the inventory dashboard UI or API.
- Importing, indexing, or exposing source repositories.
- Adding authentication, synchronization, remote persistence, or multi-user behavior.

## Decisions

- **Use two repository-scoped agents.** The orchestrator owns discovery, OpenSpec, Trello, handoffs, and acceptance verification. The developer owns approved implementation and technical validation. This is preferred over one unrestricted agent because it makes governance boundaries auditable.
- **Use dedicated local worktrees for implementation.** Each task gets a `.worktrees/inventory-dashboard-task-<number>-<kebab-case-name>` directory and matching `feature/task-<number>-<kebab-case-name>` branch. This is preferred over editing `main` because parallel work and review need independent state.
- **Use OpenSpec as the contract and Trello as execution tracking.** OpenSpec holds behavior and design; Trello holds ownership, predecessors, status, acceptance criteria, and evidence. Neither tool replaces tests or runtime output.
- **Keep Graphify and Codebase Memory read-only and worktree-local.** They provide bounded navigation context, while source files, tests, and fresh runtime output remain authoritative. Each active worktree must refresh its own graph and index rather than sharing generated artifacts.
- **Keep the future persistence decision at SQLite.** The initial product is local, single-user, and offline-friendly, so SQLite remains appropriate. MySQL is not introduced unless a future requirement adds shared access, server-managed concurrency, or operational replication.
- **Bind future services locally and protect secrets.** Runtime services will default to `127.0.0.1`; MCP credentials remain in ignored `.env.mcp`; scans must use an explicit read-only allowlist and never ingest source contents or secret-bearing files.

## Risks / Trade-offs

- [Risk] Worktree and index discipline adds setup steps → Require exact setup and evidence in every developer handoff.
- [Risk] Local generated indexes can become stale → Refresh Graphify and verify Codebase Memory status before and after implementation.
- [Risk] Governance files can drift from actual tool behavior → Validate OpenSpec, agent configuration, tests, and fresh runtime output independently.
- [Risk] SQLite would be insufficient for future collaboration → Keep persistence behind infrastructure boundaries so a later storage change does not alter domain contracts.

## Migration Plan

No runtime migration is required. The repository bootstrap is already represented by configuration and documentation. Future implementation tasks will reference this change or accepted specs, create their own worktrees, and pass the defined validation gates.

## Open Questions

- Which exact Trello lists and labels will be used for the first product backlog?
- What is the first user-facing inventory capability to formalize after this foundation?
