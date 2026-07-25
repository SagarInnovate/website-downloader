package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
	
	"website-downloader-app/models"
	"website-downloader-app/scraper"
)

// App struct
type App struct {
	ctx         context.Context
	jobStatuses map[string]*models.JobStatus
	jobsMutex   sync.RWMutex
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		jobStatuses: make(map[string]*models.JobStatus),
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// ScrapeRequest represents the scrape request from frontend
type ScrapeRequest struct {
	URL  string `json:"url"`
	Mode string `json:"mode"`
}

// ScrapeResponse represents the response after starting a scrape
type ScrapeResponse struct {
	JobID  string `json:"jobId"`
	Status string `json:"status"`
}

// StartScrape initiates a new website scraping job
func (a *App) StartScrape(url string, mode string) ScrapeResponse {
	// Generate unique job ID
	jobID := uuid.New().String()

	// Default to static mode if not specified
	if mode == "" {
		mode = "static"
	}

	// Initialize job status
	jobStatus := &models.JobStatus{
		JobID:     jobID,
		Status:    "started",
		URL:       url,
		StartTime: time.Now(),
	}

	// Store job status
	a.updateJobStatus(jobID, jobStatus)

	// Load configuration
	config, err := models.LoadConfig("config.json")
	if err != nil {
		// Use default config if loading fails
		config = models.DefaultConfig()
	}

	// Start scraping in background
	go scraper.StartJob(
		jobID,
		url,
		mode,
		config,
		a.updateJobStatus,
		a.broadcastProgress,
	)

	return ScrapeResponse{
		JobID:  jobID,
		Status: "started",
	}
}

// GetJobStatus retrieves the status of a scraping job
func (a *App) GetJobStatus(jobID string) (*models.JobStatus, error) {
	a.jobsMutex.RLock()
	defer a.jobsMutex.RUnlock()

	status, exists := a.jobStatuses[jobID]
	if !exists {
		return nil, fmt.Errorf("job not found")
	}

	return status, nil
}

// GetDownloadPath returns the path to download the scraped website
func (a *App) GetDownloadPath(jobID string) (string, error) {
	a.jobsMutex.RLock()
	defer a.jobsMutex.RUnlock()

	status, exists := a.jobStatuses[jobID]
	if !exists {
		return "", fmt.Errorf("job not found")
	}

	if status.Status != "completed" {
		return "", fmt.Errorf("job not completed yet")
	}

	if status.ZipPath == "" {
		return "", fmt.Errorf("ZIP file not found")
	}

	return status.ZipPath, nil
}

// OpenDownloadFolder opens the downloads folder in the system file manager
func (a *App) OpenDownloadFolder() error {
	config, _ := models.LoadConfig("config.json")
	if config == nil {
		config = models.DefaultConfig()
	}
	
	// Create the folder if it doesn't exist
	err := os.MkdirAll(config.OutputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create downloads folder: %w", err)
	}
	
	// Open folder based on OS
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", config.OutputDir)
	case "darwin":
		cmd = exec.Command("open", config.OutputDir)
	case "linux":
		cmd = exec.Command("xdg-open", config.OutputDir)
	default:
		return fmt.Errorf("unsupported operating system")
	}
	
	return cmd.Start()
}

// SelectSaveLocation opens a file dialog to select where to save the downloaded website
func (a *App) SelectSaveLocation(jobID string) (string, error) {
	status, err := a.GetJobStatus(jobID)
	if err != nil {
		return "", err
	}

	if status.ZipPath == "" {
		return "", fmt.Errorf("ZIP file not found")
	}

	// Open save dialog
	savePath, err := wailsruntime.SaveFileDialog(a.ctx, wailsruntime.SaveDialogOptions{
		DefaultFilename: "website.zip",
		Title:           "Save Website Archive",
		Filters: []wailsruntime.FileFilter{
			{
				DisplayName: "ZIP Archives (*.zip)",
				Pattern:     "*.zip",
			},
		},
	})

	if err != nil {
		return "", err
	}

	if savePath == "" {
		return "", fmt.Errorf("save cancelled")
	}

	// TODO: Copy file from status.ZipPath to savePath
	
	return savePath, nil
}

// updateJobStatus updates the status of a job
func (a *App) updateJobStatus(jobID string, status *models.JobStatus) {
	a.jobsMutex.Lock()
	a.jobStatuses[jobID] = status
	a.jobsMutex.Unlock()
}

// broadcastProgress sends progress updates to the frontend
func (a *App) broadcastProgress(jobID string, update models.ProgressUpdate) {
	// Ensure jobID is set in the update
	update.JobID = jobID
	// Emit event to frontend
	wailsruntime.EventsEmit(a.ctx, "progress:"+jobID, update)
}
