## Context

The project detail page currently gives the Graphify Report a small card treatment while the demo video spans the primary content width. Catalog cards also expose long descriptions and technology strings that include URLs, making the grid difficult to scan.

## Goals / Non-Goals

**Goals:**

- Give the single report image and demo video consistent full-width media presence without distorting images.
- Keep catalog cards compact with predictable description height.
- Keep pills strictly name-only in the catalog while preserving detail-page links.
- Preserve semantic HTML, keyboard access, existing data, and responsive behavior.

**Non-Goals:**

- No API or data-model changes.
- No deletion or mutation of media/technology records.
- No change to detail-page metadata link availability.
- No new dependencies or dashboard navigation redesign.

## Decisions

- **Use intrinsic image sizing.** The Graphify image will use the primary content width with `height:auto` and a bounded maximum, matching the video footprint without stretching or cropping report content.
- **Truncate only at the card presentation boundary.** Apply CSS line clamping to card descriptions so the stored description and detail page remain complete. The card will use a stable compact line count and a responsive fallback.
- **Use display-only technology names in cards.** Reuse the existing metadata parsing/name extraction, but render only names as pills in the catalog. Detail pages remain the source for associated URLs.
- **Keep media deduplication presentation-only.** The existing one-representative Graphify behavior remains; no records are removed or rewritten.

## Risks / Trade-offs

- [Risk] Full-width report images may be very tall → Mitigation: preserve aspect ratio but constrain the rendered max-height with `object-fit:contain` only if needed, and test real report dimensions.
- [Risk] Line clamping can hide important card context → Mitigation: retain full text on detail pages and use an accessible link to the project detail.
- [Risk] Parsing unusual technology strings may lose labels → Mitigation: fallback to the original text as a name when no structured name is available.

## Migration Plan

No migration is required. This is a frontend-only change and can be rolled back by reverting its PR.
