#!/data/data/com.termux/files/usr/bin/bash

echo "🔄 Restarting Enterprise Parallel Processing Scraper"
echo "📅 $(date '+%Y-%m-%d %H:%M:%S IST')"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Acquire wake lock
termux-wake-lock

# Stop existing services
echo "🛑 Stopping existing services..."
FASTAPI_PIDS=$(pgrep -f "uvicorn amazon_api:app")
GO_PIDS=$(pgrep -f "./scraper_service")

if [ ! -z "$FASTAPI_PIDS" ]; then
    for PID in $FASTAPI_PIDS; do
        echo "   Stopping FastAPI (PID: $PID)"
        kill -TERM $PID 2>/dev/null
    done
fi

if [ ! -z "$GO_PIDS" ]; then
    for PID in $GO_PIDS; do
        echo "   Stopping Go service (PID: $PID)"
        kill -TERM $PID 2>/dev/null
    done
fi

echo "⏳ Waiting 3 seconds for clean shutdown..."
sleep 3

# Navigate to project directory
cd ~/termux-scraper-service

# Create logs directory
mkdir -p logs

echo ""
echo "🚀 Starting services..."

# Start FastAPI
echo "🐍 Starting FastAPI with 4 workers..."
nohup uvicorn amazon_api:app \
  --host 0.0.0.0 \
  --port 8081 \
  --workers 4 \
  --limit-concurrency 100 \
  --access-log \
  > logs/fastapi.log 2>&1 &

FASTAPI_PID=$!
echo "✅ FastAPI started (PID: $FASTAPI_PID)"

# Start Go service
echo "🔧 Starting Go service..."
nohup ./scraper_service > logs/go_service.log 2>&1 &
GO_PID=$!
echo "✅ Go service started (PID: $GO_PID)"

# Save PIDs
echo $FASTAPI_PID > .fastapi.pid
echo $GO_PID > .go_service.pid

echo ""
echo "🔍 Performing post-restart health checks..."
sleep 5

# Health checks
FASTAPI_HEALTHY=false
GO_HEALTHY=false

if curl -s http://localhost:8081/health > /dev/null 2>&1; then
    echo "✅ FastAPI health check: PASSED"
    FASTAPI_HEALTHY=true
else
    echo "❌ FastAPI health check: FAILED"
fi

if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ Go service health check: PASSED"
    GO_HEALTHY=true
else
    echo "❌ Go service health check: FAILED"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ "$FASTAPI_HEALTHY" = true ] && [ "$GO_HEALTHY" = true ]; then
    echo "🎉 Restart completed successfully! All services healthy."
    echo "📊 System ready for concurrent processing"
elif [ "$FASTAPI_HEALTHY" = true ] || [ "$GO_HEALTHY" = true ]; then
    echo "⚠️  Restart partially successful. Check logs for issues."
else
    echo "❌ Restart failed. Both services unhealthy."
fi

echo ""
read -p "Press Enter to exit..."
