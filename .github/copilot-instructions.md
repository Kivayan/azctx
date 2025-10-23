# azctx Development Guidelines

**Project**: Azure CLI Account Context Switcher - A fast, convenient CLI tool for managing and switching between multiple Azure CLI account contexts with friendly names.

**Nature**: Personal POC project prioritizing simplicity, speed, and ease of use.

## Core Principles

### Simplicity First
- Start with the most straightforward solution - no over-engineering
- YAGNI (You Aren't Gonna Need It) strictly enforced
- Code must be readable and self-documenting
- Minimal configuration always preferred

### Modular Architecture
- Every feature is a self-contained component with single responsibility
- Explicit and minimal dependencies between modules
- No tight coupling - favor composition over inheritance
- Project structure: `src/` with `models/`, `services/`, `utils/`, `cli/`

### CLI-Centric Design
- All functionality accessible via fast, readable CLI commands
- Human-readable output by default
- Support both interactive and non-interactive modes
- Clear help text and error messages required

### UV-Managed Workflow
- **All** dependency management via UV exclusively
- No mixing of pip, poetry, or other package managers
- Virtual environments and script execution via UV
- Dependencies defined in `pyproject.toml`

### Pragmatic Quality
- **No formal test suite required** - manual testing is acceptable for POC
- Code must work correctly for intended use cases
- Critical bugs fixed before new features
- Documentation only for non-obvious functionality
- Type hints where they improve clarity (not mandatory)

## Technology Stack

**Core**:
- Python 3.11+ (managed via UV)
- Typer (CLI framework and argument parsing)
- Questionary (interactive prompts and selections)
- Rich (terminal output formatting)
- PyYAML (context storage)

**Build & Distribution**:
- PyInstaller (standalone executable builds)
- GitHub Actions (automated release workflow on version tags)

**Platform**: Cross-platform (Windows, macOS, Linux) with primary Windows focus for .exe distribution

## Development Standards

### Code Organization
```
src/
├── cli.py           # CLI entry point
├── models/          # Data models (Context entity)
├── services/        # Business logic (context_manager, azure_cli, storage)
└── utils/           # Utility functions (console, errors)
```

- Flat module structure preferred over deep hierarchies
- No `tests/` directory - manual testing only
- Configuration in project root or `.config/`

### Error Handling
- Catch errors and provide clear, actionable messages using Rich formatting
- Meaningful exit codes (0 = success, non-zero = failure)
- Hide stack traces unless debug mode enabled

### Performance
- "Fast enough" is sufficient - no premature optimization
- All operations should complete within 2 seconds on standard hardware
- Prioritize speed and convenience in user experience

### Distribution
- Local build: `uv run pyinstaller --onefile src/cli.py`
- Automated builds on version tags (v1.0.0, etc.) via GitHub Actions
- Release artifacts: `azctx-{version}-windows.zip` containing `azctx.exe`

## What NOT to Do

❌ Don't create formal test suites (manual testing sufficient for POC)
❌ Don't over-engineer solutions - simplest approach wins
❌ Don't use package managers other than UV
❌ Don't create deep module hierarchies
❌ Don't add features before fixing critical bugs
❌ Don't optimize prematurely
❌ Don't write extensive documentation for obvious functionality

## Common Commands

- `uv sync` - Install/update dependencies
- `uv run azctx` - Run the CLI application
- `uv run pyinstaller --onefile src/cli.py` - Build standalone executable

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
