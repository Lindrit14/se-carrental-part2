"""JSON logging setup."""
from __future__ import annotations

import logging

from pythonjsonlogger import jsonlogger


def configure_logging(level: str = "info") -> None:
    root = logging.getLogger()
    # Clear existing handlers (uvicorn installs its own)
    for h in list(root.handlers):
        root.removeHandler(h)

    handler = logging.StreamHandler()
    handler.setFormatter(
        jsonlogger.JsonFormatter(
            "%(asctime)s %(levelname)s %(name)s %(message)s",
            rename_fields={"asctime": "time", "levelname": "level"},
        )
    )
    root.addHandler(handler)
    root.setLevel(level.upper())
