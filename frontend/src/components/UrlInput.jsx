import React, { useState } from 'react';
import './UrlInput.css';

const UrlInput = ({ onSubmit, isLoading }) => {
    const [url, setUrl] = useState('');
    const [error, setError] = useState('');

    const handleSubmit = (e) => {
        e.preventDefault();

        // Validate URL
        if (!url.trim()) {
            setError('Please enter a URL');
            return;
        }

        try {
            const urlObj = new URL(url);
            if (!['http:', 'https:'].includes(urlObj.protocol)) {
                setError('URL must use HTTP or HTTPS protocol');
                return;
            }
        } catch (err) {
            setError('Please enter a valid URL');
            return;
        }

        setError('');
        onSubmit(url);
    };

    return (
        <div className="url-input-container">
            <form onSubmit={handleSubmit} className="url-form">
                <div className="input-wrapper">
                    <input
                        type="text"
                        value={url}
                        onChange={(e) => setUrl(e.target.value)}
                        placeholder="Enter website URL (e.g., https://example.com)"
                        className="url-input"
                        disabled={isLoading}
                    />
                    <button
                        type="submit"
                        className="submit-btn"
                        disabled={isLoading}
                    >
                        {isLoading ? (
                            <>
                                <span className="spinner"></span>
                                Downloading...
                            </>
                        ) : (
                            'Download Website'
                        )}
                    </button>
                </div>
                {error && <div className="error-message">{error}</div>}
            </form>
        </div>
    );
};

export default UrlInput;
