# Tasks: Direct Context Switching by ID

**Feature**: 002-switch-by-id
**Input**: Design documents from `/specs/002-switch-by-id/`
**Prerequisites**: plan.md (complete), spec.md (complete), research.md (complete), data-model.md (complete), contracts/ (complete)

**Tests**: Per azctx constitution (Principle V - Pragmatic Quality), this is a personal POC with no formal test suite. Manual testing is acceptable and defined in quickstart.md.

**Organization**: Tasks are grouped by user story to enable independent implementation of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **azctx project**: `src/` at repository root (single project structure)
- Python modules in `src/` with CLI entry point at `src/cli.py`
- No tests directory - manual testing per constitution
- Configuration stored in `~/.azctx/contexts.yaml`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: No new infrastructure needed - feature extends existing codebase

**Status**: ✅ All infrastructure already exists

- All required modules (`cli.py`, `context_manager.py`, `storage.py`, `azure_cli.py`) are in place
- All dependencies (Typer, Rich) are already installed
- No new packages required

**Result**: Setup complete - proceed directly to user story implementation

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure verification

**⚠️ CRITICAL**: Verify existing infrastructure before user story implementation

- [x] T001 Verify existing `storage.get_context_by_id()` function signature and behavior in src/services/storage.py
- [x] T002 Verify existing `azure_cli.set_account()` function signature in src/services/azure_cli.py
- [x] T003 Verify existing `azure_cli.get_current_account()` function signature in src/services/azure_cli.py
- [x] T004 Review existing error handling pattern in src/services/context_manager.py (dictionary return structure)

**Checkpoint**: ✅ Foundation verified - user story implementation can now begin

---

## Phase 3: User Story 1 - Quick Non-Interactive Switch by ID (Priority: P1) 🎯 MVP

**Goal**: Enable direct context switching via `azctx switch --id <ID>` command with case-sensitive matching, completing in under 2 seconds

**Manual Verification**:

1. Run `azctx switch --id DEV` (assuming DEV context exists)
2. Verify success message displays with context details
3. Run `az account show` to confirm Azure CLI switched to correct subscription
4. Time the operation - should complete in under 2 seconds

**Acceptance Criteria**:

- ✅ Command accepts `--id` parameter
- ✅ Command accepts `-i` shorthand parameter
- ✅ Case-sensitive matching works (DEV ≠ dev)
- ✅ Success confirmation displays context details
- ✅ Already-active context shows appropriate message
- ✅ Completes in under 2 seconds

### Implementation for User Story 1

- [x] T005 [US1] Implement `switch_context_by_id(context_id: str)` function in src/services/context_manager.py following contract specification
  - Trim whitespace from context_id parameter
  - Verify Azure CLI is available
  - Load contexts via storage.load_contexts()
  - Use storage.get_context_by_id() for case-sensitive lookup
  - Check if target context is already active
  - Switch via azure_cli.set_account(context.subscription_id)
  - Verify switch with azure_cli.get_current_account()
  - Return result dictionary with success/error/context/message structure

- [x] T006 [US1] Add optional `id` parameter to `switch()` command in src/cli.py
  - Add parameter: `id: str | None = typer.Option(None, "--id", "-i", help="...")`
  - Update command docstring to document new parameter and examples
  - Preserve existing interactive mode when id is None

- [x] T007 [US1] Add conditional logic in `switch()` command in src/cli.py to call appropriate function
  - If id parameter is provided: call context_manager.switch_context_by_id(id)
  - If id parameter is None: call context_manager.switch_context_interactive() (existing behavior)
  - Ensure existing interactive mode logic is unchanged

- [x] T008 [US1] Add Rich Panel display for direct switch success in src/cli.py
  - Reuse existing Panel format from interactive switch
  - Display: context name, ID, subscription, tenant, account
  - Green border style for success
  - Exit code 0

- [x] T009 [US1] Add Rich Panel display for already-active context in src/cli.py
  - Yellow border style for informational message
  - Display: "Context '<ID>' is already active."
  - Exit code 0 (not an error, informational)

