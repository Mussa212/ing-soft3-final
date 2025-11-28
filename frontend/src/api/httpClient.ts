import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

export const httpClient = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Interceptor to add X-User-ID header if user is logged in
httpClient.interceptors.request.use((config) => {
    const storedUser = localStorage.getItem('currentUser');
    if (storedUser) {
        try {
            const user = JSON.parse(storedUser);
            if (user && user.id) {
                config.headers['X-User-ID'] = user.id.toString();
            }
        } catch (e) {
            console.error('Error parsing currentUser from localStorage', e);
        }
    }
    return config;
});
