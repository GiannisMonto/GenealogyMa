/**
 * 用户管理 API
 */
import { apiClient } from './client';
import type { UserDTO, CreateUserRequest, UpdateUserRequest, ResetPasswordRequest } from '@shared/api/types';

export interface GetUsersRequest {
  keyword?: string;
  status?: string;
  page?: number;
  page_size?: number;
}

/**
 * 获取用户列表
 */
export async function getUsers(params?: GetUsersRequest): Promise<UserDTO[]> {
  const response = await apiClient.get<UserDTO[]>('/users', { params });
  return response.data;
}

/**
 * 获取用户详情
 */
export async function getUser(id: number): Promise<UserDTO> {
  const response = await apiClient.get<UserDTO>(`/users/${id}`);
  return response.data;
}

/**
 * 创建用户
 */
export async function createUser(data: CreateUserRequest): Promise<UserDTO> {
  const response = await apiClient.post<UserDTO>('/users', data);
  return response.data;
}

/**
 * 更新用户
 */
export async function updateUser(id: number, data: UpdateUserRequest): Promise<UserDTO> {
  const response = await apiClient.put<UserDTO>(`/users/${id}`, data);
  return response.data;
}

/**
 * 删除用户
 */
export async function deleteUser(id: number): Promise<void> {
  await apiClient.delete(`/users/${id}`);
}

/**
 * 重置用户密码
 */
export async function resetUserPassword(id: number, data: ResetPasswordRequest): Promise<void> {
  await apiClient.post(`/users/${id}/reset-password`, data);
}

/**
 * 更新用户状态
 */
export async function updateUserStatus(id: number, status: 'active' | 'inactive' | 'banned'): Promise<UserDTO> {
  const response = await apiClient.put<UserDTO>(`/users/${id}/status`, { status });
  return response.data;
}

// 导出所有 API
export const userApi = {
  getUsers,
  getUser,
  createUser,
  updateUser,
  deleteUser,
  resetPassword: resetUserPassword,
  updateStatus: updateUserStatus,
};