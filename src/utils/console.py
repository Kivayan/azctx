"""Rich console instance and helper functions for formatted output."""

from rich.console import Console
from rich.panel import Panel

# Global console instance for consistent output formatting
console = Console()


def print_success(message: str, title: str = "Success") -> None:
    """Print a success message in a green panel with checkmark.

    Args:
        message: The success message to display.
        title: Optional title for the panel (default: "Success").
    """
    console.print(
        Panel(
            f"✓ {message}",
            title=title,
            border_style="green",
            padding=(1, 2),
        )
    )


def print_error(message: str, title: str = "Error") -> None:
    """Print an error message in a red panel with X symbol.

    Args:
        message: The error message to display.
        title: Optional title for the panel (default: "Error").
    """
    console.print(
        Panel(
            f"✗ {message}",
            title=title,
            border_style="red",
            padding=(1, 2),
        )
    )


def print_warning(message: str, title: str = "Warning") -> None:
    """Print a warning message in a yellow panel with warning symbol.

    Args:
        message: The warning message to display.
        title: Optional title for the panel (default: "Warning").
    """
    console.print(
        Panel(
            f"⚠ {message}",
            title=title,
            border_style="yellow",
            padding=(1, 2),
        )
    )


def print_info(message: str, title: str = "Info") -> None:
    """Print an info message in a cyan panel.

    Args:
        message: The info message to display.
        title: Optional title for the panel (default: "Info").
    """
    console.print(
        Panel(
            message,
            title=title,
            border_style="cyan",
            padding=(1, 2),
        )
    )
