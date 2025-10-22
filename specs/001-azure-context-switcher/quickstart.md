# Quickstart Guide: Azure CLI Account Context Switcher

**Feature**: 001-azure-context-switcher
**Date**: October 17, 2025
**For**: Developers and implementers

## Overview

This quickstart provides step-by-step guidance for implementing the Azure CLI context switcher tool. Follow this guide to build the feature from scratch.

## Prerequisites

Before starting implementation:

1. **Azure CLI installed**: Verify with `az --version`
2. **Python 3.11+**: Verify with `python --version`
3. **UV installed**: Verify with `uv --version`
4. **Active Azure session**: Run `az login` and verify with `az account show`

## Project Setup

### Step 1: Add Dependencies

Add required packages to `pyproject.toml`:

```bash
uv add typer rich questionary pyyaml
```

This adds:
- `typer` - CLI framework
- `rich` - Terminal formatting and output
- `questionary` - Interactive prompts
- `pyyaml` - YAML file handling

### Step 2: Create Module Structure

Create the following directory structure:

```bash
# Create directories
New-Item -ItemType Directory -Path src/models -Force
New-Item -ItemType Directory -Path src/services -Force
New-Item -ItemType Directory -Path src/utils -Force

# Create __init__.py files
New-Item -ItemType File -Path src/models/__init__.py -Force
New-Item -ItemType File -Path src/services/__init__.py -Force
New-Item -ItemType File -Path src/utils/__init__.py -Force
```

Final structure:
```
src/
├── __init__.py
├── cli.py                    # CLI entry point (create this)
├── models/
│   ├── __init__.py
│   └── context.py            # Context dataclass (create this)
├── services/
│   ├── __init__.py
│   ├── azure_cli.py          # Azure CLI integration (create this)
│   ├── storage.py            # YAML file operations (create this)
│   └── context_manager.py    # Business logic (create this)
└── utils/
    ├── __init__.py
    ├── errors.py             # Custom exceptions (create this)
    └── console.py            # Rich console instance (create this)
```

## Implementation Order

Follow this sequence to build the feature incrementally and test as you go.

### Phase 1: Foundation (30 minutes)

**Goal**: Set up basic structure and utilities

1. **Create `src/utils/console.py`**:
   - Initialize Rich Console instance
   - Export for use across modules
   - Add helper functions for success/error messages

2. **Create `src/utils/errors.py`**:
   - Define custom exception classes:
     - `AzureCliNotFoundError`
     - `NoActiveSessionError`
     - `DuplicateContextError`
     - `ContextNotFoundError`
     - `StorageError`

3. **Create `src/models/context.py`**:
   - Define `Context` dataclass with all fields
   - Add validation methods
   - Add to_dict/from_dict for YAML serialization

**Test**: Import modules in Python REPL to verify no syntax errors

### Phase 2: Storage Layer (45 minutes)

**Goal**: Implement YAML file read/write operations

4. **Create `src/services/storage.py`**:
   - `get_storage_path()` - Return Path to `~/.azctx/contexts.yaml`
   - `load_contexts()` - Read and parse YAML file, return list of Context objects
   - `save_contexts(contexts)` - Write list of Context objects to YAML file
   - `add_context(context)` - Append new context
   - `delete_context(context_id)` - Remove context by ID
   - Handle file not found, directory creation, corrupted YAML

**Test**: Manual test in REPL:
```python
from src.services.storage import load_contexts, save_contexts, add_context
from src.models.context import Context
from datetime import datetime

# Test add
ctx = Context(
    context_id="test",
    context_name="Test Context",
    subscription_id="test-sub-id",
    subscription_name="Test Subscription",
    tenant_id="test-tenant-id",
    tenant_name="test.onmicrosoft.com",
    username="test@example.com",
    created_at=datetime.now()
)
add_context(ctx)

# Test load
contexts = load_contexts()
print(contexts)

# Test delete
delete_context("test")
```

### Phase 3: Azure CLI Integration (45 minutes)

**Goal**: Execute Azure CLI commands and parse results

5. **Create `src/services/azure_cli.py`**:
   - `check_azure_cli_installed()` - Verify `az` command exists
   - `get_current_account()` - Execute `az account show`, return dict or None
   - `set_account(subscription_id)` - Execute `az account set`, return success bool
   - Handle subprocess errors, timeouts, JSON parsing

