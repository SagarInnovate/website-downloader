# 🌐 Website Downloader

A powerful, production-ready tool to download entire websites for offline viewing. Built with **React** and **Golang**, featuring real-time progress tracking, concurrent downloads, and intelligent asset extraction.

![Website Downloader Interface](file:///C:/Users/Sagar%20Shinde/.gemini/antigravity/artifacts/app_interface.webp)

## ✨ Features

### Core Functionality
- 🎯 **Exact Domain Filtering** - Downloads only from the specified domain (excludes subdomains)
- 📄 **Complete HTML Crawling** - Recursively downloads all linked pages
- 🎨 **CSS Asset Extraction** - Parses CSS files to download fonts, background images, and imports
- 🖼️ **Asset Download** - Downloads images, JavaScript, fonts, and other resources
- 🔗 **Link Rewriting** - Converts absolute URLs to relative paths for offline viewing
- 📦 **ZIP Packaging** - Creates a downloadable archive with proper directory structure

### Advanced Features
- ⚡ **Concurrent Downloads** - Uses goroutines with worker pools for maximum speed
- 🔄 **Real-time Progress** - WebSocket-based live updates during scraping
- 🛡️ **Rate Limiting** - Configurable requests per second to avoid overwhelming servers
- 📊 **Progress Tracking** - Shows pages discovered, downloaded, assets, and total size
- ⚙️ **Configurable** - Customize depth, page limits, timeouts, and more
- 🎯 **Smart Parsing** - Extracts assets from HTML, CSS, and nested imports

## 🚀 Quick Start

### Prerequisites
- **Go** 1.20 or higher
- **Node.js** 18 or higher
- **npm** or **yarn**

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/yourusername/website-downloader.git
cd website-downloader
```

2. **Install backend dependencies**
```bash
cd backend
go mod download
```

3. **Install frontend dependencies**
```bash
cd ../frontend
npm install
```

### Running Locally

**Option 1: Using the test scripts (Recommended)**

Windows:
```bash
.\test.bat
```

Linux/Mac:
```bash
chmod +x test.sh
./test.sh
```

**Option 2: Manual startup**

Terminal 1 (Backend):
```bash
cd backend
go run main.go
```

Terminal 2 (Frontend):
```bash
cd frontend
npm run dev
```

Open http://localhost:5174 in your browser.

## 📖 Usage

1. **Enter a website URL** in the input field (e.g., `https://example.com`)
2. **Click "Download Website"**
3. **Watch real-time progress** as it:
   - Discovers and crawls pages
   - Downloads HTML, CSS, and JavaScript
   - Extracts fonts and images from CSS
   - Downloads all assets
4. **Download the ZIP** when complete
5. **Extract and open** `index.html` for offline viewing

## 🏗️ Project Structure

```
website-downloader/
├── backend/                    # Golang backend
│   ├── handlers/              # API endpoints
│   ├── scraper/               # Core scraping logic
│   │   ├── crawler.go        # URL crawling & queue management
│   │   ├── downloader.go     # Asset downloading
│   │   ├── rewriter.go       # Link rewriting
│   │   ├── zipper.go         # ZIP packaging
│   │   └── job.go            # Job orchestration
│   ├── models/                # Data structures
│   ├── utils/                 # Utilities (URL, logger, rate limiter)
│   ├── config.json           # Backend configuration
│   └── main.go               # Entry point
│
├── frontend/                  # React frontend
│   ├── src/
│   │   ├── components/       # UI components
│   │   ├── services/         # API client
│   │   ├── hooks/            # WebSocket hooks
│   │   └── App.jsx           # Main component
│   └── package.json
│
├── docker-compose.yml         # Docker orchestration
├── .gitignore
└── README.md
```

## ⚙️ Configuration

Edit `backend/config.json` to customize:

```json
{
  "maxDepth": 5,              // Maximum crawl depth
  "maxPages": 1000,           // Maximum pages to download
  "workerCount": 10,          // Concurrent download workers
  "requestTimeout": 30,       // Request timeout (seconds)
  "rateLimit": 5,            // Requests per second
  "allowedExtensions": [     // File types to download
    ".html", ".css", ".js",
    ".jpg", ".png", ".gif",
    ".woff", ".woff2", ".ttf"
  ]
}
```

## 🐳 Docker Deployment

```bash
# Build and run with Docker Compose
docker-compose up -d

# Access the application
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
```

## 📡 API Reference

### Start Scraping
```http
POST /api/scrape
Content-Type: application/json

{
  "url": "https://example.com"
}

Response: {
  "jobId": "uuid",
  "status": "started"
}
```

### Get Status
```http
GET /api/status/:jobId

Response: {
  "status": "processing",
  "pagesDownloaded": 10,
  "assetsDownloaded": 45
}
```

### Download ZIP
```http
GET /api/download/:jobId

Returns: ZIP file
```

### WebSocket Progress
```javascript
ws://localhost:8080/api/progress/:jobId

Message: {
  "type": "progress",
  "pagesDiscovered": 15,
  "pagesDownloaded": 10,
  "assetsDownloaded": 45,
  "totalSize": 1024000,
  "currentPage": "https://example.com/page",
  "progress": 67
}
```

## 🎯 What Gets Downloaded

When you download `https://example.com`:

✅ **Included:**
- HTML pages from `example.com`
- CSS files and their fonts/background images
- JavaScript files
- Images in `<img>` tags
- Assets referenced in CSS via `url()`
- All files from the exact same domain

❌ **Excluded:**
- `www.example.com` (subdomain)
- `blog.example.com` (subdomain)
- `external-site.com` (different domain)
- External CDN resources

## 🛠️ Technology Stack

### Backend
- **Language:** Go 1.20+
- **Framework:** Gin (HTTP server & routing)
- **HTML Parsing:** golang.org/x/net/html
- **WebSocket:** gorilla/websocket
- **Concurrency:** Goroutines with worker pools

### Frontend
- **Framework:** React 18
- **Build Tool:** Vite
- **HTTP Client:** Axios
- **Real-time:** Native WebSocket API
- **Styling:** Modern CSS

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with ❤️ using React and Golang
- Inspired by the need for reliable offline website viewing
- Thanks to all contributors!

## 📞 Support

If you encounter any issues or have questions:
- 🐛 [Report a Bug](https://github.com/yourusername/website-downloader/issues)
- 💡 [Request a Feature](https://github.com/yourusername/website-downloader/issues)
- 📧 Email: your.email@example.com

## 🌟 Star History

If you find this project useful, please consider giving it a star! ⭐

---

**Made with 💙 by [Your Name]**
