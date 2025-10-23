# Implementation Plan: Direct Context Switching by ID

**Branch**: `002-switch-by-id` | **Date**: October 23, 2025 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-switch-by-id/spec.md`

## Summary

Add `--id` (with `-i` shorthand) parameter to the existing `azctx switch` command to enable direct, non-interactive context switching. Users can switch to a specific context by providing its ID (case-sensitive) in a single command, making the tool faster for power users and enabling scriptable automation. When an invalid ID is provided, the system displays helpful error messages listing all available context IDs in alphabetical order.

## Technical Context

**Language/Version**: Python 3.11+ (managed via UV)

**Primary Dependencies**: Typer (CLI), Rich (output) - No new dependencies required

**Package Manager**: UV (all dependency and environment management)

**Testing**: Manual testing (no formal test suite per constitution)

**Target Platform**: Cross-platform (Windows, macOS, Linux)

**Project Type**: Single Python project with CLI interface

**Performance Goals**: Direct switching via ID must complete in under 2 seconds

**Constraints**: Simplicity first - extend existing `switch` command rather than creating new command

**Scale/Scope**: Personal POC - optimize for velocity and experimentation

**Existing Architecture**:

- `src/cli.py`: Contains `switch()` command - needs parameter addition
- `src/services/context_manager.py`: Contains `switch_context_interactive()` - needs new function for direct switching
- `src/services/storage.py`: Already has `get_context_by_id()` - ready to use
- `src/services/azure_cli.py`: Already has `set_account()` - ready to use
- Error handling pattern: Return dictionaries with success/error structure

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Principle I - Simplicity First**:

- [x] Feature uses the simplest possible implementation approach
  - Extends existing command with optional parameter
  - Reuses existing storage and Azure CLI functions
  - No new external dependencies
- [x] No premature optimization or over-engineering present
  - Simple parameter check: if ID provided, use direct switch; else use interactive
- [x] Code will be readable and self-documenting
  - Clear function names, straightforward control flow
- [x] Configuration is minimal and intuitive
  - No new configuration needed

**Principle II - Modular Architecture**:

- [x] Feature is designed as self-contained, modular component(s)
  - New function `switch_context_by_id()` in `context_manager.py`
  - Separated from interactive switch logic
- [x] Module responsibilities are clear and single-purpose
  - CLI layer: Parse parameters, handle display
  - Service layer: Business logic for ID-based switching
  - Storage layer: Already provides context retrieval
- [x] Dependencies between modules are explicit and minimal
  - Same dependency pattern as existing switch command
- [x] No tight coupling introduced
  - Uses existing storage and Azure CLI abstractions

**Principle III - CLI-Centric Design**:

- [x] All functionality accessible via CLI commands
  - `azctx switch --id DEV` and `azctx switch -i DEV`
- [x] Uses Typer for command structure
  - Add optional `id` parameter to existing `@app.command()` decorator
- [x] Uses Questionary for interactive prompts (if needed)
  - Not needed for this feature (non-interactive by design)
- [x] Uses Rich for terminal output formatting
  - Reuse existing Panel display for success/error messages
- [x] Commands have clear help text
  - Typer automatically includes parameter in `--help`
- [x] Supports both interactive and non-interactive modes (where applicable)
  - Preserves interactive mode (no `--id`), adds non-interactive mode (`--id` provided)

**Principle IV - UV-Managed Workflow**:

- [x] All new dependencies will be added via UV
  - No new dependencies required
- [x] No alternative package managers used
  - N/A - no new packages
- [x] Dependencies listed in pyproject.toml
  - All existing dependencies already listed

**Principle V - Pragmatic Quality**:

- [x] Feature includes documentation for non-obvious functionality
  - Update command docstring with `--id` parameter explanation
  - Update README with examples
- [x] Manual testing approach is defined
  - Test scenarios in quickstart.md
- [x] No formal test suite required (personal POC)
  - Consistent with project approach
- [x] Type hints used where they improve clarity
  - Add type hints to new function signature

**Result**: ✅ All constitution gates passed. No violations.

## Project Structure

### Documentation (this feature)

```plaintext
specs/002-switch-by-id/
├── plan.md              # This file
├── research.md          # Phase 0 output (minimal - no research needed)
├── data-model.md        # Phase 1 output (uses existing Context model)
├── quickstart.md        # Phase 1 output (test scenarios)
├── contracts/           # Phase 1 output (function contracts)
└── tasks.md             # Phase 2 output (via /speckit.tasks command)
```

### Source Code (repository root)

```plaintext
# Existing azctx structure - files to modify:
src/
├── cli.py                      # MODIFY: Add --id parameter to switch()
├── services/
│   ├── context_manager.py      # MODIFY: Add switch_context_by_id()
│   ├── storage.py              # USE: get_context_by_id() already exists
│   └── azure_cli.py            # USE: set_account() already exists
└── utils/
    └── console.py              # USE: Existing Rich console instance
