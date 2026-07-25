# Website Downloader Desktop App

A powerful desktop application built with [Wails](https://wails.io/) that downloads entire websites for offline viewing, with special support for Single Page Applications (SPAs).

## Features

✨ **Single Executable** - No backend setup required, just download and run  
🚀 **Fast Static Site Downloads** - Traditional websites download in seconds  
⚛️ **SPA Support** - Intelligent crawling for React, Vue, Angular apps  
🎯 **Smart Route Discovery** - Automatically finds all pages in SPAs  
🌐 **Network Idle Detection** - Waits for complete page rendering  
💾 **Offline Viewing** - Downloads include all assets and rewritten links  
🖥️ **Native Desktop Experience** - Works on Windows, macOS, and Linux  

## Installation

### For Users (Pre-built Binary)

1. Download the latest release for your platform:
   - **Windows**: `website-downloader-windows-amd64.exe`
   - **macOS**: `website-downloader-darwin-universal`
   - **Linux**: `website-downloader-linux-amd64`

2. Double-click to run - no installation needed!

### For Developers (Build from Source)

#### Prerequisites

- Go 1.23 or later
- Node.js 16+ and npm
- Wails CLI v2.12.0+

**Install Wails CLI:**
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

#### Build Instructions

1. **Clone the repository:**
```bash
git clone <your-repo-url>
cd website-downloader-app
```

2. **Install dependencies:**
```bash
wails build
```

This will:
- Install Go dependencies
- Install npm dependencies
- Build the frontend
- Build the executable

3. **Find your executable:**
- Windows: `build/bin/website-downloader-app.exe`
- macOS: `build/bin/website-downloader-app.app`
- Linux: `build/bin/website-downloader-app`

## Development

### Run in Development Mode

```bash
wails dev
```

This starts the app with:
- Hot reload for frontend changes
- Auto-restart for Go code changes  
- DevTools enabled

### Project Structure

```
website-downloader-app/
├── frontend/           # Frontend UI (HTML/CSS/JS)
│   ├── src/
│   ├── dist/          # Built frontend (auto-generated)
│   └── package.json
├── models/            # Data models
├── scraper/           # Website scraping logic
├── utils/             # Utility functions
├── app.go             # Main app logic (Wails bindings)
├── main.go            # Entry point
├── config.json        # Default configuration
└── wails.json         # Wails project configuration
```

## Usage

### Basic Usage

1. **Launch the app**
2. **Enter a URL** (e.g., `https://example.com`)
3. **Select mode:**
   - **Static Mode** (faster) - For traditional websites
   - **Browser Mode** - For SPAs (React, Vue, Angular)
4. **Click "Download"**
5. **Wait for completion** - Progress shown in real-time
6. **Open downloaded files** - Click "Open Folder" or "Save As"

### Download Modes

#### Static Mode (Default)
- ✅ **Fast** - 10-20x faster than browser mode
- ✅ **Traditional websites** - Server-rendered HTML
- ❌ **Not suitable for SPAs** - May miss dynamic content

#### Browser Mode
- ✅ **Complete SPA support** - React, Vue, Angular, etc.
- ✅ **Dynamic content** - Waits for JavaScript to load
- ✅ **Smart route discovery** - Finds all SPA routes automatically
- ❌ **Slower** - Uses headless Chrome

### Configuration

The app uses `config.json` for settings:

```json
{
  "maxPages": 1000,           // Maximum pages to download
  "requestTimeout": 30,        // HTTP timeout (seconds)
  "browserWaitTime": 5000,    // Max wait for JS rendering (ms)
  "browserTimeout": 60,       // Browser page timeout (seconds)
  "tempDir": "./temp",        // Temporary file directory
  "outputDir": "./downloads"  // Download output directory
}
```

You can modify these settings before building to change defaults.

## Building for Distribution

### Windows

```bash
wails build -platform windows/amd64
```

Optional: Create an installer
```bash
wails build -platform windows/amd64 -nsis
```

### macOS

```bash
wails build -platform darwin/universal
```

Optional: Code sign (requires Apple Developer account)
```bash
wails build -platform darwin/universal -sign
```

### Linux

```bash
wails build -platform linux/amd64
```

### Cross-Platform Build

```bash
wails build -platform windows/amd64,darwin/universal,linux/amd64
```

## Technical Details

### How SPA Crawling Works

1. **Initial Load** - Opens the base URL in headless Chrome
2. **Route Discovery** - Analyzes JavaScript and DOM for route patterns
3. **Sequential Navigation** - Visits each route using History API
4. **Network Idle Detection** - Waits for all HTTP requests to complete
5. **Content Capture** - Extracts fully rendered HTML
6. **Asset Download** - Downloads CSS, JS, images, fonts
7. **Link Rewriting** - Makes links relative for offline viewing
8. **ZIP Creation** - Packages everything into a single archive

### Supported Frameworks

- ✅ React (React Router)
- ✅ Vue (Vue Router)
- ✅ Angular (Angular Router)
- ✅ Next.js (static exports)
- ✅ Svelte/SvelteKit
- ✅ Any client-side routing framework

## Troubleshooting

### Chrome/Chromium Required

The browser mode requires Chrome/Chromium to be installed:
- **Windows**: Usually pre-installed or download from google.com/chrome
- **macOS**: Download from google.com/chrome
- **Linux**: `sudo apt install chromium-browser` or `google-chrome-stable`

### "Job not found" Error

This means the download hasn't started. Possible causes:
- Invalid URL
- Network connectivity issues
- Chrome/Chromium not found (browser mode only)

### Incomplete Downloads

If downloads are incomplete:
1. Increase `browserWaitTime` in config.json
2. Increase `browserTimeout` for slow sites
3. Check if the site requires authentication
4. Try browser mode if using static mode

### App Won't Start

- Ensure all dependencies are installed: `go mod tidy`
- Rebuild: `wails build`
- Check logs in console output

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

[Specify your license here]

## Credits

Built with:
- [Wails](https://wails.io/) - Go + Web framework
- [chromedp](https://github.com/chromedp/chromedp) - Chrome DevTools Protocol
- Modern web technologies

## Support

For issues, questions, or feature requests:
- Open an issue on GitHub
- Check existing documentation
- Review troubleshooting section

---

**Made with ❤️ for the open-source community**
