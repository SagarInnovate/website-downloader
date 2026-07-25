package scraper

import (
	"bytes"
	"path/filepath"
	"strings"
	"website-downloader-app/utils"

	"golang.org/x/net/html"
)

// RewriteLinks parses HTML and rewrites links to be relative to the local file system
// baseURL is the root URL of the site (e.g. https://example.com)
// pageURL is the URL of the current page being processed (e.g. https://example.com/about)
func RewriteLinks(htmlContent []byte, baseURL, pageURL string) ([]byte, error) {
	doc, err := html.Parse(strings.NewReader(string(htmlContent)))
	if err != nil {
		return nil, err
	}

	// Calculate local path for the current page
	pageLocalPath, err := utils.BuildLocalPath(baseURL, pageURL)
	if err != nil {
		return nil, err
	}
	pageDir := filepath.Dir(pageLocalPath)

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			var attrKey string
			switch n.Data {
			case "a", "link":
				attrKey = "href"
			case "script", "img", "iframe", "source", "embed", "track":
				attrKey = "src"
			}

			if attrKey != "" {
				for i, attr := range n.Attr {
					if attr.Key == attrKey {
						// Resolve raw URL to absolute URL
						resolvedURL, err := utils.ResolveURL(baseURL, attr.Val)
						if err != nil {
							continue
						}

						// Determine if we should rewrite this link
						shouldRewrite := false
						if attrKey == "href" && n.Data == "a" {
							// For anchors, only rewrite if same domain
							if utils.IsSameDomain(baseURL, resolvedURL) {
								shouldRewrite = true
							}
						} else {
							// For assets, we download everything now, so rewrite everything
							shouldRewrite = true
						}

						if shouldRewrite {
							// Calculate where this file was/will be saved
							targetLocalPath, err := utils.BuildLocalPath(baseURL, resolvedURL)
							if err == nil {
								// Calculate relative path from pageDir to targetLocalPath
								relPath, err := filepath.Rel(pageDir, targetLocalPath)
								if err == nil {
									// Normalize path for URL (use forward slashes)
									relPath = filepath.ToSlash(relPath)
									
									// Ensure it looks like a relative path if in same dir
									if !strings.HasPrefix(relPath, "..") && !strings.HasPrefix(relPath, "/") {
										relPath = "./" + relPath
									}

									n.Attr[i].Val = relPath
								}
							}
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

	// Render back to string
	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
