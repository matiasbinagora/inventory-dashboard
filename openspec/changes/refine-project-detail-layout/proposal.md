## Why

The first detail-page refinement improved labels and hierarchy but still renders long descriptions and metadata as dense, difficult-to-scan blocks. A follow-up visual pass is needed to make the page feel intentional: the description should use the full content width, structured names and links should be separated, and duplicate report imagery should not be shown twice.

## What Changes

- Make the project description span the full detail content width across both columns.
- Present technology and agentic-platform entries as bullets with a name and a smaller blue URL below it.
- Keep public links visually actionable without displaying raw URL-heavy paragraphs.
- Display only one representative image in `Graphify Report` when equivalent images are duplicated.
- Preserve the current title scale, `Graphify Report`, `Demo Video`, media behavior, navigation, and responsive support.
- Reference implementation task `DAY-1-FEATURE-019`.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `editorial-project-catalog`: refine the responsive detail layout and structured metadata presentation.

## Impact

- Affected systems: Next.js project-detail feature, detail-page styles, and frontend/Playwright tests.
- Affected behavior: user-visible layout, metadata parsing/presentation, and duplicate image presentation only.
- No API, persistence, media-serving, security, or external dependency changes.
