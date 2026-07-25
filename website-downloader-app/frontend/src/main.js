import './style.css';
import './app.css';

import {StartScrape, GetJobStatus, OpenDownloadFolder, CancelJob, GetHistory, OpenFileLocation, ClearHistory} from '../wailsjs/go/main/App';
import {EventsOn} from '../wailsjs/runtime/runtime';

let currentJobId = null;
let downloadHistory = [];

document.querySelector('#app').innerHTML = `
    <div class="app">
        <div class="container">
            <header class="header">
                <div class="brand-bar">
                    <div class="brand">
                        <div class="brand-icon">
                            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
                                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z" />
                            </svg>
                        </div>
                        <span class="brand-name">Website Downloader</span>
                    </div>
                    <button class="history-btn" onclick="toggleHistory()">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                            <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
                            <polyline points="9 22 9 12 15 12 15 22" />
                        </svg>
                        History
                    </button>
                </div>

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
                <!-- History Panel -->
                <div id="historyPanel" class="card history-panel hidden">
                    <div class="history-header">
                        <h3>📜 Download History</h3>
                        <div class="history-actions">
                            <button class="btn-small" onclick="refreshHistory()">🔄 Refresh</button>
                            <button class="btn-small" onclick="confirmClearHistory()">🗑️ Clear</button>
                            <button class="btn-small" onclick="toggleHistory()">✖ Close</button>
                        </div>
                    </div>
                    <div id="historyList" class="history-list">
                        <div class="empty-state">No downloads yet</div>
                    </div>
                </div>

                <!-- Main Form -->
                <div id="mainForm" class="card">
                    <form class="url-form" onsubmit="startDownload(event)">
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
                        <button class="stop-btn" onclick="stopDownload()">
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
                        <button onclick="openFolder()" class="download-btn">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" width="20" height="20" strokeWidth="2">
                                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4M7 10l5 5 5-5M12 15V3" strokeLinecap="round" strokeLinejoin="round" />
                            </svg>
                            Open Downloads
                        </button>

                        <button onclick="startNew()" class="new-download-btn">
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
                    <button onclick="startNew()" class="new-download-btn">
                        Try Again
                    </button>
                </div>
            </main>
        </div>
    </div>
`;

// Start download
window.startDownload = async function(event) {
    event.preventDefault();
    
    const url = document.getElementById('url').value.trim();
    const mode = document.querySelector('input[name="mode"]:checked').value;

    if (!url) {
        alert('Please enter a URL');
        return;
    }

    // Validate URL
    try {
        new URL(url);
    } catch (e) {
        alert('Please enter a valid URL (e.g., https://example.com)');
        return;
    }

    // Show progress view
    document.getElementById('mainForm').classList.add('hidden');
    document.getElementById('historyPanel').classList.add('hidden');
    document.getElementById('progressView').classList.remove('hidden');
    document.getElementById('successView').classList.add('hidden');
    document.getElementById('errorView').classList.add('hidden');

    try {
        const response = await StartScrape(url, mode);
        currentJobId = response.jobId;

        // Listen for progress updates
        EventsOn('progress:' + currentJobId, (update) => {
            updateProgress(update);
        });

    } catch (error) {
        console.error('Error starting download:', error);
        showError('Failed to start download: ' + error);
    }
};

// Update progress
function updateProgress(update) {
    const progress = update.progress || 0;
    const pagesDiscovered = update.pagesDiscovered || 0;
    const pagesDownloaded = update.pagesDownloaded || 0;
    const assetsDownloaded = update.assetsDownloaded || 0;
    const totalSize = update.totalSize || 0;
    const currentPage = update.currentPage || '';

    // Update progress bar
    document.getElementById('progressFill').style.width = progress + '%';
    document.getElementById('progressPercent').textContent = Math.round(progress) + '%';

    // Update stats
    document.getElementById('pagesFound').textContent = pagesDiscovered;
    document.getElementById('pagesDownloaded').textContent = pagesDownloaded;
    document.getElementById('assetsDownloaded').textContent = assetsDownloaded;
    document.getElementById('totalSize').textContent = formatBytes(totalSize);

    // Update current page
    if (currentPage) {
        document.getElementById('currentStatus').classList.remove('hidden');
        document.getElementById('statusText').textContent = currentPage;
    }

    // Check for completion or error
    if (update.type === 'complete') {
        showSuccess();
    } else if (update.type === 'error') {
        showError(update.error || 'An error occurred');
    }
}

// Show success
function showSuccess() {
    document.getElementById('progressView').classList.add('hidden');
    document.getElementById('successView').classList.remove('hidden');
    
    // Refresh history
    refreshHistory();
}

