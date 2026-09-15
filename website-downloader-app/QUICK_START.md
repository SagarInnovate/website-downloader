# Quick Start Guide

## 🎉 Welcome to Website Downloader!

This guide will help you start downloading websites in less than 5 minutes.

## 📥 Installation

### Option 1: Download Pre-built (Recommended)
1. Download `website-downloader-app.exe` from the latest release
2. Double-click to run - that's it! No installation needed.

### Option 2: Build from Source
```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Clone and run
git clone <repository-url>
cd website-downloader-app
wails dev
```

## 🚀 Your First Download

### 1. Launch the App
- Double-click the executable
- The app window opens immediately

### 2. Enter a Website URL
- Type or paste any website URL (e.g., `https://example.com`)
- Must include `https://` or `http://`

### 3. Choose Download Mode

**⚡ Static Fetch** - Choose this for:
- Blogs and news sites
- Documentation websites
- Traditional HTML websites
- **Fastest option!**

**🌐 Browser Render** - Choose this for:
- React, Vue, Angular apps
- Single Page Applications (SPAs)
- JavaScript-heavy websites
- Websites that load content dynamically

### 4. Start Download
- Click "Start Archiving"
- Watch the progress in real-time
- See pages, assets, and size statistics

### 5. Save Your Download
- When complete, click "Save As..."
- Choose where to save your ZIP file
- File automatically named: `domain_timestamp.zip`

## 💡 Tips for Best Results

### For Static Websites
✅ Use **Static Fetch** mode  
✅ Much faster  
✅ Uses less resources  
✅ Perfect for most sites  

### For JavaScript Apps
✅ Use **Browser Render** mode  
✅ Waits for content to load  
✅ Discovers all SPA routes  
✅ Handles dynamic content  

### General Tips
- Start with smaller sites to test
- Use **Stop** button if download takes too long
- Check History to access previous downloads
- Settings let you customize timeouts

## 📂 Where Are My Downloads?

By default, downloads are saved to:
- `./downloads` folder (next to the app)
- Or wherever you choose in "Save As..." dialog

Files are named: `example_com_20260915_143022.zip`
- Domain name from URL
- Date and time of download
- Easy to identify later

## 🎯 Navigation

The app has 4 main sections in the sidebar:

### 🏠 Home
- Download new websites
- See real-time progress
- Stop active downloads

### 📜 History
- View all past downloads
- Quick access to files
- See download statistics
- Clear history if needed

### ⚙️ Settings
- Maximum pages to download
- Browser wait time
- Page timeout settings
- Downloads folder location

### ℹ️ About
- App version and info
- Feature list
- Technology details

## ⚠️ Common Issues

### "Invalid URL" Error
- Make sure to include `https://` or `http://`
- Example: `https://example.com` not `example.com`

### Download Fails
- Check internet connection
- Try the other mode (Static vs Browser)
- Some sites may block scrapers
- Increase timeout in Settings

### Slow Downloads
- Large sites take time (normal)
- Check number of pages being downloaded
- Use Static mode if possible (faster)
- Adjust max pages in Settings

### "Save As" Cancelled
- You cancelled the save dialog
- Download is still available
- Try "Save As..." again

## 🆘 Need Help?

1. **Check the README** - Full documentation available
2. **Visit History** - See what worked before
3. **Try Different Mode** - Switch between Static/Browser
4. **Adjust Settings** - Increase timeouts for slow sites
5. **Open an Issue** - Report bugs on GitHub

## 🎓 Examples

### Download a Blog
```
URL: https://blog.example.com
Mode: Static Fetch ⚡
Result: Fast download, all posts included
```

### Download React Documentation
```
URL: https://react.dev
Mode: Browser Render 🌐
Result: All docs pages, properly rendered
```

### Download Company Website
```
URL: https://company.com
Mode: Static Fetch ⚡
Result: Homepage and linked pages
```

## 🔥 Pro Tips

1. **Test with Small Sites First** - Get familiar with the app
2. **Use History** - Quickly re-download updated sites
3. **Adjust Max Pages** - Prevent accidentally downloading huge sites
4. **Static vs Browser** - Start with Static, switch to Browser if needed
5. **Bookmark Frequently Used Sites** - Keep URLs handy

## ✅ Success Indicators

You'll know it worked when:
- ✓ Progress reaches 100%
- ✓ "Archive Ready" screen appears
- ✓ ZIP file is saved to your chosen location
- ✓ Download appears in History
- ✓ File opens and shows website structure

## 🎊 You're All Set!

Start downloading websites and enjoy offline access to your favorite content!

---

**Questions?** Check the full README.md for detailed information.  
**Issues?** Open a GitHub issue with details.  
**Ideas?** We welcome feature suggestions!

Happy archiving! 🚀
