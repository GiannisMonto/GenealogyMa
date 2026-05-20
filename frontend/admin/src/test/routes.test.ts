import { describe, it, expect } from 'vitest';
import { routes, getFlattenRoutes } from '@/routes/routes';

describe('路由配置', () => {
  it('应该有登录路由', () => {
    const loginRoute = routes.find(r => r.path === '/login');
    expect(loginRoute).toBeDefined();
    expect(loginRoute?.meta?.requiresAuth).toBe(false);
  });

  it('应该需要认证的路由', () => {
    const dashboardRoute = routes.find(r => r.path === '/dashboard');
    expect(dashboardRoute).toBeDefined();
    expect(dashboardRoute?.meta?.requiresAuth).toBe(true);
  });

  it('getFlattenRoutes 应该返回所有扁平化路由', () => {
    const flatRoutes = getFlattenRoutes(routes);
    expect(flatRoutes.length).toBeGreaterThan(0);
    expect(flatRoutes.some(r => r.path === '/login')).toBe(true);
    expect(flatRoutes.some(r => r.path === '/dashboard')).toBe(true);
  });

  it('成员路由应该需要特定权限', () => {
    const membersRoute = routes.find(r => r.path === '/members');
    expect(membersRoute?.meta?.permission).toBe('person:read');
  });

  it('404路由应该在最后', () => {
    const lastRoute = routes[routes.length - 1];
    expect(lastRoute?.path).toBe('*');
  });
});