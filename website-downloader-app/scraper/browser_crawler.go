package scraper

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"
	"website-downloader-app/models"
	"website-downloader-app/utils"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// BrowserCrawler manages browser-based crawling for SPAs
type BrowserCrawler struct {
	baseURL string
	config  *models.Config
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewBrowserCrawler creates a new browser crawler
func NewBrowserCrawler(baseURL string, config *models.Config) (*BrowserCrawler, error) {
	// Create Chrome context with options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.UserAgent(config.UserAgent),
	)
	
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocCtx)

	return &BrowserCrawler{
		baseURL: baseURL,
		config:  config,
		ctx:     ctx,
		cancel:  cancel,
	}, nil
}

// Close closes the browser context
func (bc *BrowserCrawler) Close() {
	if bc.cancel != nil {
		bc.cancel()
	}
}

// FetchPage fetches a page using headless browser and returns rendered HTML
func (bc *BrowserCrawler) FetchPage(url string) ([]byte, error) {
	utils.LogInfo("Browser mode: Fetching %s", url)

	// Set timeout for this operation
	ctx, cancel := context.WithTimeout(bc.ctx, time.Duration(bc.config.BrowserTimeout)*time.Second)
	defer cancel()

	var htmlContent string

	// Navigate and wait for page to render with network idle
	err := chromedp.Run(ctx,
		// Enable network tracking
		network.Enable(),
		
		// Navigate to URL
		chromedp.Navigate(url),

		// Wait for body to be visible (ensure basic rendering)
		chromedp.WaitVisible("body", chromedp.ByQuery),
		
		// Wait for network to be idle (no requests for 500ms)
		chromedp.ActionFunc(func(ctx context.Context) error {
			return waitForNetworkIdle(ctx, 500*time.Millisecond, time.Duration(bc.config.BrowserWaitTime)*time.Millisecond)
		}),
		
		// Additional wait for any final rendering
		chromedp.Sleep(1*time.Second),
		
		// Get the rendered HTML
		chromedp.OuterHTML("html", &htmlContent),
	)

	if err != nil {
		utils.LogError("Browser fetch failed for %s: %v", url, err)
		return nil, fmt.Errorf("browser fetch failed: %w", err)
	}

	utils.LogInfo("Browser mode: Successfully rendered %s (%d bytes)", url, len(htmlContent))
	return []byte(htmlContent), nil
}

// DiscoverSPARoutes discovers routes in a SPA by analyzing the JavaScript bundles and navigation
func (bc *BrowserCrawler) DiscoverSPARoutes() ([]string, error) {
	utils.LogInfo("Browser mode: Discovering SPA routes from %s", bc.baseURL)

	ctx, cancel := context.WithTimeout(bc.ctx, time.Duration(bc.config.BrowserTimeout)*time.Second)
	defer cancel()

	var routes []string
	var jsContent string

	// First, load the main page
	err := chromedp.Run(ctx,
		network.Enable(),
		chromedp.Navigate(bc.baseURL),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			return waitForNetworkIdle(ctx, 500*time.Millisecond, 10*time.Second)
		}),
		chromedp.Sleep(2*time.Second),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to load initial page: %w", err)
	}

	// Try to extract routes from common navigation patterns
	var navLinks []string
	err = chromedp.Run(ctx,
		// Extract href attributes from navigation links
		chromedp.Evaluate(`
			Array.from(document.querySelectorAll('a[href], nav a, [role="navigation"] a, .nav a, .menu a'))
				.map(a => a.getAttribute('href'))
				.filter(href => href && !href.startsWith('http') && !href.startsWith('#') && href !== '/')
		`, &navLinks),
	)

	if err == nil && len(navLinks) > 0 {
		utils.LogInfo("Found %d navigation links", len(navLinks))
		routes = append(routes, navLinks...)
	}

	// Try to extract routes from JavaScript (common SPA router patterns)
	err = chromedp.Run(ctx,
		chromedp.Evaluate(`
			(() => {
				// Try to get all script tags content
				const scripts = Array.from(document.scripts);
				return scripts.map(s => s.textContent || s.src).join('\n');
			})()
		`, &jsContent),
	)

	if err == nil && jsContent != "" {
		extractedRoutes := extractRoutesFromJS(jsContent)
		utils.LogInfo("Extracted %d routes from JavaScript", len(extractedRoutes))
		routes = append(routes, extractedRoutes...)
	}

	// Deduplicate and clean routes
	routes = deduplicateAndCleanRoutes(routes, bc.baseURL)

	utils.LogInfo("Browser mode: Discovered %d unique routes", len(routes))
	return routes, nil
}

