## MODIFIED Requirements

### Requirement: Project detail presentation
The system SHALL show recorded project metadata, technologies, agentic platform, curated screenshots, demo video references, and public GitHub/Trello links. When dedicated repository references are present, the detail page SHALL show the GitHub repository and Trello backlog as separate actionable links and SHALL omit absent values.

#### Scenario: View project detail
- **WHEN** the user opens a project
- **THEN** the system presents recorded fields and omits empty fields

#### Scenario: Open dedicated repository references
- **WHEN** the project has a saved GitHub repository URL or Trello backlog URL
- **THEN** the detail page shows each present URL as a distinct keyboard-accessible HTTPS link with its destination preserved and does not fetch remote content

#### Scenario: Enlarge an image
- **WHEN** the user selects a thumbnail
- **THEN** the associated original or full-size local image is displayed
