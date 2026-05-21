/**
 * 族谱树数据处理工具
 */

import type { FamilyTreeData, TreeNode, PersonDTO } from './types';

/**
 * 从 PersonDTO 转换为 TreeNode
 */
export function personToTreeNode(person: PersonDTO): TreeNode {
  return {
    id: person.id,
    name: person.name,
    styleName: person.style_name,
    gender: person.gender,
    generation: person.generation,
    isAlive: person.is_alive,
    fatherId: person.father_id,
  };
}

/**
 * 将族谱数据构建为树形结构
 * @param data 族谱数据
 * @returns 根节点和所有节点的映射
 */
export function buildTree(data: FamilyTreeData): { rootNode: TreeNode; nodeMap: Map<number, TreeNode> } {
  const nodeMap = new Map<number, TreeNode>();
  const allNodes: TreeNode[] = [data.centerPerson, ...data.ancestors, ...data.descendants];

  // 将所有节点转为 TreeNode 并建立映射
  allNodes.forEach((node) => {
    nodeMap.set(node.id, { ...node, children: [] });
  });

  // 构建父子关系
  const rootNode = nodeMap.get(data.centerPerson.id)!;
  nodeMap.forEach((node) => {
    if (node.fatherId && nodeMap.has(node.fatherId)) {
      const parent = nodeMap.get(node.fatherId)!;
      if (!parent.children) {
        parent.children = [];
      }
      // 按出生顺序排序
      const siblings = parent.children;
      const insertIndex = siblings.findIndex(
        (sibling) => (sibling.generation || 0) > (node.generation || 0) ||
          ((sibling.generation || 0) === (node.generation || 0) && (sibling as any)._birthOrderNum > (node as any)._birthOrderNum)
      );
      if (insertIndex === -1) {
        siblings.push(node);
      } else {
        siblings.splice(insertIndex, 0, node);
      }
    }
  });

  return { rootNode, nodeMap };
}

/**
 * 计算树的深度和宽度
 */
export function calculateTreeDimensions(node: TreeNode, config: { width: number; height: number; verticalSpacing: number }): { depth: number; maxWidth: number } {
  if (!node.children || node.children.length === 0) {
    return { depth: 1, maxWidth: config.width };
  }

  let maxWidth = 0;
  let maxChildDepth = 0;

  node.children.forEach((child) => {
    const childDimensions = calculateTreeDimensions(child, config);
    maxWidth = Math.max(maxWidth, childDimensions.maxWidth);
    maxChildDepth = Math.max(maxChildDepth, childDimensions.depth);
  });

  // 当前层级的宽度 = 所有子节点宽度 + 间距
  const childrenTotalWidth = node.children.reduce((sum, child) => {
    const childDim = calculateTreeDimensions(child, config);
    return sum + childDim.maxWidth;
  }, 0) + (node.children.length - 1) * 20; // 节点间距

  return {
    depth: maxChildDepth + 1,
    maxWidth: Math.max(config.width, childrenTotalWidth),
  };
}

/**
 * 将树结构转换为 D3 层级数据格式
 */
export function treeToHierarchyData(node: TreeNode): any {
  return {
    name: node.name,
    data: node,
    children: node.children?.map(child => treeToHierarchyData(child)),
  };
}