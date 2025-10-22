# Research: Azure CLI Account Context Switcher

**Feature**: 001-azure-context-switcher
**Date**: October 17, 2025
**Status**: Complete

## Overview

This document contains research findings and technical decisions for implementing the Azure CLI context switching tool. All decisions prioritize simplicity, speed, and maintainability per project constitution.

## Technology Stack Research

### Decision 1: CLI Framework - Typer

**What was chosen**: Typer (v0.12+)

**Rationale**:
- Built on Click, provides excellent type safety and auto-completion
- Automatic help text generation from type hints and docstrings
- Native support for subcommands (azctx add, azctx switch, etc.)
- Rich integration for beautiful output
- Minimal boilerplate - aligns with simplicity principle

**Alternatives considered**:
- **Click**: More verbose, requires decorators for everything
- **argparse**: Standard library but too low-level, more code required
- **Fire**: Too magical, less explicit control over CLI structure

**Best practices**:
- Use `typer.Typer()` app instance for command grouping
- Leverage type hints for automatic validation
- Use `typer.Option()` and `typer.Argument()` for clear parameter definition
- Set `help` parameter on all commands for auto-generated documentation

### Decision 2: Interactive Prompts - Questionary

**What was chosen**: Questionary (v2.0+)

**Rationale**:
- Clean, simple API for interactive prompts
- Built on prompt_toolkit - robust and well-maintained
- Supports list selection with arrow key navigation (core requirement)
- Customizable styling
- Works cross-platform (Windows, Linux, macOS)

**Alternatives considered**:
- **PyInquirer**: Abandoned, last update 2019
- **prompt_toolkit directly**: Too low-level, more complex
- **simple-term-menu**: Linux/macOS only, no Windows support

**Best practices**:
- Use `questionary.select()` for list navigation (switch, delete commands)
- Use `questionary.text()` for text input (context names, IDs)
- Use `questionary.confirm()` for yes/no prompts (delete confirmation)
- Handle KeyboardInterrupt (Ctrl+C) gracefully
- Provide clear prompt messages and instructions

### Decision 3: Terminal Output - Rich

**What was chosen**: Rich (v13.0+)

**Rationale**:
- Beautiful, styled console output with minimal effort
- Table formatting perfect for `azctx list` command
- Color support for better UX (success=green, error=red)
- Progress indicators if needed for slower operations
- Syntax highlighting for potential JSON/YAML output

**Alternatives considered**:
- **colorama**: Too basic, manual color codes
- **tabulate**: Tables only, no rich formatting
- **blessed**: Lower-level, more complex

**Best practices**:
- Use `rich.console.Console()` for output management
- Use `rich.table.Table()` for structured data display
- Use `rich.panel.Panel()` for important messages/status
- Color conventions: green for success, red for errors, yellow for warnings, cyan for info
- Use `console.print()` instead of `print()` throughout

### Decision 4: Data Persistence - YAML

**What was chosen**: PyYAML (v6.0+) with simple file-based storage

**Rationale**:
- Human-readable format - users can manually edit if needed
- Simple serialization/deserialization
- Standard library support via PyYAML
- No database overhead for small dataset (< 100 contexts expected)
- Fast read/write for CLI performance requirements

**Alternatives considered**:
- **JSON**: Less human-readable, no comments support
- **SQLite**: Over-engineered for simple key-value storage
- **TOML**: Less common, similar to YAML but less flexible
- **Pickle**: Not human-readable, security concerns

**Best practices**:
- Store in `~/.azctx/contexts.yaml` (cross-platform home directory)
- Use `safe_load()` and `safe_dump()` for security
- Create directory if it doesn't exist
- Handle file corruption gracefully
- Keep file structure simple: list of dictionaries

