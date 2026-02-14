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
                <h3>
                    <div className="spinner-pulse"></div>
                    Archiving in Progress
                </h3>
            </div>

            <div className="progress-bar-container">
                <div className="progress-bar-bg">
                    <div
                        className="progress-bar-fill"
                        style={{ width: `${progress || 0}%` }}
                    ></div>
                </div>
                <div className="progress-meta">
                    <span>Processing resources...</span>
                    <span className="progress-percent">{Math.round(progress || 0)}%</span>
                </div>
            </div>

            <div className="stats-grid">
                <div className="stat-item">
                    <div className="stat-label">Pages Found</div>
                    <div className="stat-value">{pagesDiscovered || 0}</div>
                </div>
                <div className="stat-item">
                    <div className="stat-label">Downloaded</div>
                    <div className="stat-value">{pagesDownloaded || 0}</div>
                </div>
                <div className="stat-item">
                    <div className="stat-label">Assets</div>
                    <div className="stat-value">{assetsDownloaded || 0}</div>
                </div>
                <div className="stat-item">
                    <div className="stat-label">Size</div>
                    <div className="stat-value">{formatSize(totalSize)}</div>
                </div>
            </div>

            {currentPage && (
                <div className="current-status">
                    <svg className="status-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
                        <path d="M12 2v4m0 12v4M4.93 4.93l2.83 2.83m8.48 8.48l2.83 2.83M2 12h4m12 0h4M4.93 19.07l2.83-2.83m8.48-8.48l2.83-2.83" />
                    </svg>
                    <span className="status-text">{currentPage}</span>
                </div>
            )}
        </div>
    );
};

export default Progress;
