## ADDED Requirements

### Requirement: Enforce privacy boundaries
The system SHALL store only safe user-facing metadata and explicitly curated media, excluding source code, credentials, private URLs, transcripts, and uncurated artifacts.

#### Scenario: Exclude private content
- **WHEN** private content is encountered during curation
- **THEN** the system does not copy, index, expose, or persist it

### Requirement: Operate locally by default
The system SHALL bind local services to `127.0.0.1` by default and SHALL not require authentication for the initial local runtime.

#### Scenario: Start locally
- **WHEN** the user starts the application with defaults
- **THEN** services listen only on the local host and administration is locally available