**File structure**:
```yaml
contexts:
  - context_id: "prod"
    context_name: "Production Environment"
    subscription_id: "abc-123-def-456"
    subscription_name: "MyCompany Production"
    tenant_id: "xyz-789"
    tenant_name: "mycompany.onmicrosoft.com"
    username: "user@example.com"
    created_at: "2025-10-17T10:30:00"
  - context_id: "dev"
    context_name: "Development Environment"
    # ... more fields
```

### Decision 5: Azure CLI Integration - Subprocess

**What was chosen**: Python `subprocess.run()` with JSON output parsing

**Rationale**:
- Azure CLI already installed (dependency assumption)
- `az` commands support `--output json` for structured data
- Subprocess is standard library - no extra dependencies
- Simple to execute and parse results
- Error handling via return codes

**Alternatives considered**:
- **Azure SDK for Python**: Too heavy, requires authentication management
- **Direct API calls**: Bypasses Azure CLI, defeats the purpose
- **Shell scripting**: Less portable, harder to maintain

**Best practices**:
- Always use `--output json` flag with az commands
- Set timeout for subprocess calls (prevent hanging)
- Check return code before parsing output
- Use `subprocess.run()` with `capture_output=True`
- Handle cases where `az` CLI is not installed
- Parse JSON output with standard library `json` module

**Key commands**:
```bash
# Get current context
az account show --output json

# Set context by subscription ID
az account set --subscription <subscription_id>
```

## Performance Optimization

### Decision 6: Lazy Loading & Caching

**What was chosen**: Minimal caching, direct file reads

**Rationale**:
- YAML file read is fast enough (< 50ms for reasonable file size)
- Simplicity over premature optimization
- No stale data issues
- Meets < 2 second performance requirement easily

**Alternatives considered**:
- **In-memory cache**: Adds complexity, not needed for file-based storage
- **Background refresh**: Over-engineered for single-user tool

**Best practices**:
- Read YAML file on each command execution
- Cache only during single command execution if needed
- No persistent daemon or background processes

## Cross-Platform Considerations

### Decision 7: Path Handling

**What was chosen**: `pathlib.Path` for all file operations

**Rationale**:
- Cross-platform path handling (Windows vs Unix)
- Modern, object-oriented API
- Built-in methods for home directory (`Path.home()`)
- Better than `os.path` string manipulation

**Best practices**:
- Use `Path.home()` for user directory
- Use `/` operator for path joining
- Call `.mkdir(parents=True, exist_ok=True)` for directory creation
- Use `.exists()`, `.is_file()` for checks

**Storage paths**:
- Windows: `%USERPROFILE%\.azctx\contexts.yaml`
- Unix: `~/.azctx/contexts.yaml`
- Implementation: `Path.home() / ".azctx" / "contexts.yaml"`

### Decision 8: Error Handling Strategy

**What was chosen**: Rich-formatted error messages with exit codes

**Rationale**:
- Clear, actionable error messages improve UX
- Exit codes allow scripting (0 = success, 1+ = errors)
- Rich formatting makes errors easy to spot
- No stack traces in normal operation (only in debug mode)

**Best practices**:
- Catch specific exceptions, not bare `except`
- Format errors with Rich Panel in red
- Provide next steps in error messages
- Exit with appropriate code (0 success, 1 general error, 2 validation error)
- Log stack traces only if `--debug` flag is set

**Error categories**:
1. Azure CLI not found
2. No active Azure session
3. YAML file corrupted
4. Duplicate context ID
5. Context not found
6. Subprocess timeout

## Security Considerations

### Decision 9: No Credential Storage

**What was chosen**: Store only Azure metadata, no credentials

**Rationale**:
- Tool relies on Azure CLI authentication
- Never store tokens, passwords, or secrets
- Minimize security risk surface
- Align with Azure CLI's credential management

**Best practices**:
- Only store subscription/tenant IDs and names
- Rely on `az login` for authentication
- No custom credential handling
- Document that users must use `az login` separately

## Testing Strategy

### Decision 10: Manual Testing Approach

