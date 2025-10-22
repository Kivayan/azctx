"""Custom exception classes for azctx."""


class AzureCliNotFoundError(Exception):
    """Raised when Azure CLI is not installed or not accessible in PATH."""

    pass


class NoActiveSessionError(Exception):
    """Raised when no Azure account is logged in."""

    pass


class DuplicateContextError(Exception):
    """Raised when attempting to add a context with an ID that already exists."""

    pass


class ContextNotFoundError(Exception):
    """Raised when a requested context ID cannot be found."""

    pass


class StorageError(Exception):
    """Raised when there's an error reading or writing the contexts YAML file."""

    pass
