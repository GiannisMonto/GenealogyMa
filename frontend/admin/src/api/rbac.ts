/**
 * 角色与权限 API
 * 基于 @shared/api/types 中的类型定义
 */
import { apiClient } from './client';
import type { RoleDTO, PermissionDTO } from '@shared/api/types';

export interface CreateRoleRequest {
  code: string;
  name: string;
  description?: string;
}

export interface UpdateRoleRequest {
  name?: string;
  description?: string;
}

export interface AssignPermissionsRequest {
  permission_ids: number[];
}

export interface CreatePermissionRequest {
  code: string;
  name: string;
  module: string;
  description?: string;
}

export interface UpdatePermissionRequest {
  name?: string;
  description?: string;
}

/**
 * 获取角色列表
 */
export async function getRoles(): Promise<RoleDTO[]> {
  const response = await apiClient.get<RoleDTO[]>('/roles');
  return response.data;
}

/**
 * 获取角色详情
 */
export async function getRole(id: number): Promise<RoleDTO> {
  const response = await apiClient.get<RoleDTO>(`/roles/${id}`);
  return response.data;
}

/**
 * 创建角色
 */
export async function createRole(data: CreateRoleRequest): Promise<RoleDTO> {
  const response = await apiClient.post<RoleDTO>('/roles', data);
  return response.data;
}

/**
 * 更新角色
 */
export async function updateRole(id: number, data: UpdateRoleRequest): Promise<RoleDTO> {
  const response = await apiClient.put<RoleDTO>(`/roles/${id}`, data);
  return response.data;
}

/**
 * 删除角色
 */
export async function deleteRole(id: number): Promise<void> {
  await apiClient.delete(`/roles/${id}`);
}

/**
 * 获取角色权限列表
 */
export async function getRolePermissions(id: number): Promise<PermissionDTO[]> {
  const response = await apiClient.get<PermissionDTO[]>(`/roles/${id}/permissions`);
  return response.data;
}

/**
 * 分配权限给角色
 */
export async function assignPermissionsToRole(id: number, permissionIds: number[]): Promise<void> {
  await apiClient.post(`/roles/${id}/permissions`, { permission_ids: permissionIds });
}

/**
 * 移除角色权限
 */
export async function removePermissionFromRole(roleId: number, permissionId: number): Promise<void> {
  await apiClient.delete(`/roles/${roleId}/permissions/${permissionId}`);
}

/**
 * 获取权限列表
 */
export async function getPermissions(): Promise<PermissionDTO[]> {
  const response = await apiClient.get<PermissionDTO[]>('/permissions');
  return response.data;
}

/**
 * 创建权限
 */
export async function createPermission(data: CreatePermissionRequest): Promise<PermissionDTO> {
  const response = await apiClient.post<PermissionDTO>('/permissions', data);
  return response.data;
}

/**
 * 更新权限
 */
export async function updatePermission(id: number, data: UpdatePermissionRequest): Promise<PermissionDTO> {
  const response = await apiClient.put<PermissionDTO>(`/permissions/${id}`, data);
  return response.data;
}

/**
 * 删除权限
 */
export async function deletePermission(id: number): Promise<void> {
  await apiClient.delete(`/permissions/${id}`);
}

// 导出所有 API
export const roleApi = {
  getRoles,
  getRole,
  createRole,
  updateRole,
  deleteRole,
  getRolePermissions,
  assignPermissions: assignPermissionsToRole,
  removePermission: removePermissionFromRole,
};

export const permissionApi = {
  getPermissions,
  createPermission,
  updatePermission,
  deletePermission,
};