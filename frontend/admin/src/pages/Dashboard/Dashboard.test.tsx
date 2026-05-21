import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor, act } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { Dashboard } from './index';

vi.mock('@/api/dashboard', () => ({
  getDashboardStats: vi.fn().mockResolvedValue({
    totalPersons: 5950,
    totalUsers: 128,
    totalRoles: 6,
    auditLogs: 1024,
  }),
}));

describe('Dashboard', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders dashboard page title', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Dashboard />
        </BrowserRouter>
      );
    });
    expect(screen.getByText('控制台')).toBeInTheDocument();
  });

  it('renders total persons statistic', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Dashboard />
        </BrowserRouter>
      );
    });
    await waitFor(() => {
      expect(screen.getByText('总成员数')).toBeInTheDocument();
    });
  });

  it('renders total users statistic', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Dashboard />
        </BrowserRouter>
      );
    });
    await waitFor(() => {
      expect(screen.getByText('用户数')).toBeInTheDocument();
    });
  });

  it('renders total roles statistic', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Dashboard />
        </BrowserRouter>
      );
    });
    await waitFor(() => {
      expect(screen.getByText('角色数')).toBeInTheDocument();
    });
  });

  it('renders audit logs statistic', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Dashboard />
        </BrowserRouter>
      );
    });
    await waitFor(() => {
      expect(screen.getByText('审计日志')).toBeInTheDocument();
    });
  });

  it('displays correct total persons value', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Dashboard />
        </BrowserRouter>
      );
    });
    await waitFor(() => {
      expect(screen.getByText('5,950')).toBeInTheDocument();
    });
  });

  it('displays correct total users value', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Dashboard />
        </BrowserRouter>
      );
    });
    await waitFor(() => {
      expect(screen.getByText('128')).toBeInTheDocument();
    });
  });

  it('displays correct total roles value', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Dashboard />
        </BrowserRouter>
      );
    });
    await waitFor(() => {
      expect(screen.getByText('6')).toBeInTheDocument();
    });
  });

  it('displays correct audit logs value', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Dashboard />
        </BrowserRouter>
      );
    });
    await waitFor(() => {
      expect(screen.getByText('1,024')).toBeInTheDocument();
    });
  });

  it('renders four stat cards', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Dashboard />
        </BrowserRouter>
      );
    });
    await waitFor(() => {
      const cards = document.querySelectorAll('.ant-card');
      expect(cards.length).toBe(4);
    });
  });
});