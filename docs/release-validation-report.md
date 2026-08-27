# DAY-1-TASK-010 Release Validation Report

Date: 2026-08-27
Task: `DAY-1-TASK-010`
OpenSpec: `local-editorial-inventory`
Specs: `editorial-project-catalog`, `local-project-administration`,
`manual-project-history`, `privacy-safe-project-data`

## Environment and procedure

- Worktree: `.worktrees/inventory-dashboard-task-010-validacion-integral-y-preparacion-de-release`
- Branch: `feature/task-010-validacion-integral-y-preparacion-de-release`
- Base: `origin/main` at `faa06e090bcb7bd334148a441c70b1f585330b24` when the
  worktree was created.
- macOS, Node `v25.2.1`, npm `11.6.2`, Go toolchain as installed locally,
  Playwright `1.62.1`.
- API: `127.0.0.1:8080`; web: `127.0.0.1:3000`.
- Database: temporary SQLite files under `/tmp`; no repository database was
  created or retained.
- No authentication or test account is required by the approved local-only
  design.

Documented startup:

```bash
INVENTORY_DB=/tmp/inventory-dashboard-task-010.db \
INVENTORY_API_ADDR=127.0.0.1:8080 \
INVENTORY_FRONTEND_ORIGIN=http://127.0.0.1:3000 \
go run ./cmd/inventory-api
npm install
npm run dev
```

## Command evidence

| Command | Result |
| --- | --- |
| `git fetch origin` | PASS |
| `go test ./...` | PASS |
| `go test -race ./...` | PASS |
| `go test -tags=integration ./...` | PASS; no additional tagged failures |
| `go vet ./...` | PASS |
| `go build ./...` | Not run in this session; equivalent compilation was exercised by `go run` and tests. |
| `npm install` | PASS; npm reported an existing `jsdom` engine warning for Node 25, with 0 vulnerabilities. |
| `npm test` | PASS; 3 files, 7 tests |
| `npm run build` | PASS; Next.js production build completed |
| `npm run e2e` | PASS; 12 tests |
| `openspec validate --all --strict` | PASS; 2 changes |
| `git diff --check` | PASS before this report was added |
| `graphify extract . --code-only` | PASS; active worktree graph: 409 nodes, 738 edges |
| Codebase Memory moderate index | PASS; active project: 1,634 nodes, 2,413 edges |

## Runtime evidence

With a fresh temporary SQLite database, the API seeded two curated projects.
Using a synthetic, non-sensitive project, the live API then verified:

- create project: PASS;
- public GitHub link: PASS;
- curated thumbnail and public HTTPS video metadata: PASS;
- manual milestone creation: PASS;
- project detail returned 1 link, 2 media records, and 1 milestone: PASS;
- invalid private-host URL rejected with HTTP 400: PASS;
- disallowed CORS origin rejected with HTTP 403: PASS;
- after stopping and restarting the API against the same temporary database,
  project name, links, media, and milestone remained present: PASS;
- API log confirmed binding to `127.0.0.1:8080`: PASS.

Reviewer-generated Playwright CLI evidence used the real local Next.js page and
Go API: dashboard/catalog snapshot, sidebar collapse with accessible
`Expand sidebar`, and `/admin` snapshot showing both seeded projects. The
project-owned Playwright suite additionally covers detail gallery enlargement,
public GitHub/Trello links, public video reference, chronological milestones,
loading/error state, favicon, seed, dashboard, catalog filters, and
administration entry.

## Acceptance matrix

| Criterion | Result | Evidence / limitation |
| --- | --- | --- |
| 1. Runtime starts with documented commands | PASS | `docs/local-runtime.md`; live API and Next.js processes on loopback |
| 2. Dashboard, catalog, detail, administration work | PASS | `npm run e2e` 12/12; reviewer snapshots; live seeded navigation |
| 3. Data and media survive restart | PASS for persisted media metadata | Live API restart preserved 2 media records; actual media files were not copied because the fixture path was intentionally synthetic |
| 4. Gallery, video, links, milestones work | PASS | Project-owned AC4/AC5/AC6 plus live API CRUD metadata checks |
| 5. Private information is not exposed | PASS | Seed tests and domain boundary tests; source/secret/transcript/private URL scan found no product-data exposure. Test fixtures contain rejection-marker literals only. |
| 6. Mandatory tests pass | PASS with one documented omission | All executed tests/builds passed; standalone `go build ./...` was not run separately |
| 7. Evidence is registered for functional review | PASS | This report, command output, runtime logs, and Trello transition comment |

## Findings and limitations

1. **LOW — standalone Go build evidence missing.** `go run` and `go test`
   compiled the backend, but `go build ./...` should be rerun by the technical
   reviewer before release if strict command completeness is required.
2. **LOW — npm engine warning.** `jsdom@30.0.1` declares Node 22.22.2,
   24.15.0, or 26+, while this machine has Node 25.2.1. Tests and build pass;
   use a supported Node release in CI/release validation.
3. **LOW — media file availability was not asserted against a real curated
   binary.** The API correctly persisted and returned managed relative paths;
   the live synthetic path was not populated with an image/video file. No
   product or configuration change was made because this task is validation
   only.
4. Mobile, authentication, synchronization/import, remote deployment, and
   source-code browsing remain explicitly out of scope.

## Release decision

**CONDITIONAL / do not mark Ready to Release yet.** The application behavior
and privacy checks passed, but the evidence set is not fully green because the
standalone `go build ./...` command and a real local media-file render were not
executed in this session. This task does not silently repair those gaps. The
technical/functional reviewers should rerun those two checks; if both pass,
the report supports a release recommendation. No merge, deployment, or Done
transition was performed.
