## 1. Media storage boundary

- [ ] 1.1 DAY-1-TASK-013 Create a dedicated worktree and document the managed local media root and configuration.
- [ ] 1.2 DAY-1-TASK-013 Implement safe curated asset registration/storage with size, type, traversal, and symlink checks.
- [ ] 1.3 DAY-1-TASK-013 Persist only managed relative media references and preserve existing SQLite compatibility.

## 2. Media serving and integration

- [ ] 2.1 DAY-1-TASK-013 Serve existing managed image/video binaries with correct content types and safe 404/400 responses.
- [ ] 2.2 DAY-1-TASK-013 Integrate the detail gallery with real local media without duplicating Go validation in Next.js.
- [ ] 2.3 DAY-1-TASK-013 Add tests for successful delivery, missing assets, traversal, escaping symlinks, unsupported types, and restart persistence.

## 3. Verification

- [ ] 3.1 DAY-1-TASK-013 Run Go/frontend tests, integration, race, vet, build, and reviewer-generated Playwright flows with temporary non-sensitive media.
- [ ] 3.2 DAY-1-TASK-013 Run `openspec validate --all --strict`, refresh Graphify and Codebase Memory, run `git diff --check`, and document runtime commands and limitations.
- [ ] 3.3 DAY-1-TASK-013 Create a focused PR and record the Trello transition evidence; do not merge or mark Done.
