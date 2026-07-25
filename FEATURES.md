# Website Downloader - Features Overview

## 🎯 Core Features

### 📥 Download Modes

#### Static Mode (Fast)
- Traditional websites with server-side rendering
- HTTP-based crawling
- 10-20x faster than browser mode
- Best for: WordPress, static HTML sites, traditional CMSs

#### Browser Mode (SPAs)
- Full JavaScript execution via headless Chrome
- Intelligent route discovery
- Network idle detection
- Best for: React, Vue, Angular, Next.js apps

### ⏹️ Download Control

**NEW: Stop/Cancel Downloads**
- Click the **Stop button** during download
- Instantly cancels the ongoing operation
- Prevents resource waste
- Clean cancellation with proper cleanup

### 📜 Download History

**NEW: Persistent History Tracking**
- Automatically saves all successful downloads
- Stores up to 100 most recent downloads
- Persisted to disk (`downloads/history.json`)
- Survives app restarts

**History Display:**
- URL of downloaded website
- Download mode used (Static/Browser)
- Timestamp (with relative time: "2 hours ago")
- Statistics: pages, assets, total size
- Quick access to file location

**History Actions:**
- 📁 **Open File Location** - Opens file explorer at the ZIP file
- 🔄 **Refresh** - Reload history from disk
- 🗑️ **Clear All** - Remove all history (with confirmation)

### 📊 Real-Time Progress

**Live Statistics:**
- **Pages Downloaded** - Number of HTML pages saved
- **Assets Downloaded** - CSS, JS, images, fonts, etc.
- **Total Size** - Running total in KB/MB/GB
- **Current Page** - URL being processed right now
- **Progress Bar** - Visual 0-100% completion indicator

**Status Updates:**
- Starting...
- Downloading...
- Rewriting links...
- Creating ZIP archive...
- Complete! / Failed

### 🌐 Intelligent SPA Crawling

**Route Discovery:**
- Analyzes JavaScript bundles for route definitions
- Extracts navigation links from DOM
- Detects React Router, Vue Router, Angular Router
- Finds hidden routes not visible in HTML

**Smart Waiting:**
- Network idle detection (no pending requests)
- No arbitrary delays
- Waits for complete rendering
- Efficient browser context reuse

**Content Capture:**
- Fully rendered HTML for each route
- All loaded assets (CSS, JS, images)
- CDN resources included
- Link rewriting for offline viewing

### 💾 Offline Viewing

**Link Rewriting:**
- Converts absolute URLs to relative paths
- Updates asset references
- Rewrites CSS url() references
- Ensures everything works offline

**ZIP Archive:**
- Complete website in one file
- Organized folder structure
- Easy to share or backup
- Extract and view in any browser

## 🎨 User Interface

### Main Screen
- Clean, modern design
- Purple gradient theme
- Large input fields
- Clear mode selector
- Prominent download button

### Progress Screen
- Real-time updates
- Visual progress bar
- Statistics grid
- Current page display
- **Stop button** for control

### History Panel
- Slide-in panel design
- Scrollable list
- Search-friendly layout
- Quick access icons
- Bulk actions (clear all)

### Completion Screen
- Success/failure indication
- Download summary
- Quick actions:
  - Open downloads folder
  - Download another website

## 🔧 Configuration

**config.json Settings:**
```json
{
  "maxPages": 1000,          // Maximum pages per download
  "requestTimeout": 30,       // HTTP timeout (seconds)
  "browserWaitTime": 5000,   // Max wait for JS (milliseconds)
  "browserTimeout": 60,      // Browser page timeout (seconds)
  "tempDir": "./temp",       // Temporary files
  "outputDir": "./downloads" // Final ZIPs and history
}
```

## 🚀 Performance

### Static Mode
- ⚡ 0.1-0.5 seconds per page
- 💾 ~50-100 MB memory
- 🎯 Best for traditional sites

### Browser Mode
- ⏱️ 2-5 seconds per page
- 💾 ~200-500 MB memory
- 🎯 Required for SPAs

## 🔒 Privacy & Security

- ✅ Everything runs locally
- ✅ No data sent to external servers
- ✅ No telemetry or tracking
- ✅ Your downloads stay private
- ✅ Open source - verify the code yourself

## 📱 Platform Support

### Windows
- ✅ Windows 10/11
- ✅ WebView2 (automatically installed)
- ✅ Native file explorer integration

### macOS
- ✅ macOS 10.13+
- ✅ Universal binary (Intel + Apple Silicon)
- ✅ Native Finder integration

### Linux
- ✅ Most distributions
- ✅ Requires WebKit2GTK
- ✅ Native file manager integration

## 🎯 Use Cases

### Web Archiving
Save important websites before they disappear:
- News articles
- Research papers
- Documentation
- Personal projects

### Offline Access
Work without internet:
- Travel with offline docs
- Presentations without WiFi
- Development reference materials
- Educational content

### Backup & Migration
Preserve website content:
- Before redesigns
- Website migrations
- Content audits
- Legal archives

### Testing & Development
Download production sites:
- Test responsive design
- Analyze site structure
- Study competitor websites
- Learn from examples

## 🆕 Recent Updates

### v1.1.0 - Download Control & History
- ✨ Stop/Cancel button during downloads
- ✨ Download history tracking
- ✨ Persistent storage
- ✨ Open file location from history
- ✨ Clear history functionality
- 🎨 Improved UI layout
- 🐛 Bug fixes and performance improvements

### v1.0.0 - Initial Release
- ✨ Static and Browser download modes
- ✨ SPA support with intelligent route discovery
- ✨ Network idle detection
- ✨ Link rewriting for offline viewing
- ✨ Real-time progress tracking
- ✨ Cross-platform desktop app

## 🔮 Coming Soon

### Planned Features
- [ ] Auto-detect SPA vs static sites
- [ ] Authentication support (login-protected sites)
- [ ] Scheduled/automated downloads
- [ ] Custom wait conditions
- [ ] Download presets/templates
- [ ] Search in history
- [ ] Export history to CSV
- [ ] Dark mode UI theme
- [ ] Download queue (multiple sites)
- [ ] Bandwidth limiting
- [ ] Proxy support

### Community Requests
Want a feature? Open an issue on GitHub!

## 💡 Tips & Tricks

### Getting the Best Results

**For Static Sites:**
- Use Static mode (much faster)
- Increase `maxPages` if needed
- Check the ZIP for completeness

**For SPAs:**
- Always use Browser mode
- Increase `browserWaitTime` for slow sites
- Check browser console for errors
- Some sites may require authentication

**For Large Sites:**
- Increase `maxPages` limit
- Use Browser mode sparingly
- Consider downloading sections separately
- Monitor disk space

### Troubleshooting

**Download Stops Early:**
- Increase `requestTimeout`
- Check internet connection
- Try Browser mode instead

**Incomplete SPA:**
- Increase `browserWaitTime`
- Check if site requires login
- Some routes may be protected

**History Not Saving:**
- Check write permissions for `downloads/` folder
- Ensure disk space available
- Check for corrupted history.json

## 🤝 Contributing

Want to improve the app? Contributions welcome!

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📚 Learn More

- [Project README](README.md)
- [Desktop App Quick Start](website-downloader-app/QUICK-START.md)
- [SPA Technical Details](SPA-IMPROVEMENTS.md)
- [Release Guide](RELEASE-GUIDE.md)

---

**Built with ❤️ using Wails, Go, and modern web technologies**
