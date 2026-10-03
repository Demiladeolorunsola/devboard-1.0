#!/usr/bin/env bash
#
# End-to-end smoke test for DevBoard.
# Requires: curl, python3 (used only to parse JSON, no external libs needed).
#
# Usage:
#   chmod +x test_devboard.sh
#   ./test_devboard.sh
#
# Run this while `go run ./cmd/main.go` is running in another terminal.

set -uo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
PASS=0
FAIL=0

# Unique email each run so "register" doesn't collide with a previous run.
EMAIL="tester+$(date +%s)@example.com"
PASSWORD="password123"

# --- helpers ---------------------------------------------------------------

# json_get <json> <key>  -> prints the value, or nothing if missing
json_get() {
  python3 -c "
import json, sys
try:
    data = json.loads(sys.argv[1])
    keys = sys.argv[2].split('.')
    for k in keys:
        data = data[k]
    print(data)
except Exception:
    pass
" "$1" "$2"
}

# check <label> <expected_status> <actual_status> <body>
check() {
  local label="$1" expected="$2" actual="$3" body="$4"
  if [ "$actual" = "$expected" ]; then
    echo "  OK   $label (HTTP $actual)"
    PASS=$((PASS+1))
  else
    echo "  FAIL $label (expected $expected, got $actual)"
    echo "       body: $body"
    FAIL=$((FAIL+1))
  fi
}

