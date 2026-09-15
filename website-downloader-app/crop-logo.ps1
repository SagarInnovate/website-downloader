# Crop white borders from logo
Add-Type -AssemblyName System.Drawing

$inputPath = "frontend\src\assets\images\logo.png"
$outputPath = "build\appicon.png"

# Load image
$img = [System.Drawing.Image]::FromFile((Resolve-Path $inputPath))
$bitmap = New-Object System.Drawing.Bitmap($img)

# Find bounds (crop white/transparent borders)
$left = $img.Width
$top = $img.Height
$right = 0
$bottom = 0

for ($y = 0; $y -lt $img.Height; $y++) {
    for ($x = 0; $x -lt $img.Width; $x++) {
        $pixel = $bitmap.GetPixel($x, $y)
        # Check if pixel is not white/transparent
        if ($pixel.A -gt 20 -and ($pixel.R -lt 250 -or $pixel.G -lt 250 -or $pixel.B -lt 250)) {
            if ($x -lt $left) { $left = $x }
            if ($x -gt $right) { $right = $x }
            if ($y -lt $top) { $top = $y }
            if ($y -gt $bottom) { $bottom = $y }
        }
    }
}

# Add small padding
$padding = 20
$left = [Math]::Max(0, $left - $padding)
$top = [Math]::Max(0, $top - $padding)
$right = [Math]::Min($img.Width - 1, $right + $padding)
$bottom = [Math]::Min($img.Height - 1, $bottom + $padding)

$width = $right - $left + 1
$height = $bottom - $top + 1

Write-Host "Original size: $($img.Width)x$($img.Height)"
Write-Host "Cropped size: ${width}x${height}"
Write-Host "Crop bounds: Left=$left, Top=$top, Right=$right, Bottom=$bottom"

# Create cropped image
$rect = New-Object System.Drawing.Rectangle($left, $top, $width, $height)
$cropped = $bitmap.Clone($rect, $bitmap.PixelFormat)

# Save
$cropped.Save($outputPath, [System.Drawing.Imaging.ImageFormat]::Png)

# Cleanup
$cropped.Dispose()
$bitmap.Dispose()
$img.Dispose()

Write-Host "Saved to: $outputPath" -ForegroundColor Green
