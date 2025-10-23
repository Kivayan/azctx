# Function Contracts: Direct Context Switching by ID

**Feature**: 002-switch-by-id
**Date**: October 23, 2025

## Overview

This document defines the function signatures and contracts for the ID-based context switching feature.

---

## New Function: `switch_context_by_id()`

**Module**: `src/services/context_manager.py`

**Signature**:

```python
def switch_context_by_id(context_id: str) -> dict[str, Any]:
    """Switch to a saved Azure CLI context by its ID (case-sensitive).

    This function provides direct, non-interactive context switching by ID.
    It performs case-sensitive matching against saved context IDs and
    returns detailed result information for CLI display.

    Args:
        context_id: The case-sensitive ID of the context to switch to.
                    Leading and trailing whitespace will be trimmed before matching.
                    Empty strings (after trimming) will result in an error.

    Returns:
        Dictionary with operation result containing these keys:

        - success (bool): Whether the operation succeeded
        - message (str): Human-readable status message for display
        - context (Context | None): Selected Context object if successful, None on error
        - error (str | None): Error type identifier (see Error Types below)
        - available_ids (list[str] | None): Alphabetically sorted list of all context IDs
                                           (only included when error == "not_found")

    Raises:
        AzureCliNotFoundError: If Azure CLI is not installed or not accessible in PATH.
                              This is a fatal error that should terminate the program.

    Error Types (value of 'error' field):
        - None: Success (success == True)
        - "empty_list": No contexts are saved in storage
        - "not_found": Provided context_id doesn't match any saved context
        - "already_active": Target context is already the active Azure account
        - "verification_failed": Azure CLI switch command executed but verification failed
        - "no_session": No active Azure CLI session (user needs to run 'az login')
        - "unknown": Unexpected error occurred

    Examples:
        >>> # Successful switch
        >>> result = switch_context_by_id("DEV")
        >>> result
        {
            "success": True,
            "message": "Successfully switched to context: DEV",
            "context": <Context object>,
            "error": None,
            "available_ids": None
        }

        >>> # Context not found
        >>> result = switch_context_by_id("STAGING")
        >>> result
        {
            "success": False,
            "message": "Context 'STAGING' not found.",
            "context": None,
            "error": "not_found",
            "available_ids": ["DEV", "PROD", "TEST"]  # Sorted alphabetically
        }

        >>> # Already active
        >>> result = switch_context_by_id("DEV")  # DEV is current context
        >>> result
        {
            "success": False,
            "message": "Context 'DEV' is already active.",
            "context": <Context object>,
            "error": "already_active",
            "available_ids": None
        }

        >>> # No contexts saved
        >>> result = switch_context_by_id("DEV")
        >>> result
        {
            "success": False,
            "message": "No saved contexts found. Use 'azctx add' to save your current context first.",
            "context": None,
            "error": "empty_list",
            "available_ids": None
        }

    Implementation Notes:
        1. Trim whitespace from context_id parameter using .strip()
        2. Verify Azure CLI is available (raise AzureCliNotFoundError if not)
        3. Load contexts via storage.load_contexts()
        4. Check for empty list (error: "empty_list")
        5. Use storage.get_context_by_id() for case-sensitive lookup
        6. If not found, generate sorted available_ids list (error: "not_found")
        7. Check if target context is already active (error: "already_active")
        8. Switch via azure_cli.set_account(context.subscription_id)
        9. Verify switch with azure_cli.get_current_account()
        10. Return success or appropriate error result

    Performance:
        - Should complete within 2 seconds on standard hardware
        - Context lookup is O(n) where n = number of saved contexts
        - Typical usage (< 50 contexts) has negligible performance impact
    """
```

---

## Modified Function: `switch()` CLI Command

**Module**: `src/cli.py`

