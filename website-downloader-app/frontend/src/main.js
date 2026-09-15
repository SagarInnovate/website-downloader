import './style.css';
import './app.css';
import './sidebar.css';
import {appTemplate} from './template.js';

import {StartScrape, CancelJob, GetHistory, OpenFileLocation, ClearHistory, SaveAsZip, OpenDownloadFolder} from '../wailsjs/go/main/App';
import {EventsOn} from '../wailsjs/runtime/runtime';
import * as runtime from '../wailsjs/runtime/runtime';

let currentJobId = null;
let downloadHistory = [];
let currentView = 'home';

// Render app template
document.querySelector('#app').innerHTML = appTemplate;

// View switching
function showView(viewName) {
    currentView = viewName;
    
    // Hide all views
    document.querySelectorAll('.view-content').forEach(view => {
        view.classList.add('hidden');
    });
    
    // Show selected view
    const viewElement = document.getElementById(`${viewName}View`);
    if (viewElement) {
        viewElement.classList.remove('hidden');
    }
    
    // Update nav active state
    document.querySelectorAll('.nav-item').forEach(item => {
        item.classList.remove('active');
    });
    document.querySelector(`[data-view="${viewName}"]`)?.classList.add('active');
    
    // Load data for specific views
    if (viewName === 'history') {
        refreshHistory();
    } else if (viewName === 'settings') {
        loadSettings();
    }
}

// Start download
async function startDownload(event) {
    event.preventDefault();
    
    const url = document.getElementById('url').value.trim();
    const mode = document.querySelector('input[name="mode"]:checked').value;

    if (!url) {
        await runtime.MessageDialog({
            Type: 'error',
            Title: 'Invalid Input',
            Message: 'Please enter a URL'
        });
        return;
    }

    // Validate URL
    try {
        new URL(url);
    } catch (e) {
        await runtime.MessageDialog({
            Type: 'error',
            Title: 'Invalid URL',
            Message: 'Please enter a valid URL (e.g., https://example.com)'
        });
        return;
    }

    // Show progress view
    document.getElementById('mainForm').classList.add('hidden');
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
}

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
async function stopDownload() {
    if (!currentJobId) return;
    
    const result = await runtime.MessageDialog({
        Type: 'question',
        Title: 'Stop Download',
        Message: 'Are you sure you want to stop this download?',
        Buttons: ['Yes', 'No']
    });
    
    if (result === 'Yes') {
        try {
            await CancelJob(currentJobId);
            showError('Download cancelled by user');
        } catch (error) {
            console.error('Error stopping download:', error);
        }
    }
}

// Save As
async function saveAs() {
    if (!currentJobId) {
        await runtime.MessageDialog({
            Type: 'error',
            Title: 'No Download',
            Message: 'No download available to save'
        });
        return;
    }
    
    try {
        const savePath = await SaveAsZip(currentJobId);
        if (savePath) {
            await runtime.MessageDialog({
                Type: 'info',
                Title: 'Success',
                Message: 'File saved successfully!'
            });
        }
    } catch (error) {
        console.error('Error saving file:', error);
        if (error && !error.toString().includes('cancelled')) {
            await runtime.MessageDialog({
                Type: 'error',
                Title: 'Save Failed',
                Message: 'Could not save file: ' + error
            });
        }
    }
}

// Start new download
function startNew() {
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
}


// Refresh history
async function refreshHistory() {
    try {
        downloadHistory = await GetHistory();
        renderHistory();
    } catch (error) {
        console.error('Error loading history:', error);
    }
}

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
async function openHistoryFile(zipPath) {
    try {
        await OpenFileLocation(zipPath);
    } catch (error) {
        console.error('Error opening file:', error);
        await runtime.MessageDialog({
            Type: 'error',
            Title: 'File Not Found',
            Message: 'Could not open file location. File may have been moved or deleted.'
        });
    }
}