// NavigateAndCapture navigates to a route within the SPA and captures the rendered content
func (bc *BrowserCrawler) NavigateAndCapture(route string) ([]byte, error) {
	utils.LogInfo("Browser mode: Navigating to route %s", route)

	ctx, cancel := context.WithTimeout(bc.ctx, time.Duration(bc.config.BrowserTimeout)*time.Second)
	defer cancel()

	var htmlContent string

	// For SPAs, we often need to click on links or use history.pushState
	fullURL := bc.baseURL
	if !strings.HasPrefix(route, "http") {
		// Normalize route
		if !strings.HasPrefix(route, "/") {
			route = "/" + route
		}
		fullURL = strings.TrimSuffix(bc.baseURL, "/") + route
	}

	err := chromedp.Run(ctx,
		// Navigate using history API (more SPA-friendly)
		chromedp.Evaluate(fmt.Sprintf(`window.history.pushState({}, '', '%s')`, route), nil),
		
		// Trigger popstate event to notify SPA router
		chromedp.Evaluate(`window.dispatchEvent(new PopStateEvent('popstate'))`, nil),
		
		// Wait for content to update
		chromedp.Sleep(500*time.Millisecond),
		
		// Wait for network idle
		chromedp.ActionFunc(func(ctx context.Context) error {
			return waitForNetworkIdle(ctx, 500*time.Millisecond, time.Duration(bc.config.BrowserWaitTime)*time.Millisecond)
		}),
		
		// Additional wait for rendering
		chromedp.Sleep(1*time.Second),
		
		// Get the rendered HTML
		chromedp.OuterHTML("html", &htmlContent),
	)

	if err != nil {
		// If navigation via history fails, try direct navigation
		utils.LogInfo("History navigation failed, trying direct navigation to %s", fullURL)
		err = chromedp.Run(ctx,
			chromedp.Navigate(fullURL),
			chromedp.WaitVisible("body", chromedp.ByQuery),
			chromedp.ActionFunc(func(ctx context.Context) error {
				return waitForNetworkIdle(ctx, 500*time.Millisecond, time.Duration(bc.config.BrowserWaitTime)*time.Millisecond)
			}),
			chromedp.Sleep(1*time.Second),
			chromedp.OuterHTML("html", &htmlContent),
		)
		
		if err != nil {
			return nil, fmt.Errorf("failed to navigate to route: %w", err)
		}
	}

	utils.LogInfo("Browser mode: Captured route %s (%d bytes)", route, len(htmlContent))
	return []byte(htmlContent), nil
}

// waitForNetworkIdle waits until there are no network requests for a specified duration
func waitForNetworkIdle(ctx context.Context, idleDuration, maxWait time.Duration) error {
	start := time.Now()
	lastActivity := time.Now()
	activeRequests := 0

	// Listen to network events
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev.(type) {
		case *network.EventRequestWillBeSent:
			activeRequests++
			lastActivity = time.Now()
		case *network.EventLoadingFinished, *network.EventLoadingFailed:
			if activeRequests > 0 {
				activeRequests--
			}
			lastActivity = time.Now()
		}
	})

	// Wait for idle or timeout
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			// Check if we've been idle long enough
			if time.Since(lastActivity) >= idleDuration && activeRequests == 0 {
				return nil
			}
			
			// Check if we've exceeded max wait time
			if time.Since(start) >= maxWait {
				utils.LogInfo("Network idle timeout reached after %v", time.Since(start))
				return nil
			}
		}
	}
}

// extractRoutesFromJS extracts route patterns from JavaScript code
func extractRoutesFromJS(jsContent string) []string {
	var routes []string
	seen := make(map[string]bool)

	// Common route patterns in different frameworks
	patterns := []string{
		// React Router: <Route path="/about" />
		`path\s*[=:]\s*["']([^"']+)["']`,
		// Vue Router: { path: '/about' }
		`path\s*:\s*["']([^"']+)["']`,
		// Angular: { path: 'about' }
		`path\s*:\s*["']([^"']+)["']`,
		// Direct strings that look like routes
		`["'](/[\w\-/]+)["']`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(jsContent, -1)
		
		for _, match := range matches {
			if len(match) > 1 {
				route := match[1]
				// Filter out non-route patterns
				if isValidRoute(route) && !seen[route] {
					routes = append(routes, route)
					seen[route] = true
				}
			}
		}
	}

	return routes
}

// isValidRoute checks if a string looks like a valid route
func isValidRoute(route string) bool {
	// Must start with /
	if !strings.HasPrefix(route, "/") {
		return false
	}
	
	// Exclude common non-routes
	excludePatterns := []string{
		"//", "http", "https", ".js", ".css", ".png", ".jpg", ".svg",
		"api/", "cdn", "static/", "assets/", "*.",":",
	}
	
	for _, pattern := range excludePatterns {
		if strings.Contains(route, pattern) {
			return false
		}
	}
	
	// Exclude dynamic route parameters (we'll handle these separately)
	if strings.Contains(route, ":") || strings.Contains(route, "*") {
		return false
	}
	
	return true
}

// deduplicateAndCleanRoutes removes duplicates and normalizes routes
func deduplicateAndCleanRoutes(routes []string, baseURL string) []string {
	seen := make(map[string]bool)
	var cleaned []string
	
	for _, route := range routes {
		// Normalize route
		route = strings.TrimSpace(route)
		
		// Remove query strings and fragments
		if idx := strings.Index(route, "?"); idx != -1 {
			route = route[:idx]
		}
		if idx := strings.Index(route, "#"); idx != -1 {
			route = route[:idx]
		}
		
		// Ensure it starts with /
		if !strings.HasPrefix(route, "/") {
			route = "/" + route
		}
		
		// Remove trailing slash (except for root)
		if len(route) > 1 && strings.HasSuffix(route, "/") {
			route = strings.TrimSuffix(route, "/")
		}
		
		// Skip if already seen or is root
		if seen[route] || route == "/" {
			continue
		}
		
		seen[route] = true
		cleaned = append(cleaned, route)
	}
	
	return cleaned
}
