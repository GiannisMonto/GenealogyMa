/**
 * 文化管理 API
 */
import { apiClient } from './client';
import type { DocumentDTO, StoryDTO, FamilyTeachingsDTO } from '@shared/api/types';

export interface CreateDocumentRequest {
  title: string;
  content: string;
  category: 'classic' | 'genealogy' | 'memorial' | 'history';
  author?: string;
  created_year?: number;
  dynasty?: string;
  source?: string;
  image_urls?: string[];
}

export interface UpdateDocumentRequest {
  title?: string;
  content?: string;
  category?: 'classic' | 'genealogy' | 'memorial' | 'history';
  author?: string;
  created_year?: number;
  dynasty?: string;
  source?: string;
  image_urls?: string[];
}

export interface CreateStoryRequest {
  title: string;
  content: string;
  era?: string;
  category?: string;
  tags?: string[];
  audio_url?: string;
  image_url?: string;
}

export interface UpdateStoryRequest {
  title?: string;
  content?: string;
  era?: string;
  category?: string;
  tags?: string[];
  audio_url?: string;
  image_url?: string;
}

export interface CreateFamilyTeachingsRequest {
  title: string;
  content: string;
  generation?: number;
  origin_text?: string;
  meaning?: string;
}

export interface UpdateFamilyTeachingsRequest {
  title?: string;
  content?: string;
  generation?: number;
  origin_text?: string;
  meaning?: string;
}

// ===== 文献 API =====

export async function getDocuments(): Promise<DocumentDTO[]> {
  const response = await apiClient.get('/culture/documents');
  return response.data.data;
}

export async function getDocument(id: number): Promise<DocumentDTO> {
  const response = await apiClient.get(`/culture/documents/${id}`);
  return response.data.data;
}

export async function createDocument(data: CreateDocumentRequest): Promise<DocumentDTO> {
  const response = await apiClient.post('/culture/documents', data);
  return response.data.data;
}

export async function updateDocument(id: number, data: UpdateDocumentRequest): Promise<DocumentDTO> {
  const response = await apiClient.put(`/culture/documents/${id}`, data);
  return response.data.data;
}

export async function deleteDocument(id: number): Promise<void> {
  await apiClient.delete(`/culture/documents/${id}`);
}

// ===== 故事 API =====

export async function getStories(): Promise<StoryDTO[]> {
  const response = await apiClient.get('/culture/stories');
  return response.data.data;
}

export async function getStory(id: number): Promise<StoryDTO> {
  const response = await apiClient.get(`/culture/stories/${id}`);
  return response.data.data;
}

export async function createStory(data: CreateStoryRequest): Promise<StoryDTO> {
  const response = await apiClient.post('/culture/stories', data);
  return response.data.data;
}

export async function updateStory(id: number, data: UpdateStoryRequest): Promise<StoryDTO> {
  const response = await apiClient.put(`/culture/stories/${id}`, data);
  return response.data.data;
}

export async function deleteStory(id: number): Promise<void> {
  await apiClient.delete(`/culture/stories/${id}`);
}

// ===== 家训 API =====

export async function getFamilyTeachings(): Promise<FamilyTeachingsDTO[]> {
  const response = await apiClient.get('/culture/family-teachings');
  return response.data.data;
}

export async function getFamilyTeaching(id: number): Promise<FamilyTeachingsDTO> {
  const response = await apiClient.get(`/culture/family-teachings/${id}`);
  return response.data.data;
}

export async function createFamilyTeaching(data: CreateFamilyTeachingsRequest): Promise<FamilyTeachingsDTO> {
  const response = await apiClient.post('/culture/family-teachings', data);
  return response.data.data;
}

export async function updateFamilyTeaching(id: number, data: UpdateFamilyTeachingsRequest): Promise<FamilyTeachingsDTO> {
  const response = await apiClient.put(`/culture/family-teachings/${id}`, data);
  return response.data.data;
}

export async function deleteFamilyTeaching(id: number): Promise<void> {
  await apiClient.delete(`/culture/family-teachings/${id}`);
}

export const cultureApi = {
  getDocuments,
  getDocument,
  createDocument,
  updateDocument,
  deleteDocument,
  getStories,
  getStory,
  createStory,
  updateStory,
  deleteStory,
  getFamilyTeachings,
  getFamilyTeaching,
  createFamilyTeaching,
  updateFamilyTeaching,
  deleteFamilyTeaching,
};
