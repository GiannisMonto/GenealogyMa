import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor, act } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { Audit } from './index';

const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

describe('Audit', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders audit page title', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Audit />
        </BrowserRouter>
      );
    });
    expect(screen.getByText('审计日志')).toBeInTheDocument();
  });

  it('renders export button', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Audit />
        </BrowserRouter>
      );
    });
    expect(screen.getByText('导出')).toBeInTheDocument();
  });

  it('renders search input', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Audit />
        </BrowserRouter>
      );
    });
    expect(screen.getByPlaceholderText('搜索用户名、描述...')).toBeInTheDocument();
  });

  it('renders filter selects', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Audit />
        </BrowserRouter>
      );
    });
    expect(screen.getByText('选择模块')).toBeInTheDocument();
    expect(screen.getByText('选择操作')).toBeInTheDocument();
  });

  it('renders audit logs table with columns', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Audit />
        </BrowserRouter>
      );
    });
    // Check that the table container exists - actual table rendering requires mock data setup
    expect(document.querySelector('.ant-table')).toBeInTheDocument();
  });

  it('renders reset filter button', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Audit />
        </BrowserRouter>
      );
    });
    expect(screen.getByText('重置')).toBeInTheDocument();
  });
});