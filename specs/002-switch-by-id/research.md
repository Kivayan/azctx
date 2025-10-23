# Research: Direct Context Switching by ID

**Feature**: 002-switch-by-id
**Date**: October 23, 2025
**Status**: Complete

## Overview

This research document captures the technical decisions for adding ID-based context switching to the `azctx switch` command.

## Research Questions

### Q1: How to add optional parameters to existing Typer commands?

**Context**: The `azctx switch` command currently has no parameters. We need to add an optional `--id` parameter while preserving backward compatibility.

**Research Findings**:

Typer supports optional parameters via `typer.Option()` with a default value of `None`:

```python
@app.command()
def switch(
    id: str | None = typer.Option(
        None,
        "--id",
        "-i",
        help="Context ID to switch to directly"
    )
) -> None:
    if id:
        # Direct switch logic
    else:
        # Interactive logic (existing)
```

**Decision**: Use `typer.Option(None, "--id", "-i", help="...")` pattern

**Rationale**:

- `None` default preserves backward compatibility - existing users see no changes
- Multiple names (`"--id"`, `"-i"`) provide both verbose and short-form options
- Conditional logic (`if id:`) cleanly separates direct vs interactive modes
- Typer automatically includes parameter in `--help` output

**Alternatives Considered**:

- Separate command (e.g., `azctx switch-by-id`) - Rejected: Creates command proliferation
- Required parameter - Rejected: Breaks backward compatibility

**References**: Typer documentation on optional parameters

---

### Q2: Is Python string comparison case-sensitive by default?

**Context**: The spec requires case-sensitive ID matching (e.g., "DEV" ≠ "dev").

**Research Findings**:

Python's `==` operator for strings is case-sensitive by default:

```python
>>> "DEV" == "dev"
False
>>> "DEV" == "DEV"
True
```

The existing `storage.get_context_by_id()` function uses standard equality:

```python
def get_context_by_id(context_id: str) -> Context | None:
    contexts = load_contexts()
    for ctx in contexts:
        if ctx.context_id == context_id:  # Case-sensitive by default
            return ctx
    return None
```

**Decision**: Use existing `get_context_by_id()` function without modification

**Rationale**:

- No additional code needed - Python's default behavior matches requirement
- Existing storage function already implements correct behavior
- Simple and straightforward

**Alternatives Considered**:

- Case-insensitive matching with `.lower()` - Rejected: Spec explicitly requires case-sensitivity
- Regex matching - Rejected: Over-engineering for exact string match

---

### Q3: How to format error messages with available context IDs?

**Context**: When user provides invalid ID, error message should list all available IDs in alphabetical order for easy scanning.

**Research Findings**:

Python's built-in `sorted()` function provides alphabetical sorting:

```python
contexts = storage.load_contexts()
available_ids = sorted([ctx.context_id for ctx in contexts])
error_msg = f"Available contexts: {', '.join(available_ids)}"
```

For Rich Panel formatting:

```python
console.print(
    Panel(
        f"Context '{id}' not found.\n\n"
        f"Available contexts: {', '.join(available_ids)}",
        title="✗ Context Not Found",
        border_style="red",
    )
)
```

**Decision**: Use `sorted()` + `', '.join()` pattern with Rich Panel

**Rationale**:

- `sorted()` is built-in, no dependencies
- Comma-separated format is readable and conventional
- Rich Panel maintains visual consistency with other error messages
- Two-line format (error + available IDs) is scannable

**Alternatives Considered**:

- Bulleted list - Rejected: Takes more vertical space, harder to scan for ~5 IDs
- Table format - Rejected: Over-engineered for simple ID list
- Unsorted list - Rejected: Harder to find specific ID in longer lists

---

## Best Practices Applied

### Error Handling Pattern

Reuse existing pattern from `context_manager.py`:

```python
return {
    "success": bool,
    "message": str,
    "context": Context | None,
    "error": str | None,
    "available_ids": list[str] | None  # New field for not_found errors
}
```

**Rationale**: Consistent return structure across all context operations

### Input Sanitization

Trim whitespace from user input:

```python
context_id = id.strip()
```

**Rationale**: Prevents accidental failures from copy-paste with trailing spaces

### Performance Considerations

For typical usage (< 50 contexts), `sorted()` has negligible performance impact:

- Time complexity: O(n log n)
- For 50 items: < 1ms on modern hardware
- Acceptable within 2-second performance goal

---

## Summary

All research questions resolved with simple, standard Python patterns. No new dependencies required. Implementation can proceed using:

1. Typer's `Option(None, "--id", "-i")` for optional parameter
2. Existing `storage.get_context_by_id()` for case-sensitive lookup
3. Built-in `sorted()` for alphabetical ID listing

**Next Steps**: Proceed to Phase 1 (Design & Contracts) - documented in `plan.md`
