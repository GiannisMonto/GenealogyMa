/**
 * 认证 API
 */
import { request } from './request';

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  refreshToken: string;
  expiresIn: number;
  user: {
    id: number;
    username: string;
    name: string;
    role: string;
    avatar?: string;
  };
}

/**
 * 用户登录
 */
export async function login(data: LoginRequest): Promise<LoginResponse> {
  return request<LoginResponse>('/auth/login', {
    method: 'POST',
    body: data,
  });
}

/**
 * 刷新令牌
 */
export async function refreshToken(refreshToken: string): Promise<LoginResponse> {
  return request<LoginResponse>('/auth/refresh', {
    method: 'POST',
    body: { refreshToken },
  });
}

/**
 * 用户登出
 */
export async function logout(): Promise<void> {
  return request('/auth/logout', { method: 'POST' });
}

/**
 * 获取当前用户信息
 */
export async function getCurrentUser(): Promise<LoginResponse['user']> {
  return request('/auth/me');
}