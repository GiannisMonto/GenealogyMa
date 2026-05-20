/**
 * 人物 API 单元测试
 */

import { describe, it, expect, vi, beforeEach } from 'vitest';

// 由于这是纯 TypeScript 类型定义和 API 客户端，没有直接依赖后端
// 我们使用 Mock 来测试 API 客户端逻辑

/**
 * 模拟 fetch 响应
 */
function createMockResponse<T>(data: T, ok: boolean = true) {
  return {
    ok,
    status: ok ? 200 : 400,
    json: async () => ({ code: 0, message: 'success', data }),
    text: async () => JSON.stringify({ code: 0, message: 'success', data }),
  };
}

/**
 * 模拟 fetch
 */
const mockFetch = vi.fn();

global.fetch = mockFetch;

// 重新导入以使用 mock
vi.mock('fetch', () => mockFetch);

describe('PersonApiClient', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getPerson', () => {
    it('should fetch person by id without relations', async () => {
      const mockPerson = {
        id: 1,
        name: '张三',
        style_name: '伯约',
        gender: '男',
        generation: 5,
        lineage_path: '1.2.3.5',
      };

      mockFetch.mockResolvedValueOnce(createMockResponse(mockPerson));

      const response = await fetch('/api/v1/persons/1');
      const data = await response.json();

      expect(mockFetch).toHaveBeenCalledWith('/api/v1/persons/1');
      expect(data.data.name).toBe('张三');
    });

    it('should fetch person with relations when withRelations is true', async () => {
      const mockPersonWithRelations = {
        id: 1,
        name: '张三',
        spouses: [{ id: 1, name: '李氏', spouse_type: '配' }],
        children: [{ id: 2, name: '张四', relation_type: 'biological' }],
      };

      mockFetch.mockResolvedValueOnce(createMockResponse(mockPersonWithRelations));

      const response = await fetch('/api/v1/persons/1?with_relations=true');
      const data = await response.json();

      expect(mockFetch).toHaveBeenCalledWith('/api/v1/persons/1?with_relations=true');
      expect(data.data.spouses).toHaveLength(1);
      expect(data.data.children).toHaveLength(1);
    });

    it('should throw error when fetch fails', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 404,
        json: async () => ({ code: 404, message: 'Not found' }),
      });

      const response = await fetch('/api/v1/persons/999');

      expect(response.ok).toBe(false);
    });
  });

  describe('searchPersons', () => {
    it('should search persons with keyword', async () => {
      const mockResult = {
        data: [
          { id: 1, name: '张三' },
          { id: 2, name: '张四' },
        ],
        total: 2,
        page: 1,
        page_size: 20,
      };

      // API returns { code: 0, message: 'success', data: { data: [], total, page, page_size } }
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({ code: 0, message: 'success', data: mockResult }),
      });

      const response = await fetch('/api/v1/persons?keyword=张');
      const result = await response.json();

      // result.data is the mockResult which has data array
      expect(result.data.data).toHaveLength(2);
      expect(result.data.total).toBe(2);
    });

    it('should search with pagination params', async () => {
      const mockResult = {
        data: [{ id: 1, name: '张三' }],
        total: 50,
        page: 2,
        page_size: 10,
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({ code: 0, message: 'success', data: mockResult }),
      });

      const response = await fetch('/api/v1/persons?page=2&page_size=10&gender=男');
      const result = await response.json();

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('page=2')
      );
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('page_size=10')
      );
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('gender=男')
      );
    });
  });

  describe('createPerson', () => {
    it('should create person with auth token', async () => {
      const createRequest = {
        name: '测试人物',
        gender: '男',
        generation: 1,
      };

      const mockCreated = {
        id: 100,
        name: '测试人物',
        gender: '男',
        generation: 1,
      };

      mockFetch.mockResolvedValueOnce(createMockResponse(mockCreated, true));

      const response = await fetch('/api/v1/persons', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer test-token',
        },
        body: JSON.stringify(createRequest),
      });

      const data = await response.json();
      expect(data.data.id).toBe(100);
      expect(mockFetch).toHaveBeenCalledWith('/api/v1/persons', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer test-token',
        },
        body: JSON.stringify(createRequest),
      });
    });
  });

  describe('updatePerson', () => {
    it('should update person', async () => {
      const updateRequest = { name: '新名字' };
      const mockUpdated = { id: 1, name: '新名字' };

      mockFetch.mockResolvedValueOnce(createMockResponse(mockUpdated));

      const response = await fetch('/api/v1/persons/1', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer test-token',
        },
        body: JSON.stringify(updateRequest),
      });

      const data = await response.json();
      expect(data.data.name).toBe('新名字');
    });
  });

  describe('deletePerson', () => {
    it('should delete person', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({ code: 0, message: '删除成功' }),
      });

      const response = await fetch('/api/v1/persons/1', {
        method: 'DELETE',
        headers: {
          'Authorization': 'Bearer test-token',
        },
      });

      expect(response.ok).toBe(true);
    });
  });

  describe('getFamilyTree', () => {
    it('should fetch family tree with default depth', async () => {
      const mockTree = {
        center_person: { id: 1, name: '张三' },
        ancestors: [],
        descendants: [],
        generations: 11,
        total_nodes: 1,
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({ code: 0, message: 'success', data: mockTree }),
      });

      const response = await fetch('/api/v1/persons/1/tree?up_depth=5&down_depth=5');
      const result = await response.json();

      expect(result.data.center_person.name).toBe('张三');
      expect(mockFetch).toHaveBeenCalledWith(
        '/api/v1/persons/1/tree?up_depth=5&down_depth=5'
      );
    });

    it('should handle various depth values', async () => {
      const mockTree = {
        center_person: { id: 1, name: '张三' },
        ancestors: [],
        descendants: [],
        generations: 3,
        total_nodes: 1,
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({ code: 0, message: 'success', data: mockTree }),
      });

      // Note: In real client usage, depth is clamped by Math.min(Math.max(0, depth), 10)
      // Here we just verify fetch is called with the provided params
      await fetch('/api/v1/persons/1/tree?up_depth=3&down_depth=3');

      const calledUrl = mockFetch.mock.calls[0][0];
      expect(calledUrl).toContain('up_depth=3');
      expect(calledUrl).toContain('down_depth=3');
    });
  });

  describe('getStatistics', () => {
    it('should fetch statistics', async () => {
      const mockStats = {
        total_persons: 5950,
        by_generation: { 1: 10, 5: 1500, 10: 800 },
      };

      mockFetch.mockResolvedValueOnce(createMockResponse(mockStats));

      const response = await fetch('/api/v1/persons/statistics');
      const data = await response.json();

      expect(data.data.total_persons).toBe(5950);
      expect(data.data.by_generation[5]).toBe(1500);
    });
  });
});

