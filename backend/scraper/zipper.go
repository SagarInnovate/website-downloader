package scraper

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"website-downloader/utils"
)

// CreateZIP creates a ZIP archive from the downloaded website
func CreateZIP(sourceDir, outputPath string) error {
	utils.LogInfo("Creating ZIP archive: %s", outputPath)

	// Create ZIP file
	zipFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create ZIP file: %w", err)
	}
	defer zipFile.Close()

	// Create ZIP writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk through source directory
	err = filepath.Walk(sourceDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Get relative path for ZIP entry
		relPath, err := filepath.Rel(sourceDir, filePath)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Convert to forward slashes for ZIP (platform independent)
		relPath = strings.ReplaceAll(relPath, "\\", "/")

		// Create ZIP entry
		zipEntry, err := zipWriter.Create(relPath)
		if err != nil {
			return fmt.Errorf("failed to create ZIP entry: %w", err)
		}

		// Open source file
		sourceFile, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer sourceFile.Close()

		// Copy file content to ZIP
		_, err = io.Copy(zipEntry, sourceFile)
		if err != nil {
			return fmt.Errorf("failed to copy file to ZIP: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}

	utils.LogInfo("ZIP archive created successfully: %s", outputPath)
	return nil
}

// GetZIPSize returns the size of a ZIP file in bytes
func GetZIPSize(zipPath string) (int64, error) {
	info, err := os.Stat(zipPath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}
