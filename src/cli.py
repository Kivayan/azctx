"""CLI entry point for azctx.

Azure CLI context switcher - manage and switch between Azure accounts quickly.

Built with:
- Typer: Command structure and argument parsing
- Questionary: Interactive prompts
- Rich: Terminal output formatting
"""

import sys

import typer
from rich.panel import Panel
from rich.table import Table

from src import __version__
from src.services import azure_cli, context_manager
from src.utils.console import console
from src.utils.errors import AzureCliNotFoundError, DuplicateContextError, NoActiveSessionError

# Initialize Typer app
app = typer.Typer(
    name="azctx",
    help="Azure CLI context switcher - manage and switch between Azure accounts quickly.",
    add_completion=False,
)


def version_callback(value: bool) -> None:
    """Display version information."""
    if value:
        console.print(f"azctx version {__version__}")
        raise typer.Exit()


# Commands


@app.command()
def switch(
    id: str | None = typer.Option(
        None,
        "--id",
        "-i",
        help="Context ID to switch to directly (case-sensitive). Omit for interactive selection."
    )
) -> None:
    """Switch to a different saved context.

    Displays an interactive list of saved contexts (default behavior when no
    options are provided) or switches directly to a specific context when the
    --id option is used.

    This command supports two modes:

    1. Interactive Mode (default): Presents a navigable list of saved contexts
       using keyboard navigation (arrow keys + Enter). This is the original
       behavior and remains unchanged.

    2. Direct Mode (--id flag): Switches to a specific context by its ID without
       any interactive prompts. Useful for scripting and automation.

    Examples:
        azctx switch                 # Interactive selection
        azctx switch --id DEV        # Direct switch to "DEV" context
        azctx switch -i PROD         # Direct switch using short flag
    """
    try:
        # Determine which mode to use based on whether id parameter is provided
        if id:
            # Direct mode - switch by ID
            result = context_manager.switch_context_by_id(id)
        else:
            # Interactive mode - existing behavior
            result = context_manager.switch_context_interactive()

        if result["success"]:
            # Success - display green panel with context details
            context = result["context"]
            console.print(
                Panel(
                    f"[bold]Name:[/bold] {context.context_name}\n"
                    f"[bold]ID:[/bold] {context.context_id}\n"
                    f"[bold]Subscription:[/bold] {context.subscription_name}\n"
                    f"[bold]Tenant:[/bold] {context.tenant_name}\n"
                    f"[bold]Account:[/bold] {context.username}",
                    title=f"✓ Successfully Switched to {context.context_name}",
                    border_style="green",
                )
            )
            sys.exit(0)
        else:
            # Handle different error types
            error_type = result["error"]

            if error_type == "cancelled":
                console.print(
                    Panel(
                        result["message"],
                        title="Context Switch Cancelled",
                        border_style="yellow",
                    )
                )
                sys.exit(130)
            elif error_type == "already_active":
                # Already active context - informational message (yellow)
                console.print(
                    Panel(
                        result["message"],
                        title="⚠ Already Active",
                        border_style="yellow",
                    )
                )
                sys.exit(0)  # Not an error, just informational
            elif error_type == "not_found":
                # Context not found - show available IDs
                message = result["message"]
                if result.get("available_ids"):
                    available = ", ".join(result["available_ids"])
                    message += f"\n\nAvailable contexts: {available}"
                console.print(
                    Panel(
                        message,
                        title="✗ Context Not Found",
                        border_style="red",
                    )
                )
                sys.exit(1)
            elif error_type in ("empty_list", "single_context"):
                console.print(
                    Panel(
                        result["message"],
                        title="⚠ Cannot Switch Context",
                        border_style="yellow",
                    )
                )
                sys.exit(1)
            else:
                # All other errors (verification_failed, no_session, unknown, etc.)
                console.print(
                    Panel(
                        result["message"],
                        title="✗ Failed to Switch Context",
                        border_style="red",
                    )
                )
                sys.exit(1)

    except KeyboardInterrupt:
        console.print(
            Panel(
                "Context switch cancelled",
                title="Cancelled",
                border_style="yellow",
            )
        )
        sys.exit(130)
    except AzureCliNotFoundError as e:
        console.print(
            Panel(
                f"{e}\n\nPlease install Azure CLI: https://aka.ms/azure-cli",
                title="✗ Azure CLI Not Found",
                border_style="red",
            )
        )
        sys.exit(1)
    except Exception as e:
        console.print(
            Panel(
                f"Unexpected error: {e}",
                title="✗ Error",
                border_style="red",
            )
        )
        sys.exit(1)


