# Release Guide - How to Share Your Desktop App

## 🎯 Overview

Your Wails desktop app is **ready to share with users** who need zero technical knowledge. No installation, no dependencies, just download and run!

## 📦 What You Have

After running `wails build`, you get:
- **Windows**: `website-downloader-app.exe` (~25-40MB)
- **macOS**: `website-downloader-app.app` (application bundle)
- **Linux**: `website-downloader-app` (binary)

## 🚀 Distribution Methods

### Method 1: Direct Download (Simplest)

1. **Build the app:**
   ```bash
   cd website-downloader-app
   wails build
   ```

2. **Upload to GitHub Releases:**
   - Go to your GitHub repo
   - Click "Releases" → "Create new release"
   - Upload `build/bin/website-downloader-app.exe`
   - Add release notes
   - Publish

3. **Users download and run** - That's it!

### Method 2: Installer (More Professional)

#### Windows Installer (NSIS)

```bash
wails build -nsis
```

Creates an installer in `build/bin/`:
- Professional installation wizard
- Start menu shortcuts
- Uninstaller
- File associations (optional)

#### macOS DMG

```bash
wails build -platform darwin/universal -dmg
```

Creates a `.dmg` file users can drag to Applications folder.

### Method 3: Cross-Platform Build

Build for all platforms at once:

```bash
wails build -platform windows/amd64,darwin/universal,linux/amd64
```

Creates executables for:
- Windows (64-bit)
- macOS (Intel + Apple Silicon)
- Linux (64-bit)

## 📝 Release Checklist

Before releasing:

- [ ] Test the executable on a clean machine
- [ ] Test with both Static and Browser modes
- [ ] Verify Chrome/Chromium detection works
- [ ] Check downloads folder creation
- [ ] Test with real websites (not just localhost)
- [ ] Update version in `wails.json`
- [ ] Write release notes
- [ ] Add screenshots to README
- [ ] Test on all target platforms

## 📋 Release Notes Template

```markdown
## Website Downloader v1.0.0

### Features
- Download entire websites for offline viewing
- Support for modern SPAs (React, Vue, Angular)
- Fast static mode for traditional websites
- Intelligent route discovery
- Network idle detection for complete rendering

### Requirements
- Windows 10/11 (or macOS 10.13+, or Linux)
- Chrome/Chromium browser (for Browser mode only)
- Internet connection

### Download
- [Windows](link-to-exe)
- [macOS](link-to-app)
- [Linux](link-to-binary)

### How to Use
1. Download the file for your operating system
2. Double-click to run (no installation needed)
3. Enter a website URL
4. Choose mode (Static/Browser)
5. Click Download
6. Wait for completion
7. Find your files in the downloads folder

### Known Issues
- None currently

### Notes
- First run may be slow while initializing
- Browser mode requires Chrome/Chromium installed
```

## 🔐 Code Signing (Optional but Recommended)

### Windows
```bash
# Requires code signing certificate
wails build -nsis -sign
```

### macOS
```bash
# Requires Apple Developer account
wails build -platform darwin/universal -codesign
```

**Benefits:**
- Removes "Unknown Publisher" warnings
- Users trust the app more
- Required for distribution outside app stores

## 🌐 Web Distribution

### GitHub Releases (Recommended)
- Free hosting
- Automatic version tracking
- Download statistics
- Release notes built-in

### Your Own Website
```html
<a href="/downloads/website-downloader-app.exe">
  Download for Windows
</a>
```

### App Stores (Advanced)
- Windows Store (requires Microsoft Developer account)
- Mac App Store (requires Apple Developer account)
- Snapcraft (Linux)

## 📊 File Sizes

Typical sizes after compression:
- **Windows .exe**: 25-40MB
- **macOS .app**: 30-50MB
- **Linux binary**: 25-40MB

With installer:
- **Windows .exe with NSIS**: 30-45MB
- **macOS .dmg**: 35-55MB

## 🎨 Branding (Optional)

### Update App Icon

1. **Replace icons:**
   ```
   build/appicon.png          # Main icon (512x512)
   build/windows/icon.ico     # Windows icon
   build/darwin/icon.icns     # macOS icon (create with iconutil)
   ```

2. **Rebuild:**
   ```bash
   wails build
   ```

### Update App Name

Edit `wails.json`:
```json
{
  "name": "Website Downloader",
  "info": {
    "productName": "Website Downloader",
    "comments": "Download websites for offline viewing"
  }
}
```

## 📱 Platform-Specific Notes

### Windows
- Works on Windows 10/11
- WebView2 automatically downloaded if needed
- Defender SmartScreen may show warning (code signing fixes this)

### macOS
- Universal binary supports Intel + Apple Silicon
- Gatekeeper may block unsigned apps
  - Users: Right-click → Open → Allow
  - Or code sign the app

### Linux
- Works on most distributions
- Requires WebKit2GTK
- Users may need: `sudo apt install webkit2gtk-4.0`

## 🎯 Marketing Tips

### Where to Share
- GitHub Releases
- Product Hunt
- Reddit (r/opensource, r/webdev)
- Hacker News
- Dev.to / Hashnode
- Twitter / X
- LinkedIn

### What to Highlight
- "No installation required - just download and run"
- "Works with modern SPAs unlike other tools"
- "100% free and open source"
- "Cross-platform support"
- "Privacy-focused - everything stays local"

### Screenshots to Include
1. Main interface showing URL input
2. Download in progress with progress bar
3. Completed download with file browser
4. Example of downloaded website working offline

## 📈 Version Management

### Semantic Versioning
- `1.0.0` - Major release
- `1.1.0` - New features
- `1.1.1` - Bug fixes

### Update Process
1. Update version in `wails.json`
2. Update CHANGELOG.md
3. Build: `wails build`
4. Test thoroughly
5. Create git tag: `git tag v1.0.0`
6. Push: `git push --tags`
7. Create GitHub release
8. Upload binaries

## 🔄 Auto-Updates (Future Enhancement)

Consider implementing auto-update functionality:
- Check for updates on launch
- Download and install updates
- Notify users of new features

Libraries:
- [go-update](https://github.com/inconshreveable/go-update)
- [wails-updater](https://github.com/wailsapp/wails/discussions/1726)

## 💡 Next Steps

1. **Build and test** the application
2. **Create release** on GitHub
3. **Write announcement** post
4. **Share** on social media and communities
5. **Gather feedback** from users
6. **Iterate** and improve

## 🎉 You're Ready!

Your app is production-ready. Users can now:
- Download one file
- Double-click to run
- Start downloading websites
- No setup, no dependencies, no hassle

**That's the power of Wails!** 🚀

---

Need help? Check:
- [Wails Documentation](https://wails.io/docs/guides/distributing)
- [GitHub Releases Guide](https://docs.github.com/en/repositories/releasing-projects-on-github)
- [Code Signing Guide](https://wails.io/docs/guides/signing)
