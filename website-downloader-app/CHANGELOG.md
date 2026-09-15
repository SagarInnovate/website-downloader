# Changelog

All notable changes to Website Downloader will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-09-15

### 🎉 Initial Public Release

First stable release of Website Downloader - A professional desktop application for archiving websites offline.

### ✨ Added

#### Core Features
- **Two Download Modes**
  - Static Fetch mode for traditional websites (fast, low resource)
  - Browser Render mode for JavaScript-heavy SPAs (full execution)
- **Intelligent SPA Support**
  - Automatic route discovery for single-page applications
  - Network idle detection for complete content capture
  - Smart waiting strategies for dynamic content
- **Real-time Progress Tracking**
  - Pages discovered and downloaded counters
  - Assets downloaded counter
  - Total size tracking
  - Current page being processed
  - Percentage progress bar with animation

#### User Interface
- **Professional Sidebar Navigation**
  - Home view for downloads
  - History view for past downloads
  - Settings view for configuration
  - About view with project information
- **Native Desktop Feel**
  - Smooth scrolling with custom scrollbars
  - Hardware-accelerated animations
  - No text selection on UI elements
  - Proper overflow handling
  - Responsive layout
- **Custom Logo Integration**
  - Branded sidebar header
  - Logo in About section
  - Consistent visual identity

#### Download Management
- **Stop/Cancel Functionality**
  - Cancel active downloads at any time
  - Native confirmation dialog
  - Proper cleanup of resources
- **Smart File Naming**
  - Format: `domain_timestamp.zip`
  - Example: `stripe_com_20260915_143022.zip`
  - Automatic domain extraction from URL
  - Sanitized filenames
- **Native Save Dialog**
  - "Save As..." functionality
  - Choose save location
  - Native Windows file picker
  - Default to downloads folder

#### History & Persistence
- **Download History**
  - Stores last 100 downloads
  - Persistent JSON storage
  - Metadata: URL, mode, pages, assets, size, timestamp
  - Quick access to saved files
- **History Management**
  - View all past downloads
  - Open file location
  - Clear history option
  - Formatted timestamps (relative and absolute)

#### Settings & Configuration
- **Configurable Limits**
  - Maximum pages per download
  - Browser wait time (ms)
  - Page timeout (seconds)
- **Settings UI**
  - Editable configuration
  - Instant save functionality
  - Default values displayed
  - Download folder location

#### Technical Features
- **Asset Handling**
  - CSS files
  - JavaScript files
  - Images (all formats)
  - Fonts
  - Videos and media
- **Link Rewriting**
  - Converts absolute URLs to relative
  - Maintains website structure
  - Preserves functionality offline
- **ZIP Packaging**
  - Automatic compression
  - Preserves folder structure
  - Ready for offline browsing

#### Developer Features
- **Native Dialogs**
  - Error messages
  - Success notifications
  - Confirmation prompts
  - Question dialogs
  - Windows MessageDialog API
- **WebSocket Updates**
  - Real-time progress events
  - Job status tracking
  - Event-driven architecture
- **Modular Architecture**
  - Separate scraper package
  - Clean handler structure
  - Reusable components

### 🛠️ Technical Stack
- Backend: Go 1.18+
- Browser Automation: Chromedp
- Frontend: Vanilla JavaScript + Modern CSS
- Desktop Framework: Wails v2
- Platform: Windows (macOS/Linux coming soon)

### 📦 Distribution
- Single executable file
- No dependencies required
- ~20-40 MB size
- Native Windows integration

### 📚 Documentation
- Comprehensive README
- Quick Start Guide
- Release Guide
- Contributing Guidelines
- License (MIT)

### 🎨 UI/UX
- Modern, professional design
- Purple gradient color scheme (#667eea to #764ba2)
- Smooth animations and transitions
- Responsive layout
- Native-feeling scrollbars
- Accessible interface

### ⚙️ Configuration
- `config.json` for backend settings
- Settings UI for user configuration
- Persistent preferences
- Sensible defaults

### 🔒 Security
- No data collection
- Local processing only
- Secure file handling
- Safe URL validation

### Known Limitations
- Some websites may block automated scraping
- Authentication-required sites need manual handling
- Very large sites (1000+ pages) may take significant time
- Windows-only for initial release

### Future Enhancements (Planned)
- macOS and Linux builds
- Download queue management
- Scheduled downloads
- Export/Import history
- Advanced filtering options
- Proxy support
- Custom user agents
- Cloud storage integration

---

## How to Update

To update to this version:
1. Download the latest release from GitHub
2. Replace your existing executable
3. No migration needed for first release

## Support

- Report bugs: [GitHub Issues](https://github.com/SagarInnovate/website-downloader/issues)
- Feature requests: [GitHub Issues](https://github.com/SagarInnovate/website-downloader/issues)
- Questions: [GitHub Discussions](https://github.com/SagarInnovate/website-downloader/discussions)

---

**Full Changelog**: https://github.com/SagarInnovate/website-downloader/commits/v1.0.0
