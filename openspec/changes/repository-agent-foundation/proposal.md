## Why

The repository needs a predictable operating model before product implementation begins. Without explicit agent ownership, isolated worktrees, OpenSpec traceability, and privacy-safe local context, implementation work could bypass review gates or expose local project data.

## What Changes

- Establish the `inventory-orchestrator` as the primary governance agent.
- Establish the hidden `inventory-developer` as the implementation agent for frontend, backend, persistence, tests, and pull requests.
- Define OpenSpec and Trello as the sources of truth for requirements, status, acceptance criteria, and evidence.
- Require dedicated `.worktrees/` worktrees and task branches for implementation work.
- Configure local Graphify and Codebase Memory services as read-only, worktree-scoped context sources.
- Add repository instructions, ignored local artifacts, MCP gateway configuration, and supporting skills for the planned stack.
- Do not introduce product UI, HTTP endpoints, database schema, or repository scanning behavior in this change.

## Capabilities

### New Capabilities

- `agent-governance`: Defines agent responsibilities, permissions, required workflows, and review gates.
- `isolated-implementation-worktrees`: Defines branch and worktree isolation requirements for implementation tasks.
- `privacy-safe-development-context`: Defines safe use of local project context, Graphify, Codebase Memory, and ignored artifacts.

### Modified Capabilities

- None.

## Impact

- Affected tooling and documentation: `AGENTS.md`, `README.md`, `.opencode/`, `.agents/`, `.gitignore`, and OpenSpec configuration.
- Affected workflow: OpenSpec planning and validation, Trello handoffs and transitions, developer worktree setup, and review evidence.
- Affected local integrations: GitHub/Trello MCP gateways, Graphify, and Codebase Memory.
- No production API, frontend runtime, persistence layer, or external user-facing behavior is added by this change.
