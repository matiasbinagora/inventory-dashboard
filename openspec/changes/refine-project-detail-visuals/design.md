## Context

The existing Next.js project-detail page already owns the presentation of project metadata, curated media, public links, milestones, and the agentic platform. The requested change is a presentation refinement based on observed desktop output: the hero title dominates too much space, the description has a narrow reading measure, and supporting fields are rendered as dense text blocks.

The change is frontend-only. The Go API, SQLite schema, media-serving contract, public URL validation, and route structure remain authoritative and unchanged.

## Goals / Non-Goals

**Goals:**

- Establish a smaller, more technical project-detail title hierarchy.
- Improve description readability with a wider bounded measure and stable line-height.
- Make Graphify report media, technologies, public links, agentic platform, and demo video easy to scan.
- Preserve the existing data ownership, navigation, media URLs, and responsive behavior.
- Cover observable labels and layout behavior with frontend and Playwright regression tests.

**Non-Goals:**

- No API, persistence, or domain-model changes.
- No new external dependencies or icon service.
- No changes to public-link validation or remote-link destinations.
- No dashboard redesign, media import, authentication, or new project-detail capabilities.

## Decisions

- **Keep presentation inside the existing project-detail feature.** Update the existing component and stylesheet rather than adding a new design-system abstraction. This keeps route ownership and data flow unchanged and limits the surface area of the visual fix.
- **Use CSS for the technical treatment.** Reduce the existing responsive heading scale, tighten tracking, and use existing uppercase eyebrow conventions instead of introducing a new font or visual dependency. CSS is preferable to runtime layout logic because the treatment is deterministic and works during server rendering.
- **Use semantic lists and link elements.** Technologies will render as a `ul`; public links and agentic-platform references will remain real anchors with clearer grouping and hover/focus states. This improves scanning and keyboard access without changing data.
- **Use bounded readable measures.** The hero description will receive a larger `max-width` and controlled line-height; text justification may be applied only where the existing responsive layout has enough width, avoiding excessive word spacing on narrow viewports.
- **Preserve media behavior.** The Graphify Report label changes only the section heading. Existing image lightbox, local media URLs, video controls, and public HTTPS video references remain unchanged; the demo section label changes to `Demo Video`.

## Risks / Trade-offs

- [Risk] A smaller title can reduce visual impact → Mitigation: retain strong weight and responsive scale while reducing only the oversized desktop footprint.
- [Risk] Justified text can create uneven gaps on narrow screens → Mitigation: use a readable max-width and disable justification below the supported narrow breakpoint if needed.
- [Risk] Link styling could obscure the distinction between local metadata and public destinations → Mitigation: keep explicit link labels/domains and visible keyboard focus states.
- [Risk] The two frontend bug fixes recently merged may be touched by shared styles → Mitigation: inspect current `main`, preserve administration/dashboard selectors, and add focused detail-page regression coverage only.

## Migration Plan

No data or deployment migration is required. Implement the component/style/test changes, run the existing frontend and browser validation, and roll back by reverting the single feature PR if the visual result is unacceptable.
