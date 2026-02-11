# API Documentation

## Overview

The Website Downloader API provides RESTful endpoints for scraping websites and WebSocket support for real-time progress updates.

**Base URL**: `http://localhost:8080`

## Authentication

Currently, no authentication is required. For production deployments, consider adding API key authentication.

## Endpoints

### Health Check

Check if the server is running.

```http
GET /health
```

**Response**
```json
{
  "status": "ok"
}
```

**Status Codes**
- `200 OK` - Server is healthy

---

### Start Scraping Job

Initiate a new website scraping job.

```http
POST /api/scrape
```

**Request Body**
```json
{
  "url": "https://example.com"
}
```

**Parameters**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| url | string | Yes | The website URL to scrape |

**Response**
```json
{
  "jobId": "550e8400-e29b-41d4-a716-446655440000",
  "status": "started"
}
```

**Status Codes**
- `200 OK` - Job started successfully
- `400 Bad Request` - Invalid URL or request body
- `500 Internal Server Error` - Server error

**Example**
```bash
curl -X POST http://localhost:8080/api/scrape \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com"}'
```

---

### Get Job Status

Retrieve the current status of a scraping job.

```http
GET /api/status/:jobId
```

**Parameters**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| jobId | string (UUID) | Yes | The job identifier |

**Response**
```json
{
  "jobId": "550e8400-e29b-41d4-a716-446655440000",
  "status": "running",
  "progress": {
    "pagesDiscovered": 50,
    "pagesDownloaded": 25,
    "assetsDownloaded": 150,
    "totalSize": 5242880,
    "currentPage": "https://example.com/about",
    "progress": 50
  }
}
```

**Status Values**
- `started` - Job has been initiated
- `running` - Currently scraping
- `complete` - Scraping finished successfully
- `error` - An error occurred

**Status Codes**
- `200 OK` - Status retrieved successfully
- `404 Not Found` - Job ID not found
- `500 Internal Server Error` - Server error

**Example**
```bash
curl http://localhost:8080/api/status/550e8400-e29b-41d4-a716-446655440000
```

---

### Download ZIP File

Download the scraped website as a ZIP archive.

```http
GET /api/download/:jobId
```

**Parameters**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| jobId | string (UUID) | Yes | The job identifier |

**Response**
- Binary ZIP file stream

**Headers**
- `Content-Type: application/zip`
- `Content-Disposition: attachment; filename="website-{jobId}.zip"`

**Status Codes**
- `200 OK` - ZIP file downloaded successfully
- `404 Not Found` - Job ID not found or ZIP not ready
- `500 Internal Server Error` - Server error

**Example**
```bash
curl -O http://localhost:8080/api/download/550e8400-e29b-41d4-a716-446655440000
```

---

### WebSocket Progress Updates

Receive real-time progress updates for a scraping job.

```
WS /api/progress/:jobId
```

**Parameters**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| jobId | string (UUID) | Yes | The job identifier |

**Message Types**

**Progress Update**
```json
{
  "type": "progress",
  "jobId": "550e8400-e29b-41d4-a716-446655440000",
  "pagesDiscovered": 10,
  "pagesDownloaded": 5,
  "assetsDownloaded": 30,
  "totalSize": 1048576,
  "currentPage": "https://example.com/about",
  "progress": 50
}
```

**Completion**
```json
{
  "type": "complete",
  "jobId": "550e8400-e29b-41d4-a716-446655440000",
  "message": "Scraping completed successfully"
}
```

**Error**
```json
{
  "type": "error",
  "jobId": "550e8400-e29b-41d4-a716-446655440000",
  "error": "Failed to fetch page: timeout"
}
```

**JavaScript Example**
```javascript
const ws = new WebSocket('ws://localhost:8080/api/progress/550e8400-e29b-41d4-a716-446655440000');

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Progress:', data);
  
  if (data.type === 'complete') {
    console.log('Download complete!');
  } else if (data.type === 'error') {
    console.error('Error:', data.error);
  }
};

ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};

ws.onclose = () => {
  console.log('Connection closed');
};
```

---

## Data Models

### Job

```typescript
interface Job {
  jobId: string;        // UUID v4
  url: string;          // Original URL
  status: string;       // started | running | complete | error
  createdAt: string;    // ISO 8601 timestamp
  updatedAt: string;    // ISO 8601 timestamp
  progress: Progress;   // Progress information
}
```

### Progress

```typescript
interface Progress {
  pagesDiscovered: number;   // Total pages found
  pagesDownloaded: number;   // Pages successfully downloaded
  assetsDownloaded: number;  // Assets successfully downloaded
  totalSize: number;         // Total bytes downloaded
  currentPage: string;       // Currently processing page
  progress: number;          // Percentage (0-100)
}
```

### ProgressUpdate

```typescript
interface ProgressUpdate {
  type: string;              // progress | complete | error
  jobId: string;             // Job identifier
  pagesDiscovered?: number;
  pagesDownloaded?: number;
  assetsDownloaded?: number;
  totalSize?: number;
  currentPage?: string;
  progress?: number;
  message?: string;          // For complete/error types
  error?: string;            // For error type
}
```

---

## Error Handling

All errors follow this format:

```json
{
  "error": "Error message description"
}
```

### Common Error Codes

| HTTP Code | Description |
|-----------|-------------|
| 400 | Bad Request - Invalid input or parameters |
| 404 | Not Found - Resource doesn't exist |
| 500 | Internal Server Error - Server-side error |
| 503 | Service Unavailable - Server overloaded |

---

## Rate Limiting

The API implements rate limiting to prevent abuse:

- Default: 5 requests per second
- Configurable in `config.json`
- Applies to all endpoints except `/health`

**Headers** (Future feature)
- `X-RateLimit-Limit: 5`
- `X-RateLimit-Remaining: 4`
- `X-RateLimit-Reset: 1234567890`

---

## CORS

CORS is enabled for the following origins:
- `http://localhost:5173` (Vite dev server)
- `http://localhost:3000` (Common React dev port)

For production, update `main.go` with your domain.

---

## Best Practices

1. **Poll Status Sparingly**: Use WebSocket for real-time updates instead of polling `/api/status`
2. **Handle Errors**: Always implement error handling for network failures
3. **Timeout Management**: Large websites may take several minutes to scrape
4. **Validate URLs**: Validate URLs on the client side before sending to the API
5. **Use Job IDs**: Store job IDs for future reference and downloads

---

## Example Workflow

```javascript
// 1. Start scraping
const response = await fetch('http://localhost:8080/api/scrape', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ url: 'https://example.com' })
});

const { jobId } = await response.json();

// 2. Connect to WebSocket for progress
const ws = new WebSocket(`ws://localhost:8080/api/progress/${jobId}`);

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  
  if (data.type === 'progress') {
    updateProgressBar(data.progress);
  } else if (data.type === 'complete') {
    showDownloadButton(jobId);
  }
};

// 3. Download when complete
function downloadZip(jobId) {
  window.location.href = `http://localhost:8080/api/download/${jobId}`;
}
```

---

## Configuration

API behavior can be customized via `backend/config.json`:

```json
{
  "maxDepth": 5,
  "maxPages": 1000,
  "workerCount": 10,
  "requestTimeout": 30,
  "rateLimit": 5
}
```

See [Configuration Guide](../README.md#configuration) for details.

---

## Support

For issues or questions:
- Check the [main README](../README.md)
- Review [DEPLOYMENT.md](../DEPLOYMENT.md) for production setup
- Open an issue on GitHub
