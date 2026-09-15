# Release Guide for Website Downloader

## 🚀 Building for Production

### Windows Build

1. **Open PowerShell in the app directory**
   ```powershell
   cd website-downloader-app
   ```

2. **Run the build script**
   ```powershell
   .\build.ps1
   ```
   
   OR manually:
   ```powershell
   wails build -platform windows/amd64
   ```

3. **Find your executable**
   - Location: `build\bin\website-downloader-app.exe`
   - This is a single executable - no dependencies needed!

### macOS Build

```bash
wails build -platform darwin/universal
```
Output: `build/bin/website-downloader-app.app`

### Linux Build

```bash
wails build -platform linux/amd64
```
Output: `build/bin/website-downloader-app`

## 📦 What Gets Built

- **Single executable file** - No installation required
- **All dependencies bundled** - No Go, Node.js, or Chrome needed
- **Native OS integration** - Uses system dialogs and file pickers
- **Small size** - Typically 20-40 MB depending on platform

## 🎯 Distribution

### For Windows:

1. Build the `.exe` file
2. Optionally create an installer using:
   - Inno Setup
   - NSIS
   - WiX Toolset
3. Or just distribute the `.exe` directly - it works standalone!

### For macOS:

1. Build the `.app` bundle
2. Sign and notarize for distribution (optional)
3. Create a `.dmg` installer (optional)

### For Linux:

1. Build the binary
2. Create `.deb` or `.rpm` packages (optional)
3. Or distribute the binary directly

## ✅ Pre-Release Checklist

- [ ] Test on target platform
- [ ] Verify all features work:
  - [ ] Static mode downloads
  - [ ] Browser mode downloads
  - [ ] SPA detection and routing
  - [ ] Stop/Cancel functionality
  - [ ] History tracking
  - [ ] Save As dialog
  - [ ] Settings persistence
  - [ ] Native dialogs
- [ ] Check file naming (domain_timestamp.zip)
- [ ] Test with various websites
- [ ] Verify downloads folder creation
- [ ] Test history clear function
- [ ] Check About page information

## 📝 Version Update Process

1. **Update version in files:**
   - `wails.json` - name and outputfilename
   - `README.md` - badge and version references
   - `frontend/src/template.js` - About view version
   - `package.json` - version field

2. **Update CHANGELOG.md** (create if doesn't exist)

3. **Commit version changes:**
   ```bash
   git add .
   git commit -m "chore: Bump version to X.Y.Z"
   git tag vX.Y.Z
   ```

4. **Build and test**

5. **Create GitHub Release:**
   - Tag: `vX.Y.Z`
   - Title: `Website Downloader vX.Y.Z`
   - Description: Include changelog and notable features
   - Attach built executables

## 🔍 Testing Checklist

### Basic Functionality
- [ ] App starts without errors
- [ ] UI loads correctly
- [ ] All navigation works (Home, History, Settings, About)
- [ ] Logo displays properly

### Download Testing
- [ ] Static mode: Simple website (blog)
- [ ] Static mode: Documentation site
- [ ] Browser mode: React SPA
- [ ] Browser mode: Vue SPA
- [ ] Large website (100+ pages)
- [ ] Website with many assets

### Feature Testing
- [ ] Stop button cancels download
- [ ] Progress updates in real-time
- [ ] ZIP file created with correct name
- [ ] Save As dialog works
- [ ] History saves correctly
- [ ] History shows correct metadata
- [ ] Open file location works
- [ ] Clear history works
- [ ] Settings values load
- [ ] Native dialogs appear (not browser alerts)

### Edge Cases
- [ ] Invalid URL handling
- [ ] Network error handling
- [ ] Timeout scenarios
- [ ] Very slow websites
- [ ] Websites with authentication
- [ ] Websites that block scrapers

## 🐛 Known Issues to Document

- Some websites may block automated scraping
- Very large websites (1000+ pages) may take significant time
- Authentication-required sites need manual login first
- Some dynamic content may not render properly

## 📈 Post-Release

1. **Monitor for issues:**
   - Check GitHub issues
   - User feedback
   - Crash reports

2. **Document fixes needed**

3. **Plan next release:**
   - Feature requests
   - Bug fixes
   - Performance improvements

## 🎁 Release Assets to Include

For each platform release:

1. **Executable/Binary**
   - Windows: `website-downloader-app.exe`
   - macOS: `website-downloader-app.app` (in .dmg)
   - Linux: `website-downloader-app` (binary or .deb/.rpm)

2. **README**
   - Quick start guide
   - System requirements
   - Usage instructions

3. **LICENSE**
   - MIT License file

4. **CHANGELOG**
   - What's new in this version
   - Bug fixes
   - Breaking changes

## 💡 Tips

- **Test on clean machines** - Verify it works without dev tools
- **Include screenshots** - Show the UI in release notes
- **Write clear instructions** - Not everyone is technical
- **Version consistently** - Use semantic versioning (MAJOR.MINOR.PATCH)
- **Sign your binaries** - Prevents security warnings
- **Auto-update** - Consider adding in future releases

## 🔗 Useful Links

- [Wails Documentation](https://wails.io/)
- [Semantic Versioning](https://semver.org/)
- [GitHub Releases Guide](https://docs.github.com/en/repositories/releasing-projects-on-github)

---

**Current Version:** 1.0.0  
**Last Updated:** September 15, 2026  
**Developer:** sagarinnovate
