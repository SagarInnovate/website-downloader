package scraper

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
	"website-downloader/models"
	"website-downloader/utils"

	"golang.org/x/net/html"
)

// Crawler manages the website crawling process
type Crawler struct {
	baseURL          string
	baseDomain       string
	config           *models.Config
	visited          map[string]bool
	urlQueue         []string
	downloadedPages  []string
	downloadedAssets []string
	totalSize        int64
	mutex            sync.RWMutex
	client           *http.Client
	workDir          string
}

// NewCrawler creates a new crawler instance
func NewCrawler(startURL string, config *models.Config, workDir string) (*Crawler, error) {
	// Validate URL
	if !utils.IsValidURL(startURL) {
		return nil, fmt.Errorf("invalid URL: %s", startURL)
	}

	// Normalize URL
	normalizedURL, err := utils.NormalizeURL(startURL)
	if err != nil {
		return nil, fmt.Errorf("failed to normalize URL: %w", err)
	}

	// Extract domain
	domain, err := utils.GetDomain(normalizedURL)
	if err != nil {
		return nil, fmt.Errorf("failed to extract domain: %w", err)
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: time.Duration(config.RequestTimeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allow up to 10 redirects
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	crawler := &Crawler{
		baseURL:          normalizedURL,
		baseDomain:       domain,
		config:           config,
		visited:          make(map[string]bool),
		urlQueue:         []string{normalizedURL},
		downloadedPages:  []string{},
		downloadedAssets: []string{},
		totalSize:        0,
		client:           client,
		workDir:          workDir,
	}

	return crawler, nil
}

// Crawl starts the crawling process
func (c *Crawler) Crawl(progressCallback func(update models.ProgressUpdate)) error {
	utils.LogInfo("Starting crawl for: %s", c.baseURL)

	depth := 0
	for len(c.urlQueue) > 0 && len(c.downloadedPages) < c.config.MaxPages {
		// Get next URL from queue
		currentURL := c.dequeueURL()
		if currentURL == "" {
			break
		}

		// Skip if already visited
		if c.isVisited(currentURL) {
			continue
		}

		// Mark as visited
		c.markVisited(currentURL)

		// Send progress update
		if progressCallback != nil {
			progressCallback(c.getProgressUpdate(currentURL))
		}

		// Fetch and process the page
		err := c.processPage(currentURL, depth)
		if err != nil {
			utils.LogError("Error processing %s: %v", currentURL, err)
			continue
		}

		utils.LogInfo("Processed: %s", currentURL)
	}

	utils.LogInfo("Crawl complete. Pages: %d, Assets: %d", len(c.downloadedPages), len(c.downloadedAssets))
	return nil
}

// processPage fetches and processes a single page
func (c *Crawler) processPage(pageURL string, depth int) error {
	// Fetch the page
	resp, err := c.client.Get(pageURL)
	if err != nil {
		return fmt.Errorf("failed to fetch: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	// Read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	// Update total size
	c.mutex.Lock()
	c.totalSize += int64(len(body))
	c.downloadedPages = append(c.downloadedPages, pageURL)
	c.mutex.Unlock()

	// Determine content type
	contentType := resp.Header.Get("Content-Type")

	// If it's HTML, parse and extract links
	if strings.Contains(contentType, "text/html") {
		links, assets, err := c.extractLinksAndAssets(pageURL, body)
		if err != nil {
			utils.LogError("Failed to extract links from %s: %v", pageURL, err)
		} else {
			// Add new pages to queue
			for _, link := range links {
				c.enqueueURL(link)
			}
			// Add assets to download queue
			for _, asset := range assets {
				c.enqueueAsset(asset)
			}
		}
	} else if strings.Contains(contentType, "text/css") || strings.HasSuffix(pageURL, ".css") {
		// Parse CSS and extract assets (fonts, images, etc.)
		cssAssets := c.extractCSSAssets(pageURL, body)
		for _, asset := range cssAssets {
			c.enqueueAsset(asset)
		}
	}

	// Save the file
	err = SaveFile(c.workDir, pageURL, body, c.baseURL)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

// extractLinksAndAssets parses HTML and extracts links and assets
func (c *Crawler) extractLinksAndAssets(baseURL string, htmlContent []byte) ([]string, []string, error) {
	doc, err := html.Parse(strings.NewReader(string(htmlContent)))
	if err != nil {
		return nil, nil, err
	}

	var links []string
	var assets []string

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "a":
				// Extract href from anchor tags
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						resolvedURL, err := utils.ResolveURL(baseURL, attr.Val)
						if err == nil && utils.IsSameDomain(c.baseURL, resolvedURL) {
							// Remove fragment
							if u, err := url.Parse(resolvedURL); err == nil {
								u.Fragment = ""
								links = append(links, u.String())
							}
						}
					}
				}
			case "link":
				// Extract href from link tags (CSS, etc.)
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						resolvedURL, err := utils.ResolveURL(baseURL, attr.Val)
						if err == nil && utils.IsSameDomain(c.baseURL, resolvedURL) {
							assets = append(assets, resolvedURL)
						}
					}
				}
			case "script":
				// Extract src from script tags
				for _, attr := range n.Attr {
					if attr.Key == "src" {
						resolvedURL, err := utils.ResolveURL(baseURL, attr.Val)
						if err == nil && utils.IsSameDomain(c.baseURL, resolvedURL) {
							assets = append(assets, resolvedURL)
						}
					}
				}
			case "img":
				// Extract src from img tags
				for _, attr := range n.Attr {
					if attr.Key == "src" {
						resolvedURL, err := utils.ResolveURL(baseURL, attr.Val)
						if err == nil && utils.IsSameDomain(c.baseURL, resolvedURL) {
							assets = append(assets, resolvedURL)
						}
					}
				}
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}
	}

	traverse(doc)

	return links, assets, nil
}

