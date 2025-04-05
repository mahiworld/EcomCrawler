package crawler

import (
	"log"
	"sync"
)

// Worker processes URLs from the productChannel and saves them to the result map
func Worker(productChannel chan ProductURL, results map[string]map[string]struct{}, mu *sync.Mutex) {
	for product := range productChannel {
		mu.Lock()
		if _, exists := results[product.Domain]; !exists {
			results[product.Domain] = make(map[string]struct{})
		}
		results[product.Domain][product.URL] = struct{}{}
		mu.Unlock()

		log.Println("Discovered:", product.URL)
	}
}