**Signature**:

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

    Displays an interactive list of saved contexts (default behavior when no
    options are provided) or switches directly to a specific context when the
    --id option is used.

    This command supports two modes:

    1. Interactive Mode (default): Presents a navigable list of saved contexts
       using keyboard navigation (arrow keys + Enter). This is the original
       behavior and remains unchanged.

    2. Direct Mode (--id flag): Switches to a specific context by its ID without
       any interactive prompts. Useful for scripting and automation.

    Args:
        id: Optional context ID for direct switching. If None, interactive mode
            is used. If provided, the command attempts to switch directly to the
            specified context (case-sensitive match).

    Returns:
        None. Exits with appropriate exit codes:
        - 0: Success
        - 1: Error (context not found, Azure CLI issue, etc.)
        - 130: User cancellation (interactive mode only)

    Examples:
        # Interactive mode (default)
        azctx switch

        # Direct mode with long flag
        azctx switch --id DEV

        # Direct mode with short flag
        azctx switch -i PROD

    Display Format:
        Success and error messages are displayed using Rich Panel formatting
        for consistent visual presentation. See CLI contract documentation
        for specific message formats.

    Implementation Notes:
        - If id parameter is provided, call context_manager.switch_context_by_id(id)
        - If id parameter is None, call context_manager.switch_context_interactive()
        - Handle result dictionary and display using Rich Panel
        - For "not_found" errors, include available_ids in error message
        - Set appropriate exit codes based on result
    """
```

---

## Existing Functions (Used but Not Modified)

### `storage.get_context_by_id()`

**Module**: `src/services/storage.py`

**Signature**:

```python
def get_context_by_id(context_id: str) -> Context | None:
    """Find and return context by ID (case-sensitive).

    Args:
        context_id: The ID of the context to find.

    Returns:
        Context object if found, None otherwise.

    Note:
        This function performs exact case-sensitive string matching.
        Example: "DEV" will NOT match "dev".
    """
```

**Usage**: Called by `switch_context_by_id()` to perform case-sensitive context lookup.

---

### `storage.load_contexts()`

**Module**: `src/services/storage.py`

**Signature**:

```python
def load_contexts() -> list[Context]:
    """Read and parse YAML file, return list of Context objects.

    Returns:
        List of Context objects. Empty list if file doesn't exist or is empty.

    Raises:
        StorageError: If the YAML file is corrupted or cannot be read.
    """
```

**Usage**: Called by `switch_context_by_id()` to load all saved contexts.

---

### `azure_cli.set_account()`

**Module**: `src/services/azure_cli.py`

**Signature**:

```python
def set_account(subscription_id: str) -> None:
    """Switch to specified Azure subscription.

    Args:
        subscription_id: The Azure subscription GUID to switch to.

    Raises:
        NoActiveSessionError: If not logged into Azure CLI.
        Various subprocess errors: If az command fails.
    """
```

**Usage**: Called by `switch_context_by_id()` to perform the actual Azure CLI switch.

---

### `azure_cli.get_current_account()`

**Module**: `src/services/azure_cli.py`

**Signature**:

```python
def get_current_account() -> dict | None:
    """Get currently active Azure account information.

    Returns:
        Dictionary with account details including 'id', 'name', 'tenantId', etc.
        None if no active session.

    Raises:
        NoActiveSessionError: If not logged into Azure CLI.
    """
```

**Usage**: Called by `switch_context_by_id()` to verify the switch was successful.

---

### `azure_cli.check_azure_cli_installed()`

**Module**: `src/services/azure_cli.py`

**Signature**:

```python
def check_azure_cli_installed() -> bool:
    """Check if Azure CLI is installed and accessible.

    Returns:
        True if Azure CLI is available, False otherwise.
    """
```

**Usage**: Called by `switch_context_by_id()` to verify Azure CLI availability before attempting operations.

---

## Type Definitions

### Context (Dataclass)

**Module**: `src/models/context.py`

```python
@dataclass
class Context:
    context_id: str
    context_name: str
    subscription_id: str
    subscription_name: str
    tenant_id: str
    tenant_name: str
    username: str
    created_at: datetime
```

### Result Dictionary Structure

**Used by**: `switch_context_by_id()` and `switch_context_interactive()`

```python
{
    "success": bool,            # Required: operation outcome
    "message": str,             # Required: human-readable message
    "context": Context | None,  # Required: context object or None
    "error": str | None,        # Required: error type or None
    "available_ids": list[str] | None  # Optional: only for "not_found" error
}
```

---

## Integration Points

### CLI → Service Layer

```python
# In src/cli.py
from src.services import context_manager

# Direct switch
result = context_manager.switch_context_by_id(id)
```

### Service → Storage Layer

```python
# In src/services/context_manager.py
from src.services import storage

# Load contexts
contexts = storage.load_contexts()

# Lookup specific context
context = storage.get_context_by_id(context_id)
```

### Service → Azure CLI Layer

```python
# In src/services/context_manager.py
from src.services import azure_cli

# Check CLI availability
if not azure_cli.check_azure_cli_installed():
    raise AzureCliNotFoundError(...)

# Perform switch
azure_cli.set_account(context.subscription_id)

# Verify switch
current = azure_cli.get_current_account()
```

---

## Testing Contract

### Unit Test Considerations (If Implemented)

While formal tests are not required per project constitution, the function signature supports testing:

```python
# Mock storage.get_context_by_id() to return specific contexts
# Mock azure_cli functions to simulate different scenarios
# Verify result dictionary structure and content

# Test cases:
# - Valid ID returns success
# - Invalid ID returns not_found with available_ids
# - Empty storage returns empty_list error
# - Already active context returns already_active
# - Azure CLI errors are handled gracefully
```

### Manual Test Coverage

See `quickstart.md` for comprehensive manual testing scenarios covering:

- Happy path (valid ID)
- Error cases (not found, empty list, already active)
- Edge cases (whitespace, case sensitivity)
- Performance verification (< 2 seconds)

---

## Backward Compatibility

### Existing Function Behavior

**`switch_context_interactive()`**: Unchanged

- Still used when `azctx switch` is called without `--id` parameter
- Return structure remains the same
- No breaking changes

### New Function Integration

**`switch_context_by_id()`**: New function with compatible return structure

- Uses same result dictionary pattern as `switch_context_interactive()`
- Enables consistent error handling in CLI layer
- Adds optional `available_ids` field without breaking existing code

---

## Summary

This feature adds one new function (`switch_context_by_id()`) and modifies one existing function signature (`switch()` command parameter). All other functions are used as-is with no modifications required.

**Key Design Decisions**:

1. **Consistent return structure**: Both switch functions return dictionaries with same core keys
2. **Optional `available_ids` field**: Only included when relevant (not_found error)
3. **Case-sensitive by design**: Uses existing storage function without modification
4. **Whitespace trimming**: Applied at service layer, not storage layer
5. **Type hints**: Full type annotations for IDE support and clarity