// Queue management methods
func (c *Crawler) enqueueURL(url string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Check if already visited or queued
	if c.visited[url] {
		return
	}

	for _, queuedURL := range c.urlQueue {
		if queuedURL == url {
			return
		}
	}

	c.urlQueue = append(c.urlQueue, url)
}

func (c *Crawler) enqueueAsset(assetURL string) {
	c.mutex.Lock()
	
	// Check if already downloaded
	if c.visited[assetURL] {
		c.mutex.Unlock()
		return
	}
	
	// Mark as visited
	c.visited[assetURL] = true
	c.mutex.Unlock()
	
	// Download asset in background goroutine
	go func() {
		utils.LogInfo("Downloading asset: %s", assetURL)
		
		// Fetch the asset
		resp, err := c.client.Get(assetURL)
		if err != nil {
			utils.LogError("Failed to fetch asset %s: %v", assetURL, err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			utils.LogError("Bad status %d for asset %s", resp.StatusCode, assetURL)
			return
		}

		// Read content
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.LogError("Failed to read asset %s: %v", assetURL, err)
			return
		}

		// Save the file
		err = SaveFile(c.workDir, assetURL, content, c.baseURL)
		if err != nil {
			utils.LogError("Failed to save asset %s: %v", assetURL, err)
			return
		}
		
		c.mutex.Lock()
		c.downloadedAssets = append(c.downloadedAssets, assetURL)
		c.totalSize += int64(len(content))
		c.mutex.Unlock()
		
		// If it's a CSS file, parse it for additional assets
		contentType := resp.Header.Get("Content-Type")
		if strings.Contains(contentType, "text/css") || strings.HasSuffix(assetURL, ".css") {
			utils.LogInfo("Parsing CSS file for nested assets: %s", assetURL)
			cssAssets := c.extractCSSAssets(assetURL, content)
			for _, nestedAsset := range cssAssets {
				c.enqueueAsset(nestedAsset)
			}
		}
		
		utils.LogInfo("Downloaded asset: %s", assetURL)
	}()
}

func (c *Crawler) dequeueURL() string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if len(c.urlQueue) == 0 {
		return ""
	}

	url := c.urlQueue[0]
	c.urlQueue = c.urlQueue[1:]
	return url
}

func (c *Crawler) isVisited(url string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.visited[url]
}

func (c *Crawler) markVisited(url string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.visited[url] = true
}

func (c *Crawler) getProgressUpdate(currentURL string) models.ProgressUpdate {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	pagesDownloaded := len(c.downloadedPages)
	pagesDiscovered := len(c.visited)
	progress := 0
	if c.config.MaxPages > 0 {
		progress = (pagesDownloaded * 100) / c.config.MaxPages
		if progress > 100 {
			progress = 100
		}
	}

	return models.ProgressUpdate{
		Type:             "progress",
		PagesDiscovered:  pagesDiscovered,
		PagesDownloaded:  pagesDownloaded,
		AssetsDownloaded: len(c.downloadedAssets),
		TotalSize:        c.totalSize,
		CurrentPage:      currentURL,
		Progress:         progress,
	}
}

// extractCSSAssets parses CSS content and extracts asset URLs (fonts, images, etc.)
func (c *Crawler) extractCSSAssets(cssURL string, cssContent []byte) []string {
	content := string(cssContent)
	var assets []string
	seenAssets := make(map[string]bool)

	// Regex to match url() declarations
	urlRegex := regexp.MustCompile(`url\s*\(\s*['"]?([^'")]+)['"]?\s*\)`)
	matches := urlRegex.FindAllStringSubmatch(content, -1)
	
	for _, match := range matches {
		if len(match) > 1 {
			assetURL := strings.TrimSpace(match[1])
			
			// Skip data URIs
			if strings.HasPrefix(assetURL, "data:") {
				continue
			}
			
			// Resolve URL
			resolvedURL, err := utils.ResolveURL(cssURL, assetURL)
			if err != nil || !utils.IsSameDomain(c.baseURL, resolvedURL) {
				continue
			}
			
			// Add to assets if not seen before
			if !seenAssets[resolvedURL] {
				assets = append(assets, resolvedURL)
				seenAssets[resolvedURL] = true
			}
		}
	}
	
	// Also extract @import statements
	importRegex := regexp.MustCompile(`@import\s+['"]([^'"]+)['"]`)
	importMatches := importRegex.FindAllStringSubmatch(content, -1)
	
	for _, match := range importMatches {
		if len(match) > 1 {
			importURL := strings.TrimSpace(match[1])
			
			// Resolve URL
			resolvedURL, err := utils.ResolveURL(cssURL, importURL)
			if err != nil || !utils.IsSameDomain(c.baseURL, resolvedURL) {
				continue
			}
			
			// Add to assets if not seen before
			if !seenAssets[resolvedURL] {
				assets = append(assets, resolvedURL)
				seenAssets[resolvedURL] = true
			}
		}
	}
	
	utils.LogInfo("Extracted %d assets from CSS: %s", len(assets), cssURL)
	return assets
}
