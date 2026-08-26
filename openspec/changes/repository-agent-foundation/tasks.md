## 1. Repository Governance

- [ ] 1.1 Verify `inventory-orchestrator` is the primary agent and its permissions prevent production-code edits.
- [ ] 1.2 Verify `inventory-developer` is hidden, implementation-scoped, and requires `/opsx-apply` before coding.
- [ ] 1.3 Verify repository instructions define OpenSpec, Trello, worktree, review-gate, and security obligations.

## 2. Worktree And Context Isolation

- [ ] 2.1 Verify the documented task worktree and branch naming conventions with `git worktree list --porcelain`.
- [ ] 2.2 Verify `.env.mcp`, `.worktrees/`, Graphify output, Codebase Memory cache, databases, and generated artifacts are ignored.
- [ ] 2.3 Verify Graphify and Codebase Memory are configured with the active worktree as their local scope and read-only navigation role.

## 3. Tooling And Validation

- [ ] 3.1 Verify configured MCP gateways use the ignored local secrets file without exposing its contents.
- [ ] 3.2 Run `openspec validate --all --strict` and resolve all validation failures.
- [ ] 3.3 Run `git diff --check` and inspect the final tracked-file candidate list for secrets or private metadata.
- [ ] 3.4 Record Graphify and Codebase Memory status as supporting evidence; do not substitute either index for source, test, or runtime validation.
