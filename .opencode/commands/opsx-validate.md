---
description: Validate OpenSpec artifacts and acceptance readiness
---

Validate an OpenSpec change and its implementation against the approved
requirements.

Input: optionally `/opsx-validate <change-name>`.

1. Read `AGENTS.md`, the Trello card, and the referenced OpenSpec artifacts.
2. Run `openspec validate --all --strict`.
3. If a change name is provided, run `openspec status --change "<change-name>"
   --json` and inspect its proposal, specs, design, and tasks.
4. Map every acceptance criterion to fresh test, runtime, or operational
   evidence.
5. For frontend changes, run the required Playwright-backed validation.
6. Report `PASS` or `FAIL`, evidence, skipped checks, remaining risks, and the
   permitted Trello transition.

Do not modify production code, merge pull requests, or move work to Done.
