// 認証サービス
import api from './api';

// 認証関連のエンドポイント
const AUTH_ENDPOINTS = {
  SIGNUP: '/auth/signup',
  LOGIN: '/auth/login',
  CURRENT_USER: '/auth/me'
};

/**
 * サインアップ（ユーザー登録）
 * @param {Object} userData - ユーザー登録データ
 * @param {string} userData.email - メールアドレス
 * @param {string} userData.password - パスワード
 * @param {string} userData.name - ユーザー名
 * @returns {Promise} APIレスポンス
 */
export const signup = async (userData) => {
  try {
    const response = await api.post(AUTH_ENDPOINTS.SIGNUP, userData);
    // 成功したらトークンとユーザー情報を保存
    if (response.data && response.data.token) {
      localStorage.setItem('token', response.data.token);
      localStorage.setItem('user', JSON.stringify(response.data.user));
    }
    return response.data;
  } catch (error) {
    console.error('サインアップエラー:', error);
    throw error;
  }
};

/**
 * ログイン
 * @param {Object} credentials - ログイン認証情報
 * @param {string} credentials.email - メールアドレス
 * @param {string} credentials.password - パスワード
 * @returns {Promise} APIレスポンス
 */
export const login = async (credentials) => {
  try {
    const response = await api.post(AUTH_ENDPOINTS.LOGIN, credentials);
    // 成功したらトークンとユーザー情報を保存
    if (response.data && response.data.token) {
      localStorage.setItem('token', response.data.token);
      localStorage.setItem('user', JSON.stringify(response.data.user));
    }
    return response.data;
  } catch (error) {
    console.error('ログインエラー:', error);
    throw error;
  }
};

/**
 * ログアウト
 */
export const logout = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('user');
};

/**
 * 現在のユーザー情報を取得
 * @returns {Promise} APIレスポンス
 */
export const getCurrentUser = async () => {
  try {
    const response = await api.get(AUTH_ENDPOINTS.CURRENT_USER);
    return response.data;
  } catch (error) {
    console.error('ユーザー情報取得エラー:', error);
    throw error;
  }
};

/**
 * ローカルストレージからユーザー情報を取得
 * @returns {Object|null} ユーザー情報またはnull
 */
export const getStoredUser = () => {
  const userStr = localStorage.getItem('user');
  return userStr ? JSON.parse(userStr) : null;
};

/**
 * 認証トークンを取得
 * @returns {string|null} 認証トークンまたはnull
 */
export const getToken = () => {
  return localStorage.getItem('token');
};

/**
 * ユーザーが認証済みかどうか確認
 * @returns {boolean} 認証済みならtrue
 */
export const isAuthenticated = () => {
  return !!getToken();
};

export default {
  signup,
  login,
  logout,
  getCurrentUser,
  getStoredUser,
  getToken,
  isAuthenticated
};
