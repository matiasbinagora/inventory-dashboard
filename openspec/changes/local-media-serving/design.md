## Context

The application currently stores curated media metadata such as
`media/project/image.png`, but no local runtime serves the managed asset. The
detail page therefore receives a valid record and renders a broken resource.
The application is local-first, single-user, and must not expose source code or
private project content.

## Goals / Non-Goals

**Goals:**

- Serve approved local image and video assets from a managed media root.
- Keep the Go API responsible for validation and media ownership.
- Preserve relative media references across SQLite restart.
- Make the browser detail gallery render a real local binary.
- Test traversal, unsupported types, missing files, and successful delivery.

**Non-Goals:**

- Remote object storage, CDN, cloud deployment, or synchronization.
- Importing media from reference repositories automatically.
- Serving arbitrary filesystem paths.
- Authentication or public network exposure.

## Decisions

- **Managed local root:** Use a configurable local media directory, defaulting
  to a runtime-owned `media/` directory. Relative asset references are resolved
  beneath that root and never from arbitrary absolute paths.
- **Go-owned serving boundary:** Add the media route to the Go service or a
  Go-owned handler so path validation and local bind policy are centralized.
  Next.js continues to consume the Go API and does not duplicate media rules.
- **Explicit curation:** Upload/registration requires the existing curated flag;
  only approved image/video MIME types and safe relative paths are accepted.
- **Missing asset behavior:** Return HTTP 404 for a missing managed file without
  revealing filesystem paths. Invalid or escaping paths return HTTP 400/403.
- **SQLite remains appropriate:** SQLite continues to store metadata and the
  managed relative reference. Binary files remain local filesystem assets,
  which matches the single-user local-first boundary and avoids a new service.
- **Atomic file handling:** Write uploads to a temporary file under the managed
  root, validate the final path, then rename atomically to prevent partial
  assets.

## Risks / Trade-offs

- [Risk] A user can delete a referenced file outside the app → Return a clear
  404 and expose missing-asset state in the detail UI.
- [Risk] Path traversal could expose local files → Reject absolute paths,
  normalized paths containing `..`, symlinks escaping the root, and unknown
  extensions; test each boundary.
- [Risk] Large local videos consume disk → Enforce a configured size limit and
  report a validation error before persistence.
- [Risk] Existing metadata references may not have binaries → Keep metadata
  compatible and let Administration attach approved assets later.

## Migration Plan

1. Create/use the managed media root without changing existing SQLite records.
2. Add the serving and upload/registration validation paths.
3. Add temporary test assets only during tests; do not commit binaries unless
   explicitly approved.
4. Verify existing records remain readable and missing assets return 404.
5. Roll back by disabling the route and retaining metadata; no schema rollback
   is required unless implementation adds a migration.

## Open Questions

- The exact user-facing upload control and media approval workflow remain part
  of Administration and may be refined after this serving capability.
