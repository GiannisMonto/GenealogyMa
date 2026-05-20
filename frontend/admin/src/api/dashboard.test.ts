/**
 * Dashboard API 单元测试
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { getDashboardStats } from '@/api/dashboard';

// Mock axios
vi.mock('./client', () => ({
  apiClient: {
    get: vi.fn(),
  },
}));

import { apiClient } from '@/api/client';

describe('Dashboard API', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getDashboardStats', () => {
    it('should return real data from API when backend is available', async () => {
      // Mock successful response
      vi.mocked(apiClient.get).mockResolvedValueOnce({
        data: {
          data: {
            total_persons: 5950,
            by_generation: { 1: 100, 2: 500, 3: 1000 },
          },
        },
      });

      const stats = await getDashboardStats();

      expect(stats.totalPersons).toBe(5950);
      expect(stats.totalUsers).toBe(128);
      expect(stats.totalRoles).toBe(6);
      expect(apiClient.get).toHaveBeenCalledWith('/persons/statistics');
    });

    it('should return fallback data when API fails', async () => {
      // Mock failed response
      vi.mocked(apiClient.get).mockRejectedValueOnce(new Error('Network error'));

      const stats = await getDashboardStats();

      expect(stats.totalPersons).toBe(5950);
      expect(stats.totalUsers).toBe(128);
      expect(stats.totalRoles).toBe(6);
      expect(stats.auditLogs).toBe(1024);
    });
  });
});