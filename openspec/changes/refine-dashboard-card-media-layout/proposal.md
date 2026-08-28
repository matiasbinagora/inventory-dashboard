## Why

The current detail and catalog views still present media and project metadata with inconsistent scale and too much raw text. The next visual pass should give the Graphify report the same full-width presence as the demo video while making catalog cards compact and scannable.

## What Changes

- Render the single Graphify Report image at the full content width while preserving its intrinsic aspect ratio.
- Truncate catalog-card descriptions to a compact fixed line count.
- Render catalog-card technologies as name-only pills without embedded links.
- Keep full technology links and metadata available on the project detail page.
- Add responsive and browser regression coverage.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `editorial-project-catalog`: change catalog-card density and project-detail media presentation.

## Impact

- Affected systems: Next.js dashboard/catalog and project-detail features, shared styles, frontend/Playwright tests.
- Affected behavior: presentation only; no API, persistence, media-serving, or stored metadata changes.
- Reference task: `DAY-1-FEATURE-020`.
