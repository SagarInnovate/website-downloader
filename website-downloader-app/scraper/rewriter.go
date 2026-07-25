package scraper

import (
	"bytes"
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"
	"website-downloader-app/utils"

	"golang.org/x/net/html"
)

// RewriteHTML rewrites all URLs in HTML to use relative paths
func RewriteHTML(htmlContent []byte, pageURL, baseURL string) ([]byte, error) {
	doc, err := html.Parse(bytes.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Rewrite URLs in the HTML tree
	var rewrite func(*html.Node)
	rewrite = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for i, attr := range n.Attr {
				shouldRewrite := false
				switch n.Data {
				case "a", "link":
					shouldRewrite = attr.Key == "href"
				case "script", "img", "iframe":
					shouldRewrite = attr.Key == "src"
				}

				if shouldRewrite {
					// Resolve the URL
					resolvedURL, err := utils.ResolveURL(pageURL, attr.Val)
					if err == nil && utils.IsSameDomain(baseURL, resolvedURL) {
						// Convert to relative path
						relativePath, err := getRelativePathBetweenPages(pageURL, resolvedURL, baseURL)
						if err == nil {
							n.Attr[i].Val = relativePath
						}
					}
				}
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			rewrite(child)
		}
	}

	rewrite(doc)

	// Render back to HTML
	var buf bytes.Buffer
	err = html.Render(&buf, doc)
	if err != nil {
		return nil, fmt.Errorf("failed to render HTML: %w", err)
	}

	return buf.Bytes(), nil
}

// RewriteCSS rewrites URLs in CSS files (url() and @import)
func RewriteCSS(cssContent []byte, cssURL, baseURL string) ([]byte, error) {
	content := string(cssContent)

	// Regex to match url() declarations
	urlRegex := regexp.MustCompile(`url\(['"]?([^'")]+)['"]?\)`)
	content = urlRegex.ReplaceAllStringFunc(content, func(match string) string {
		// Extract the URL from url(...)
		urlMatch := regexp.MustCompile(`url\(['"]?([^'")]+)['"]?\)`).FindStringSubmatch(match)
		if len(urlMatch) < 2 {
			return match
		}

		originalURL := urlMatch[1]
		
		// Resolve URL
		resolvedURL, err := utils.ResolveURL(cssURL, originalURL)
		if err != nil || !utils.IsSameDomain(baseURL, resolvedURL) {
			return match
		}

		// Convert to relative path
		relativePath, err := getRelativePathBetweenPages(cssURL, resolvedURL, baseURL)
		if err != nil {
			return match
		}

		return fmt.Sprintf("url('%s')", relativePath)
	})

	// Regex to match @import statements
	importRegex := regexp.MustCompile(`@import\s+['"]([^'"]+)['"]`)
	content = importRegex.ReplaceAllStringFunc(content, func(match string) string {
		importMatch := regexp.MustCompile(`@import\s+['"]([^'"]+)['"]`).FindStringSubmatch(match)
		if len(importMatch) < 2 {
			return match
		}

		originalURL := importMatch[1]
		
		// Resolve URL
		resolvedURL, err := utils.ResolveURL(cssURL, originalURL)
		if err != nil || !utils.IsSameDomain(baseURL, resolvedURL) {
			return match
		}

		// Convert to relative path
		relativePath, err := getRelativePathBetweenPages(cssURL, resolvedURL, baseURL)
		if err != nil {
			return match
		}

		return fmt.Sprintf("@import '%s'", relativePath)
	})

	return []byte(content), nil
}

// getRelativePathBetweenPages calculates the relative path from one page to another
func getRelativePathBetweenPages(fromURL, toURL, baseURL string) (string, error) {
	// Parse URLs
	from, err := url.Parse(fromURL)
	if err != nil {
		return "", err
	}
	to, err := url.Parse(toURL)
	if err != nil {
		return "", err
	}

	// Get paths
	fromPath := from.Path
	toPath := to.Path

	// Normalize paths
	if fromPath == "" || fromPath == "/" {
		fromPath = "/index.html"
	} else if strings.HasSuffix(fromPath, "/") {
		fromPath += "index.html"
	}

	if toPath == "" || toPath == "/" {
		toPath = "/index.html"
	} else if strings.HasSuffix(toPath, "/") {
		toPath += "index.html"
	}

	// Calculate relative path
	fromDir := path.Dir(fromPath)
	relativePath, err := relPath(fromDir, toPath)
	if err != nil {
		// Fallback to absolute path from root
		return strings.TrimPrefix(toPath, "/"), nil
	}

	// Convert backslashes to forward slashes (Windows compatibility)
	relativePath = strings.ReplaceAll(relativePath, "\\", "/")

	return relativePath, nil
}

// relPath calculates relative path for URL paths (since filepath.Rel is OS-dependent)
func relPath(basepath, targpath string) (string, error) {
	baseSlash := strings.HasSuffix(basepath, "/")
	base := path.Clean(basepath)
	targ := path.Clean(targpath)
	
	if base == targ {
		return ".", nil
	}

	// Split paths
	baseparts := strings.Split(base, "/")
	targparts := strings.Split(targ, "/")

	// Find common prefix
	commonLen := 0
	for i := 0; i < len(baseparts) && i < len(targparts); i++ {
		if baseparts[i] != targparts[i] {
			break
		}
		commonLen++
	}

	// Build relative path
	var relparts []string
	
	// Add ".." for each directory in base after common prefix
	for i := commonLen; i < len(baseparts); i++ {
		relparts = append(relparts, "..")
	}

	// Add remaining target path parts
	relparts = append(relparts, targparts[commonLen:]...)

	if len(relparts) == 0 {
		return ".", nil
	}

	result := strings.Join(relparts, "/")
	if baseSlash && !strings.HasSuffix(result, "/") {
		result += "/"
	}

	return result, nil
}
