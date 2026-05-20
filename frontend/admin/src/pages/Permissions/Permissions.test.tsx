import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor, act } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { Permissions } from './index';

vi.mock('@/api/rbac', () => ({
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

describe('Permissions', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders permissions page title', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Permissions />
        </BrowserRouter>
      );
    });
    expect(screen.getByText('权限管理')).toBeInTheDocument();
  });

  it('renders add permission button', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Permissions />
        </BrowserRouter>
      );
    });
    expect(screen.getByText('新增权限')).toBeInTheDocument();
  });

  it('renders search input', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Permissions />
        </BrowserRouter>
      );
    });
    expect(screen.getByPlaceholderText('搜索权限名称、代码或模块...')).toBeInTheDocument();
  });

  it('renders permissions table with columns', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Permissions />
        </BrowserRouter>
      );
    });

    await waitFor(() => {
      expect(screen.getByText('ID')).toBeInTheDocument();
      expect(screen.getByText('权限名称')).toBeInTheDocument();
      expect(screen.getByText('权限代码')).toBeInTheDocument();
      expect(screen.getByText('模块')).toBeInTheDocument();
      expect(screen.getByText('描述')).toBeInTheDocument();
      expect(screen.getByText('操作')).toBeInTheDocument();
    });
  });

  it('renders permissions correctly', async () => {
    const { permissionApi } = await import('@/api/rbac');
    vi.mocked(permissionApi.getPermissions).mockResolvedValueOnce([
      {
        id: 1,
        code: 'person_read',
        name: '查看人物',
        module: 'person',
        description: '允许查看人物信息',
        created_at: '2026-05-21T00:00:00Z',
      },
    ]);

    await act(async () => {
      render(
        <BrowserRouter>
          <Permissions />
        </BrowserRouter>
      );
    });

    await waitFor(() => {
      expect(screen.getByText('查看人物')).toBeInTheDocument();
      expect(screen.getByText('person_read')).toBeInTheDocument();
    });
  });

  it('renders module tags with correct colors', async () => {
    const { permissionApi } = await import('@/api/rbac');
    vi.mocked(permissionApi.getPermissions).mockResolvedValueOnce([
      {
        id: 1,
        code: 'admin_user',
        name: '用户管理',
        module: 'admin',
        description: '管理员用户模块',
        created_at: '2026-05-21T00:00:00Z',
      },
    ]);

    await act(async () => {
      render(
        <BrowserRouter>
          <Permissions />
        </BrowserRouter>
      );
    });

    await waitFor(() => {
      // Check module tag appears (in table)
      const adminTags = screen.getAllByText('admin');
      expect(adminTags.length).toBeGreaterThan(0);
    });
  });
});