# Data Model: Direct Context Switching by ID

**Feature**: 002-switch-by-id
**Date**: October 23, 2025

## Overview

This feature does **not introduce new data models**. It extends existing functionality to use the established `Context` entity for ID-based lookup and switching.

## Existing Entity: Context

**Location**: `src/models/context.py`

**Definition**:

```python
@dataclass
class Context:
    """Represents a saved Azure CLI account context."""

    context_id: str          # Unique identifier (1-20 chars, alphanumeric/hyphens/underscores)
    context_name: str        # User-friendly name (1-100 chars)
    subscription_id: str     # Azure subscription GUID
    subscription_name: str   # Azure subscription display name
    tenant_id: str          # Azure tenant GUID
    tenant_name: str        # Azure tenant display name (simplified to tenant_id)
    username: str           # Azure account email/username
    created_at: datetime    # Timestamp when context was added
```

## Key Field for This Feature

### `context_id` (str)

**Purpose**: Unique identifier used for the `--id` parameter matching

**Characteristics**:

- **Case-sensitive**: "DEV" and "dev" are distinct identifiers
- **Unique**: Enforced by `storage.context_id_exists()` check during add operation
- **Format**: 1-20 characters, alphanumeric with hyphens and underscores
- **Validation**: Performed by `Context.validate_context_id()` during context creation

**Usage in This Feature**:

```python
# User command
azctx switch --id DEV

# Lookup (case-sensitive)
context = storage.get_context_by_id("DEV")  # Returns Context object or None

# Comparison
if ctx.context_id == "DEV":  # Exact match required
    # Found
```

## Data Flow

### Direct Switch by ID

```plaintext
User Input (CLI)
    ↓
azctx switch --id DEV
    ↓
cli.switch(id="DEV")
    ↓
context_manager.switch_context_by_id("DEV")
    ↓
storage.get_context_by_id("DEV")  ← Case-sensitive lookup
    ↓
storage.load_contexts()           ← Read YAML file
    ↓
List[Context]                     ← All saved contexts
    ↓
Find ctx where ctx.context_id == "DEV"  ← Linear search with ==
    ↓
Context object or None
    ↓
azure_cli.set_account(context.subscription_id)
    ↓
Success/Error result dict
```

## Storage Format

**File**: `~/.azctx/contexts.yaml`

**Format** (unchanged):

```yaml
contexts:
  - context_id: "DEV"
    context_name: "Development Environment"
    subscription_id: "sub-guid-1"
    subscription_name: "My Dev Subscription"
    tenant_id: "tenant-guid"
    tenant_name: "tenant-guid"
    username: "user@example.com"
    created_at: "2025-10-20T10:30:00"

  - context_id: "PROD"
    context_name: "Production Environment"
    subscription_id: "sub-guid-2"
    subscription_name: "My Prod Subscription"
    tenant_id: "tenant-guid"
    tenant_name: "tenant-guid"
    username: "user@example.com"
    created_at: "2025-10-21T14:15:00"
```

**Key Point**: The `context_id` field already exists in storage - no schema changes required.

## Validation Rules

### Input Validation

**Parameter**: `--id` / `-i` value

**Rules** (applied by new function):

1. **Whitespace trimming**: Leading/trailing spaces removed via `.strip()`
2. **Empty check**: After trimming, must not be empty string
3. **Case preservation**: Input case is preserved for lookup (no `.lower()` or `.upper()`)

**Example**:

```python
# User input with accidental space
azctx switch --id " DEV "

# After .strip()
context_id = "DEV"  # Correctly matches stored ID
```

### Lookup Logic

**Function**: `storage.get_context_by_id(context_id: str)`

**Algorithm** (existing implementation):

```python
def get_context_by_id(context_id: str) -> Context | None:
    contexts = load_contexts()
    for ctx in contexts:
        if ctx.context_id == context_id:  # Exact case-sensitive match
            return ctx
    return None  # Not found
```

**Time Complexity**: O(n) where n = number of saved contexts

**Performance**: Acceptable for typical usage (< 50 contexts, < 1ms lookup time)

## Error Cases

### Not Found

**Condition**: No context with matching `context_id` exists

**Response**:

```python
{
    "success": False,
    "error": "not_found",
    "message": "Context 'STAGING' not found.",
    "available_ids": ["DEV", "PROD", "TEST"],  # Sorted alphabetically
    "context": None
}
```

### Empty List

**Condition**: No contexts saved in storage

**Response**:

```python
{
    "success": False,
    "error": "empty_list",
    "message": "No saved contexts found. Use 'azctx add' to save your current context first.",
    "available_ids": None,
    "context": None
}
```

### Already Active

**Condition**: Target context is already the active Azure CLI account

**Response**:

```python
{
    "success": True,  # Not an error, but informational
    "error": "already_active",
    "message": "Context 'DEV' is already active.",
    "available_ids": None,
    "context": <Context object>
}
```

## No Schema Changes

**Confirmation**:

- ✅ No new fields added to `Context` dataclass
- ✅ No changes to YAML storage format
- ✅ No database migrations needed (file-based storage)
- ✅ No changes to `Context.to_dict()` or `Context.from_dict()`

**Rationale**: Feature leverages existing `context_id` field, which was designed for this exact purpose.

## Summary

This feature is a **pure behavior extension** with zero data model changes. The existing `Context.context_id` field provides all necessary functionality for case-sensitive ID-based lookup and switching.
