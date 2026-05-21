/**
 * 人物相关类型定义
 */

export interface PersonDTO {
  id: number;
  name: string;
  styleName?: string;
  gender: '男' | '女';
  birthDate?: string;
  deathDate?: string;
  generation?: number;
  isAlive: boolean;
  fatherId?: number;
  motherId?: number;
  birthPlace?: string;
  burialPlace?: string;
  biography?: string;
  portrait?: string;
}

export interface SpouseDTO {
  id: number;
  personId: number;
  spouseId: number;
  spouseName: string;
  spouseGender: '男' | '女';
  marriageDate?: string;
  divorceDate?: string;
  isActive: boolean;
}

export interface ChildDTO {
  id: number;
  personId: number;
  childId: number;
  childName: string;
  childGender: '男' | '女';
  birthDate?: string;
  isAlive: boolean;
}

export interface SearchPersonRequest {
  keyword?: string;
  name?: string;
  gender?: '男' | '女';
  generation?: number;
  birthPlace?: string;
  page?: number;
  pageSize?: number;
}

export interface PageResponse<T> {
  data: T[];
  total: number;
  page: number;
  pageSize: number;
}

export interface FamilyTreeNode {
  id: number;
  name: string;
  styleName?: string;
  gender: '男' | '女';
  generation?: number;
  isAlive: boolean;
  spouseId?: number;
  spouseName?: string;
  children?: FamilyTreeNode[];
}

export interface FamilyTreeResponse {
  centerPerson: FamilyTreeNode;
  ancestors: FamilyTreeNode[];
  descendants: FamilyTreeNode[];
}