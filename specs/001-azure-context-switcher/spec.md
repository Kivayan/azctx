# Feature Specification: Azure CLI Account Context Switcher

**Feature Branch**: `001-azure-context-switcher`
**Created**: October 17, 2025
**Status**: Draft
**Input**: User description: "I want to create a basic CLI tool that is very fast to use and convenient, which is used for switching Azure CLI accounts. I'm working on multiple different Azure CLI accounts and switching between them using Azure CLI is tedious and inconvenient. I want to have a tool in which I am able to register my Azure CLI account context, assign them friendly names, and then using interactive CLI to be able to switch between them, manage them, and check active one as well."

## Clarifications

### Session 2025-10-17

- Q: Which tool should be used to build standalone .exe files? → A: PyInstaller
- Q: When should the GitHub Actions workflow trigger to build the .exe? → A: Trigger on version tags (e.g., v1.0.0, v1.1.0) - Build only for releases
- Q: How should the built .exe be distributed? → A: Auto-create GitHub Release with .exe attached - Standard distribution method
- Q: What should the distribution package filename and format be? → A: Package as `azctx-{version}-windows.zip` containing `azctx.exe`
- Q: What command should developers use to build the .exe locally? → A: uv run pyinstaller --onefile src/cli.py

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Switch Between Saved Contexts (Priority: P1)

A developer working on multiple Azure projects needs to quickly switch between different Azure CLI account contexts without remembering account IDs or typing lengthy commands.

**Why this priority**: This is the core value proposition - fast context switching. Without this, there's no reason to use the tool.

**Independent Test**: Can be fully tested by saving at least two contexts and switching between them. Delivers immediate value by eliminating the need to manually run `az account set` commands.

**Acceptance Scenarios**:

1. **Given** I have saved multiple Azure contexts, **When** I run the switch command, **Then** I see an interactive list of all saved contexts with their friendly names
2. **Given** the interactive list is displayed, **When** I select a context using keyboard navigation, **Then** the selected Azure account becomes active in Azure CLI
3. **Given** I select a context to switch to, **When** the switch completes, **Then** I receive confirmation that the context is now active
4. **Given** I have only one saved context, **When** I attempt to switch, **Then** I am informed that there's only one context available

---

### User Story 2 - Add New Context (Priority: P2)

A developer who is currently logged into an Azure account via Azure CLI wants to save this context with a friendly name for easy future access.

**Why this priority**: Essential for populating the tool with contexts to switch between, but switching is more valuable since users need to add contexts only once.

**Independent Test**: Can be tested by logging into Azure CLI and running the add command. Delivers value by allowing users to organize their contexts with meaningful names.

**Acceptance Scenarios**:

1. **Given** I am logged into an Azure account via Azure CLI, **When** I run the add command, **Then** I am prompted to provide a friendly name and unique ID for the current context
2. **Given** I provide a friendly name and ID, **When** the context is saved, **Then** I receive confirmation that the context was successfully registered
3. **Given** I try to add a context with an ID that already exists, **When** I submit, **Then** I receive an error message and am asked to choose a different ID
4. **Given** I am not logged into any Azure account, **When** I attempt to add a context, **Then** I receive an error message indicating no active Azure session was found
5. **Given** I provide a friendly name with special characters, **When** the context is saved, **Then** the name is properly stored and displayed in future operations

---

### User Story 3 - Check Active Context (Priority: P3)

A developer needs to quickly verify which Azure account context is currently active before performing operations.

**Why this priority**: Helpful for confirmation but less critical than switching and adding contexts. Users can also verify via Azure CLI directly.

**Independent Test**: Can be tested by running the status command after switching contexts. Delivers value by providing quick confirmation without needing to remember Azure CLI commands.

**Acceptance Scenarios**:

1. **Given** I have an active Azure context that was set via this tool, **When** I run the status command, **Then** I see the friendly name and ID of the currently active context
2. **Given** I have an active Azure context that was not set via this tool, **When** I run the status command, **Then** I see the raw Azure account information with a note that it's not managed by this tool
3. **Given** no Azure context is active, **When** I run the status command, **Then** I receive a message indicating no active Azure session

---

### User Story 4 - List All Contexts (Priority: P4)

A developer wants to see all saved contexts with basic or detailed information to review what's available.

**Why this priority**: Useful for management but not critical for core workflow. The switch command already shows available contexts.

**Independent Test**: Can be tested by adding multiple contexts and running the list command. Delivers value by providing an overview of all saved contexts without entering interactive mode.