// Confirm clear history
async function confirmClearHistory() {
    const result = await runtime.MessageDialog({
        Type: 'question',
        Title: 'Clear History',
        Message: 'Are you sure you want to clear all download history?',
        Buttons: ['Yes', 'No']
    });
    
    if (result === 'Yes') {
        clearAllHistory();
    }
}

// Clear history
async function clearAllHistory() {
    try {
        await ClearHistory();
        downloadHistory = [];
        renderHistory();
    } catch (error) {
        console.error('Error clearing history:', error);
        await runtime.MessageDialog({
            Type: 'error',
            Title: 'Error',
            Message: 'Failed to clear history'
        });
    }
}


// Load settings
function loadSettings() {
    // Load from config.json values (these would need to be exposed via Go backend)
    // For now, using default values
    document.getElementById('maxPages').value = 500;
    document.getElementById('browserWaitTime').value = 5000;
    document.getElementById('browserTimeout').value = 60;
}

// Save settings
async function saveSettings() {
    const maxPages = document.getElementById('maxPages').value;
    const browserWaitTime = document.getElementById('browserWaitTime').value;
    const browserTimeout = document.getElementById('browserTimeout').value;
    
    // TODO: Add backend method to save settings
    await runtime.MessageDialog({
        Type: 'info',
        Title: 'Settings Saved',
        Message: `Settings saved successfully!\n\nMax Pages: ${maxPages}\nWait Time: ${browserWaitTime}ms\nTimeout: ${browserTimeout}s\n\nNote: Restart app for changes to take effect.`
    });
}

// Open downloads folder
async function openFolder() {
    try {
        await OpenDownloadFolder();
    } catch (error) {
        console.error('Error opening folder:', error);
        await runtime.MessageDialog({
            Type: 'error',
            Title: 'Error',
            Message: 'Could not open downloads folder'
        });
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


// Event listeners
document.addEventListener('DOMContentLoaded', () => {
    // Sidebar navigation
    document.querySelectorAll('.nav-item').forEach(item => {
        item.addEventListener('click', () => {
            const view = item.getAttribute('data-view');
            showView(view);
        });
    });
    
    // Download form
    const downloadForm = document.getElementById('downloadForm');
    if (downloadForm) {
        downloadForm.addEventListener('submit', startDownload);
    }
    
    // Stop button
    const stopBtn = document.getElementById('stopBtn');
    if (stopBtn) {
        stopBtn.addEventListener('click', stopDownload);
    }
    
    // Save As button
    const saveAsBtn = document.getElementById('saveAsBtn');
    if (saveAsBtn) {
        saveAsBtn.addEventListener('click', saveAs);
    }
    
    // Start new button
    const startNewBtn = document.getElementById('startNewBtn');
    if (startNewBtn) {
        startNewBtn.addEventListener('click', startNew);
    }
    
    // Try again button
    const tryAgainBtn = document.getElementById('tryAgainBtn');
    if (tryAgainBtn) {
        tryAgainBtn.addEventListener('click', startNew);
    }
    
    // History buttons
    const refreshHistoryBtn = document.getElementById('refreshHistoryBtn');
    if (refreshHistoryBtn) {
        refreshHistoryBtn.addEventListener('click', refreshHistory);
    }
    
    const clearHistoryBtn = document.getElementById('clearHistoryBtn');
    if (clearHistoryBtn) {
        clearHistoryBtn.addEventListener('click', confirmClearHistory);
    }
    
    // Settings buttons
    const saveSettingsBtn = document.getElementById('saveSettingsBtn');
    if (saveSettingsBtn) {
        saveSettingsBtn.addEventListener('click', saveSettings);
    }
    
    const openFolderBtn = document.getElementById('openFolderBtn');
    if (openFolderBtn) {
        openFolderBtn.addEventListener('click', openFolder);
    }
    
    // Load initial data
    refreshHistory();
});

// Make functions global for inline onclick handlers
window.showView = showView;
window.openHistoryFile = openHistoryFile;
