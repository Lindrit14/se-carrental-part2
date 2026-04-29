#!/usr/bin/env bash
# Codegen for currency.proto -> Python (grpcio + protobuf).
# Output is committed to keep the runtime image free of build tools.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
PROTO_DIR="$ROOT/proto"
OUT_DIR="$ROOT/app/infrastructure/grpc/generated"

mkdir -p "$OUT_DIR"
touch "$OUT_DIR/__init__.py"

python -m grpc_tools.protoc \
    -I "$PROTO_DIR" \
    --python_out="$OUT_DIR" \
    --grpc_python_out="$OUT_DIR" \
    "$PROTO_DIR/currency/v1/currency.proto"

# protoc emits absolute imports (`from currency.v1 import currency_pb2 as ...`)
# which only resolve if `proto/` is on sys.path. We rewrite to relative imports
# rooted at the generated package so it works as a normal Python package.
GRPC_PY="$OUT_DIR/currency/v1/currency_pb2_grpc.py"
if [ -f "$GRPC_PY" ]; then
    # macOS/BSD sed needs the empty -i argument
    sed -i '' 's|^from currency.v1 import currency_pb2 as|from . import currency_pb2 as|' "$GRPC_PY" 2>/dev/null \
      || sed -i 's|^from currency.v1 import currency_pb2 as|from . import currency_pb2 as|' "$GRPC_PY"
fi

# Make every emitted package directory a Python package.
find "$OUT_DIR" -type d -exec touch {}/__init__.py \;

echo "Generated proto stubs in $OUT_DIR"
