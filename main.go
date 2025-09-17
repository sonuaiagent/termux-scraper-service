package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
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
	Timing      TimingInfo  `json:"timing,omitempty"`
}

type TimingInfo struct {
	StartTime     string `json:"start_time"`
	EndTime       string `json:"end_time"`
	Duration      string `json:"duration"`
	ScrapingMethod string `json:"scraping_method"`
	ConcurrentRequests int `json:"concurrent_requests"`
	RequestID     string `json:"request_id"`
}

var (
	requestCounter int
	requestMutex   sync.Mutex
	activeRequests = make(map[string]time.Time)
)

func main() {
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/scrape", handleScrape)

	port := ":8080"
	fmt.Println("🚀 Go Scraper Service v2.1 Started (Parallel Processing Monitor)")
	fmt.Println("📅 Service running on port 8080")
	fmt.Println("🛒 Flipkart scraping: Go ChromeDriver (flipkart.go)")
	fmt.Println("📦 Amazon scraping: FastAPI (amazon_api.py)")
	fmt.Println("⚡ Parallel processing monitoring ENABLED!")

	log.Fatal(http.ListenAndServe(port, nil))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("💚 Health check requested")

	response := HealthResponse{
		Status:    "ok",
		Service:   "Parallel Processing Monitor Service (Go+FastAPI)",
		Version:   "2.1-parallel-monitor",
		Timestamp: time.Now().Format("2006-01-02 15:04:05 IST"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleScrape(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	
	// Generate unique request ID
	requestMutex.Lock()
	requestCounter++
	requestID := fmt.Sprintf("REQ_%d_%d", requestCounter, startTime.Unix())
	activeRequests[requestID] = startTime
	concurrentCount := len(activeRequests)
	requestMutex.Unlock()

	defer func() {
		requestMutex.Lock()
		delete(activeRequests, requestID)
		requestMutex.Unlock()
	}()

	fmt.Printf("🔥 [%s] NEW REQUEST - Concurrent requests: %d\n", requestID, concurrentCount)

	if r.Method != "POST" {
		fmt.Printf("❌ [%s] Invalid method: %s\n", requestID, r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		fmt.Printf("❌ [%s] Bad request body: %v\n", requestID, err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	command, _ := requestData["command"].(string)
	url, _ := requestData["url"].(string)
	chatID, _ := requestData["chat_id"].(string)
	username, _ := requestData["username"].(string)

	fmt.Printf("🤖 [%s] Scrape request: user='%s', chat='%s', concurrent=%d\n", requestID, username, chatID, concurrentCount)

	if command == "go" || url == "go" {
		fmt.Printf("⚡ [%s] Processing 'go' command\n", requestID)

		duration := time.Since(startTime)
		timing := TimingInfo{
			StartTime:          startTime.Format("15:04:05.000"),
			EndTime:            time.Now().Format("15:04:05.000"),
			Duration:           duration.String(),
			ScrapingMethod:     "Go Service Test",
			ConcurrentRequests: concurrentCount,
			RequestID:          requestID,
		}

		response := ScrapeResponse{
			Success: true,
			Message: fmt.Sprintf("🚀 Parallel Processing Monitor v2.1!\n⚡ FastAPI ENABLED\n🔥 Request ID: %s\n📊 Concurrent: %d requests\n⏱️ Response time: %s", 
				requestID, concurrentCount, duration),
			ProductInfo: map[string]interface{}{
				"title":     "Parallel Processing Monitor - ACTIVE",
				"price":     fmt.Sprintf("Response in %s", duration),
				"timestamp": time.Now().Format("2006-01-02 15:04:05 IST"),
			},
			Debug: []string{
				fmt.Sprintf("Request ID: %s", requestID),
				fmt.Sprintf("Concurrent requests: %d", concurrentCount),
				"FastAPI enabled for Amazon",
				"Go ChromeDriver for Flipkart",
			},
			Timing: timing,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	var scrapingMethod string
	var productInfo map[string]string
	var err error

	if strings.Contains(url, "flipkart.com") || strings.Contains(url, "dl.flipkart.com") {
		fmt.Printf("🛒 [%s] Flipkart URL detected - using Go ChromeDriver\n", requestID)
		scrapingMethod = "Go ChromeDriver (flipkart.go)"
		productInfo, err = scrapeFlipkartViaGo(url, requestID)
	} else {
		fmt.Printf("📦 [%s] Amazon URL assumed - using FastAPI\n", requestID)
		scrapingMethod = "FastAPI (amazon_api.py)"
		productInfo, err = scrapeAmazonViaFastAPI(url, requestID)
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	timing := TimingInfo{
		StartTime:          startTime.Format("15:04:05.000"),
		EndTime:            endTime.Format("15:04:05.000"),
		Duration:           duration.String(),
		ScrapingMethod:     scrapingMethod,
		ConcurrentRequests: concurrentCount,
		RequestID:          requestID,
	}

	fmt.Printf("✅ [%s] COMPLETED in %s - Concurrent was: %d\n", requestID, duration, concurrentCount)

	if err != nil {
		fmt.Printf("❌ [%s] Scraping failed: %v\n", requestID, err)
		response := ScrapeResponse{
			Success: false,
			Error:   fmt.Sprintf("Scraping failed: %v", err),
			Message: fmt.Sprintf("❌ Request failed\n🔍 Method: %s\n⏱️ Duration: %s\n📊 Concurrent: %d", 
				scrapingMethod, duration, concurrentCount),
			Timing: timing,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Build success message with timing and method info
	var productMessage string
	if strings.Contains(scrapingMethod, "FastAPI") {
		productMessage = fmt.Sprintf("📦 Amazon Product Found!\n\n📱 Name: %s\n💰 Price: %s\n🏷️ MRP: %s\n💸 Discount: %s\n⭐ Rating: %s\n📦 Availability: %s\n\n🔧 Method: %s ⚡\n⏱️ Duration: %s\n📊 Concurrent: %d requests\n🆔 ID: %s",
			productInfo["title"], productInfo["price"], productInfo["mrp"], 
			productInfo["discount"], productInfo["rating"], productInfo["availability"],
			scrapingMethod, duration, concurrentCount, requestID)
	} else {
		productMessage = fmt.Sprintf("🛒 Flipkart Product Found!\n\n📱 Name: %s\n💰 Price: %s\n⭐ Rating: %s\n\n🔧 Method: %s\n⏱️ Duration: %s\n📊 Concurrent: %d requests\n🆔 ID: %s",
			productInfo["name"], productInfo["price"], productInfo["rating"],
			scrapingMethod, duration, concurrentCount, requestID)
	}

	response := ScrapeResponse{
		Success:     true,
		Message:     productMessage,
		ProductInfo: productInfo,
		Debug: []string{
			fmt.Sprintf("Method: %s", scrapingMethod),
			fmt.Sprintf("Duration: %s", duration),
			fmt.Sprintf("Concurrent requests: %d", concurrentCount),
			fmt.Sprintf("Request ID: %s", requestID),
		},
		Timing: timing,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func scrapeFlipkartViaGo(productURL, requestID string) (map[string]string, error) {
	fmt.Printf("🛒 [%s] Starting Flipkart scraping via Go ChromeDriver\n", requestID)

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
		return nil, fmt.Errorf("failed to parse Go output: %v", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("Go scraping failed: %s", result.Error)
	}

	product := map[string]string{
		"name":   result.Name,
		"price":  result.Price,
		"rating": result.Rating,
	}

	fmt.Printf("✅ [%s] Flipkart scraping completed via Go ChromeDriver\n", requestID)
	return product, nil
}

func scrapeAmazonViaFastAPI(productURL, requestID string) (map[string]string, error) {
	fmt.Printf("📦 [%s] Starting Amazon scraping via FastAPI\n", requestID)

	payload := map[string]string{"url": productURL}
	payloadJSON, _ := json.Marshal(payload)

	resp, err := http.Post("http://localhost:8081/scrape", "application/json", strings.NewReader(string(payloadJSON)))
	if err != nil {
		return nil, fmt.Errorf("FastAPI request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("FastAPI returned status %d, body: %s", resp.StatusCode, string(body))
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("FastAPI scraping failed: %s", result.Error)
	}

	product := map[string]string{
		"title":        result.Title,
		"mrp":          result.MRP,
		"discount":     result.Discount,
		"price":        result.Price,
		"rating":       result.Rating,
		"availability": result.Availability,
	}

	fmt.Printf("✅ [%s] Amazon FastAPI scraping completed\n", requestID)
	return product, nil
}
