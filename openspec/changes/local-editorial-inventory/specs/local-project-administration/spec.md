## ADDED Requirements

### Requirement: Manage projects and curated media
The system SHALL allow the local single user to create and edit projects and associate approved local thumbnails, originals, screenshots, and videos without authentication.

#### Scenario: Create a project
- **WHEN** the user submits a valid project form
- **THEN** the project is persisted locally and appears in the catalog

#### Scenario: Add curated media
- **WHEN** the user associates an approved local file with a project
- **THEN** it is stored within the local media boundary and shown on the detail page