**Test**: Manual test in terminal (with active Azure session):
```python
from src.services.azure_cli import check_azure_cli_installed, get_current_account, set_account

# Test check
print(check_azure_cli_installed())  # Should be True

# Test get
account = get_current_account()
print(account)  # Should show current Azure account

# Test set (use your subscription ID)
result = set_account(account['id'])
print(result)  # Should be True
```

### Phase 4: Business Logic (1 hour)

**Goal**: Implement core context management operations

6. **Create `src/services/context_manager.py`**:
   - `add_context_interactive()` - Orchestrate add command
     - Get current Azure account
     - Prompt for context_name and context_id
     - Validate uniqueness
     - Save to storage
     - Return success/error

   - `switch_context_interactive()` - Orchestrate switch command
     - Load contexts
     - Show interactive list
     - Execute Azure CLI set
     - Verify and return result

   - `list_contexts(verbose)` - Return formatted list of contexts

   - `get_status(verbose)` - Get current context status
     - Get current Azure account
     - Match against saved contexts
     - Return formatted status

   - `delete_context_interactive()` - Orchestrate delete command
     - Load contexts
     - Show interactive list
     - Confirm deletion
     - Delete from storage

**Test**: Manual test each function in isolation before integrating with CLI

### Phase 5: CLI Commands (1 hour)

**Goal**: Wire up Typer commands to business logic

7. **Create `src/cli.py`**:
   - Initialize Typer app: `app = typer.Typer()`

   - Define commands:
     - `@app.command()` for `add`
     - `@app.command()` for `switch`
     - `@app.command()` for `list`
     - `@app.command()` for `status`
     - `@app.command()` for `delete`

   - Add `if __name__ == "__main__":` block to run app

   - Each command:
     - Add docstring (becomes help text)
     - Call appropriate function from context_manager
     - Handle errors and display with Rich
     - Set exit codes appropriately

**Test**: Run each command from terminal:
```bash
uv run python -m src.cli --help
uv run python -m src.cli add
uv run python -m src.cli list
uv run python -m src.cli status
uv run python -m src.cli switch
uv run python -m src.cli delete
```

### Phase 6: Polish & Error Handling (30 minutes)

**Goal**: Improve UX and handle edge cases

8. **Enhance error messages**:
   - Use Rich Panels for all errors
   - Add helpful suggestions in error messages
   - Test all error scenarios:
     - Azure CLI not installed
     - No active Azure session
     - Duplicate context ID
     - Empty context list
     - Corrupted YAML file

9. **Add color and formatting**:
   - Success messages in green with ✓
   - Errors in red with ✗
   - Warnings in yellow with ⚠
   - Tables for list command
   - Panels for status command

10. **Keyboard interrupt handling**:
    - Wrap interactive prompts in try/except
    - Catch `KeyboardInterrupt` and exit gracefully
    - Display "Cancelled" message

## Manual Testing Checklist

Test these scenarios manually before considering the feature complete:

### Basic Flow

- [ ] Add first context (empty state)
- [ ] List contexts (should show one)
- [ ] Add second context
- [ ] Switch between contexts
- [ ] Check status (should match switched context)
- [ ] List with --verbose flag
- [ ] Status with --verbose flag
- [ ] Delete one context
- [ ] Verify deletion in list

### Edge Cases

- [ ] Try to add duplicate context_id (should reject)
- [ ] Try to add with invalid context_id format (should reject)
- [ ] Try to switch with no contexts (should error)
- [ ] Try to switch with only one context (should inform)
- [ ] Try to delete with no contexts (should error)
- [ ] Cancel add with Ctrl+C (should exit cleanly)
- [ ] Cancel switch with Esc (should exit cleanly)
- [ ] Cancel delete with Esc (should exit cleanly)
- [ ] Run commands with no Azure session (should error)
- [ ] Delete contexts.yaml and run list (should show empty)
- [ ] Manually corrupt contexts.yaml (should handle gracefully)

### Cross-Platform (if possible)

