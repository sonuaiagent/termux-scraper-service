#!/bin/bash

echo "🛑 Stopping all scraper services..."
echo "⏰ Timestamp: $(date)"

# Kill the process running on port 8081 (FastAPI server)
PID=$(lsof -ti:8081 2>/dev/null)
if [ -z "$PID" ]; then
  echo "✅ No process running on port 8081 (FastAPI)"
else
  echo "🔥 Stopping FastAPI server on port 8081, PID: $PID"
  kill -9 $PID
  echo "✅ FastAPI server stopped"
fi

# Kill the Go scraper service running on port 8080
PID_GO=$(lsof -ti:8080 2>/dev/null)
if [ -z "$PID_GO" ]; then
  echo "✅ No process running on port 8080 (Go Service)"
else
  echo "🔥 Stopping Go scraper service on port 8080, PID: $PID_GO"
  kill -9 $PID_GO
  echo "✅ Go scraper service stopped"
fi

# Alternative method using pkill
pkill -f "uvicorn amazon_api" 2>/dev/null && echo "🔥 Killed uvicorn processes"
pkill -f "scraper_service" 2>/dev/null && echo "🔥 Killed scraper_service processes"

echo "🎯 All services stopped successfully!"
echo "⏰ Completed at: $(date)"
