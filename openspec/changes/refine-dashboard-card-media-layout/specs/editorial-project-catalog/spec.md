## MODIFIED Requirements

### Requirement: Dashboard and catalog discovery
The system SHALL provide a dashboard summary and desktop project catalog with editorial cards and an optional compact table. Each catalog card SHALL show the project name, a compact truncated description, a thumbnail when available, recorded status without inventing values, and technology pills containing names only; embedded technology URLs SHALL NOT be rendered in the card.

#### Scenario: Browse compact project cards
- **WHEN** the user opens the catalog
- **THEN** each project card shows its name, thumbnail when available, recorded status, and description limited to the card's compact line count

#### Scenario: Scan card technologies
- **WHEN** a project has technology metadata containing names and URLs
- **THEN** the card displays name-only technology pills and does not display the embedded URLs

### Requirement: Project detail presentation
The system SHALL show recorded project metadata, technologies, agentic platform, curated screenshots, demo video references, and public GitHub/Trello links. The single representative Graphify Report image SHALL render at the full primary content width with its intrinsic aspect ratio preserved, matching the demo video's content presence; full technology and platform links SHALL remain available on the detail page.

#### Scenario: View full-width report media
- **WHEN** the user opens a project with one representative Graphify Report image
- **THEN** the image spans the primary content width, preserves its intrinsic aspect ratio, and remains selectable for the existing full-size preview behavior

#### Scenario: Preserve detail metadata links
- **WHEN** the user opens a project detail page
- **THEN** the complete technology and agentic-platform names and associated links remain available there even though catalog cards use name-only pills

#### Scenario: Use a narrow viewport
- **WHEN** the user opens catalog cards or project detail on a supported narrow viewport
- **THEN** descriptions remain compact, pills remain readable, and report media remains contained without clipping or horizontal overflow
