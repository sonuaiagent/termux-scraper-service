package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os/exec"
    "strings"
    "time"

    "github.com/tebeka/selenium"
    "github.com/tebeka/selenium/chrome"
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

var chromedriverCmd *exec.Cmd

func main() {
    http.HandleFunc("/health", handleHealth)
    http.HandleFunc("/scrape", handleScrape)

    port := ":8080"
    fmt.Println("🚀 Go Scraper Service v2.0 Started (Hybrid)")
    fmt.Println("📅 Service running on port 8080")
    fmt.Println("🛒 Flipkart scraping: Go ChromeDriver (4-6s)")
    fmt.Println("📦 Amazon scraping: Python requests (2-3s)")
    fmt.Println("🔍 Optimized for speed and reliability")

    log.Fatal(http.ListenAndServe(port, nil))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
    fmt.Println("💚 Health check requested")

    response := HealthResponse{
        Status:    "ok",
        Service:   "Hybrid Go+Python Scraper Service (Flipkart+Amazon)",
        Version:   "2.0-hybrid",
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
            Message: fmt.Sprintf("🚀 Hello from Hybrid Scraper Service v2.0! Chat ID: %s User: @%s ⚡ Ultra Fast Performance - Flipkart: Go ChromeDriver, Amazon: Python requests!", chatID, username),
            ProductInfo: map[string]interface{}{
                "title":     "Hybrid Scraper Service (Flipkart + Amazon)",
                "price":     "Active & Ultra Fast",
                "timestamp": time.Now().Format("2006-01-02 15:04:05 IST"),
            },
            Debug: []string{"Hybrid service running on port 8080", "Flipkart: Go ChromeDriver", "Amazon: Python requests"},
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
        return
    }

    if strings.Contains(url, "flipkart.com") {
        fmt.Println("🛒 Flipkart URL detected - using Go ChromeDriver...")

        productInfo, err := scrapeFlipkartProduct(url)
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

🔧 Scraped with: Go ChromeDriver`, productInfo["name"], productInfo["price"], productInfo["rating"]),
            ProductInfo: productInfo,
            Debug: []string{"Flipkart product scraped with Go ChromeDriver"},
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
        return
    }

    if strings.Contains(url, "amazon.") {
        fmt.Println("📦 Amazon URL detected - using Python requests...")

        productInfo, err := scrapeAmazonProduct(url)
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

🔧 Scraped with: Python requests`, productInfo["title"], productInfo["price"], productInfo["mrp"], productInfo["discount"], productInfo["rating"], productInfo["availability"]),
            ProductInfo: productInfo,
            Debug: []string{"Amazon product scraped with Python requests"},
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
        return
    }

    fmt.Printf("❓ Unknown command or unsupported URL: %s\n", url)

    response := ScrapeResponse{
        Success: false,
        Error:   "Send a Flipkart or Amazon product URL or 'go' command to test the service",
        Message: fmt.Sprintf("Hybrid Scraper Service v2.0 - Received: %s", command),
        Debug:   []string{fmt.Sprintf("Service running at %s", time.Now().Format("2006-01-02 15:04:05 IST"))},
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func scrapeAmazonProduct(productURL string) (map[string]string, error) {
    fmt.Println("📦 Starting Amazon scraping with Python...")
    
    // Path to Python script
    pythonScript := "/data/data/com.termux/files/home/termux-scraper-service/amazon_scraper.py"
    
    // Call Python script via exec
    cmd := exec.Command("python3", pythonScript, "--url", productURL)
    
    output, err := cmd.CombinedOutput()
    if err != nil {
        return nil, fmt.Errorf("Python script failed: %v, output: %s", err, string(output))
    }
    
    // Parse JSON output from Python script
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
    
    // Convert to map format expected by your system
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

func scrapeFlipkartProduct(productURL string) (map[string]interface{}, error) {
    fmt.Println("🚀 Starting ChromeDriver for Flipkart...")
    chromedriverPath := "/data/data/com.termux/files/usr/lib/chromium/chromedriver"
    chromedriverCmd = exec.Command(chromedriverPath, "--port=9515")
    err := chromedriverCmd.Start()
    if err != nil {
        return nil, fmt.Errorf("failed to start ChromeDriver: %v", err)
    }

    fmt.Println("⏳ ChromeDriver started, waiting...")
    time.Sleep(3 * time.Second)

    defer func() {
        if chromedriverCmd != nil && chromedriverCmd.Process != nil {
            chromedriverCmd.Process.Kill()
        }
    }()

    fmt.Println("🌐 Connecting to ChromeDriver...")

    caps := selenium.Capabilities{"browserName": "chrome"}
    chromeCaps := chrome.Capabilities{
        Args: []string{
            "--headless",
            "--disable-gpu",
            "--no-sandbox",
            "--disable-dev-shm-usage",
            "--disable-blink-features=AutomationControlled",
            "--window-size=1920,1080",
            "--disable-notifications",
            "--disable-infobars",
            "--disable-extensions",
            "--user-agent=Mozilla/5.0 (Linux; Android 10; SM-G973F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Mobile Safari/537.36",
        },
    }
    caps.AddChrome(chromeCaps)

    wd, err := selenium.NewRemote(caps, "http://localhost:9515")
    if err != nil {
        return nil, fmt.Errorf("failed to open Chrome session: %v", err)
    }
    defer wd.Quit()

    fmt.Println("✅ Connected successfully")

    fmt.Printf("📄 Loading Flipkart page: %s\n", productURL)
    if err := wd.Get(productURL); err != nil {
        return nil, fmt.Errorf("failed to load page: %v", err)
    }

    time.Sleep(3 * time.Second)

    result := make(map[string]interface{})

    fmt.Println("🔍 Extracting product name...")
    productName := "Name not found"
    nameSelectors := []string{"span.B_NuCI", "h1", "span._35KyD6"}
    for _, sel := range nameSelectors {
        elem, err := wd.FindElement(selenium.ByCSSSelector, sel)
        if err == nil {
            name, _ := elem.Text()
            if name != "" {
                productName = strings.TrimSpace(name)
                break
            }
        }
    }
    result["name"] = productName

    fmt.Println("💰 Extracting product price...")
    productPrice := "Price not found"
    priceSelectors := []string{"div._30jeq3._16Jk6d", "div.Nx9bqj.CxhGGd", "div._30jeq3", "div._1_WHN1"}
    for _, sel := range priceSelectors {
        elem, err := wd.FindElement(selenium.ByCSSSelector, sel)
        if err == nil {
            price, _ := elem.Text()
            if price != "" {
                productPrice = strings.TrimSpace(price)
                break
            }
        }
    }
    result["price"] = productPrice

    fmt.Println("⭐ Extracting product rating...")
    productRating := "Rating not available"

    ratingSelectors := []string{"div._3LWZlK", "div._1lRcqv", "span._1lRcqv", "div.gUuXy-"}
    for _, sel := range ratingSelectors {
        elem, err := wd.FindElement(selenium.ByCSSSelector, sel)
        if err == nil {
            rating, _ := elem.Text()
            if rating != "" {
                productRating = strings.TrimSpace(rating)
                break
            }
        }
    }

    result["rating"] = productRating

    fmt.Println("✅ Product information extracted successfully!")
    fmt.Printf("📦 Name: %s\n", productName)
    fmt.Printf("💰 Price: %s\n", productPrice)
    fmt.Printf("⭐ Rating: %s\n", productRating)

    return result, nil
}
