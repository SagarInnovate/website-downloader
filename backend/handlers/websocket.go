package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
	"website-downloader/models"
	"website-downloader/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins (adjust for production)
		},
	}
	
	// Store WebSocket connections per job
	wsConnections = make(map[string][]*websocket.Conn)
	wsMutex       sync.RWMutex
)

// WebSocketHandler handles WebSocket connections for progress updates
func WebSocketHandler(c *gin.Context) {
	jobID := c.Param("jobId")

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.LogError("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// Register connection
	wsMutex.Lock()
	wsConnections[jobID] = append(wsConnections[jobID], conn)
	wsMutex.Unlock()

	utils.LogInfo("WebSocket connected for job: %s", jobID)

	// Send initial status
	if status, exists := GetJobStatus(jobID); exists {
		update := models.ProgressUpdate{
			Type:             "progress",
			JobID:            jobID,
			PagesDiscovered:  status.PagesDiscovered,
			PagesDownloaded:  status.PagesDownloaded,
			AssetsDownloaded: status.AssetsDownloaded,
			TotalSize:        status.TotalSize,
			CurrentPage:      status.CurrentPage,
			Progress:         status.Progress,
		}
		conn.WriteJSON(update)
	}

	// Keep connection alive and listen for close
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			// Connection closed
			removeWebSocketConnection(jobID, conn)
			break
		}
	}
}

// BroadcastProgress sends progress updates to all connected WebSocket clients
func BroadcastProgress(jobID string, update models.ProgressUpdate) {
	wsMutex.RLock()
	connections := wsConnections[jobID]
	wsMutex.RUnlock()

	for _, conn := range connections {
		err := conn.WriteJSON(update)
		if err != nil {
			utils.LogError("Failed to send WebSocket message: %v", err)
			removeWebSocketConnection(jobID, conn)
		}
	}
}

// removeWebSocketConnection removes a connection from the map
func removeWebSocketConnection(jobID string, conn *websocket.Conn) {
	wsMutex.Lock()
	defer wsMutex.Unlock()

	connections := wsConnections[jobID]
	for i, c := range connections {
		if c == conn {
			wsConnections[jobID] = append(connections[:i], connections[i+1:]...)
			break
		}
	}
}

// SendProgressUpdate is a helper to broadcast progress updates
func SendProgressUpdate(jobID string, updateType string, status *models.JobStatus, message string) {
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

	// Broadcast to WebSocket clients
	BroadcastProgress(jobID, update)

	// Log the update
	data, _ := json.Marshal(update)
	utils.LogInfo("Progress update: %s", string(data))
}