**Acceptance Scenarios**:

1. **Given** I have multiple saved contexts, **When** I run the list command without verbose flag, **Then** I see a simple list of IDs and friendly names
2. **Given** I have multiple saved contexts, **When** I run the list command with verbose flag, **Then** I see detailed information including Azure subscription IDs, tenant IDs, and account names
3. **Given** I have no saved contexts, **When** I run the list command, **Then** I receive a message indicating no contexts have been registered

---

### User Story 5 - Delete Context (Priority: P5)

A developer wants to remove contexts that are no longer needed from the saved list.

**Why this priority**: Maintenance feature that's only needed occasionally. Lowest priority as it doesn't directly support the main workflow.

**Independent Test**: Can be tested by adding a context and then deleting it. Delivers value by keeping the context list clean and manageable.

**Acceptance Scenarios**:

1. **Given** I have saved contexts, **When** I run the delete command with a context ID, **Then** the specified context is removed and I receive confirmation
2. **Given** I specify an ID that doesn't exist, **When** I run the delete command, **Then** I receive an error message indicating the context was not found
3. **Given** I delete a context that is currently active, **When** the deletion completes, **Then** the Azure CLI context remains active but is no longer managed by this tool
4. **Given** I have multiple contexts, **When** I delete one, **Then** the remaining contexts are unaffected

---

### User Story 6 - Build and Distribute Standalone Executable (Priority: P6)

A developer wants to distribute azctx as a standalone Windows executable that can be dropped into any PATH directory for immediate use, without requiring Python installation.

**Why this priority**: Deployment convenience feature that enables wider adoption. Users can install by simply downloading and placing the .exe file, but core functionality works fine via Python/UV for development.

**Independent Test**: Can be tested by building the .exe locally with PyInstaller, placing it in a PATH directory, and running `azctx` commands from any location. Also verified through GitHub Actions automated builds on version tags.

**Acceptance Scenarios**:

1. **Given** I am a developer, **When** I run the local build command, **Then** a standalone `azctx.exe` is created in the dist folder
2. **Given** a version tag is pushed to GitHub (e.g., v1.0.0), **When** the GitHub Actions workflow runs, **Then** a versioned zip file (`azctx-v1.0.0-windows.zip`) containing `azctx.exe` is automatically created
3. **Given** a successful build completes on a version tag, **When** the workflow finishes, **Then** a GitHub Release is automatically created with the zip file attached
4. **Given** I download the release zip, **When** I extract `azctx.exe` to a directory in my PATH, **Then** I can run `azctx` commands from any terminal location without Python installed
5. **Given** the .exe is built, **When** I run it, **Then** it behaves identically to running via `uv run python -m src.cli`

---

### Edge Cases

- What happens when the Azure CLI is not installed or not in the system PATH?
- How does the system handle corrupted context storage data?
- What happens when a saved context references an Azure account that no longer exists or has been revoked?
- How does the system handle concurrent usage (multiple terminal sessions attempting to switch contexts simultaneously)?
- What happens when friendly names or IDs contain unusual characters or are extremely long?
- How does the system handle contexts saved on different machines (if storage is synced)?
- What happens when Azure CLI version changes and the account structure changes?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide a command to add the current Azure CLI context with a user-provided unique ID and friendly name
- **FR-002**: System MUST provide an interactive command to switch between saved contexts using keyboard navigation
- **FR-003**: System MUST provide a command to display the currently active Azure context
- **FR-004**: System MUST provide a command to list all saved contexts with basic information (ID and friendly name)
- **FR-005**: System MUST provide a command to list all saved contexts with detailed information (verbose mode including subscription ID, tenant ID, account name)
- **FR-006**: System MUST provide a command to delete a saved context by its ID
- **FR-007**: System MUST validate that each context ID is unique when adding new contexts
- **FR-008**: System MUST persist saved contexts between sessions
- **FR-009**: System MUST detect when Azure CLI is not installed or not accessible and provide a clear error message
- **FR-010**: System MUST detect when no Azure account is logged in and provide appropriate error messages
- **FR-011**: System MUST retrieve current Azure context using `az account show` command
- **FR-012**: System MUST set Azure context using `az account set` command
- **FR-013**: System MUST complete all operations within 2 seconds on standard hardware (prioritize speed)
- **FR-014**: Interactive switching MUST support keyboard-only navigation (up/down arrows, enter to select, escape to cancel)
- **FR-015**: System MUST display clear confirmation messages for all successful operations
- **FR-016**: System MUST display clear error messages for all failure scenarios
- **FR-017**: When adding a context, system MUST extract and store Azure subscription ID, tenant ID, and account name from `az account show`
- **FR-018**: System MUST prevent deletion of contexts while displaying appropriate confirmation
- **FR-019**: System MUST gracefully handle cases where stored contexts reference non-existent Azure accounts
- **FR-020**: System MUST support building a standalone Windows executable using PyInstaller with the command `uv run pyinstaller --onefile src/cli.py`
- **FR-021**: The standalone executable MUST be named `azctx.exe` for easy PATH installation
- **FR-022**: GitHub Actions workflow MUST trigger on version tag pushes (e.g., v1.0.0, v1.1.0) to build the executable
- **FR-023**: GitHub Actions workflow MUST package the `azctx.exe` into a versioned zip file named `azctx-{version}-windows.zip`
- **FR-024**: GitHub Actions workflow MUST automatically create a GitHub Release with the zip file attached as a downloadable asset
- **FR-025**: The standalone executable MUST function identically to running via `uv run python -m src.cli` without requiring Python installation on the target machine
- **FR-026**: Local build instructions MUST be documented for developers who want to build the executable on their own machines

