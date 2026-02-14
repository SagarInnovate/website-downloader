package models

import "time"

// ScrapeRequest represents the incoming scrape job request
type ScrapeRequest struct {
	URL  string `json:"url" binding:"required"`
	Mode string `json:"mode"` // "static" (default) or "browser" for SPA support
}

// ScrapeResponse represents the API response after starting a scrape job
type ScrapeResponse struct {
	JobID  string `json:"jobId"`
	Status string `json:"status"`
}

// JobStatus represents the current status of a scrape job
type JobStatus struct {
	JobID           string    `json:"jobId"`
	Status          string    `json:"status"` // "running", "completed", "failed"
	URL             string    `json:"url"`
	PagesDiscovered int       `json:"pagesDiscovered"`
	PagesDownloaded int       `json:"pagesDownloaded"`
	AssetsDownloaded int      `json:"assetsDownloaded"`
	TotalSize       int64     `json:"totalSize"` // in bytes
	CurrentPage     string    `json:"currentPage"`
	Progress        int       `json:"progress"` // 0-100
	Error           string    `json:"error,omitempty"`
	StartTime       time.Time `json:"startTime"`
	EndTime         *time.Time `json:"endTime,omitempty"`
	ZipPath         string    `json:"-"` // Internal use only
}

// ProgressUpdate represents real-time progress updates sent via WebSocket
type ProgressUpdate struct {
	Type             string `json:"type"` // "progress", "complete", "error"
	JobID            string `json:"jobId"`
	PagesDiscovered  int    `json:"pagesDiscovered"`
	PagesDownloaded  int    `json:"pagesDownloaded"`
	AssetsDownloaded int    `json:"assetsDownloaded"`
	TotalSize        int64  `json:"totalSize"`
	CurrentPage      string `json:"currentPage"`
	Progress         int    `json:"progress"`
	Message          string `json:"message,omitempty"`
	Error            string `json:"error,omitempty"`
}