- [ ] Test on Windows
- [ ] Test on Linux (if available)
- [ ] Test on macOS (if available)
- [ ] Verify path handling works on each platform

## Configuration

### Entry Point

Add to `pyproject.toml` for easy execution:

```toml
[project.scripts]
azctx = "src.cli:app"
```

After adding, install in editable mode:
```bash
uv pip install -e .
```

Now you can run:
```bash
azctx --help
azctx add
# ... etc
```

### Alternative: Direct Python Module

Without scripts entry, use:
```bash
uv run python -m src.cli [command]
```

## Performance Verification

Measure command execution time:

```bash
# PowerShell
Measure-Command { uv run python -m src.cli list }
Measure-Command { uv run python -m src.cli status }
```

All commands should complete < 2 seconds (excluding user input time).

## Documentation

After implementation, update `README.md` with:

1. **What it does**: One-paragraph description
2. **Installation**: How to set up (UV commands)
3. **Usage**: Examples of each command
4. **Requirements**: Azure CLI, Python version
5. **Configuration**: Where contexts are stored

## Troubleshooting

Common issues during development:

**Import errors**:
- Ensure all `__init__.py` files exist
- Run from repo root, not from `src/`

**Azure CLI not found**:
- Verify `az --version` works in terminal
- Check PATH includes Azure CLI

**YAML parse errors**:
- Use `yaml.safe_load()` not `yaml.load()`
- Handle empty file case (returns None)

**Questionary not displaying**:
- Ensure terminal supports interactive input
- Test in standard terminal, not IDE integrated terminal if issues

**Rich formatting not showing**:
- Check terminal supports ANSI colors
- Test with `rich.print("[green]test[/green]")`

## Next Steps

After completing implementation:

1. Run full manual test checklist
2. Update README.md with usage examples
3. Create release (tag version)
4. Test installation from scratch in clean environment
5. Consider future enhancements (see Out of Scope in spec)

## Estimated Time

- **Total implementation time**: 4-5 hours
- **Testing time**: 1 hour
- **Documentation time**: 30 minutes
- **Total**: 5.5-6.5 hours for complete feature

Follow phases in order for best results. Test each phase before moving to the next.

---

## Phase 7: Building and Distributing Executables (User Story 6)

**Goal**: Create standalone Windows executable and automate releases

### Step 1: Add PyInstaller Dependency

Add PyInstaller as a development dependency:

```bash
uv add --dev pyinstaller
```

This adds to `[tool.uv.dev-dependencies]` in `pyproject.toml`.

### Step 2: Local Executable Build

Build the standalone executable locally:

```bash
# Build single-file executable
uv run pyinstaller --onefile src/cli.py --name azctx

# Output: dist/azctx.exe
```

**What happens**:
- PyInstaller analyzes imports and bundles Python runtime
- Creates `dist/azctx.exe` (typically 20-30MB)
- Creates `build/` folder with build artifacts
- May generate `azctx.spec` file for customization

### Step 3: Test Local Build

Test the executable works without Python:

```bash
# Test from dist folder
.\dist\azctx.exe --help
.\dist\azctx.exe status
.\dist\azctx.exe list

# Copy to PATH location for system-wide access
Copy-Item .\dist\azctx.exe -Destination C:\Windows\System32\
azctx --help  # Should work from any directory
```

**Verification**:
- Commands work identically to `uv run python -m src.cli`
- No Python environment required
- Executable size is under 50MB

### Step 4: Update .gitignore

Add build artifacts to `.gitignore`:

```gitignore
# PyInstaller build artifacts
dist/
build/
*.spec
```

### Step 5: Create GitHub Actions Workflow

Create `.github/workflows/release.yml`:

