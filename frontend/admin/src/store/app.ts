import { create } from 'zustand';

interface AppState {
  sidebarCollapsed: boolean;
  breadcrumbs: { path: string; name: string }[];
  toggleSidebar: () => void;
  setSidebarCollapsed: (collapsed: boolean) => void;
  setBreadcrumbs: (breadcrumbs: { path: string; name: string }[]) => void;
}

export const useAppStore = create<AppState>((set) => ({
  sidebarCollapsed: false,
  breadcrumbs: [],
  toggleSidebar: () => set((state) => ({ sidebarCollapsed: !state.sidebarCollapsed })),
  setSidebarCollapsed: (collapsed) => set({ sidebarCollapsed: collapsed }),
  setBreadcrumbs: (breadcrumbs) => set({ breadcrumbs }),
}));