@app.command()
def status(
    verbose: bool = typer.Option(False, "--verbose", "-v", help="Show detailed information")
) -> None:
    """Show currently active Azure context.

    Displays the current Azure CLI session and indicates whether it's managed by azctx.
    """
    try:
        result = context_manager.get_status(verbose)

        if not result["success"]:
            # Handle no active session
            console.print(
                Panel(
                    f"{result['message']}\n\n"
                    "Please run 'az login' to authenticate with Azure first.",
                    title="✗ No Active Azure Session",
                    border_style="red",
                )
            )
            sys.exit(1)

        # Display status based on whether context is managed
        if result["is_managed"]:
            # Managed context - show with green border and ✓ symbol
            context = result["context"]
            current = result["current_account"]

            # Build content
            content_lines = [
                f"[bold]Name:[/bold] {context.context_name}",
                f"[bold]ID:[/bold] {context.context_id}",
                f"[bold]Subscription:[/bold] {context.subscription_name}",
                f"[bold]Tenant:[/bold] {context.tenant_name}",
                f"[bold]Account:[/bold] {context.username}",
            ]

            # Add verbose details if requested
            if verbose:
                content_lines.extend(
                    [
                        "",
                        "[dim]Detailed Information:[/dim]",
                        f"[dim]Subscription ID:[/dim] {context.subscription_id}",
                        f"[dim]Tenant ID:[/dim] {context.tenant_id}",
                        f"[dim]Created:[/dim] {context.created_at}",
                    ]
                )

            console.print(
                Panel(
                    "\n".join(content_lines),
                    title=f"✓ Managed Context: {context.context_name}",
                    border_style="green",
                )
            )
            sys.exit(0)
        else:
            # Unmanaged context - show with yellow border and ⚠ symbol
            current = result["current_account"]

            # Build content
            content_lines = [
                f"[bold]Subscription:[/bold] {current['name']}",
                f"[bold]Tenant:[/bold] {current['tenantId']}",
                f"[bold]Account:[/bold] {current['user']['name']}",
                "",
                "[yellow]This context is not managed by azctx.[/yellow]",
                "Run [cyan]azctx add[/cyan] to save it with a friendly name.",
            ]

            # Add verbose details if requested
            if verbose:
                content_lines.extend(
                    [
                        "",
                        "[dim]Detailed Information:[/dim]",
                        f"[dim]Subscription ID:[/dim] {current['id']}",
                        f"[dim]Tenant ID:[/dim] {current['tenantId']}",
                    ]
                )

            console.print(
                Panel(
                    "\n".join(content_lines),
                    title="⚠ Unmanaged Context",
                    border_style="yellow",
                )
            )
            sys.exit(0)

    except AzureCliNotFoundError as e:
        console.print(
            Panel(
                f"{e}\n\nPlease install Azure CLI: https://aka.ms/azure-cli",
                title="✗ Azure CLI Not Found",
                border_style="red",
            )
        )
        sys.exit(1)
    except NoActiveSessionError as e:
        console.print(
            Panel(
                f"{e}\n\nPlease run 'az login' to authenticate with Azure first.",
                title="✗ No Active Azure Session",
                border_style="red",
            )
        )
        sys.exit(1)
    except Exception as e:
        console.print(
            Panel(
                f"Unexpected error: {e}",
                title="✗ Error",
                border_style="red",
            )
        )
        sys.exit(1)


