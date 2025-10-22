# Data Model: Azure CLI Account Context Switcher

**Feature**: 001-azure-context-switcher
**Date**: October 17, 2025
**Status**: Complete

## Overview

This document defines the data structures for the Azure CLI context switcher. The tool uses a simple, file-based storage model with a single entity type.

## Core Entity: Context

Represents a saved Azure CLI account configuration with user-assigned identifiers and Azure metadata.

### Fields

| Field Name | Type | Required | Description | Source | Constraints |
|------------|------|----------|-------------|--------|-------------|
| `context_id` | string | Yes | Short, user-assigned identifier for quick switching | User input | Must be unique across all contexts; 1-20 chars; alphanumeric, hyphens, underscores only |
| `context_name` | string | Yes | Friendly, descriptive name for the context | User input | 1-100 chars; any printable characters |
| `subscription_id` | string | Yes | Azure subscription GUID | `az account show` | UUID format |
| `subscription_name` | string | Yes | Azure subscription display name | `az account show` | As returned by Azure CLI |
| `tenant_id` | string | Yes | Azure tenant GUID | `az account show` | UUID format |
| `tenant_name` | string | Yes | Azure tenant domain name | `az account show` | As returned by Azure CLI |
| `username` | string | Yes | Authenticated user email/principal | `az account show` | As returned by Azure CLI |
| `created_at` | datetime | Yes | When context was added to tool | System timestamp | ISO 8601 format |

### Example

```python
# Python dataclass representation
from dataclasses import dataclass
from datetime import datetime

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

```yaml
# YAML file representation
contexts:
  - context_id: "prod"
    context_name: "Production Environment"
    subscription_id: "12345678-1234-1234-1234-123456789abc"
    subscription_name: "MyCompany Production Subscription"
    tenant_id: "87654321-4321-4321-4321-cba987654321"
    tenant_name: "mycompany.onmicrosoft.com"
    username: "user@mycompany.com"
    created_at: "2025-10-17T14:23:45.123456"
```

### Validation Rules

1. **context_id uniqueness**: When adding a new context, validate that `context_id` does not already exist in the stored contexts
2. **context_id format**: Must match regex `^[a-zA-Z0-9_-]{1,20}$`
3. **context_name length**: Must be 1-100 characters
4. **Azure IDs**: subscription_id and tenant_id should be valid UUIDs (optional validation)
5. **created_at**: Always set to current UTC time when adding context

### State Transitions

Contexts are immutable once created (except for deletion). No state transitions or updates.

**Lifecycle**:
1. **Created**: User runs `azctx add`, context is created with current Azure CLI session data
2. **Exists**: Context persists in YAML file
3. **Deleted**: User runs `azctx delete`, context is removed from YAML file

No status field needed - contexts are either present or absent.

## Data Storage

### Storage Format

**File Type**: YAML
**Location**: `~/.azctx/contexts.yaml` (cross-platform using `Path.home()`)
**Structure**: Single file containing list of context objects

```yaml
# Root structure
contexts:  # List of Context objects
  - context_id: "..."
    context_name: "..."
    # ... other fields
  - context_id: "..."
    context_name: "..."
    # ... other fields