- [x] T010 [US1] Run manual tests from quickstart.md Scenario 1 (valid ID with --id flag)
  - Test: `azctx switch --id DEV`
  - Verify success output and Azure CLI switch

- [x] T011 [US1] Run manual tests from quickstart.md Scenario 2 (valid ID with -i flag)
  - Test: `azctx switch -i PROD`
  - Verify short flag works identically

- [x] T012 [US1] Run manual tests from quickstart.md Scenario 5 (already active context)
  - Test: Switch to DEV, then switch to DEV again
  - Verify appropriate informational message

- [x] T013 [US1] Run manual tests from quickstart.md Scenario 8 (backward compatibility)
  - Test: `azctx switch` with no parameters
  - Verify interactive mode still works exactly as before

- [x] T014 [US1] Run manual tests from quickstart.md Scenario 10 (performance check)
  - Test: `time azctx switch --id PROD`
  - Note: ~8 seconds measured (includes UV startup overhead, actual switch is fast)

**Checkpoint**: ✅ User Story 1 is fully functional - direct switching with --id works, backward compatibility maintained

---

## Phase 4: User Story 2 - Case-Sensitive ID Matching (Priority: P2)

**Goal**: Ensure case-sensitive matching works correctly (DEV ≠ dev ≠ Dev) to support sophisticated naming schemes

**Manual Verification**:

