package scraper

import (
	"context"
	"fmt"
	"time"
	"website-downloader/models"
	"website-downloader/utils"

	"github.com/chromedp/chromedp"
)

// BrowserCrawler manages browser-based crawling for SPAs
type BrowserCrawler struct {
	baseURL string
	config  *models.Config
}

// NewBrowserCrawler creates a new browser crawler
func NewBrowserCrawler(baseURL string, config *models.Config) (*BrowserCrawler, error) {
	return &BrowserCrawler{
		baseURL: baseURL,
		config:  config,
	}, nil
}

// FetchPage fetches a page using headless browser and returns rendered HTML
func (bc *BrowserCrawler) FetchPage(url string) ([]byte, error) {
	utils.LogInfo("Browser mode: Fetching %s", url)

	// Create Chrome context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set timeout
	ctx, cancel = context.WithTimeout(ctx, time.Duration(bc.config.BrowserTimeout)*time.Second)
	defer cancel()

	var htmlContent string

	// Navigate and wait for page to render
	err := chromedp.Run(ctx,
		// Navigate to URL
		chromedp.Navigate(url),

		// Wait for body to be visible (ensure basic rendering)
		chromedp.WaitVisible("body", chromedp.ByQuery),
		
		// Wait for network idle or timeout (extra wait for JS execution)
		chromedp.Sleep(time.Duration(bc.config.BrowserWaitTime + 2000)*time.Millisecond), // Add extra 2s buffer
		
		// Get the rendered HTML
		chromedp.OuterHTML("html", &htmlContent),
	)

	if err != nil {
		utils.LogError("Browser fetch failed for %s: %v", url, err)
		return nil, fmt.Errorf("browser fetch failed: %w", err)
	}

	utils.LogInfo("Browser mode: Successfully rendered %s (%d bytes). Mode: SPA/Browser", url, len(htmlContent))
	return []byte(htmlContent), nil
}
