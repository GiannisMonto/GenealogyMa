import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { Members } from './index';

vi.mock('@shared/api/person', () => ({
  personApi: {
    searchPersons: vi.fn().mockResolvedValue({
      data: [],
      total: 0,
      page: 1,
      pageSize: 20,
    }),
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

describe('Members', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders members page title', () => {
    render(
      <BrowserRouter>
        <Members />
      </BrowserRouter>
    );
    expect(screen.getByText('成员管理')).toBeInTheDocument();
  });

  it('renders add member button', () => {
    render(
      <BrowserRouter>
        <Members />
      </BrowserRouter>
    );
    expect(screen.getByText('新增成员')).toBeInTheDocument();
  });

  it('renders search input', () => {
    render(
      <BrowserRouter>
        <Members />
      </BrowserRouter>
    );
    expect(screen.getByPlaceholderText('搜索姓名、字号...')).toBeInTheDocument();
  });

  it('renders reset button', () => {
    render(
      <BrowserRouter>
        <Members />
      </BrowserRouter>
    );
    expect(screen.getByText('重置')).toBeInTheDocument();
  });
});