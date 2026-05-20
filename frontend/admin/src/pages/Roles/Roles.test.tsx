import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor, act } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { Roles } from './index';

vi.mock('@/api/rbac', () => ({
  roleApi: {
    getRoles: vi.fn().mockResolvedValue([]),
    getRole: vi.fn(),
    createRole: vi.fn(),
    updateRole: vi.fn(),
    deleteRole: vi.fn(),
    getRolePermissions: vi.fn(),
    assignPermissions: vi.fn(),
  },
  permissionApi: {
    getPermissions: vi.fn().mockResolvedValue([]),
    createPermission: vi.fn(),
    updatePermission: vi.fn(),
    deletePermission: vi.fn(),
  },
}));

const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

describe('Roles', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders roles page title', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Roles />
        </BrowserRouter>
      );
    });
    expect(screen.getByText('角色管理')).toBeInTheDocument();
  });

  it('renders add role button', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Roles />
        </BrowserRouter>
      );
    });
    expect(screen.getByText('新增角色')).toBeInTheDocument();
  });

  it('renders search input', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Roles />
        </BrowserRouter>
      );
    });
    expect(screen.getByPlaceholderText('搜索角色名称或代码...')).toBeInTheDocument();
  });

  it('renders roles table with columns', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Roles />
        </BrowserRouter>
      );
    });

    await waitFor(() => {
      expect(screen.getByText('ID')).toBeInTheDocument();
      expect(screen.getByText('角色代码')).toBeInTheDocument();
      expect(screen.getByText('角色名称')).toBeInTheDocument();
      expect(screen.getByText('描述')).toBeInTheDocument();
      expect(screen.getByText('系统角色')).toBeInTheDocument();
      expect(screen.getByText('操作')).toBeInTheDocument();
    });
  });

  it('renders system roles correctly', async () => {
    const { roleApi } = await import('@/api/rbac');
    vi.mocked(roleApi.getRoles).mockResolvedValueOnce([
      {
        id: 1,
        code: 'super_admin',
        name: '超级管理员',
        description: '系统超级管理员',
        is_system: true,
        created_at: '2026-05-21T00:00:00Z',
        updated_at: '2026-05-21T00:00:00Z',
      },
    ]);

    await act(async () => {
      render(
        <BrowserRouter>
          <Roles />
        </BrowserRouter>
      );
    });

    await waitFor(() => {
      expect(screen.getByText('超级管理员')).toBeInTheDocument();
      expect(screen.getByText('super_admin')).toBeInTheDocument();
      expect(screen.getByText('是')).toBeInTheDocument();
    });
  });

  it('renders permission button for each role', async () => {
    const { roleApi } = await import('@/api/rbac');
    vi.mocked(roleApi.getRoles).mockResolvedValueOnce([
      {
        id: 1,
        code: 'admin',
        name: '管理员',
        description: '系统管理员',
        is_system: true,
        created_at: '2026-05-21T00:00:00Z',
        updated_at: '2026-05-21T00:00:00Z',
      },
    ]);

    await act(async () => {
      render(
        <BrowserRouter>
          <Roles />
        </BrowserRouter>
      );
    });

    await waitFor(() => {
      const permissionButtons = screen.getAllByText('权限');
      expect(permissionButtons.length).toBeGreaterThan(0);
    });
  });
});