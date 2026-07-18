import axios from 'axios';
import { loadConfig } from './config';

const config = loadConfig();

const api = axios.create({
    baseURL: config.apiUrl,
    headers: { 'Content-Type': 'application/json' }
});

// Attach JWT
api.interceptors.request.use((req) => {
    if (config.token) {
        req.headers.Authorization = `Bearer ${config.token}`;
    }
    return req;
});

// Handle 401
api.interceptors.response.use(
    (res) => res,
    (err) => {
        if (err.response?.status === 401) {
            console.error('Session expired. Please login again.');
            process.exit(1);
        }
        return Promise.reject(err);
    }
);

export default api;