#!/usr/bin/env bash
# Полная проверка life-assistant API.
# Использование: ./test_api.sh [BASE_URL]
# По умолчанию BASE_URL=http://localhost:8080

set -uo pipefail

BASE="${1:-http://localhost:8080}/api"
PASS=0
FAIL=0

# ---- утилиты ----

req() {
  # req METHOD PATH [BODY]
  local method="$1" path="$2" body="${3:-}"
  if [ -n "$body" ]; then
    curl -s -w "\n%{http_code}" -X "$method" "$BASE$path" \
      -H "Content-Type: application/json" -d "$body"
  else
    curl -s -w "\n%{http_code}" -X "$method" "$BASE$path"
  fi
}

check() {
  # check "название" ожидаемый_код факт_код тело
  local name="$1" expected="$2" actual="$3" body="$4"
  if [ "$actual" = "$expected" ]; then
    echo "OK   [$name] -> $actual"
    PASS=$((PASS+1))
  else
    echo "FAIL [$name] -> ожидали $expected, получили $actual"
    echo "     тело: $body"
    FAIL=$((FAIL+1))
  fi
}

extract() {
  # extract JSON_STRING KEY -> достаёт число/строку по простому ключу верхнего уровня
  echo "$1" | grep -o "\"$2\":[ ]*[\"0-9.]*" | head -1 | sed -E 's/.*:[ ]*"?([^",}]*)"?/\1/'
}

echo "=== Базовый URL: $BASE ==="
echo

# ---- 1. Health ----
resp=$(req GET /health); code=$(echo "$resp" | tail -1); body=$(echo "$resp" | sed '$d')
check "GET /health" 200 "$code" "$body"

resp=$(req GET /health/db); code=$(echo "$resp" | tail -1); body=$(echo "$resp" | sed '$d')
check "GET /health/db" 200 "$code" "$body"
echo

# ---- 2. Scenario CRUD ----
resp=$(req POST /scenario '{"title":"Test scenario","description":"desc","category":"test","start_node_id":0}')
code=$(echo "$resp" | tail -1); body=$(echo "$resp" | sed '$d')
check "POST /scenario" 200 "$code" "$body"
SCENARIO_ID=$(extract "$body" "id")
echo "  -> scenario_id=$SCENARIO_ID"

resp=$(req GET "/scenario?id=$SCENARIO_ID"); code=$(echo "$resp" | tail -1); body=$(echo "$resp" | sed '$d')
check "GET /scenario?id=" 200 "$code" "$body"

resp=$(req GET /scenarios); code=$(echo "$resp" | tail -1); body=$(echo "$resp" | sed '$d')
check "GET /scenarios" 200 "$code" "$body"
echo

# ---- 3. Node CRUD (используем scenario_id=0, т.к. движок сейчас читает только его — см. пояснение в чате) ----
resp=$(req POST /node '{"scenario_id":0,"type":"question","text":"Вы изменили паспорт?","order":1}')
code=$(echo "$resp" | tail -1); body=$(echo "$resp" | sed '$d')
check "POST /node (start)" 200 "$code" "$body"
NODE1_ID=$(extract "$body" "id")
echo "  -> node1_id=$NODE1_ID"

resp=$(req POST /node '{"scenario_id":0,"type":"action","text":"Уведомите банк","action_type":"notify","organization":"Bank","order":2}')
code=$(echo "$resp" | tail -1); body=$(echo "$resp" | sed '$d')
check "POST /node (action)" 200 "$code" "$body"
NODE2_ID=$(extract "$body" "id")
echo "  -> node2_id=$NODE2_ID"

resp=$(req GET "/node?id=$NODE1_ID"); code=$(echo "$resp" | tail -1); body=$(echo "$resp" | sed '$d')
check "GET /node?id=" 200 "$code" "$body"

resp=$(req GET "/nodes?scenario_id=0"); code=$(echo "$resp" | tail -1); body=$(echo "$resp" | sed '$d')
check "GET /nodes?scenario_id=" 200 "$code" "$body"
echo

# ---- 4. Edge CRUD ----
resp=$(req POST /edge "{\"from_node\":$NODE1_ID,\"to_node\":$NODE2_ID,\"condition\":\"yes\"}")
code=$(echo "$resp" | tail -1); body=$(echo "$resp" | sed '$d')
check "POST /edge" 200 "$code" "$body"
EDGE_ID=$(extract "$body" "id")
echo "  -> edge_id=$EDGE_ID"

resp=$(req GET "/edge?id=$EDGE_ID"); code=$(echo "$resp" | tail -1); body=$(echo "$resp" | sed '$d')
check "GET /edge?id=" 200 "$code" "$body"

resp=$(req GET "/edges?scenario_id=0"); code=$(echo "$resp" | tail -1); body=$(echo "$resp" | sed '$d')
check "GET /edges?scenario_id=" 200 "$code" "$body"
echo

echo "!!! ВАЖНО: движок (GraphEngine) собирается один раз при старте сервера"
echo "!!! из nodeRepo.GetByScenario(0) / edgeRepo.GetByScenario(0)."
echo "!!! Если узлы/рёбра выше созданы ПОСЛЕ запуска сервера — перезапустите"
echo "!!! сервер прямо сейчас и запустите скрипт заново, иначе шаги ниже упадут."
echo

# ---- 5. Engine flow ----
resp=$(req POST /user/scenario/start '{"user_id":1,"scenario_id":0}')
code=$(echo "$resp" | tail -1); body=$(echo "$resp" | sed '$d')
check "POST /user/scenario/start" 200 "$code" "$body"
STATE_ID=$(extract "$body" "state_id")
echo "  -> state_id=$STATE_ID"

resp=$(req POST /user/scenario/step "{\"user_state_id\":$STATE_ID,\"answer\":\"yes\"}")
code=$(echo "$resp" | tail -1); body=$(echo "$resp" | sed '$d')
check "POST /user/scenario/step" 200 "$code" "$body"

resp=$(req GET "/user/scenario/status?user_id=1&scenario_id=0")
code=$(echo "$resp" | tail -1); body=$(echo "$resp" | sed '$d')
check "GET /user/scenario/status" 200 "$code" "$body"
echo

# ---- 6. Проверка ошибок (негативные сценарии) ----
resp=$(req GET "/node?id=999999"); code=$(echo "$resp" | tail -1); body=$(echo "$resp" | sed '$d')
check "GET /node?id=999999 (ожидаем 404)" 404 "$code" "$body"

resp=$(req GET "/scenario?id=abc"); code=$(echo "$resp" | tail -1); body=$(echo "$resp" | sed '$d')
check "GET /scenario?id=abc (ожидаем 400)" 400 "$code" "$body"

resp=$(req POST /user/scenario/step '{"user_state_id":999999,"answer":"yes"}')
code=$(echo "$resp" | tail -1); body=$(echo "$resp" | sed '$d')
check "POST /step с несуществующим state (ожидаем 404)" 404 "$code" "$body"
echo

echo "=== ИТОГО: PASS=$PASS FAIL=$FAIL ==="
[ "$FAIL" -eq 0 ]
