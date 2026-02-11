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
        console.log('WebSocket message:', data);

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

    // Use WebSocket hook
    useWebSocket(wsUrl, handleWebSocketMessage);

    // Handle form submission
    const handleSubmit = async (url) => {
        setIsLoading(true);
        setError(null);
        setIsComplete(false);
        setStatus(null);

        try {
            const response = await api.startScrape(url);
            setJobId(response.jobId);
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to start download');
            setIsLoading(false);
        }
    };

    // Handle new download
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
                    <h1 className="title">
                        <span className="title-icon">🌐</span>
                        Website Downloader
                    </h1>
                    <p className="subtitle">
                        Download entire websites with all assets for offline viewing
                    </p>
                </header>

                {error && (
                    <div className="error-banner">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" className="error-icon">
                            <circle cx="12" cy="12" r="10" strokeWidth="2" />
                            <line x1="12" y1="8" x2="12" y2="12" strokeWidth="2" strokeLinecap="round" />
                            <line x1="12" y1="16" x2="12.01" y2="16" strokeWidth="2" strokeLinecap="round" />
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

                <footer className="footer">
                    <p>Built with React & Golang • Download websites while preserving structure</p>
                </footer>
            </div>
        </div>
    );
}

export default App;
