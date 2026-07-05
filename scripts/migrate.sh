#!/usr/bin/env bash
# Применяет .sql миграции из backend/migrations по порядку, ровно один раз каждую.
# Работает как с локальным psql, так и (если его нет) через docker exec в контейнер postgres.
#
# Использование:
#   ./scripts/migrate.sh            # применить все неприменённые миграции
#   ./scripts/migrate.sh status     # показать, что применено, а что нет
#
# Переменные окружения (необязательно, иначе берутся из backend/.env):
#   POSTGRES_DSN            - строка подключения (postgres://user:pass@host:port/db?sslmode=disable)
#   POSTGRES_CONTAINER      - имя docker-контейнера с postgres (по умолчанию: life-postgres)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
MIGRATIONS_DIR="$ROOT_DIR/backend/migrations"
ENV_FILE="$ROOT_DIR/backend/.env"
CONTAINER="${POSTGRES_CONTAINER:-life-postgres}"
MODE="${1:-apply}"

if [ -z "${POSTGRES_DSN:-}" ] && [ -f "$ENV_FILE" ]; then
  set -a
  # shellcheck disable=SC1090
  source "$ENV_FILE"
  set +a
fi

DSN="${POSTGRES_DSN:-postgres://postgres:postgres@localhost:5432/life?sslmode=disable}"

if [ ! -d "$MIGRATIONS_DIR" ]; then
  echo "Директория с миграциями не найдена: $MIGRATIONS_DIR" >&2
  exit 1
fi

# --- Выбираем способ выполнения SQL: локальный psql или docker exec ---
if command -v psql >/dev/null 2>&1; then
  RUNNER="local"
  echo "Использую локальный psql, DSN: $DSN"
elif command -v docker >/dev/null 2>&1 && docker ps --format '{{.Names}}' | grep -qx "$CONTAINER"; then
  RUNNER="docker"
  echo "Локальный psql не найден, использую docker exec -> контейнер '$CONTAINER'"
else
  echo "Не найден ни локальный psql, ни запущенный контейнер '$CONTAINER'." >&2
  echo "Запустите 'make db-up' или установите psql, либо укажите POSTGRES_CONTAINER=<имя>." >&2
  exit 1
fi

# psql_exec_file FILE  — выполнить файл целиком в одной транзакции (-1), падать на первой ошибке
psql_exec_file() {
  local file="$1"
  if [ "$RUNNER" = "local" ]; then
    psql "$DSN" -v ON_ERROR_STOP=1 -q -1 -f "$file"
  else
    docker exec -i "$CONTAINER" psql -U postgres -d life -v ON_ERROR_STOP=1 -q -1 < "$file"
  fi
}

# psql_exec_sql SQL — выполнить строку SQL, вернуть stdout (используется для служебных запросов)
psql_exec_sql() {
  local sql="$1"
  if [ "$RUNNER" = "local" ]; then
    psql "$DSN" -t -A -v ON_ERROR_STOP=1 -q -c "$sql"
  else
    docker exec -i "$CONTAINER" psql -U postgres -d life -t -A -v ON_ERROR_STOP=1 -q -c "$sql"
  fi
}

# --- Таблица учёта применённых миграций ---
psql_exec_sql "
CREATE TABLE IF NOT EXISTS schema_migrations (
    version    TEXT PRIMARY KEY,
    applied_at TIMESTAMP NOT NULL DEFAULT NOW()
);" >/dev/null

if [ "$MODE" = "status" ]; then
  echo
  echo "Файл миграции                | Статус"
  echo "------------------------------|-----------"
  for file in "$MIGRATIONS_DIR"/*.sql; do
    name="$(basename "$file")"
    applied="$(psql_exec_sql "SELECT 1 FROM schema_migrations WHERE version = '$name';" | tr -d '[:space:]')"
    if [ "$applied" = "1" ]; then
      printf "%-30s| applied\n" "$name"
    else
      printf "%-30s| pending\n" "$name"
    fi
  done
  exit 0
fi

applied_any=false

for file in "$MIGRATIONS_DIR"/*.sql; do
  name="$(basename "$file")"
  already="$(psql_exec_sql "SELECT 1 FROM schema_migrations WHERE version = '$name';" | tr -d '[:space:]')"

  if [ "$already" = "1" ]; then
    echo "skip   $name (уже применена)"
    continue
  fi

  echo "apply  $name"
  psql_exec_file "$file"
  psql_exec_sql "INSERT INTO schema_migrations(version) VALUES ('$name');" >/dev/null
  applied_any=true
done

echo
if [ "$applied_any" = false ]; then
  echo "Все миграции уже применены — схема актуальна."
else
  echo "Миграции успешно применены."
fi
