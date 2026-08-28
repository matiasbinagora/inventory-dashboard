## 1. Contract And Persistence

- [x] 1.1 DAY-1-FEATURE-021 Create a dedicated worktree from `origin/main` and record Graphify/Codebase Memory context.
- [x] 1.2 Extend the project domain/API contract with optional GitHub repository and Trello backlog URLs and validate allowed HTTPS hosts.
- [x] 1.3 Add a migration-safe SQLite schema update and persist/read both optional fields across restart.
- [x] 1.4 Add API and persistence tests for valid, empty, invalid, and replacement values.

## 2. Administration And Detail

- [x] 2.1 Add optional GitHub repository and Trello backlog inputs to Administration and submit them through the existing project save flow.
- [x] 2.2 Show controlled validation errors and preserve values when the form reloads.
- [x] 2.3 Render present dedicated links on project detail with accessible labels and no remote fetching.

## 3. Verification

- [x] 3.1 Add frontend tests for form fields, validation feedback, and detail link rendering.
- [x] 3.2 Add integration and reviewer-generated Playwright coverage for create/edit/reload and link navigation behavior.
- [x] 3.3 Run Go tests, integration, race, vet, build, frontend tests/build, Playwright, `openspec validate --all --strict`, Graphify, Codebase Memory, and `git diff --check`.
- [ ] 3.4 Create the focused commit and PR, document evidence, and move the Trello card to `Code Review` without merging.
