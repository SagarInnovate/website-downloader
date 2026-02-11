package utils

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

// NormalizeURL normalizes a URL by removing fragments and sorting query parameters
func NormalizeURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// Remove fragment
	u.Fragment = ""

	// Convert to lowercase scheme and host
	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)

	return u.String(), nil
}

// IsSameDomain checks if two URLs belong to the EXACT same domain (no subdomains)
func IsSameDomain(baseURL, targetURL string) bool {
	base, err := url.Parse(baseURL)
	if err != nil {
		return false
	}

	target, err := url.Parse(targetURL)
	if err != nil {
		return false
	}

	baseHost := strings.ToLower(base.Host)
	targetHost := strings.ToLower(target.Host)

	// Only exact match - no subdomains allowed
	return baseHost == targetHost
}

// ResolveURL resolves a potentially relative URL against a base URL
func ResolveURL(baseURL, href string) (string, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	ref, err := url.Parse(href)
	if err != nil {
		return "", err
	}

	resolved := base.ResolveReference(ref)
	return resolved.String(), nil
}

// GetRelativePath converts an absolute URL to a relative path suitable for local storage
func GetRelativePath(baseURL, targetURL string) (string, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return "", err
	}

	// Use URL path as relative path
	relativePath := strings.TrimPrefix(target.Path, "/")
	
	// If path is empty or ends with /, append index.html
	if relativePath == "" || strings.HasSuffix(relativePath, "/") {
		relativePath = path.Join(relativePath, "index.html")
	}

	return relativePath, nil
}

// IsValidURL checks if a URL is valid and uses HTTP/HTTPS
func IsValidURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	scheme := strings.ToLower(u.Scheme)
	return (scheme == "http" || scheme == "https") && u.Host != ""
}

// GetFileExtension returns the file extension from a URL path
func GetFileExtension(urlPath string) string {
	return strings.ToLower(path.Ext(urlPath))
}

// SanitizeFilename removes invalid characters from filenames
func SanitizeFilename(filename string) string {
	// Replace invalid characters with underscores
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := filename
	for _, char := range invalid {
		result = strings.ReplaceAll(result, char, "_")
	}
	return result
}

// GetDomain extracts the domain from a URL
func GetDomain(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	return strings.ToLower(u.Host), nil
}

// BuildLocalPath creates a local file path from a URL
func BuildLocalPath(baseURL, targetURL string) (string, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return "", fmt.Errorf("invalid target URL: %w", err)
	}

	// Start with the host
	localPath := target.Host

	// Add the path
	urlPath := target.Path
	if urlPath == "" || urlPath == "/" {
		urlPath = "/index.html"
	} else if strings.HasSuffix(urlPath, "/") {
		urlPath += "index.html"
	}

	localPath = path.Join(localPath, urlPath)

	return localPath, nil
}
