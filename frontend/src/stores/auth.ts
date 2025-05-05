import { defineStore } from 'pinia';
import axios from 'axios';

interface LoginPayload {
  email: string;
  password: string;
}

interface RegisterPayload {
  email: string;
  password: string;
}

interface User {
  id: string;
  email: string;
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as User | null,
    token: null as string | null,
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
  },

  actions: {
    async login(payload: LoginPayload) {
      try {
        const response = await axios.post('/api/auth/login', payload);
        this.token = response.data.token;
        this.user = response.data.user;
        
        // 保存token到localStorage
        localStorage.setItem('token', this.token);
      } catch (error) {
        console.error('Login failed:', error);
        throw error;
      }
    },

    async register(payload: RegisterPayload) {
      try {
        const response = await axios.post('/api/auth/register', payload);
        return response.data;
      } catch (error) {
        console.error('Registration failed:', error);
        throw error;
      }
    },

    async logout() {
      try {
        await axios.post('/api/auth/logout');
      } catch (error) {
        console.error('Logout failed:', error);
      } finally {
        this.token = null;
        this.user = null;
        localStorage.removeItem('token');
      }
    },

    async checkAuth() {
      const token = localStorage.getItem('token');
      if (!token) return;

      try {
        const response = await axios.get('/api/auth/me');
        this.user = response.data.user;
        this.token = token;
      } catch (error) {
        console.error('Auth check failed:', error);
        this.token = null;
        this.user = null;
        localStorage.removeItem('token');
      }
    },
  },
});