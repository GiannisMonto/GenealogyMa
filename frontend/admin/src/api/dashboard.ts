/**
 * 成员统计 API
 */
import { apiClient } from './client';
import type { StatisticsResponse } from '@shared/api/types';

export interface DashboardStats {
  totalPersons: number;
  totalUsers: number;
  totalRoles: number;
  auditLogs: number;
}

/**
 * 获取仪表板统计数据
 */
export async function getDashboardStats(): Promise<DashboardStats> {
  try {
    // 尝试从后端获取真实数据
    const response = await apiClient.get('/persons/statistics');
    const stats: StatisticsResponse = response.data.data;

    return {
      totalPersons: stats.total_persons,
      totalUsers: 128, // TODO: 从用户服务获取
      totalRoles: 6,
      auditLogs: 1024, // TODO: 从审计服务获取
    };
  } catch {
    // 后端不可用时返回默认值
    return {
      totalPersons: 5950,
      totalUsers: 128,
      totalRoles: 6,
      auditLogs: 1024,
    };
  }
}