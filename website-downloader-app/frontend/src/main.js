import './style.css';
import './app.css';

import {StartScrape, GetJobStatus, OpenDownloadFolder, CancelJob, GetHistory, OpenFileLocation, ClearHistory} from '../wailsjs/go/main/App';
import {EventsOn} from '../wailsjs/runtime/runtime';

let currentJobId = null;
let downloadHistory = [];

document.querySelector('#app').innerHTML = `
    <div class="container">
        <div class="header">
            <h1>🌐 Website Downloader</h1>
            <p class="subtitle">Download entire websites for offline viewing</p>
            <button class="btn-history" onclick="toggleHistory()">📜 History</button>
        </div>

        <div id="historySection" class="history-section hidden">
            <div class="history-header">
                <h3>📜 Download History</h3>
                <div class="history-actions">
                    <button class="btn-small" onclick="refreshHistory()">🔄 Refresh</button>
                    <button class="btn-small btn-danger" onclick="confirmClearHistory()">🗑️ Clear All</button>
                    <button class="btn-small" onclick="toggleHistory()">✖ Close</button>
                </div>
            </div>
            <div id="historyList" class="history-list">
                <p class="empty-state">No downloads yet</p>
            </div>
        </div>

        <div class="form-section">
            <div class="input-group">
                <label for="url">Website URL</label>
                <input 
                    type="text" 
                    id="url" 
                    placeholder="https://example.com" 
                    autocomplete="off"
                />
            </div>

            <div class="input-group">
                <label for="mode">Download Mode</label>
                <select id="mode">
                    <option value="static">Static Mode (Fast - Traditional Websites)</option>
                    <option value="browser">Browser Mode (SPAs - React, Vue, Angular)</option>
                </select>
            </div>

            <button class="btn-primary" id="startBtn" onclick="startDownload()">
                📥 Start Download
            </button>
        </div>

        <div id="progressSection" class="progress-section hidden">
            <div class="progress-header">
                <h3>Downloading...</h3>
                <div class="header-actions">
                    <span id="statusText" class="status-text">Starting...</span>
                    <button class="btn-stop" onclick="stopDownload()">⏹️ Stop</button>
                </div>
            </div>
            
            <div class="progress-bar">
                <div id="progressFill" class="progress-fill"></div>
                <span id="progressText" class="progress-text">0%</span>
            </div>

            <div class="stats">
                <div class="stat-item">
                    <span class="stat-label">Pages:</span>
                    <span id="pagesCount" class="stat-value">0</span>
                </div>
                <div class="stat-item">
                    <span class="stat-label">Assets:</span>
                    <span id="assetsCount" class="stat-value">0</span>
                </div>
                <div class="stat-item">
                    <span class="stat-label">Size:</span>
                    <span id="sizeCount" class="stat-value">0 KB</span>
                </div>
            </div>

            <div id="currentPage" class="current-page"></div>
        </div>

        <div id="completeSection" class="complete-section hidden">
            <div class="success-icon">✅</div>
            <h3>Download Complete!</h3>
            <p id="completeMessage">Your website has been downloaded successfully.</p>
            <div class="button-group">
                <button class="btn-success" onclick="openFolder()">
                    📁 Open Downloads Folder
                </button>
                <button class="btn-secondary" onclick="startNew()">
                    🔄 Download Another
                </button>
            </div>
        </div>

        <div id="errorSection" class="error-section hidden">
            <div class="error-icon">❌</div>
            <h3>Download Failed</h3>
            <p id="errorMessage" class="error-message"></p>
            <button class="btn-secondary" onclick="startNew()">
                🔄 Try Again
            </button>
        </div>
    </div>
`;

// Focus URL input
document.getElementById('url').focus();

// Start download function
window.startDownload = async function() {
    const url = document.getElementById('url').value.trim();
    const mode = document.getElementById('mode').value;

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

    // Hide form, show progress
    document.querySelector('.form-section').classList.add('hidden');
    document.getElementById('progressSection').classList.remove('hidden');
    document.getElementById('completeSection').classList.add('hidden');
    document.getElementById('errorSection').classList.add('hidden');

    try {
        const response = await StartScrape(url, mode);
        currentJobId = response.jobId;

        // Listen for progress updates
        EventsOn('progress:' + currentJobId, (update) => {
            console.log('Progress update:', update);
            updateProgress(update);
        });

    } catch (error) {
        console.error('Error starting download:', error);
        showError('Failed to start download: ' + error);
    }
};

