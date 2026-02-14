import React from 'react';
import './Download.css';

const Download = ({ jobId, onNewDownload }) => {
    const handleDownload = () => {
        window.location.href = `http://localhost:8080/api/download/${jobId}`;
    };

    return (
        <div className="download-container">
            <div className="success-icon-wrapper">
                <div className="success-ring"></div>
                <div className="success-icon">
                    <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="3">
                        <path d="M20 6L9 17l-5-5" strokeLinecap="round" strokeLinejoin="round" />
                    </svg>
                </div>
            </div>

            <h2 className="success-title">Archive Ready</h2>
            <p className="success-message">
                Your website has been successfully captured. All assets, scripts, and styles are packaged and ready for offline use.
            </p>

            <div className="button-group">
                <button onClick={handleDownload} className="download-btn">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" className="btn-icon" strokeWidth="2">
                        <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4M7 10l5 5 5-5M12 15V3" strokeLinecap="round" strokeLinejoin="round" />
                    </svg>
                    Download ZIP
                </button>

                <button onClick={onNewDownload} className="new-download-btn">
                    Archive Another
                </button>
            </div>
        </div>
    );
};

export default Download;
