#!/data/data/com.termux/files/usr/bin/bash

echo "📊 Enterprise Parallel Processing Scraper Status"
echo "📅 $(date '+%Y-%m-%d %H:%M:%S IST')"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Check FastAPI status
echo "🐍 FastAPI Service Status:"
FASTAPI_PIDS=$(pgrep -f "uvicorn amazon_api:app")

if [ ! -z "$FASTAPI_PIDS" ]; then
    echo "   Status: ✅ RUNNING"
    echo "   PIDs: $FASTAPI_PIDS"
    echo "   Port: 8081"
    echo "   Workers: 4 (parallel processing enabled)"
    
    # Health check
    if curl -s http://localhost:8081/health > /dev/null 2>&1; then
        echo "   Health: ✅ HEALTHY"
    else
        echo "   Health: ❌ UNHEALTHY (not responding)"
    fi
    
    # Resource usage
    for PID in $FASTAPI_PIDS; do
        CPU_MEM=$(ps -o pid,pcpu,pmem --no-headers -p $PID 2>/dev/null)
        if [ ! -z "$CPU_MEM" ]; then
            echo "   Resources (PID $PID): $CPU_MEM"
        fi
    done
else
    echo "   Status: ❌ NOT RUNNING"
    echo "   Port: 8081 (not listening)"
fi

echo ""

# Check Go service status
echo "🔧 Go Service Status:"
GO_PIDS=$(pgrep -f "./scraper_service")

if [ ! -z "$GO_PIDS" ]; then
    echo "   Status: ✅ RUNNING"
    echo "   PIDs: $GO_PIDS"
    echo "   Port: 8080"
    echo "   Features: Concurrent monitoring"
    
    # Health check
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        echo "   Health: ✅ HEALTHY"
    else
        echo "   Health: ❌ UNHEALTHY (not responding)"
    fi
    
    # Resource usage
    for PID in $GO_PIDS; do
        CPU_MEM=$(ps -o pid,pcpu,pmem --no-headers -p $PID 2>/dev/null)
        if [ ! -z "$CPU_MEM" ]; then
            echo "   Resources (PID $PID): $CPU_MEM"
        fi
    done
else
    echo "   Status: ❌ NOT RUNNING"
    echo "   Port: 8080 (not listening)"
fi

echo ""

# Network status
echo "🌐 Network Status:"
FASTAPI_LISTEN=$(netstat -tlnp 2>/dev/null | grep :8081 | head -1)
GO_LISTEN=$(netstat -tlnp 2>/dev/null | grep :8080 | head -1)

if [ ! -z "$FASTAPI_LISTEN" ]; then
    echo "   FastAPI Port 8081: ✅ LISTENING"
else
    echo "   FastAPI Port 8081: ❌ NOT LISTENING"
fi

if [ ! -z "$GO_LISTEN" ]; then
    echo "   Go Service Port 8080: ✅ LISTENING"
else
    echo "   Go Service Port 8080: ❌ NOT LISTENING"
fi

# Connection counts
FASTAPI_CONNECTIONS=$(netstat -an 2>/dev/null | grep :8081 | grep ESTABLISHED | wc -l)
GO_CONNECTIONS=$(netstat -an 2>/dev/null | grep :8080 | grep ESTABLISHED | wc -l)

echo "   Active connections FastAPI: $FASTAPI_CONNECTIONS"
echo "   Active connections Go: $GO_CONNECTIONS"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Overall system status
FASTAPI_RUNNING=$(pgrep -f "uvicorn amazon_api:app" | wc -l)
GO_RUNNING=$(pgrep -f "./scraper_service" | wc -l)

if [ $FASTAPI_RUNNING -gt 0 ] && [ $GO_RUNNING -gt 0 ]; then
    echo "🎉 System Status: ✅ FULLY OPERATIONAL"
    echo "🚀 Ready for concurrent users"
elif [ $FASTAPI_RUNNING -gt 0 ] || [ $GO_RUNNING -gt 0 ]; then
    echo "⚠️  System Status: ⚠️ PARTIALLY OPERATIONAL"
    echo "🔧 Some services need attention"
else
    echo "❌ System Status: ❌ NOT OPERATIONAL"
    echo "🚀 Use start_scraper.sh to start services"
fi

echo ""
read -p "Press Enter to exit..."
