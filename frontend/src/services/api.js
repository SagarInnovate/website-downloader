import axios from 'axios';

// Use relative URL to leverage Vite's proxy configuration
const API_BASE_URL = '/api';

export const api = {
    // Start a new scraping job
    startScrape: async (url, mode = 'static') => {
        const response = await axios.post(`${API_BASE_URL}/scrape`, { url, mode });
        return response.data;
    },

    // Get job status
    getStatus: async (jobId) => {
        const response = await axios.get(`${API_BASE_URL}/status/${jobId}`);
        return response.data;
    },

    // Get download URL
    getDownloadUrl: (jobId) => {
        return `${API_BASE_URL}/download/${jobId}`;
    },

    // Get WebSocket URL - WebSocket doesn't go through Vite proxy
    getWebSocketUrl: (jobId) => {
        return `ws://localhost:8080/api/progress/${jobId}`;
    }
};

export default api;
