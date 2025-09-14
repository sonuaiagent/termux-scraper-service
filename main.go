package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
)

type HealthResponse struct {
    Status    string `json:"status"`
    Service   string `json:"service"`
    Version   string `json:"version"`
    Timestamp string `json:"timestamp"`
}

type ScrapeResponse struct {
    Success     bool        `json:"success"`
    Message     string      `json:"message"`
    ProductInfo interface{} `json:"product_info,omitempty"`
    Error       string      `json:"error,omitempty"`
    Debug       []string    `json:"debug,omitempty"`
}

func main() {
    http.HandleFunc("/health", handleHealth)
    http.HandleFunc("/scrape", handleScrape)

    port := ":8080"
    fmt.Println("🚀 Go Scraper Service v1.0 Started")
    fmt.Println("📅 Service running on port 8080")
    fmt.Println("💡 Send 'go' command to test")

    log.Fatal(http.ListenAndServe(port, nil))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
    response := HealthResponse{
        Status:    "ok",
        Service:   "Go Termux Scraper Service",
        Version:   "1.0",
        Timestamp: time.Now().Format("2006-01-02 15:04:05 IST"),
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func handleScrape(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var requestData map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }

    command, _ := requestData["command"].(string)
    url, _ := requestData["url"].(string)
    chatID, _ := requestData["chat_id"].(string)
    username, _ := requestData["username"].(string)

    if command == "go" || url == "go" {
        response := ScrapeResponse{
            Success: true,
            Message: fmt.Sprintf("🚀 Hello from Go Scraper Service! Chat ID: %s User: @%s ⚡ Ultra Fast Performance Ready for Amazon scraping!", chatID, username),
            ProductInfo: map[string]interface{}{
                "title":     "Go Scraper Service",
                "price":     "Active & Ultra Fast",
                "timestamp": time.Now().Format("2006-01-02 15:04:05 IST"),
            },
            Debug: []string{"Go service running on port 8080", "Amazon scraping ready"},
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
        return
    }

    response := ScrapeResponse{
        Success: false,
        Error:   "Send 'go' command to test the Go scraper service",
        Message: fmt.Sprintf("🤖 Go Scraper Service - Received: %s", command),
        Debug:   []string{fmt.Sprintf("Service running at %s", time.Now().Format("2006-01-02 15:04:05 IST"))},
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
