# CLI Contract: Azure CLI Account Context Switcher

**Feature**: 001-azure-context-switcher
**Date**: October 17, 2025
**Tool**: azctx

## Overview

This document defines the command-line interface contract for the `azctx` tool. All commands follow Typer conventions with auto-generated help text.

## Global Behavior

**Command Name**: `azctx`
**Entry Point**: `python -m src.cli` or `azctx` (if installed)
**Exit Codes**:

- `0`: Success
- `1`: General error (Azure CLI not found, file I/O error, etc.)
- `2`: Validation error (duplicate ID, context not found, etc.)
- `130`: User cancelled (Ctrl+C)

**Global Options**: None (keep simple per constitution)

## Commands

### 1. azctx add

**Purpose**: Add the current Azure CLI context with user-assigned identifiers

**Signature**:

```bash
azctx add
```

**Behavior**:

1. Execute `az account show --output json` to get current Azure session
2. If no active session, display error and exit with code 1
3. Prompt interactively for `context_name` (text input)
4. Prompt interactively for `context_id` (text input)
5. Validate `context_id` uniqueness and format
6. If validation fails, display error and re-prompt
7. Save context to YAML file
8. Display success message with context details

**Interactive Prompts**:

```
Context Name: [user input - friendly name]
Context ID: [user input - short identifier]
```

**Success Output**:

```
✓ Context added successfully!

  ID: prod
  Name: Production Environment
  Subscription: MyCompany Production
  Username: user@mycompany.com
```

**Error Scenarios**:

- Azure CLI not installed → "Error: Azure CLI not found. Please install Azure CLI."
- No active session → "Error: No active Azure session. Run 'az login' first."
- Duplicate context_id → "Error: Context ID 'prod' already exists. Please choose a different ID."
- Invalid context_id format → "Error: Context ID must be 1-20 characters (alphanumeric, hyphens, underscores)."

**Arguments**: None
**Options**: None
**Exit Codes**: 0 (success), 1 (Azure CLI error), 2 (validation error), 130 (cancelled)

---

### 2. azctx switch

**Purpose**: Interactively switch to a saved context

**Signature**:

```bash
azctx switch
```

**Behavior**:

1. Load saved contexts from YAML file
2. If no contexts exist, display error and exit with code 1
3. If only one context exists, inform user and exit
4. Display interactive list with arrow key navigation
5. User selects context (Enter) or cancels (Esc/Ctrl+C)
6. If cancelled, exit with code 130
7. Execute `az account set --subscription <subscription_id>`
8. Verify switch with `az account show`
9. Display success confirmation

**Interactive Selection**:

```
? Select a context:
❯ prod - Production Environment
  dev - Development Environment
  staging - Staging Environment
```

**Success Output**:

```
✓ Switched to context: prod

  Name: Production Environment
  Subscription: MyCompany Production
  Username: user@mycompany.com
```

**Error Scenarios**:

- No saved contexts → "Error: No saved contexts. Use 'azctx add' to create your first context."
- Only one context → "Info: Only one context available ('prod'). No need to switch."
- Azure CLI error during switch → "Error: Failed to switch context. Azure CLI returned an error."

**Arguments**: None
**Options**: None
**Exit Codes**: 0 (success), 1 (error), 130 (cancelled)

---

### 3. azctx list

**Purpose**: List all saved contexts

**Signature**:

```bash
azctx list [--verbose]
```

**Behavior**:

1. Load saved contexts from YAML file
2. If no contexts exist, display message and exit with code 0
3. If `--verbose` flag is NOT set, display simple table (ID and Name only)
4. If `--verbose` flag IS set, display detailed table with all fields
5. Exit with code 0

**Default Output** (without --verbose):

```
Saved Contexts:

ID       Name
──────   ───────────────────────
prod     Production Environment
dev      Development Environment
staging  Staging Environment

Use 'azctx list --verbose' for detailed information.
```

**Verbose Output** (with --verbose):

```
Saved Contexts (Detailed):

Context: prod
  Name: Production Environment
  Subscription: MyCompany Production (12345678-1234-1234-1234-123456789abc)
  Tenant: 87654321-4321-4321-4321-cba987654321
  Username: user@mycompany.com
  Created: 2025-10-17 14:23:45

Context: dev
  Name: Development Environment
  Subscription: MyCompany Development (abcdef12-3456-7890-abcd-ef1234567890)
  Tenant: 87654321-4321-4321-4321-cba987654321
  Username: user@mycompany.com
  Created: 2025-10-17 15:10:22
```

**Empty State**:

```
No saved contexts found.

Use 'azctx add' to create your first context.
```

**Arguments**: None
**Options**:

- `--verbose` / `-v`: Show detailed information including subscription IDs, tenant IDs, usernames, and creation timestamps

**Exit Codes**: 0 (always succeeds)

---

### 4. azctx status

**Purpose**: Display the currently active Azure context

**Signature**:

```bash
azctx status [--verbose]
```

**Behavior**:

1. Execute `az account show --output json` to get current Azure session
2. If no active session, display error and exit with code 1
3. Load saved contexts from YAML file
4. Attempt to match current subscription_id against saved contexts
5. If match found, display as "managed context" with friendly names
6. If no match, display as "unmanaged context" with raw Azure data
7. If `--verbose` flag is set, include additional details

