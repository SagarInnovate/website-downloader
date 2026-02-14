import React, { useState, useCallback } from 'react';
import UrlInput from './components/UrlInput';
import Progress from './components/Progress';
import Download from './components/Download';
import { api } from './services/api';
import { useWebSocket } from './hooks/useWebSocket';
import './App.css';

function App() {
    const [jobId, setJobId] = useState(null);
    const [isLoading, setIsLoading] = useState(false);
    const [status, setStatus] = useState(null);
    const [error, setError] = useState(null);
    const [isComplete, setIsComplete] = useState(false);

    // WebSocket URL
    const wsUrl = jobId ? api.getWebSocketUrl(jobId) : null;

    // Handle WebSocket messages
    const handleWebSocketMessage = useCallback((data) => {
        setStatus({
            pagesDiscovered: data.pagesDiscovered || 0,
            pagesDownloaded: data.pagesDownloaded || 0,
            assetsDownloaded: data.assetsDownloaded || 0,
            totalSize: data.totalSize || 0,
            currentPage: data.currentPage || '',
            progress: data.progress || 0,
        });

        if (data.type === 'complete') {
            setIsComplete(true);
            setIsLoading(false);
        } else if (data.type === 'error') {
            setError(data.error || 'An error occurred');
            setIsLoading(false);
        }
    }, []);

    useWebSocket(wsUrl, handleWebSocketMessage);

    const handleSubmit = async (url, mode) => {
        setIsLoading(true);
        setError(null);
        setIsComplete(false);
        setStatus(null);

        try {
            const response = await api.startScrape(url, mode);
            setJobId(response.jobId);
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to start download');
            setIsLoading(false);
        }
    };

    const handleNewDownload = () => {
        setJobId(null);
        setIsLoading(false);
        setStatus(null);
        setError(null);
        setIsComplete(false);
    };

    return (
        <div className="app">
            <div className="container">
                <header className="header">
                    <div className="brand-bar">
                        <div className="brand">
                            <div className="brand-icon">
                                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
                                    <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z" />
                                </svg>
                            </div>
                            <span className="brand-name">Website Downloader</span>
                        </div>
                        <a href="https://github.com/sagarinnovate/website-downloader" target="_blank" rel="noopener noreferrer" className="github-btn">
                            <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                                <path d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" />
                            </svg>
                            Star on GitHub
                        </a>
                    </div>

                    <div className="hero-text">
                        <h1 className="title-main">
                            Download <span className="text-primary">Any Website</span>
                        </h1>
                        <p className="subtitle">
                            Production-ready offline scraping with structure preservation.
                            <br />
                            Built by <strong>sagarinnovate</strong>.
                        </p>
                    </div>
                </header>

                <main>
                    {error && (
                        <div className="error-banner">
                            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                <circle cx="12" cy="12" r="10" />
                                <line x1="12" y1="8" x2="12" y2="12" />
                                <line x1="12" y1="16" x2="12.01" y2="16" />
                            </svg>
                            <span>{error}</span>
                        </div>
                    )}

                    {!isComplete ? (
                        <>
                            <UrlInput onSubmit={handleSubmit} isLoading={isLoading} />
                            {isLoading && status && <Progress status={status} />}
                        </>
                    ) : (
                        <Download jobId={jobId} onNewDownload={handleNewDownload} />
                    )}
                </main>
            </div>

            <footer className="footer">
                <div className="container" style={{ padding: '0' }}>
                    <p>
                        Made with 💙 by <a href="https://github.com/sagarinnovate" target="_blank" rel="noopener noreferrer" className="footer-link">Sagar Shinde (@sagarinnovate)</a>
                    </p>
                </div>
            </footer>
        </div>
    );
}

export default App;
