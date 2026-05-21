/**
 * 族谱树组件类型定义
 */

/**
 * 族谱树节点
 */
export interface TreeNode {
  id: number;
  name: string;
  styleName?: string;
  gender: '男' | '女' | '';
  generation: number;
  isAlive: boolean;
  fatherId?: number | null;
  children?: TreeNode[];
}

/**
 * 族谱树数据
 */
export interface FamilyTreeData {
  centerPerson: TreeNode;
  ancestors: TreeNode[];
  descendants: TreeNode[];
}

/**
 * 树节点渲染配置
 */
export interface NodeRenderConfig {
  width: number;
  height: number;
  horizontalSpacing: number;
  verticalSpacing: number;
  nodeColor: {
    male: string;
    female: string;
    selected: string;
    alive: string;
    deceased: string;
  };
}

/**
 * 族谱树组件 Props
 */
export interface GenealogyTreeProps {
  data: FamilyTreeData | null;
  selectedNodeId?: number | null;
  onNodeClick?: (node: TreeNode) => void;
  config?: Partial<NodeRenderConfig>;
}