// APIサービス
import axios from 'axios';

// API基本設定
const api = axios.create({
  baseURL: process.env.VUE_APP_API_URL || 'http://localhost:8080/api',
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  }
});

// リクエストインターセプター（認証トークンの追加）
api.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  error => {
    return Promise.reject(error);
  }
);

// レスポンスインターセプター（エラーハンドリング）
api.interceptors.response.use(
  response => response,
  error => {
    // 認証エラーの場合はローカルストレージをクリアしてログイン画面にリダイレクト
    if (error.response && error.response.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      if (window.location.pathname !== '/') {
        window.location.href = '/';
      }
    }
    return Promise.reject(error);
  }
);

export default api;
