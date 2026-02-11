import React from 'react';
import './Download.css';

const Download = ({ jobId, onNewDownload }) => {
    const handleDownload = () => {
        window.location.href = `http://localhost:8080/api/download/${jobId}`;
    };

    return (
        <div className="download-container">
            <div className="success-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
            </div>

            <h2 className="success-title">Download Complete!</h2>
            <p className="success-message">
                Your website has been successfully downloaded and packaged.
            </p>

            <div className="button-group">
                <button onClick={handleDownload} className="download-btn">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" className="btn-icon">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                            d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                    Download ZIP
                </button>

                <button onClick={onNewDownload} className="new-download-btn">
                    Download Another Website
                </button>
            </div>
        </div>
    );
};

export default Download;
