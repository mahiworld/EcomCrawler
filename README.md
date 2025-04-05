# EcomCrawler

A web crawler written in Go, designed to extract product URLs from popular e-commerce websites.

## Overview

This project crawls a list of e-commerce domains, discovers potential product page URLs, and writes the result to a results.json file. It uses [Colly](https://github.com/gocolly/colly) for web crawling and Go's built-in concurrency features for performance.

## How Product URLs Are Detected

The crawler checks for product-like URLs by looking for common patterns in hyperlinks, such as:

```
/product/, /products/, /p/, /item/, /shop/, /details/, /dp/
```

These are defined in constants/constants.go as:

```
var ProductIndicators = []string{
  "/product/", "/products/", "/p/", "/item/", "/shop/", "/details/", "/dp/",
}
```

When any link on the site contains one of these substrings, it's considered a potential product page.

## Features

- Concurrent crawling using Go routines and worker pools

- Duplicate product URL filtering per domain

- Domain list loaded dynamically from configs/domains.json

- Output saved as a structured results.json

## Project Structure 
```
EcomCrawler/
├── configs/
│   └── domains.json         # List of domains to crawl
├── constants/
│   └── constants.go         # Product URL patterns
├── crawler/
│   ├── crawler.go           # Crawler logic using Colly
│   └── worker.go            # Worker logic for processing URLs
├── results.json             # Output of crawled product URLs
└── main.go                  # Entry point, orchestration
```

## Usage

1. Configure Domains

Add the e-commerce domains to configs/domains.json:
```
[
  "https://www.virgio.com/",
  "https://www.tatacliq.com/",
  "https://nykaafashion.com/",
  "https://www.westside.com/"
]
```

2. Run the Crawler
```
go run main.go
```

It will:

- Load domains from domains.json

- Start concurrent crawlers

- Filter duplicate product URLs

- Save to results.json

  
## Sample Output (results.json)
```
{
  "https://www.virgio.com/": [
    "https://www.virgio.com/products/veronicas-chic-cotton-co-ords",
    "https://www.virgio.com/products/date-chic-red-solid-co-ords",
    "https://www.virgio.com/products/date-breezy-black-relaxed-trousers",
    "https://www.virgio.com/products/veronicas-playful-cotton-checked-skirt",
    "https://www.virgio.com/products/date-comfy-beige-relaxed-trousers",
    "https://www.virgio.com/products/veronicas-retro-cotton-pleated-skirt"
  ]
}
```

## Dependencies

github.com/gocolly/colly - Elegant scraping framework for Go

Install with:
```
go get github.com/gocolly/colly
```