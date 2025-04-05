package crawler

import (
	"log"
	"strings"

	"EcomCrawler/constants" // import your constants package

	"github.com/gocolly/colly"
)

// ProductURL URL structure
type ProductURL struct {
	Domain string
	URL    string
}

// Crawl starts crawling a given domain and sends product URLs to the channel
func Crawl(domain string, productChannel chan ProductURL) {
	// Create a new collector for each domain
	c := colly.NewCollector(
		colly.AllowedDomains(getDomain(domain)),
		colly.Async(true),
	)

	// Use the shared constants
	productIndicators := constants.ProductIndicators

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))

		if isProductURL(link, productIndicators) {
			productChannel <- ProductURL{
				Domain: domain,
				URL:    link,
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	c.Visit(domain)
	c.Wait()
}

// getDomain extracts the base domain for AllowedDomains
func getDomain(url string) string {
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	return strings.TrimSuffix(url, "/")
}

// isProductURL checks if the URL looks like a product page
func isProductURL(url string, indicators []string) bool {
	for _, keyword := range indicators {
		if strings.Contains(url, keyword) {
			return true
		}
	}
	return false
}
