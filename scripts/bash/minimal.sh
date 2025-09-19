#!/usr/bin/env bash
set -euo pipefail

# Абсолютные пути с учётом пробелов в директориях
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"   # подняться из scripts/bash на корень
BIN="$ROOT_DIR/bin/shell"
VFS="$ROOT_DIR/vfs_variants/minimal.json"
SCRIPT="$ROOT_DIR/tmp/_tmp_minimal.script"

mkdir -p "$(dirname "$VFS")" "$(dirname "$SCRIPT")" "$ROOT_DIR/bin"

# Собрать бинарь, если нет
if [[ ! -x "$BIN" ]]; then
  echo "Building binary at: $BIN"
  go build -buildvcs=false -o "$BIN" "$ROOT_DIR/cmd/shell"
fi

# VFS (минимальный)
cat > "$VFS" <<'JSON'
{
  "name": "/",
  "type": "dir",
  "children": [
    {
      "name": "home",
      "type": "dir",
      "children": [
        {
          "name": "readme.txt",
          "type": "file",
          "contentText": "Hello from minimal VFS!"
        }
      ]
    }
  ]
}
JSON

# Стартовый скрипт
cat > "$SCRIPT" <<'SCRIPT'
# Проверяем структуру
ls
cd /home
ls
cat readme.txt
exit
SCRIPT

# Запуск
"$BIN" -vfs "$VFS" -script "$SCRIPT"