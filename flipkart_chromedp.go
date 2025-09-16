package main

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "strings"
    "time"

    "github.com/chromedp/chromedp"
)

type FlipkartResult struct {
    Success bool   `json:"success"`
    Name    string `json:"name"`
    Price   string `json:"price"`
    Rating  string `json:"rating"`
    Error   string `json:"error,omitempty"`
}

func scrapeFlipkartChromedp(url string) FlipkartResult {
    fmt.Fprintf(os.Stderr, "🚀 Starting chromedp for Flipkart...\n")
    
    // Create optimized browser options
    opts := append(chromedp.DefaultExecAllocatorOptions[:],
        chromedp.DisableGPU,
        chromedp.NoDefaultBrowserCheck,
        chromedp.Flag("disable-images", true),
        chromedp.Flag("disable-javascript", false), // Keep JS for dynamic content
        chromedp.Flag("disable-plugins", true),
        chromedp.Flag("disable-extensions", true),
        chromedp.Flag("no-sandbox", true),
        chromedp.Flag("disable-dev-shm-usage", true),
    )
    
    allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
    defer cancel()
    
    ctx, cancel := chromedp.NewContext(allocCtx)
    defer cancel()
    
    // Set timeout for entire operation
    ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
    defer cancel()
    
    var name, price, rating string
    
    err := chromedp.Run(ctx,
        // Navigate to page
        chromedp.Navigate(url),
        
        // Wait for and extract title
        chromedp.WaitVisible("span.B_NuCI", chromedp.ByQuery),
        chromedp.Text("span.B_NuCI", &name, chromedp.ByQuery),
        
        // Extract price
        chromedp.Text("div._30jeq3._16Jk6d", &price, chromedp.ByQuery),
        
        // Try to extract rating with fallback selectors
        chromedp.ActionFunc(func(ctx context.Context) error {
            selectors := []string{"div._3LWZlK", "div._1lRcqv", "span._1lRcqv"}
            for _, sel := range selectors {
                if err := chromedp.Text(sel, &rating, chromedp.ByQuery).Do(ctx); err == nil && rating != "" {
                    break
                }
            }
            return nil
        }),
    )
    
    if err != nil {
        fmt.Fprintf(os.Stderr, "❌ chromedp error: %v\n", err)
        return FlipkartResult{Success: false, Error: err.Error()}
    }
    
    fmt.Fprintf(os.Stderr, "✅ chromedp scraping completed!\n")
    fmt.Fprintf(os.Stderr, "📦 Name: %s\n", strings.TrimSpace(name))
    fmt.Fprintf(os.Stderr, "💰 Price: %s\n", strings.TrimSpace(price))
    fmt.Fprintf(os.Stderr, "⭐ Rating: %s\n", strings.TrimSpace(rating))
    
    return FlipkartResult{
        Success: true,
        Name:    strings.TrimSpace(name),
        Price:   strings.TrimSpace(price),
        Rating:  strings.TrimSpace(rating),
    }
}

func main() {
    if len(os.Args) < 2 {
        result := FlipkartResult{Success: false, Error: "Usage: go run flipkart_chromedp.go <url>"}
        jsonOutput, _ := json.Marshal(result)
        fmt.Println(string(jsonOutput))
        os.Exit(1)
    }
    
    url := os.Args[1]
    result := scrapeFlipkartChromedp(url)
    
    jsonOutput, _ := json.Marshal(result)
    fmt.Println(string(jsonOutput))
}
