export const appTemplate = `
    <div class="app-layout">
        <!-- Sidebar -->
        <aside class="sidebar">
            <div class="sidebar-header">
                <div class="brand-icon">
                    <img src="/src/assets/images/logo.png" alt="Website Downloader" class="logo-image" />
                </div>
                <h2 class="sidebar-title">Website Downloader</h2>
            </div>
            
            <nav class="sidebar-nav">
                <button class="nav-item active" data-view="home">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                        <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
                        <polyline points="9 22 9 12 15 12 15 22" />
                    </svg>
                    Home
                </button>
                
                <button class="nav-item" data-view="history">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                        <circle cx="12" cy="12" r="10" />
                        <polyline points="12 6 12 12 16 14" />
                    </svg>
                    History
                </button>
                
                <button class="nav-item" data-view="settings">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                        <circle cx="12" cy="12" r="3" />
                        <path d="M12 1v6m0 6v6m5.2-13.2l-4.2 4.2m-2 2l-4.2 4.2M23 12h-6m-6 0H5m13.2 5.2l-4.2-4.2m-2-2l-4.2-4.2" />
                    </svg>
                    Settings
                </button>
                
                <button class="nav-item" data-view="about">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                        <circle cx="12" cy="12" r="10" />
                        <path d="M12 16v-4m0-4h.01" />
                    </svg>
                    About
                </button>
            </nav>
            
            <div class="sidebar-footer">
                <a href="https://github.com/SagarInnovate/website-downloader" target="_blank" class="github-star-btn">
                    <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
                        <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/>
                    </svg>
                    <span>Star on GitHub</span>
                </a>
                <p class="footer-text">
                    <strong>v1.0.0</strong><br />
                    Built with Wails
                </p>
            </div>
        </aside>
        
        <!-- Main Content -->
        <div class="main-content">
            <div class="container">
                <!-- Home View -->
                <div id="homeView" class="view-content">
                    <header class="header">
                        <div class="hero-text">
                            <h1 class="title-main">
                                Download <span class="text-primary">Any Website</span>
                            </h1>
                            <p class="subtitle">
                                Production-ready offline scraping with structure preservation.
                                <br />
                                Built by <strong>sagarinnovate</strong>.
                            </p>
                        </div>
                    </header>

                    <main>
                        <!-- Main Form -->
                        <div id="mainForm" class="card">
                            <form class="url-form" id="downloadForm">
                                <div class="input-group">
                                    <label class="input-label">
                                        Target Website
                                        <span class="input-label-hint">Full URL including https://</span>
                                    </label>
                                    <input
                                        type="text"
                                        id="url"
                                        class="url-input"
                                        placeholder="e.g. https://stripe.com/docs"
                                        autocomplete="off"
                                        autofocus
                                    />
                                </div>


                                <div class="input-group">
                                    <label class="input-label">Extraction Method</label>
                                    <div class="mode-options">
                                        <label class="mode-card">
                                            <input type="radio" name="mode" value="static" checked />
                                            <div class="mode-content">
                                                <span class="mode-title">
                                                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
                                                        <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z" />
                                                    </svg>
                                                    Static Fetch
                                                </span>
                                                <span class="mode-desc">Lightning fast. Ideal for blogs, documentation, and traditional server-rendered sites.</span>
                                            </div>
                                        </label>
                                        <label class="mode-card">
                                            <input type="radio" name="mode" value="browser" />
                                            <div class="mode-content">
                                                <span class="mode-title">
                                                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
                                                        <circle cx="12" cy="12" r="10" />
                                                        <line x1="2" y1="12" x2="22" y2="12" />
                                                        <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z" />
                                                    </svg>
                                                    Browser Render
                                                </span>
                                                <span class="mode-desc">Full simulation. Executes JavaScript for SPAs like React, Vue, and Angular applications.</span>
                                            </div>
                                        </label>
                                    </div>
                                </div>

                                <button type="submit" class="submit-btn" id="submitBtn">
                                    Start Archiving
                                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
                                        <path d="M5 12h14M12 5l7 7-7 7" />
                                    </svg>
                                </button>
                            </form>
                        </div>


                        <!-- Progress View -->
                        <div id="progressView" class="card hidden">
                            <div class="progress-header">
                                <h3>
                                    <div class="spinner-pulse"></div>
                                    Archiving in Progress
                                </h3>
                                <button class="stop-btn" id="stopBtn">
                                    ⏹️ Stop
                                </button>
                            </div>

                            <div class="progress-bar-container">
                                <div class="progress-bar-bg">
                                    <div id="progressFill" class="progress-bar-fill" style="width: 0%"></div>
                                </div>
                                <div class="progress-meta">
                                    <span>Processing resources...</span>
                                    <span class="progress-percent" id="progressPercent">0%</span>
                                </div>
                            </div>

                            <div class="stats-grid">
                                <div class="stat-item">
                                    <div class="stat-label">Pages Found</div>
                                    <div class="stat-value" id="pagesFound">0</div>
                                </div>
                                <div class="stat-item">
                                    <div class="stat-label">Downloaded</div>
                                    <div class="stat-value" id="pagesDownloaded">0</div>
                                </div>
                                <div class="stat-item">
                                    <div class="stat-label">Assets</div>
                                    <div class="stat-value" id="assetsDownloaded">0</div>
                                </div>
                                <div class="stat-item">
                                    <div class="stat-label">Size</div>
                                    <div class="stat-value" id="totalSize">0 B</div>
                                </div>
                            </div>

                            <div id="currentStatus" class="current-status hidden">
                                <svg class="status-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
                                    <path d="M12 2v4m0 12v4M4.93 4.93l2.83 2.83m8.48 8.48l2.83 2.83M2 12h4m12 0h4M4.93 19.07l2.83-2.83m8.48-8.48l2.83-2.83" />
                                </svg>
                                <span class="status-text" id="statusText"></span>
                            </div>
                        </div>


                        <!-- Success View -->
                        <div id="successView" class="card hidden">
                            <div class="success-icon-wrapper">
                                <div class="success-ring"></div>
                                <div class="success-icon">
                                    <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="3">
                                        <path d="M20 6L9 17l-5-5" strokeLinecap="round" strokeLinejoin="round" />
                                    </svg>
                                </div>
                            </div>

                            <h2 class="success-title">Archive Ready</h2>
                            <p class="success-message">
                                Your website has been successfully captured. All assets, scripts, and styles are packaged and ready for offline use.
                            </p>

                            <div class="button-group">
                                <button id="saveAsBtn" class="download-btn">
                                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" width="20" height="20" strokeWidth="2">
                                        <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z" />
                                        <polyline points="17 21 17 13 7 13 7 21" />
                                        <polyline points="7 3 7 8 15 8" />
                                    </svg>
                                    Save As...
                                </button>

                                <button id="startNewBtn" class="new-download-btn">
                                    Archive Another
                                </button>
                            </div>
                        </div>

                        <!-- Error View -->
                        <div id="errorView" class="card error-screen hidden">
                            <div class="error-icon-wrapper">
                                <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                    <circle cx="12" cy="12" r="10" />
                                    <line x1="12" y1="8" x2="12" y2="12" />
                                    <line x1="12" y1="16" x2="12.01" y2="16" />
                                </svg>
                            </div>
                            <h2 class="error-title">Download Failed</h2>
                            <p class="error-message" id="errorMessage"></p>
                            <button id="tryAgainBtn" class="new-download-btn">
                                Try Again
                            </button>
                        </div>
                    </main>
                </div>


                <!-- History View -->
                <div id="historyView" class="view-content hidden">
                    <div class="page-header">
                        <h1 class="title-main">📜 Download History</h1>
                        <p class="subtitle">View and manage your previous downloads</p>
                    </div>
                    
                    <div class="card">
                        <div class="history-actions-bar">
                            <button class="btn-small" id="refreshHistoryBtn">
                                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                    <polyline points="23 4 23 10 17 10" />
                                    <path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10" />
                                </svg>
                                Refresh
                            </button>
                            <button class="btn-small btn-danger" id="clearHistoryBtn">
                                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                    <polyline points="3 6 5 6 21 6" />
                                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
                                </svg>
                                Clear All
                            </button>
                        </div>
                        
                        <div id="historyList" class="history-list">
                            <div class="empty-state">No downloads yet</div>
                        </div>
                    </div>
                </div>


                <!-- Settings View -->
                <div id="settingsView" class="view-content hidden">
                    <div class="page-header">
                        <h1 class="title-main">⚙️ Settings</h1>
                        <p class="subtitle">Configure your download preferences</p>
                    </div>
                    
                    <div class="card">
                        <div class="settings-section">
                            <h3 class="section-title">Scraping Limits</h3>
                            
                            <div class="setting-item">
                                <label class="setting-label">
                                    <span>Maximum Pages</span>
                                    <span class="setting-hint">Maximum number of pages to download per job</span>
                                </label>
                                <input type="number" id="maxPages" class="setting-input" value="500" min="1" max="10000" />
                            </div>
                        </div>
                        
                        <div class="settings-section">
                            <h3 class="section-title">Browser Mode Settings</h3>
                            
                            <div class="setting-item">
                                <label class="setting-label">
                                    <span>Wait Time (ms)</span>
                                    <span class="setting-hint">Time to wait for JavaScript execution after page load</span>
                                </label>
                                <input type="number" id="browserWaitTime" class="setting-input" value="5000" min="1000" max="30000" step="1000" />
                            </div>
                            
                            <div class="setting-item">
                                <label class="setting-label">
                                    <span>Page Timeout (seconds)</span>
                                    <span class="setting-hint">Maximum time to wait for a page to load</span>
                                </label>
                                <input type="number" id="browserTimeout" class="setting-input" value="60" min="10" max="300" step="5" />
                            </div>
                        </div>
                        
                        <div class="settings-section">
                            <h3 class="section-title">Storage</h3>
                            
                            <div class="setting-item">
                                <label class="setting-label">
                                    <span>Downloads Folder</span>
                                    <span class="setting-hint">Default location for downloaded websites</span>
                                </label>
                                <div class="setting-path">
                                    <input type="text" id="downloadPath" class="setting-input" value="./downloads" readonly />
                                    <button class="btn-small" id="openFolderBtn">
                                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                            <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z" />
                                        </svg>
                                        Open
                                    </button>
                                </div>
                            </div>
                        </div>
                        
                        <button id="saveSettingsBtn" class="submit-btn">
                            💾 Save Settings
                        </button>
                    </div>
                </div>


                <!-- About View -->
                <div id="aboutView" class="view-content hidden">
                    <div class="about-content">
                        <div class="about-logo">
                            <img src="/src/assets/images/logo.png" alt="Website Downloader" class="logo-image-large" />
                        </div>
                        
                        <h1 class="about-title">Website Downloader</h1>
                        <p class="about-version">Version 1.0.0</p>
                        
                        <div class="card">
                            <div class="about-section">
                                <h3>🎯 What It Does</h3>
                                <p>
                                    Website Downloader is a professional tool for archiving entire websites for offline use. 
                                    It captures HTML, CSS, JavaScript, images, and all other assets while preserving the 
                                    original structure and functionality.
                                </p>
                            </div>
                            
                            <div class="about-section">
                                <h3>✨ Features</h3>
                                <ul class="about-list">
                                    <li>⚡ Two download modes: Static (fast) and Browser (JavaScript-enabled)</li>
                                    <li>🎯 Intelligent SPA support with route discovery</li>
                                    <li>📦 Automatic ZIP packaging with smart naming</li>
                                    <li>📊 Real-time progress tracking</li>
                                    <li>🔄 Download history with persistent storage</li>
                                    <li>⏹️ Cancel downloads at any time</li>
                                    <li>💾 Native file dialogs and Windows integration</li>
                                </ul>
                            </div>
                            
                            <div class="about-section">
                                <h3>🛠️ Technology</h3>
                                <ul class="about-list">
                                    <li><strong>Backend:</strong> Go with Chromedp for browser automation</li>
                                    <li><strong>Frontend:</strong> Vanilla JavaScript with modern CSS</li>
                                    <li><strong>Framework:</strong> Wails v2 for native desktop experience</li>
                                    <li><strong>Platform:</strong> Windows, macOS, Linux</li>
                                </ul>
                            </div>
                            
                            <div class="about-section">
                                <h3>👨‍💻 Developer</h3>
                                <p>
                                    Created by <strong>sagarinnovate</strong> as an open-source project.
                                    Contributions and feedback are welcome!
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
`;
