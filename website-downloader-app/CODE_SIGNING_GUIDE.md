# Code Signing Guide - Fix "Unverified Publisher" Warning

## Why This Happens
Windows shows "Unverified Publisher" warning because the .exe is not digitally signed with a trusted certificate.

## Solution Options

### Option 1: Purchase Code Signing Certificate (Recommended for Public Release)

**Where to Buy:**
1. **Sectigo/Comodo** (~$70-200/year)
2. **DigiCert** (~$400/year) 
3. **GlobalSign** (~$250/year)
4. **Certum** (~€100/year - cheapest)

**Process:**
1. Purchase certificate from CA
2. Verify your identity (may need business documents)
3. Receive certificate (.pfx file)
4. Sign your executable

**How to Sign:**

```powershell
# Install SignTool (comes with Windows SDK)
# Then sign your executable:

signtool sign /f "your-certificate.pfx" /p "password" /t http://timestamp.digicert.com "website-downloader.exe"
```

### Option 2: Self-Sign (For Testing - Still Shows Warning)

```powershell
# Create self-signed certificate
$cert = New-SelfSignedCertificate -Type CodeSigningCert -Subject "CN=SagarInnovate" -CertStoreLocation "Cert:\CurrentUser\My"

# Export certificate
$pwd = ConvertTo-SecureString -String "YourPassword" -Force -AsPlainText
Export-PfxCertificate -Cert $cert -FilePath "self-signed.pfx" -Password $pwd

# Sign executable
& "C:\Program Files (x86)\Windows Kits\10\bin\10.0.22621.0\x64\signtool.exe" sign /f "self-signed.pfx" /p "YourPassword" "website-downloader.exe"
```

**Note:** Self-signed certs still show warnings until users manually trust them.

### Option 3: Add to wails.json (For Future Builds)

Add to `wails.json`:

```json
{
  "info": {
    "companyName": "SagarInnovate",
    "productName": "Website Downloader",
    "productVersion": "1.0.0",
    "copyright": "Copyright © 2026 SagarInnovate",
    "comments": "Professional desktop app for archiving websites offline"
  },
  "windows": {
    "certificate": {
      "file": "path/to/certificate.pfx",
      "password": "your-password"
    }
  }
}
```

Then Wails will automatically sign during build.

### Option 4: Accept the Warning (Temporary Solution)

For now, you can:
1. Inform users in README that Windows will show a warning
2. Provide instructions to click "More info" → "Run anyway"
3. Explain it's normal for new/unsigned apps
4. Add note that you're working on code signing

## What Users See

**Without Signing:**
```
Windows protected your PC
Microsoft Defender SmartScreen prevented an unrecognized app from starting.
```

**With Signing:**
```
(No warning - app runs directly)
```

## Costs Comparison

| Provider | Annual Cost | Validation Time |
|----------|------------|-----------------|
| Certum | €100 | 1-2 days |
| Sectigo | $180 | 1-3 days |
| DigiCert | $400 | 1-5 days |
| Self-Signed | Free | Immediate (but shows warnings) |

## My Recommendation

**For v1.0.0:**
1. Add disclaimer in README about security warning
2. Show users how to bypass (More info → Run anyway)
3. Most open-source apps don't have signing initially

**For v1.1.0+:**
1. If app gets popular, invest in certificate (~$100-200/year)
2. Sign all future releases
3. Builds trust and professionalism

## Temporary Fix (README)

Add this to your README and release notes:

```markdown
## ⚠️ Windows Security Warning

When first running the app, Windows may show a SmartScreen warning:
"Windows protected your PC - Unrecognized app"

**This is normal for new applications.** To run the app:
1. Click "More info"
2. Click "Run anyway"

We're working on code signing for future releases.
```

## Future: Automatic Signing with GitHub Actions

Once you have a certificate, add it to GitHub Secrets and update the workflow to sign automatically.

---

**Bottom Line:** Code signing costs $100-400/year. For v1.0.0, just inform users about the warning. Get a certificate when the app gains traction!
