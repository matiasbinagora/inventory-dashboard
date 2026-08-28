## MODIFIED Requirements

### Requirement: Manage projects and curated media
The system SHALL allow the local single user to create and edit projects, associate approved local thumbnails, originals, screenshots, and videos, and optionally record one HTTPS GitHub repository URL and one HTTPS Trello backlog URL per project without authentication. Empty optional URL fields SHALL remain absent and the system SHALL NOT contact the referenced services.

#### Scenario: Create a project
- **WHEN** the user submits a valid project form with zero or one valid GitHub repository URL and zero or one valid Trello backlog URL
- **THEN** the project and supplied references are persisted locally and the project appears in the catalog

#### Scenario: Edit project repository references
- **WHEN** the user saves changed optional GitHub or Trello URLs for an existing project
- **THEN** the existing values are replaced atomically, omitted values are cleared, and the saved project returns the new values

#### Scenario: Reject invalid repository references
- **WHEN** the user submits a non-HTTPS URL, unsupported host, malformed URL, URL with userinfo, or URL without a path in either dedicated field
- **THEN** the API rejects the project change without persisting the invalid value and the administration UI shows a controlled validation error

#### Scenario: Add curated media
- **WHEN** the user associates an approved local file with a project
- **THEN** it is stored within the local media boundary and shown on the detail page
