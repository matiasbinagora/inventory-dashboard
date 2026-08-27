## Why

The release validation found that curated media metadata persists correctly but
the referenced local binary path returns HTTP 404, so the detail gallery cannot
render real images or videos. This change is needed before the local dashboard
can be considered release-ready while preserving the privacy-safe, single-user
design.

## What Changes

- Add a local media storage/serving capability for explicitly curated assets.
- Make media available to the detail gallery using managed relative paths.
- Validate path ownership, file type, and traversal boundaries before serving or
  storing an asset.
- Provide a clear not-found response for missing managed assets.
- Add runtime and Playwright evidence using non-sensitive local test media.
- Reference implementation task `DAY-1-TASK-013`.

## Capabilities

### New Capabilities

- `local-media-serving`: Store and serve explicitly curated local media within
  the local application boundary.

### Modified Capabilities

None. The new capability contains the additional media-serving contract while
existing privacy and catalog requirements remain in force.

## Impact

- Affected systems: Go API/application, SQLite media metadata, Next.js detail
  route, and local runtime configuration.
- Affected behavior: local media upload/registration and browser serving.
- No remote storage, authentication, synchronization, source import, or public
  deployment is introduced.
