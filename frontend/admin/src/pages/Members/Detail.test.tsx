import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { MembersDetail } from './Detail';

vi.mock('@shared/api/person', () => ({
  personApi: {
    getPerson: vi.fn().mockResolvedValue({
      id: 1,
      name: '张三',
      style_name: '字云峰',
      gender: '男',
      generation: 1,
      birth_order: '长',
      father_id: null,
      lineage_path: '1',
      detail_text: '测试生平',
      birth_time_text: '1900-01-01',
      death_time_text: '1970-01-01',
      birth_place: '浙江绍兴',
      burial_place: '浙江绍兴',
      son_count: 3,
      daughter_count: 2,
      adopted_heir_count: 0,
      total_children_count: 5,
      age: 70,
      is_alive: false,
      full_name: '张三',
      created_at: '2024-01-01 10:00:00',
      updated_at: '2024-01-01 10:00:00',
      spouses: [],
      children: [],
    }),
  },
}));

const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
    useParams: () => ({ id: '1' }),
  };
});

describe('MembersDetail', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders back button', async () => {
    render(
      <BrowserRouter>
        <MembersDetail />
      </BrowserRouter>
    );
    await waitFor(() => {
      expect(screen.getByText('返回')).toBeInTheDocument();
    });
  });

  it('renders edit button', async () => {
    render(
      <BrowserRouter>
        <MembersDetail />
      </BrowserRouter>
    );
    await waitFor(() => {
      expect(screen.getByText('编辑')).toBeInTheDocument();
    });
  });

  it('renders basic info card', async () => {
    render(
      <BrowserRouter>
        <MembersDetail />
      </BrowserRouter>
    );
    await waitFor(() => {
      expect(screen.getByText('基本信息')).toBeInTheDocument();
    });
  });

  it('renders family statistics card', async () => {
    render(
      <BrowserRouter>
        <MembersDetail />
      </BrowserRouter>
    );
    await waitFor(() => {
      expect(screen.getByText('家庭统计')).toBeInTheDocument();
    });
  });

  it('renders system info card', async () => {
    render(
      <BrowserRouter>
        <MembersDetail />
      </BrowserRouter>
    );
    await waitFor(() => {
      expect(screen.getByText('系统信息')).toBeInTheDocument();
    });
  });
});