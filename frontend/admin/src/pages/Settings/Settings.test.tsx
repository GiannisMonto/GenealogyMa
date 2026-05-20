import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor, act } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { Settings } from './index';

vi.mock('@/api/config', () => ({
  configApi: {
    getConfig: vi.fn().mockResolvedValue(null),
    updateConfig: vi.fn(),
  },
}));

describe('Settings', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders settings page title', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Settings />
        </BrowserRouter>
      );
    });
    expect(screen.getByText('系统设置')).toBeInTheDocument();
  });

  it('renders save button', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Settings />
        </BrowserRouter>
      );
    });
    expect(screen.getByText('保存设置')).toBeInTheDocument();
  });

  it('renders reset button', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Settings />
        </BrowserRouter>
      );
    });
    expect(screen.getByText('重置')).toBeInTheDocument();
  });

  it('renders basic info card', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Settings />
        </BrowserRouter>
      );
    });
    expect(screen.getByText('基本信息')).toBeInTheDocument();
  });

  it('renders feature switches card', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Settings />
        </BrowserRouter>
      );
    });
    expect(screen.getByText('功能开关')).toBeInTheDocument();
  });

  it('renders upload settings card', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Settings />
        </BrowserRouter>
      );
    });
    expect(screen.getByText('上传设置')).toBeInTheDocument();
  });

  it('renders log settings card', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Settings />
        </BrowserRouter>
      );
    });
    expect(screen.getByText('日志设置')).toBeInTheDocument();
  });
});