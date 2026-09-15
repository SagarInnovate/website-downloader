# Website Downloader

A professional desktop application for archiving entire websites for offline use. Built with Wails v2, Go, and modern web technologies.

![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey.svg)

## 🎯 Features

### Two Download Modes
- **⚡ Static Fetch**: Lightning-fast downloads for traditional websites (blogs, documentation, static sites)
- **🌐 Browser Render**: Full JavaScript execution for modern SPAs (React, Vue, Angular applications)

### Intelligent SPA Support
- Automatically discovers all routes in single-page applications
- Waits for network idle before capturing content
- Handles dynamic content with intelligent waiting strategies

### User-Friendly Interface
- Native desktop experience (no browser needed)
- Real-time progress tracking with detailed stats
- Stop/Cancel downloads at any time
- Professional sidebar navigation
- Download history with persistent storage (last 100 downloads)
- Native OS dialogs (Windows, macOS, Linux)

### Professional Output
- Creates ZIP files with smart naming: `domain_timestamp.zip`
- Native "Save As..." dialog to choose save location
- Preserves complete website structure
- Includes all assets (HTML, CSS, JS, images, fonts)

## 📸 Screenshots

*Professional UI with sidebar navigation, real-time progress tracking, and native feel*

## 🚀 Quick Start

### Prerequisites
- Go 1.18 or higher
- Node.js 16 or higher
- Wails CLI v2

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd website-downloader-app
   ```

2. **Install Wails CLI** (if not already installed)
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```

3. **Run development server**
   ```bash
   wails dev
   ```

### Building for Production

**Windows**
```bash
wails build -platform windows/amd64
```

**macOS**
```bash
wails build -platform darwin/universal
```

**Linux**
```bash
wails build -platform linux/amd64
```

The executable will be created in the `build/bin` directory.

## 🛠️ Technology Stack

- **Backend**: Go with Chromedp for browser automation
- **Frontend**: Vanilla JavaScript + Modern CSS
- **Framework**: Wails v2 for native desktop apps
- **Platform**: Cross-platform (Windows, macOS, Linux)

## 📖 Usage

### Basic Workflow

1. **Enter URL**: Type or paste the website URL (e.g., `https://example.com`)
2. **Choose Mode**:
   - Select **Static Fetch** for traditional websites
   - Select **Browser Render** for JavaScript-heavy SPAs
3. **Start Download**: Click "Start Archiving"
4. **Monitor Progress**: Watch real-time stats (pages, assets, size)
5. **Save**: Use "Save As..." to choose where to save the ZIP file

### Download Modes Explained

#### Static Fetch
Best for:
- Blogs and news sites
- Documentation sites
- Traditional server-rendered websites
- Static HTML sites

Characteristics:
- Fastest performance
- Lower resource usage
- No JavaScript execution

#### Browser Render
Best for:
- Single Page Applications (SPAs)
- React, Vue, Angular apps
- Websites with heavy JavaScript
- Dynamic content loading

Characteristics:
- Full JavaScript execution
- Route discovery for SPAs
- Network idle detection
- Higher resource usage

### Configuration

Edit `backend/config.json` to customize:

```json
{
  "maxPages": 500,
  "browserWaitTime": 5000,
  "browserTimeout": 60,
  "downloadPath": "./downloads"
}
```

- `maxPages`: Maximum pages to download per job
- `browserWaitTime`: Time to wait after page load (ms)
- `browserTimeout`: Page load timeout (seconds)
- `downloadPath`: Default save location

## 🎨 Features Overview

### Home View
- Clean, professional interface
- Quick URL input
- Mode selection (Static/Browser)
- Real-time progress tracking

### History View
- View all past downloads
- Quick access to saved files
- Metadata (pages, assets, size, date)
- Open file location with one click

### Settings View
- Configure download limits
- Adjust browser timing
- Manage storage location
- All settings in one place

### About View
- Version information
- Feature list
- Technology stack
- Developer info

## 🔧 Development

### Project Structure
```
website-downloader-app/
├── backend/
│   ├── handlers/          # HTTP & WebSocket handlers
│   ├── models/           # Data models
│   ├── scraper/          # Core scraping logic
│   │   ├── crawler.go
│   │   ├── browser_crawler.go
│   │   ├── downloader.go
│   │   ├── link_rewriter.go
│   │   └── zipper.go
│   ├── downloads/        # Downloaded ZIP files
│   ├── temp/             # Temporary scraping files
│   └── main.go          # Backend entry point
├── frontend/
│   ├── src/
│   │   ├── main.js      # Frontend logic
│   │   ├── template.js  # HTML templates
│   │   ├── app.css      # Main styles
│   │   ├── sidebar.css  # Sidebar styles
│   │   └── style.css    # Base styles
│   └── index.html       # Entry point
├── wails.json           # Wails configuration
└── go.mod              # Go dependencies
```

### Key Backend Functions

- `StartScrape(url, mode)` - Initiates download
- `CancelJob(jobID)` - Stops active download
- `SaveAsZip(jobID)` - Opens native save dialog
- `GetHistory()` - Returns download history
- `OpenFileLocation(path)` - Opens file in explorer

### Frontend Architecture

- **View Management**: Sidebar navigation with Home, History, Settings, About
- **Real-time Updates**: WebSocket progress tracking
- **Native Dialogs**: Windows/macOS/Linux system dialogs
- **Responsive**: Smooth scrolling and native feel

## 🐛 Troubleshooting

### App won't start
- Ensure Go and Node.js are installed
- Run `go mod tidy` in the backend directory
- Run `npm install` in the frontend directory

### Download fails
- Check internet connection
- Verify URL is accessible
- Try different mode (Static vs Browser)
- Check console for error messages

### Browser mode not working
- Ensure Chromium/Chrome is installed
- Check browser timeout settings
- Increase wait time for slow sites

## 📝 License

MIT License - feel free to use this project for personal or commercial purposes.

## 👨‍💻 Developer

Created by **sagarinnovate**

Contributions and feedback are welcome! Open an issue or submit a pull request.

## 🌟 Acknowledgments

- Built with [Wails](https://wails.io/) - Amazing Go + Web framework
- Powered by [Chromedp](https://github.com/chromedp/chromedp) - Headless Chrome automation
- Inspired by the need for professional website archiving tools

## 🔮 Future Enhancements

Potential features for future releases:
- [ ] Download queue management
- [ ] Scheduled downloads
- [ ] Export/Import history
- [ ] Advanced filtering options
- [ ] Custom user agents
- [ ] Proxy support
- [ ] Cloud storage integration
- [ ] Multi-language support

---

**Made with ❤️ for the open-source community**
