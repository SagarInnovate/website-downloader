import React, { useState } from 'react';
import './UrlInput.css';

const UrlInput = ({ onSubmit, isLoading }) => {
    const [url, setUrl] = useState('');
    const [mode, setMode] = useState('static');
    const [error, setError] = useState('');

    const handleSubmit = (e) => {
        e.preventDefault();
        if (!url.trim()) {
            setError('Please enter a valid URL');
            return;
        }
        try {
            const urlObj = new URL(url);
            if (!['http:', 'https:'].includes(urlObj.protocol)) {
                setError('URL must use HTTP or HTTPS');
                return;
            }
        } catch (err) {
            setError('Please enter a valid URL');
            return;
        }
        setError('');
        onSubmit(url, mode);
    };

    return (
        <div className="url-input-container">
            <form onSubmit={handleSubmit} className="url-form">
                <div className="input-group">
                    <label className="input-label">
                        Target Website
                        <span className="input-label-hint">Full URL including https://</span>
                    </label>
                    <input
                        type="text"
                        value={url}
                        onChange={(e) => setUrl(e.target.value)}
                        placeholder="e.g. https://stripe.com/docs"
                        className="url-input"
                        disabled={isLoading}
                        autoFocus
                    />
                </div>

                <div className="mode-group">
                    <label className="input-label">Extraction Method</label>
                    <div className="mode-options">
                        <label className="mode-card">
                            <input
                                type="radio"
                                name="mode"
                                value="static"
                                checked={mode === 'static'}
                                onChange={(e) => setMode(e.target.value)}
                                disabled={isLoading}
                            />
                            <div className="mode-content">
                                <span className="mode-title">
                                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
                                        <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z" />
                                    </svg>
                                    Static Fetch
                                </span>
                                <span className="mode-desc">Lightning fast. Ideal for blogs, documentation, and traditional server-rendered sites.</span>
                            </div>
                        </label>
                        <label className="mode-card">
                            <input
                                type="radio"
                                name="mode"
                                value="browser"
                                checked={mode === 'browser'}
                                onChange={(e) => setMode(e.target.value)}
                                disabled={isLoading}
                            />
                            <div className="mode-content">
                                <span className="mode-title">
                                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
                                        <circle cx="12" cy="12" r="10" />
                                        <line x1="2" y1="12" x2="22" y2="12" />
                                        <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z" />
                                    </svg>
                                    Browser Render
                                </span>
                                <span className="mode-desc">Full simulation. Executes JavaScript for SPAs like React, Vue, and Angular applications.</span>
                            </div>
                        </label>
                    </div>
                </div>

                {error && (
                    <div className="error-msg">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                            <circle cx="12" cy="12" r="10"></circle>
                            <line x1="12" y1="8" x2="12" y2="12"></line>
                            <line x1="12" y1="16" x2="12.01" y2="16"></line>
                        </svg>
                        {error}
                    </div>
                )}

                <button type="submit" className="submit-btn" disabled={isLoading}>
                    {isLoading ? (
                        <>
                            <div className="spinner"></div>
                            <span>Starting Engine...</span>
                        </>
                    ) : (
                        <>
                            Start Archiving
                            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
                                <path d="M5 12h14M12 5l7 7-7 7" />
                            </svg>
                        </>
                    )}
                </button>
            </form>
        </div>
    );
};

export default UrlInput;
