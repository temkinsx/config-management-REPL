#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"   # подняться из scripts/bash на корень
BIN="$ROOT_DIR/bin/shell"
VFS="$ROOT_DIR/vfs_variants/deeptree.json"
SCRIPT="$ROOT_DIR/tmp/_tmp_deeptree.script"

mkdir -p "$(dirname "$VFS")" "$(dirname "$SCRIPT")" "$ROOT_DIR/bin"

# Собрать бинарь, если нет
if [[ ! -x "$BIN" ]]; then
  echo "Building binary at: $BIN"
  go build -buildvcs=false -o "$BIN" "$ROOT_DIR/cmd/shell"
fi

# VFS (глубокое дерево, минимум 3 уровня)
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
          "name": "user",
          "type": "dir",
          "children": [
            { "name": "docs", "type": "dir",
              "children": [
                { "name": "notes.txt", "type": "file", "contentText": "line1\nline1\nline2" },
                { "name": "guide.md",  "type": "file", "contentText": "# Guide\nUse VFS" }
              ]
            },
            { "name": "projects", "type": "dir",
              "children": [
                { "name": "demo", "type": "dir",
                  "children": [
                    { "name": "main.go", "type": "file", "contentText": "package main\nfunc main(){}" }
                  ]
                }
              ]
            }
          ]
        }
      ]
    },
    {
      "name": "var",
      "type": "dir",
      "children": [
        { "name": "data", "type": "dir",
          "children": [
            { "name": "cache.db", "type": "file", "contentBase64": "AAEC" }
          ]
        }
      ]
    }
  ]
}
JSON

# Стартовый скрипт (включает ошибочные команды для демонстрации)
cat > "$SCRIPT" <<'SCRIPT'
# Начало демонстрации (нормальные операции)
ls
cd /home/user/docs
ls
cat notes.txt
uniq notes.txt
cat notes.txt

# Ошибочные операции для демонстрации обработки ошибок:
# 1) Переход в несуществующий каталог
cd /no/such
# 2) Попытка прочитать несуществующий файл
cat missing.txt
# 3) Попытка cd в файл
cd notes.txt
# 4) Вызов uniq на директории
uniq /home/user
# 5) Некорректная загрузка VFS (не существует файл)
vfs-load /no/such/path.json
# 6) Попытка touch с некорректным именем
touch ../..
# 7) Попытка создать файл в несуществующем пути
touch /no/such/dir/newfile.txt

# Вернуться к рабочему сценарию
cd /home/user/projects/demo
ls
touch README
ls
cd /
ls

exit
SCRIPT

# Запуск
"$BIN" -vfs "$VFS" -script "$SCRIPT"