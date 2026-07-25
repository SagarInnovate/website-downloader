package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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
	ctx            context.Context
	jobStatuses    map[string]*models.JobStatus
	jobCancellers  map[string]context.CancelFunc
	downloadHistory []HistoryItem
	jobsMutex      sync.RWMutex
}

// HistoryItem represents a download history entry
type HistoryItem struct {
	JobID      string    `json:"jobId"`
	URL        string    `json:"url"`
	Mode       string    `json:"mode"`
	Status     string    `json:"status"`
	ZipPath    string    `json:"zipPath"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
	Pages      int       `json:"pages"`
	Assets     int       `json:"assets"`
	TotalSize  int64     `json:"totalSize"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		jobStatuses:     make(map[string]*models.JobStatus),
		jobCancellers:   make(map[string]context.CancelFunc),
		downloadHistory: []HistoryItem{},
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// Load history from disk
	a.loadHistory()
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

// SaveAsZip opens a save dialog and copies the ZIP to user-selected location
func (a *App) SaveAsZip(jobID string) (string, error) {
	status, err := a.GetJobStatus(jobID)
	if err != nil {
		return "", err
	}

	if status.ZipPath == "" {
		return "", fmt.Errorf("ZIP file not found")
	}

	// Check if source file exists
	if _, err := os.Stat(status.ZipPath); os.IsNotExist(err) {
		return "", fmt.Errorf("ZIP file no longer exists")
	}

	// Get filename from path
	filename := filepath.Base(status.ZipPath)

	// Open save dialog
	savePath, err := wailsruntime.SaveFileDialog(a.ctx, wailsruntime.SaveDialogOptions{
		DefaultFilename: filename,
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

	// Copy file to selected location
	err = copyFile(status.ZipPath, savePath)
	if err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	return savePath, nil
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

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
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
	
	// If job completed, add to history
	if update.Type == "complete" {
		a.jobsMutex.RLock()
		status, exists := a.jobStatuses[jobID]
		a.jobsMutex.RUnlock()
		
		if exists {
			// Get mode from the job (we need to track this)
			mode := "static" // Default, should be tracked properly
			a.addToHistory(status, mode)
		}
	}
	
	// Emit event to frontend
	wailsruntime.EventsEmit(a.ctx, "progress:"+jobID, update)
}

// CancelJob cancels a running scraping job
func (a *App) CancelJob(jobID string) error {
	a.jobsMutex.Lock()
	defer a.jobsMutex.Unlock()

	// Get the cancel function
	cancelFunc, exists := a.jobCancellers[jobID]
	if !exists {
		return fmt.Errorf("job not found or already completed")
	}

	// Cancel the job
	cancelFunc()

	// Update status
	if status, ok := a.jobStatuses[jobID]; ok {
		status.Status = "cancelled"
		status.Error = "Cancelled by user"
		endTime := time.Now()
		status.EndTime = &endTime
	}

	// Remove canceller
	delete(a.jobCancellers, jobID)

	// Broadcast cancellation
	a.broadcastProgress(jobID, models.ProgressUpdate{
		Type:   "error",
		JobID:  jobID,
		Error:  "Download cancelled by user",
	})

	return nil
}

// GetHistory returns the download history
func (a *App) GetHistory() []HistoryItem {
	a.jobsMutex.RLock()
	defer a.jobsMutex.RUnlock()
	
	// Return a copy to avoid concurrent modification
	history := make([]HistoryItem, len(a.downloadHistory))
	copy(history, a.downloadHistory)
	return history
}

// ClearHistory clears all download history
func (a *App) ClearHistory() error {
	a.jobsMutex.Lock()
	defer a.jobsMutex.Unlock()
	
	a.downloadHistory = []HistoryItem{}
	return a.saveHistory()
}

// DeleteHistoryItem removes a specific item from history
func (a *App) DeleteHistoryItem(jobID string) error {
	a.jobsMutex.Lock()
	defer a.jobsMutex.Unlock()
	
	for i, item := range a.downloadHistory {
		if item.JobID == jobID {
			a.downloadHistory = append(a.downloadHistory[:i], a.downloadHistory[i+1:]...)
			return a.saveHistory()
		}
	}
	
	return fmt.Errorf("history item not found")
}

// OpenFileLocation opens the file explorer at the downloaded file location
func (a *App) OpenFileLocation(zipPath string) error {
	// Check if file exists
	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		return fmt.Errorf("file not found")
	}
	
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", "/select,", zipPath)
	case "darwin":
		cmd = exec.Command("open", "-R", zipPath)
	case "linux":
		// Open the folder containing the file
		cmd = exec.Command("xdg-open", zipPath)
	default:
		return fmt.Errorf("unsupported operating system")
	}
	
	return cmd.Start()
}

// addToHistory adds a completed job to history
func (a *App) addToHistory(status *models.JobStatus, mode string) {
	a.jobsMutex.Lock()
	defer a.jobsMutex.Unlock()
	
	if status.Status != "completed" {
		return
	}
	
	endTime := time.Now()
	if status.EndTime != nil {
		endTime = *status.EndTime
	}
	
	historyItem := HistoryItem{
		JobID:     status.JobID,
		URL:       status.URL,
		Mode:      mode,
		Status:    status.Status,
		ZipPath:   status.ZipPath,
		StartTime: status.StartTime,
		EndTime:   endTime,
		Pages:     status.PagesDownloaded,
		Assets:    status.AssetsDownloaded,
		TotalSize: status.TotalSize,
	}
	
	// Add to beginning of history (most recent first)
	a.downloadHistory = append([]HistoryItem{historyItem}, a.downloadHistory...)
	
	// Keep only last 100 items
	if len(a.downloadHistory) > 100 {
		a.downloadHistory = a.downloadHistory[:100]
	}
	
	// Save to disk
	a.saveHistory()
}

// saveHistory saves download history to disk
func (a *App) saveHistory() error {
	config, _ := models.LoadConfig("config.json")
	if config == nil {
		config = models.DefaultConfig()
	}
	
	// Create history file path
	historyPath := config.OutputDir + "/history.json"
	
	// Marshal history to JSON
	data, err := json.Marshal(a.downloadHistory)
	if err != nil {
		return err
	}
	
	// Write to file
	return os.WriteFile(historyPath, data, 0644)
}

// loadHistory loads download history from disk
func (a *App) loadHistory() {
	config, _ := models.LoadConfig("config.json")
	if config == nil {
		config = models.DefaultConfig()
	}
	
	historyPath := config.OutputDir + "/history.json"
	
	// Read file
	data, err := os.ReadFile(historyPath)
	if err != nil {
		// File doesn't exist or can't be read - start with empty history
		a.downloadHistory = []HistoryItem{}
		return
	}
	
	// Unmarshal JSON
	var history []HistoryItem
	if err := json.Unmarshal(data, &history); err != nil {
		// Invalid JSON - start with empty history
		a.downloadHistory = []HistoryItem{}
		return
	}
	
	a.downloadHistory = history
}
