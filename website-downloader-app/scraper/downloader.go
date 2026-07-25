package scraper

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"website-downloader-app/utils"
)

// SaveFile saves downloaded content to disk with proper directory structure
func SaveFile(workDir, fileURL string, content []byte, baseURL string) error {
	// Parse URL to get path
	// u, err := url.Parse(fileURL)
	// if err != nil {
	// 	return fmt.Errorf("invalid URL: %w", err)
	// }

	// Build local file path
	localPath, err := utils.BuildLocalPath(baseURL, fileURL)
	if err != nil {
		return fmt.Errorf("failed to build local path: %w", err)
	}

	// Full path including work directory
	fullPath := filepath.Join(workDir, localPath)

	// Create directory if needed
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(fullPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	utils.LogInfo("Saved: %s -> %s", fileURL, fullPath)
	return nil
}

// DownloadAsset downloads a single asset (image, CSS, JS, etc.)
func DownloadAsset(client *http.Client, assetURL, workDir, baseURL string) error {
	// Fetch the asset
	resp, err := client.Get(assetURL)
	if err != nil {
		return fmt.Errorf("failed to fetch asset: %w", err)
	}
	defer resp.Body.Close()

	// Check status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	// Read content
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read asset: %w", err)
	}

	// Save to disk
	return SaveFile(workDir, assetURL, content, baseURL)
}

// GetAssetType determines the type of asset from URL or content type
func GetAssetType(assetURL string) string {
	u, err := url.Parse(assetURL)
	if err != nil {
		return "unknown"
	}

	ext := strings.ToLower(path.Ext(u.Path))
	
	switch ext {
	case ".css":
		return "css"
	case ".js":
		return "javascript"
	case ".jpg", ".jpeg", ".png", ".gif", ".svg", ".webp", ".ico":
		return "image"
	case ".woff", ".woff2", ".ttf", ".eot", ".otf":
		return "font"
	case ".pdf":
		return "pdf"
	case ".json":
		return "json"
	case ".xml":
		return "xml"
	default:
		return "other"
	}
}

// ShouldDownloadAsset checks if an asset should be downloaded based on config
func ShouldDownloadAsset(assetURL string, allowedExtensions []string) bool {
	if len(allowedExtensions) == 0 {
		return true // Download everything if no restrictions
	}

	u, err := url.Parse(assetURL)
	if err != nil {
		return false
	}

	ext := strings.ToLower(path.Ext(u.Path))
	
	for _, allowed := range allowedExtensions {
		if ext == allowed {
			return true
		}
	}

	return false
}
