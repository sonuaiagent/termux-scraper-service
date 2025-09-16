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
    fmt.Fprintf(os.Stderr, "🚀 Starting optimized chromedp for Flipkart...\n")
    
    // More aggressive optimization for Termux
    opts := append(chromedp.DefaultExecAllocatorOptions[:],
        chromedp.DisableGPU,
        chromedp.NoSandbox,
        chromedp.Flag("headless", true),
        chromedp.Flag("disable-images", true),
        chromedp.Flag("disable-css", true),           // Block CSS
        chromedp.Flag("disable-plugins", true),
        chromedp.Flag("disable-extensions", true),
        chromedp.Flag("disable-dev-shm-usage", true),
        chromedp.Flag("single-process", true),        // Use single process
        chromedp.Flag("no-zygote", true),            // No zygote process
        chromedp.Flag("disable-background-timer-throttling", true),
        chromedp.Flag("disable-renderer-backgrounding", true),
    )
    
    allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
    defer cancel()
    
    ctx, cancel := chromedp.NewContext(allocCtx)
    defer cancel()
    
    // ✅ INCREASE TIMEOUT for Termux (30s instead of 15s)
    ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
    defer cancel()
    
    var name, price, rating string
    
    // Use simpler, more direct approach
    err := chromedp.Run(ctx,
        chromedp.Navigate(url),
        chromedp.Sleep(3*time.Second),  // Give page time to load
        
        // Extract name with fallback
        chromedp.ActionFunc(func(ctx context.Context) error {
            selectors := []string{"span.B_NuCI", "h1", "span._35KyD6"}
            for _, sel := range selectors {
                if err := chromedp.Text(sel, &name, chromedp.ByQuery).Do(ctx); err == nil && name != "" {
                    break
                }
            }
            return nil
        }),
        
        // Extract price with fallback
        chromedp.ActionFunc(func(ctx context.Context) error {
            selectors := []string{"div._30jeq3._16Jk6d", "div.Nx9bqj.CxhGGd", "div._30jeq3"}
            for _, sel := range selectors {
                if err := chromedp.Text(sel, &price, chromedp.ByQuery).Do(ctx); err == nil && price != "" {
                    break
                }
            }
            return nil
        }),
        
        // Extract rating with fallback
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
        result := FlipkartResult{Success: false, Error: "Usage: go run flipkart_chromedp_fixed.go <url>"}
        jsonOutput, _ := json.Marshal(result)
        fmt.Println(string(jsonOutput))
        os.Exit(1)
    }
    
    url := os.Args[1]
    result := scrapeFlipkartChromedp(url)
    
    jsonOutput, _ := json.Marshal(result)
    fmt.Println(string(jsonOutput))
}
