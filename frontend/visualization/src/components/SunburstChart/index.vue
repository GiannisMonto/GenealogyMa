<template>
  <div class="sunburst-chart" ref="containerRef">
    <svg ref="svgRef" :width="width" :height="height"></svg>
    <div class="sunburst-tooltip" ref="tooltipRef" v-show="tooltipVisible">
      <div class="tooltip-title">{{ tooltipData.name }}</div>
      <div class="tooltip-value">{{ tooltipData.value }}人</div>
      <div class="tooltip-percent">{{ tooltipData.percent }}%</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, onUnmounted } from 'vue'
import * as d3 from 'd3'

interface SunburstNode {
  name: string
  value?: number
  children?: SunburstNode[]
}

const props = withDefaults(
  defineProps<{
    data: SunburstNode
    width?: number
    height?: number
  }>(),
  {
    width: 400,
    height: 400
  }
)

const emit = defineEmits<{
  (e: 'nodeClick', node: SunburstNode): void
}>()

const svgRef = ref<SVGSVGElement | null>(null)
const containerRef = ref<HTMLDivElement | null>(null)
const tooltipRef = ref<HTMLDivElement | null>(null)
const tooltipVisible = ref(false)
const tooltipData = ref({ name: '', value: 0, percent: 0 })

let svg: d3.Selection<SVGSVGElement, unknown, null, undefined>
let g: d3.Selection<SVGGElement, unknown, null, undefined>

function createSunburst() {
  if (!svgRef.value || !props.data) return

  // Clear previous content
  d3.select(svgRef.value).selectAll('*').remove()

  svg = d3.select(svgRef.value)
  g = svg.append('g').attr('transform', `translate(${props.width / 2},${props.height / 2})`)

  // Create hierarchy
  const root = d3
    .hierarchy<SunburstNode>(props.data)
    .sum((d) => d.value || 0)
    .sort((a, b) => (b.value || 0) - (a.value || 0))

  // Create partition layout
  const radius = Math.min(props.width, props.height) / 2
  const partition = d3.partition<SunburstNode>().size([2 * Math.PI, radius])

  partition(root)

  // Color scale
  const color = d3
    .scaleOrdinal<string>()
    .domain(root.children ? root.children.map((d) => d.data.name) : [])
    .range(['#1976d2', '#388e3c', '#d32f2f', '#7b1fa2', '#f57c00', '#00838f'])

  // Arc generator
  const arc = d3
    .arc<d3.HierarchyRectangularNode<SunburstNode>>()
    .startAngle((d) => d.x0)
    .endAngle((d) => d.x1)
    .padAngle((d) => Math.min((d.x1 - d.x0) / 2, 0.005))
    .padRadius(radius / 2)
    .innerRadius((d) => d.y0)
    .outerRadius((d) => d.y1 - 1)

  // Total value for percentage calculation
  const totalValue = root.value || 1

  // Create arcs
  g.selectAll('path')
    .data(root.descendants().filter((d) => d.depth > 0))
    .join('path')
    .attr('fill', (d) => {
      let node = d
      while (node.depth > 1) node = node.parent!
      return color(node.data.name)
    })
    .attr('fill-opacity', (d) => 1 - (d.depth - 1) * 0.15)
    .attr('d', arc)
    .style('cursor', 'pointer')
    .on('mouseover', function (event, d) {
      d3.select(this).attr('fill-opacity', 1)
      tooltipData.value = {
        name: d.data.name,
        value: d.value || 0,
        percent: (((d.value || 0) / totalValue) * 100).toFixed(1)
      }
      tooltipVisible.value = true
    })
    .on('mousemove', function (event) {
      if (tooltipRef.value && containerRef.value) {
        const rect = containerRef.value.getBoundingClientRect()
        tooltipRef.value.style.left = event.clientX - rect.left + 10 + 'px'
        tooltipRef.value.style.top = event.clientY - rect.top - 10 + 'px'
      }
    })
    .on('mouseout', function (_, d) {
      d3.select(this).attr('fill-opacity', 1 - (d.depth - 1) * 0.15)
      tooltipVisible.value = false
    })
    .on('click', function (_, d) {
      emit('nodeClick', d.data)
    })

  // Add labels for top-level nodes
  g.selectAll('text')
    .data(root.descendants().filter((d) => d.depth === 1 && (d.y1 - d.y0) > 20))
    .join('text')
    .attr('transform', (d) => {
      const angle = ((d.x0 + d.x1) / 2) * (180 / Math.PI) - 90
      const radius = (d.y0 + d.y1) / 2
      return `rotate(${angle}) translate(${radius},0) rotate(${angle > 90 ? 180 : 0})`
    })
    .attr('dy', '0.35em')
    .attr('text-anchor', 'middle')
    .attr('fill', '#fff')
    .attr('font-size', '12px')
    .text((d) => d.data.name)
}

onMounted(() => {
  createSunburst()
})

watch(
  () => [props.data, props.width, props.height],
  () => createSunburst(),
  { deep: true }
)

onUnmounted(() => {
  if (svgRef.value) {
    d3.select(svgRef.value).selectAll('*').remove()
  }
})
</script>

<style scoped>
.sunburst-chart {
  position: relative;
  width: 100%;
  height: 100%;
}

.sunburst-chart svg {
  display: block;
}

.sunburst-tooltip {
  position: absolute;
  background: rgba(0, 0, 0, 0.8);
  color: white;
  padding: 8px 12px;
  border-radius: 4px;
  pointer-events: none;
  font-size: 12px;
  z-index: 100;
}

.tooltip-title {
  font-weight: 600;
  margin-bottom: 4px;
}

.tooltip-value {
  color: #ccc;
}

.tooltip-percent {
  color: #1976d2;
  font-weight: 500;
}
</style>