@app.command()
def add() -> None:
    """Add current Azure CLI context with friendly identifiers.

    Saves the currently active Azure CLI session with a user-assigned name and ID.
    """
    try:
        # First, get and display the current Azure context
        try:
            current_account = azure_cli.get_current_account()
            if current_account:
                console.print(
                    Panel(
                        f"[bold]Subscription:[/bold] {current_account['name']}\n"
                        f"[bold]Tenant:[/bold] {current_account['tenantId']}\n"
                        f"[bold]Account:[/bold] {current_account['user']['name']}",
                        title="Current Azure Context",
                        border_style="cyan",
                    )
                )
                console.print()  # Add blank line for spacing
        except (NoActiveSessionError, Exception):
            # If we can't get current account, the interactive function will handle it
            pass

        result = context_manager.add_context_interactive()

        if result["success"]:
            # Success - display green panel with context details
            context = result["context"]
            console.print(
                Panel(
                    f"[bold]Name:[/bold] {context.context_name}\n"
                    f"[bold]ID:[/bold] {context.context_id}\n"
                    f"[bold]Subscription:[/bold] {context.subscription_name}\n"
                    f"[bold]Tenant:[/bold] {context.tenant_name}\n"
                    f"[bold]Account:[/bold] {context.username}",
                    title=f"✓ Successfully Added Context: {context.context_name}",
                    border_style="green",
                )
            )
            sys.exit(0)
        else:
            # Handle different error types
            error_type = result["error"]

            if error_type == "cancelled":
                console.print(
                    Panel(
                        result["message"],
                        title="Context Addition Cancelled",
                        border_style="yellow",
                    )
                )
                sys.exit(130)
            elif error_type == "already_exists":
                console.print(
                    Panel(
                        result["message"],
                        title="⚠ Context Already Managed",
                        border_style="yellow",
                    )
                )
                sys.exit(1)
            elif error_type == "no_session":
                console.print(
                    Panel(
                        f"{result['message']}\n\n"
                        "Please run 'az login' to authenticate with Azure first.",
                        title="✗ No Active Azure Session",
                        border_style="red",
                    )
                )
                sys.exit(1)
            else:
                console.print(
                    Panel(
                        result["message"],
                        title="✗ Failed to Add Context",
                        border_style="red",
                    )
                )
                sys.exit(1)

    except KeyboardInterrupt:
        console.print(
            Panel(
                "Context addition cancelled",
                title="Cancelled",
                border_style="yellow",
            )
        )
        sys.exit(130)
    except AzureCliNotFoundError as e:
        console.print(
            Panel(
                f"{e}\n\nPlease install Azure CLI: https://aka.ms/azure-cli",
                title="✗ Azure CLI Not Found",
                border_style="red",
            )
        )
        sys.exit(1)
    except DuplicateContextError as e:
        console.print(
            Panel(
                f"{e}\n\nPlease choose a different context ID.",
                title="✗ Duplicate Context ID",
                border_style="red",
            )
        )
        sys.exit(1)
    except NoActiveSessionError as e:
        console.print(
            Panel(
                f"{e}\n\nPlease run 'az login' to authenticate with Azure first.",
                title="✗ No Active Azure Session",
                border_style="red",
            )
        )
        sys.exit(1)
    except Exception as e:
        console.print(
            Panel(
                f"Unexpected error: {e}",
                title="✗ Error",
                border_style="red",
            )
        )
        sys.exit(1)


