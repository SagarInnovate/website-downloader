package handlers

import (
	"net/http"
	"sync"
	"time"
	"website-downloader/models"
	"website-downloader/scraper"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	// In-memory storage for job statuses (in production, use Redis or database)
	jobStatuses = make(map[string]*models.JobStatus)
	jobsMutex   sync.RWMutex
)

// StartScrapeHandler handles POST /api/scrape
func StartScrapeHandler(c *gin.Context) {
	var req models.ScrapeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Generate unique job ID
	jobID := uuid.New().String()

	// Initialize job status
	jobStatus := &models.JobStatus{
		JobID:     jobID,
		Status:    "started",
		URL:       req.URL,
		StartTime: time.Now(),
	}

	// Store job status
	UpdateJobStatus(jobID, jobStatus)

	// Load configuration
	config, err := models.LoadConfig("config.json")
	if err != nil {
		// Use default config if loading fails
		config = models.DefaultConfig()
	}

	// Start scraping in background goroutine with callbacks
	go scraper.StartJob(
		jobID,
		req.URL,
		config,
		UpdateJobStatus,
		BroadcastProgress,
	)

	// Return response immediately
	c.JSON(http.StatusOK, models.ScrapeResponse{
		JobID:  jobID,
		Status: "started",
	})
}

// GetStatusHandler handles GET /api/status/:jobId
func GetStatusHandler(c *gin.Context) {
	jobID := c.Param("jobId")

	jobsMutex.RLock()
	status, exists := jobStatuses[jobID]
	jobsMutex.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	c.JSON(http.StatusOK, status)
}

// DownloadHandler handles GET /api/download/:jobId
func DownloadHandler(c *gin.Context) {
	jobID := c.Param("jobId")

	jobsMutex.RLock()
	status, exists := jobStatuses[jobID]
	jobsMutex.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	if status.Status != "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Job not completed yet"})
		return
	}

	if status.ZipPath == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ZIP file not found"})
		return
	}

	// Serve the ZIP file
	c.FileAttachment(status.ZipPath, "website.zip")
}

// GetJobStatus retrieves job status (helper function)
func GetJobStatus(jobID string) (*models.JobStatus, bool) {
	jobsMutex.RLock()
	defer jobsMutex.RUnlock()
	status, exists := jobStatuses[jobID]
	return status, exists
}

// UpdateJobStatus updates job status (helper function)
func UpdateJobStatus(jobID string, status *models.JobStatus) {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()
	jobStatuses[jobID] = status
}
