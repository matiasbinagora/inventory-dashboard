## ADDED Requirements

### Requirement: Agent responsibilities are separated
The repository MUST define an orchestrator responsible for discovery, OpenSpec governance, Trello coordination, handoffs, and acceptance verification, and a developer responsible for approved implementation and technical validation.

#### Scenario: Planning request is handled by the orchestrator
- **WHEN** a request introduces or changes product behavior, API, persistence, UI, security, or workflow
- **THEN** the orchestrator MUST begin discovery through OpenSpec exploration and MUST NOT implement production code

#### Scenario: Approved implementation is handed to the developer
- **WHEN** a complete Trello handoff references an OpenSpec change, spec, or documented exception
- **THEN** the developer MUST implement only that approved scope and MUST report validation evidence before requesting code review

### Requirement: Workflow gates are enforced
The repository MUST require OpenSpec validation, technical review, and functional review before an implementation task can be considered complete.

#### Scenario: Missing contract blocks implementation
- **WHEN** an implementation task has no OpenSpec reference or documented exception
- **THEN** the task MUST NOT be considered ready for implementation

#### Scenario: Review evidence is incomplete
- **WHEN** required tests or OpenSpec validation fail
- **THEN** the task MUST remain blocked from code review or release progression until the failure is resolved
