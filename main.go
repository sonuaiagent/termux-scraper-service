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
    fmt.Println("🚀 Go Scraper Service v1.0 Started")
    fmt.Println("📅 Service running on port 8080")
    fmt.Println("🛒 Flipkart scraping enabled")
    fmt.Println("🔍 ChromeDriver ready")

    log.Fatal(http.ListenAndServe(port, nil))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
    fmt.Println("💚 Health check requested")

    response := HealthResponse{
        Status:    "ok",
        Service:   "Go Termux Scraper Service with Flipkart Support",
        Version:   "1.0",
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
            Message: fmt.Sprintf("🚀 Hello from Go Scraper Service! Chat ID: %s User: @%s ⚡ Ultra Fast Performance Ready for Flipkart scraping!", chatID, username),
            ProductInfo: map[string]interface{}{
                "title":     "Go Scraper Service with Flipkart Support",
                "price":     "Active & Ultra Fast",
                "timestamp": time.Now().Format("2006-01-02 15:04:05 IST"),
            },
            Debug: []string{"Go service running on port 8080", "Flipkart scraping ready"},
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
        return
    }

    if strings.Contains(url, "flipkart.com") {
        fmt.Println("🛒 Flipkart URL detected - starting scraping process...")

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
            Message: fmt.Sprintf(`📦 Flipkart Product Found!

📱 Name: %s
💰 Price: %s
⭐ Rating: %s`, productInfo["name"], productInfo["price"], productInfo["rating"]),
            ProductInfo: productInfo,
            Debug: []string{"Flipkart product scraped successfully"},
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
        return
    }

    fmt.Printf("❓ Unknown command or unsupported URL: %s\n", url)

    response := ScrapeResponse{
        Success: false,
        Error:   "Send a Flipkart product URL or 'go' command to test the service",
        Message: fmt.Sprintf("Go Scraper Service - Received: %s", command),
        Debug:   []string{fmt.Sprintf("Service running at %s", time.Now().Format("2006-01-02 15:04:05 IST"))},
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func scrapeFlipkartProduct(productURL string) (map[string]interface{}, error) {
    fmt.Println("🚀 Starting ChromeDriver...")
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

    if productRating == "Rating not available" {
        currentURL, err := wd.CurrentURL()
        if err == nil {
            pid, lid := "", ""

            if p := strings.Index(currentURL, "pid="); p != -1 {
                end := strings.IndexAny(currentURL[p:], "&")
                if end != -1 {
                    pid = currentURL[p+4 : p+end]
                } else {
                    pid = currentURL[p+4:]
                }
            }

            if l := strings.Index(currentURL, "lid="); l != -1 {
                end := strings.IndexAny(currentURL[l:], "&")
                if end != -1 {
                    lid = currentURL[l+4 : l+end]
                } else {
                    lid = currentURL[l+4:]
                }
            }

            if pid != "" && lid != "" {
                ratingElems, err := wd.FindElements(selenium.ByXPATH, "//*[starts-with(@id, 'productRating_')]")
                if err == nil {
                    for _, elem := range ratingElems {
                        elemID, _ := elem.GetAttribute("id")
                        if strings.Contains(elemID, pid) && strings.Contains(elemID, lid) {
                            text, _ := elem.Text()
                            if text != "" {
                                productRating = strings.TrimSpace(text)
                                break
                            }
                        }
                    }
                }
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