**What was chosen**: Manual testing with documented test scenarios

**Rationale**:
- Constitution allows manual testing for POC projects
- Automated tests for CLI tools require mocking complexity
- Fast iteration more valuable than test coverage
- Document test cases for regression checking

**Manual test scenarios**:
1. Add first context (empty state)
2. Add duplicate context ID (validation)
3. Switch between 2+ contexts
4. List contexts (verbose and non-verbose)
5. Delete context (with confirmation)
6. Check status on managed/unmanaged contexts
7. Azure CLI not installed scenario
8. Corrupted YAML file recovery
9. Cross-platform path handling (Windows/Linux)
10. Large context list (50+ items) performance

## Implementation Notes

### Module Structure

```
src/
├── __init__.py
├── cli.py                 # Entry point, Typer app definition
├── models/
│   ├── __init__.py
│   └── context.py         # Context dataclass
├── services/
│   ├── __init__.py
│   ├── azure_cli.py       # Azure CLI interaction (subprocess)
│   ├── storage.py         # YAML file read/write
│   └── context_manager.py # Business logic (add, switch, delete)
└── utils/
    ├── __init__.py
    ├── errors.py          # Custom exception classes
    └── console.py         # Rich console instance
```

### Data Flow

1. **Add Context**:
   - Execute `az account show` → parse JSON
   - Prompt for context_id and context_name (Questionary)
   - Validate uniqueness
   - Append to YAML file
   - Display success (Rich)

2. **Switch Context**:
   - Load contexts from YAML
   - Display interactive list (Questionary)
   - Execute `az account set --subscription <id>`
   - Verify switch with `az account show`
   - Display confirmation (Rich)

3. **List Contexts**:
   - Load contexts from YAML
   - Format as table (Rich Table)
   - Show subset or all fields based on --verbose

4. **Status**:
   - Execute `az account show`
   - Match against saved contexts
   - Display formatted output (Rich Panel)

5. **Delete Context**:
   - Load contexts from YAML
   - Display interactive list (Questionary)
   - Confirm deletion (Questionary)
   - Remove from list and save YAML
   - Display confirmation (Rich)

## Open Questions & Assumptions

**Assumptions made**:
- Users will have Azure CLI >= 2.0 installed
- Users authenticate via `az login` before using tool
- Context file will remain < 1MB (~ 1000 contexts max)
- Terminal supports UTF-8 and ANSI colors
- No concurrent access to contexts.yaml file

**No clarifications needed** - all requirements are clear from spec and user input.

## Dependencies Summary

**Required packages** (to be added via UV):
```toml
[project.dependencies]
typer = "^0.12.0"
questionary = "^2.0.0"
rich = "^13.0.0"
pyyaml = "^6.0.0"
```

**Standard library usage**:
- `subprocess` - Azure CLI execution
- `json` - Parse Azure CLI output
- `pathlib` - Cross-platform path handling
- `datetime` - Timestamp context creation
- `sys` - Exit codes

**Total new dependencies**: 4 packages (all align with constitution - simple, focused, well-maintained)

---

## User Story 6: Build and Distribute Standalone Executable (Added 2025-10-17)

### Decision 6: Build Tool - PyInstaller

**What was chosen**: PyInstaller with `--onefile` mode

**Rationale**:
- Most popular Python packaging tool (14k+ GitHub stars, actively maintained)
- Well-documented support for Typer, Questionary, Rich applications
- `--onefile` mode creates single executable (simpler distribution than `--onedir`)
- Cross-platform build capability (though focusing on Windows for this story)
- Simple command-line interface: `pyinstaller --onefile src/cli.py`
- Works seamlessly with UV as dev dependency

**Alternatives considered**:
- **cx_Freeze**: More complex configuration, less popular for single-file builds
- **Nuitka**: Produces smaller executables but significantly slower build times (10-30 minutes vs 1-2 minutes)
- **py2exe**: Windows-only, less active development than PyInstaller
- **Briefcase**: More suited for GUI applications, overkill for CLI tools

