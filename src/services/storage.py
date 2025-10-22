"""YAML file storage operations for contexts."""

from pathlib import Path

import yaml

from src.models.context import Context
from src.utils.errors import StorageError


def get_storage_path() -> Path:
    """Get the path to the contexts YAML file.

    Returns:
        Path to ~/.azctx/contexts.yaml (cross-platform).
    """
    return Path.home() / ".azctx" / "contexts.yaml"


def _ensure_storage_dir() -> None:
    """Create the .azctx directory if it doesn't exist."""
    storage_dir = Path.home() / ".azctx"
    storage_dir.mkdir(parents=True, exist_ok=True)


def load_contexts() -> list[Context]:
    """Read and parse YAML file, return list of Context objects.

    Returns:
        List of Context objects. Empty list if file doesn't exist or is empty.

    Raises:
        StorageError: If the YAML file is corrupted or cannot be read.
    """
    storage_path = get_storage_path()

    # Return empty list if file doesn't exist
    if not storage_path.exists():
        return []

    try:
        with open(storage_path, encoding="utf-8") as f:
            data = yaml.safe_load(f)

        # Handle empty file or file with no contexts
        if not data or "contexts" not in data or not data["contexts"]:
            return []

        # Parse each context from YAML
        contexts = []
        for ctx_dict in data["contexts"]:
            try:
                contexts.append(Context.from_dict(ctx_dict))
            except (KeyError, ValueError) as e:
                # Skip invalid entries but continue loading others
                print(f"Warning: Skipping invalid context entry: {e}")
                continue

        return contexts

    except yaml.YAMLError as e:
        raise StorageError(f"Failed to parse contexts file: {e}")
    except Exception as e:
        raise StorageError(f"Failed to read contexts file: {e}")


def save_contexts(contexts: list[Context]) -> None:
    """Write list of Context objects to YAML file.

    Args:
        contexts: List of Context objects to save.

    Raises:
        StorageError: If the file cannot be written.
    """
    _ensure_storage_dir()
    storage_path = get_storage_path()

    try:
        # Convert contexts to dictionary format
        data = {"contexts": [ctx.to_dict() for ctx in contexts]}

        with open(storage_path, "w", encoding="utf-8") as f:
            yaml.safe_dump(data, f, default_flow_style=False, sort_keys=False)

    except Exception as e:
        raise StorageError(f"Failed to write contexts file: {e}")


def add_context(context: Context) -> None:
    """Append new context to YAML file.

    Args:
        context: Context object to add.

    Raises:
        StorageError: If the file cannot be read or written.
    """
    contexts = load_contexts()
    contexts.append(context)
    save_contexts(contexts)


def delete_context(context_id: str) -> None:
    """Remove context by ID and save YAML file.

    Args:
        context_id: The ID of the context to delete.

    Raises:
        StorageError: If the file cannot be read or written.
    """
    contexts = load_contexts()
    contexts = [ctx for ctx in contexts if ctx.context_id != context_id]
    save_contexts(contexts)


def get_context_by_id(context_id: str) -> Context | None:
    """Find and return context by ID.

    Args:
        context_id: The ID of the context to find.

    Returns:
        Context object if found, None otherwise.
    """
    contexts = load_contexts()
    for ctx in contexts:
        if ctx.context_id == context_id:
            return ctx
    return None


def context_id_exists(context_id: str) -> bool:
    """Check if a context with the given ID exists.

    Args:
        context_id: The ID to check.

    Returns:
        True if a context with this ID exists, False otherwise.
    """
    return get_context_by_id(context_id) is not None
