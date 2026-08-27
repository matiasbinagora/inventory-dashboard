# QA Progress Report

Date: 2026-08-26
Scope: merged progress through DAY-1-TASK-005
OpenSpec: `local-editorial-inventory`

## Decision

**FAIL / release blocked**

The implemented flows are usable and the automated checks pass, but the
release gate remains blocked by findings listed below. This report is evidence
for review; it does not modify production code or change workflow state.

## Reviewed Progress

| Task | Scope | Observed status |
| --- | --- | --- |
| DAY-1-TASK-001 | Visual exploration and direction C | Merged and complete |
| DAY-1-TASK-002 | Project, link, media, and milestone contracts | Merged and complete |
| DAY-1-TASK-003 | SQLite persistence | Merged and complete |
| DAY-1-TASK-004 | Go local HTTP API | PR #4 merged; reviewed |
| DAY-1-TASK-005 | Desktop dashboard and editorial catalog | PR #5 merged; reviewed |
| DAY-1-TASK-006 | Desktop project detail | Not implemented; intentionally next scope |
| DAY-1-TASK-007 | Local administration | Not implemented; intentionally next scope |
| DAY-1-TASK-008 | Privacy-safe seed metadata and boundary tests | Not implemented; intentionally next scope |

## Acceptance Validation

### DAY-1-TASK-003 and DAY-1-TASK-004

- SQLite schema and persistence packages are present and reproducible: PASS.
- Project CRUD API is available: PASS.
- Invalid project names return HTTP 400: PASS via API tests and live request.
- Invalid or non-allowlisted public URLs return HTTP 400: PASS via live request.
- Local media paths are constrained to managed `media/` paths: PASS via live request.
- Media and milestone endpoints exist: PASS by route and package inspection.
- No remote importer or synchronization was executed: PASS within reviewed scope.
- API defaults to `127.0.0.1:8080`: PASS by `cmd/inventory-api/main.go` and runtime log.

### DAY-1-TASK-005

- Dashboard renders data from `/api/projects`: PASS in live browser; response was
  obtained through the Next.js rewrite and showed the seeded safe QA record.
- Differentiated metric cards render and reach their final values: PASS; live
  snapshot showed `1`, `0`, `0`, and `1` for the seeded record.
- Sidebar collapses while retaining accessible expand/collapse semantics: PASS;
  Playwright snapshot showed `aria-expanded="false"` and the `Expand sidebar`
  accessible name.
- Cards/table presentation switch works: PASS; live browser rendered the table
  with project, description, technologies, and links columns.
- Search and technology/platform filters update visible results: PASS; live
  search for `not found` rendered the empty filtered state and visible results
  changed to `0`.
- Card/table navigation reaches the project route: PASS; generated link was
  `/projects/c59e572ee42bb74f1d5dcb5e8be669b2`.
- Loading, API error, and empty catalog states are represented: PARTIAL. Loading
  and empty states were observed; the error branch exists and is covered by
  project E2E assertions, but a live API-down browser run was not isolated in
  this session because another local process owned port 8080.

## Commands And Evidence

Executed from the TASK-005 worktree:

- `npm test`: PASS, 1 file and 3 tests.
- `npm run build`: PASS.
- `go test ./...`: PASS.
- `go test -race ./...`: PASS.
- `go vet ./...`: PASS.
- `openspec validate --all --strict`: PASS, 2 changes.
- `git diff --check`: PASS.
- `npx --no-install playwright cli open http://127.0.0.1:3000`: PASS.
- Reviewer-generated browser actions: reload, inspect API-backed catalog, switch
  table view, collapse sidebar, fill search, and inspect snapshots.
- Live API checks against temporary SQLite runtime: safe project creation PASS;
  invalid GitHub host `400`; absolute private media path `400`.

Runtime services were local only: Next.js at `127.0.0.1:3000` and Go API at
`127.0.0.1:8080`. Test data was safe synthetic metadata (`QA Atlas`) and was
stored in `/tmp/inventory-dashboard-qa.db`.

## Findings

### HIGH: API has no explicit CORS policy

`internal/api/http.go` wraps the mux with JSON headers but does not set or
validate an `Origin` policy. The architecture and security instructions require
CORS to be restricted to the local frontend origin. The Next.js rewrite makes
normal browser traffic same-origin, but it does not enforce the API boundary
for direct callers.

Recommendation: add an explicit localhost-only CORS policy, including a tested
preflight path, or document and enforce that the API is never directly exposed.

### MEDIUM: Child writes can leave partially-created projects

`internal/application/inventory.go` creates the project and then writes links,
media, and milestones individually. If a child write fails, the project and
earlier children remain persisted while the request returns an error. This can
leave the catalog inconsistent.

Recommendation: make aggregate creation transactional in the SQLite repository
or add an application-level rollback contract and tests.

### LOW: Runtime emits a missing favicon request

The live browser console reported `404` for `/favicon.ico`. This does not block
the current acceptance criteria, but it creates avoidable console noise and
should be resolved before release polish.

### LOW: Frontend test command emits a Vite configuration warning

`npm test` passes but warns that `vitest.config.ts` uses ESM syntax while the
package is treated as CommonJS. This is not currently a functional failure,
but future Vite behavior may make the warning actionable.

## Not Validated

- Full project detail presentation and image enlargement: reserved for TASK-006.
- Administration workflows and primary-thumbnail selection: reserved for TASK-007.
- Curated seed import and complete privacy-boundary test suite: reserved for TASK-008.
- Mobile behavior: explicitly out of scope.
- Accessibility audit beyond the reviewed navigation controls: not performed.
- Browser run with the API process intentionally stopped: not isolated because
  port 8080 was already occupied by a local runtime process.

## Release Recommendation

Do not mark the progress as release-ready until the HIGH and MEDIUM findings
have an accepted fix or documented exception. TASK-006 may proceed as planned,
but these findings should be tracked before final release acceptance.