**Best practices**:
- Add PyInstaller to dev dependencies: `uv add --dev pyinstaller`
- Use `--name azctx` to control executable name
- Generate .spec file for reproducible builds if customization needed
- Exclude unnecessary modules to reduce executable size
- Test on clean Windows VM without Python installed

**Implementation details**:
```bash
# Add to pyproject.toml [tool.uv.dev-dependencies]
pyinstaller = ">=6.0.0"

# Build command (local development)
uv run pyinstaller --onefile src/cli.py --name azctx

# Output: dist/azctx.exe (typically 20-30MB for CLI apps with Rich/Typer)
```

### Decision 7: GitHub Actions Workflow Trigger

**What was chosen**: Version tags (`v*` pattern)

**Rationale**:
- Version tags (v1.0.0, v1.1.0) are semantic and standard for releases
- Separates development commits from release builds (no build on every push)
- Aligns with semantic versioning best practices
- GitHub automatically links tags to releases
- Tag name can be used in zip filename: `azctx-{tag}-windows.zip`

**Alternatives considered**:
- **Manual builds**: Error-prone, inconsistent, requires local Windows machine
- **Build on every commit**: Wasteful CI time, pollutes releases
- **Build on PR merge to main**: Decoupled from versioning
- **Manual workflow dispatch**: Extra step, defeats automation purpose

**Best practices**:
- Use semantic versioning: v1.0.0, v1.1.0, v2.0.0
- Create annotated tags: `git tag -a v1.0.0 -m "Release v1.0.0"`
- Document tagging process in README for contributors
- Use pre-release tags for testing: v1.0.0-rc1, v1.0.0-beta

**Implementation details**:
```yaml
# .github/workflows/release.yml
on:
  push:
    tags:
      - 'v*'  # Triggers on v1.0.0, v2.1.3, etc.
```

### Decision 8: Release Distribution Method

**What was chosen**: Automated GitHub Release creation with zip attachment

**Rationale**:
- GitHub Releases are the standard distribution method for GitHub projects
- Users expect to find downloads in Releases section
- `softprops/action-gh-release` action handles all complexity
- Automatically creates release from tag with zip file attached
- Discoverable via GitHub UI and API
- No need for separate hosting or CDN

**Alternatives considered**:
- **GitHub Packages**: More complex, requires authentication to download
- **External hosting (S3, Azure Blob)**: Unnecessary complexity for POC
- **Artifacts on Actions**: Not user-facing, 90-day retention limit
- **Release notes only**: Requires manual upload of .exe

**Best practices**:
- Package as zip file: `azctx-{version}-windows.zip` containing `azctx.exe`
- Include version in filename for clarity
- Auto-generate release notes from commits (GitHub feature)
- Mark pre-release versions appropriately