```yaml
name: Build and Release

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    name: Build Windows Executable
    runs-on: windows-latest

    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Set up Python
        uses: actions/setup-python@v5
        with:
          python-version: '3.11'

      - name: Install UV
        run: pip install uv

      - name: Install dependencies
        run: uv sync --all-extras --dev

      - name: Build executable
        run: uv run pyinstaller --onefile src/cli.py --name azctx

      - name: Check executable size
        run: |
          $size = (Get-Item dist/azctx.exe).Length / 1MB
          Write-Output "Executable size: $size MB"
          if ($size -gt 50) {
            Write-Error "Executable exceeds 50MB limit"
            exit 1
          }

      - name: Package executable
        run: |
          Compress-Archive -Path dist/azctx.exe -DestinationPath azctx-${{ github.ref_name }}-windows.zip

      - name: Create GitHub Release
        uses: softprops/action-gh-release@v1
        with:
          files: azctx-${{ github.ref_name }}-windows.zip
          generate_release_notes: true
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### Step 6: Create and Push Version Tag

Trigger automated build by creating a version tag:

```bash
# Ensure all changes are committed
git add .
git commit -m "Implement User Story 6: Executable distribution"

# Create annotated version tag
git tag -a v1.0.0 -m "Release v1.0.0 - Initial stable release"

# Push tag to trigger GitHub Actions
git push origin v1.0.0
```

**What happens**:
1. GitHub Actions detects tag push
2. Workflow runs on windows-latest runner
3. PyInstaller builds `azctx.exe`
4. Size is verified (must be <50MB)
5. Executable packaged as `azctx-v1.0.0-windows.zip`
6. GitHub Release automatically created
7. Zip file attached to release as downloadable asset

### Step 7: Verify GitHub Release

1. Go to repository on GitHub
2. Navigate to **Releases** section
3. Verify release `v1.0.0` was created
4. Download `azctx-v1.0.0-windows.zip`
5. Extract and test `azctx.exe`

### Step 8: Update Documentation

Update `README.md` with installation from release:

```markdown
## Installation

### Option 1: From Release (Windows)

1. Download latest release: [azctx-v1.0.0-windows.zip](https://github.com/youruser/azctx/releases)
2. Extract `azctx.exe`
3. Place in a directory in your PATH (e.g., `C:\Windows\System32` or `C:\Users\YourName\bin`)
4. Run `azctx --help` to verify

### Option 2: From Source (Cross-platform)

1. Clone repository: `git clone https://github.com/youruser/azctx`
2. Install dependencies: `uv sync`
3. Run with: `uv run python -m src.cli [command]`

## Building Executable Locally

Developers can build the executable locally:

```bash
# Install dev dependencies
uv sync --dev

# Build executable
uv run pyinstaller --onefile src/cli.py --name azctx

# Test
.\dist\azctx.exe --help
```

## Creating Releases

To create a new release:

1. Update version in code if applicable
2. Create annotated tag: `git tag -a v1.1.0 -m "Release v1.1.0"`
3. Push tag: `git push origin v1.1.0`
4. GitHub Actions automatically builds and publishes release
```

### Troubleshooting Build Issues

**PyInstaller import errors**:
- If hidden imports are missing, create `azctx.spec` and add to `hiddenimports` list
- Common additions: `['yaml', 'pkg_resources.py2_warn']`

**Executable too large**:
- Add exclusions to `.spec` file: `excludes=['tkinter', 'matplotlib', 'IPython']`
- Enable UPX compression: `upx=True` in EXE() section

**GitHub Actions fails**:
- Check workflow logs in Actions tab
- Verify UV sync completes successfully
- Ensure all dependencies are in pyproject.toml
- Test build locally first

**Executable doesn't run**:
- Test on clean Windows VM without Python
- Check for missing DLLs or dependencies
- Verify console mode is enabled in .spec

### Success Criteria Verification

- ✅ **SC-009**: Local build completes without errors
- ✅ **SC-010**: GitHub Actions completes within 10 minutes
- ✅ **SC-011**: Executable is under 50MB
- ✅ **FR-020**: Can build with `uv run pyinstaller --onefile src/cli.py`
- ✅ **FR-021**: Executable is named `azctx.exe`
- ✅ **FR-022**: Workflow triggers on version tags
- ✅ **FR-023**: Creates versioned zip file
- ✅ **FR-024**: Auto-creates GitHub Release
- ✅ **FR-025**: Executable works without Python installed
- ✅ **FR-026**: Build instructions documented

## Updated Estimated Time

- **Total implementation time (US1-5)**: 4-5 hours
- **Executable setup (US6)**: 1 hour
- **Testing time**: 1.5 hours
- **Documentation time**: 45 minutes
- **Total**: 7-8 hours for complete feature including distribution

