#!/data/data/com.termux/files/usr/bin/bash

echo "🚀 Starting Enterprise Parallel Processing Scraper v2.1"
echo "📅 $(date '+%Y-%m-%d %H:%M:%S IST')"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Navigate to project directory
cd ~/termux-scraper-service

# Check if services are already running
FASTAPI_PID=$(pgrep -f "uvicorn amazon_api:app")
GO_PID=$(pgrep -f "./scraper_service")

if [ ! -z "$FASTAPI_PID" ]; then
    echo "⚠️  FastAPI already running (PID: $FASTAPI_PID)"
    echo "   Use stop_scraper.sh first if you want to restart"
    read -p "Press Enter to exit..."
    exit 1
fi

if [ ! -z "$GO_PID" ]; then
    echo "⚠️  Go service already running (PID: $GO_PID)"
    echo "   Use stop_scraper.sh first if you want to restart"
    read -p "Press Enter to exit..."
    exit 1
fi

# Create logs directory
mkdir -p logs

echo ""
echo "🐍 Starting FastAPI with 4 workers (parallel processing)..."

# Acquire wake lock to prevent Android sleep
termux-wake-lock

# Start FastAPI in background with logging
nohup uvicorn amazon_api:app \
  --host 0.0.0.0 \
  --port 8081 \
  --workers 4 \
  --limit-concurrency 100 \
  --access-log \
  > logs/fastapi.log 2>&1 &

FASTAPI_PID=$!
echo "✅ FastAPI started (PID: $FASTAPI_PID)"

# Wait for FastAPI to initialize
echo ""
echo "⏳ Waiting for FastAPI to initialize..."
sleep 3

# Check if FastAPI is responding
if curl -s http://localhost:8081/health > /dev/null 2>&1; then
    echo "✅ FastAPI health check passed"
else
    echo "⚠️  FastAPI health check failed (may still be starting)"
fi

echo ""
echo "🔧 Starting Go service (concurrent request monitor)..."

# Start Go service in background with logging
nohup ./scraper_service > logs/go_service.log 2>&1 &
GO_PID=$!
echo "✅ Go service started (PID: $GO_PID)"

# Save PIDs for management
echo $FASTAPI_PID > .fastapi.pid
echo $GO_PID > .go_service.pid

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎉 Enterprise Parallel Processing Scraper Started Successfully!"
echo "📊 FastAPI: http://localhost:8081 (4 workers)"
echo "🔧 Go Service: http://localhost:8080"
echo "📝 Logs: logs/fastapi.log, logs/go_service.log"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
read -p "Press Enter to exit..."
