## Why

The project detail page presents useful metadata, media, links, and agentic context, but its current hierarchy is too oversized and its supporting fields are difficult to scan. This refinement will make the detail view more technical, readable, and explicit without changing its data or navigation behavior.

## What Changes

- Reduce and refine the project detail title treatment for a more technical visual hierarchy.
- Expand and improve the project description reading measure.
- Rename `Visual evidence` to `Graphify Report`.
- Present technologies as scannable bullet items.
- Improve the visual treatment of public links and the agentic platform field.
- Rename `Watch the workflow` to `Demo Video`.
- Add frontend and browser coverage for the revised labels, layout, and responsive readability.
- Reference implementation task `DAY-1-FEATURE-018`.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `editorial-project-catalog`: refine the project detail presentation while preserving existing project metadata, media, links, and navigation behavior.

## Impact

- Affected systems: Next.js project-detail feature, shared detail-page styles, and frontend/Playwright tests.
- Affected behavior: user-visible labels and presentation only; no API, persistence, or media contract changes.
- Non-goals: redesigning the dashboard, changing project data, changing link destinations, adding remote services, or changing responsive support boundaries.
