import React from 'react';
import './Progress.css';

const Progress = ({ status }) => {
    if (!status) return null;

    const {
        pagesDiscovered,
        pagesDownloaded,
        assetsDownloaded,
        currentPage,
        progress,
        totalSize
    } = status;

    const formatSize = (bytes) => {
        if (!bytes) return '0 B';
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
    };

    return (
        <div className="progress-container">
            <div className="progress-header">
                <h3>Downloading Website...</h3>
            </div>

            <div className="progress-bar-wrapper">
                <div className="progress-bar">
                    <div
                        className="progress-bar-fill"
                        style={{ width: `${progress || 0}%` }}
                    >
                        <span className="progress-text">{progress || 0}%</span>
                    </div>
                </div>
            </div>

            <div className="stats-grid">
                <div className="stat-card">
                    <div className="stat-label">Pages Discovered</div>
                    <div className="stat-value">{pagesDiscovered || 0}</div>
                </div>
                <div className="stat-card">
                    <div className="stat-label">Pages Downloaded</div>
                    <div className="stat-value">{pagesDownloaded || 0}</div>
                </div>
                <div className="stat-card">
                    <div className="stat-label">Assets Downloaded</div>
                    <div className="stat-value">{assetsDownloaded || 0}</div>
                </div>
                <div className="stat-card">
                    <div className="stat-label">Total Size</div>
                    <div className="stat-value">{formatSize(totalSize)}</div>
                </div>
            </div>

            {currentPage && (
                <div className="current-page">
                    <span className="current-page-label">Currently processing:</span>
                    <span className="current-page-url">{currentPage}</span>
                </div>
            )}
        </div>
    );
};

export default Progress;
