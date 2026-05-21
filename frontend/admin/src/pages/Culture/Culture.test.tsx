import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor, act, fireEvent } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { Culture } from './index';

vi.mock('@/api/culture', () => ({
  cultureApi: {
    getDocuments: vi.fn().mockResolvedValue([]),
    getStories: vi.fn().mockResolvedValue([]),
    getFamilyTeachings: vi.fn().mockResolvedValue([]),
    getDocument: vi.fn(),
    createDocument: vi.fn(),
    updateDocument: vi.fn(),
    deleteDocument: vi.fn(),
    getStory: vi.fn(),
    createStory: vi.fn(),
    updateStory: vi.fn(),
    deleteStory: vi.fn(),
    getFamilyTeaching: vi.fn(),
    createFamilyTeaching: vi.fn(),
    updateFamilyTeaching: vi.fn(),
    deleteFamilyTeaching: vi.fn(),
  },
}));

const mockDocuments = [
  {
    id: 1,
    title: '族谱序',
    content: '族谱序内容',
    category: 'genealogy',
    author: '张三',
    created_year: 1800,
    dynasty: '清朝',
    source: '',
    image_urls: [],
    view_count: 10,
    collect_count: 2,
    created_at: '2026-01-15T08:00:00Z',
    updated_at: '2026-05-20T10:30:00Z',
  },
  {
    id: 2,
    title: '家训原文',
    content: '家训内容',
    category: 'classic',
    author: '李四',
    created_year: null,
    dynasty: '明朝',
    source: '',
    image_urls: [],
    view_count: 5,
    collect_count: 1,
    created_at: '2026-02-10T09:00:00Z',
    updated_at: '2026-05-19T14:20:00Z',
  },
];

describe('Culture', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('1. renders culture page title', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Culture />
        </BrowserRouter>
      );
    });
    expect(screen.getByText('文化管理')).toBeInTheDocument();
  });

  it('2. renders add button', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Culture />
        </BrowserRouter>
      );
    });
    expect(screen.getByText('新增文献')).toBeInTheDocument();
  });

  it('3. displays document list when data is loaded', async () => {
    const { cultureApi } = await import('@/api/culture');
    (cultureApi.getDocuments as ReturnType<typeof vi.fn>).mockResolvedValueOnce(mockDocuments);

    await act(async () => {
      render(
        <BrowserRouter>
          <Culture />
        </BrowserRouter>
      );
    });

    await waitFor(() => {
      expect(screen.getByText('族谱序')).toBeInTheDocument();
      expect(screen.getByText('家训原文')).toBeInTheDocument();
    });
  });

  it('4. opens modal when add button is clicked', async () => {
    const { cultureApi } = await import('@/api/culture');
    (cultureApi.getDocuments as ReturnType<typeof vi.fn>).mockResolvedValueOnce([]);

    await act(async () => {
      render(
        <BrowserRouter>
          <Culture />
        </BrowserRouter>
      );
    });

    const addButton = screen.getByText('新增文献');
    await act(async () => {
      fireEvent.click(addButton);
    });

    await waitFor(() => {
      expect(screen.getByRole('dialog')).toBeInTheDocument();
    });
  });

  it('5. renders with empty data without crashing', async () => {
    const { cultureApi } = await import('@/api/culture');
    (cultureApi.getDocuments as ReturnType<typeof vi.fn>).mockResolvedValueOnce([]);

    await act(async () => {
      render(
        <BrowserRouter>
          <Culture />
        </BrowserRouter>
      );
    });

    expect(screen.getByText('文化管理')).toBeInTheDocument();
  });

  it('6. displays tabs for documents, stories, and teachings', async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Culture />
        </BrowserRouter>
      );
    });

    expect(screen.getByText('文献')).toBeInTheDocument();
    expect(screen.getByText('故事')).toBeInTheDocument();
    expect(screen.getByText('家训')).toBeInTheDocument();
  });
});
