/**
 * API 配置
 */
const API_BASE_URL = process.env.API_BASE_URL || 'http://localhost:8080/api/v1';

interface RequestOptions {
  headers?: Record<string, string>;
}

/**
 * 通用请求方法
 */
async function request<T>(
  url: string,
  options: RequestOptions & { method?: string; body?: unknown } = {}
): Promise<T> {
  const token = uni.getStorageSync('token');

  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...options.headers,
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const response = await uni.request({
    url: `${API_BASE_URL}${url}`,
    method: options.method || 'GET',
    header: headers,
    data: options.body,
  });

  if ((response.statusCode as number) >= 400) {
    throw new Error(`请求失败: ${response.statusCode}`);
  }

  const data = response.data as { data?: T; code?: number; message?: string };
  if (data.code !== 0 && data.code !== 200) {
    throw new Error(data.message || '请求失败');
  }

  return data.data as T;
}

export { API_BASE_URL, request };
export type { RequestOptions };