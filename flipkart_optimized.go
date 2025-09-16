package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

type FlipkartResult struct {
	Success bool   `json:"success"`
	Name    string `json:"name"`
	Price   string `json:"price"`
	Rating  string `json:"rating"`
	Error   string `json:"error,omitempty"`
}

var chromedriverCmd *exec.Cmd

func scrapeFlipkartProduct(url string) FlipkartResult {
	fmt.Fprintf(os.Stderr, "🚀 Starting OPTIMIZED ChromeDriver for Flipkart...\n")
	chromedriverPath := "/data/data/com.termux/files/usr/lib/chromium/chromedriver"
	chromedriverCmd = exec.Command(chromedriverPath, "--port=9515")
	err := chromedriverCmd.Start()
	if err != nil {
		return FlipkartResult{Success: false, Error: fmt.Sprintf("failed to start ChromeDriver: %v", err)}
	}

	fmt.Fprintf(os.Stderr, "⏳ ChromeDriver started, waiting...\n")
	time.Sleep(3 * time.Second)

	defer func() {
		if chromedriverCmd != nil && chromedriverCmd.Process != nil {
			chromedriverCmd.Process.Kill()
		}
	}()

	fmt.Fprintf(os.Stderr, "🌐 Connecting to ChromeDriver...\n")

	caps := selenium.Capabilities{"browserName": "chrome"}
	chromeCaps := chrome.Capabilities{
		Args: []string{
			"--headless",
			"--disable-gpu",
			"--no-sandbox",
			"--disable-dev-shm-usage",
			"--disable-blink-features=AutomationControlled",
			"--window-size=1920,1080",
			// ✅ SPEED OPTIMIZATIONS - Keep CSS but block heavy resources
			"--disable-images",                    // Block images (major speed gain)
			"--disable-plugins",                   // No plugins
			"--disable-extensions",                // No extensions
			"--no-first-run",                     // Skip setup
			"--disable-default-apps",             // No default apps
			"--single-process",                   // Single process mode
			"--disable-background-timer-throttling",
			"--disable-renderer-backgrounding",
			"--disable-backgrounding-occluded-windows",
			"--disable-ipc-flooding-protection",  // Faster IPC
			"--no-zygote",                        // No zygote process
		},
	}
	caps.AddChrome(chromeCaps)

	wd, err := selenium.NewRemote(caps, "http://localhost:9515")
	if err != nil {
		return FlipkartResult{Success: false, Error: fmt.Sprintf("failed to open Chrome session: %v", err)}
	}
	defer wd.Quit()

	fmt.Fprintf(os.Stderr, "✅ Connected successfully\n")
	fmt.Fprintf(os.Stderr, "📄 Loading Flipkart page...\n")
	
	if err := wd.Get(url); err != nil {
		return FlipkartResult{Success: false, Error: fmt.Sprintf("failed to load page: %v", err)}
	}

	time.Sleep(3 * time.Second)

	// ----- Product Name -----
	fmt.Fprintf(os.Stderr, "🔍 Extracting product name...\n")
	productName := "Name not found"
	nameSelectors := []string{"span.B_NuCI", "h1"}
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

	// ----- Product Price -----
	fmt.Fprintf(os.Stderr, "💰 Extracting product price...\n")
	productPrice := "Price not found"
	priceSelectors := []string{"div._30jeq3._16Jk6d", "div.Nx9bqj.CxhGGd", "div._30jeq3"}
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

	// ----- Product Rating -----
	fmt.Fprintf(os.Stderr, "⭐ Extracting product rating...\n")
	productRating := "Rating not available"

	currentURL, err := wd.CurrentURL()
	if err == nil {
		pid, lid := "", ""

		if p := strings.Index(currentURL, "pid="); p != -1 {
			end := strings.Index(currentURL[p:], "&")
			if end != -1 {
				pid = currentURL[p+4 : p+end]
			} else {
				pid = currentURL[p+4:]
			}
		}

		if l := strings.Index(currentURL, "lid="); l != -1 {
			end := strings.Index(currentURL[l:], "&")
			if end != -1 {
				lid = currentURL[l+4 : l+end]
			} else {
				lid = currentURL[l+4:]
			}
		}

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

	fmt.Fprintf(os.Stderr, "✅ Optimized Selenium scraping completed!\n")
	fmt.Fprintf(os.Stderr, "📦 Name: %s\n", productName)
	fmt.Fprintf(os.Stderr, "💰 Price: %s\n", productPrice)
	fmt.Fprintf(os.Stderr, "⭐ Rating: %s\n", productRating)

	return FlipkartResult{
		Success: true,
		Name:    productName,
		Price:   productPrice,
		Rating:  productRating,
	}
}

func main() {
	if len(os.Args) < 2 {
		result := FlipkartResult{Success: false, Error: "Usage: go run flipkart_optimized.go <url>"}
		jsonOutput, _ := json.Marshal(result)
		fmt.Println(string(jsonOutput))
		os.Exit(1)
	}
	
	url := os.Args[1]
	result := scrapeFlipkartProduct(url)
	
	jsonOutput, _ := json.Marshal(result)
	fmt.Println(string(jsonOutput))
}
