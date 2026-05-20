/**
 * 审计日志 API
 */
import { apiClient } from './client';

export interface AuditLogDTO {
  id: number;
  user_id: number;
  username: string;
  module: string;
  action: string;
  resource_type: string;
  resource_id: string;
  old_value?: string;
  new_value?: string;
  ip_address?: string;
  user_agent?: string;
  description?: string;
  created_at: string;
}

export interface AuditFilters {
  keyword?: string;
  module?: string;
  action?: string;
  start_date?: string;
  end_date?: string;
  page?: number;
  page_size?: number;
}

export interface PageResponse<T> {
  data: T[];
  total: number;
  page: number;
  page_size: number;
}

/**
 * 获取审计日志列表
 */
export async function getAuditLogs(
  filters?: AuditFilters
): Promise<PageResponse<AuditLogDTO>> {
  const params = new URLSearchParams();

  if (filters) {
    if (filters.keyword) params.append('keyword', filters.keyword);
    if (filters.module) params.append('module', filters.module);
    if (filters.action) params.append('action', filters.action);
    if (filters.start_date) params.append('start_date', filters.start_date);
    if (filters.end_date) params.append('end_date', filters.end_date);
    if (filters.page) params.append('page', String(filters.page));
    if (filters.page_size) params.append('page_size', String(filters.page_size));
  }

  const response = await apiClient.get(`/audit-logs?${params.toString()}`);
  return response.data.data;
}

/**
 * 获取审计日志详情
 */
export async function getAuditLogDetail(id: number): Promise<AuditLogDTO> {
  const response = await apiClient.get(`/audit-logs/${id}`);
  return response.data.data;
}

/**
 * 导出审计日志
 */
export async function exportAuditLogs(
  filters?: AuditFilters
): Promise<Blob> {
  const params = new URLSearchParams();

  if (filters) {
    if (filters.keyword) params.append('keyword', filters.keyword);
    if (filters.module) params.append('module', filters.module);
    if (filters.action) params.append('action', filters.action);
    if (filters.start_date) params.append('start_date', filters.start_date);
    if (filters.end_date) params.append('end_date', filters.end_date);
  }

  const response = await apiClient.get(`/audit-logs/export?${params.toString()}`, {
    responseType: 'blob',
  });
  return response.data;
}