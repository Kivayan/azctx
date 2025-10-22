![Python](https://img.shields.io/badge/python-3.11+-brightgreen)
![License](https://img.shields.io/badge/license-MIT-green)
![GitHub Tag](https://img.shields.io/github/v/tag/Kivayan/azctx)

# Azure CLI Context Manager (azctx)

**Azure CLI Account Context Switcher** - A fast, convenient CLI tool for managing and switching between multiple Azure CLI account contexts with friendly names.

## Overview

`azctx` is a personal proof-of-concept (POC) tool that wraps Azure CLI commands to provide a better user experience when working with multiple Azure subscriptions and accounts. Instead of remembering subscription IDs or account details, you can save contexts with memorable names and quickly switch between them.

### Key Features

- 🚀 **Fast context switching** - Switch between Azure accounts in seconds
- 🏷️ **Friendly names** - Use memorable names instead of subscription IDs
- 📋 **Easy management** - Add, list, delete, and check status of saved contexts
- 🎨 **Beautiful UI** - Interactive prompts with Rich terminal formatting
- ⚡ **Performance focused** - All operations complete in under 2 seconds

## Technology Stack

- **Python**: 3.11+
- **Package Manager**: [UV](https://github.com/astral-sh/uv) - Fast Python package installer and resolver
- **CLI Framework**: [Typer](https://typer.tiangolo.com/) - Modern CLI framework
- **User Input**: [Questionary](https://questionary.readthedocs.io/) - Interactive prompts
- **Terminal Output**: [Rich](https://rich.readthedocs.io/) - Beautiful terminal formatting
- **Data Storage**: YAML files in `~/.azctx/contexts.yaml`

## Prerequisites

- Python 3.11 or higher
- [UV package manager](https://github.com/astral-sh/uv)
- [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli) installed and configured

## Installation

### Option 1: From Release (Windows - Recommended)

Download the standalone executable - no Python installation required!

1. Go to [Releases](https://github.com/[your-username]/azctx/releases)
2. Download the latest `azctx-v*-windows.zip` file
3. Extract `azctx.exe` from the zip
4. Place the executable in a directory in your PATH:
   - **Quick option**: `C:\Windows\System32` (requires admin)
   - **User option**: `C:\Users\[YourUsername]\bin` (add to PATH if needed)
5. Open a new terminal and run: `azctx --help`

**Adding to PATH (if needed)**:

```powershell
# PowerShell (add user bin directory to PATH)
$userBin = "$env:USERPROFILE\bin"
New-Item -ItemType Directory -Path $userBin -Force
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";$userBin", "User")
```

## Quick Start

```bash
# 1. Login to Azure (if not already logged in)
az login

# 2. Add your current Azure context with a friendly name
azctx add

# 3. Login to another Azure account
az login

# 4. Add that context too
azctx add

# 5. Switch between contexts
azctx switch

# 6. Check which context is active
azctx status

# 7. List all saved contexts
azctx list
```

## Usage

### Add a New Context

Save your currently active Azure CLI session with a friendly name:

```bash
azctx add
```

You'll be prompted to:

1. Enter a friendly name (e.g., "Production", "Development")
2. Enter a short ID (e.g., "prod", "dev")

**Example:**

```
Current Azure Context
┌─────────────────────────────────────┐
│ Subscription: My Production Sub     │
│ Tenant: my-tenant-id                │
│ Account: user@example.com           │
└─────────────────────────────────────┘

Enter a friendly name: Production
Enter a short ID: prod

✓ Successfully Added Context: Production
```

### Switch Between Contexts

Use the interactive menu to choose a context:

```bash
azctx switch
```

**Example:**

```
? Select a context to switch to: (Use arrow keys)
❯ prod - Production
  dev - Development
  staging - Staging Environment
```

Use arrow keys to navigate, Enter to select, or Esc to cancel.

**Example:**

```
? Select a context to switch to:
❯ Production (prod)
  Development (dev)
  Testing (test)

✓ Successfully Switched to Production
```

### Check Active Context

Display the currently active Azure CLI context:

```bash
azctx status
```

Add `--verbose` or `-v` for detailed information including IDs and timestamps.

**Example:**

```
✓ Managed Context: Production
┌─────────────────────────────────────┐
│ Name: Production                    │
│ ID: prod                            │
│ Subscription: My Production Sub     │
│ Tenant: my-tenant-id                │
│ Account: user@example.com           │
└─────────────────────────────────────┘
```

### List All Contexts

Display all saved contexts in a table:

```bash
azctx list
```

Add `--verbose` or `-v` to see full details for each context.

**Example:**

```
           Saved Contexts (3)
┏━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━┓
┃ ID      ┃ Name              ┃
┡━━━━━━━━━╇━━━━━━━━━━━━━━━━━━━┩
│ prod    │ Production        │
│ dev     │ Development       │
│ test    │ Testing           │
└─────────┴───────────────────┘

Use 'azctx list --verbose' for detailed information
```

### Delete a Context

Interactively select and delete a saved context:

```bash
azctx delete
```

You'll be asked to confirm before deletion. **Note:** This only removes the saved context, your active Azure CLI session remains unchanged.

**Example:**

```
? Select a context to delete:
❯ Production (prod)
  Development (dev)
  Testing (test)

? Are you sure you want to delete context 'Production' (prod)? No

✓ Successfully Deleted Context: Production
Note: Your active Azure CLI session remains unchanged.
```

## Storage

Contexts are saved in a YAML file at:

- **Windows**: `%USERPROFILE%\.azctx\contexts.yaml`
- **macOS/Linux**: `~/.azctx/contexts.yaml`

The file is created automatically when you add your first context.

## Command Reference

| Command | Description | Options |
|---------|-------------|---------|
| `azctx add` | Add current Azure CLI context | None |
| `azctx switch` | Switch to a different saved context | None |
| `azctx status` | Show currently active context | `--verbose, -v` |
| `azctx list` | List all saved contexts | `--verbose, -v` |
| `azctx delete` | Delete a saved context | None |

All commands support `--help` for detailed information.

## Development

### Project Structure

```
azctx/
├── src/                    # Source code
│   ├── __init__.py
│   ├── cli.py             # CLI entry point (Typer commands)
│   ├── models/            # Data models (Context)
│   │   ├── __init__.py
│   │   └── context.py
│   ├── services/          # Business logic
│   │   ├── __init__.py
│   │   ├── azure_cli.py   # Azure CLI subprocess wrapper
│   │   ├── context_manager.py  # Interactive operations
│   │   └── storage.py     # YAML persistence
│   └── utils/             # Utilities
│       ├── __init__.py
│       ├── console.py     # Rich console helpers
│       └── errors.py      # Custom exceptions
├── specs/                 # Feature specifications
├── .specify/              # Speckit framework files
├── pyproject.toml         # Project dependencies
└── README.md              # This file
```

### Development Philosophy

This project follows the principles defined in [`.specify/memory/constitution.md`](.specify/memory/constitution.md):

1. **Simplicity First** - Favor straightforward solutions over complex ones
2. **Modular Architecture** - Self-contained components with clear responsibilities
3. **CLI-Centric Design** - Fast, readable, and easy-to-use command-line interface
4. **UV-Managed Workflow** - All dependencies managed through UV
5. **Pragmatic Quality** - Manual testing, no formal test suite (POC nature)

### Running Directly

For development, you can run the CLI directly without installation:

```bash
# Using UV
uv run python -m src.cli [command]

# Examples
uv run python -m src.cli add
uv run python -m src.cli switch
uv run python -m src.cli status --verbose
```

## Troubleshooting

### Azure CLI Not Found

If you see "Azure CLI is not installed or not in PATH":

1. Install Azure CLI: <https://aka.ms/azure-cli>
2. Verify installation: `az --version`
3. Restart your terminal

### No Active Azure Session

If you see "No active Azure session":

1. Login to Azure: `az login`
2. Verify login: `az account show`
3. Try your azctx command again

### Corrupted Contexts File

If you encounter errors loading contexts:

1. Backup your file: `copy %USERPROFILE%\.azctx\contexts.yaml contexts.yaml.bak`
2. Delete the file: `del %USERPROFILE%\.azctx\contexts.yaml`
3. Re-add your contexts with `azctx add`

## Contributing

This is a personal project, but suggestions and feedback are welcome via issues.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
