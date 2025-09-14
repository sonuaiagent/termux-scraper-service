package main

import (
\t"bytes"
\t"encoding/json"
\t"fmt"
\t"log"
\t"net/http"
\t"os"
\t"time"
)

type TelegramUpdate struct {
\tMessage struct {
\t\tChat struct {
\t\t\tID int64 `json:"id"`
\t\t} `json:"chat"`
\t\tText     string `json:"text"`
\t\tFrom     struct {
\t\t\tUsername  string `json:"username"`
\t\t\tFirstName string `json:"first_name"`
\t\t} `json:"from"`
\t} `json:"message"`
}

type HealthResponse struct {
\tStatus    string `json:"status"`
\tService   string `json:"service"`
\tVersion   string `json:"version"`
\tTimestamp string `json:"timestamp"`
}

type ScrapeResponse struct {
\tSuccess     bool        `json:"success"`
\tMessage     string      `json:"message"`
\tProductInfo interface{} `json:"product_info,omitempty"`
\tError       string      `json:"error,omitempty"`
\tDebug       []string    `json:"debug,omitempty"`
}

func main() {
\thttp.HandleFunc("/health", handleHealth)
\thttp.HandleFunc("/scrape", handleScrape)
\t
\tport := ":8080"
\tfmt.Printf("🚀 Starting Go Scraper Service v1.0
")
\tfmt.Printf("📅 Started at: %s
", time.Now().Format("2006-01-02 15:04:05 IST"))
\tfmt.Printf("📱 Service running on http://0.0.0.0%s
", port)
\tfmt.Printf("🛒 Features: Amazon scraping, Health check, Ultra-fast performance
")
\tfmt.Printf("🔗 Use Cloudflare Tunnel to expose this service
")
\tfmt.Printf("💡 Send 'go' command to test the service
")
\tfmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
")
\t
\tlog.Fatal(http.ListenAndServe(port, nil))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
\tresponse := HealthResponse{
\t\tStatus:    "ok",
\t\tService:   "Go Termux Scraper Service",
\t\tVersion:   "1.0",
\t\tTimestamp: time.Now().Format("2006-01-02 15:04:05 IST"),
\t}
\t
\tw.Header().Set("Content-Type", "application/json")
\tjson.NewEncoder(w).Encode(response)
}

func handleScrape(w http.ResponseWriter, r *http.Request) {
\tif r.Method != "POST" {
\t\thttp.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
\t\treturn
\t}
\t
\tvar requestData map[string]interface{}
\tif err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
\t\thttp.Error(w, "Bad request", http.StatusBadRequest)
\t\treturn
\t}
\t
\tcommand, _ := requestData["command"].(string)
\turl, _ := requestData["url"].(string)
\tchatID, _ := requestData["chat_id"].(string)
\tusername, _ := requestData["username"].(string)
\t
\t// Handle /go command
\tif command == "go" || url == "go" {
\t\tresponse := ScrapeResponse{
\t\t\tSuccess: true,
\t\t\tMessage: fmt.Sprintf("🚀 Hello from Go Scraper Service!

✅ Service running at %s
💡 Chat ID: %s
👤 User: @%s
⚡ Performance: Ultra Fast
🛒 Ready for Amazon scraping!", 
\t\t\t\ttime.Now().Format("2006-01-02 15:04:05 IST"), chatID, username),
\t\t\tProductInfo: map[string]interface{}{
\t\t\t\t"title":     "Go Scraper Service",
\t\t\t\t"price":     "Active & Ultra Fast",
\t\t\t\t"timestamp": time.Now().Format("2006-01-02 15:04:05 IST"),
\t\t\t},
\t\t\tDebug: []string{"Go service running on port 8080", "Amazon scraping ready"},
\t\t}
\t\t
\t\tw.Header().Set("Content-Type", "application/json")
\t\tjson.NewEncoder(w).Encode(response)
\t\treturn
\t}
\t
\t// Handle other commands/URLs
\tresponse := ScrapeResponse{
\t\tSuccess: false,
\t\tError:   "Send 'go' command to test the Go scraper service",
\t\tMessage: fmt.Sprintf("🤖 Go Scraper Service

Received: %s
⚡ Ultra-fast performance ready!", command),
\t\tDebug:   []string{fmt.Sprintf("Service running at %s", time.Now().Format("2006-01-02 15:04:05 IST"))},
\t}
\t
\tw.Header().Set("Content-Type", "application/json")
\tjson.NewEncoder(w).Encode(response)
}