# request <method> <path> <body_json> <auth_header_or_empty>
# sets globals: RESP_BODY, RESP_STATUS
request() {
  local method="$1" path="$2" body="$3" auth="$4"
  local response
  if [ -n "$auth" ]; then
    response=$(curl -s -w '\n%{http_code}' -X "$method" "$BASE_URL$path" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $auth" \
      -d "$body")
  else
    response=$(curl -s -w '\n%{http_code}' -X "$method" "$BASE_URL$path" \
      -H "Content-Type: application/json" \
      -d "$body")
  fi
  RESP_STATUS=$(echo "$response" | tail -n1)
  RESP_BODY=$(echo "$response" | sed '$d')
}

# --- tests -------------------------------------------------------------------

echo "== Health check =="
health=$(curl -s -w '\n%{http_code}' "$BASE_URL/health")
health_status=$(echo "$health" | tail -n1)
check "GET /health" "200" "$health_status" "$(echo "$health" | sed '$d')"

echo ""
echo "== Auth =="

request POST "/api/auth/register" \
  "{\"name\":\"Test User\",\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}" ""
check "POST /api/auth/register" "201" "$RESP_STATUS" "$RESP_BODY"

request POST "/api/auth/login" \
  "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}" ""
check "POST /api/auth/login" "200" "$RESP_STATUS" "$RESP_BODY"
TOKEN=$(json_get "$RESP_BODY" "token")

if [ -z "$TOKEN" ]; then
  echo ""
  echo "No token returned from login — stopping here, nothing after this point can run."
  echo "Final response body was: $RESP_BODY"
  exit 1
fi

echo ""
echo "== Profile =="

request GET "/api/profile" "" "$TOKEN"
check "GET /api/profile" "200" "$RESP_STATUS" "$RESP_BODY"

request PUT "/api/profile" '{"name":"Updated Name","email":""}' "$TOKEN"
check "PUT /api/profile" "200" "$RESP_STATUS" "$RESP_BODY"

echo ""
echo "== Projects =="

request POST "/api/projects" '{"name":"Test Project","description":"created by test script"}' "$TOKEN"
check "POST /api/projects" "201" "$RESP_STATUS" "$RESP_BODY"
PROJECT_ID=$(json_get "$RESP_BODY" "project.id")

request GET "/api/projects" "" "$TOKEN"
check "GET /api/projects" "200" "$RESP_STATUS" "$RESP_BODY"

request GET "/api/projects/$PROJECT_ID" "" "$TOKEN"
check "GET /api/projects/:project_id" "200" "$RESP_STATUS" "$RESP_BODY"

request PUT "/api/projects/$PROJECT_ID" '{"name":"","description":"","status":"active"}' "$TOKEN"
check "PUT /api/projects/:project_id" "200" "$RESP_STATUS" "$RESP_BODY"

echo ""
echo "== Tasks =="

request POST "/api/projects/$PROJECT_ID/tasks" '{"title":"Test Task","description":"first task","priority":"high"}' "$TOKEN"
check "POST /api/projects/:project_id/tasks" "201" "$RESP_STATUS" "$RESP_BODY"
TASK_ID=$(json_get "$RESP_BODY" "task.id")

request GET "/api/projects/$PROJECT_ID/tasks" "" "$TOKEN"
check "GET /api/projects/:project_id/tasks" "200" "$RESP_STATUS" "$RESP_BODY"

request GET "/api/tasks/$TASK_ID" "" "$TOKEN"
check "GET /api/tasks/:task_id" "200" "$RESP_STATUS" "$RESP_BODY"

request PUT "/api/tasks/$TASK_ID" '{"title":"","description":"","status":"in_progress","priority":"high"}' "$TOKEN"
check "PUT /api/tasks/:task_id" "200" "$RESP_STATUS" "$RESP_BODY"

echo ""
echo "== Time tracking =="

request POST "/api/time/start" "{\"project_id\":$PROJECT_ID,\"task_id\":$TASK_ID}" "$TOKEN"
check "POST /api/time/start" "201" "$RESP_STATUS" "$RESP_BODY"
TIME_ID=$(json_get "$RESP_BODY" "time_entry.id")

sleep 1

request POST "/api/time/stop/$TIME_ID" "" "$TOKEN"
check "POST /api/time/stop/:id" "200" "$RESP_STATUS" "$RESP_BODY"

request GET "/api/time" "" "$TOKEN"
check "GET /api/time" "200" "$RESP_STATUS" "$RESP_BODY"

request GET "/api/projects/$PROJECT_ID/time" "" "$TOKEN"
check "GET /api/projects/:project_id/time" "200" "$RESP_STATUS" "$RESP_BODY"

request GET "/api/tasks/$TASK_ID/time" "" "$TOKEN"
check "GET /api/tasks/:task_id/time" "200" "$RESP_STATUS" "$RESP_BODY"

echo ""
echo "== Reports =="

request GET "/api/reports/projects/$PROJECT_ID/time" "" "$TOKEN"
check "GET /api/reports/projects/:project_id/time" "200" "$RESP_STATUS" "$RESP_BODY"

request GET "/api/reports/tasks/productivity" "" "$TOKEN"
check "GET /api/reports/tasks/productivity" "200" "$RESP_STATUS" "$RESP_BODY"

request GET "/api/reports/activity" "" "$TOKEN"
check "GET /api/reports/activity" "200" "$RESP_STATUS" "$RESP_BODY"

request GET "/api/reports/dashboard" "" "$TOKEN"
check "GET /api/reports/dashboard" "200" "$RESP_STATUS" "$RESP_BODY"

echo ""
echo "== Cleanup =="

request DELETE "/api/tasks/$TASK_ID" "" "$TOKEN"
check "DELETE /api/tasks/:task_id" "200" "$RESP_STATUS" "$RESP_BODY"

request DELETE "/api/projects/$PROJECT_ID" "" "$TOKEN"
check "DELETE /api/projects/:project_id" "200" "$RESP_STATUS" "$RESP_BODY"

request DELETE "/api/profile" "" "$TOKEN"
check "DELETE /api/profile" "200" "$RESP_STATUS" "$RESP_BODY"

echo ""
echo "=================================="
echo "Passed: $PASS   Failed: $FAIL"
echo "=================================="

if [ "$FAIL" -gt 0 ]; then
  exit 1
fi