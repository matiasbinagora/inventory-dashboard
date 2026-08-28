## Context

The project contract already supports validated public GitHub and Trello links as a collection, but the administration workflow lacks dedicated one-per-project fields. The requested behavior is local metadata management only: users enter one repository URL and one backlog URL, and the application stores and displays the validated references without contacting either service.

## Goals / Non-Goals

**Goals:**

- Add optional `github_repository_url` and `trello_backlog_url` project fields.
- Preserve existing projects and allow empty values.
- Validate each field using the existing public URL rules and enforce one value per kind.
- Persist fields through SQLite migration and return them through existing project endpoints.
- Present saved links in administration and detail with existing safe link behavior.

**Non-Goals:**

- No remote API calls, repository inspection, Trello synchronization, or authentication.
- No replacement of existing arbitrary public-link records unless required by the current contract.
- No support for multiple GitHub repositories or Trello boards in these dedicated fields.

## Decisions

- **Use dedicated nullable project columns.** These fields have one-to-one cardinality and are edited with the project, so columns are clearer and simpler than introducing another relation. Existing public-link records remain compatible.
- **Reuse domain URL validation.** Validate HTTPS, allowed host, non-empty path, and forbidden userinfo/private destinations in the domain/application seam. The server never fetches the URL.
- **Use an idempotent SQLite migration.** Add nullable columns to the existing projects table without changing or rewriting existing rows. Reads map NULL to omitted/empty optional JSON fields according to current frontend types.
- **Keep one shared form flow.** Add two optional controls to the existing Administration project form and submit them through the existing create/update project API client. The detail page renders only present values as actionable anchors.

## Risks / Trade-offs

- [Risk] Existing SQLite databases lack the new columns → Mitigation: apply a versioned/idempotent schema migration before queries.
- [Risk] Dedicated fields could diverge from existing public-link records → Mitigation: preserve the existing collection and define dedicated fields as the canonical one-per-kind administration values without remote synchronization.
- [Risk] Users may enter a valid host with an unusable path → Mitigation: require HTTPS and non-empty path using the existing validator.
- [Risk] Long URLs can affect layout → Mitigation: use existing wrapping/focus styles and responsive browser coverage.

## Migration Plan

On API startup, apply the SQLite schema migration for nullable project URL columns. Existing rows remain valid with empty values. Rollback is a code revert; the extra nullable columns are harmless to older code.
