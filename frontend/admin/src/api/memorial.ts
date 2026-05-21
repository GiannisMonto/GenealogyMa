/**
 * 宗祠管理 API
 */
import { apiClient } from './client';
import type { MemorialHallDTO, TabletDTO } from '@shared/api/types';

export interface CreateHallRequest {
  name: string;
  description?: string;
  province: string;
  city: string;
  district?: string;
  address?: string;
  latitude?: number;
  longitude?: number;
  build_year?: number;
  style?: string;
  image_url?: string;
  total_tablet?: number;
}

export interface UpdateHallRequest {
  name?: string;
  description?: string;
  province?: string;
  city?: string;
  district?: string;
  address?: string;
  latitude?: number;
  longitude?: number;
  build_year?: number;
  style?: string;
  image_url?: string;
  total_tablet?: number;
}

export interface CreateTabletRequest {
  hall_id: number;
  person_id?: number;
  person_name: string;
  generation?: number;
  tablet_type?: 'ancestor' | 'martyr' | 'sage' | 'founder';
  position?: string;
  floor?: number;
  row?: number;
  number: number;
  entronement?: string;
  note?: string;
}

export interface UpdateTabletRequest {
  person_id?: number;
  person_name?: string;
  generation?: number;
  tablet_type?: 'ancestor' | 'martyr' | 'sage' | 'founder';
  position?: string;
  floor?: number;
  row?: number;
  number?: number;
  entronement?: string;
  note?: string;
}

// ===== 宗祠 API =====

export async function getHalls(): Promise<MemorialHallDTO[]> {
  const response = await apiClient.get('/halls');
  return response.data.data;
}

export async function getHall(id: number): Promise<MemorialHallDTO> {
  const response = await apiClient.get(`/halls/${id}`);
  return response.data.data;
}

export async function getHallWithTablets(id: number): Promise<MemorialHallDTO> {
  const response = await apiClient.get(`/halls/${id}/tablets`);
  return response.data.data;
}

export async function createHall(data: CreateHallRequest): Promise<MemorialHallDTO> {
  const response = await apiClient.post('/halls', data);
  return response.data.data;
}

export async function updateHall(id: number, data: UpdateHallRequest): Promise<MemorialHallDTO> {
  const response = await apiClient.put(`/halls/${id}`, data);
  return response.data.data;
}

export async function deleteHall(id: number): Promise<void> {
  await apiClient.delete(`/halls/${id}`);
}

// ===== 牌位 API =====

export async function getTablet(id: number): Promise<TabletDTO> {
  const response = await apiClient.get(`/tablets/${id}`);
  return response.data.data;
}

export async function getTabletByPersonID(personId: number): Promise<TabletDTO> {
  const response = await apiClient.get(`/tablets/person/${personId}`);
  return response.data.data;
}

export async function createTablet(data: CreateTabletRequest): Promise<TabletDTO> {
  const response = await apiClient.post('/tablets', data);
  return response.data.data;
}

export async function updateTablet(id: number, data: UpdateTabletRequest): Promise<TabletDTO> {
  const response = await apiClient.put(`/tablets/${id}`, data);
  return response.data.data;
}

export async function deleteTablet(id: number): Promise<void> {
  await apiClient.delete(`/tablets/${id}`);
}

export const memorialApi = {
  getHalls,
  getHall,
  getHallWithTablets,
  createHall,
  updateHall,
  deleteHall,
  getTablet,
  getTabletByPersonID,
  createTablet,
  updateTablet,
  deleteTablet,
};