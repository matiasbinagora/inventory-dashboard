# DAY-1-TASK-001 · Visual exploration

This is a throwaway, static desktop prototype for the `local-editorial-inventory`
P0 discovery gate. It contains three comparable directions in one reviewable
artifact:

- **A · Field notes** — editorial dashboard, card-led catalog with a table-ready
  secondary mode, and story-first detail.
- **B · Gallery wall** — image-led browsing for projects with strong visual
  evidence.
- **C · Mission control** — dense dashboard and table-first catalog for fast
  maintenance. This is the selected direction, with the review refinements
  below.

All directions share the same detail and administration contract: fictional
project metadata, explicit public display-only GitHub/Trello links, technologies,
agentic platform, local video reference, expandable full-size placeholder SVGs,
and manual milestones. No repository source, credentials, private URLs,
transcripts, or reference-project media are included.

## Review locally

From this directory run:

```bash
python3 -m http.server 4173 --bind 127.0.0.1
```

Open <http://127.0.0.1:4173/> at a desktop viewport (recommended 1440×1000).
Select A/B/C at the top, inspect the shared detail and administration sections,
and click any detail thumbnail to open its full-size local SVG placeholder.
In C, the four dashboard metrics are rounded, color-coded cards whose numbers
count up on load. Use the icon sidebar button to collapse and expand the rail;
its labels are visible when expanded and remain available as accessible names.
In Administration, choose a project and an approved primary thumbnail: the
choice immediately updates that project's `Primary` cell in the C catalog and
the preview text. The form buttons only show static-study feedback; they do not
persist data or call an API.

## Recommendation and decision gate

User selection: **C · Mission control**. The approved refinements are rounded,
color-coded metric cards with incremental number animation, an explicit
administration-to-catalog primary-thumbnail choice, and a collapsible icon
sidebar with expanded accessible labels. The accepted trade-off is that the
catalog is denser and less portfolio-like than A.

A is rejected as the selected direction because its editorial warmth makes
maintenance scanning slower. B is rejected because image-first browsing can
hide status and metadata. The C refinements remain throwaway visual evidence;
they do not define production implementation yet.

## Scope and risks

This artifact intentionally does not create a Next.js app, Go API, SQLite
database, sync behavior, authentication, mobile layout, or production styling.
The main follow-up risks are final status/category vocabulary, managed media
directory conventions, and whether the selected card/table balance remains
directory conventions, whether the selected card/table balance remains clear
with real curated content, and whether motion should respect a future reduced
motion preference.