**Default Output** (managed context):

```
✓ Active Context: prod

  Name: Production Environment
  Subscription: MyCompany Production
  Username: user@mycompany.com
```

**Verbose Output** (managed context):

```
✓ Active Context: prod

  Name: Production Environment
  Subscription: MyCompany Production (12345678-1234-1234-1234-123456789abc)
  Tenant: 87654321-4321-4321-4321-cba987654321
  Username: user@mycompany.com
  Created: 2025-10-17 14:23:45
```

**Unmanaged Context Output**:

```
⚠ Active Azure session (not managed by azctx)

  Subscription: MyCompany Production (12345678-1234-1234-1234-123456789abc)
  Username: user@mycompany.com

Use 'azctx add' to save this context.
```

**No Session Output**:

```
✗ No active Azure session

Run 'az login' to authenticate.
```

**Arguments**: None
**Options**:

- `--verbose` / `-v`: Show detailed information including full IDs and timestamps

**Exit Codes**: 0 (active session found), 1 (no active session)

---

### 5. azctx delete

**Purpose**: Interactively delete a saved context

**Signature**:

```bash
azctx delete
```

**Behavior**:

1. Load saved contexts from YAML file
2. If no contexts exist, display error and exit with code 1
3. Display interactive list with arrow key navigation
4. User selects context to delete (Enter) or cancels (Esc/Ctrl+C)
5. If cancelled, exit with code 130
6. Prompt for confirmation (yes/no)
7. If confirmed, remove context from YAML file
8. Display success message
9. If context was currently active in Azure CLI, note that Azure session remains active

**Interactive Selection**:

```
? Select a context to delete:
❯ prod - Production Environment
  dev - Development Environment
  staging - Staging Environment
```

**Confirmation Prompt**:

```
? Are you sure you want to delete context 'prod'? (y/N)
```

**Success Output**:

```
✓ Context 'prod' deleted successfully.

Note: Azure CLI session remains active. Run 'azctx switch' to change contexts.
```

**Cancelled Output**:

```
Deletion cancelled.
```

**Error Scenarios**:

- No saved contexts → "Error: No saved contexts to delete."

**Arguments**: None
**Options**: None
**Exit Codes**: 0 (success), 1 (no contexts), 130 (cancelled)

---

### 6. azctx --help

**Purpose**: Display help information (auto-generated by Typer)

**Signature**:

```bash
azctx --help
azctx [command] --help
```

**Behavior**: Typer auto-generates help text from command docstrings and type hints

**Example Output**:

```
Usage: azctx [OPTIONS] COMMAND [ARGS]...

  Azure CLI context switcher - manage and switch between Azure accounts quickly.

Options:
  --help  Show this message and exit.

Commands:
  add     Add current Azure CLI context with friendly identifiers
  delete  Delete a saved context
  list    List all saved contexts
  status  Show currently active Azure context
  switch  Switch to a different saved context
```

---

## Input Validation

### context_id Validation

**Rules**:

- Must be 1-20 characters
- Allowed characters: `a-z`, `A-Z`, `0-9`, `-`, `_`
- Must be unique across all saved contexts
- Case-sensitive

**Regex**: `^[a-zA-Z0-9_-]{1,20}$`

**Examples**:

- Valid: `prod`, `dev-west`, `staging_2`, `azure-prod-01`
- Invalid: ``, `prod env` (space), `prod!` (special char), `very-long-identifier-name` (> 20 chars)

### context_name Validation

**Rules**:

- Must be 1-100 characters
- Any printable characters allowed (including spaces and special chars)
- No uniqueness requirement (IDs are unique, names can duplicate)

**Examples**:

- Valid: `Production Environment`, `Dev (West US)`, `Staging @ Azure`, `🚀 Production`
- Invalid: `` (empty string)

## Color Scheme

Using Rich for colored output:

- **Success** (green): `✓` symbol, confirmation messages
- **Error** (red): `✗` symbol, error messages
- **Warning** (yellow): `⚠` symbol, warnings (unmanaged context)
- **Info** (cyan): General information, table headers
- **Highlight** (bold): Context names, IDs

## Performance Targets

All commands must complete within:

- `add`: < 2 seconds (excluding user input time)
- `switch`: < 2 seconds (excluding user selection time)
- `list`: < 1 second
- `status`: < 1 second
- `delete`: < 1 second (excluding user input and confirmation time)

## Error Message Format

All error messages use Rich Panel with red border:

```
╭─ Error ──────────────────────────────────────╮
│ Azure CLI not found. Please install Azure    │
│ CLI from https://aka.ms/azure-cli             │
╰───────────────────────────────────────────────╯
```

## Non-Interactive Mode Support

**Current version**: All commands are interactive (Questionary prompts)

**Future consideration**: Could add non-interactive flags for scripting:

```bash
azctx switch --id prod  # Non-interactive switch
azctx add --id dev --name "Development"  # Non-interactive add
```

**Not implemented in v1.0** per constitution (simplicity first, add if needed).