```

**Structure Decision**: Minimal changes to existing structure. No new modules needed. Extends existing `cli.py` command and adds one new service function in `context_manager.py`.

## Complexity Tracking

**No violations** - All constitution gates passed without exceptions.

---

## Phase 0: Research

### Research Questions

Since this feature extends existing functionality using established patterns, minimal research is needed:

1. **Typer Optional Parameters**: How to add optional parameter to existing Typer command
2. **Case-Sensitive Matching**: Confirm Python string comparison is case-sensitive by default
3. **Alphabetical Sorting**: Standard approach for sorting context IDs

### Research Findings

**Decision 1: Typer Optional Parameter Pattern**

```python
@app.command()
def switch(
    id: str | None = typer.Option(None, "--id", "-i", help="Context ID to switch to")
) -> None:
    """Switch to a different saved context."""
    if id:
        # Direct switch by ID
        switch_by_id_logic(id)
    else:
        # Interactive switch (existing code)
        switch_interactive_logic()
```

**Rationale**: Typer's `Option` with `None` default preserves backward compatibility. Existing interactive mode works when `--id` is omitted.

**Decision 2: Case-Sensitive Matching**

Python string equality (`==`) is case-sensitive by default. No special handling needed.

**Decision 3: Error Message with Available IDs**

```python
available_ids = sorted([ctx.context_id for ctx in contexts])
error_msg = f"Context '{id}' not found.\nAvailable contexts: {', '.join(available_ids)}"
```

**Rationale**: Use built-in `sorted()` for alphabetical ordering. Join with commas for readability.

---

## Phase 1: Design & Contracts

### Data Model

**No new entities** - Feature uses existing `Context` model from `src/models/context.py`:

```python
@dataclass
class Context:
    context_id: str          # Used for --id parameter matching (case-sensitive)
    context_name: str        # Displayed in success messages
    subscription_id: str     # Used by azure_cli.set_account()
    subscription_name: str   # Displayed in success confirmation
    tenant_id: str
    tenant_name: str
    username: str
    created_at: datetime
```

**Key Field**: `context_id` - The unique identifier used for `--id` parameter matching.

### Function Contracts

#### New Function: `switch_context_by_id()`

**Location**: `src/services/context_manager.py`

**Signature**:

```python
def switch_context_by_id(context_id: str) -> dict[str, Any]:
    """Switch to a saved Azure CLI context by its ID (case-sensitive).

    Args:
        context_id: The case-sensitive ID of the context to switch to.
                    Leading/trailing whitespace will be trimmed.

    Returns:
        Dictionary with operation result:
        - success (bool): Whether the operation succeeded
        - message (str): Status message
        - context (Context | None): Selected context if successful
        - error (str | None): Error type if failed
        - available_ids (list[str] | None): Sorted list of available IDs (only on 'not_found' error)

    Raises:
        AzureCliNotFoundError: If Azure CLI is not installed.
    """