### Key Entities

- **Context**: Represents a saved Azure CLI account configuration, containing:
  - Unique ID (user-provided short identifier)
  - Friendly name (user-provided descriptive name)
  - Azure subscription ID
  - Azure tenant ID
  - Azure account name/email
  - Timestamp of when context was added

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can switch between Azure contexts in under 5 seconds from command invocation to confirmation
- **SC-002**: Users can add a new context in under 10 seconds
- **SC-003**: The tool launches and displays interactive selection within 1 second
- **SC-004**: All commands execute and complete within 2 seconds (excluding network latency from Azure CLI)
- **SC-005**: Users can successfully navigate and select from a list of 50+ saved contexts without performance degradation
- **SC-006**: 100% of valid Azure CLI contexts can be imported and switched to successfully
- **SC-007**: Users require fewer than 5 keystrokes to switch contexts (command shortcut + 1-2 keys for selection)
- **SC-008**: Error messages clearly identify the problem in 100% of failure cases
- **SC-009**: Standalone executable builds successfully on Windows using PyInstaller without errors
- **SC-010**: GitHub Actions workflow completes and creates releases within 10 minutes of version tag push
- **SC-011**: The standalone `azctx.exe` is under 50MB in size for fast download and distribution

## Assumptions

- Users already have Azure CLI installed and configured
- Users are familiar with basic terminal/command-line usage
- Users have already authenticated to Azure using `az login` before attempting to add contexts
- Context storage will be local to the machine (no cloud sync)
- Standard terminal capabilities are available (color support, arrow key navigation)
- Azure CLI commands (`az account show`, `az account set`) maintain backward compatibility
- Users will manage contexts from a single machine (no multi-machine synchronization required)
- Context switching only affects the Azure CLI global context, not application-specific contexts
- Users need to switch contexts at most a few times per hour (not continuous rapid switching)

## Dependencies

- Azure CLI must be installed and accessible in system PATH
- Users must have valid Azure accounts and active subscriptions
- Terminal must support interactive input (not in non-interactive/CI environments)
- PyInstaller must be available for building standalone executables
- GitHub repository must have Actions enabled for automated builds
- Version tags must follow semantic versioning format (e.g., v1.0.0) for release automation

## Scope

### In Scope

- Adding current Azure CLI context with friendly identifiers
- Interactive switching between saved contexts
- Listing saved contexts (basic and verbose)
- Viewing currently active context
- Deleting saved contexts
- Local persistence of context metadata
- Keyboard navigation for context selection
- Basic validation and error handling
- Building standalone Windows executable using PyInstaller
- Local build instructions for developers
- Automated GitHub Actions workflow for building executables on version tags
- Automated GitHub Release creation with executable distribution
- Packaging executable in versioned zip files

### Out of Scope

- Authentication or login to Azure (handled by Azure CLI)
- Managing Azure resources or subscriptions
- Cloud synchronization of saved contexts
- GUI or web interface
- Context switching for other cloud providers (AWS, GCP)
- Automatic context switching based on directory or project
- Context grouping or tagging
- Export/import of context configurations
- Multi-tenancy or team sharing of contexts
- Integration with IDEs or other tools
- Historical tracking of context switches
- Backup and restore of contexts
