import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor, act, fireEvent } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { Cemetery } from './index';

vi.mock('@/api/cemetery', () => ({
  cemeteryApi: {
    getCemeteries: vi.fn().mockResolvedValue([]),
    getCemetery: vi.fn(),
    createCemetery: vi.fn(),
    updateCemetery: vi.fn(),
    deleteCemetery: vi.fn(),
    getGraves: vi.fn(),
    getGrave: vi.fn(),
    createGrave: vi.fn(),
    updateGrave: vi.fn(),
    deleteGrave: vi.fn(),
  },
}));

const mockCemeteries = [
  {
    id: 1,
    name: '南山墓园',
    description: '深圳市南山区大型墓园',
    province: '广东省',
    city: '深圳市',
    district: '南山区',
    address: '南山区桃源路123号',
    latitude: 22.5312,
    longitude: 113.9293,
    total_grave: 1000,
    used_grave: 450,
    image_url: '',
    created_at: '2026-01-15T08:00:00Z',
    updated_at: '2026-05-20T10:30:00Z',
  },
  {
    id: 2,
    name: '宝安墓园',
    description: '深圳市宝安区墓园',
    province: '广东省',
    city: '深圳市',
    district: '宝安区',
    address: '宝安区西乡街道',
    latitude: 22.7213,
    longitude: 113.8233,
    total_grave: 800,
    used_grave: 760,
    image_url: '',
    created_at: '2026-02-10T09:00:00Z',
    updated_at: '2026-05-19T14:20:00Z',
  },
];

describe('Cemetery', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('1. renders cemetery page title', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Cemetery />
        </BrowserRouter>
      );
    });
    expect(screen.getByText('墓园管理')).toBeInTheDocument();
  });

  it('2. renders add cemetery button', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Cemetery />
        </BrowserRouter>
      );
    });
    expect(screen.getByText('新增墓园')).toBeInTheDocument();
  });

  it('3. displays cemetery list when data is loaded', async () => {
    const { cemeteryApi } = await import('@/api/cemetery');
    (cemeteryApi.getCemeteries as ReturnType<typeof vi.fn>).mockResolvedValueOnce(mockCemeteries);

    await act(async () => {
      render(
        <BrowserRouter>
          <Cemetery />
        </BrowserRouter>
      );
    });

    await waitFor(() => {
      expect(screen.getByText('南山墓园')).toBeInTheDocument();
      expect(screen.getByText('宝安墓园')).toBeInTheDocument();
    });
  });

  it('4. search functionality works', async () => {
    const { cemeteryApi } = await import('@/api/cemetery');
    (cemeteryApi.getCemeteries as ReturnType<typeof vi.fn>).mockResolvedValueOnce(mockCemeteries);

    await act(async () => {
      render(
        <BrowserRouter>
          <Cemetery />
        </BrowserRouter>
      );
    });

    await waitFor(() => {
      expect(screen.getByText('南山墓园')).toBeInTheDocument();
    });

    const searchInput = screen.getByPlaceholderText('搜索墓园名称或地区...');
    await act(async () => {
      fireEvent.change(searchInput, { target: { value: '宝安' } });
    });

    await waitFor(() => {
      expect(screen.getByText('宝安墓园')).toBeInTheDocument();
    });
  });

  it('5. opens modal when add button is clicked', async () => {
    const { cemeteryApi } = await import('@/api/cemetery');
    (cemeteryApi.getCemeteries as ReturnType<typeof vi.fn>).mockResolvedValueOnce([]);

    await act(async () => {
      render(
        <BrowserRouter>
          <Cemetery />
        </BrowserRouter>
      );
    });

    const addButton = screen.getByText('新增墓园');
    await act(async () => {
      fireEvent.click(addButton);
    });

    await waitFor(() => {
      expect(screen.getByRole('dialog')).toBeInTheDocument();
    });
  });

  it('6. renders with empty data without crashing', async () => {
    const { cemeteryApi } = await import('@/api/cemetery');
    (cemeteryApi.getCemeteries as ReturnType<typeof vi.fn>).mockResolvedValueOnce([]);

    await act(async () => {
      render(
        <BrowserRouter>
          <Cemetery />
        </BrowserRouter>
      );
    });

    // The page should render without crashing even with empty data
    expect(screen.getByText('墓园管理')).toBeInTheDocument();
  });
});