```

**Implementation Logic**:

1. Verify Azure CLI is available (reuse existing check)
2. Trim whitespace from `context_id` parameter
3. Load saved contexts via `storage.load_contexts()`
4. Check if list is empty → return error with guidance to add contexts
5. Use `storage.get_context_by_id(context_id)` for case-sensitive lookup
6. If not found:
   - Get all context IDs, sort alphabetically
   - Return error with `available_ids` list
7. If found but context_id matches current active account → return already-active message
8. Switch via `azure_cli.set_account(context.subscription_id)`
9. Verify switch with `azure_cli.get_current_account()`
10. Return success with context details

**Error Types**:

- `empty_list`: No contexts saved
- `not_found`: Provided ID doesn't exist (includes `available_ids`)
- `already_active`: Target context is already active
- `verification_failed`: Azure CLI switch failed
- `no_session`: Azure CLI not logged in
- `unknown`: Unexpected error

#### Modified Function: `switch()` command

**Location**: `src/cli.py`

**Modified Signature**:

```python
@app.command()
def switch(
    id: str | None = typer.Option(
        None,
        "--id",
        "-i",
        help="Context ID to switch to directly (case-sensitive). Omit for interactive selection."
    )
) -> None:
    """Switch to a different saved context.

    Displays an interactive list of saved contexts (default) or switches
    directly when --id is provided.

    Examples:
        azctx switch              # Interactive selection
        azctx switch --id DEV     # Direct switch to "DEV" context
        azctx switch -i PROD      # Direct switch using short flag
    """
```

**Implementation Logic**:

```python
if id:
    # Direct switch by ID
    result = context_manager.switch_context_by_id(id)
    # Handle result with appropriate Rich Panel output
    # Include available_ids in error message if present
else:
    # Interactive switch (existing code unchanged)
    result = context_manager.switch_context_interactive()
    # Existing display logic
```

### CLI Contracts

#### Command Usage

**Interactive Mode** (existing, unchanged):

```bash
azctx switch
# Displays interactive list
```

**Direct Mode** (new):

```bash
azctx switch --id DEV
azctx switch -i DEV
```

#### Success Output

```plaintext
╭─ ✓ Successfully Switched to Development ─────────────────╮
│ Name: Development                                         │
│ ID: DEV                                                   │
│ Subscription: My Azure Subscription                       │
│ Tenant: tenant-guid                                       │
│ Account: user@example.com                                 │
╰───────────────────────────────────────────────────────────╯
```

#### Error Output - Context Not Found

```plaintext
╭─ ✗ Context Not Found ─────────────────────────────────────╮
│ Context 'STAGING' not found.                              │
│                                                            │
│ Available contexts: DEV, PROD, TEST                        │
╰───────────────────────────────────────────────────────────╯
```

#### Error Output - No Contexts Saved

```plaintext
╭─ ⚠ No Contexts Available ────────────────────────────────╮
│ No saved contexts found.                                  │
│ Use 'azctx add' to save your current context first.      │
╰───────────────────────────────────────────────────────────╯
```

#### Error Output - Already Active

```plaintext
╭─ ⚠ Already Active ───────────────────────────────────────╮
│ Context 'DEV' is already active.                          │
╰───────────────────────────────────────────────────────────╯
```

### Testing Scenarios (Quickstart)

See `quickstart.md` for detailed manual test scenarios.

---

## Phase 2: Task Breakdown

Phase 2 (task creation) will be handled by the `/speckit.tasks` command, not this planning phase.

Expected task categories:

1. **Code Changes**: Modify `cli.py` and add function to `context_manager.py`
2. **Documentation**: Update README with examples
3. **Manual Testing**: Execute test scenarios from quickstart.md

---

## Implementation Notes

### Backward Compatibility

- Existing interactive mode remains default behavior
- No breaking changes to current users
- Help text updated to reflect new option

### Performance Considerations

- Direct switch bypasses interactive prompt → ~3 second time savings
- Alphabetical sorting of IDs has negligible performance impact for typical usage (< 50 contexts)

### Future Enhancements (Out of Scope)

- Case-insensitive matching (intentionally not included per spec)
- Partial ID matching / fuzzy search
- Auto-completion for context IDs

### Code Review Checklist

- [ ] Type hints added to new function
- [ ] Docstrings follow existing format
- [ ] Error messages use Rich Panel formatting consistently
- [ ] Exit codes match existing patterns (0=success, 1=error, 130=cancelled)
- [ ] Whitespace trimming applied to `context_id` parameter
- [ ] Alphabetical sorting applied to available IDs in error messages
- [ ] Manual testing scenarios from quickstart.md executed successfully

