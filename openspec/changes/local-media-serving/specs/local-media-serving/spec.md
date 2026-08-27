## ADDED Requirements

### Requirement: Serve managed curated media
The system SHALL serve explicitly curated local media referenced by a managed
relative path, and SHALL keep the asset inside the configured local media root.

#### Scenario: Serve an approved image
- **WHEN** a browser requests an existing curated image through its managed media path
- **THEN** the local service returns the binary with an appropriate image content type

#### Scenario: Serve an approved video
- **WHEN** a browser requests an existing curated video through its managed media path
- **THEN** the local service returns the binary with an appropriate video content type

#### Scenario: Missing managed asset
- **WHEN** a browser requests a managed media path whose file does not exist
- **THEN** the service returns HTTP 404 without exposing an absolute filesystem path

### Requirement: Enforce media path safety
The system MUST reject media paths that are absolute, traverse outside the managed
root, resolve through an escaping symlink, or use an unsupported media type.

#### Scenario: Path traversal attempt
- **WHEN** a client submits or requests a path containing traversal outside the managed root
- **THEN** the service rejects the operation and does not read or persist the target file

#### Scenario: Unsupported media type
- **WHEN** a client submits a media asset with an unsupported extension or content type
- **THEN** the service returns a validation error and does not persist the asset reference

### Requirement: Preserve local media across restart
The system SHALL persist the managed relative media reference in SQLite and SHALL
serve the corresponding existing binary after the local API restarts.

#### Scenario: Restart with an existing asset
- **WHEN** the API restarts using the same local database and media root
- **THEN** the recorded media remains associated with its project and remains requestable

### Requirement: Keep media local by default
The system SHALL bind media serving to the local runtime by default and SHALL
not fetch or synchronize media from remote repositories.

#### Scenario: Local-only media runtime
- **WHEN** the application starts with default configuration
- **THEN** media is served only through the local loopback runtime and no remote import is performed
