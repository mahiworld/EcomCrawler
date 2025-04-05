package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sync"

	"EcomCrawler/crawler"
)

func loadDomainsFromFile(path string) []string {
	var domains []string
	data, err := os.ReadFile(path)
	if err != nil {
		panic("Unable to read domain file: " + err.Error())
	}
	err = json.Unmarshal(data, &domains)
	if err != nil {
		panic("Invalid JSON in domain file: " + err.Error())
	}
	return domains
}

func main() {
	domains := loadDomainsFromFile("configs/domains.json")

	rawResults := make(map[string]map[string]struct{})

	// Use CPU count for optimal workers
	numWorkers := runtime.NumCPU() * 2
	fmt.Printf("Starting %d workers...\n", numWorkers)

	// Buffered channel to handle incoming product URLs
	bufferSize := len(domains) * 1000
	productChannel := make(chan crawler.ProductURL, bufferSize)

	// Waits for all worker goroutines to finish.
	var wg sync.WaitGroup
	// Ensures safe concurrent writes to the results map.
	var mu sync.Mutex

	// Start worker pool
	wg.Add(numWorkers)
	for range numWorkers {
		go func() {
			defer wg.Done()
			crawler.Worker(productChannel, rawResults, &mu)
		}()
	}

	// Start crawling each domain
	var crawlWg sync.WaitGroup
	crawlWg.Add(len(domains))
	for _, domain := range domains {
		go func(d string) {
			defer crawlWg.Done()
			crawler.Crawl(d, productChannel)
		}(domain)
	}

	// Wait for crawling to finish
	crawlWg.Wait()
	// Close the product channel to signal workers to finish
	// This is done after all crawlers are done to ensure no data is lost.
	close(productChannel)

	// Wait for all workers to finish
	wg.Wait()

	// Convert map[string]map[string]struct{} to map[string][]string
	finalResults := make(map[string][]string)
	for domain, urls := range rawResults {
		for url := range urls {
			finalResults[domain] = append(finalResults[domain], url)
		}
	}

	// Save to file
	file, _ := json.MarshalIndent(finalResults, "", "  ")
	_ = os.WriteFile("./results.json", file, 0644)

	fmt.Println("Crawling completed. Results saved to results.json")
	fmt.Println("Total products crawled:", len(finalResults))
}
