<template>
  <div class="genealogy-tree" ref="container">
    <svg ref="svgRef"></svg>
    <div class="tree-controls">
      <button @click="zoomIn" title="放大">+</button>
      <button @click="zoomOut" title="缩小">-</button>
      <button @click="resetZoom" title="重置">⟲</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, nextTick } from 'vue'
import * as d3 from 'd3'
import type { Person } from '@/types'

interface TreeNode {
  id: number
  name: string
  generation: number
  gender: string
  styleName?: string
  birthYear?: number
  deathYear?: number
  children?: TreeNode[]
}

const props = defineProps<{
  data: TreeNode | null
}>()

const emit = defineEmits<{
  (e: 'node-click', person: Person): void
}>()

const container = ref<HTMLElement | null>(null)
const svgRef = ref<SVGSVGElement | null>(null)
const zoom = ref<d3.ZoomBehavior<SVGSVGElement, unknown> | null>(null)

const NODE_WIDTH = 120
const NODE_HEIGHT = 60
const HORIZONTAL_SPACING = 40
const VERTICAL_SPACING = 100

function buildTree(data: TreeNode): d3.HierarchyNode<TreeNode> {
  return d3.hierarchy(data)
}

function renderTree() {
  if (!svgRef.value || !container.value || !props.data) return

  const svg = d3.select(svgRef.value)
  svg.selectAll('*').remove()

  const width = container.value.clientWidth
  const height = container.value.clientHeight

  svg.attr('width', width).attr('height', height)

  const g = svg.append('g')

  const root = buildTree(props.data)
  const treeLayout = d3.tree<TreeNode>()
    .nodeSize([NODE_WIDTH + HORIZONTAL_SPACING, NODE_HEIGHT + VERTICAL_SPACING])

  treeLayout(root)

  g.selectAll('.link')
    .data(root.links())
    .enter()
    .append('path')
    .attr('class', 'link')
    .attr('d', d3.linkVertical<d3.HierarchyLink<TreeNode>, d3.HierarchyPointNode<TreeNode>>()
      .x(d => d.x + NODE_WIDTH / 2)
      .y(d => d.y + NODE_HEIGHT / 2)
    )
    .attr('fill', 'none')
    .attr('stroke', '#DEB887')
    .attr('stroke-width', 2)

  const nodes = g.selectAll('.node')
    .data(root.descendants())
    .enter()
    .append('g')
    .attr('class', 'node')
    .attr('transform', d => `translate(${d.x}, ${d.y})`)
    .style('cursor', 'pointer')
    .on('click', (_event, d) => {
      emit('node-click', {
        member_id: d.data.id,
        name: d.data.name,
        gender: d.data.gender,
        generation: d.data.generation,
        lineage_path: '',
        style_name: d.data.styleName,
        birth_year: d.data.birthYear,
        death_year: d.data.deathYear
      })
    })

  nodes.append('rect')
    .attr('width', NODE_WIDTH)
    .attr('height', NODE_HEIGHT)
    .attr('x', 0)
    .attr('y', 0)
    .attr('rx', 8)
    .attr('fill', d => d.data.gender === 'M' ? '#E8F4FD' : '#FDE8F4')
    .attr('stroke', d => d.data.gender === 'M' ? '#1890FF' : '#EB2F96')
    .attr('stroke-width', 2)

  nodes.append('text')
    .attr('x', NODE_WIDTH / 2)
    .attr('y', NODE_HEIGHT / 2 - 8)
    .attr('text-anchor', 'middle')
    .attr('dominant-baseline', 'middle')
    .attr('font-size', '14px')
    .attr('font-weight', '600')
    .attr('fill', '#333')
    .text(d => d.data.name)

  nodes.append('text')
    .attr('x', NODE_WIDTH / 2)
    .attr('y', NODE_HEIGHT / 2 + 12)
    .attr('text-anchor', 'middle')
    .attr('dominant-baseline', 'middle')
    .attr('font-size', '11px')
    .attr('fill', '#666')
    .text(d => `第${d.data.generation}代`)

  zoom.value = d3.zoom<SVGSVGElement, unknown>()
    .scaleExtent([0.1, 4])
    .on('zoom', (event) => {
      g.attr('transform', event.transform)
    })

  svg.call(zoom.value)

  const initialX = width / 2
  const initialY = 60
  svg.call(zoom.value.transform, d3.zoomIdentity.translate(initialX, initialY))
}

function zoomIn() {
  if (svgRef.value && zoom.value) {
    d3.select(svgRef.value).transition().call(zoom.value.scaleBy, 1.3)
  }
}

function zoomOut() {
  if (svgRef.value && zoom.value) {
    d3.select(svgRef.value).transition().call(zoom.value.scaleBy, 0.7)
  }
}

function resetZoom() {
  if (svgRef.value && zoom.value && container.value) {
    const width = container.value.clientWidth
    d3.select(svgRef.value).transition().call(
      zoom.value.transform,
      d3.zoomIdentity.translate(width / 2, 60)
    )
  }
}

onMounted(() => {
  nextTick(() => {
    renderTree()
  })
})

watch(() => props.data, () => {
  nextTick(() => {
    renderTree()
  })
})
</script>

<style scoped>
.genealogy-tree {
  width: 100%;
  height: 100%;
  position: relative;
  overflow: hidden;
}

.genealogy-tree svg {
  width: 100%;
  height: 100%;
}

.tree-controls {
  position: absolute;
  bottom: 20px;
  right: 20px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.tree-controls button {
  width: 40px;
  height: 40px;
  border: none;
  border-radius: 50%;
  background: var(--color-primary);
  color: white;
  font-size: 1.2rem;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  transition: transform 0.2s;
}

.tree-controls button:hover {
  transform: scale(1.1);
}
</style>
