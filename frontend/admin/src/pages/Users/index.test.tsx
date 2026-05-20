import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor, act } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { Users } from './index';

vi.mock('@/api/user', () => ({
  userApi: {
    getUsers: vi.fn().mockResolvedValue([]),
    getUser: vi.fn(),
    createUser: vi.fn(),
    updateUser: vi.fn(),
    deleteUser: vi.fn(),
    resetPassword: vi.fn(),
    updateStatus: vi.fn(),
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

describe('Users', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders users page title', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Users />
        </BrowserRouter>
      );
    });
    expect(screen.getByText('用户管理')).toBeInTheDocument();
  });

  it('renders add user button', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Users />
        </BrowserRouter>
      );
    });
    expect(screen.getByText('新增用户')).toBeInTheDocument();
  });

  it('renders search input', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Users />
        </BrowserRouter>
      );
    });
    expect(screen.getByPlaceholderText('搜索用户名、邮箱或昵称...')).toBeInTheDocument();
  });

  it('renders users table with columns', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Users />
        </BrowserRouter>
      );
    });

    await waitFor(() => {
      expect(screen.getByText('ID')).toBeInTheDocument();
      expect(screen.getByText('用户')).toBeInTheDocument();
      expect(screen.getByText('邮箱')).toBeInTheDocument();
      expect(screen.getByText('状态')).toBeInTheDocument();
      expect(screen.getByText('角色')).toBeInTheDocument();
      expect(screen.getByText('最后登录')).toBeInTheDocument();
      expect(screen.getByText('操作')).toBeInTheDocument();
    });
  });

  it('renders users correctly', async () => {
    const { userApi } = await import('@/api/user');
    vi.mocked(userApi.getUsers).mockResolvedValueOnce([
      {
        id: 1,
        username: 'admin',
        email: 'admin@example.com',
        display_name: '管理员',
        avatar_url: '',
        status: 'active' as const,
        last_login_at: '2026-05-21T00:00:00Z',
        created_at: '2026-05-21T00:00:00Z',
        updated_at: '2026-05-21T00:00:00Z',
        roles: [{ id: 1, code: 'super_admin', name: '超级管理员', is_system: true, created_at: '', description: '' }],
      },
    ]);

    await act(async () => {
      render(
        <BrowserRouter>
          <Users />
        </BrowserRouter>
      );
    });

    await waitFor(() => {
      expect(screen.getByText('管理员')).toBeInTheDocument();
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
      expect(screen.getByText('正常')).toBeInTheDocument();
      expect(screen.getByText('超级管理员')).toBeInTheDocument();
    });
  });

  it('renders status tags correctly', async () => {
    const { userApi } = await import('@/api/user');
    vi.mocked(userApi.getUsers).mockResolvedValueOnce([
      {
        id: 1,
        username: 'user1',
        email: 'user1@example.com',
        display_name: '用户1',
        avatar_url: '',
        status: 'active' as const,
        last_login_at: null,
        created_at: '2026-05-21T00:00:00Z',
        updated_at: '2026-05-21T00:00:00Z',
        roles: [],
      },
      {
        id: 2,
        username: 'user2',
        email: 'user2@example.com',
        display_name: '用户2',
        avatar_url: '',
        status: 'inactive' as const,
        last_login_at: null,
        created_at: '2026-05-21T00:00:00Z',
        updated_at: '2026-05-21T00:00:00Z',
        roles: [],
      },
      {
        id: 3,
        username: 'user3',
        email: 'user3@example.com',
        display_name: '用户3',
        avatar_url: '',
        status: 'banned' as const,
        last_login_at: null,
        created_at: '2026-05-21T00:00:00Z',
        updated_at: '2026-05-21T00:00:00Z',
        roles: [],
      },
    ]);

    await act(async () => {
      render(
        <BrowserRouter>
          <Users />
        </BrowserRouter>
      );
    });

    await waitFor(() => {
      expect(screen.getByText('正常')).toBeInTheDocument();
      expect(screen.getByText('未激活')).toBeInTheDocument();
      expect(screen.getByText('已禁用')).toBeInTheDocument();
    });
  });
});