// Show error
function showError(message) {
    document.getElementById('progressView').classList.add('hidden');
    document.getElementById('mainForm').classList.add('hidden');
    document.getElementById('errorView').classList.remove('hidden');
    document.getElementById('errorMessage').textContent = message;
}

// Stop download
window.stopDownload = async function() {
    if (!currentJobId) return;
    
    if (confirm('Are you sure you want to stop this download?')) {
        try {
            await CancelJob(currentJobId);
            showError('Download cancelled by user');
        } catch (error) {
            console.error('Error stopping download:', error);
        }
    }
};

// Open folder
window.openFolder = async function() {
    try {
        await OpenDownloadFolder();
    } catch (error) {
        console.error('Error opening folder:', error);
        alert('Could not open downloads folder');
    }
};

// Start new download
window.startNew = function() {
    currentJobId = null;
    document.getElementById('url').value = '';
    document.querySelector('input[name="mode"][value="static"]').checked = true;
    
    document.getElementById('mainForm').classList.remove('hidden');
    document.getElementById('progressView').classList.add('hidden');
    document.getElementById('successView').classList.add('hidden');
    document.getElementById('errorView').classList.add('hidden');
    
    // Reset progress
    document.getElementById('progressFill').style.width = '0%';
    document.getElementById('progressPercent').textContent = '0%';
    document.getElementById('pagesFound').textContent = '0';
    document.getElementById('pagesDownloaded').textContent = '0';
    document.getElementById('assetsDownloaded').textContent = '0';
    document.getElementById('totalSize').textContent = '0 B';
    document.getElementById('currentStatus').classList.add('hidden');
    
    document.getElementById('url').focus();
};

// Toggle history
window.toggleHistory = async function() {
    const historyPanel = document.getElementById('historyPanel');
    const mainForm = document.getElementById('mainForm');
    
    if (historyPanel.classList.contains('hidden')) {
        await refreshHistory();
        historyPanel.classList.remove('hidden');
        mainForm.classList.add('hidden');
    } else {
        historyPanel.classList.add('hidden');
        mainForm.classList.remove('hidden');
    }
};

// Refresh history
window.refreshHistory = async function() {
    try {
        downloadHistory = await GetHistory();
        renderHistory();
    } catch (error) {
        console.error('Error loading history:', error);
    }
};

// Render history
function renderHistory() {
    const historyList = document.getElementById('historyList');
    
    if (!downloadHistory || downloadHistory.length === 0) {
        historyList.innerHTML = '<div class="empty-state">No downloads yet</div>';
        return;
    }
    
    historyList.innerHTML = downloadHistory.map(item => `
        <div class="history-item">
            <div class="history-info">
                <div class="history-url">${escapeHtml(item.url)}</div>
                <div class="history-meta">
                    <span class="history-mode">${item.mode === 'browser' ? '🌐 Browser' : '⚡ Static'}</span>
                    <span>${formatDate(item.endTime)}</span>
                    <span>${item.pages} pages · ${item.assets} assets · ${formatBytes(item.totalSize)}</span>
                </div>
            </div>
            <button class="btn-icon" onclick="openHistoryFile('${escapeHtml(item.zipPath)}')" title="Open file location">
                📁
            </button>
        </div>
    `).join('');
}

// Open history file
window.openHistoryFile = async function(zipPath) {
    try {
        await OpenFileLocation(zipPath);
    } catch (error) {
        console.error('Error opening file:', error);
        alert('Could not open file location. File may have been moved or deleted.');
    }
};

// Confirm clear history
window.confirmClearHistory = function() {
    if (confirm('Are you sure you want to clear all download history?')) {
        clearAllHistory();
    }
};

// Clear history
async function clearAllHistory() {
    try {
        await ClearHistory();
        downloadHistory = [];
        renderHistory();
    } catch (error) {
        console.error('Error clearing history:', error);
        alert('Failed to clear history');
    }
}

// Format bytes
function formatBytes(bytes) {
    if (!bytes || bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
}

// Format date
function formatDate(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now - date;
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);
    
    if (diffMins < 1) return 'Just now';
    if (diffMins < 60) return `${diffMins} min ago`;
    if (diffHours < 24) return `${diffHours} hour${diffHours > 1 ? 's' : ''} ago`;
    if (diffDays < 7) return `${diffDays} day${diffDays > 1 ? 's' : ''} ago`;
    
    return date.toLocaleDateString();
}

// Escape HTML
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Load history on startup
refreshHistory();
