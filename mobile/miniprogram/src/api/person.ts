/**
 * 人物 API
 */
import { request } from './request';
import type {
  PersonDTO,
  FamilyTreeResponse,
  SearchPersonRequest,
  PageResponse,
} from '../types/person';

export interface CreatePersonRequest {
  name: string;
  styleName?: string;
  gender: '男' | '女';
  birthDate?: string;
  deathDate?: string;
  generation?: number;
  fatherId?: number;
  motherId?: number;
  birthPlace?: string;
  burialPlace?: string;
  biography?: string;
}

export interface UpdatePersonRequest extends Partial<CreatePersonRequest> {}

/**
 * 获取人物详情
 */
export async function getPerson(id: number, withRelations = false): Promise<PersonDTO> {
  const params = withRelations ? '?with_relations=true' : '';
  return request<PersonDTO>(`/persons/${id}${params}`);
}

/**
 * 搜索人物
 */
export async function searchPersons(
  params: SearchPersonRequest
): Promise<PageResponse<PersonDTO>> {
  const query = new URLSearchParams();
  if (params.keyword) query.append('keyword', params.keyword);
  if (params.name) query.append('name', params.name);
  if (params.gender) query.append('gender', params.gender);
  if (params.generation) query.append('generation', String(params.generation));
  if (params.page) query.append('page', String(params.page));
  if (params.pageSize) query.append('page_size', String(params.pageSize));

  return request<PageResponse<PersonDTO>>(`/persons?${query.toString()}`);
}

/**
 * 获取族谱树
 */
export async function getFamilyTree(
  id: number,
  upDepth = 5,
  downDepth = 5
): Promise<FamilyTreeResponse> {
  const depth = Math.min(upDepth, 10) + Math.min(downDepth, 10);
  return request<FamilyTreeResponse>(
    `/persons/${id}/tree?up_depth=${Math.min(upDepth, 10)}&down_depth=${Math.min(downDepth, 10)}`
  );
}

/**
 * 获取统计数据
 */
export async function getStatistics(): Promise<{
  totalMembers: number;
  totalGenerations: number;
  malesCount: number;
  femalesCount: number;
  aliveCount: number;
  deceasedCount: number;
}> {
  return request('/persons/statistics');
}