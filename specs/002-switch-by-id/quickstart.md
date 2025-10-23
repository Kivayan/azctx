# Quickstart: Direct Context Switching by ID

**Feature**: 002-switch-by-id
**Date**: October 23, 2025
**Purpose**: Manual testing scenarios for ID-based context switching

## Prerequisites

Before testing this feature, ensure:

1. Azure CLI is installed (`az --version`)
2. You are logged into Azure (`az login`)
3. You have `azctx` built and available
4. You have at least 2-3 contexts saved via `azctx add`

## Setup Test Data

If you don't have saved contexts, create test contexts:

```bash
# Switch to first Azure account
az account set --subscription <subscription-1>

# Add first context
azctx add
# Enter name: "Development Environment"
# Enter ID: "DEV"

# Switch to second Azure account
az account set --subscription <subscription-2>

# Add second context
azctx add
# Enter name: "Production Environment"
# Enter ID: "PROD"

# (Optional) Add third context
az account set --subscription <subscription-3>
azctx add
# Enter name: "Test Environment"
# Enter ID: "TEST"
```

## Test Scenarios

### Scenario 1: Direct Switch with Valid ID (Happy Path)

**Test**: Switch to a known context using long-form flag

```bash
azctx switch --id DEV
```

**Expected Output**:

```plaintext
╭─ ✓ Successfully Switched to Development Environment ──────╮
│ Name: Development Environment                             │
│ ID: DEV                                                    │
│ Subscription: <subscription-name>                          │
│ Tenant: <tenant-id>                                        │
│ Account: user@example.com                                  │
╰────────────────────────────────────────────────────────────╯
```

**Verification**:

```bash
az account show --query "[name, id]" -o tsv
# Should show the DEV subscription details
```

**Exit Code**: `0` (success)

---

### Scenario 2: Direct Switch with Short Flag

**Test**: Switch using `-i` shorthand

```bash
azctx switch -i PROD
```

**Expected Output**: Same format as Scenario 1, but for PROD context

**Exit Code**: `0` (success)

---

### Scenario 3: Case-Sensitive Matching

**Test**: Attempt to switch using wrong case

```bash
azctx switch --id dev
```

**Expected Output**:

```plaintext
╭─ ✗ Context Not Found ──────────────────────────────────────╮
│ Context 'dev' not found.                                   │
│                                                             │
│ Available contexts: DEV, PROD, TEST                         │
╰─────────────────────────────────────────────────────────────╯
```

**Exit Code**: `1` (error)

**Verification**: Current Azure account should be unchanged

---

### Scenario 4: Invalid Context ID

**Test**: Provide ID that doesn't exist

```bash
azctx switch --id STAGING
```

**Expected Output**:

```plaintext
╭─ ✗ Context Not Found ──────────────────────────────────────╮
│ Context 'STAGING' not found.                               │
│                                                             │
│ Available contexts: DEV, PROD, TEST                         │
╰─────────────────────────────────────────────────────────────╯
```

**Verification**:

- Error message lists all available IDs
- IDs are sorted alphabetically
- Current Azure account is unchanged

**Exit Code**: `1` (error)

---

### Scenario 5: Already Active Context

**Test**: Switch to currently active context

```bash
# First, switch to DEV
azctx switch --id DEV

# Then try to switch to DEV again
azctx switch --id DEV
```

**Expected Output**:

```plaintext
╭─ ⚠ Already Active ─────────────────────────────────────────╮
│ Context 'DEV' is already active.                           │
╰─────────────────────────────────────────────────────────────╯
```

**Exit Code**: `0` or `1` (check implementation - informational message)

---

### Scenario 6: No Contexts Saved

**Test**: Attempt switch when no contexts exist

**Setup**:

```bash
# Backup existing contexts
mv ~/.azctx/contexts.yaml ~/.azctx/contexts.yaml.backup

# Try to switch
azctx switch --id DEV
```

**Expected Output**:

```plaintext
╭─ ⚠ No Contexts Available ──────────────────────────────────╮
│ No saved contexts found.                                   │
│ Use 'azctx add' to save your current context first.       │
╰─────────────────────────────────────────────────────────────╯
```

**Cleanup**:

```bash
# Restore backup
mv ~/.azctx/contexts.yaml.backup ~/.azctx/contexts.yaml
```

**Exit Code**: `1` (error)

---

### Scenario 7: Whitespace Handling

**Test**: Provide ID with leading/trailing spaces

```bash
azctx switch --id " DEV "
```

**Expected Behavior**: Should trim whitespace and match "DEV" context

**Expected Output**: Success message (same as Scenario 1)

**Exit Code**: `0` (success)

---

### Scenario 8: Interactive Mode Still Works (Backward Compatibility)

**Test**: Ensure existing interactive mode is unchanged

```bash
azctx switch
```

**Expected Behavior**:

