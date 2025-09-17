package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os/exec"
    "strings"
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
    fmt.Println("🚀 Go Scraper Service v2.0 Started (Optimal Hybrid)")
    fmt.Println("📅 Service running on port 8080")
    fmt.Println("🛒 Flipkart scraping: Go ChromeDriver (flipkart.go)")
    fmt.Println("📦 Amazon scraping: Python requests (amazon.py)")
    fmt.Println("⚡ Best performance for each platform!")

    log.Fatal(http.ListenAndServe(port, nil))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
    fmt.Println("💚 Health check requested")

    response := HealthResponse{
        Status:    "ok",
        Service:   "Optimal Hybrid Scraper Service (Go+Python)",
        Version:   "2.0-optimal-fixed",
        Timestamp: time.Now().Format("2006-01-02 15:04:05 IST"),
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func handleScrape(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        fmt.Println("❌ Invalid method:", r.Method)
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var requestData map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
        fmt.Println("❌ Bad request body:", err)
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }

    command, _ := requestData["command"].(string)
    url, _ := requestData["url"].(string)
    chatID, _ := requestData["chat_id"].(string)
    username, _ := requestData["username"].(string)

    fmt.Printf("🤖 Scrape request: command='%s', user='%s', chatID='%s'\n", command, username, chatID)

    if command == "go" || url == "go" {
        fmt.Println("⚡ Processing 'go' command - sending ultra-fast response!")

        response := ScrapeResponse{
            Success: true,
            Message: fmt.Sprintf("🚀 Hello from Optimal Hybrid Scraper v2.0! Chat ID: %s User: @%s ⚡ Flipkart: Go ChromeDriver, Amazon: Python requests - FIXED!", chatID, username),
            ProductInfo: map[string]interface{}{
                "title":     "Optimal Hybrid Scraper Service - FIXED",
                "price":     "Active & Ultra Fast",
                "timestamp": time.Now().Format("2006-01-02 15:04:05 IST"),
            },
            Debug: []string{"Service running on port 8080", "Flipkart: flipkart.go", "Amazon: amazon.py", "JSON parsing FIXED"},
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
        return
    }

    // Route to Flipkart scraper (no URL validation - trust worker.js)
    if strings.Contains(url, "flipkart.com") || strings.Contains(url, "dl.flipkart.com") {
        fmt.Println("🛒 Flipkart URL detected - using flipkart.go...")

        productInfo, err := scrapeFlipkartViaGo(url)
        if err != nil {
            fmt.Printf("❌ Flipkart scraping failed: %v\n", err)
            response := ScrapeResponse{
                Success: false,
                Error:   fmt.Sprintf("Flipkart scraping failed: %v", err),
                Message: "Could not extract product information",
            }
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(response)
            return
        }

        fmt.Printf("✅ Flipkart product scraped successfully: %s\n", productInfo["name"])

        response := ScrapeResponse{
            Success: true,
            Message: fmt.Sprintf(`🛒 Flipkart Product Found!

📱 Name: %s
💰 Price: %s
⭐ Rating: %s

🔧 Scraped with: Go ChromeDriver (flipkart.go)`, productInfo["name"], productInfo["price"], productInfo["rating"]),
            ProductInfo: productInfo,
            Debug: []string{"Flipkart product scraped with Go ChromeDriver"},
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
        return
    }

    // Default to Amazon scraper for all other URLs (no validation - trust worker.js)
    fmt.Println("📦 Amazon URL assumed - using amazon.py...")

    productInfo, err := scrapeAmazonViaPython(url)
    if err != nil {
        fmt.Printf("❌ Amazon scraping failed: %v\n", err)
        response := ScrapeResponse{
            Success: false,
            Error:   fmt.Sprintf("Amazon scraping failed: %v", err),
            Message: "Could not extract product information",
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
        return
    }

    fmt.Printf("✅ Amazon product scraped successfully: %s\n", productInfo["title"])

    response := ScrapeResponse{
        Success: true,
        Message: fmt.Sprintf(`📦 Amazon Product Found!

📱 Name: %s
💰 Price: %s
🏷️ MRP: %s
💸 Discount: %s
⭐ Rating: %s
📦 Availability: %s

🔧 Scraped with: Python requests (amazon.py)`, productInfo["title"], productInfo["price"], productInfo["mrp"], productInfo["discount"], productInfo["rating"], productInfo["availability"]),
        ProductInfo: productInfo,
        Debug: []string{"Amazon product scraped with Python requests"},
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func scrapeFlipkartViaGo(productURL string) (map[string]string, error) {
    fmt.Println("🛒 Starting Flipkart scraping via Go script...")
    
    flipkartScript := "/data/data/com.termux/files/home/termux-scraper-service/flipkart.go"
    cmd := exec.Command("go", "run", flipkartScript, productURL)
    
    output, err := cmd.Output()
    if err != nil {
        return nil, fmt.Errorf("Go script failed: %v", err)
    }
    
    var result struct {
        Success bool   `json:"success"`
        Name    string `json:"name"`
        Price   string `json:"price"`
        Rating  string `json:"rating"`
        Error   string `json:"error"`
    }
    
    if err := json.Unmarshal(output, &result); err != nil {
        return nil, fmt.Errorf("failed to parse Go output: %v, raw output: %s", err, string(output))
    }
    
    if !result.Success {
        return nil, fmt.Errorf("Go scraping failed: %s", result.Error)
    }
    
    product := map[string]string{
        "name":   result.Name,
        "price":  result.Price,
        "rating": result.Rating,
    }
    
    fmt.Println("✅ Flipkart product extraction completed via Go!")
    fmt.Printf("📦 Name: %s\n", product["name"])
    fmt.Printf("💰 Price: %s\n", product["price"])
    fmt.Printf("⭐ Rating: %s\n", product["rating"])
    
    return product, nil
}

func scrapeAmazonViaPython(productURL string) (map[string]string, error) {
    fmt.Println("📦 Starting Amazon scraping via Python script...")
    
    pythonScript := "/data/data/com.termux/files/home/termux-scraper-service/amazon.py"
    cmd := exec.Command("python3", pythonScript, productURL)
    
    output, err := cmd.Output()
    if err != nil {
        return nil, fmt.Errorf("Python script failed: %v", err)
    }
    
    var result struct {
        Success      bool   `json:"success"`
        Title        string `json:"title"`
        MRP          string `json:"mrp"`
        Discount     string `json:"discount"`
        Price        string `json:"price"`
        Rating       string `json:"rating"`
        Availability string `json:"availability"`
        Error        string `json:"error"`
    }
    
    if err := json.Unmarshal(output, &result); err != nil {
        return nil, fmt.Errorf("failed to parse Python output: %v, raw output: %s", err, string(output))
    }
    
    if !result.Success {
        return nil, fmt.Errorf("Python scraping failed: %s", result.Error)
    }
    
    product := map[string]string{
        "title":        result.Title,
        "mrp":          result.MRP,
        "discount":     result.Discount,
        "price":        result.Price,
        "rating":       result.Rating,
        "availability": result.Availability,
    }
    
    fmt.Println("✅ Amazon product extraction completed via Python!")
    fmt.Printf("📱 Title: %s\n", product["title"])
    fmt.Printf("💰 Price: %s\n", product["price"])
    fmt.Printf("🏷️  MRP: %s\n", product["mrp"])
    fmt.Printf("💸 Discount: %s\n", product["discount"])
    fmt.Printf("⭐ Rating: %s\n", product["rating"])
    fmt.Printf("📦 Availability: %s\n", product["availability"])
    
    return product, nil
}
