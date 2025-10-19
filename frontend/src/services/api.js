import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add token to requests if available
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Handle 401 errors by clearing token
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Auth APIs
export const register = (email, password) => {
  return api.post('/api/auth/register', { email, password });
};

export const login = (email, password) => {
  return api.post('/api/auth/login', { email, password });
};

// LinkedIn Connection APIs
export const connectLinkedInWithCookie = (cookie) => {
  return api.post('/api/linkedin/connect/cookie', { cookie });
};

export const connectLinkedInWithCredentials = (username, password) => {
  return api.post('/api/linkedin/connect/credentials', { username, password });
};

// Account APIs
export const getAccounts = () => {
  return api.get('/api/accounts');
};

export const deleteAccount = (accountId) => {
  return api.delete(`/api/accounts/${accountId}`);
};

export default api;

