# Tasks: Azure CLI Account Context Switcher

**Input**: Design documents from `/specs/001-azure-context-switcher/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/cli-contract.md

**Tests**: Per azctx constitution (Principle V - Pragmatic Quality), this is a personal POC with no formal test suite. Manual testing is acceptable. Test tasks are NOT included per constitution.

**Organization**: Tasks are grouped by user story to enable independent implementation of each story in priority order.

## Format: `- [ ] [ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **azctx project**: `src/` at repository root (single project structure)
- No `tests/` directory - manual testing per constitution
- Python modules in `src/` with CLI entry point at `src/cli.py`
- Storage: `~/.azctx/contexts.yaml` (created at runtime)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and dependency setup

- [x] T001 Add dependencies via UV: `uv add typer rich questionary pyyaml`
- [x] T002 [P] Create `src/models/__init__.py` (empty module initializer)
- [x] T003 [P] Create `src/services/__init__.py` (empty module initializer)
- [x] T004 [P] Create `src/utils/__init__.py` (empty module initializer)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T005 Create Context dataclass in `src/models/context.py` with all 8 fields per data-model.md
- [x] T006 Add `to_dict()` and `from_dict()` methods to Context for YAML serialization in `src/models/context.py`
- [x] T007 Add `validate_context_id()` static method to Context with regex `^[a-zA-Z0-9_-]{1,20}$` in `src/models/context.py`
- [x] T008 [P] Create custom exception classes in `src/utils/errors.py`: AzureCliNotFoundError, NoActiveSessionError, DuplicateContextError, ContextNotFoundError, StorageError
- [x] T009 [P] Create Rich Console instance in `src/utils/console.py` and export for reuse
- [x] T010 [P] Add helper functions in `src/utils/console.py`: `print_success()`, `print_error()`, `print_warning()`, `print_info()`
- [x] T011 Create `get_storage_path()` function in `src/services/storage.py` returning `Path.home() / ".azctx" / "contexts.yaml"`
- [x] T012 Create `_ensure_storage_dir()` helper in `src/services/storage.py` to create `.azctx` directory if needed
- [x] T013 Create `load_contexts()` function in `src/services/storage.py` to read YAML and return list of Context objects
- [x] T014 Create `save_contexts(contexts)` function in `src/services/storage.py` to write list of Context objects to YAML
- [x] T015 Add error handling for corrupted YAML in `load_contexts()` in `src/services/storage.py`
- [x] T016 Create `check_azure_cli_installed()` function in `src/services/azure_cli.py` using subprocess
- [x] T017 Create `get_current_account()` function in `src/services/azure_cli.py` executing `az account show --output json`
- [x] T018 Create `set_account(subscription_id)` function in `src/services/azure_cli.py` executing `az account set --subscription <id>`
- [x] T019 Add timeout handling (5 seconds) for subprocess calls in `src/services/azure_cli.py`
- [x] T020 Add JSON parsing and error handling in `get_current_account()` in `src/services/azure_cli.py`
- [x] T021 Initialize Typer app in `src/cli.py`: `app = typer.Typer(help="Azure CLI context switcher")`
- [x] T022 Add `if __name__ == "__main__": app()` block to `src/cli.py` for direct execution

**Checkpoint**: Foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - Switch Between Saved Contexts (Priority: P1) 🎯 MVP

**Goal**: Enable interactive switching between saved Azure CLI contexts with keyboard navigation

**Independent Test**: Add at least 2 contexts manually to YAML file, run `uv run python -m src.cli switch`, navigate with arrow keys, select a context, verify `az account show` reflects the switch

**Manual Verification**:

1. Create test contexts in `~/.azctx/contexts.yaml`
2. Run `uv run python -m src.cli switch`
3. Use arrow keys to navigate list
4. Press Enter to select
5. Verify success message displays
6. Run `az account show` to confirm switch

### Implementation for User Story 1

- [x] T023 [US1] Create `get_context_by_id(context_id)` helper in `src/services/storage.py` to find context by ID
- [x] T024 [US1] Create `switch_context_interactive()` function in `src/services/context_manager.py`
- [x] T025 [US1] In `switch_context_interactive()`: Load contexts using `storage.load_contexts()`
- [x] T026 [US1] In `switch_context_interactive()`: Handle empty contexts list - return error message
- [x] T027 [US1] In `switch_context_interactive()`: Handle single context - inform user and exit
- [x] T028 [US1] In `switch_context_interactive()`: Use `questionary.select()` to display context list with format "ID - Name"
- [x] T029 [US1] In `switch_context_interactive()`: Call `azure_cli.set_account()` with selected subscription_id
- [x] T030 [US1] In `switch_context_interactive()`: Verify switch with `azure_cli.get_current_account()`
- [x] T031 [US1] In `switch_context_interactive()`: Return success with context details
- [x] T032 [US1] Add `@app.command()` decorator for `switch` command in `src/cli.py`
- [x] T033 [US1] In `switch` command: Add docstring "Switch to a different saved context" in `src/cli.py`
- [x] T034 [US1] In `switch` command: Call `context_manager.switch_context_interactive()`
- [x] T035 [US1] In `switch` command: Use Rich Panel to display success message (green) with context details
- [x] T036 [US1] In `switch` command: Handle KeyboardInterrupt and exit with code 130
- [x] T037 [US1] In `switch` command: Handle errors and display with `console.print_error()`
- [x] T038 [US1] In `switch` command: Set appropriate exit codes (0=success, 1=error, 130=cancelled)

**Checkpoint**: User Story 1 complete - can switch between contexts interactively

---

## Phase 4: User Story 2 - Add New Context (Priority: P2)

**Goal**: Enable users to save the current Azure CLI context with user-assigned identifiers

**Independent Test**: Login to Azure with `az login`, run `uv run python -m src.cli add`, provide name and ID when prompted, verify context saved to YAML file

**Manual Verification**:

1. Run `az login` and authenticate
2. Run `uv run python -m src.cli add`
3. Enter context name when prompted
4. Enter context ID when prompted
5. Verify success message
6. Check `~/.azctx/contexts.yaml` contains new context
7. Run `uv run python -m src.cli list` to confirm

### Implementation for User Story 2

- [x] T039 [US2] Create `add_context(context)` function in `src/services/storage.py` to append context and save YAML
- [x] T040 [US2] Create `context_id_exists(context_id)` function in `src/services/storage.py` to check uniqueness
- [x] T041 [US2] Create `add_context_interactive()` function in `src/services/context_manager.py`
- [x] T042 [US2] In `add_context_interactive()`: Check Azure CLI installed with `azure_cli.check_azure_cli_installed()`
- [x] T043 [US2] In `add_context_interactive()`: Get current account with `azure_cli.get_current_account()`
- [x] T044 [US2] In `add_context_interactive()`: Handle no active session - return error
- [x] T045 [US2] In `add_context_interactive()`: Use `questionary.text()` to prompt for context_name
- [x] T046 [US2] In `add_context_interactive()`: Validate context_name length (1-100 chars)
- [x] T047 [US2] In `add_context_interactive()`: Use `questionary.text()` to prompt for context_id
- [x] T048 [US2] In `add_context_interactive()`: Validate context_id with `Context.validate_context_id()`
- [x] T049 [US2] In `add_context_interactive()`: Check uniqueness with `storage.context_id_exists()`
- [x] T050 [US2] In `add_context_interactive()`: Re-prompt if validation fails
- [x] T051 [US2] In `add_context_interactive()`: Create Context object with Azure data and user input
- [x] T052 [US2] In `add_context_interactive()`: Set created_at to `datetime.now()`
- [x] T053 [US2] In `add_context_interactive()`: Use tenant_id for tenant_name (simple approach per data-model.md)
- [x] T054 [US2] In `add_context_interactive()`: Call `storage.add_context()` to save
- [x] T055 [US2] In `add_context_interactive()`: Return success with context details
- [x] T056 [US2] Add `@app.command()` decorator for `add` command in `src/cli.py`
- [x] T057 [US2] In `add` command: Add docstring "Add current Azure CLI context with friendly identifiers"
- [x] T058 [US2] In `add` command: Call `context_manager.add_context_interactive()`
- [x] T059 [US2] In `add` command: Use Rich Panel to display success message (green) with context details
- [x] T060 [US2] In `add` command: Handle KeyboardInterrupt and exit with code 130
- [x] T061 [US2] In `add` command: Handle AzureCliNotFoundError with helpful message
- [x] T062 [US2] In `add` command: Handle NoActiveSessionError with "run az login first" message
- [x] T063 [US2] In `add` command: Handle DuplicateContextError with re-prompt suggestion
- [x] T064 [US2] In `add` command: Set appropriate exit codes

**Checkpoint**: User Stories 1 AND 2 complete - can add and switch contexts

---

## Phase 5: User Story 3 - Check Active Context (Priority: P3)

**Goal**: Display currently active Azure context with friendly names if managed

**Independent Test**: Run `uv run python -m src.cli status`, verify it shows active context details

**Manual Verification**:

1. Switch to a saved context using `azctx switch`
2. Run `uv run python -m src.cli status`
3. Verify it shows friendly name and managed context marker
4. Run `az account set --subscription <other-id>` manually
5. Run `uv run python -m src.cli status` again
6. Verify it shows unmanaged context warning

### Implementation for User Story 3

- [x] T065 [US3] Create `get_status(verbose)` function in `src/services/context_manager.py`
- [x] T066 [US3] In `get_status()`: Get current account with `azure_cli.get_current_account()`
- [x] T067 [US3] In `get_status()`: Handle no active session - return error message
- [x] T068 [US3] In `get_status()`: Load saved contexts with `storage.load_contexts()`
- [x] T069 [US3] In `get_status()`: Match current subscription_id against saved contexts
- [x] T070 [US3] In `get_status()`: If match found, return managed context data with friendly names
- [x] T071 [US3] In `get_status()`: If no match, return unmanaged context data with suggestion to run `azctx add`
- [x] T072 [US3] In `get_status()`: Include verbose details if `verbose=True` (subscription_id, tenant_id, created_at)
- [x] T073 [US3] Add `@app.command()` decorator for `status` command in `src/cli.py`
- [x] T074 [US3] In `status` command: Add docstring "Show currently active Azure context"
- [x] T075 [US3] In `status` command: Add `--verbose` flag using `typer.Option(False, "--verbose", "-v")`
- [x] T076 [US3] In `status` command: Call `context_manager.get_status(verbose)`
- [x] T077 [US3] In `status` command: Use Rich Panel with green border for managed contexts (✓ symbol)
- [x] T078 [US3] In `status` command: Use Rich Panel with yellow border for unmanaged contexts (⚠ symbol)
- [x] T079 [US3] In `status` command: Display context_id and context_name for managed contexts
- [x] T080 [US3] In `status` command: Display subscription and username for all contexts
- [x] T081 [US3] In `status` command: Add verbose details (IDs, created_at) if flag set
- [x] T082 [US3] In `status` command: Handle NoActiveSessionError with red panel (✗ symbol)
- [x] T083 [US3] In `status` command: Set exit code 1 for no session, 0 for active session

**Checkpoint**: User Stories 1, 2, AND 3 complete - full context lifecycle with status

---

## Phase 6: User Story 4 - List All Contexts (Priority: P4)

**Goal**: Display all saved contexts in tabular format with optional verbose details

**Independent Test**: Add multiple contexts, run `uv run python -m src.cli list`, verify table displays; run with `--verbose` and verify all fields shown

**Manual Verification**:

1. Ensure multiple contexts saved (use `azctx add`)
2. Run `uv run python -m src.cli list`
3. Verify simple table with ID and Name columns
4. Run `uv run python -m src.cli list --verbose`
5. Verify detailed information for each context
6. Delete all contexts and run `azctx list`
7. Verify empty state message

### Implementation for User Story 4

- [x] T084 [US4] Create `list_contexts(verbose)` function in `src/services/context_manager.py`
- [x] T085 [US4] In `list_contexts()`: Load contexts with `storage.load_contexts()`
- [x] T086 [US4] In `list_contexts()`: Handle empty list - return empty state message
- [x] T087 [US4] In `list_contexts()`: If not verbose, return list of (context_id, context_name) tuples
- [x] T088 [US4] In `list_contexts()`: If verbose, return full Context objects
- [x] T089 [US4] Add `@app.command()` decorator for `list` command in `src/cli.py`
- [x] T090 [US4] In `list` command: Add docstring "List all saved contexts"
- [x] T091 [US4] In `list` command: Add `--verbose` flag using `typer.Option(False, "--verbose", "-v")`
- [x] T092 [US4] In `list` command: Call `context_manager.list_contexts(verbose)`
- [x] T093 [US4] In `list` command: For empty list, display "No saved contexts" message with suggestion to run `azctx add`
- [x] T094 [US4] In `list` command: For non-verbose, create Rich Table with columns: "ID", "Name"
- [x] T095 [US4] In `list` command: For non-verbose, add table rows from tuples
- [x] T096 [US4] In `list` command: For non-verbose, add footer text "Use 'azctx list --verbose' for detailed information"
- [x] T097 [US4] In `list` command: For verbose, display each context in separate Rich Panel
- [x] T098 [US4] In `list` command: For verbose, include: Name, Subscription (name + ID), Tenant ID, Username, Created timestamp
- [x] T099 [US4] In `list` command: Use cyan color for table headers and panel titles

**Checkpoint**: User Stories 1-4 complete - full read operations (list, status, switch)

---

## Phase 7: User Story 5 - Delete Context (Priority: P5)

**Goal**: Interactively select and delete saved contexts with confirmation

**Independent Test**: Run `uv run python -m src.cli delete`, select context from list, confirm deletion, verify context removed from YAML

**Manual Verification**:

1. Ensure multiple contexts saved
2. Run `uv run python -m src.cli delete`
3. Navigate list with arrow keys
4. Select context to delete
5. Confirm when prompted
6. Verify success message
7. Run `azctx list` to confirm deletion
8. Try cancelling with Esc, verify no deletion

### Implementation for User Story 5

- [x] T100 [US5] Create `delete_context(context_id)` function in `src/services/storage.py` to remove and save
- [x] T101 [US5] Create `delete_context_interactive()` function in `src/services/context_manager.py`
- [x] T102 [US5] In `delete_context_interactive()`: Load contexts with `storage.load_contexts()`
- [x] T103 [US5] In `delete_context_interactive()`: Handle empty list - return error message
- [x] T104 [US5] In `delete_context_interactive()`: Use `questionary.select()` to display context list for deletion
- [x] T105 [US5] In `delete_context_interactive()`: Format list items as "ID - Name"
- [x] T106 [US5] In `delete_context_interactive()`: Use `questionary.confirm()` to prompt "Are you sure you want to delete context 'X'?"
- [x] T107 [US5] In `delete_context_interactive()`: If not confirmed, return cancelled message
- [x] T108 [US5] In `delete_context_interactive()`: Call `storage.delete_context(context_id)` if confirmed
- [x] T109 [US5] In `delete_context_interactive()`: Return success with deleted context_id
- [x] T110 [US5] In `delete_context_interactive()`: Add note that Azure CLI session remains active
- [x] T111 [US5] Add `@app.command()` decorator for `delete` command in `src/cli.py`
- [x] T112 [US5] In `delete` command: Add docstring "Delete a saved context"
- [x] T113 [US5] In `delete` command: Call `context_manager.delete_context_interactive()`
- [x] T114 [US5] In `delete` command: Display success message with context_id
- [x] T115 [US5] In `delete` command: Display note about Azure session remaining active
- [x] T116 [US5] In `delete` command: Handle KeyboardInterrupt and display "Deletion cancelled"
- [x] T117 [US5] In `delete` command: Handle ContextNotFoundError
- [x] T118 [US5] In `delete` command: Set exit code 130 for cancel, 0 for success

**Checkpoint**: All 5 user stories complete - full CRUD operations on contexts

---

## Phase 9: User Story 6 - Build and Distribute Standalone Executable (Priority: P6)

**Purpose**: Enable distribution as standalone Windows executable without requiring Python installation

**User Story**: A developer wants to distribute azctx as a standalone Windows executable that can be dropped into any PATH directory for immediate use, without requiring Python installation.

**Independent Test**: Build .exe locally, place in PATH, run commands from any location without Python installed. Verify GitHub Actions creates releases on version tags.

**Acceptance Criteria**:

1. Local build command creates `azctx.exe` in dist folder
2. Version tag push triggers automated GitHub Actions build
3. GitHub Release auto-created with versioned zip file
4. Downloaded .exe works from any PATH location
5. Executable behaves identically to `uv run python -m src.cli`

### Setup Tasks

- [x] T143 [P] [US6] Add PyInstaller to dev dependencies in `pyproject.toml`: `uv add --dev pyinstaller`
- [x] T144 [P] [US6] Add `dist/` to `.gitignore` for PyInstaller output directory
- [x] T145 [P] [US6] Add `build/` to `.gitignore` for PyInstaller build artifacts
- [x] T146 [P] [US6] Add `*.spec` to `.gitignore` for PyInstaller spec files (if not using custom spec)

### Local Build Implementation

- [x] T147 [US6] Test local build command: `uv run pyinstaller --onefile src/cli.py --name azctx`
- [x] T148 [US6] Verify `dist/azctx.exe` is created and size is under 50MB (SC-011)
- [x] T149 [US6] Test executable functionality: `.\dist\azctx.exe --help` works correctly
- [x] T150 [US6] Test executable commands: Run `.\dist\azctx.exe status` and verify output matches Python version
- [x] T151 [US6] Test executable from different directory to verify no path dependencies

### GitHub Actions Workflow

- [x] T152 [P] [US6] Create `.github/workflows/` directory if it doesn't exist
- [x] T153 [US6] Create `.github/workflows/release.yml` with trigger `on: push: tags: 'v*'`
- [x] T154 [US6] Add checkout step using `actions/checkout@v4` in workflow
- [x] T155 [US6] Add Python setup step using `actions/setup-python@v5` with version 3.11 in workflow
- [x] T156 [US6] Add UV installation step: `pip install uv` in workflow
- [x] T157 [US6] Add dependency installation step: `uv sync --all-extras --dev` in workflow
- [x] T158 [US6] Add PyInstaller build step: `uv run pyinstaller --onefile src/cli.py --name azctx` in workflow
- [x] T159 [US6] Add size verification step with PowerShell script to fail if executable > 50MB in workflow
- [x] T160 [US6] Add packaging step to create `azctx-${{ github.ref_name }}-windows.zip` using Compress-Archive in workflow
- [x] T161 [US6] Add release creation step using `softprops/action-gh-release@v1` with zip file attachment in workflow

### Documentation

- [x] T162 [P] [US6] Add "Building Executable" section to README.md with local build instructions
- [x] T163 [P] [US6] Add "Installation from Release" section to README.md with download and PATH setup instructions
- [x] T164 [P] [US6] Add "Creating Releases" section to README.md explaining version tag workflow
- [x] T165 [P] [US6] Document troubleshooting for common PyInstaller issues in README.md (hidden imports, size optimization)
- [x] T166 [P] [US6] Update README.md installation section to list both .exe and source installation options

### Testing & Validation

**Note**: These tasks require manual execution as they involve GitHub repository operations with version tags.

- [ ] T167 [US6] Create test version tag (e.g., v0.1.0-test): `git tag -a v0.1.0-test -m "Test release workflow"`
- [ ] T168 [US6] Push test tag and verify GitHub Actions workflow triggers: `git push origin v0.1.0-test`
- [ ] T169 [US6] Monitor GitHub Actions logs to ensure build completes within 10 minutes (SC-010)
- [ ] T170 [US6] Verify GitHub Release was auto-created with correct zip filename
- [ ] T171 [US6] Download release zip and verify it contains `azctx.exe`
- [ ] T172 [US6] Extract `azctx.exe` and verify size is under 50MB (SC-011)
- [ ] T173 [US6] Test downloaded executable on clean Windows machine without Python installed (if possible)
- [ ] T174 [US6] Verify executable commands work: `azctx --help`, `azctx status`, `azctx list`
- [ ] T175 [US6] Compare executable behavior with Python version to ensure identical functionality (FR-025)
- [ ] T176 [US6] Clean up test release and tag if workflow successful

**Manual Testing Steps**:

1. **Ensure repository is pushed to GitHub**:
   ```bash
   git add .
   git commit -m "Implement User Story 6: Executable distribution"
   git push origin 001-azure-context-switcher
   ```

2. **Create and push test version tag**:
   ```bash
   git tag -a v0.1.0-test -m "Test release workflow"
   git push origin v0.1.0-test
   ```

3. **Monitor GitHub Actions**:
   - Go to your repository on GitHub
   - Navigate to **Actions** tab
   - Click on the "Build and Release" workflow run
   - Verify all steps complete successfully
   - Check total runtime is under 10 minutes

4. **Verify Release**:
   - Navigate to **Releases** section
   - Verify `v0.1.0-test` release was created
   - Download `azctx-v0.1.0-test-windows.zip`
   - Extract and verify `azctx.exe` is present
   - Check file size: should be under 50MB (typically 20-30MB)

5. **Test Executable**:
   ```bash
   # Test from extracted location
   .\azctx.exe --help
   .\azctx.exe status
   .\azctx.exe list

   # Compare with Python version
   uv run python -m src.cli --help
   uv run python -m src.cli status
   ```

6. **Clean Up Test Release** (optional):
   - Delete the test release from GitHub Releases UI
   - Delete local tag: `git tag -d v0.1.0-test`
   - Delete remote tag: `git push origin :refs/tags/v0.1.0-test`

**Parallel Opportunities**: Tasks T143-T146 (setup), T162-T166 (documentation) can run in parallel with workflow creation

**Time Estimate**: ~1.5-2.5 hours (1-2 hours implementation + 30 min testing)

**Outcome**: Standalone Windows executable distributed via GitHub Releases, automated on version tags

**Checkpoint**: User Story 6 complete - executable distribution enabled

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories and overall UX

- [x] T119 [P] Add detailed docstrings to all public functions in `src/services/` modules
- [x] T120 [P] Add type hints to all function signatures across the codebase
- [x] T121 [P] Update README.md with project description, installation instructions (UV), and usage examples
- [x] T122 [P] Add usage examples to README.md for all 5 commands (add, switch, list, status, delete)
- [x] T123 [P] Document YAML storage location in README.md
- [x] T124 Verify all error messages are clear and actionable per cli-contract.md
- [x] T125 Ensure all success messages use green color and ✓ symbol
- [x] T126 Ensure all error messages use red color and ✗ symbol
- [x] T127 Ensure warning messages use yellow color and ⚠ symbol
- [x] T128 Add help text validation - run `uv run python -m src.cli --help` and verify output
- [x] T129 Add help text validation - run `uv run python -m src.cli [command] --help` for each command
- [ ] T130 Run complete manual test checklist from quickstart.md
- [ ] T131 Test edge case: Azure CLI not installed (mock by renaming `az` temporarily)
- [ ] T132 Test edge case: No active Azure session
- [ ] T133 Test edge case: Corrupted YAML file (manually edit to invalid syntax)
- [ ] T134 Test edge case: Duplicate context_id on add
- [ ] T135 Test edge case: Empty contexts list for switch/delete
- [ ] T136 Test edge case: Single context for switch
- [x] T137 Measure command performance - all should complete < 2 seconds (excluding user input)
- [x] T138 [P] Add project.scripts entry in pyproject.toml: `azctx = "src.cli:app"`
- [ ] T139 Test installation: `uv pip install -e .` and verify `azctx` command works
- [x] T140 Code cleanup: Remove any debug print statements
- [x] T141 Code cleanup: Ensure consistent error handling patterns
- [x] T142 Add import statements organization (group stdlib, third-party, local)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup (Phase 1) completion - BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational (Phase 2) completion
- **User Story 2 (Phase 4)**: Depends on Foundational (Phase 2) completion - Can run in parallel with US1
- **User Story 3 (Phase 5)**: Depends on Foundational (Phase 2) completion - Can run in parallel with US1, US2
- **User Story 4 (Phase 6)**: Depends on Foundational (Phase 2) completion - Can run in parallel with US1, US2, US3
- **User Story 5 (Phase 7)**: Depends on Foundational (Phase 2) completion - Can run in parallel with US1, US2, US3, US4
- **User Story 6 (Phase 9)**: Depends on US1-US5 completion (needs working CLI to build) - Can run in parallel with Polish
- **Polish (Phase 8)**: Depends on completion of desired user stories (typically all)

### User Story Dependencies

**Critical Path**: Setup → Foundational → User Stories (any order) → Distribution (US6)

- **User Story 1 (P1)**: No dependencies on other user stories - Can start immediately after Foundational
- **User Story 2 (P2)**: No dependencies on other user stories - Independent (but required for US1 to be useful)
- **User Story 3 (P3)**: No dependencies on other user stories - Independent
- **User Story 4 (P4)**: No dependencies on other user stories - Independent
- **User Story 5 (P5)**: No dependencies on other user stories - Independent
- **User Story 6 (P6)**: Depends on US1-US5 completion (requires working CLI) - Builds the application for distribution

**Note**: While technically independent, practical workflow is: Add contexts (US2) → Switch (US1) → Status/List (US3/US4) → Delete (US5) → Build distribution (US6)

### Within Each Phase

**Foundational (Phase 2)**:

- T005-T007 (Context model) can run in sequence
- T008 (errors), T009-T010 (console) can run in parallel with each other
- T011-T015 (storage) must run in sequence
- T016-T020 (Azure CLI) must run in sequence
- T021-T022 (CLI initialization) must wait for T008-T010
- Storage and Azure CLI modules can be built in parallel

**User Story Phases**: Within each user story, tasks generally run sequentially due to dependencies (models → services → CLI commands)

### Parallel Opportunities

- **Phase 1 (Setup)**: All tasks T002, T003, T004 can run in parallel
- **Phase 2 (Foundational)**: Tasks T008, T009 (utils) can run in parallel with storage/Azure CLI development
- **Phase 3-7 (User Stories)**: Different user stories can be implemented in parallel by different developers
- **Phase 8 (Polish)**: Most tasks marked [P] can run in parallel (documentation, testing)
- **Phase 9 (US6)**: Setup tasks (T143-T146) and documentation tasks (T162-T166) can run in parallel with workflow creation

**Single Developer**: Recommend sequential execution in priority order: P1 → P2 → P3 → P4 → P5 → P6

**Team of 2-3**: After Foundational phase, split user stories (e.g., Dev1: US1+US2, Dev2: US3+US4+US5, Dev3: US6)

---

## Parallel Example: Foundational Phase

```bash
# After Setup (Phase 1) completes, these can start in parallel:

Group A (Context Model):
- T005: Create Context dataclass
- T006: Add to_dict/from_dict methods
- T007: Add validate_context_id method

Group B (Utilities - parallel with Group A):
- T008: Create custom exceptions
- T009: Create Rich Console instance
- T010: Add console helper functions

Group C (Storage - starts after Group A completes):
- T011: get_storage_path()
- T012: _ensure_storage_dir()
- T013: load_contexts()
- T014: save_contexts()
- T015: Error handling

Group D (Azure CLI - parallel with Group C):
- T016: check_azure_cli_installed()
- T017: get_current_account()
- T018: set_account()
- T019: Timeout handling
- T020: JSON parsing

Group E (CLI Init - after Group B completes):
- T021: Initialize Typer app
- T022: Add main block
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

**Goal**: Get to a working prototype as fast as possible

**Scope**: Phase 1 (Setup) + Phase 2 (Foundational) + Phase 3 (User Story 1)

**Time Estimate**: ~3-4 hours

**Outcome**: Working `azctx switch` command (requires manual YAML file creation for testing)

### MVP + Essential (User Stories 1 & 2)

**Goal**: Complete the core workflow

**Scope**: MVP + Phase 4 (User Story 2)

**Time Estimate**: ~4-5 hours total

**Outcome**: Can add contexts and switch between them - fully self-sufficient tool

### Full Feature (All User Stories 1-5)

**Goal**: Complete implementation with all core features

**Scope**: All phases (1-8)

**Time Estimate**: ~6-7 hours total

**Outcome**: Production-ready tool with add, switch, status, list, delete capabilities

### Full Feature + Distribution (All User Stories 1-6)

**Goal**: Complete implementation with distribution capability

**Scope**: All phases (1-9) including executable builds

**Time Estimate**: ~8.5-10 hours total

**Outcome**: Production-ready tool with GitHub Release distribution via standalone .exe

### Incremental Delivery

1. **Sprint 1**: MVP (US1) - Deliver switchable contexts
2. **Sprint 2**: MVP + US2 - Deliver self-sufficient tool
3. **Sprint 3**: US3 + US4 - Add visibility features
4. **Sprint 4**: US5 + Polish - Complete feature set
5. **Sprint 5**: US6 - Add executable distribution (optional, for wider adoption)

Each sprint delivers independently testable value.

---

## Task Summary

- **Total Tasks**: 176 (34 new tasks for US6)
- **Setup Phase**: 4 tasks
- **Foundational Phase**: 18 tasks (blocking)
- **User Story 1 (P1)**: 16 tasks
- **User Story 2 (P2)**: 26 tasks
- **User Story 3 (P3)**: 19 tasks
- **User Story 4 (P4)**: 16 tasks
- **User Story 5 (P5)**: 19 tasks
- **User Story 6 (P6)**: 34 tasks
- **Polish Phase**: 24 tasks

**Parallel Tasks**: 21 tasks marked [P] (can execute in parallel)

**Estimated Total Time**: 8.5-10 hours for complete implementation (including US6)

**MVP Time** (US1 only): ~3-4 hours
**Core Workflow** (US1+US2): ~4-5 hours
**With Distribution** (US1-US6): ~8.5-10 hours
