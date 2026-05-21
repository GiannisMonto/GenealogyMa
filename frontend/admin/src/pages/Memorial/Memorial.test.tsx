import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor, act } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { Memorial } from './index';
import { memorialApi } from '@/api/memorial';

// Mock the API
vi.mock('@/api/memorial', () => ({
  memorialApi: {
    getHalls: vi.fn(),
    getHall: vi.fn(),
    getHallWithTablets: vi.fn(),
    createHall: vi.fn(),
    updateHall: vi.fn(),
    deleteHall: vi.fn(),
  },
}));

const mockHalls = [
  {
    id: 1,
    name: '陈氏宗祠',
    description: '陈氏家族宗祠',
    province: '广东省',
    city: '深圳市',
    district: '南山区',
    address: '南山路123号',
    latitude: 22.543,
    longitude: 114.057,
    build_year: 2000,
    style: '岭南风格',
    image_url: '',
    total_tablet: 100,
    used_tablet: 45,
    usage_rate: 45,
    created_at: '2024-01-01 10:00:00',
    updated_at: '2024-01-01 10:00:00',
    tablets: [],
  },
  {
    id: 2,
    name: '王氏宗祠',
    description: '王氏家族宗祠',
    province: '福建省',
    city: '福州市',
    district: '鼓楼区',
    address: '鼓楼路456号',
    latitude: 26.075,
    longitude: 119.296,
    build_year: 1995,
    style: '闽南风格',
    image_url: '',
    total_tablet: 200,
    used_tablet: 180,
    usage_rate: 90,
    created_at: '2024-01-02 10:00:00',
    updated_at: '2024-01-02 10:00:00',
    tablets: [],
  },
];

describe('Memorial', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(memorialApi.getHalls).mockResolvedValue(mockHalls);
  });

  it('renders memorial page title', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Memorial />
        </BrowserRouter>
      );
    });
    expect(screen.getByText('宗祠管理')).toBeInTheDocument();
  });

  it('renders add hall button', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Memorial />
        </BrowserRouter>
      );
    });
    expect(screen.getByText('新增宗祠')).toBeInTheDocument();
  });

  it('renders search input', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Memorial />
        </BrowserRouter>
      );
    });
    expect(screen.getByPlaceholderText('搜索宗祠名称或地区...')).toBeInTheDocument();
  });

  it('renders halls table with columns', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Memorial />
        </BrowserRouter>
      );
    });

    await waitFor(() => {
      expect(screen.getByText('宗祠名称')).toBeInTheDocument();
      expect(screen.getByText('位置')).toBeInTheDocument();
      expect(screen.getByText('牌位使用')).toBeInTheDocument();
      expect(screen.getByText('建筑风格')).toBeInTheDocument();
      expect(screen.getByText('操作')).toBeInTheDocument();
    });
  });

  it('renders hall data correctly', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Memorial />
        </BrowserRouter>
      );
    });

    await waitFor(() => {
      expect(screen.getByText('陈氏宗祠')).toBeInTheDocument();
      expect(screen.getByText('王氏宗祠')).toBeInTheDocument();
    });
  });

  it('shows pagination info', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Memorial />
        </BrowserRouter>
      );
    });

    await waitFor(() => {
      expect(screen.getByText('共 2 条')).toBeInTheDocument();
    });
  });
});