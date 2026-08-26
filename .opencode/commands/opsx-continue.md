---
description: Continue completing an OpenSpec change whose artifacts are incomplete
---

Continue an incomplete OpenSpec change without implementing production code.

Input: `/opsx-continue <change-name>`.

1. Run `openspec status --change "<change-name>" --json`.
2. Read the existing proposal, design, specs, and tasks identified by the
   status output.
3. Run `openspec instructions <artifact> --change "<change-name>" --json` for
   the next incomplete artifact.
4. Complete only the missing planning artifact using its schema instructions.
5. Repeat until the change is implementation-ready or a product decision is
   required.
6. Run `openspec validate --all --strict` and report the result.

Do not implement application code, change Trello scope, or archive the change.
