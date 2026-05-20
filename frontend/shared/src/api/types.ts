/**
 * API 类型定义
 * 人物相关 API 的 TypeScript 类型定义
 */

/**
 * 人物 DTO
 */
export interface PersonDTO {
  id: number;
  legacy_id: string;
  name: string;
  style_name: string;
  gender: '男' | '女' | '';
  generation: number;
  birth_order: string;
  father_id: number | null;
  lineage_path: string;
  detail_text: string;
  birth_time_text: string;
  death_time_text: string;
  birth_place: string;
  burial_place: string;
  son_count: number;
  daughter_count: number;
  adopted_heir_count: number;
  total_children_count: number;
  age: number;
  is_alive: boolean;
  full_name: string;
  created_at: string;
  updated_at: string;

  // 关联数据
  spouses?: SpouseDTO[];
  children?: ChildDTO[];
  father?: PersonDTO;
}

/**
 * 配偶 DTO
 */
export interface SpouseDTO {
  id: number;
  spouse_type: '配' | '继' | '妣' | '三';
  name: string;
  birth_time_text: string;
  death_time_text: string;
  birth_place: string;
  burial_place: string;
}

/**
 * 子女 DTO
 */
export interface ChildDTO extends PersonDTO {
  relation_type: 'biological' | 'adoptive' | 'step';
  birth_order_num: number;
  is_primary: boolean;
}

/**
 * 创建人物请求
 */
export interface CreatePersonRequest {
  name: string;
  style_name?: string;
  gender: '男' | '女';
  generation?: number;
  birth_order?: string;
  father_id?: number;
  detail_text?: string;
  birth_time_text?: string;
  death_time_text?: string;
  birth_place?: string;
  burial_place?: string;
}

/**
 * 更新人物请求
 */
export interface UpdatePersonRequest {
  name?: string;
  style_name?: string;
  gender?: '男' | '女';
  generation?: number;
  birth_order?: string;
  father_id?: number;
  detail_text?: string;
  birth_time_text?: string;
  death_time_text?: string;
  birth_place?: string;
  burial_place?: string;
}

/**
 * 搜索人物请求
 */
export interface SearchPersonRequest {
  keyword?: string;
  name?: string;
  gender?: '男' | '女';
  generation?: number;
  page?: number;
  page_size?: number;
  sort_by?: string;
  sort_desc?: boolean;
}

/**
 * 族谱树响应
 */
export interface FamilyTreeResponse {
  center_person: PersonDTO;
  ancestors: PersonDTO[];
  descendants: PersonDTO[];
  generations: number;
  total_nodes: number;
}

/**
 * 统计响应
 */
export interface StatisticsResponse {
  total_persons: number;
  by_generation: Record<number, number>;
}

/**
 * 分页响应
 */
export interface PageResponse<T> {
  data: T[];
  total: number;
  page: number;
  page_size: number;
}

/**
 * API 统一响应格式
 */
export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

/**
 * 角色 DTO
 */
export interface RoleDTO {
  id: number;
  code: string;
  name: string;
  description?: string;
  is_system: boolean;
  created_at: string;
  updated_at: string;
}

/**
 * 权限 DTO
 */
export interface PermissionDTO {
  id: number;
  code: string;
  name: string;
  module: string;
  description?: string;
  created_at: string;
}