// Update progress display
function updateProgress(update) {
    const progress = update.progress || 0;
    const pagesDownloaded = update.pagesDownloaded || 0;
    const assetsDownloaded = update.assetsDownloaded || 0;
    const totalSize = update.totalSize || 0;
    const currentPage = update.currentPage || '';

    // Update progress bar
    document.getElementById('progressFill').style.width = progress + '%';
    document.getElementById('progressText').textContent = progress + '%';

    // Update stats
    document.getElementById('pagesCount').textContent = pagesDownloaded;
    document.getElementById('assetsCount').textContent = assetsDownloaded;
    document.getElementById('sizeCount').textContent = formatBytes(totalSize);

    // Update current page
    if (currentPage) {
        document.getElementById('currentPage').textContent = 'Current: ' + currentPage;
    }

    // Update status
    if (update.type === 'complete') {
        showComplete(update.message || 'Download complete!');
    } else if (update.type === 'error') {
        showError(update.error || 'An error occurred');
    } else {
        document.getElementById('statusText').textContent = 
            pagesDownloaded > 0 ? 'Downloading...' : 'Starting...';
    }
}

// Show completion screen
function showComplete(message) {
    document.getElementById('progressSection').classList.add('hidden');
    document.getElementById('completeSection').classList.remove('hidden');
    document.getElementById('completeMessage').textContent = message;
}

// Show error screen
function showError(error) {
    document.getElementById('progressSection').classList.add('hidden');
    document.getElementById('errorSection').classList.remove('hidden');
    document.getElementById('errorMessage').textContent = error;
}

// Open downloads folder
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
    document.getElementById('mode').value = 'static';
    
    document.querySelector('.form-section').classList.remove('hidden');
    document.getElementById('progressSection').classList.add('hidden');
    document.getElementById('completeSection').classList.add('hidden');
    document.getElementById('errorSection').classList.add('hidden');
    
    // Reset progress
    document.getElementById('progressFill').style.width = '0%';
    document.getElementById('progressText').textContent = '0%';
    document.getElementById('pagesCount').textContent = '0';
    document.getElementById('assetsCount').textContent = '0';
    document.getElementById('sizeCount').textContent = '0 KB';
    document.getElementById('currentPage').textContent = '';
    
    document.getElementById('url').focus();
};

// Format bytes to readable size
function formatBytes(bytes) {
    if (bytes === 0) return '0 KB';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
}

// Handle Enter key in URL input
document.getElementById('url').addEventListener('keypress', (e) => {
    if (e.key === 'Enter') {
        startDownload();
    }
});


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

// Toggle history view
window.toggleHistory = async function() {
    const historySection = document.getElementById('historySection');
    const formSection = document.querySelector('.form-section');
    
    if (historySection.classList.contains('hidden')) {
        // Show history
        await refreshHistory();
        historySection.classList.remove('hidden');
        formSection.classList.add('hidden');
    } else {
        // Hide history
        historySection.classList.add('hidden');
        formSection.classList.remove('hidden');
    }
};

// Refresh history list
window.refreshHistory = async function() {
    try {
        downloadHistory = await GetHistory();
        renderHistory();
    } catch (error) {
        console.error('Error loading history:', error);
    }
};

// Render history list
function renderHistory() {
    const historyList = document.getElementById('historyList');
    
    if (!downloadHistory || downloadHistory.length === 0) {
        historyList.innerHTML = '<p class="empty-state">No downloads yet</p>';
        return;
    }
    
    historyList.innerHTML = downloadHistory.map(item => `
        <div class="history-item">
            <div class="history-info">
                <div class="history-url">${escapeHtml(item.url)}</div>
                <div class="history-meta">
                    <span class="history-mode">${item.mode === 'browser' ? '🌐 Browser' : '⚡ Static'}</span>
                    <span class="history-date">${formatDate(item.endTime)}</span>
                    <span class="history-stats">
                        ${item.pages} pages · ${item.assets} assets · ${formatBytes(item.totalSize)}
                    </span>
                </div>
            </div>
            <div class="history-actions">
                <button class="btn-icon" onclick="openHistoryFile('${escapeHtml(item.zipPath)}')" title="Open file location">
                    📁
                </button>
            </div>
        </div>
    `).join('');
}

// Open file from history
window.openHistoryFile = async function(zipPath) {
    try {
        await OpenFileLocation(zipPath);
    } catch (error) {
        console.error('Error opening file:', error);
        alert('Could not open file location. File may have been moved or deleted.');
    }
};

// Confirm and clear history
window.confirmClearHistory = function() {
    if (confirm('Are you sure you want to clear all download history? This cannot be undone.')) {
        clearAllHistory();
    }
};

// Clear all history
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
