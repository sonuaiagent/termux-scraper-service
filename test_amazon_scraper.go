package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
    "time"

    "github.com/PuerkitoBio/goquery"
)

func scrapeAmazonProduct(url string) (map[string]string, error) {
    fmt.Println("🔄 Making request to Amazon...")
    
    client := &http.Client{
        Timeout: 30 * time.Second,
    }

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %v", err)
    }

    // Use EXACT same headers as successful Python version
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
    req.Header.Set("Accept-Language", "en-US,en;q=0.9")
    req.Header.Set("Accept-Encoding", "gzip, deflate, br")

    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("request error: %v", err)
    }
    defer resp.Body.Close()

    fmt.Printf("📊 Response status: %d\n", resp.StatusCode)

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %v", err)
    }

    fmt.Printf("📐 Response length: %d characters\n", len(bodyBytes))

    doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyBytes)))
    if err != nil {
        return nil, fmt.Errorf("failed to parse HTML: %v", err)
    }

    // Helper function to get text by selector (like Python's select_one)
    getText := func(selector string) string {
        elem := doc.Find(selector).First()
        if elem.Length() > 0 {
            return strings.TrimSpace(elem.Text())
        }
        return ""
    }

    product := make(map[string]string)

    // Extract product title - EXACT same selector as Python
    fmt.Println("🔍 Extracting product title...")
    product["title"] = getText("span#productTitle")
    if product["title"] == "" {
        product["title"] = "Title Not Found"
    }

    // Extract MRP - EXACT same selectors as Python
    fmt.Println("💰 Extracting MRP...")
    mrpSelectors := []string{
        "span.a-price.a-text-price span.a-offscreen",
        "span#priceblock_mrp",
        "span.a-text-price",
        "span.a-price.a-text-price",
    }
    
    product["mrp"] = "MRP Not Found"
    for _, selector := range mrpSelectors {
        mrp := getText(selector)
        if mrp != "" {
            product["mrp"] = mrp
            break
        }
    }

    // Extract discount - EXACT same selectors as Python
    fmt.Println("💸 Extracting discount...")
    discount := getText("span.savingsPercentage")
    if discount == "" {
        discount = getText("span.a-color-price")
    }
    if discount == "" {
        discount = "Discount Not Found"
    }
    product["discount"] = discount

    // Extract rating - EXACT same selectors as Python
    fmt.Println("⭐ Extracting rating...")
    rating := getText("span.a-icon-alt")
    if rating == "" {
        rating = getText("span#acrPopover")
    }
    if rating == "" {
        rating = "Rating Not Found"
    }
    product["rating"] = rating

    // Extract availability - EXACT same selectors as Python
    fmt.Println("📦 Extracting availability...")
    availability := getText("div#availability span")
    if availability == "" {
        availability = getText("span#availability")
    }
    if availability == "" {
        availability = "Availability Not Found"
    }
    product["availability"] = availability

    // Extract price - using core price selector like Python
    fmt.Println("💲 Extracting final price...")
    priceSelectors := []string{
        "div#corePrice_feature_div span.a-offscreen",
        "span.a-price span.a-offscreen",
        "span#priceblock_ourprice",
        "span#priceblock_dealprice",
        "span#price_inside_buybox",
    }

    product["price"] = "Price Not Found"
    for _, selector := range priceSelectors {
        price := getText(selector)
        if price != "" && strings.Contains(price, "₹") {
            product["price"] = price
            break
        }
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

func main() {
    // Use EXACT same URL as successful Python test
    testURL := "https://www.amazon.in/LEOTUDE-Cottonblend-Graphic-Oversized-FS49_Navy_Boston_P_Navy/dp/B0FKGKJ41P?ref_=Oct_d_onr_d_1968123031_0&pd_rd_w=HAflx&content-id=amzn1.sym.6994f97b-af7d-4405-a303-8aac5f8b11eb&pf_rd_p=6994f97b-af7d-4405-a303-8aac5f8b11eb&pf_rd_r=7AS90A2KJN17V7RH9ADC&pd_rd_wg=9hq3I&pd_rd_r=e3e371f7-a0a9-47ce-a939-f3ba7df16660&pd_rd_i=B0FHDHSFYY&th=1&psc=1"

    fmt.Println("🧪 Testing Amazon scraper with Go net/http + goquery")
    fmt.Println("🔗 URL:", testURL)
    fmt.Println(strings.Repeat("=", 80))

    startTime := time.Now()
    result, err := scrapeAmazonProduct(testURL)
    duration := time.Since(startTime)

    if err != nil {
        fmt.Println("❌ Error:", err)
    } else {
        fmt.Println("✅ Scraping Results:")
        fmt.Println("📱 Product Title:", result["title"])
        fmt.Println("🏷️  MRP:", result["mrp"])
        fmt.Println("💸 Discount:", result["discount"])
        fmt.Println("💰 Final Price:", result["price"])
        fmt.Println("⭐ Rating:", result["rating"])
        fmt.Println("📦 Availability:", result["availability"])
        fmt.Printf("⏱️  Total time: %v\n", duration)
    }
}
