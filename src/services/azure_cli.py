"""Azure CLI integration via subprocess."""

import json
import os
import shutil
import subprocess
import sys
from typing import Any

from src.utils.errors import AzureCliNotFoundError, NoActiveSessionError


def _find_az_executable() -> str | None:
    """Find the Azure CLI executable path.

    Returns:
        Full path to the az executable, or None if not found.
    """
    # First try shutil.which which searches PATH
    az_path = shutil.which("az")
    if az_path:
        return az_path

    # On Windows, also check common installation locations
    if sys.platform == "win32":
        common_paths = [
            r"C:\Program Files\Microsoft SDKs\Azure\CLI2\wbin\az.cmd",
            r"C:\Program Files (x86)\Microsoft SDKs\Azure\CLI2\wbin\az.cmd",
        ]
        for path in common_paths:
            if os.path.exists(path):
                return path

    return None


def _run_az_command(args: list[str], timeout: int = 5) -> subprocess.CompletedProcess:
    """Run an Azure CLI command with proper environment handling.

    Args:
        args: Command arguments (e.g., ["account", "show"]).
        timeout: Command timeout in seconds.

    Returns:
        CompletedProcess object from subprocess.run.

    Raises:
        FileNotFoundError: If Azure CLI is not found.
        subprocess.TimeoutExpired: If command times out.
    """
    # Find the Azure CLI executable
    az_executable = _find_az_executable()
    if not az_executable:
        raise FileNotFoundError("Azure CLI executable not found")

    # Build command with full path to az
    cmd = [az_executable] + args

    # On Windows, use shell=True when using .cmd files
    use_shell = sys.platform == "win32" and az_executable.endswith(".cmd")

    # Get current environment to ensure PATH is inherited
    env = os.environ.copy()

    return subprocess.run(
        cmd,
        capture_output=True,
        text=True,
        timeout=timeout,
        shell=use_shell,
        env=env,
    )


def check_azure_cli_installed() -> bool:
    """Verify Azure CLI is installed and accessible.

    Returns:
        True if Azure CLI is installed, False otherwise.
    """
    try:
        result = _run_az_command(["--version"], timeout=30)
        return result.returncode == 0
    except (FileNotFoundError, subprocess.TimeoutExpired):
        return False


def get_current_account() -> dict[str, Any] | None:
    """Execute 'az account show' and return parsed JSON.

    Returns:
        Dictionary containing Azure account information, or None if no active session.

    Raises:
        AzureCliNotFoundError: If Azure CLI is not installed.
        NoActiveSessionError: If no Azure account is logged in.
    """
    if not check_azure_cli_installed():
        raise AzureCliNotFoundError("Azure CLI is not installed or not in PATH")

    try:
        result = _run_az_command(["account", "show", "--output", "json"], timeout=30)

        # Check if command failed (no active session)
        if result.returncode != 0:
            # Check error message for "Please run 'az login'"
            if "az login" in result.stderr.lower():
                raise NoActiveSessionError("No active Azure session. Run 'az login' first.")
            raise NoActiveSessionError(f"Failed to get Azure account: {result.stderr}")

        # Parse JSON output
        try:
            account_data = json.loads(result.stdout)
            return account_data
        except json.JSONDecodeError as e:
            raise NoActiveSessionError(f"Failed to parse Azure CLI output: {e}")

    except subprocess.TimeoutExpired:
        raise NoActiveSessionError("Azure CLI command timed out")


def set_account(subscription_id: str) -> bool:
    """Execute 'az account set --subscription <id>'.

    Args:
        subscription_id: The Azure subscription ID to set as active.

    Returns:
        True if successful, False otherwise.

    Raises:
        AzureCliNotFoundError: If Azure CLI is not installed.
    """
    if not check_azure_cli_installed():
        raise AzureCliNotFoundError("Azure CLI is not installed or not in PATH")

    try:
        result = _run_az_command(["account", "set", "--subscription", subscription_id])
        return result.returncode == 0

    except subprocess.TimeoutExpired:
        return False