- Interactive list appears
- Arrow keys navigate
- Enter selects
- Esc/Ctrl+C cancels

**Expected Output**: Interactive questionary prompt (unchanged from before)

**Exit Code**:

- `0` on successful selection
- `130` on cancellation (Ctrl+C)

---

### Scenario 9: Help Text Includes New Parameter

**Test**: Check command help

```bash
azctx switch --help
```

**Expected Output** (includes):

```plaintext
Usage: azctx switch [OPTIONS]

  Switch to a different saved context.

  Displays an interactive list of saved contexts (default) or switches
  directly when --id is provided.

Options:
  -i, --id TEXT  Context ID to switch to directly (case-sensitive). Omit for
                 interactive selection.
  --help         Show this message and exit.
```

**Verification**: Both `-i` and `--id` are documented

---

### Scenario 10: Performance Check

**Test**: Measure switch time

```bash
time azctx switch --id PROD
```

**Expected**: Total time < 2 seconds (per success criteria SC-001)

**Comparison**: Should be faster than interactive mode (~5 seconds)

---

## Error Handling Tests

### Test: Empty ID Parameter

**Command**:

```bash
azctx switch --id ""
```

**Expected**: Error message about invalid/empty ID

---

### Test: Azure CLI Not Installed

**Setup**: Temporarily rename `az` command or set invalid PATH

**Command**:

```bash
azctx switch --id DEV
```

**Expected**: Clear error message about Azure CLI not found

---

### Test: Azure CLI Not Logged In

**Setup**:

```bash
az logout
```

**Command**:

```bash
azctx switch --id DEV
```

**Expected**: Error message about no active Azure session

**Cleanup**:

```bash
az login
```

---

## Scripting Test

**Test**: Use in a shell script

```bash
#!/bin/bash
# deploy.sh - Switch to PROD and run deployment

azctx switch --id PROD
if [ $? -eq 0 ]; then
    echo "Switched to PROD. Running deployment..."
    # Deployment commands here
else
    echo "Failed to switch to PROD. Aborting."
    exit 1
fi
```

**Expected**: Exit codes enable reliable scripting logic

---

## Regression Tests

Ensure existing functionality still works:

- [ ] `azctx add` - Adding new contexts
- [ ] `azctx list` - Listing contexts (basic and verbose)
- [ ] `azctx status` - Checking current context
- [ ] `azctx delete` - Deleting contexts
- [ ] `azctx switch` (no parameters) - Interactive mode

---

## Acceptance Criteria Verification

Map test scenarios to spec requirements:

| Requirement | Test Scenario | Pass/Fail |
|-------------|---------------|-----------|
| FR-001: Support `--id` flag | Scenario 1 | ⬜ |
| FR-002: Support `-i` shorthand | Scenario 2 | ⬜ |
| FR-003: Case-sensitive matching | Scenario 3 | ⬜ |
| FR-004: Success confirmation | Scenario 1, 2 | ⬜ |
| FR-005: Error for invalid ID | Scenario 4 | ⬜ |
| FR-006: List available IDs in error | Scenario 4 | ⬜ |
| FR-007: Alphabetical sorting of IDs | Scenario 4 | ⬜ |
| FR-009: Complete within 2 seconds | Scenario 10 | ⬜ |
| FR-010: Already-active message | Scenario 5 | ⬜ |
| FR-011: Guide users when no contexts | Scenario 6 | ⬜ |
| FR-014: Trim whitespace | Scenario 7 | ⬜ |
| FR-015: Prioritize `--id` over interactive | Scenario 1-7 | ⬜ |

---

## Success Criteria Verification

| Success Criterion | Test Method | Pass/Fail |
|-------------------|-------------|-----------|
| SC-001: Switch in under 2 seconds | Scenario 10 (timing) | ⬜ |
| SC-002: 100% valid IDs switch successfully | Scenarios 1, 2 | ⬜ |
| SC-003: Error messages include available IDs | Scenario 4 | ⬜ |
| SC-004: Scriptable (exit codes) | Scripting Test | ⬜ |
| SC-005: Case-sensitive matching works | Scenario 3 | ⬜ |
| SC-006: Help text documents flags | Scenario 9 | ⬜ |

---

## Completion Checklist

After running all scenarios:

- [ ] All test scenarios pass
- [ ] All functional requirements verified
- [ ] All success criteria met
- [ ] No regressions in existing commands
- [ ] Help text is accurate
- [ ] Performance goals achieved
- [ ] Error messages are clear and helpful
- [ ] Feature works on target platform (Windows/macOS/Linux)

---

## Known Limitations (Expected Behavior)

- No fuzzy matching - exact ID required
- No auto-completion (shell completion not implemented)
- No case-insensitive mode (intentional per spec)
- IDs with spaces require quotes: `azctx switch --id "MY ID"` (but not recommended per validation rules)
