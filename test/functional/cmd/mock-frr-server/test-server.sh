#!/bin/bash

# Test script for Mock FRR Server
# This script starts the server, runs basic tests, and stops it

set -e

echo "=== Mock FRR Server Test ==="
echo

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Configuration
CONFIG_PATH="../../config/mock-frr-config.yaml"
HTTP_PORT=51051
TIMEOUT=5

# Start the server in background
echo "Starting mock FRR server..."
./mock-frr-server -config "$CONFIG_PATH" > /tmp/mock-frr-server.log 2>&1 &
SERVER_PID=$!

# Function to cleanup on exit
cleanup() {
    echo
    echo "Stopping mock FRR server (PID: $SERVER_PID)..."
    kill $SERVER_PID 2>/dev/null || true
    wait $SERVER_PID 2>/dev/null || true
    echo "Server stopped."
}

trap cleanup EXIT

# Wait for server to start
echo "Waiting for server to start..."
sleep 2

# Check if server is running
if ! kill -0 $SERVER_PID 2>/dev/null; then
    echo -e "${RED}✗ Server failed to start${NC}"
    cat /tmp/mock-frr-server.log
    exit 1
fi

echo -e "${GREEN}✓ Server started (PID: $SERVER_PID)${NC}"
echo

# Test 1: Health check
echo "Test 1: Health check"
if curl -s -f "http://localhost:$HTTP_PORT/health" > /dev/null; then
    echo -e "${GREEN}✓ Health check passed${NC}"
else
    echo -e "${RED}✗ Health check failed${NC}"
    exit 1
fi
echo

# Test 2: Get stats
echo "Test 2: Get stats"
STATS=$(curl -s "http://localhost:$HTTP_PORT/stats")
if echo "$STATS" | grep -q "total_peers"; then
    echo -e "${GREEN}✓ Stats endpoint working${NC}"
    echo "Stats: $STATS"
else
    echo -e "${RED}✗ Stats endpoint failed${NC}"
    exit 1
fi
echo

# Test 3: Add a peer
echo "Test 3: Add a BGP peer"
ADD_RESULT=$(curl -s -X POST "http://localhost:$HTTP_PORT/peers/add" \
    -H "Content-Type: application/json" \
    -d '{
        "IPAddress": "192.168.1.1",
        "ASN": 65000,
        "RemoteASN": 65001
    }')

if echo "$ADD_RESULT" | grep -q "success.*true"; then
    echo -e "${GREEN}✓ Peer added successfully${NC}"
    echo "Result: $ADD_RESULT"
else
    echo -e "${RED}✗ Failed to add peer${NC}"
    echo "Result: $ADD_RESULT"
    exit 1
fi
echo

# Test 4: List peers
echo "Test 4: List all peers"
PEERS=$(curl -s "http://localhost:$HTTP_PORT/peers")
if echo "$PEERS" | grep -q "192.168.1.1"; then
    echo -e "${GREEN}✓ Peer listed successfully${NC}"
    echo "Peers: $PEERS"
else
    echo -e "${RED}✗ Peer not found in list${NC}"
    exit 1
fi
echo

# Test 5: Wait for session establishment
echo "Test 5: Wait for session establishment (2 seconds)..."
sleep 2

# Test 6: Get session state
echo "Test 6: Get session state"
SESSION=$(curl -s "http://localhost:$HTTP_PORT/sessions/state?ip=192.168.1.1")
if echo "$SESSION" | grep -q "State"; then
    echo -e "${GREEN}✓ Session state retrieved${NC}"
    echo "Session: $SESSION"
else
    echo -e "${RED}✗ Failed to get session state${NC}"
    exit 1
fi
echo

# Test 7: Get all sessions
echo "Test 7: Get all sessions"
ALL_SESSIONS=$(curl -s "http://localhost:$HTTP_PORT/sessions")
if echo "$ALL_SESSIONS" | grep -q "192.168.1.1"; then
    echo -e "${GREEN}✓ All sessions retrieved${NC}"
    echo "Sessions: $ALL_SESSIONS"
else
    echo -e "${RED}✗ Failed to get all sessions${NC}"
    exit 1
fi
echo

# Test 8: Get running config
echo "Test 8: Get running config"
CONFIG=$(curl -s "http://localhost:$HTTP_PORT/config")
if echo "$CONFIG" | grep -q "FRR Mock Configuration"; then
    echo -e "${GREEN}✓ Running config retrieved${NC}"
    echo "Config preview:"
    echo "$CONFIG" | head -10
else
    echo -e "${RED}✗ Failed to get running config${NC}"
    exit 1
fi
echo

# Test 9: Remove peer
echo "Test 9: Remove BGP peer"
REMOVE_RESULT=$(curl -s -X POST "http://localhost:$HTTP_PORT/peers/remove" \
    -H "Content-Type: application/json" \
    -d '{"ip_address": "192.168.1.1"}')

if echo "$REMOVE_RESULT" | grep -q "success.*true"; then
    echo -e "${GREEN}✓ Peer removed successfully${NC}"
    echo "Result: $REMOVE_RESULT"
else
    echo -e "${RED}✗ Failed to remove peer${NC}"
    echo "Result: $REMOVE_RESULT"
    exit 1
fi
echo

# Test 10: Verify peer removed
echo "Test 10: Verify peer removed"
PEERS_AFTER=$(curl -s "http://localhost:$HTTP_PORT/peers")
if echo "$PEERS_AFTER" | grep -q "192.168.1.1"; then
    echo -e "${RED}✗ Peer still exists after removal${NC}"
    exit 1
else
    echo -e "${GREEN}✓ Peer successfully removed${NC}"
fi
echo

echo "=== All tests passed! ==="