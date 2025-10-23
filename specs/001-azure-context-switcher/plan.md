# Implementation Plan: Azure CLI Account Context Switcher - User Story 6

**Branch**: `001-azure-context-switcher` | **Date**: 2025-10-17 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-azure-context-switcher/spec.md` - User Story 6

**Note**: This plan focuses on User Story 6: Build and Distribute Standalone Executable. User Stories 1-5 are already implemented (96% complete).

## Summary

Add standalone Windows executable distribution using PyInstaller with automated GitHub Actions workflow. Enables users to install azctx without Python by downloading a versioned .exe from GitHub Releases. Local builds via `uv run pyinstaller --onefile src/cli.py`, automated builds triggered on version tags (v1.0.0, etc.), packaged as `azctx-{version}-windows.zip` with auto-created GitHub Release.

## Technical Context

**Language/Version**: Python 3.11+ (managed via UV)
**Primary Dependencies**: Typer (CLI), Questionary (prompts), Rich (output), PyYAML (storage)
**Build Tool**: PyInstaller (standalone executable generation)
**Package Manager**: UV (all dependency and environment management)
**Testing**: Manual testing (no formal test suite per constitution)
**Target Platform**: Windows (primary), cross-platform Python source code
**Distribution**: GitHub Releases with automated builds via GitHub Actions
**Executable Format**: Single-file .exe (PyInstaller --onefile)
**Package Format**: Versioned zip files (azctx-{version}-windows.zip)
**Trigger Mechanism**: Git version tags (1.0.0, 1.1.0, etc.)
**Project Type**: Single Python project with CLI interface
**Performance Goals**: Executable under 50MB, builds complete within 10 minutes
**Constraints**: Simplicity first - minimal build configuration
**Scale/Scope**: Personal POC with wider distribution capability

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Principle I - Simplicity First**:

- [x] Feature uses the simplest possible implementation approach (PyInstaller is standard, well-documented)
- [x] No premature optimization or over-engineering present (single .spec file, minimal config)
- [x] Code will be readable and self-documenting (build script + GitHub Actions YAML)
- [x] Configuration is minimal and intuitive (version tags trigger builds automatically)

**Principle II - Modular Architecture**:

- [x] Feature is designed as self-contained, modular component(s) (build separate from core CLI)
- [x] Module responsibilities are clear and single-purpose (PyInstaller only for builds)
- [x] Dependencies between modules are explicit and minimal (no runtime impact on core app)
- [x] No tight coupling introduced (build artifacts don't affect source code)

**Principle III - CLI-Centric Design**:

- [x] All functionality accessible via CLI commands (local build: uv run pyinstaller)
- [x] Uses Typer for command structure (existing CLI unchanged)
- [x] Uses Questionary for interactive prompts (not applicable for builds)
- [x] Uses Rich for terminal output formatting (existing CLI unchanged)
- [x] Commands have clear help text (build documentation in README)
- [x] Supports both interactive and non-interactive modes (builds are non-interactive)

**Principle IV - UV-Managed Workflow**:

- [x] All new dependencies will be added via UV (PyInstaller as dev dependency)
- [x] No alternative package managers used (UV only)
- [x] Dependencies listed in pyproject.toml (PyInstaller in [tool.uv.dev-dependencies])

**Principle V - Pragmatic Quality**:

- [x] Feature includes documentation for non-obvious functionality (README build instructions)
- [x] Manual testing approach is defined (local .exe testing, GitHub Actions validation)
- [x] No formal test suite required (personal POC - manual validation sufficient)
- [x] Type hints used where they improve clarity (not applicable for YAML/config files)

**GATE STATUS**: ✅ PASSED - All constitutional principles satisfied

**POST-PHASE 1 RE-CHECK**: ✅ PASSED - Design artifacts confirm simplicity-first approach with minimal configuration

## Project Structure

### Documentation (this feature)

```
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```
# azctx project structure (Python single project)
src/
├── __init__.py
├── main.py or cli.py      # CLI entry point
├── models/                 # Data models and entities
├── services/               # Business logic
├── cli/                    # CLI command modules (Typer)
└── utils/                  # Utility functions

pyproject.toml              # UV-managed dependencies
README.md                   # Project documentation
.config/                    # Optional configuration directory
```

**Structure Decision**: Using Python single project layout with modular structure.
CLI commands organized in `src/cli/`, models in `src/models/`, business logic in
`src/services/`. No tests directory per constitution (manual testing only).

## Complexity Tracking

**No constitutional violations** - all principles satisfied. PyInstaller is the standard tool for Python executable builds, GitHub Actions is the standard CI/CD platform, and implementation follows simplicity-first approach with minimal configuration.

---

## Planning Completion Summary

### Phase 0: Research ✅ COMPLETE

**Output**: `research.md` (updated with User Story 6 research)

Completed research on:

1. PyInstaller configuration for Typer/Questionary/Rich applications
2. GitHub Actions workflow structure and triggers
3. .spec file customization approach
4. Executable size optimization strategies
5. Local development build workflow

**Key Decisions**:

- Build Tool: PyInstaller with `--onefile` mode
- Trigger: Version tags (`v*` pattern)
- Distribution: Automated GitHub Releases with zip packaging
- Configuration: Command-line first, .spec file if needed
- Size Target: <50MB (expected 20-30MB without optimization)

### Phase 1: Design & Contracts ✅ COMPLETE

**Outputs**:

1. `data-model.md` - No changes needed (no new entities for build process)
2. `contracts/github-actions-release.yml` - GitHub Actions workflow contract
3. `quickstart.md` - Updated with Phase 7 (build and distribution instructions)
4. `.github/copilot-instructions.md` - Updated with PyInstaller and GitHub Actions

**Artifacts Created**:

- GitHub Actions workflow contract with full step-by-step specification
- Quickstart guide section covering local builds, testing, releases
- Documentation for troubleshooting build issues
- Success criteria verification checklist

### Constitution Check

**Pre-Phase 0**: ✅ PASSED
**Post-Phase 1**: ✅ PASSED

All constitutional principles satisfied. Implementation approach aligns with:

- Simplicity First: Minimal configuration, standard tools
- Modular Architecture: Build process separate from core application
- CLI-Centric Design: Simple command-line interface for builds
- UV-Managed Workflow: PyInstaller added as dev dependency via UV
- Pragmatic Quality: Manual testing approach, clear documentation

### Next Steps

**Command**: `/speckit.tasks` or run `.specify/scripts/powershell/setup-tasks.ps1`

This will generate `tasks.md` with detailed implementation tasks for User Story 6:

- Adding PyInstaller dependency to pyproject.toml
- Creating GitHub Actions workflow file
- Updating .gitignore
- Testing local builds
- Creating and pushing version tags
- Verifying automated releases
- Updating README with installation instructions

**Estimated Implementation Time**: 1-2 hours
**Estimated Testing Time**: 30 minutes
**Total**: 1.5-2.5 hours for User Story 6 completion
