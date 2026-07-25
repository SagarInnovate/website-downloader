package scraper

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"website-downloader-app/models"
	"website-downloader-app/utils"
)

// Job represents a scraping job
type Job struct {
	ID        string
	URL       string
	Config    *models.Config
	Status    *models.JobStatus
	WorkDir   string
	OutputDir string
}

// StatusUpdateFunc is a function type for updating job status
type StatusUpdateFunc func(jobID string, status *models.JobStatus)

// ProgressBroadcastFunc is a function type for broadcasting progress
type ProgressBroadcastFunc func(jobID string, update models.ProgressUpdate)

// StartJob starts a new scraping job
func StartJob(
	jobID, url, mode string,
	config *models.Config,
	updateStatus StatusUpdateFunc,
	broadcastProgress ProgressBroadcastFunc,
) {
	utils.LogInfo("Starting job %s for URL: %s (mode: %s)", jobID, url, mode)

	// Create work directory
	workDir := filepath.Join(config.TempDir, jobID)
	if err := os.MkdirAll(workDir, 0755); err != nil {
		updateJobError(jobID, fmt.Sprintf("Failed to create work directory: %v", err), updateStatus, broadcastProgress)
		return
	}

	// Create output directory
	outputDir := config.OutputDir
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		updateJobError(jobID, fmt.Sprintf("Failed to create output directory: %v", err), updateStatus, broadcastProgress)
		return
	}

	// Initialize job status
	status := &models.JobStatus{
		JobID:     jobID,
		Status:    "running",
		URL:       url,
		StartTime: time.Now(),
	}
	updateStatus(jobID, status)

	// Send initial progress
	sendProgressUpdate(jobID, "progress", status, "Starting download...", broadcastProgress)

	// Create crawler with mode
	crawler, err := NewCrawler(url, config, workDir, mode)
	if err != nil {
		updateJobError(jobID, fmt.Sprintf("Failed to create crawler: %v", err), updateStatus, broadcastProgress)
		return
	}

	// Progress callback
	progressCallback := func(update models.ProgressUpdate) {
		// Update status
		status.PagesDiscovered = update.PagesDiscovered
		status.PagesDownloaded = update.PagesDownloaded
		status.AssetsDownloaded = update.AssetsDownloaded
		status.TotalSize = update.TotalSize
		status.CurrentPage = update.CurrentPage
		status.Progress = update.Progress

		updateStatus(jobID, status)
		broadcastProgress(jobID, update)
	}

	// Start crawling
	err = crawler.Crawl(progressCallback)
	if err != nil {
		updateJobError(jobID, fmt.Sprintf("Crawling failed: %v", err), updateStatus, broadcastProgress)
		return
	}

	// Rewrite all HTML and CSS files
	utils.LogInfo("Rewriting links...")
	status.CurrentPage = "Rewriting links for offline use..."
	status.Progress = 90
	updateStatus(jobID, status)
	sendProgressUpdate(jobID, "progress", status, "Rewriting links...", broadcastProgress)

	err = rewriteAllFiles(workDir, url)
	if err != nil {
		utils.LogError("Failed to rewrite files: %v", err)
		// Don't fail the job, just log the error
	}

	// Create ZIP archive with proper naming
	utils.LogInfo("Creating ZIP archive...")
	status.CurrentPage = "Creating ZIP archive..."
	status.Progress = 95
	updateStatus(jobID, status)
	sendProgressUpdate(jobID, "progress", status, "Creating ZIP archive...", broadcastProgress)

	// Generate filename: domain_timestamp.zip
	zipFilename := generateZipFilename(url)
	zipPath := filepath.Join(outputDir, zipFilename)
	err = CreateZIP(workDir, zipPath)
	if err != nil {
		updateJobError(jobID, fmt.Sprintf("Failed to create ZIP: %v", err), updateStatus, broadcastProgress)
		return
	}

	// Update final status
	endTime := time.Now()
	status.Status = "completed"
	status.EndTime = &endTime
	status.ZipPath = zipPath
	status.Progress = 100
	status.CurrentPage = "Complete!"

	updateStatus(jobID, status)

	// Send completion message
	sendProgressUpdate(jobID, "complete", status, "Download complete!", broadcastProgress)

	utils.LogInfo("Job %s completed successfully", jobID)

	// Cleanup work directory
	os.RemoveAll(workDir)
}

// rewriteAllFiles rewrites URLs in all HTML and CSS files
func rewriteAllFiles(workDir, baseURL string) error {
	return filepath.Walk(workDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Get file extension
		ext := filepath.Ext(path)

		// Read file
		content, err := os.ReadFile(path)
		if err != nil {
			utils.LogError("Failed to read file %s: %v", path, err)
			return nil // Continue with other files
		}

		var rewritten []byte

		// Determine file URL from path
		relPath, _ := filepath.Rel(workDir, path)
		fileURL := baseURL + "/" + filepath.ToSlash(relPath)

		// Rewrite based on file type
		switch ext {
		case ".html", ".htm":
			rewritten, err = RewriteHTML(content, fileURL, baseURL)
		case ".css":
			rewritten, err = RewriteCSS(content, fileURL, baseURL)
		default:
			return nil // Skip other files
		}

		if err != nil {
			utils.LogError("Failed to rewrite %s: %v", path, err)
			return nil // Continue with other files
		}

		// Write back
		err = os.WriteFile(path, rewritten, 0644)
		if err != nil {
			utils.LogError("Failed to write file %s: %v", path, err)
		}

		return nil
	})
}

// sendProgressUpdate sends progress update via callback
func sendProgressUpdate(
	jobID, updateType string,
	status *models.JobStatus,
	message string,
	broadcastProgress ProgressBroadcastFunc,
) {
	update := models.ProgressUpdate{
		Type:             updateType,
		JobID:            jobID,
		PagesDiscovered:  status.PagesDiscovered,
		PagesDownloaded:  status.PagesDownloaded,
		AssetsDownloaded: status.AssetsDownloaded,
		TotalSize:        status.TotalSize,
		CurrentPage:      status.CurrentPage,
		Progress:         status.Progress,
		Message:          message,
		Error:            status.Error,
	}

	broadcastProgress(jobID, update)
}

// updateJobError updates job status with error using callbacks
func updateJobError(
	jobID, errorMsg string,
	updateStatus StatusUpdateFunc,
	broadcastProgress ProgressBroadcastFunc,
) {
	utils.LogError("Job %s error: %s", jobID, errorMsg)

	status := &models.JobStatus{
		JobID:  jobID,
		Status: "failed",
		Error:  errorMsg,
	}
	endTime := time.Now()
	status.EndTime = &endTime

	updateStatus(jobID, status)
	sendProgressUpdate(jobID, "error", status, errorMsg, broadcastProgress)
}


// generateZipFilename creates a filename from URL: domain_timestamp.zip
func generateZipFilename(urlStr string) string {
	// Parse URL to get domain
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		// Fallback to timestamp if URL parsing fails
		return fmt.Sprintf("website_%s.zip", time.Now().Format("20060102_150405"))
	}
	
	// Get domain without www
	domain := parsedURL.Hostname()
	domain = strings.TrimPrefix(domain, "www.")
	
	// Clean domain for filename (replace dots and special chars)
	domain = strings.ReplaceAll(domain, ".", "_")
	domain = regexp.MustCompile(`[^a-zA-Z0-9_-]`).ReplaceAllString(domain, "_")
	
	// Create timestamp
	timestamp := time.Now().Format("20060102_150405")
	
	// Return formatted filename
	return fmt.Sprintf("%s_%s.zip", domain, timestamp)
}
