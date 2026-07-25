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
	"website-downloader-app/models"
	"website-downloader-app/utils"

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
	mode             string // "static" or "browser"
}

// NewCrawler creates a new crawler instance
func NewCrawler(startURL string, config *models.Config, workDir string, mode string) (*Crawler, error) {
	// Default to static mode if not specified
	if mode == "" {
		mode = "static"
	}
	
	// Validate mode
	if mode != "static" && mode != "browser" {
		return nil, fmt.Errorf("invalid mode: %s (must be 'static' or 'browser')", mode)
	}

	utils.LogInfo("Initializing crawler in %s mode", mode)

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
		mode:             mode,
	}

	return crawler, nil
}

// Crawl starts the crawling process
func (c *Crawler) Crawl(progressCallback func(update models.ProgressUpdate)) error {
	utils.LogInfo("Starting crawl for: %s (mode: %s)", c.baseURL, c.mode)

	// If browser mode, use SPA-aware crawling
	if c.mode == "browser" {
		return c.crawlSPA(progressCallback)
	}

	// Otherwise, use traditional crawling for static sites
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

// crawlSPA handles SPA-specific crawling with route discovery
func (c *Crawler) crawlSPA(progressCallback func(update models.ProgressUpdate)) error {
	utils.LogInfo("Starting SPA crawl for: %s", c.baseURL)

	// Create browser crawler (persistent for the entire session)
	browserCrawler, err := NewBrowserCrawler(c.baseURL, c.config)
	if err != nil {
		return fmt.Errorf("failed to initialize browser: %w", err)
	}
	defer browserCrawler.Close()

	// Step 1: Discover all routes in the SPA
	utils.LogInfo("Step 1: Discovering SPA routes...")
	routes, err := browserCrawler.DiscoverSPARoutes()
	if err != nil {
		utils.LogError("Failed to discover routes: %v", err)
		// Continue with just the base URL
		routes = []string{"/"}
	}

	// Add base URL first
	allURLs := []string{c.baseURL}
	
	// Add discovered routes
	for _, route := range routes {
		if len(allURLs) >= c.config.MaxPages {
			break
		}
		// Convert route to full URL
		fullURL := strings.TrimSuffix(c.baseURL, "/") + route
		allURLs = append(allURLs, fullURL)
	}

	utils.LogInfo("Step 2: Crawling %d routes...", len(allURLs))

	// Step 2: Process each route sequentially
	for i, pageURL := range allURLs {
		if len(c.downloadedPages) >= c.config.MaxPages {
			break
		}

		// Skip if already visited
		if c.isVisited(pageURL) {
			continue
		}

		// Mark as visited
		c.markVisited(pageURL)

		// Send progress update
		if progressCallback != nil {
			progress := (i * 100) / len(allURLs)
			update := c.getProgressUpdate(pageURL)
			update.Progress = progress
			progressCallback(update)
		}

		utils.LogInfo("Processing SPA route %d/%d: %s", i+1, len(allURLs), pageURL)

		// Fetch and process the page using browser
		err := c.processSPAPage(browserCrawler, pageURL)
		if err != nil {
			utils.LogError("Error processing SPA route %s: %v", pageURL, err)
			continue
		}

		utils.LogInfo("Completed SPA route: %s", pageURL)

		// Small delay between routes to avoid overwhelming the browser
		time.Sleep(500 * time.Millisecond)
	}

	utils.LogInfo("SPA crawl complete. Pages: %d, Assets: %d", len(c.downloadedPages), len(c.downloadedAssets))
	return nil
}

// processSPAPage processes a single page in SPA mode
func (c *Crawler) processSPAPage(browserCrawler *BrowserCrawler, pageURL string) error {
	var body []byte
	var err error

	// Extract route from URL for navigation
	route := strings.TrimPrefix(pageURL, strings.TrimSuffix(c.baseURL, "/"))
	if route == "" {
		route = "/"
	}

	// If this is the first page (root), use FetchPage
	if route == "/" || pageURL == c.baseURL {
		body, err = browserCrawler.FetchPage(pageURL)
	} else {
		// For other routes, use NavigateAndCapture (more efficient for SPAs)
		body, err = browserCrawler.NavigateAndCapture(route)
	}

	if err != nil {
		return fmt.Errorf("browser fetch failed: %w", err)
	}

	// Update total size
	c.mutex.Lock()
	c.totalSize += int64(len(body))
	c.downloadedPages = append(c.downloadedPages, pageURL)
	c.mutex.Unlock()

	// Parse and extract assets from the rendered HTML
	_, assets, err := c.extractLinksAndAssets(pageURL, body)
	if err != nil {
		utils.LogError("Failed to extract assets from %s: %v", pageURL, err)
	} else {
		// Add assets to download queue
		for _, asset := range assets {
			c.enqueueAsset(asset)
		}
		utils.LogInfo("Extracted %d assets from %s", len(assets), pageURL)
	}

	// Rewrite links to be relative for offline viewing
	rewrittenBody, err := RewriteLinks(body, c.baseURL, pageURL)
	if err != nil {
		utils.LogError("Failed to rewrite links for %s: %v", pageURL, err)
		// Continue with original body if rewriting fails
	} else {
		body = rewrittenBody
		utils.LogInfo("Rewrote links for %s", pageURL)
	}

	// Save the file
	err = SaveFile(c.workDir, pageURL, body, c.baseURL)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

// processPage fetches and processes a single page
func (c *Crawler) processPage(pageURL string, depth int) error {
	var body []byte
	var err error

	// Fetch using appropriate mode
	if c.mode == "browser" {
		// Use headless browser for SPAs
		browserCrawler, err := NewBrowserCrawler(c.baseURL, c.config)
		if err != nil {
			return fmt.Errorf("failed to initialize browser: %w", err)
		}
		
		body, err = browserCrawler.FetchPage(pageURL)
		if err != nil {
			return fmt.Errorf("browser fetch failed: %w", err)
		}
	} else {
		// Use HTTP client for static sites (fast mode)
		resp, err := c.client.Get(pageURL)
		if err != nil {
			return fmt.Errorf("failed to fetch: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("bad status: %d", resp.StatusCode)
		}

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read body: %w", err)
		}
	}

	// Update total size
	c.mutex.Lock()
	c.totalSize += int64(len(body))
	c.downloadedPages = append(c.downloadedPages, pageURL)
	c.mutex.Unlock()

	// Determine if content is HTML (check for HTML tags in content)
	isHTML := strings.Contains(string(body), "<html") || strings.Contains(string(body), "<!DOCTYPE")
	
	// If it's HTML, parse and extract links
	if isHTML {
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

		// Rewrite links to be relative for offline viewing
		rewrittenBody, err := RewriteLinks(body, c.baseURL, pageURL)
		if err != nil {
			utils.LogError("Failed to rewrite links for %s: %v", pageURL, err)
			// Continue with original body if rewriting fails
		} else {
			body = rewrittenBody
			utils.LogInfo("Rewrote links for %s", pageURL)
		}
	} else if strings.HasSuffix(pageURL, ".css") {
		// Parse CSS and extract assets (fonts, images, etc.)
		cssAssets := c.extractCSSAssets(pageURL, body)
		for _, asset := range cssAssets {
			c.enqueueAsset(asset)
		}
		
		// TODO: Rewrite CSS URLs too? (For now, relative paths in CSS usually work if assets are downloaded relative to CSS)
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
						if err == nil {
							// Allow assets from any domain
							assets = append(assets, resolvedURL)
						}
					}
				}
			case "script":
				// Extract src from script tags
				for _, attr := range n.Attr {
					if attr.Key == "src" {
						resolvedURL, err := utils.ResolveURL(baseURL, attr.Val)
						if err == nil {
							// Allow assets from any domain
							assets = append(assets, resolvedURL)
						}
					}
				}
			case "img":
				// Extract src from img tags
				for _, attr := range n.Attr {
					if attr.Key == "src" {
						resolvedURL, err := utils.ResolveURL(baseURL, attr.Val)
						if err == nil {
							// Allow assets from any domain
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
			if err != nil {
				continue
			}
			
			// Allow assets from any domain
			
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
