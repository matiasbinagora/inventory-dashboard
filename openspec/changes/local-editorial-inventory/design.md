## Context

The first product is a local, single-user inventory for current and historical development projects. It combines an editorial portfolio feel with dashboard discovery, curated local media, public GitHub/Trello links, and manually authored history. Mobile layouts are deferred.

## Goals / Non-Goals

**Goals:** dashboard summary, desktop catalog with cards or table, project detail, local administration, curated media, manual milestones, and privacy-safe data boundaries.

**Non-Goals:** mobile layout, automatic Git/GitHub/Trello synchronization, source-code browsing, remote hosting, authentication, and multi-user behavior.

## Decisions

- Next.js routes and features will call a Go API; the API owns validation and behavior; SQLite owns local persistence.
- Public GitHub and Trello URLs are inert display links only and do not trigger imports.
- Curated media is stored in a local application-managed directory; thumbnails may point to originals for enlargement.
- A P0 visual exploration must present 2 or 3 alternatives before production implementation.
- SQLite fits the offline-friendly single-user product; storage remains behind infrastructure boundaries for future evolution.

## Conceptual Data Flow

```text
Read-only repositories -> safe manual curation -> local admin -> Go API -> SQLite + local media -> Next.js views
Source code / secrets / private content / uncurated artifacts -> excluded
```

## Risks / Trade-offs

- [Risk] Media paths break when files move -> use managed relative paths and validate existence.
- [Risk] Local service exposure -> bind to `127.0.0.1` by default.
- [Risk] Visual exploration expands into implementation -> keep P0 static or throwaway only.

## Migration Plan

No runtime migration exists. Future schema and media initialization require explicit implementation tasks.

## Open Questions

- Final status and category vocabularies, media directory layout, and post-P0 frontend/backend task split.
