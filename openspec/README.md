# Inventory Dashboard OpenSpec Workflow

OpenSpec is the source of truth for product requirements, API and data
contracts, design decisions, implementation tasks, and traceability. Trello is
the source of truth for workflow state, ownership, acceptance criteria, and
review evidence.

## Required Flow

```text
/opsx-explore -> /opsx-propose -> /opsx-apply -> /opsx-validate
                                      |
                                      v
                             /opsx-sync -> /opsx-archive
```

Every implementation card must reference an active OpenSpec change, an
accepted spec, or a documented exception. New user-visible behavior, API
contracts, persistence behavior, security rules, and workflow changes require
a formal change.

## Validation

```bash
openspec validate --all --strict
```

Only the orchestrator may sync or archive changes, and only after technical and
functional review gates pass.
