import { describe, it, expect } from 'vitest';
import { personToTreeNode, buildTree, calculateTreeDimensions, treeToHierarchyData } from './utils';
import type { FamilyTreeData, TreeNode } from './types';

describe('GenealogyTree utils', () => {
  const mockPerson = {
    id: 1,
    name: '张三',
    style_name: '伯高',
    gender: '男' as const,
    generation: 10,
    is_alive: true,
    father_id: null,
  };

  const mockTreeNode: TreeNode = {
    id: 1,
    name: '张三',
    styleName: '伯高',
    gender: '男',
    generation: 10,
    isAlive: true,
  };

  describe('personToTreeNode', () => {
    it('converts PersonDTO to TreeNode correctly', () => {
      const result = personToTreeNode(mockPerson);
      expect(result.id).toBe(1);
      expect(result.name).toBe('张三');
      expect(result.styleName).toBe('伯高');
      expect(result.gender).toBe('男');
      expect(result.generation).toBe(10);
      expect(result.isAlive).toBe(true);
    });

    it('handles optional fatherId', () => {
      const personWithFather = { ...mockPerson, father_id: 2 };
      const result = personToTreeNode(personWithFather);
      expect(result.fatherId).toBe(2);
    });
  });

  describe('buildTree', () => {
    it('builds tree from family data', () => {
      const data: FamilyTreeData = {
        centerPerson: mockTreeNode,
        ancestors: [
          { id: 2, name: '张二', styleName: '仲明', gender: '男' as const, generation: 9, isAlive: false, fatherId: null },
        ],
        descendants: [
          { id: 3, name: '张小', styleName: '季实', gender: '男' as const, generation: 11, isAlive: true, fatherId: 1 },
        ],
      };

      const { rootNode, nodeMap } = buildTree(data);
      expect(rootNode.id).toBe(1);
      expect(nodeMap.size).toBe(3);
    });

    it('handles empty ancestors and descendants', () => {
      const data: FamilyTreeData = {
        centerPerson: mockTreeNode,
        ancestors: [],
        descendants: [],
      };

      const { rootNode, nodeMap } = buildTree(data);
      expect(rootNode.id).toBe(1);
      expect(nodeMap.size).toBe(1);
    });

    it('builds correct parent-child relationships', () => {
      const data: FamilyTreeData = {
        centerPerson: mockTreeNode,
        ancestors: [],
        descendants: [
          { id: 3, name: '张小', styleName: '季实', gender: '男' as const, generation: 11, isAlive: true, fatherId: 1 },
        ],
      };

      const { nodeMap } = buildTree(data);
      const child = nodeMap.get(3);
      const parent = nodeMap.get(1);
      expect(parent?.children).toBeDefined();
      expect(parent?.children?.some(c => c.id === 3)).toBe(true);
    });
  });

  describe('calculateTreeDimensions', () => {
    it('calculates dimensions for leaf node', () => {
      const leafNode: TreeNode = {
        id: 1,
        name: '叶节点',
        gender: '男',
        generation: 10,
        isAlive: true,
      };

      const result = calculateTreeDimensions(leafNode, { width: 120, height: 60, verticalSpacing: 80 });
      expect(result.depth).toBe(1);
    });

    it('calculates dimensions for tree with children', () => {
      const parentNode: TreeNode = {
        id: 1,
        name: '父亲',
        gender: '男',
        generation: 10,
        isAlive: true,
        children: [
          { id: 2, name: '孩子', gender: '男', generation: 11, isAlive: true },
        ],
      };

      const result = calculateTreeDimensions(parentNode, { width: 120, height: 60, verticalSpacing: 80 });
      expect(result.depth).toBe(2);
    });
  });

  describe('treeToHierarchyData', () => {
    it('converts tree to D3 hierarchy format', () => {
      const treeNode: TreeNode = {
        id: 1,
        name: '根节点',
        gender: '男',
        generation: 10,
        isAlive: true,
        children: [
          { id: 2, name: '子节点', gender: '女', generation: 11, isAlive: true },
        ],
      };

      const result = treeToHierarchyData(treeNode);
      expect(result.name).toBe('根节点');
      expect(result.data).toEqual(treeNode);
      expect(result.children).toHaveLength(1);
      expect(result.children![0].name).toBe('子节点');
    });

    it('handles node without children', () => {
      const treeNode: TreeNode = {
        id: 1,
        name: '叶节点',
        gender: '男',
        generation: 10,
        isAlive: true,
      };

      const result = treeToHierarchyData(treeNode);
      expect(result.name).toBe('叶节点');
      expect(result.children).toBeUndefined();
    });
  });
});