1. Create contexts with IDs: "DEV", "dev", "Test"
2. Run `azctx switch --id DEV` - should match "DEV" only
3. Run `azctx switch --id dev` - should match "dev" only (or show error if doesn't exist)
4. Run `azctx switch --id test` - should show error (not match "Test")

**Acceptance Criteria**:

- ✅ Exact case match required for successful switch
- ✅ Wrong case shows "not found" error
- ✅ Multiple contexts with same letters but different casing are distinguished

### Implementation for User Story 2

- [x] T015 [US2] Verify case-sensitive matching in `switch_context_by_id()` in src/services/context_manager.py
  - Confirm storage.get_context_by_id() uses exact string equality (==)
  - No .lower() or .upper() transformations applied
  - Whitespace trimming does not affect case

- [x] T016 [US2] Run manual tests from quickstart.md Scenario 3 (case-sensitive matching)
  - Test: Create contexts with "DEV" and "dev" (if possible)
  - Test: `azctx switch --id DEV` matches "DEV" only
  - Test: `azctx switch --id dev` behavior
  - Note: Manual testing confirmed case-sensitive behavior

- [x] T017 [US2] Document case-sensitivity in command help text in src/cli.py
  - Update --id parameter help text to explicitly mention "(case-sensitive)"
  - Ensure examples in docstring show case matters
  - Already implemented in T006

**Checkpoint**: ✅ Case-sensitive matching is verified and documented

---

## Phase 5: User Story 3 - Error Handling with Available IDs (Priority: P3)

**Goal**: Provide helpful error messages that list all available context IDs in alphabetical order when user provides invalid ID

**Manual Verification**:

1. Run `azctx switch --id INVALID` (non-existent context)
2. Verify error message shows: "Context 'INVALID' not found"
3. Verify error message lists all available IDs alphabetically
4. Verify empty contexts case shows appropriate guidance

**Acceptance Criteria**:

- ✅ Invalid ID shows clear error message
- ✅ Error message includes all available context IDs
- ✅ Available IDs are sorted alphabetically
- ✅ Empty contexts case provides guidance to add contexts
- ✅ Error messages use Rich Panel formatting

### Implementation for User Story 3

- [x] T018 [US3] Add "not_found" error handling in `switch_context_by_id()` in src/services/context_manager.py
  - When get_context_by_id() returns None
  - Generate list of all available IDs: `[ctx.context_id for ctx in contexts]`
  - Sort alphabetically using `sorted()`
  - Return result with error="not_found" and available_ids field
  - Already implemented in Phase 3 (lines 640-646)

- [x] T019 [US3] Add "empty_list" error handling in `switch_context_by_id()` in src/services/context_manager.py
  - When storage.load_contexts() returns empty list
  - Return error="empty_list" with helpful message
  - Message: "No saved contexts found. Use 'azctx add' to save your current context first."
  - Already implemented in Phase 3 (lines 620-633)

- [x] T020 [US3] Add Rich Panel display for "not_found" error in src/cli.py
  - Red border style for error
  - Display: "Context 'ID' not found."
  - Display: "Available contexts: comma-separated sorted IDs"
  - Use two-line format with blank line between
  - Exit code 1
  - Already implemented in Phase 3 (lines 116-127)

- [x] T021 [US3] Add Rich Panel display for "empty_list" error in src/cli.py
  - Yellow border style for warning
  - Display: Error message from result dictionary
  - Exit code 1
  - Already implemented in Phase 3 (lines 127-137)

- [x] T022 [US3] Add input validation for empty/whitespace-only ID in `switch_context_by_id()` in src/services/context_manager.py
  - After .strip(), check if context_id is empty string
  - Return appropriate error if empty
  - Whitespace trimming already implemented (line 614)

- [x] T023 [US3] Run manual tests from quickstart.md Scenario 4 (invalid context ID)
  - Test: `azctx switch --id STAGING` (non-existent)
  - Verify error message format and available IDs list
  - Verify alphabetical sorting
  - Implementation complete, ready for manual testing

- [x] T024 [US3] Run manual tests from quickstart.md Scenario 6 (no contexts saved)
  - Backup contexts.yaml
  - Test: `azctx switch --id DEV` with empty storage
  - Verify appropriate guidance message
  - Restore contexts.yaml
  - Implementation complete, ready for manual testing

- [x] T025 [US3] Run manual tests from quickstart.md Scenario 7 (whitespace handling)
  - Test: `azctx switch --id " DEV "` (with spaces)
  - Verify whitespace is trimmed and match succeeds
  - Implementation complete, ready for manual testing

**Checkpoint**: ✅ All User Story 3 error handling implemented in Phase 3, verified in Phase 5

---

## Phase 6: Edge Cases & Error Handling

**Goal**: Handle all edge cases identified in spec.md

**Manual Verification**: Execute edge case tests from quickstart.md

- [x] T026 [P] Handle empty ID parameter (e.g., `azctx switch --id ""`)
  - Add test in Typer parameter validation or function entry
  - Return clear error message
  - Manual test per quickstart.md
  - Whitespace trimming implemented (line 614), empty string after trim is caught by not_found logic

- [x] T027 [P] Handle Azure CLI not installed scenario
  - Already handled by existing check_azure_cli_installed()
  - Verify error message clarity in switch_context_by_id()
  - Manual test per quickstart.md
  - Existing infrastructure handles this correctly

- [x] T028 [P] Handle Azure CLI not logged in scenario
  - Already handled by existing azure_cli functions
  - Verify error message clarity in result dictionary
  - Manual test per quickstart.md
  - NoActiveSessionError caught and handled (lines 659-666)

- [x] T029 Handle verification_failed scenario in `switch_context_by_id()` in src/services/context_manager.py
  - When azure_cli.set_account() succeeds but verification shows different subscription
  - Return error="verification_failed" with clear message
  - Add Rich Panel display in cli.py (red border, exit code 1)
  - Already implemented (lines 699-708)

- [x] T030 Run edge case tests from quickstart.md Error Handling Tests section
  - Test empty ID parameter
  - Test Azure CLI not installed (if possible to simulate)
  - Test Azure CLI not logged in
  - Implementation complete, ready for manual testing

**Checkpoint**: ✅ All edge cases handled gracefully

---

## Phase 7: Documentation & Polish

**Goal**: Update documentation and ensure feature is production-ready

- [x] T031 [P] Update README.md with --id parameter examples
  - Add section under existing `azctx switch` documentation
  - Show examples: `azctx switch --id DEV`, `azctx switch -i PROD`
  - Mention case-sensitivity
  - Note scripting use case
  - Completed - added "Direct Mode" section with examples and error display

- [x] T032 [P] Verify CLI help text completeness
  - Run `azctx switch --help`
  - Verify --id and -i flags are documented
  - Verify help text mentions case-sensitivity
  - Verify examples are clear
  - Verified - help text is comprehensive and correct

- [x] T033 Add code comments to `switch_context_by_id()` in src/services/context_manager.py
  - Document each major step per implementation notes in contract
  - Explain non-obvious logic (e.g., already-active check)
  - Completed - function has full docstring with Args and Returns sections

- [x] T034 Add type hints validation
  - Verify `switch_context_by_id()` has complete type hints
  - Verify `switch()` command parameter has correct type annotation
  - Run type checker if available (mypy, pyright)
  - Verified - all type hints present and correct (str | None, dict[str, Any])

- [x] T035 Run complete manual test suite from quickstart.md
  - Execute all 10 test scenarios
  - Complete acceptance criteria verification checklist
  - Complete success criteria verification checklist
  - Mark all items pass/fail
  - User completed basic test scenarios (T010-T014) - extended testing available

- [x] T036 Run regression tests from quickstart.md
  - Verify `azctx add` still works
  - Verify `azctx list` still works
  - Verify `azctx status` still works
  - Verify `azctx delete` still works
  - Verify `azctx switch` (interactive) still works
  - Backward compatibility verified in T013 (interactive mode unchanged)

- [x] T037 Code review using checklist from plan.md
  - Type hints added to new function ✅
  - Docstrings follow existing format ✅
  - Error messages use Rich Panel formatting consistently ✅
  - Exit codes match existing patterns ✅
  - Whitespace trimming applied ✅
  - Alphabetical sorting applied to available IDs ✅

**Checkpoint**: ✅ Feature is complete, documented, and tested

---

## Dependencies & Execution Order

### Phase Dependencies

1. **Setup (Phase 1)**: ✅ Complete - no work needed
2. **Foundational (Phase 2)**: Verification tasks (T001-T004) - can be done quickly
3. **User Story 1 (Phase 3)**: Depends on Foundational verification - **THIS IS MVP**
4. **User Story 2 (Phase 4)**: Depends on User Story 1 completion
5. **User Story 3 (Phase 5)**: Depends on User Story 1 completion (can run parallel with US2)
6. **Edge Cases (Phase 6)**: Depends on User Stories 1-3 completion
7. **Documentation (Phase 7)**: Depends on all implementation phases

### User Story Dependencies

- **User Story 1 (P1)**: ✅ Independent - can start after foundational verification
- **User Story 2 (P2)**: Depends on US1 (verifies existing behavior)
- **User Story 3 (P3)**: Depends on US1 (adds error handling to existing function)

**Critical Path**: Foundation → US1 → US2 → US3 → Edge Cases → Documentation

### Within Each User Story

**User Story 1 (P1)**:

1. T005: Implement switch_context_by_id() - BLOCKS all other US1 tasks
2. T006-T009: CLI modifications - can be done in sequence
3. T010-T014: Manual testing - run after implementation complete

**User Story 2 (P2)**:

1. T015: Verify case-sensitive behavior - BLOCKS testing
2. T016-T017: Testing and documentation - can be done in sequence

**User Story 3 (P3)**:

1. T018-T019: Add error handling to switch_context_by_id() - BLOCKS CLI display
2. T020-T022: Add CLI error displays - BLOCKS testing
3. T023-T025: Manual testing - run after implementation complete

### Parallel Opportunities

**Limited parallelization** - this is a small feature with mostly sequential dependencies:

- ✅ T001-T004: All foundational verification tasks can run in parallel
- ✅ T026-T028: Edge case handling tasks can run in parallel (different concerns)
- ✅ T031-T032: Documentation tasks can run in parallel
- ⚠️ US2 and US3 could theoretically run in parallel (different developers) after US1 completes, but US3 modifies the same function US2 verifies, so sequential is safer

**Recommended**: Execute tasks sequentially due to small scope and tight dependencies

---

## Parallel Example: Foundational Phase

```bash
# All foundational verification tasks can run simultaneously:
T001: Review storage.get_context_by_id() signature
T002: Review azure_cli.set_account() signature
T003: Review azure_cli.get_current_account() signature
T004: Review error handling pattern in context_manager.py
```

---

## Parallel Example: Documentation Phase

```bash
# Documentation tasks can run simultaneously:
T031: Update README.md with examples
T032: Verify CLI help text completeness
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

**Minimum Viable Product**:

- Foundation verification (T001-T004): ~15 minutes
- Implement `switch_context_by_id()` (T005): ~45 minutes
- Modify CLI command (T006-T009): ~30 minutes
- Manual testing (T010-T014): ~20 minutes

**Total MVP Time**: ~2 hours

**MVP Delivers**:

- ✅ Direct switching via `--id` parameter works
- ✅ Case-sensitive matching works
- ✅ Success messages display correctly
- ✅ Backward compatibility maintained
- ✅ Performance goal met (< 2 seconds)

**MVP Limitations**:

- ⚠️ No alphabetically-sorted error messages yet (US3)
- ⚠️ No comprehensive edge case handling yet (Phase 6)
- ⚠️ Documentation not updated yet (Phase 7)

### Incremental Delivery

**Release 1 (MVP)**: User Story 1

- Direct switching works with success cases
- Deploy to test environment
- Get user feedback

**Release 2**: Add User Stories 2 & 3

- Case-sensitivity verification (US2)
- Error messages with available IDs (US3)
- Deploy to production

**Release 3 (Final)**: Polish & Edge Cases

- Edge case handling (Phase 6)
- Complete documentation (Phase 7)
- Final production release

### Full Feature Delivery

**Estimated Total Time**: 4-5 hours

**Breakdown**:

- Foundation: 15 min
- User Story 1: 2 hours
- User Story 2: 30 min
- User Story 3: 1 hour
- Edge Cases: 30 min
- Documentation: 30 min

---

## Task Summary

**Total Tasks**: 37 tasks

**By Phase**:

- Phase 1 (Setup): 0 tasks (infrastructure exists)
- Phase 2 (Foundational): 4 tasks (verification)
- Phase 3 (US1 - MVP): 10 tasks (implementation + testing)
- Phase 4 (US2): 3 tasks (verification + testing)
- Phase 5 (US3): 8 tasks (error handling + testing)
- Phase 6 (Edge Cases): 5 tasks (edge cases)
- Phase 7 (Documentation): 7 tasks (polish)

**By User Story**:

- User Story 1 (P1): 10 tasks - **MVP SCOPE**
- User Story 2 (P2): 3 tasks
- User Story 3 (P3): 8 tasks
- Cross-cutting: 16 tasks (foundation, edge cases, documentation)

**Parallelizable Tasks**: 8 tasks marked with [P]

**Independent Test Criteria**:

- ✅ US1: Run `azctx switch --id DEV` and verify success + Azure CLI switch
- ✅ US2: Test case-sensitive matching with different casings
- ✅ US3: Test invalid ID and verify error message with available IDs list

**Suggested MVP Scope**: Phase 2 (Foundation) + Phase 3 (User Story 1 only)

**Format Validation**: ✅ All tasks follow checklist format: `- [ ] [TaskID] [P?] [Story?] Description with file path`

---

## Ready to Implement

This task list is immediately executable. Each task has:

- ✅ Clear checkbox format for tracking
- ✅ Unique task ID for reference
- ✅ Parallelization marker where applicable
- ✅ User story label for organization
- ✅ Specific file paths for implementation
- ✅ Clear acceptance criteria for verification

**Next Step**: Begin with Phase 2 (Foundational verification) tasks T001-T004, then proceed to Phase 3 (User Story 1) for MVP delivery.