@app.command()
def list(
    verbose: bool = typer.Option(False, "--verbose", "-v", help="Show detailed information")
) -> None:
    """List all saved contexts.

    Displays a table of all saved Azure CLI contexts with their IDs and names.
    Use --verbose to see full details for each context.
    """
    try:
        result = context_manager.list_contexts(verbose)

        # Handle empty list
        if result["error"] == "empty_list":
            console.print(
                Panel(
                    result["message"],
                    title="No Saved Contexts",
                    border_style="yellow",
                )
            )
            sys.exit(0)

        contexts = result["contexts"]

        if verbose:
            # Verbose mode: Display each context in a separate panel
            console.print(
                f"\n[cyan bold]Saved Contexts ({len(contexts)}):[/cyan bold]\n"
            )
            for context in contexts:
                console.print(
                    Panel(
                        f"[bold]Name:[/bold] {context.context_name}\n"
                        f"[bold]Subscription:[/bold] {context.subscription_name} "
                        f"({context.subscription_id})\n"
                        f"[bold]Tenant ID:[/bold] {context.tenant_id}\n"
                        f"[bold]Username:[/bold] {context.username}\n"
                        f"[bold]Created:[/bold] {context.created_at}",
                        title=f"[cyan]{context.context_id}[/cyan]",
                        border_style="cyan",
                    )
                )
            sys.exit(0)
        else:
            # Simple mode: Display table with ID and Name
            table = Table(
                title=f"[cyan bold]Saved Contexts ({len(contexts)})[/cyan bold]",
                show_header=True,
                header_style="cyan bold",
            )
            table.add_column("ID", style="cyan")
            table.add_column("Name", style="white")

            # Add rows from context tuples
            for context_id, context_name in contexts:
                table.add_row(context_id, context_name)

            console.print()
            console.print(table)
            console.print(
                "\n[dim]Use 'azctx list --verbose' for detailed information[/dim]\n"
            )
            sys.exit(0)

    except Exception as e:
        console.print(
            Panel(
                f"Unexpected error: {e}",
                title="✗ Error",
                border_style="red",
            )
        )
        sys.exit(1)


@app.command()
def delete() -> None:
    """Delete a saved context.

    Displays an interactive list of saved contexts, allows selection,
    and confirms before deleting. Note: This does not affect the active Azure CLI session.
    """
    try:
        result = context_manager.delete_context_interactive()

        if result["success"]:
            # Success - display green panel
            console.print(
                Panel(
                    f"[bold]Context ID:[/bold] {result['context_id']}\n"
                    f"[bold]Context Name:[/bold] {result['context_name']}\n\n"
                    "[dim]Note: Your active Azure CLI session remains unchanged.[/dim]",
                    title=f"✓ Successfully Deleted Context: {result['context_name']}",
                    border_style="green",
                )
            )
            sys.exit(0)
        else:
            # Handle different error types
            error_type = result["error"]

            if error_type == "cancelled":
                console.print(
                    Panel(
                        result["message"],
                        title="Deletion Cancelled",
                        border_style="yellow",
                    )
                )
                sys.exit(130)
            elif error_type == "empty_list":
                console.print(
                    Panel(
                        result["message"],
                        title="⚠ No Contexts to Delete",
                        border_style="yellow",
                    )
                )
                sys.exit(0)
            elif error_type == "not_found":
                console.print(
                    Panel(
                        result["message"],
                        title="✗ Context Not Found",
                        border_style="red",
                    )
                )
                sys.exit(1)
            else:
                console.print(
                    Panel(
                        result["message"],
                        title="✗ Failed to Delete Context",
                        border_style="red",
                    )
                )
                sys.exit(1)

    except KeyboardInterrupt:
        console.print(
            Panel(
                "Deletion cancelled",
                title="Cancelled",
                border_style="yellow",
            )
        )
        sys.exit(130)
    except Exception as e:
        console.print(
            Panel(
                f"Unexpected error: {e}",
                title="✗ Error",
                border_style="red",
            )
        )
        sys.exit(1)


@app.callback(invoke_without_command=True)
def callback(
    ctx: typer.Context,
    version: bool = typer.Option(
        False,
        "--version",
        "-V",
        help="Show version information.",
        callback=version_callback,
        is_eager=True,
    )
) -> None:
    """
    Azure CLI context switcher.

    Run with --help to see available commands.
    """
    if ctx.invoked_subcommand is None:
        # When no command is provided, show help
        typer.echo(ctx.get_help())


def main() -> None:
    """Entry point for the CLI application."""
    app()


if __name__ == "__main__":
    app()
