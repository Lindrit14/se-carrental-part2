"""Shared pytest fixtures + auto-codegen for the gRPC stubs.

Generated proto stubs aren't committed (see .gitignore). To keep `pytest` a
one-step command for new contributors, regenerate them at collection time if
they're missing.
"""
from __future__ import annotations

import subprocess
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent
GENERATED_PB2 = (
    REPO_ROOT / "app" / "infrastructure" / "grpc" / "generated"
    / "currency" / "v1" / "currency_pb2.py"
)


def _ensure_proto_stubs() -> None:
    if GENERATED_PB2.exists():
        return
    script = REPO_ROOT / "scripts" / "gen-proto.sh"
    subprocess.run(["bash", str(script)], check=True, cwd=REPO_ROOT)


_ensure_proto_stubs()
