#!/data/data/com.termux/files/usr/bin/bash

echo "🛑 Stopping Enterprise Parallel Processing Scraper"
echo "📅 $(date '+%Y-%m-%d %H:%M:%S IST')"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

STOPPED_SERVICES=0

# Stop FastAPI
echo "🐍 Stopping FastAPI service..."
FASTAPI_PIDS=$(pgrep -f "uvicorn amazon_api:app")

if [ ! -z "$FASTAPI_PIDS" ]; then
    for PID in $FASTAPI_PIDS; do
        echo "   Killing FastAPI worker (PID: $PID)"
        kill -TERM $PID 2>/dev/null || kill -9 $PID 2>/dev/null
    done
    
    # Wait for graceful shutdown
    sleep 2
    echo "✅ FastAPI stopped"
    STOPPED_SERVICES=$((STOPPED_SERVICES + 1))
else
    echo "ℹ️  FastAPI not running"
fi

echo ""

# Stop Go service
echo "🔧 Stopping Go service..."
GO_PIDS=$(pgrep -f "./scraper_service")

if [ ! -z "$GO_PIDS" ]; then
    for PID in $GO_PIDS; do
        echo "   Killing Go service (PID: $PID)"
        kill -TERM $PID 2>/dev/null || kill -9 $PID 2>/dev/null
    done
    
    sleep 2
    echo "✅ Go service stopped"
    STOPPED_SERVICES=$((STOPPED_SERVICES + 1))
else
    echo "ℹ️  Go service not running"
fi

# Clean up PID files
cd ~/termux-scraper-service
rm -f .fastapi.pid .go_service.pid

# Release wake lock
termux-wake-unlock

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ $STOPPED_SERVICES -gt 0 ]; then
    echo "✅ Stopped $STOPPED_SERVICES service(s) successfully"
    echo "📝 Logs preserved in logs/ directory"
else
    echo "ℹ️  No services were running"
fi

echo ""
read -p "Press Enter to exit..."
