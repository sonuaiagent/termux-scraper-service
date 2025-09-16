package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os/exec"
    "strconv"
    "strings"
    "time"

    "github.com/PuerkitoBio/goquery"
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
    fmt.Println("🚀 Go Scraper Service v2.0 Started")
    fmt.Println("📅 Service running on port 8080")
    fmt.Println("🛒 Flipkart scraping enabled")
    fmt.Println("📦 Amazon scraping enabled")
    fmt.Println("🔍 ChromeDriver ready")

    log.Fatal(http.ListenAndServe(port, nil))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
    fmt.Println("💚 Health check requested")

    response := HealthResponse{
        Status:    "ok",
        Service:   "Go Termux Scraper Service with Flipkart & Amazon Support",
        Version:   "2.0",
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
            Message: fmt.Sprintf("🚀 Hello from Go Scraper Service v2.0! Chat ID: %s User: @%s ⚡ Ultra Fast Performance Ready for Flipkart & Amazon scraping!", chatID, username),
            ProductInfo: map[string]interface{}{
                "title":     "Go Scraper Service with Flipkart & Amazon Support",
                "price":     "Active & Ultra Fast",
                "timestamp": time.Now().Format("2006-01-02 15:04:05 IST"),
            },
            Debug: []string{"Go service running on port 8080", "Flipkart & Amazon scraping ready"},
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
            Message: fmt.Sprintf(`🛒 Flipkart Product Found!

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

    if strings.Contains(url, "amazon.") {
        fmt.Println("📦 Amazon URL detected - starting scraping process...")

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
🏷️  MRP: %s
💸 Discount: %s
⭐ Rating: %s
📦 Availability: %s`, productInfo["title"], productInfo["price"], productInfo["mrp"], productInfo["discount"], productInfo["rating"], productInfo["availability"]),
            ProductInfo: productInfo,
            Debug: []string{"Amazon product scraped successfully"},
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
        return
    }

    fmt.Printf("❓ Unknown command or unsupported URL: %s\n", url)

    response := ScrapeResponse{
        Success: false,
        Error:   "Send a Flipkart or Amazon product URL or 'go' command to test the service",
        Message: fmt.Sprintf("Go Scraper Service v2.0 - Received: %s", command),
        Debug:   []string{fmt.Sprintf("Service running at %s", time.Now().Format("2006-01-02 15:04:05 IST"))},
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func scrapeAmazonProduct(productURL string) (map[string]string, error) {
    fmt.Println("📦 Starting Amazon scraping...")

    client := &http.Client{
        Timeout: 30 * time.Second,
    }

    req, err := http.NewRequest("GET", productURL, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %v", err)
    }

    // Set headers to mimic browser
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
    req.Header.Set("Accept-Language", "en-US,en;q=0.9")
    req.Header.Set("Accept-Encoding", "gzip, deflate, br")

    fmt.Println("🌐 Making request to Amazon...")
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("request error: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to parse HTML: %v", err)
    }

    fmt.Println("🔍 Parsing Amazon product page...")

    // Helper function to get trimmed text by selector
    getText := func(selector string) string {
        sel := doc.Find(selector).First()
        return strings.TrimSpace(sel.Text())
    }

    product := make(map[string]string)

    // Extract Title
    fmt.Println("📱 Extracting product title...")
    product["title"] = getText("span#productTitle")
    if product["title"] == "" {
        product["title"] = "Title Not Found"
    }

    // Extract MRP
    fmt.Println("💰 Extracting MRP...")
    mrpSelectors := []string{
        "span.a-price.a-text-price span.a-offscreen",
        "span#priceblock_mrp",
        "span.a-text-price",
        "span.a-price.a-text-price",
    }
    product["mrp"] = "MRP Not Found"
    for _, sel := range mrpSelectors {
        mrp := getText(sel)
        if mrp != "" {
            product["mrp"] = mrp
            break
        }
    }

    // Extract Discount
    fmt.Println("💸 Extracting discount...")
    discount := getText("span.savingsPercentage")
    if discount == "" {
        discount = getText("span.a-color-price")
    }
    if discount == "" {
        discount = "Discount Not Found"
    }
    product["discount"] = discount

    // Extract Rating
    fmt.Println("⭐ Extracting rating...")
    rating := getText("span.a-icon-alt")
    if rating == "" {
        rating = getText("span#acrPopover")
    }
    if rating == "" {
        rating = "Rating Not Found"
    }
    product["rating"] = rating

    // Extract Availability
    fmt.Println("📦 Extracting availability...")
    availability := getText("div#availability span")
    if availability == "" {
        availability = getText("span#availability")
    }
    if availability == "" {
        availability = "Availability Not Found"
    }
    product["availability"] = availability

    // Extract all possible prices
    fmt.Println("💲 Extracting all prices...")
    allPrices := make([]int, 0)

    priceSelectors := []string{
        "span.a-price span.a-offscreen",
        "span#priceblock_ourprice",
        "span#priceblock_dealprice",
        "span#price_inside_buybox",
        "span.a-color-price",
        "span.offer-price",
        "span.a-price-whole",
        "span.a-price-fraction",
        "div#corePrice_feature_div span.a-offscreen",
        "span#priceblock_saleprice",
        "span#priceblock_regularprice",
        "div#averageCustomerReviews span.a-price",
    }

    for _, sel := range priceSelectors {
        doc.Find(sel).Each(func(i int, s *goquery.Selection) {
            text := strings.TrimSpace(s.Text())
            if len(text) > 0 && strings.HasPrefix(text, "₹") {
                // Remove ₹ and commas, parse value
                clean := strings.ReplaceAll(strings.ReplaceAll(text, "₹", ""), ",", "")
                clean = strings.Split(clean, ".")[0]
                val, err := strconv.Atoi(clean)
                if err == nil {
                    allPrices = append(allPrices, val)
                }
            }
        })
    }

    // Find core price
    corePrice := getText("div#corePrice_feature_div span.a-offscreen")
    if corePrice == "" {
        corePrice = "Price Not Found"
    }
    product["price"] = corePrice

    // Determine highest price as MRP if none found earlier
    if product["mrp"] == "MRP Not Found" && len(allPrices) > 0 {
        maxPrice := 0
        for _, p := range allPrices {
            if p > maxPrice {
                maxPrice = p
            }
        }
        product["mrp"] = fmt.Sprintf("₹%d", maxPrice)
    }

    fmt.Println("✅ Amazon product extraction completed!")
    fmt.Printf("📱 Title: %s\n", product["title"])
    fmt.Printf("💰 Price: %s\n", product["price"])
    fmt.Printf("🏷️  MRP: %s\n", product["mrp"])
    fmt.Printf("💸 Discount: %s\n", product["discount"])
    fmt.Printf("⭐ Rating: %s\n", product["rating"])
    fmt.Printf("📦 Availability: %s\n", product["availability"])

    return product, nil
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

    result["rating"] = productRating

    fmt.Println("✅ Product information extracted successfully!")
    fmt.Printf("📦 Name: %s\n", productName)
    fmt.Printf("💰 Price: %s\n", productPrice)
    fmt.Printf("⭐ Rating: %s\n", productRating)

    return result, nil
}
