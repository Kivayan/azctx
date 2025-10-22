<!--
=== SYNC IMPACT REPORT ===
Version Change: Initial → 1.0.0
Modified Principles: N/A (initial version)
Added Sections:
  - Core Principles (5 principles)
  - Technology Stack
  - Development Standards
  - Governance
Removed Sections: N/A (initial version)
Templates Requiring Updates:
  ✅ plan-template.md - Updated constitution check gates
  ✅ spec-template.md - Aligned with simplicity and modularity principles
  ✅ tasks-template.md - Removed test requirements, aligned with POC nature
Follow-up TODOs: None
=== END REPORT ===
-->

# azctx Constitution

## Core Principles

### I. Simplicity First

This is a personal POC project. Every decision MUST favor simplicity over complexity.

- Start with the most straightforward solution
- Avoid premature optimization
- No over-engineering - YAGNI (You Aren't Gonna Need It) strictly enforced
- Code MUST be readable and self-documenting
- Configuration MUST be minimal and intuitive

**Rationale**: Personal projects benefit from quick iteration and clarity. Complex
solutions increase maintenance burden and slow down experimentation.

### II. Modular Architecture

Every feature MUST be implemented as a modular, self-contained component.

- Modules MUST have clear, single responsibilities
- Dependencies between modules MUST be explicit and minimal
- Each module MUST be independently understandable
- Modules SHOULD be reusable in different contexts
- Tight coupling is prohibited - favor composition over inheritance

**Rationale**: Modularity enables easier debugging, refactoring, and feature additions.
It also allows parts of the codebase to evolve independently.

### III. CLI-Centric Design

All functionality MUST be accessible through a fast, readable, and easy-to-use CLI.

- Use Typer for command structure and argument parsing
- Use Questionary for interactive prompts and user input
- Use Rich for terminal output formatting and progress indication
- Commands MUST provide clear help text and usage examples
- Output MUST be human-readable by default, with optional structured formats (JSON)
- All operations MUST support both interactive and non-interactive modes

**Rationale**: CLI-first design ensures the tool is scriptable, testable via command
invocation, and provides immediate user feedback.

### IV. UV-Managed Workflow

UV MUST be used for all project management tasks.

- Dependency management via UV exclusively
- Virtual environment creation and activation via UV
- Package installation and updates via UV
- Script execution via UV run where applicable
- No mixing of pip, poetry, or other package managers

**Rationale**: UV provides fast, reliable dependency resolution and environment
management. Standardizing on one tool eliminates tooling complexity.

### V. Pragmatic Quality

Quality standards MUST be appropriate for a personal POC.

- No formal test suite required (manual testing acceptable)
- Code MUST work correctly for intended use cases
- Critical bugs MUST be fixed before adding new features
- Documentation MUST exist for non-obvious functionality
- Type hints SHOULD be used where they improve clarity
- Linting and formatting tools MAY be used but are not mandatory

**Rationale**: For a POC, velocity and experimentation matter more than comprehensive
test coverage. Quality comes from clear code and practical validation, not ceremony.

## Technology Stack

**Language**: Python 3.11+
**Project Manager**: UV (all dependency and environment management)
**CLI Framework**: Typer (command structure and argument parsing)
**User Input**: Questionary (interactive prompts and selections)
**Terminal Output**: Rich (formatting, tables, progress bars, syntax highlighting)
**Platform**: Cross-platform (Windows, macOS, Linux)

All dependencies MUST be managed via `pyproject.toml` and installed using UV.

## Development Standards

### Code Organization

- Single project structure: `src/` for source code
- Flat module structure preferred over deep hierarchies
- Entry point MUST be clearly identified (e.g., `src/main.py` or `src/cli.py`)
- Configuration files SHOULD live in project root or `.config/` directory

### Documentation

- README.md MUST explain what the project does and how to use it
- Each module SHOULD have a docstring explaining its purpose
- Complex functions SHOULD have docstrings with parameter descriptions
- No separate documentation site required

### Error Handling

- Errors MUST be caught and reported with clear, actionable messages
- Use Rich to format error output for readability
- Exit codes MUST be meaningful (0 for success, non-zero for failures)
- Stack traces MAY be hidden from users unless debug mode is enabled

### Performance

- No specific performance targets - "fast enough" is sufficient
- Avoid obviously inefficient algorithms
- Lazy loading and caching MAY be used if they simplify user experience

## Governance

### Amendment Process

- Constitution can be amended at any time by updating this file
- Version MUST be bumped according to semantic versioning rules
- Changes MUST be documented in the Sync Impact Report comment
- Dependent templates MUST be reviewed and updated if affected

### Versioning Policy

- **MAJOR**: Principles removed or fundamentally changed
- **MINOR**: New principles added or existing ones materially expanded
- **PATCH**: Clarifications, typo fixes, non-semantic refinements

### Compliance

- This is a personal project - self-review is sufficient
- Constitution serves as a reference, not a bureaucratic checklist
- Deviations are allowed if justified by practical constraints

**Version**: 1.0.0 | **Ratified**: 2025-10-17 | **Last Amended**: 2025-10-17
