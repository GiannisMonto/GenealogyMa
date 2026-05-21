import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor, act, fireEvent } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { Backup } from './index';

vi.mock('@/api/backup', () => ({
  backupApi: {
    getBackups: vi.fn().mockResolvedValue({ data: [], total: 0 }),
    getBackup: vi.fn(),
    createBackup: vi.fn(),
    restoreBackup: vi.fn(),
    deleteBackup: vi.fn(),
    cleanOldBackups: vi.fn(),
  },
}));

const mockBackups = [
  {
    id: 1,
    name: '全量备份-20260520',
    type: 'full',
    status: 'completed',
    file_path: '/backups/backup_20260520.sql',
    file_size: 104857600,
    file_size_str: '100 MB',
    database: 'genealogy',
    error_message: '',
    started_at: '2026-05-20T10:00:00Z',
    completed_at: '2026-05-20T10:15:00Z',
    created_at: '2026-05-20T10:00:00Z',
    created_by: 1,
  },
  {
    id: 2,
    name: '增量备份-20260521',
    type: 'incremental',
    status: 'running',
    file_path: '',
    file_size: 0,
    file_size_str: '',
    database: 'genealogy',
    error_message: '',
    started_at: '2026-05-21T10:00:00Z',
    completed_at: '',
    created_at: '2026-05-21T10:00:00Z',
    created_by: 1,
  },
  {
    id: 3,
    name: '差异备份-20260519',
    type: 'differential',
    status: 'failed',
    file_path: '',
    file_size: 0,
    file_size_str: '',
    database: 'genealogy',
    error_message: 'Connection timeout',
    started_at: '2026-05-19T10:00:00Z',
    completed_at: '2026-05-19T10:05:00Z',
    created_at: '2026-05-19T10:00:00Z',
    created_by: 1,
  },
];

describe('Backup', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('1. renders backup page title', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Backup />
        </BrowserRouter>
      );
    });
    expect(screen.getByText('备份管理')).toBeInTheDocument();
  });

  it('2. renders create backup button', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Backup />
        </BrowserRouter>
      );
    });
    expect(screen.getByText('创建备份')).toBeInTheDocument();
  });

  it('3. renders statistics cards', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Backup />
        </BrowserRouter>
      );
    });
    expect(screen.getByText('总备份数')).toBeInTheDocument();
    expect(screen.getByText('已完成')).toBeInTheDocument();
    expect(screen.getByText('进行中')).toBeInTheDocument();
    expect(screen.getByText('失败')).toBeInTheDocument();
  });

  it('4. displays backup list when data is loaded', async () => {
    const { backupApi } = await import('@/api/backup');
    (backupApi.getBackups as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      data: mockBackups,
      total: 3,
    });

    await act(async () => {
      render(
        <BrowserRouter>
          <Backup />
        </BrowserRouter>
      );
    });

    await waitFor(() => {
      expect(screen.getByText('全量备份-20260520')).toBeInTheDocument();
      expect(screen.getByText('增量备份-20260521')).toBeInTheDocument();
    });
  });

  it('5. search functionality works', async () => {
    const { backupApi } = await import('@/api/backup');
    (backupApi.getBackups as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      data: mockBackups,
      total: 3,
    });

    await act(async () => {
      render(
        <BrowserRouter>
          <Backup />
        </BrowserRouter>
      );
    });

    await waitFor(() => {
      expect(screen.getByText('全量备份-20260520')).toBeInTheDocument();
    });

    const searchInput = screen.getByPlaceholderText('搜索备份名称或数据库...');
    await act(async () => {
      fireEvent.change(searchInput, { target: { value: '差异' } });
    });

    await waitFor(() => {
      expect(screen.getByText('差异备份-20260519')).toBeInTheDocument();
    });
  });

  it('6. opens modal when create button is clicked', async () => {
    const { backupApi } = await import('@/api/backup');
    (backupApi.getBackups as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      data: [],
      total: 0,
    });

    await act(async () => {
      render(
        <BrowserRouter>
          <Backup />
        </BrowserRouter>
      );
    });

    const createButton = screen.getByText('创建备份');
    await act(async () => {
      fireEvent.click(createButton);
    });

    await waitFor(() => {
      expect(screen.getByRole('dialog')).toBeInTheDocument();
    });
  });

  it('7. opens drawer when view button is clicked', async () => {
    const { backupApi } = await import('@/api/backup');
    (backupApi.getBackups as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      data: mockBackups,
      total: 3,
    });

    await act(async () => {
      render(
        <BrowserRouter>
          <Backup />
        </BrowserRouter>
      );
    });

    await waitFor(() => {
      expect(screen.getByText('全量备份-20260520')).toBeInTheDocument();
    });

    const viewButton = screen.getAllByText('详情')[0];
    await act(async () => {
      fireEvent.click(viewButton);
    });

    await waitFor(() => {
      expect(screen.getByText('备份详情')).toBeInTheDocument();
    });
  });

  it('8. renders with empty data without crashing', async () => {
    const { backupApi } = await import('@/api/backup');
    (backupApi.getBackups as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      data: [],
      total: 0,
    });

    await act(async () => {
      render(
        <BrowserRouter>
          <Backup />
        </BrowserRouter>
      );
    });

    expect(screen.getByText('备份管理')).toBeInTheDocument();
    expect(screen.getByText('总备份数')).toBeInTheDocument();
  });
});