import { describe, it, expect } from 'vitest';
import { useAppStore } from './app';

describe('useAppStore', () => {
  beforeEach(() => {
    // Reset store state before each test
    useAppStore.setState({
      sidebarCollapsed: false,
      breadcrumbs: [],
    });
  });

  describe('initial state', () => {
    it('should have sidebar collapsed as false initially', () => {
      const { sidebarCollapsed } = useAppStore.getState();
      expect(sidebarCollapsed).toBe(false);
    });

    it('should have empty breadcrumbs initially', () => {
      const { breadcrumbs } = useAppStore.getState();
      expect(breadcrumbs).toEqual([]);
    });
  });

  describe('toggleSidebar', () => {
    it('should toggle sidebarCollapsed from false to true', () => {
      useAppStore.getState().toggleSidebar();
      expect(useAppStore.getState().sidebarCollapsed).toBe(true);
    });

    it('should toggle sidebarCollapsed from true to false', () => {
      useAppStore.setState({ sidebarCollapsed: true });
      useAppStore.getState().toggleSidebar();
      expect(useAppStore.getState().sidebarCollapsed).toBe(false);
    });
  });

  describe('setSidebarCollapsed', () => {
    it('should set sidebarCollapsed to true', () => {
      useAppStore.getState().setSidebarCollapsed(true);
      expect(useAppStore.getState().sidebarCollapsed).toBe(true);
    });

    it('should set sidebarCollapsed to false', () => {
      useAppStore.setState({ sidebarCollapsed: true });
      useAppStore.getState().setSidebarCollapsed(false);
      expect(useAppStore.getState().sidebarCollapsed).toBe(false);
    });
  });

  describe('setBreadcrumbs', () => {
    it('should set breadcrumbs', () => {
      const crumbs = [
        { path: '/dashboard', name: '仪表盘' },
        { path: '/members', name: '成员管理' },
      ];
      useAppStore.getState().setBreadcrumbs(crumbs);
      expect(useAppStore.getState().breadcrumbs).toEqual(crumbs);
    });

    it('should set empty breadcrumbs', () => {
      useAppStore.setState({
        breadcrumbs: [{ path: '/dashboard', name: '仪表盘' }],
      });
      useAppStore.getState().setBreadcrumbs([]);
      expect(useAppStore.getState().breadcrumbs).toEqual([]);
    });
  });
});