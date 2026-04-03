import axios from 'axios';

const baseURL = (import.meta.env.VITE_API_URL || 'http://localhost:3000').replace(/\/$/, '');
console.log('API Base URL:', baseURL);

const api = axios.create({
  baseURL,
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  console.log('Interceptor - Token found:', !!token);
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('API Error:', error.response?.status, error.response?.data);
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      // window.location.href = '/login'; // Optional: force redirect on 401
    }
    return Promise.reject(error);
  }
);

export default api;
