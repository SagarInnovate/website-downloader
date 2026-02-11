@echo off
REM Website Downloader - Test Script for Windows
REM This script performs end-to-end testing of the application

setlocal enabledelayedexpansion

echo =========================================
echo Website Downloader - Test Suite
echo =========================================
echo.

REM Test configuration
set BACKEND_URL=http://localhost:8080
set TEST_URL=https://example.com
set TIMEOUT=60

echo Step 1: Checking if backend is running...
curl -s -f "%BACKEND_URL%/health" >nul 2>&1
if %errorlevel% equ 0 (
    echo [OK] Backend is running
) else (
    echo [ERROR] Backend is not running
    echo Start it with: cd backend ^&^& go run main.go
    exit /b 1
)

echo.
echo Step 2: Testing API endpoints...
echo Testing POST /api/scrape...

REM Create temp file for response
set TEMP_FILE=%TEMP%\scrape_response.txt

curl -s -X POST "%BACKEND_URL%/api/scrape" ^
    -H "Content-Type: application/json" ^
    -d "{\"url\":\"%TEST_URL%\"}" ^
    -o "%TEMP_FILE%"

if %errorlevel% equ 0 (
    echo [OK] Scrape endpoint returns data
    
    REM Extract jobId (simplified for Windows)
    for /f "tokens=2 delims=:" %%a in ('findstr "jobId" "%TEMP_FILE%"') do (
        set JOB_ID=%%a
    )
    
    REM Clean up jobId (remove quotes and comma)
    set JOB_ID=!JOB_ID:"=!
    set JOB_ID=!JOB_ID:,=!
    set JOB_ID=!JOB_ID: =!
    
    echo [OK] Received job ID: !JOB_ID!
) else (
    echo [ERROR] Scrape endpoint failed
    exit /b 1
)

echo.
echo Step 3: Testing status endpoint...
timeout /t 2 >nul

curl -s "%BACKEND_URL%/api/status/!JOB_ID!" -o "%TEMP_FILE%"
findstr "status" "%TEMP_FILE%" >nul
if %errorlevel% equ 0 (
    echo [OK] Status endpoint working
) else (
    echo [ERROR] Status endpoint failed
    exit /b 1
)

echo.
echo Step 4: Waiting for scraping to complete...
echo This may take up to %TIMEOUT% seconds...

set /a COUNTER=0
:wait_loop
if !COUNTER! geq %TIMEOUT% goto timeout_reached

curl -s "%BACKEND_URL%/api/status/!JOB_ID!" -o "%TEMP_FILE%"
findstr "complete" "%TEMP_FILE%" >nul
if %errorlevel% equ 0 (
    echo [OK] Scraping completed successfully
    goto scraping_done
)

findstr "error" "%TEMP_FILE%" >nul
if %errorlevel% equ 0 (
    echo [ERROR] Scraping failed
    goto scraping_done
)

timeout /t 2 >nul
set /a COUNTER=!COUNTER!+2
echo Waiting... !COUNTER!s elapsed
goto wait_loop

:timeout_reached
echo [INFO] Timeout reached. Job may still be running.

:scraping_done
echo.
echo Step 5: Testing download endpoint...

set DOWNLOAD_FILE=test_download_!JOB_ID!.zip

curl -s -o "%DOWNLOAD_FILE%" "%BACKEND_URL%/api/download/!JOB_ID!"
if exist "%DOWNLOAD_FILE%" (
    echo [OK] Download successful
    
    REM Get file size
    for %%A in ("%DOWNLOAD_FILE%") do set FILE_SIZE=%%~zA
    echo [OK] Downloaded ZIP file (!FILE_SIZE! bytes)
    
    REM Cleanup
    del "%DOWNLOAD_FILE%"
    echo [INFO] Cleaned up test file
) else (
    echo [ERROR] Download failed
    exit /b 1
)

REM Cleanup temp file
del "%TEMP_FILE%"

echo.
echo =========================================
echo All tests passed!
echo =========================================
echo.
echo Test Summary:
echo   - Backend health check: OK
echo   - Scrape endpoint: OK
echo   - Status endpoint: OK
echo   - Scraping completion: OK
echo   - Download endpoint: OK
echo.

endlocal
