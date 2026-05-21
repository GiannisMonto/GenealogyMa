/**
 * 备份管理 API
 */
import { apiClient } from './client';

export interface BackupDTO {
  id: number;
  name: string;
  type: string;
  status: string;
  file_path: string;
  file_size: number;
  file_size_str: string;
  database: string;
  error_message?: string;
  started_at?: string;
  completed_at?: string;
  created_at: string;
  created_by: number;
}

export interface BackupListDTO {
  data: BackupDTO[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface CreateBackupRequest {
  name: string;
  type: string;
  database?: string;
  tables?: string[];
}

export interface RestoreBackupRequest {
  backup_id: number;
  target_db?: string;
}

export interface BackupFilterRequest {
  type?: string;
  status?: string;
  database?: string;
  start_date?: string;
  end_date?: string;
  page?: number;
  page_size?: number;
}

/**
 * 获取备份列表
 */
export async function getBackups(params?: BackupFilterRequest): Promise<BackupListDTO> {
  const response = await apiClient.get('/backups', { params });
  return response.data.data;
}

/**
 * 获取备份详情
 */
export async function getBackup(id: number): Promise<BackupDTO> {
  const response = await apiClient.get(`/backups/${id}`);
  return response.data.data;
}

/**
 * 创建备份
 */
export async function createBackup(data: CreateBackupRequest): Promise<BackupDTO> {
  const response = await apiClient.post('/backups', data);
  return response.data.data;
}

/**
 * 恢复备份
 */
export async function restoreBackup(id: number, data?: RestoreBackupRequest): Promise<void> {
  await apiClient.post(`/backups/${id}/restore`, data || { backup_id: id });
}

/**
 * 删除备份
 */
export async function deleteBackup(id: number): Promise<void> {
  await apiClient.delete(`/backups/${id}`);
}

/**
 * 清理旧备份
 */
export async function cleanOldBackups(days: number): Promise<{ deleted_count: number }> {
  const response = await apiClient.delete('/backups/clean-old', { data: { days } });
  return response.data.data;
}

export const backupApi = {
  getBackups,
  getBackup,
  createBackup,
  restoreBackup,
  deleteBackup,
  cleanOldBackups,
};