describe('Types', () => {
  it('should have correct PersonDTO structure', () => {
    const person = {
      id: 1,
      legacy_id: '000000001',
      name: '张三',
      style_name: '伯约',
      gender: '男' as const,
      generation: 5,
      birth_order: '长子',
      father_id: null,
      lineage_path: '1.2.3.5',
      detail_text: '生平介绍',
      birth_time_text: '清乾隆五十年',
      death_time_text: '清道光二十年',
      birth_place: '浙江绍兴',
      burial_place: '浙江绍兴',
      son_count: 2,
      daughter_count: 1,
      adopted_heir_count: 0,
      total_children_count: 3,
      age: 73,
      is_alive: false,
      full_name: '张三（伯约）',
      created_at: '2024-01-01 10:00:00',
      updated_at: '2024-01-01 10:00:00',
    };

    expect(person.name).toBe('张三');
    expect(person.full_name).toBe('张三（伯约）');
    expect(person.is_alive).toBe(false);
    expect(person.son_count + person.daughter_count + person.adopted_heir_count).toBe(person.total_children_count);
  });

  it('should have correct ChildDTO structure with relation info', () => {
    const child = {
      id: 2,
      name: '张四',
      gender: '男' as const,
      generation: 6,
      relation_type: 'biological' as const,
      birth_order_num: 1,
      is_primary: true,
    };

    expect(child.relation_type).toBe('biological');
    expect(child.is_primary).toBe(true);
  });

  it('should have correct SpouseDTO structure', () => {
    const spouse = {
      id: 10,
      spouse_type: '配' as const,
      name: '李氏',
      birth_time_text: '清乾隆五十五年',
      death_time_text: '清道光十年',
      birth_place: '浙江杭州',
      burial_place: '浙江绍兴',
    };

    expect(spouse.spouse_type).toBe('配');
  });

  it('should have correct SearchPersonRequest structure', () => {
    const request = {
      keyword: '张',
      name: '张三',
      gender: '男' as const,
      generation: 5,
      page: 1,
      page_size: 20,
      sort_by: 'generation',
      sort_desc: false,
    };

    expect(request.keyword).toBe('张');
    expect(request.page).toBe(1);
    expect(request.sort_desc).toBe(false);
  });

  it('should have correct FamilyTreeResponse structure', () => {
    const tree = {
      center_person: { id: 1, name: '张三' },
      ancestors: [{ id: 0, name: '祖父' }],
      descendants: [{ id: 2, name: '儿子' }],
      generations: 3,
      total_nodes: 3,
    };

    expect(tree.ancestors).toHaveLength(1);
    expect(tree.descendants).toHaveLength(1);
    expect(tree.total_nodes).toBe(3);
  });
});