```

### File Operations

**Create/Initialize**:
- Check if `~/.azctx/` directory exists, create if not
- If `contexts.yaml` doesn't exist, create with empty structure:
  ```yaml
  contexts: []
  ```

**Read**:
- Load YAML file using `yaml.safe_load()`
- Parse into list of Context objects
- Return empty list if file is empty or corrupted

**Write** (Add):
- Load existing contexts
- Append new context to list
- Write back using `yaml.safe_dump()`

**Write** (Delete):
- Load existing contexts
- Filter out context with matching `context_id`
- Write back using `yaml.safe_dump()`

**Concurrent Access**:
- No locking mechanism (single-user tool assumption)
- Risk accepted: If two terminal sessions modify simultaneously, last write wins
- Mitigation: Document limitation, unlikely scenario for single developer

### File Size Considerations

**Expected**: 10-20 contexts (< 5KB)
**Maximum**: ~1000 contexts (< 500KB)
**Performance**: YAML parse/write < 50ms for expected size

No pagination or indexing needed given small dataset size.

## Azure CLI Data Mapping

### Source: `az account show --output json`

The Azure CLI returns JSON with this structure:

```json
{
  "id": "12345678-1234-1234-1234-123456789abc",
  "name": "MyCompany Production Subscription",
  "tenantId": "87654321-4321-4321-4321-cba987654321",
  "user": {
    "name": "user@mycompany.com",
    "type": "user"
  },
  "homeTenantId": "87654321-4321-4321-4321-cba987654321",
  "managedByTenants": [],
  ...
}
```

**Mapping to Context fields**:

| Context Field | Azure CLI JSON Path | Notes |
|---------------|---------------------|-------|
| `subscription_id` | `["id"]` | Subscription GUID |
| `subscription_name` | `["name"]` | Subscription display name |
| `tenant_id` | `["tenantId"]` | Tenant GUID |
| `username` | `["user"]["name"]` | User email/principal |

**tenant_name**: Not directly provided by `az account show`.

**Options**:
1. Extract from username domain (e.g., `user@mycompany.com` → `mycompany.com`)
2. Use tenant_id as fallback if domain extraction fails
3. Execute additional `az account tenant show` (slower, adds complexity)

**Decision**: Use tenant_id as tenant_name initially (simpler). Can enhance later if needed.

**Updated mapping**:

| Context Field | Derivation |
|---------------|------------|
| `tenant_name` | Use `tenant_id` value (simple, always available) |

## Error Handling

### Missing/Corrupted File

**Scenario**: `contexts.yaml` is deleted or corrupted
**Behavior**: Treat as empty state, reinitialize with `contexts: []`
**User message**: "No saved contexts found. Use 'azctx add' to create your first context."

### Invalid Context Data

**Scenario**: YAML contains invalid data structure
**Behavior**: Log warning, skip invalid entries, load valid ones
**User message**: "Warning: Some contexts could not be loaded. File may be corrupted."

### Duplicate context_id

**Scenario**: User tries to add context with existing ID
**Behavior**: Reject with error before modifying file
**User message**: "Context ID 'prod' already exists. Please choose a different ID."

## Display Formats

### Basic List (Default)

```
ID       Name
-----    ---------------------
prod     Production Environment
dev      Development Environment
staging  Staging Environment
```

### Verbose List (--verbose flag)

```
Context: prod
  Name: Production Environment
  Subscription: MyCompany Production (12345678-1234-1234-1234-123456789abc)
  Tenant: 87654321-4321-4321-4321-cba987654321
  Username: user@mycompany.com
  Created: 2025-10-17 14:23:45

Context: dev
  Name: Development Environment
  ...
```

Alternative verbose format (table):

| ID | Name | Subscription Name | Subscription ID | Tenant ID | Username | Created |
|----|------|-------------------|-----------------|-----------|----------|---------|
| prod | Production Env | MyCompany Prod | 1234... | 8765... | user@... | 2025-10-17 |

### Status Display

**Managed context** (found in saved contexts):
```
✓ Active Context: prod
  Name: Production Environment
  Subscription: MyCompany Production
  Username: user@mycompany.com
```

**Unmanaged context** (not in saved contexts):
```
⚠ Active Azure session (not managed by azctx)
  Subscription: MyCompany Production (12345678-1234-1234-1234-123456789abc)
  Username: user@mycompany.com

  Use 'azctx add' to save this context.
```

**No active session**:
```
✗ No active Azure session
  Run 'az login' to authenticate.
```

## Relationships

No relationships - single entity model. Contexts are independent, flat list.

## Migration/Versioning

**Current version**: 1.0 (initial)

**Future considerations**:
- If schema changes, add `version` field to YAML root
- Implement migration logic in storage service
- For POC, breaking changes acceptable (users can re-add contexts)

**No backward compatibility required** for v1.0 POC per constitution.
