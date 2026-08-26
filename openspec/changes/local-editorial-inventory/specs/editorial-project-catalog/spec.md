## ADDED Requirements

### Requirement: Dashboard and catalog discovery
The system SHALL provide a dashboard summary and desktop project catalog with editorial cards and an optional compact table.

#### Scenario: Browse projects
- **WHEN** the user opens the catalog
- **THEN** the system shows each project with name, short description, thumbnail when available, and recorded status without inventing values

### Requirement: Project detail presentation
The system SHALL show recorded project metadata, technologies, agentic platform, curated screenshots, demo video references, and public GitHub/Trello links.

#### Scenario: View project detail
- **WHEN** the user opens a project
- **THEN** the system presents recorded fields and omits empty fields

#### Scenario: Enlarge an image
- **WHEN** the user selects a thumbnail
- **THEN** the associated original or full-size local image is displayed
