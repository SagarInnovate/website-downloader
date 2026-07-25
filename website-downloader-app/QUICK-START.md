# Quick Start Guide

## ✨ You just built a desktop application!

Your executable is ready at:
```
build/bin/website-downloader-app.exe
```

## 🚀 Run the App

Simply double-click the executable, or run:
```bash
.\build\bin\website-downloader-app.exe
```

## 📦 What You Get

- **Single executable** - No backend needed, no complicated setup
- **Zero dependencies** - Everything is bundled
- **Native performance** - Uses your system's WebView
- **Cross-platform** - Can build for Windows, macOS, Linux

## 🔧 Development Mode

Want to make changes? Run in dev mode with hot reload:
```bash
wails dev
```

Changes to the frontend reload instantly. Changes to Go code restart the app automatically.

## 📱 Building for Other Platforms

### macOS
```bash
wails build -platform darwin/universal
```

### Linux
```bash
wails build -platform linux/amd64
```

### All Platforms at Once
```bash
wails build -platform windows/amd64,darwin/universal,linux/amd64
```

## 📝 Next Steps

1. **Customize the UI** - Edit files in `frontend/`
2. **Add features** - Modify `app.go` for new functionality
3. **Configure defaults** - Edit `config.json`
4. **Create installer** - Use `wails build -nsis` (Windows)

## 🎯 Usage

1. Launch the app
2. Enter a website URL
3. Choose mode (Static for regular sites, Browser for SPAs)
4. Click Download
5. Wait for completion
6. Open the downloaded ZIP file

## 💡 Tips

- **Static Mode**: 10-20x faster, works for most websites
- **Browser Mode**: Use for React/Vue/Angular apps
- **Max Pages**: Configurable in config.json (default: 1000)
- **Downloads Location**: `./downloads/` folder (opens automatically)

## 🐛 Troubleshooting

### App won't start
- Make sure Chrome/Chromium is installed (for Browser mode)
- Check if antivirus is blocking the exe

### Downloads not working
- Check internet connection
- Try a different website
- Check console output for errors

## 📚 Learn More

- [Wails Documentation](https://wails.io)
- [Go Documentation](https://golang.org/doc/)
- [Project README](README.md)

---

**You're ready to share this with users! No installation required, just send them the .exe file** 🎉
