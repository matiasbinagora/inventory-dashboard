## Context

The current project detail route receives technologies and agentic-platform text from the existing project contract. Some values contain several named tools followed by URLs in a single string, which currently creates dense paragraphs. The detail layout also places the hero in a constrained width and renders equivalent curated images as separate cards.

## Goals / Non-Goals

**Goals:**

- Let the hero description use the full content width while remaining readable.
- Split URL-bearing metadata into scannable name/link entries.
- Use consistent blue, smaller links below each metadata name.
- Deduplicate equivalent Graphify report images at presentation time.
- Preserve all existing data, URLs, media controls, and responsive behavior.

**Non-Goals:**

- No API or database migration.
- No change to stored metadata or URL validation.
- No remote fetching, image processing, or deletion of media records.
- No redesign of dashboard or administration.

## Decisions

- **Parse presentation strings locally.** Add a small detail-feature parser for the established `Name: URL` and `Name — description: URL` forms. This avoids changing the API contract while allowing semantic markup. Segments without URLs remain safe text.
- **Use semantic list markup.** Technologies and platform entries render as `ul`/`li` items; names remain primary text and links render as separate smaller anchors below. Existing HTTPS targets and `rel="noreferrer"` behavior remain unchanged.
- **Make the hero span the grid.** Keep the detail grid for content below the hero, but set the description/hero width to the full detail container rather than the primary-column width. Use responsive max-width and line-height instead of aggressive justification.
- **Deduplicate only the rendered report.** Filter image media by a stable source identity, preferring the first curated representative. Do not delete or mutate persisted media records, and preserve original-media lightbox lookup for the selected item.

## Risks / Trade-offs

- [Risk] Existing free-form metadata may not match the parser → Mitigation: preserve unmatched text as a plain item and never discard content.
- [Risk] Deduplicating by source could hide intentionally different records → Mitigation: deduplicate only exact equivalent sources in the Graphify image presentation, retaining the first representative.
- [Risk] Full-width descriptions may become too long to read → Mitigation: apply a bounded character measure and stable line-height within the full two-column container.
- [Risk] Long URLs can overflow → Mitigation: use `overflow-wrap:anywhere`, smaller blue link styling, and narrow-viewport regression coverage.

## Migration Plan

No migration is required. Deploy as a frontend-only change; rollback is a PR revert if the visual or responsive result is not accepted.
