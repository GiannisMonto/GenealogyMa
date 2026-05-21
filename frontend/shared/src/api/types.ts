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

/**
 * 用户 DTO
 */
export interface UserDTO {
  id: number;
  username: string;
  email: string;
  display_name: string;
  avatar_url: string;
  status: 'active' | 'inactive' | 'banned';
  last_login_at: string | null;
  created_at: string;
  updated_at: string;
  roles?: RoleDTO[];
}

/**
 * 创建用户请求
 */
export interface CreateUserRequest {
  username: string;
  email: string;
  password: string;
  display_name?: string;
  role_ids?: number[];
}

/**
 * 更新用户请求
 */
export interface UpdateUserRequest {
  email?: string;
  display_name?: string;
  avatar_url?: string;
  status?: 'active' | 'inactive' | 'banned';
  role_ids?: number[];
}

/**
 * 重置密码请求
 */
export interface ResetPasswordRequest {
  new_password: string;
}

/**
 * 墓园 DTO
 */
export interface CemeteryDTO {
  id: number;
  name: string;
  description: string;
  province: string;
  city: string;
  district: string;
  address: string;
  latitude: number;
  longitude: number;
  total_grave: number;
  used_grave: number;
  image_url: string;
  created_at: string;
  updated_at: string;
}

/**
 * 墓位 DTO
 */
export interface GraveDTO {
  id: number;
  cemetery_id: number;
  person_id: number | null;
  section: string;
  row: number;
  number: number;
  status: 'available' | 'occupied' | 'reserved';
  buried_name: string;
  buried_date: string | null;
  buried_year: number | null;
  note: string;
  created_at: string;
  updated_at: string;
}

/**
 * 文献 DTO
 */
export interface DocumentDTO {
  id: number;
  title: string;
  content: string;
  category: 'classic' | 'genealogy' | 'memorial' | 'history';
  author: string;
  created_year: number | null;
  dynasty: string;
  source: string;
  image_urls: string[];
  view_count: number;
  collect_count: number;
  created_at: string;
  updated_at: string;
}

/**
 * 故事 DTO
 */
export interface StoryDTO {
  id: number;
  title: string;
  content: string;
  era: string;
  category: string;
  tags: string[];
  audio_url: string;
  image_url: string;
  view_count: number;
  created_at: string;
  updated_at: string;
}

/**
 * 家训 DTO
 */
export interface FamilyTeachingsDTO {
  id: number;
  title: string;
  content: string;
  generation: number;
  origin_text: string;
  meaning: string;
  usage_count: number;
  created_at: string;
  updated_at: string;
}