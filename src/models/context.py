"""Context model for Azure CLI account configurations."""

import re
from dataclasses import dataclass
from datetime import datetime
from typing import ClassVar


@dataclass
class Context:
    """Represents a saved Azure CLI account configuration.

    Attributes:
        context_id: Short, user-assigned identifier (1-20 chars, alphanumeric, hyphens, underscores)
        context_name: Friendly, descriptive name (1-100 chars)
        subscription_id: Azure subscription GUID
        subscription_name: Azure subscription display name
        tenant_id: Azure tenant GUID
        tenant_name: Azure tenant domain name
        username: Authenticated user email/principal
        created_at: When context was added to tool (ISO 8601 format)
    """

    context_id: str
    context_name: str
    subscription_id: str
    subscription_name: str
    tenant_id: str
    tenant_name: str
    username: str
    created_at: datetime

    # Class variable for validation regex
    CONTEXT_ID_PATTERN: ClassVar[str] = r"^[a-zA-Z0-9_-]{1,20}$"

    def to_dict(self) -> dict:
        """Convert Context to dictionary for YAML serialization.

        Returns:
            Dictionary representation of the Context with datetime as ISO string.
        """
        return {
            "context_id": self.context_id,
            "context_name": self.context_name,
            "subscription_id": self.subscription_id,
            "subscription_name": self.subscription_name,
            "tenant_id": self.tenant_id,
            "tenant_name": self.tenant_name,
            "username": self.username,
            "created_at": self.created_at.isoformat(),
        }

    @classmethod
    def from_dict(cls, data: dict) -> "Context":
        """Create Context from dictionary (YAML deserialization).

        Args:
            data: Dictionary containing context fields.

        Returns:
            Context instance.
        """
        # Parse ISO format datetime string
        created_at = data["created_at"]
        if isinstance(created_at, str):
            created_at = datetime.fromisoformat(created_at)

        return cls(
            context_id=data["context_id"],
            context_name=data["context_name"],
            subscription_id=data["subscription_id"],
            subscription_name=data["subscription_name"],
            tenant_id=data["tenant_id"],
            tenant_name=data["tenant_name"],
            username=data["username"],
            created_at=created_at,
        )

    @staticmethod
    def validate_context_id(context_id: str) -> bool:
        """Validate context_id format.

        Context ID must be 1-20 characters containing only:
        - Alphanumeric characters (a-z, A-Z, 0-9)
        - Hyphens (-)
        - Underscores (_)

        Args:
            context_id: The context ID to validate.

        Returns:
            True if valid, False otherwise.
        """
        if not context_id:
            return False
        return bool(re.match(Context.CONTEXT_ID_PATTERN, context_id))
