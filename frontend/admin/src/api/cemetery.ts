/**
 * 墓园管理 API
 */
import { apiClient } from './client';
import type { CemeteryDTO, GraveDTO } from './types';

export interface CreateCemeteryRequest {
  name: string;
  description?: string;
  province: string;
  city: string;
  district?: string;
  address?: string;
  latitude?: number;
  longitude?: number;
  total_grave?: number;
  image_url?: string;
}

export interface UpdateCemeteryRequest {
  name?: string;
  description?: string;
  province?: string;
  city?: string;
  district?: string;
  address?: string;
  latitude?: number;
  longitude?: number;
  total_grave?: number;
  image_url?: string;
}

export interface CreateGraveRequest {
  cemetery_id: number;
  section: string;
  row: number;
  number: number;
  status?: 'available' | 'occupied' | 'reserved';
  buried_name?: string;
  buried_date?: string;
  buried_year?: number;
  note?: string;
}

export interface UpdateGraveRequest {
  section?: string;
  row?: number;
  number?: number;
  status?: 'available' | 'occupied' | 'reserved';
  buried_name?: string;
  buried_date?: string;
  buried_year?: number;
  note?: string;
}

/**
 * 获取所有墓园
 */
export async function getCemeteries(): Promise<CemeteryDTO[]> {
  const response = await apiClient.get('/cemeteries');
  return response.data.data;
}

/**
 * 获取墓园详情
 */
export async function getCemetery(id: number): Promise<CemeteryDTO> {
  const response = await apiClient.get(`/cemeteries/${id}`);
  return response.data.data;
}

/**
 * 创建墓园
 */
export async function createCemetery(data: CreateCemeteryRequest): Promise<CemeteryDTO> {
  const response = await apiClient.post('/cemeteries', data);
  return response.data.data;
}

/**
 * 更新墓园
 */
export async function updateCemetery(id: number, data: UpdateCemeteryRequest): Promise<CemeteryDTO> {
  const response = await apiClient.put(`/cemeteries/${id}`, data);
  return response.data.data;
}

/**
 * 删除墓园
 */
export async function deleteCemetery(id: number): Promise<void> {
  await apiClient.delete(`/cemeteries/${id}`);
}

/**
 * 获取墓园下的墓位列表
 */
export async function getGraves(cemeteryId: number): Promise<GraveDTO[]> {
  const response = await apiClient.get(`/cemeteries/${cemeteryId}/graves`);
  return response.data.data;
}

/**
 * 获取墓位详情
 */
export async function getGrave(id: number): Promise<GraveDTO> {
  const response = await apiClient.get(`/graves/${id}`);
  return response.data.data;
}

/**
 * 创建墓位
 */
export async function createGrave(data: CreateGraveRequest): Promise<GraveDTO> {
  const response = await apiClient.post('/graves', data);
  return response.data.data;
}

/**
 * 更新墓位
 */
export async function updateGrave(id: number, data: UpdateGraveRequest): Promise<GraveDTO> {
  const response = await apiClient.put(`/graves/${id}`, data);
  return response.data.data;
}

/**
 * 删除墓位
 */
export async function deleteGrave(id: number): Promise<void> {
  await apiClient.delete(`/graves/${id}`);
}

export const cemeteryApi = {
  getCemeteries,
  getCemetery,
  createCemetery,
  updateCemetery,
  deleteCemetery,
  getGraves,
  getGrave,
  createGrave,
  updateGrave,
  deleteGrave,
};