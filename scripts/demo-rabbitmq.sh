#!/usr/bin/env bash
# Demo script for showing inter-service RabbitMQ communication during a defense.
# Each section prints what's being shown and pauses; press Enter to advance.
#
# Usage: ./scripts/demo-rabbitmq.sh
set -euo pipefail

RABBIT=platform-rabbitmq
AUTH=http://localhost:8080
BOOK=http://localhost:8082

bold() { printf "\033[1m%s\033[0m\n" "$1"; }
hr()   { printf "\033[2m%s\033[0m\n" "------------------------------------------------------------"; }
pause(){ printf "\n\033[2m[Enter to continue]\033[0m"; read -r; }

# ── 0. Make sure the demo sniffer queue exists ────────────────────────
docker exec "$RABBIT" sh -c '
  rabbitmqadmin declare queue name=demo.sniffer durable=true
  rabbitmqadmin declare binding source=user.events destination=demo.sniffer routing_key="user.*"
  rabbitmqadmin declare binding source=booking.events destination=demo.sniffer routing_key="booking.*"
  rabbitmqadmin purge queue name=demo.sniffer
' >/dev/null
echo "Sniffer-Queue 'demo.sniffer' ist bereit (gebunden an user.* und booking.*)."
pause

# ── 1. Active connections per language ────────────────────────────────
bold "1) Drei Services, drei Sprachen — alle an RabbitMQ verbunden"
hr
docker exec "$RABBIT" rabbitmqctl list_connections client_properties \
  | grep -oE '"product","[^"]+"|"platform","[^"]+"' \
  | paste -d' ' - - \
  | sed -e 's/"product",//g' -e 's/"platform",//g' -e 's/"//g'
pause

# ── 2. Topology: exchanges + queues + bindings ───────────────────────
bold "2) Topologie — was ist wo deklariert?"
hr
echo "Exchanges (publisher schreiben hier rein):"
docker exec "$RABBIT" rabbitmqctl -q list_exchanges name type durable \
  | grep -vE '^(amq\.|^$)' | column -t

echo
echo "Queues (subscriber/RPC handler hängen dran):"
docker exec "$RABBIT" rabbitmqctl -q list_queues name consumers messages | column -t

echo
echo "Bindings (welche Queue hört auf welchen Exchange):"
docker exec "$RABBIT" rabbitmqctl -q list_bindings source_name routing_key destination_name \
  | column -t
pause

# ── 3. Trigger real events through the frontend's APIs ───────────────
bold "3) Aktionen auslösen — User registriert sich + bucht ein Auto"
hr

EMAIL="prof$(date +%s)@uni.de"
PW="hunter2hunter2"
echo "Register $EMAIL"
curl -s -X POST "$AUTH/api/v1/auth/register" \
  -H 'Content-Type: application/json' \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PW\"}" \
  | python3 -m json.tool

echo
echo "Login"
TOKENS=$(curl -s -X POST "$AUTH/api/v1/auth/login" \
  -H 'Content-Type: application/json' \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PW\"}")
ACCESS=$(echo "$TOKENS" | python3 -c 'import json,sys;print(json.load(sys.stdin)["access_token"])')
echo "Got access token: ${ACCESS:0:30}..."

echo
echo "Warte kurz, damit booking den user.registered-Event konsumiert..."
sleep 2

echo
echo "Booking erstellen (löst RabbitMQ-RPC zu currency-converter aus):"
BOOKING=$(curl -s -X POST "$BOOK/api/v1/bookings" \
  -H "Authorization: Bearer $ACCESS" \
  -H 'Content-Type: application/json' \
  -d '{"carId":"22222222-2222-2222-2222-222222222222","startDate":"2026-06-01","endDate":"2026-06-04","targetCurrency":"USD"}')
echo "$BOOKING" | python3 -m json.tool

echo
echo "→ totalSourceAmount = 3 Tage × 85 EUR = 255 EUR"
echo "→ totalTargetAmount kommt vom currency-converter via RabbitMQ-RPC"
pause

# ── 4. Read the events out of the sniffer queue ──────────────────────
bold "4) Was ist in den Topic-Exchanges geflossen?"
hr
docker exec "$RABBIT" rabbitmqadmin get queue=demo.sniffer count=10 ackmode=ack_requeue_true \
  | grep -E '^\|' \
  | awk -F'|' 'NR>1 { gsub(/^ +| +$/,"",$2); gsub(/^ +| +$/,"",$3); gsub(/^ +| +$/,"",$5); printf "\n● %s   (%s)\n%s\n", $2, $3, $5 }'
pause

# ── 5. Negative test: kill the Python service, see booking fail ──────
bold "5) Beweis dass die Inter-Service-Comm WIRKLICH über RabbitMQ läuft"
hr
echo "Wir stoppen den currency-converter Container."
echo "Wenn booking → currency direkt im Prozess wäre, würde nichts ändern."
echo "Wenn RabbitMQ-RPC läuft, schlägt das nächste Booking mit 503 fehl (RPC-Timeout)."
pause

docker stop platform-currency-converter >/dev/null
echo "currency-converter gestoppt. Versuch ein Booking..."
echo
curl -s -X POST "$BOOK/api/v1/bookings" \
  -H "Authorization: Bearer $ACCESS" \
  -H 'Content-Type: application/json' \
  -w "\nHTTP-Status: %{http_code}\n" \
  -d '{"carId":"33333333-3333-3333-3333-333333333333","startDate":"2026-07-01","endDate":"2026-07-03","targetCurrency":"EUR"}'

echo
echo "→ 503 + 'currency-converter timed out': der RPC-Call ist blockiert."
echo "→ Beweis: booking redet NICHT direkt mit currency-converter, sondern wartet auf RabbitMQ."
pause

echo "Service wieder hochfahren..."
docker start platform-currency-converter >/dev/null
sleep 3
echo "Booking nochmal versuchen, sollte wieder klappen:"
curl -s -X POST "$BOOK/api/v1/bookings" \
  -H "Authorization: Bearer $ACCESS" \
  -H 'Content-Type: application/json' \
  -d '{"carId":"33333333-3333-3333-3333-333333333333","startDate":"2026-07-01","endDate":"2026-07-03","targetCurrency":"EUR"}' \
  | python3 -m json.tool

echo
bold "Demo durch. RabbitMQ Management UI: http://localhost:15672 (guest/guest)"
