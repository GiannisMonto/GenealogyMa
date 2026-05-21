/**
 * 族谱树可视化组件导出
 */

export { GenealogyTree } from './index';
export type { GenealogyTreeProps, TreeNode, FamilyTreeData, NodeRenderConfig } from './types';
export { personToTreeNode, buildTree, calculateTreeDimensions, treeToHierarchyData } from './utils';