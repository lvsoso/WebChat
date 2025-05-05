import { createApp } from 'vue';
import { createPinia } from 'pinia';
import App from './App.vue';
import router from './router';
import axios from 'axios';

// Import CSS
import './assets/main.css';

// Configure axios
axios.defaults.baseURL = 'http://localhost:8080';
axios.interceptors.request.use(config => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Create app
const app = createApp(App);

// Use plugins
app.use(createPinia());
app.use(router);

// Initialize auth check
import { useAuthStore } from './stores/auth';
const authStore = useAuthStore();
authStore.checkAuth();

// Mount app
app.mount('#app');