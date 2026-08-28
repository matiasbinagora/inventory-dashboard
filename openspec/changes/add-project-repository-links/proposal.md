## Why

Projects currently have public-link records, but Administration does not provide dedicated fields for the one GitHub repository and one Trello backlog that identify a project. Adding optional first-class fields will make these references easier to maintain and consistently available in the project detail view.

## What Changes

- Add optional GitHub repository and Trello backlog URL fields to project administration.
- Persist and return the two URLs with the project contract and SQLite storage.
- Validate HTTPS host/path rules without fetching remote destinations.
- Show saved links on the project detail page when present.
- Add migration-safe persistence and backend/frontend/browser coverage.
- Reference implementation task `DAY-1-FEATURE-021`.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `local-project-administration`: allow editing one optional GitHub repository URL and one optional Trello backlog URL per project.
- `editorial-project-catalog`: display the optional repository/backlog references on project detail without inventing values.

## Impact

- Affected systems: Go domain/API/application, SQLite project schema, Next.js administration and detail features, and tests.
- Affected behavior: project API/data contract, local persistence, administration forms, and detail presentation.
- Non-goals: remote integrations, fetching repository/Trello content, authentication, synchronization, or multiple URLs per kind.
