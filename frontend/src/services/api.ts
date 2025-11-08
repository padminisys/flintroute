import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

// Create axios instance
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('access_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor to handle token refresh
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const refreshToken = localStorage.getItem('refresh_token');
        if (refreshToken) {
          const response = await axios.post(`${API_BASE_URL}/auth/refresh`, {
            refresh_token: refreshToken,
          });

          const { access_token, refresh_token } = response.data;
          localStorage.setItem('access_token', access_token);
          localStorage.setItem('refresh_token', refresh_token);

          originalRequest.headers.Authorization = `Bearer ${access_token}`;
          return api(originalRequest);
        }
      } catch (refreshError) {
        // Refresh failed, logout user
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);

// Auth API
export const authAPI = {
  login: async (username: string, password: string) => {
    const response = await api.post('/auth/login', { username, password });
    return response.data;
  },

  logout: async () => {
    const response = await api.post('/auth/logout');
    return response.data;
  },

  refresh: async (refreshToken: string) => {
    const response = await api.post('/auth/refresh', { refresh_token: refreshToken });
    return response.data;
  },
};

// BGP Peers API
export const bgpPeersAPI = {
  list: async () => {
    const response = await api.get('/bgp/peers');
    return response.data.peers;
  },

  get: async (id: number) => {
    const response = await api.get(`/bgp/peers/${id}`);
    return response.data;
  },

  create: async (peer: any) => {
    const response = await api.post('/bgp/peers', peer);
    return response.data;
  },

  update: async (id: number, peer: any) => {
    const response = await api.put(`/bgp/peers/${id}`, peer);
    return response.data;
  },

  delete: async (id: number) => {
    const response = await api.delete(`/bgp/peers/${id}`);
    return response.data;
  },
};

// BGP Sessions API
export const bgpSessionsAPI = {
  list: async () => {
    const response = await api.get('/bgp/sessions');
    return response.data.sessions;
  },

  get: async (id: number) => {
    const response = await api.get(`/bgp/sessions/${id}`);
    return response.data;
  },
};

// Config API
export const configAPI = {
  listVersions: async () => {
    const response = await api.get('/config/versions');
    return response.data.versions;
  },

  backup: async (description: string) => {
    const response = await api.post('/config/backup', { description });
    return response.data;
  },

  restore: async (id: number) => {
    const response = await api.post(`/config/restore/${id}`);
    return response.data;
  },
};

// Alerts API
export const alertsAPI = {
  list: async (params?: { acknowledged?: boolean; severity?: string }) => {
    const response = await api.get('/alerts', { params });
    return response.data.alerts;
  },

  acknowledge: async (id: number) => {
    const response = await api.post(`/alerts/${id}/acknowledge`);
    return response.data;
  },
};

export default api;