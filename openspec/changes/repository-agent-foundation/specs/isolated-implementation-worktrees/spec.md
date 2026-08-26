## ADDED Requirements

### Requirement: Implementation uses a dedicated worktree
The developer MUST perform implementation changes inside a dedicated repository-local worktree under `.worktrees/` and on a matching task feature branch.

#### Scenario: Task worktree is created
- **WHEN** a developer starts an implementation task numbered `<number>` with kebab-case name `<name>`
- **THEN** the active worktree MUST be `.worktrees/inventory-dashboard-task-<number>-<name>` and the branch MUST be `feature/task-<number>-<name>`

#### Scenario: Root checkout remains untouched by implementation
- **WHEN** the task worktree exists
- **THEN** the developer MUST NOT edit implementation files in the repository root checkout

### Requirement: Worktree identity is reported
The developer MUST verify and report the active absolute worktree path and branch before editing and in the final handoff.

#### Scenario: Developer begins work
- **WHEN** implementation is about to start
- **THEN** the developer MUST run `git worktree list --porcelain` and record the matching worktree and branch
