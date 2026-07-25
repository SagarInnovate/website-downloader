# Website Downloader

Download entire websites for offline viewing with intelligent support for both traditional websites and modern Single Page Applications (SPAs).

## 🚀 Two Ways to Use

### Option 1: Desktop App (Recommended for End Users)
**One-click installation. No setup required.**

The Wails desktop application provides a native, user-friendly experience:
- ✅ Single executable file
- ✅ No backend setup needed
- ✅ No Go, Node, or other dependencies required
- ✅ Works on Windows, macOS, and Linux
- ✅ Native system integration

👉 **[Go to Desktop App →](./website-downloader-app/)**

### Option 2: Backend + Frontend (For Developers)
Traditional client-server architecture for maximum flexibility:
- Backend API server (Go)
- Separate frontend (React/Vue/etc.)
- Full REST API access
- WebSocket real-time updates
- Ideal for integration into existing systems

👉 **[Go to Backend API →](./backend/)** | **[Frontend →](./frontend/)**

## ✨ Features

### Core Capabilities
- 🌐 **Download Complete Websites** - All pages, assets, and resources
- ⚡ **Fast Static Mode** - Traditional websites download in seconds
- ⚛️ **SPA Support** - Intelligent crawling for React, Vue, Angular apps
- 🎯 **Smart Route Discovery** - Automatically finds all pages in SPAs
- 🕷️ **Network Idle Detection** - Waits for complete page rendering
- 💾 **Offline Viewing** - Rewritten links work without internet
- 📦 **ZIP Archives** - Everything packaged for easy sharing

### Advanced Features
- **Headless Browser Support** - Full JavaScript execution for dynamic content
- **Cross-Domain Assets** - Downloads CSS, JS, images from CDNs
- **Link Rewriting** - Converts absolute URLs to relative paths
- **Progress Tracking** - Real-time updates via WebSocket
- **Configurable Limits** - Control max pages, timeout, wait times

## 🎯 Use Cases

- **Web Archiving** - Preserve websites for posterity
- **Offline Documentation** - Access docs without internet
- **Testing** - Download production sites for local testing
- **Research** - Save research materials permanently
- **Backup** - Create snapshots of your websites
- **Education** - Study website structure and design offline

## 🔧 Technology Stack

- **Backend**: Go 1.23+ with chromedp for browser automation
- **Frontend**: Modern JavaScript/TypeScript
- **Desktop**: Wails v2 (Go + Web)
- **Browser Automation**: Chrome DevTools Protocol
- **Packaging**: Single executable with embedded assets

## 📊 Performance

| Mode | Speed | Memory | Use Case |
|------|-------|--------|----------|
| **Static** | 0.1-0.5s/page | ~50MB | Traditional websites |
| **Browser** | 2-5s/page | ~200-500MB | SPAs & dynamic content |

## 🛠️ Quick Start

### Desktop App (Easiest)

1. **Download** the pre-built executable for your platform
2. **Double-click** to run
3. **Enter URL** and click Download
4. Done! 🎉

### From Source (Desktop App)

```bash
# Clone repository
git clone <your-repo-url>
cd website-downloader/website-downloader-app

# Build (installs dependencies automatically)
wails build

# Run
./build/bin/website-downloader-app
```

### Backend + Frontend (Developer Setup)

```bash
# Start backend
cd backend
go run main.go

# Start frontend (separate terminal)
cd frontend
npm install
npm run dev
```

## 📚 Documentation

- **[Desktop App README](./website-downloader-app/README.md)** - Full desktop app documentation
- **[Desktop App Quick Start](./website-downloader-app/QUICK-START.md)** - Get started in 5 minutes
- **[API Documentation](./API.md)** - REST API reference
- **[SPA Improvements](./SPA-IMPROVEMENTS.md)** - Technical details on SPA crawling

## 🎨 Screenshots

*Coming soon - Add your app screenshots here*

## 🔍 How It Works

### Static Mode (Traditional Websites)
1. HTTP GET request to starting URL
2. Parse HTML for links and assets
3. Queue discovered pages
4. Download assets in parallel
5. Rewrite links for offline use
6. Package into ZIP

### Browser Mode (SPAs)
1. Launch headless Chrome
2. Load initial page with JavaScript
3. Discover routes from navigation and JS bundles
4. Navigate to each route using History API
5. Wait for network idle (no pending requests)
6. Capture fully rendered HTML
7. Extract and download all assets
8. Rewrite links and package into ZIP

## 🌟 Why This Project?

Most website downloaders fail with modern SPAs because they:
- Don't execute JavaScript
- Miss dynamically loaded content
- Can't discover client-side routes
- Use fixed delays instead of smart waiting

This project solves all of these issues with:
- ✅ Full JavaScript execution via headless Chrome
- ✅ Intelligent route discovery from multiple sources
- ✅ Network-aware waiting (no arbitrary delays)
- ✅ Support for React Router, Vue Router, Angular Router
- ✅ Efficient browser context reuse

## 🤝 Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

[Specify your license here]

## 🙏 Acknowledgments

- [Wails](https://wails.io/) - Amazing Go + Web framework
- [chromedp](https://github.com/chromedp/chromedp) - Chrome DevTools Protocol
- Open source community

## 💬 Support

- 🐛 **Issues**: [GitHub Issues](your-repo/issues)
- 💡 **Discussions**: [GitHub Discussions](your-repo/discussions)
- 📧 **Email**: [your-email]

## 🗺️ Roadmap

- [ ] Auto-detect SPA vs static sites
- [ ] Authentication support for private sites
- [ ] Infinite scroll detection
- [ ] Custom wait conditions
- [ ] Browser extensions
- [ ] Cloud deployment options
- [ ] Scheduled/automated downloads

---

**Made with ❤️ for the open-source community**

*Star ⭐ this repo if you find it useful!*
