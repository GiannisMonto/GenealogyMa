import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useAuthStore } from './auth';

// Mock axios
vi.mock('axios', () => {
  const mockAxios = {
    create: vi.fn(() => ({
      interceptors: {
        request: { use: vi.fn((cb) => cb) },
        response: { use: vi.fn((cb) => cb) },
      },
      post: vi.fn(),
      get: vi.fn(),
    })),
  };
  return { default: mockAxios, ...mockAxios };
});

describe('useAuthStore', () => {
  beforeEach(() => {
    // Reset store state before each test
    useAuthStore.setState({
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: false,
    });
  });

  describe('initial state', () => {
    it('should have null user initially', () => {
      const { user } = useAuthStore.getState();
      expect(user).toBeNull();
    });

    it('should have null token initially', () => {
      const { token } = useAuthStore.getState();
      expect(token).toBeNull();
    });

    it('should not be authenticated initially', () => {
      const { isAuthenticated } = useAuthStore.getState();
      expect(isAuthenticated).toBe(false);
    });

    it('should not be loading initially', () => {
      const { isLoading } = useAuthStore.getState();
      expect(isLoading).toBe(false);
    });
  });

  describe('hasPermission', () => {
    it('should return false when user is null', () => {
      const hasPermission = useAuthStore.getState().hasPermission('person:read');
      expect(hasPermission).toBe(false);
    });

    it('should return true when user has the specific permission', () => {
      useAuthStore.setState({
        user: {
          id: '1',
          username: 'admin',
          email: 'admin@example.com',
          role: 'admin',
          permissions: ['person:read', 'person:write'],
        },
      });

      expect(useAuthStore.getState().hasPermission('person:read')).toBe(true);
      expect(useAuthStore.getState().hasPermission('person:write')).toBe(true);
    });

    it('should return false when user does not have the permission', () => {
      useAuthStore.setState({
        user: {
          id: '1',
          username: 'admin',
          email: 'admin@example.com',
          role: 'admin',
          permissions: ['person:read'],
        },
      });

      expect(useAuthStore.getState().hasPermission('person:delete')).toBe(false);
    });

    it('should return true when user has admin:* permission', () => {
      useAuthStore.setState({
        user: {
          id: '1',
          username: 'superadmin',
          email: 'superadmin@example.com',
          role: 'super_admin',
          permissions: ['admin:*'],
        },
      });

      expect(useAuthStore.getState().hasPermission('person:read')).toBe(true);
      expect(useAuthStore.getState().hasPermission('person:write')).toBe(true);
      expect(useAuthStore.getState().hasPermission('person:delete')).toBe(true);
      expect(useAuthStore.getState().hasPermission('any:permission')).toBe(true);
    });
  });

  describe('logout', () => {
    it('should clear user, token and isAuthenticated', async () => {
      // First set some state
      useAuthStore.setState({
        user: {
          id: '1',
          username: 'admin',
          email: 'admin@example.com',
          role: 'admin',
          permissions: ['person:read'],
        },
        token: 'some-token',
        isAuthenticated: true,
      });

      // Perform logout
      await useAuthStore.getState().logout();

      // Verify state is cleared
      const { user, token, isAuthenticated } = useAuthStore.getState();
      expect(user).toBeNull();
      expect(token).toBeNull();
      expect(isAuthenticated).toBe(false);
    });
  });
});