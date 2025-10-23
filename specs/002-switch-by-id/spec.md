# Feature Specification: Direct Context Switching by ID

**Feature Branch**: `002-switch-by-id`
**Created**: October 23, 2025
**Status**: Draft
**Input**: User description: "azctx switch should support parameter --id (with -i shorthand) followed by user-assigned context id to enable quick, non-interactive switch. Case Sensitive ID. For example: azctx switch --id DEV switches to Context which user assigned id DEV. Add standard error handling for non-existing ids (in feedback message give actual existing ids)"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Quick Non-Interactive Switch by ID (Priority: P1)

A developer knows the ID of the context they want to switch to and wants to switch directly without entering interactive mode, making scripting and automation possible.

**Why this priority**: This is the core value of the feature - enabling fast, scriptable context switches. Power users who know their context IDs can switch in a single command, significantly faster than interactive selection.

**Independent Test**: Can be fully tested by saving contexts with known IDs and running `azctx switch --id <ID>`. Delivers immediate value by reducing switch time from ~5 seconds (interactive) to under 2 seconds (direct).

**Acceptance Scenarios**:

1. **Given** I have saved a context with ID "DEV", **When** I run `azctx switch --id DEV`, **Then** the Azure CLI switches to that context and displays a success confirmation
2. **Given** I have saved a context with ID "PROD", **When** I run `azctx switch -i PROD`, **Then** the Azure CLI switches to that context using the short-form flag
3. **Given** multiple contexts exist, **When** I run `azctx switch --id <valid-id>`, **Then** the switch completes in under 2 seconds
4. **Given** I am already on the target context, **When** I run `azctx switch --id <current-id>`, **Then** I receive a message indicating the context is already active

---

### User Story 2 - Case-Sensitive ID Matching (Priority: P2)

A developer uses IDs with specific casing (e.g., "DEV", "dev", "Dev" as different contexts) and expects the system to respect the exact case when switching.

**Why this priority**: Critical for accuracy and preventing accidental switches, but slightly lower priority than basic functionality. Ensures users can organize contexts with meaningful case distinctions.

**Independent Test**: Can be tested by creating contexts with IDs that differ only in casing (e.g., "DEV" and "dev") and verifying each switches correctly. Delivers value by supporting sophisticated naming schemes.

**Acceptance Scenarios**:

1. **Given** I have contexts with IDs "DEV" and "dev", **When** I run `azctx switch --id DEV`, **Then** only the "DEV" context is activated (not "dev")
2. **Given** I have a context with ID "Staging", **When** I run `azctx switch --id staging`, **Then** I receive an error indicating the ID was not found
3. **Given** I have a context with ID "PROD", **When** I run `azctx switch --id prod`, **Then** I receive an error with the correct available ID shown ("PROD")

---

### User Story 3 - Error Handling with Available IDs (Priority: P3)

A developer mistypes or forgets an ID and needs helpful feedback showing what IDs actually exist.

**Why this priority**: Improves user experience and reduces frustration, but not critical for core functionality. Users who know their IDs won't need this often.

**Independent Test**: Can be tested by attempting to switch to non-existent IDs and verifying the error message lists all available IDs. Delivers value by serving as a quick reference without running a separate list command.

**Acceptance Scenarios**:

1. **Given** I have saved contexts with IDs "DEV", "TEST", "PROD", **When** I run `azctx switch --id STAGING`, **Then** I receive an error message stating "Context 'STAGING' not found" followed by "Available contexts: DEV, TEST, PROD"
2. **Given** I have only one saved context with ID "DEV", **When** I run `azctx switch --id PROD`, **Then** the error message shows "Available contexts: DEV"
3. **Given** I have no saved contexts, **When** I run `azctx switch --id DEV`, **Then** I receive a message indicating no contexts have been registered and instructions to add one
4. **Given** I have 10+ saved contexts, **When** I provide an invalid ID, **Then** all available context IDs are listed in alphabetical order for easy scanning

---

### Edge Cases

- What happens when the --id parameter is provided but no ID value follows (e.g., `azctx switch --id`)?
- How does the system handle both --id and interactive mode being requested simultaneously?
- What happens when special characters or spaces are used in the ID parameter (e.g., `azctx switch --id "DEV 1"`)?
- How does the system handle extremely long ID values (100+ characters)?
- What happens if the context storage is corrupted and IDs cannot be retrieved?
- How does the system behave when Azure CLI fails during the switch operation?
- What happens when the specified context exists in storage but the Azure account no longer exists?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST support `--id` flag for the switch command to enable direct context switching by ID
- **FR-002**: System MUST support `-i` as a shorthand alias for `--id`
- **FR-003**: System MUST perform case-sensitive matching when comparing provided ID to saved context IDs
- **FR-004**: System MUST display a success confirmation message when switching via ID completes successfully
- **FR-005**: System MUST display an error message when the provided ID does not match any saved context
- **FR-006**: Error messages MUST include the list of all available context IDs to help users correct their input
- **FR-007**: The list of available IDs in error messages MUST be sorted alphabetically for easy scanning
- **FR-008**: System MUST handle the case where no ID value is provided after `--id` flag with a clear error message
- **FR-009**: System MUST complete ID-based switching within 2 seconds on standard hardware
- **FR-010**: When switching to an already-active context via ID, system MUST inform the user that the context is already active
- **FR-011**: When no contexts are registered, error message MUST guide users to add contexts first
- **FR-012**: System MUST validate that the ID parameter is not empty or whitespace-only
- **FR-013**: System MUST handle Azure CLI failures during ID-based switching with appropriate error messages
- **FR-014**: System MUST trim leading/trailing whitespace from the provided ID before matching
- **FR-015**: When both `--id` and interactive mode are applicable, system MUST prioritize the `--id` parameter

### Key Entities

This feature extends the existing **Context** entity (no new entities):

- **Context**: Already contains the unique ID field that will be matched against the `--id` parameter

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can switch contexts via ID in under 2 seconds from command invocation to confirmation (faster than interactive mode's ~5 seconds)
- **SC-002**: 100% of valid context IDs successfully switch when provided via `--id` parameter
- **SC-003**: Error messages for invalid IDs include the complete list of available IDs in 100% of cases
- **SC-004**: Users can successfully script context switches using the `--id` parameter without requiring interactive input
- **SC-005**: Case-sensitive matching correctly distinguishes between IDs differing only in case 100% of the time
- **SC-006**: Command-line help text clearly documents both `--id` and `-i` flags with examples

## Assumptions

- Users are already familiar with the existing `azctx switch` interactive command
- Users know or can easily reference their context IDs (either from memory or previous `azctx list` output)
- Context IDs are unique and already enforced by the existing add command
- Azure CLI is installed and accessible (inherited from base tool requirements)
- Users may want to use this feature in scripts or automation workflows
- Most users have fewer than 20 contexts, making the "available IDs" list readable in error messages
- The existing context storage mechanism is reliable and accessible
- Users prefer speed and directness over interactive selection when they know the target ID
