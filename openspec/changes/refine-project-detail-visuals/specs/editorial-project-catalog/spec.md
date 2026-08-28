## MODIFIED Requirements

### Requirement: Present project detail
The system SHALL present a selected project's identity, description, Graphify Report media, technologies, public links, agentic platform, demo video, and manual history in a readable responsive detail view. The project title SHALL use a compact technical hierarchy, the description SHALL use a wider bounded reading measure, technologies SHALL be presented as a semantic bullet list, links and agentic platform references SHALL be visually distinct and keyboard accessible, the media section SHALL be labeled `Graphify Report`, and the video section SHALL be labeled `Demo Video`.

#### Scenario: Read the project identity and description
- **WHEN** the user opens a project detail page on a supported desktop viewport
- **THEN** the project title is prominent but does not overwhelm the page, and the description is presented in a wider readable measure with stable line spacing

#### Scenario: Scan the Graphify report
- **WHEN** the selected project has curated image media
- **THEN** the media section is labeled `Graphify Report` and the existing image preview/lightbox behavior remains available

#### Scenario: Scan technologies
- **WHEN** the selected project has technologies
- **THEN** the technologies are displayed as individual semantic bullet items that can be read separately

#### Scenario: Use public links
- **WHEN** the selected project has public links
- **THEN** each link is presented as a distinct readable, keyboard-accessible anchor with its existing destination preserved

#### Scenario: Read the agentic platform
- **WHEN** the selected project has an agentic platform value containing public references
- **THEN** the platform information is displayed as a distinct readable block and its public references remain actionable anchors

#### Scenario: Watch the demo video
- **WHEN** the selected project has local or public video media
- **THEN** the video section is labeled `Demo Video` and the existing local controls or public-link behavior remains available

#### Scenario: Use a narrow viewport
- **WHEN** the user opens the project detail page on a supported narrow viewport
- **THEN** the title, description, lists, links, media, and demo video remain readable without unintended overlap, clipping, or horizontal overflow
