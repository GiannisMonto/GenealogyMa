/**
 * 人物 API 客户端
 * 人物 CRUD、搜索、族谱树查询等 API 封装
 */

import type {
  PersonDTO,
  SpouseDTO,
  ChildDTO,
  FamilyTreeResponse,
  StatisticsResponse,
  CreatePersonRequest,
  UpdatePersonRequest,
  SearchPersonRequest,
  PageResponse,
} from './types';

/**
 * 默认配置
 */
const DEFAULT_PAGE_SIZE = 20;
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1';

/**
 * 通用请求选项
 */
interface RequestOptions {
  headers?: Record<string, string>;
}

/**
 * 人物 API 客户端类
 */
export class PersonApiClient {
  private baseUrl: string;

  constructor(baseUrl: string = API_BASE_URL) {
    this.baseUrl = baseUrl;
  }

  /**
   * 获取人物详情
   * @param id 人物ID
   * @param withRelations 是否加载关联数据（配偶、子女、父亲）
   */
  async getPerson(id: number, withRelations: boolean = false): Promise<PersonDTO> {
    const params = withRelations ? '?with_relations=true' : '';
    const response = await fetch(`${this.baseUrl}/persons/${id}${params}`);
    if (!response.ok) {
      throw new Error(`获取人物详情失败: ${response.status}`);
    }
    const data = await response.json();
    return data.data;
  }

  /**
   * 搜索人物列表
   * @param request 搜索请求参数
   */
  async searchPersons(request: SearchPersonRequest): Promise<PageResponse<PersonDTO>> {
    const params = new URLSearchParams();
    if (request.keyword) params.append('keyword', request.keyword);
    if (request.name) params.append('name', request.name);
    if (request.gender) params.append('gender', request.gender);
    if (request.generation) params.append('generation', String(request.generation));
    if (request.page) params.append('page', String(request.page));
    if (request.pageSize) params.append('page_size', String(request.pageSize));
    if (request.sortBy) params.append('sort_by', request.sortBy);
    if (request.sortDesc) params.append('sort_desc', 'true');

    const response = await fetch(`${this.baseUrl}/persons?${params.toString()}`);
    if (!response.ok) {
      throw new Error(`搜索人物失败: ${response.status}`);
    }
    const data = await response.json();
    return {
      data: data.data,
      total: data.data.total,
      page: data.data.page,
      pageSize: data.data.page_size,
    };
  }

  /**
   * 创建人物
   * @param request 创建人物请求
   * @param token 认证令牌
   */
  async createPerson(request: CreatePersonRequest, token: string): Promise<PersonDTO> {
    const response = await fetch(`${this.baseUrl}/persons`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify(request),
    });
    if (!response.ok) {
      throw new Error(`创建人物失败: ${response.status}`);
    }
    const data = await response.json();
    return data.data;
  }

  /**
   * 更新人物
   * @param id 人物ID
   * @param request 更新人物请求
   * @param token 认证令牌
   */
  async updatePerson(id: number, request: UpdatePersonRequest, token: string): Promise<PersonDTO> {
    const response = await fetch(`${this.baseUrl}/persons/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify(request),
    });
    if (!response.ok) {
      throw new Error(`更新人物失败: ${response.status}`);
    }
    const data = await response.json();
    return data.data;
  }

  /**
   * 删除人物
   * @param id 人物ID
   * @param token 认证令牌
   */
  async deletePerson(id: number, token: string): Promise<void> {
    const response = await fetch(`${this.baseUrl}/persons/${id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });
    if (!response.ok) {
      throw new Error(`删除人物失败: ${response.status}`);
    }
  }

  /**
   * 获取族谱树
   * @param id 人物ID
   * @param upDepth 向上几代（默认5代，最高10代）
   * @param downDepth 向下几代（默认5代，最高10代）
   */
  async getFamilyTree(id: number, upDepth: number = 5, downDepth: number = 5): Promise<FamilyTreeResponse> {
    const clampedUpDepth = Math.min(Math.max(0, upDepth), 10);
    const clampedDownDepth = Math.min(Math.max(0, downDepth), 10);
    const response = await fetch(
      `${this.baseUrl}/persons/${id}/tree?up_depth=${clampedUpDepth}&down_depth=${clampedDownDepth}`
    );
    if (!response.ok) {
      throw new Error(`获取族谱树失败: ${response.status}`);
    }
    const result = await response.json();
    return result.data;
  }

  /**
   * 获取人物统计数据
   */
  async getStatistics(): Promise<StatisticsResponse> {
    const response = await fetch(`${this.baseUrl}/persons/statistics`);
    if (!response.ok) {
      throw new Error(`获取统计数据失败: ${response.status}`);
    }
    const data = await response.json();
    return data.data;
  }
}

// 导出默认实例
export const personApi = new PersonApiClient();

// 导出类型
export type {
  PersonDTO,
  SpouseDTO,
  ChildDTO,
  FamilyTreeResponse,
  StatisticsResponse,
  CreatePersonRequest,
  UpdatePersonRequest,
  SearchPersonRequest,
};