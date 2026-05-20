/**
 * API 客户端配置
 * 基于 axios 的统一封装
 */
import axios from 'axios';
import type { AxiosInstance, InternalAxiosRequestConfig } from 'axios';

// API 基础配置
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1';
const REQUEST_TIMEOUT = 30000;

// 创建 axios 实例
const createApiClient = (): AxiosInstance => {
  const client = axios.create({
    baseURL: API_BASE_URL,
    timeout: REQUEST_TIMEOUT,
    headers: {
      'Content-Type': 'application/json',
    },
  });

  // 请求拦截器 - 添加认证令牌
  client.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
      // 从 localStorage 获取 token（Zustand persist 会存储在这里）
      const token = localStorage.getItem('auth-storage');
      if (token) {
        try {
          const parsed = JSON.parse(token);
          if (parsed.state?.token) {
            config.headers.Authorization = `Bearer ${parsed.state.token}`;
          }
        } catch {
          // 无效的 token 数据
        }
      }
      return config;
    },
    (error) => Promise.reject(error)
  );

  // 响应拦截器 - 统一错误处理
  client.interceptors.response.use(
    (response) => response,
    (error) => {
      if (error.response?.status === 401) {
        // 清除认证状态
        localStorage.removeItem('auth-storage');
        window.location.href = '/login';
      }
      return Promise.reject(error);
    }
  );

  return client;
};

// 导出 API 客户端单例
export const apiClient = createApiClient();

// 导出 API 基础 URL
export { API_BASE_URL };