**Implementation details**:
```yaml
- name: Create Release
  uses: softprops/action-gh-release@v1
  with:
    files: azctx-${{ github.ref_name }}-windows.zip
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### Decision 9: PyInstaller Configuration Approach

**What was chosen**: Start with command-line flags, create .spec file if customization needed

**Rationale**:
- Command-line approach is simpler for straightforward builds
- Typer/Questionary/Rich don't require special PyInstaller configuration
- .spec file can be generated automatically and committed if needed
- Start simple, add complexity only when necessary (constitution principle)

**When .spec file is needed**:
- Hidden imports not auto-detected (e.g., PyYAML in some cases)
- Want to add application icon
- Need to embed version information
- Want to exclude specific modules to reduce size

**Best practices**:
- First try: `uv run pyinstaller --onefile src/cli.py --name azctx`
- If issues, generate .spec: PyInstaller creates azctx.spec automatically
- Commit .spec file to version control for reproducibility
- Document any manual .spec changes in comments

**Example .spec customization if needed**:
```python
# azctx.spec
a = Analysis(
    ['src/cli.py'],
    hiddenimports=['yaml'],  # If PyYAML not auto-detected
    excludes=['tkinter', 'matplotlib', 'IPython'],  # Reduce size
)
exe = EXE(..., upx=True, console=True, name='azctx')
```

### Decision 10: Executable Size Optimization

**What was chosen**: Use exclusions and optional UPX compression to stay under 50MB

**Rationale**:
- Base PyInstaller build for CLI app with Rich/Typer typically 20-30MB
- SC-011 requires <50MB, so base build should meet requirement
- Exclusions can save 5-10MB by removing unused stdlib modules
- UPX compression can reduce size by 40-60% with minimal performance impact
- Only optimize if approaching limit (simplicity principle)

**Size breakdown estimate**:
- Python runtime: ~8MB
- Typer + dependencies: ~3MB
- Questionary + prompt_toolkit: ~5MB
- Rich: ~2MB
- PyYAML: ~1MB
- Application code: <1MB
- **Total estimate**: 20-25MB (well under 50MB target)

**Best practices**:
- Monitor executable size in GitHub Actions logs
- Add size assertion to fail workflow if over 50MB
- Document size in release notes
- Consider UPX only if approaching limit (introduces slight startup delay)

**Implementation if needed**:
```python
# In .spec file:
excludes=['tkinter', 'matplotlib', 'IPython', 'notebook', 'jupyter']
exe = EXE(..., upx=True, ...)  # Only if needed
```

### Decision 11: Local Development Build Workflow

**What was chosen**: Simple one-line command documented in README

**Rationale**:
- UV already manages PyInstaller as dev dependency
- PyInstaller outputs to `dist/` by default
- Developers can test locally before pushing tags
- No need for build scripts or Make targets (Windows compatibility)

**Best practices**:
- Add `dist/` and `build/` to .gitignore
- Document in README under "Building Executable" section
- Include testing steps (run commands with built .exe)
- Explain difference between development (`uv run python -m src.cli`) and distribution (`.exe`)

**Developer workflow**:
```bash
# 1. Build executable locally
uv run pyinstaller --onefile src/cli.py --name azctx

# 2. Test the executable
./dist/azctx.exe --help
./dist/azctx.exe status
./dist/azctx.exe switch

# 3. If satisfied, create version tag and push
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# 4. GitHub Actions builds and creates release automatically
```

## Updated Dependencies Summary

**Updated dev dependencies** (to be added via UV):
```toml
[tool.uv.dev-dependencies]
pyinstaller = ">=6.0.0"
```

**New files to create**:
- `.github/workflows/release.yml` - GitHub Actions workflow
- `.gitignore` additions: `dist/`, `build/`, `*.spec` (if not using custom spec)

**Documentation updates needed**:
- README.md: Add "Building Executable" section
- README.md: Add "Releases" section explaining download/install
- README.md: Update installation instructions to include .exe option

## Implementation Risks and Mitigations

| Risk | Impact | Likelihood | Mitigation |
|------|--------|------------|-----------|
| PyInstaller doesn't detect all imports | Runtime errors in .exe | Medium | Add hidden imports to .spec, test thoroughly |
| Executable over 50MB | Fails SC-011 | Low | Use exclusions, monitor in CI |
| GitHub Actions fails on tag push | No release created | Low | Test with v0.1.0-test tag first |
| .exe doesn't work without Python | Defeats purpose | Low | Test on clean VM without Python |
| Build time exceeds 10 minutes | Fails SC-010 | Very Low | PyInstaller typically 1-2 min, use caching |
| Windows-only distribution | Limited platform support | Accepted | Document that source works cross-platform |

## Next Steps for Phase 1

1. Update data-model.md: No new entities (build process is external)
2. Create contracts/: GitHub Actions YAML workflow contract
3. Update quickstart.md: Add build/release instructions
4. Update agent context with PyInstaller technology

