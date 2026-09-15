# Fix App Icon - Show Custom Logo Instead of Wails Logo

## Problem
The app is currently showing the default Wails "W" logo in the taskbar instead of your custom logo.

## Solution

### Step 1: Update wails.json

Edit `wails.json` and replace its contents with:

```json
{
  "$schema": "https://wails.io/schemas/config.v2.json",
  "name": "Website Downloader",
  "outputfilename": "website-downloader",
  "frontend:install": "npm install",
  "frontend:build": "npm run build",
  "frontend:dev:watcher": "npm run dev",
  "frontend:dev:serverUrl": "auto",
  "author": {
    "name": "Sagar Shinde",
    "email": "sagarinnovate@gmail.com"
  },
  "info": {
    "companyName": "SagarInnovate",
    "productName": "Website Downloader",
    "productVersion": "1.0.0",
    "copyright": "Copyright © 2026 SagarInnovate",
    "comments": "Professional desktop app for archiving websites offline"
  }
}
```

### Step 2: Ensure Icon Files Exist

Make sure these files exist:
- `build/appicon.png` (already exists - your logo)
- `frontend/src/assets/images/logo.png` (already exists)

### Step 3: Generate Windows Icon

Run this command to let Wails generate the icon:

```powershell
cd website-downloader-app
wails build
```

Wails will automatically use `build/appicon.png` to generate Windows `.ico` files.

### Step 4: For Manual Icon Generation (Optional)

If you want to manually create icons:

1. **Windows Icon (.ico)**
   - Use an online converter or tool like ImageMagick
   - Convert `logo.png` to `appicon.ico`
   - Place in `build/windows/` directory

2. **macOS Icon (.icns)**
   - Use `iconutil` or online converter
   - Convert `logo.png` to `appicon.icns`
   - Place in `build/darwin/` directory

### Step 5: Clean Build

For a fresh build with new icon:

```powershell
# Remove old build
Remove-Item -Recurse -Force build/bin

# Rebuild
wails build
```

### Step 6: Verify

After building:
1. Run the new executable
2. Check taskbar - should show your custom logo
3. Check window title bar icon
4. Check Alt+Tab icon

## Quick Fix Commands

```powershell
# Navigate to app directory
cd website-downloader-app

# Update wails.json manually (open in editor)
notepad wails.json

# Clean and rebuild
Remove-Item -Recurse -Force build/bin -ErrorAction SilentlyContinue
wails build
```

## Why This Happens

Wails uses `build/appicon.png` as the source for all platform-specific icons. If the icon looks wrong:
1. The PNG might be too small (should be at least 512x512)
2. The build cache might be stale
3. The wails.json needs proper info section

## Expected Result

After fixing:
- Taskbar shows your custom logo ✅
- Window icon shows your custom logo ✅
- Alt+Tab shows your custom logo ✅
- About dialog shows proper app name ✅

---

**Note:** The icon is already copied to `build/appicon.png` from `frontend/src/assets/images/logo.png`.
Just update `wails.json` and rebuild!
