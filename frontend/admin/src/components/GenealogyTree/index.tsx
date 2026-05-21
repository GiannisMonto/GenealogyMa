import { ReactElement, useEffect, useRef, useMemo } from 'react';
import * as d3 from 'd3';
import type { GenealogyTreeProps, TreeNode, NodeRenderConfig } from './types';
import { personToTreeNode, buildTree } from './utils';
import styles from './index.module.css';

const DEFAULT_CONFIG: NodeRenderConfig = {
  width: 120,
  height: 60,
  horizontalSpacing: 30,
  verticalSpacing: 80,
  nodeColor: {
    male: '#e6f7ff',
    female: '#fff1f0',
    selected: '#1890ff',
    alive: '#52c41a',
    deceased: '#8c8c8c',
  },
};

export function GenealogyTree({
  data,
  selectedNodeId,
  onNodeClick,
  config,
}: GenealogyTreeProps): ReactElement {
  const svgRef = useRef<SVGSVGElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  const mergedConfig = useMemo(() => ({ ...DEFAULT_CONFIG, ...config }), [config]);

  // 转换数据为树结构
  const treeData = useMemo(() => {
    if (!data) return null;
    const centerNode = personToTreeNode(data.centerPerson);
    const ancestors = data.ancestors.map(personToTreeNode);
    const descendants = data.descendants.map(personToTreeNode);
    return buildTree({ centerPerson: centerNode, ancestors, descendants });
  }, [data]);

  useEffect(() => {
    if (!svgRef.current || !containerRef.current || !treeData) return;

    const svg = d3.select(svgRef.current);
    svg.selectAll('*').remove();

    const container = containerRef.current;
    const { width, height } = container.getBoundingClientRect();

    // 创建缩放组
    const g = svg.append('g')
      .attr('class', 'zoom-group')
      .attr('transform', `translate(${width / 2}, 60)`);

    // 创建树布局
    const hierarchy = d3.hierarchy<TreeNode>(treeData.rootNode, (d) => d.children);
    const treeLayout = d3.tree<TreeNode>()
      .nodeSize([mergedConfig.width + mergedConfig.horizontalSpacing, mergedConfig.height + mergedConfig.verticalSpacing]);

    const treeRoot = treeLayout(hierarchy);

    // 绘制连接线
    const linksGroup = g.append('g').attr('class', 'links');
    linksGroup.selectAll('path')
      .data(treeRoot.links())
      .enter()
      .append('path')
      .attr('class', styles.link)
      .attr('d', (d) => {
        const sourceX = d.source.x;
        const sourceY = d.source.y + mergedConfig.height / 2;
        const targetX = d.target.x;
        const targetY = d.target.y - mergedConfig.height / 2;
        const midY = (sourceY + targetY) / 2;
        return `M ${sourceX} ${sourceY} C ${sourceX} ${midY}, ${targetX} ${midY}, ${targetX} ${targetY}`;
      })
      .attr('fill', 'none')
      .attr('stroke', '#d9d9d9')
      .attr('stroke-width', 1.5);

    // 绘制节点
    const nodesGroup = g.append('g').attr('class', 'nodes');
    const nodes = nodesGroup.selectAll('g')
      .data(treeRoot.descendants())
      .enter()
      .append('g')
      .attr('class', styles.node)
      .attr('transform', (d) => `translate(${d.x - mergedConfig.width / 2}, ${d.y - mergedConfig.height / 2})`)
      .style('cursor', 'pointer')
      .on('click', (_, d) => {
        onNodeClick?.(d.data);
      });

    // 节点背景
    nodes.append('rect')
      .attr('width', mergedConfig.width)
      .attr('height', mergedConfig.height)
      .attr('rx', 4)
      .attr('fill', (d) => {
        const isSelected = d.data.id === selectedNodeId;
        if (isSelected) return mergedConfig.nodeColor.selected;
        return d.data.gender === '男' ? mergedConfig.nodeColor.male : mergedConfig.nodeColor.female;
      })
      .attr('stroke', (d) => {
        const isSelected = d.data.id === selectedNodeId;
        return isSelected ? mergedConfig.nodeColor.selected : '#d9d9d9';
      })
      .attr('stroke-width', (d) => (d.data.id === selectedNodeId ? 2 : 1));

    // 姓名文字
    nodes.append('text')
      .attr('x', mergedConfig.width / 2)
      .attr('y', mergedConfig.height / 2 - 6)
      .attr('text-anchor', 'middle')
      .attr('dominant-baseline', 'middle')
      .attr('font-size', 14)
      .attr('font-weight', 500)
      .attr('fill', '#262626')
      .text((d) => d.data.name);

    // 字号文字
    nodes.append('text')
      .attr('x', mergedConfig.width / 2)
      .attr('y', mergedConfig.height / 2 + 10)
      .attr('text-anchor', 'middle')
      .attr('dominant-baseline', 'middle')
      .attr('font-size', 11)
      .attr('fill', '#8c8c8c')
      .text((d) => d.data.styleName || '');

    // 性别图标
    nodes.append('text')
      .attr('x', 8)
      .attr('y', 12)
      .attr('font-size', 12)
      .attr('fill', '#8c8c8c')
      .text((d) => (d.data.gender === '男' ? '♂' : '♀'));

    // 存活状态图标
    nodes.append('circle')
      .attr('cx', mergedConfig.width - 12)
      .attr('cy', 12)
      .attr('r', 5)
      .attr('fill', (d) => (d.data.isAlive ? mergedConfig.nodeColor.alive : mergedConfig.nodeColor.deceased));

    // 设置缩放
    const zoom = d3.zoom<SVGSVGElement, unknown>()
      .scaleExtent([0.3, 3])
      .on('zoom', (event) => {
        g.attr('transform', event.transform);
      });

    svg.call(zoom);

    // 初始居中
    const initialTransform = d3.zoomIdentity.translate(width / 2, 80).scale(0.8);
    svg.call(zoom.transform, initialTransform);

  }, [treeData, selectedNodeId, onNodeClick, mergedConfig]);

  return (
    <div ref={containerRef} className={styles.container}>
      <svg ref={svgRef} className={styles.svg} />
    </div>
  );
}