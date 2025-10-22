# Implementation Plan: [FEATURE]

**Branch**: `[###-feature-name]` | **Date**: [DATE] | **Spec**: [link]
**Input**: Feature specification from `/specs/[###-feature-name]/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

[Extract from feature spec: primary requirement + technical approach from research]

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: Python 3.11+ (managed via UV)
**Primary Dependencies**: Typer (CLI), Questionary (prompts), Rich (output)
**Package Manager**: UV (all dependency and environment management)
**Testing**: Manual testing (no formal test suite per constitution)
**Target Platform**: Cross-platform (Windows, macOS, Linux)
**Project Type**: Single Python project with CLI interface
**Performance Goals**: "Fast enough" for CLI operations (no specific targets)
**Constraints**: Simplicity first - avoid over-engineering
**Scale/Scope**: Personal POC - optimize for velocity and experimentation

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Principle I - Simplicity First**:
- [ ] Feature uses the simplest possible implementation approach
- [ ] No premature optimization or over-engineering present
- [ ] Code will be readable and self-documenting
- [ ] Configuration is minimal and intuitive

**Principle II - Modular Architecture**:
- [ ] Feature is designed as self-contained, modular component(s)
- [ ] Module responsibilities are clear and single-purpose
- [ ] Dependencies between modules are explicit and minimal
- [ ] No tight coupling introduced

**Principle III - CLI-Centric Design**:
- [ ] All functionality accessible via CLI commands
- [ ] Uses Typer for command structure
- [ ] Uses Questionary for interactive prompts (if needed)
- [ ] Uses Rich for terminal output formatting
- [ ] Commands have clear help text
- [ ] Supports both interactive and non-interactive modes (where applicable)

**Principle IV - UV-Managed Workflow**:
- [ ] All new dependencies will be added via UV
- [ ] No alternative package managers used
- [ ] Dependencies listed in pyproject.toml

**Principle V - Pragmatic Quality**:
- [ ] Feature includes documentation for non-obvious functionality
- [ ] Manual testing approach is defined
- [ ] No formal test suite required (personal POC)
- [ ] Type hints used where they improve clarity

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

*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |

