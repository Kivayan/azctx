"""Context management service for interactive operations."""

from datetime import datetime
from typing import Any

import questionary

from src.models.context import Context
from src.services import azure_cli, storage
from src.utils.errors import (
    AzureCliNotFoundError,
    DuplicateContextError,
    NoActiveSessionError,
)


def switch_context_interactive() -> dict[str, Any]:
    """Interactively switch between saved Azure CLI contexts.

    Loads saved contexts, presents them in an interactive list,
    and switches to the selected context.

    Returns:
        Dictionary with operation result:
        - success (bool): Whether the operation succeeded
        - message (str): Status message
        - context (Context | None): Selected context if successful
        - error (str | None): Error type if failed

    Raises:
        AzureCliNotFoundError: If Azure CLI is not installed.
    """
    # Verify Azure CLI is available
    if not azure_cli.check_azure_cli_installed():
        raise AzureCliNotFoundError("Azure CLI is not installed or not in PATH")

    # Load saved contexts
    contexts = storage.load_contexts()

    # Handle empty contexts list
    if not contexts:
        return {
            "success": False,
            "message": (
                "No saved contexts found. "
                "Use 'azctx add' to save your current context first."
            ),
            "context": None,
            "error": "empty_list",
        }

    # Handle single context - inform user
    if len(contexts) == 1:
        ctx = contexts[0]
        return {
            "success": False,
            "message": (
                f"Only one context saved: {ctx.context_name} ({ctx.context_id}). "
                "Use 'azctx add' to save more contexts before switching."
            ),
            "context": contexts[0],
            "error": "single_context",
        }

    # Create choices for questionary - format: "Name (ID)"
    choices = [f"{ctx.context_name} ({ctx.context_id})" for ctx in contexts]

    # Interactive selection
    try:
        selected = questionary.select(
            "Select a context to switch to:",
            choices=choices,
            use_arrow_keys=True,
        ).ask()

        # Handle cancellation (Ctrl+C or Esc)
        if selected is None:
            return {
                "success": False,
                "message": "Context switch cancelled",
                "context": None,
                "error": "cancelled",
            }

        # Extract context_id from selected choice (format: "Name (ID)")
        # Find the ID between the last pair of parentheses
        context_id = selected.rsplit("(", 1)[1].rstrip(")")
        selected_context = storage.get_context_by_id(context_id)

        if not selected_context:
            return {
                "success": False,
                "message": f"Context '{context_id}' not found",
                "context": None,
                "error": "not_found",
            }

        # Switch to selected context
        azure_cli.set_account(selected_context.subscription_id)

        # Verify the switch was successful
        current_account = azure_cli.get_current_account()
        if current_account and current_account["id"] != selected_context.subscription_id:
            return {
                "success": False,
                "message": (
                    f"Failed to switch to context '{context_id}'. "
                    "Azure CLI returned different subscription."
                ),
                "context": selected_context,
                "error": "verification_failed",
            }

        # Success
        return {
            "success": True,
            "message": f"Successfully switched to context: {selected_context.context_id}",
            "context": selected_context,
            "error": None,
        }

    except KeyboardInterrupt:
        return {
            "success": False,
            "message": "Context switch cancelled",
            "context": None,
            "error": "cancelled",
        }
    except NoActiveSessionError as e:
        return {
            "success": False,
            "message": str(e),
            "context": None,
            "error": "no_session",
        }
    except Exception as e:
        return {
            "success": False,
            "message": f"Failed to switch context: {e}",
            "context": None,
            "error": "unknown",
        }


def add_context_interactive() -> dict[str, Any]:
    """Interactively add the current Azure CLI context with user-assigned identifiers.

    Prompts user for context name and ID, then saves the current Azure session.

    Returns:
        Dictionary with operation result:
        - success (bool): Whether the operation succeeded
        - message (str): Status message
        - context (Context | None): Added context if successful
        - current_account (dict | None): Current Azure account info
        - error (str | None): Error type if failed

    Raises:
        AzureCliNotFoundError: If Azure CLI is not installed.
    """
    # Check Azure CLI is installed
    if not azure_cli.check_azure_cli_installed():
        raise AzureCliNotFoundError("Azure CLI is not installed or not in PATH")

    # Get current Azure account
    try:
        current_account = azure_cli.get_current_account()
    except NoActiveSessionError as e:
        return {
            "success": False,
            "message": str(e),
            "context": None,
            "current_account": None,
            "error": "no_session",
        }

    if not current_account:
        return {
            "success": False,
            "message": "No active Azure session. Run 'az login' first.",
            "context": None,
            "current_account": None,
            "error": "no_session",
        }

    # Check if this context already exists (by subscription, tenant, and user)
    existing_contexts = storage.load_contexts()
    for ctx in existing_contexts:
        if (
            ctx.subscription_id == current_account["id"]
            and ctx.tenant_id == current_account["tenantId"]
            and ctx.username == current_account["user"]["name"]
        ):
            return {
                "success": False,
                "message": (
                    f"This Azure context is already managed as "
                    f"'{ctx.context_name}' ({ctx.context_id}).\n"
                    f"Subscription: {ctx.subscription_name}\n"
                    f"Tenant: {ctx.tenant_id}\n"
                    f"Account: {ctx.username}"
                ),
                "context": ctx,
                "current_account": current_account,
                "error": "already_exists",
            }

    # Prompt for context name
    context_name = None
    while not context_name:
        try:
            name_input = questionary.text(
                "Enter a friendly name for this context (1-100 characters):",
                validate=lambda text: (
                    len(text) >= 1 and len(text) <= 100
                    or "Name must be between 1 and 100 characters"
                ),
            ).ask()

            if name_input is None:  # User cancelled
                return {
                    "success": False,
                    "message": "Context addition cancelled",
                    "context": None,
                    "current_account": current_account,
                    "error": "cancelled",
                }

            context_name = name_input.strip()

        except KeyboardInterrupt:
            return {
                "success": False,
                "message": "Context addition cancelled",
                "context": None,
                "current_account": current_account,
                "error": "cancelled",
            }

    # Prompt for context ID with validation
    context_id = None
    while not context_id:
        try:
            id_input = questionary.text(
                "Enter a short ID for this context (1-20 chars, alphanumeric/hyphens/underscores):",
                validate=lambda text: (
                    Context.validate_context_id(text)
                    or "ID must be 1-20 characters (alphanumeric, hyphens, underscores only)"
                ),
            ).ask()

            if id_input is None:  # User cancelled
                return {
                    "success": False,
                    "message": "Context addition cancelled",
                    "context": None,
                    "current_account": current_account,
                    "error": "cancelled",
                }

            # Check uniqueness
            if storage.context_id_exists(id_input):
                raise DuplicateContextError(
                    f"Context ID '{id_input}' already exists. Please choose a different ID."
                )

            context_id = id_input.strip()

        except KeyboardInterrupt:
            return {
                "success": False,
                "message": "Context addition cancelled",
                "context": None,
                "current_account": current_account,
                "error": "cancelled",
            }
        except DuplicateContextError as e:
            # Show error and re-prompt
            print(f"Error: {e}")
            continue

    # Create Context object
    new_context = Context(
        context_id=context_id,
        context_name=context_name,
        subscription_id=current_account["id"],
        subscription_name=current_account["name"],
        tenant_id=current_account["tenantId"],
        tenant_name=current_account["tenantId"],  # Simple approach per data-model.md
        username=current_account["user"]["name"],
        created_at=datetime.now(),
    )

    # Save context
    try:
        storage.add_context(new_context)
        return {
            "success": True,
            "message": f"Successfully added context: {context_name}",
            "context": new_context,
            "current_account": current_account,
            "error": None,
        }
    except Exception as e:
        return {
            "success": False,
            "message": f"Failed to save context: {e}",
            "context": None,
            "current_account": current_account,
            "error": "storage_error",
        }


def get_status(verbose: bool = False) -> dict[str, Any]:
    """Get the status of the currently active Azure context.

    Checks if the current Azure session matches a saved context.

    Args:
        verbose: Whether to include detailed information.

    Returns:
        Dictionary with operation result:
        - success (bool): Whether the operation succeeded
        - message (str): Status message
        - context (Context | None): Matched context if managed
        - current_account (dict | None): Current Azure account info
        - is_managed (bool): Whether context is managed by azctx
        - verbose (bool): Echo of verbose parameter
        - error (str | None): Error type if failed

    Raises:
        AzureCliNotFoundError: If Azure CLI is not installed.
    """
    # Check Azure CLI is installed
    if not azure_cli.check_azure_cli_installed():
        raise AzureCliNotFoundError("Azure CLI is not installed or not in PATH")

    # Get current Azure account
    try:
        current_account = azure_cli.get_current_account()
    except NoActiveSessionError as e:
        return {
            "success": False,
            "message": str(e),
            "context": None,
            "current_account": None,
            "is_managed": False,
            "verbose": verbose,
            "error": "no_session",
        }

    if not current_account:
        return {
            "success": False,
            "message": "No active Azure session. Run 'az login' first.",
            "context": None,
            "current_account": None,
            "is_managed": False,
            "verbose": verbose,
            "error": "no_session",
        }

    # Load saved contexts
    contexts = storage.load_contexts()

    # Match current subscription against saved contexts
    matched_context = None
    for ctx in contexts:
        if ctx.subscription_id == current_account["id"]:
            matched_context = ctx
            break

    # Return result based on whether context is managed
    if matched_context:
        return {
            "success": True,
            "message": "Current context is managed by azctx",
            "context": matched_context,
            "current_account": current_account,
            "is_managed": True,
            "verbose": verbose,
            "error": None,
        }
    else:
        return {
            "success": True,
            "message": "Current context is not managed. Use 'azctx add' to save it.",
            "context": None,
            "current_account": current_account,
            "is_managed": False,
            "verbose": verbose,
            "error": None,
        }


def list_contexts(verbose: bool = False) -> dict[str, Any]:
    """List all saved Azure CLI contexts.

    Args:
        verbose: Whether to return full Context objects or just ID/Name tuples.

    Returns:
        Dictionary with operation result:
        - success (bool): Whether the operation succeeded
        - message (str): Status message
        - contexts (list[Context] | list[tuple[str, str]] | None):
            Full Context objects if verbose=True,
            List of (context_id, context_name) tuples if verbose=False,
            None if empty
        - verbose (bool): Echo of verbose parameter
        - error (str | None): Error type if failed
    """
    # Load saved contexts
    contexts = storage.load_contexts()

    # Handle empty list
    if not contexts:
        return {
            "success": True,
            "message": "No saved contexts found. Use 'azctx add' to save your current context.",
            "contexts": None,
            "verbose": verbose,
            "error": "empty_list",
        }

    # Return based on verbose flag
    if verbose:
        # Return full Context objects for detailed view
        return {
            "success": True,
            "message": f"Found {len(contexts)} saved context(s)",
            "contexts": contexts,
            "verbose": verbose,
            "error": None,
        }
    else:
        # Return tuples of (context_id, context_name) for simple table
        context_tuples = [(ctx.context_id, ctx.context_name) for ctx in contexts]
        return {
            "success": True,
            "message": f"Found {len(contexts)} saved context(s)",
            "contexts": context_tuples,
            "verbose": verbose,
            "error": None,
        }


def delete_context_interactive() -> dict[str, Any]:
    """Interactively select and delete a saved context with confirmation.

    Returns:
        Dictionary with operation result:
        - success (bool): Whether the operation succeeded
        - message (str): Status message
        - context_id (str | None): Deleted context ID if successful
        - context_name (str | None): Deleted context name if successful
        - error (str | None): Error type if failed
    """
    # Load saved contexts
    contexts = storage.load_contexts()

    # Handle empty list
    if not contexts:
        return {
            "success": False,
            "message": "No saved contexts found. Nothing to delete.",
            "context_id": None,
            "context_name": None,
            "error": "empty_list",
        }

    # Create choices for questionary - format: "Name (ID)"
    choices = [f"{ctx.context_name} ({ctx.context_id})" for ctx in contexts]

    # Interactive selection
    try:
        selected = questionary.select(
            "Select a context to delete:",
            choices=choices,
            use_arrow_keys=True,
        ).ask()

        # Handle cancellation (Ctrl+C or Esc)
        if selected is None:
            return {
                "success": False,
                "message": "Deletion cancelled",
                "context_id": None,
                "context_name": None,
                "error": "cancelled",
            }

        # Extract context_id from selected choice (format: "Name (ID)")
        context_id = selected.rsplit("(", 1)[1].rstrip(")")
        selected_context = storage.get_context_by_id(context_id)

        if not selected_context:
            return {
                "success": False,
                "message": f"Context '{context_id}' not found",
                "context_id": None,
                "context_name": None,
                "error": "not_found",
            }

        # Confirmation prompt
        confirmed = questionary.confirm(
            f"Are you sure you want to delete context '{selected_context.context_name}' "
            f"({selected_context.context_id})?",
            default=False,
        ).ask()

        # Handle cancellation or no confirmation
        if confirmed is None or not confirmed:
            return {
                "success": False,
                "message": "Deletion cancelled",
                "context_id": None,
                "context_name": None,
                "error": "cancelled",
            }

        # Delete the context
        storage.delete_context(context_id)

        # Success
        return {
            "success": True,
            "message": (
                f"Successfully deleted context: {selected_context.context_name} "
                f"({selected_context.context_id})"
            ),
            "context_id": selected_context.context_id,
            "context_name": selected_context.context_name,
            "error": None,
        }

    except KeyboardInterrupt:
        return {
            "success": False,
            "message": "Deletion cancelled",
            "context_id": None,
            "context_name": None,
            "error": "cancelled",
        }
    except Exception as e:
        return {
            "success": False,
            "message": f"Failed to delete context: {e}",
            "context_id": None,
            "context_name": None,
            "error": "unknown",
        }


def switch_context_by_id(context_id: str) -> dict[str, Any]:
    """Switch to a saved Azure CLI context by its ID (case-sensitive).

    This function provides direct, non-interactive context switching by ID.
    It performs case-sensitive matching against saved context IDs and
    returns detailed result information for CLI display.

    Args:
        context_id: The case-sensitive ID of the context to switch to.
                    Leading and trailing whitespace will be trimmed before matching.
                    Empty strings (after trimming) will result in an error.

    Returns:
        Dictionary with operation result containing these keys:

        - success (bool): Whether the operation succeeded
        - message (str): Human-readable status message for display
        - context (Context | None): Selected Context object if successful, None on error
        - error (str | None): Error type identifier (see Error Types below)
        - available_ids (list[str] | None): Alphabetically sorted list of all context IDs
                                           (only included when error == "not_found")

    Raises:
        AzureCliNotFoundError: If Azure CLI is not installed or not accessible in PATH.
                              This is a fatal error that should terminate the program.

    Error Types (value of 'error' field):
        - None: Success (success == True)
        - "empty_list": No contexts are saved in storage
        - "not_found": Provided context_id doesn't match any saved context
        - "already_active": Target context is already the active Azure account
        - "verification_failed": Azure CLI switch command executed but verification failed
        - "no_session": No active Azure CLI session (user needs to run 'az login')
        - "unknown": Unexpected error occurred
    """
    # Trim whitespace from context_id parameter
    context_id = context_id.strip()

    # Validate context_id is not empty after trimming
    if not context_id:
        return {
            "success": False,
            "message": "Context ID cannot be empty",
            "context": None,
            "error": "empty_id",
            "available_ids": None,
        }

    # Verify Azure CLI is available
    if not azure_cli.check_azure_cli_installed():
        raise AzureCliNotFoundError("Azure CLI is not installed or not in PATH")

    # Load saved contexts
    try:
        contexts = storage.load_contexts()
    except Exception as e:
        return {
            "success": False,
            "message": f"Failed to load contexts: {e}",
            "context": None,
            "error": "storage_error",
            "available_ids": None,
        }

    # Check for empty list
    if not contexts:
        return {
            "success": False,
            "message": (
                "No saved contexts found. "
                "Use 'azctx add' to save your current context first."
            ),
            "context": None,
            "error": "empty_list",
            "available_ids": None,
        }

    # Use storage.get_context_by_id() for case-sensitive lookup
    selected_context = storage.get_context_by_id(context_id)

    # If not found, generate sorted available_ids list
    if not selected_context:
        available_ids = sorted([ctx.context_id for ctx in contexts])
        return {
            "success": False,
            "message": f"Context '{context_id}' not found.",
            "context": None,
            "error": "not_found",
            "available_ids": available_ids,
        }

    # Check if target context is already active
    try:
        current_account = azure_cli.get_current_account()
        if current_account and current_account.get("id") == selected_context.subscription_id:
            return {
                "success": False,
                "message": f"Context '{context_id}' is already active.",
                "context": selected_context,
                "error": "already_active",
                "available_ids": None,
            }
    except NoActiveSessionError:
        return {
            "success": False,
            "message": "No active Azure session. Run 'az login' first.",
            "context": None,
            "error": "no_session",
            "available_ids": None,
        }
    except Exception:
        # Continue with switch attempt even if current account check fails
        pass

    # Switch via azure_cli.set_account(context.subscription_id)
    try:
        success = azure_cli.set_account(selected_context.subscription_id)
        if not success:
            return {
                "success": False,
                "message": f"Failed to switch to context '{context_id}'.",
                "context": selected_context,
                "error": "switch_failed",
                "available_ids": None,
            }
    except NoActiveSessionError as e:
        return {
            "success": False,
            "message": str(e),
            "context": selected_context,
            "error": "no_session",
            "available_ids": None,
        }
    except Exception as e:
        return {
            "success": False,
            "message": f"Failed to switch context: {e}",
            "context": selected_context,
            "error": "unknown",
            "available_ids": None,
        }

    # Verify switch with azure_cli.get_current_account()
    try:
        current_account = azure_cli.get_current_account()
        if current_account and current_account.get("id") != selected_context.subscription_id:
            return {
                "success": False,
                "message": (
                    f"Failed to verify switch to context '{context_id}'. "
                    "Azure CLI returned different subscription."
                ),
                "context": selected_context,
                "error": "verification_failed",
                "available_ids": None,
            }
    except Exception:
        # If verification fails but switch command succeeded, consider it a warning
        # Still return success but note the verification issue
        pass

    # Return success
    return {
        "success": True,
        "message": f"Successfully switched to context: {context_id}",
        "context": selected_context,
        "error": None,
        "available_ids": None,
    }
