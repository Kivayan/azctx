"""
azctx - A personal Python CLI tool

This is the main entry point for the azctx CLI application.
"""

try:
    from ._version import __version__
except ImportError:
    # Fallback version for development environments without installed package
    __version__ = "0.0.0+dev"
