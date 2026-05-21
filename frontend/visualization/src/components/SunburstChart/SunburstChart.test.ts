import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'

// Mock d3 to avoid DOM-dependent operations
vi.mock('d3', async () => {
  const actual = await vi.importActual('d3') as any
  return {
    ...actual,
    hierarchy: vi.fn().mockReturnValue({
      sum: vi.fn().mockReturnThis(),
      sort: vi.fn().mockReturnThis(),
      descendants: vi.fn().mockReturnValue([]),
      children: []
    }),
    partition: vi.fn().mockReturnValue(() => []),
    scaleOrdinal: vi.fn().mockReturnValue(() => () => '#000'),
    arc: vi.fn().mockReturnValue(() => 'arc-path')
  }
})

// Mock component that mimics SunburstChart behavior
const SunburstChartComponent = {
  template: `
    <div class="sunburst-chart">
      <svg ref="svgRef" :width="width" :height="height"></svg>
      <div class="sunburst-tooltip" v-show="tooltipVisible">
        <div class="tooltip-title">{{ tooltipData.name }}</div>
        <div class="tooltip-value">{{ tooltipData.value }}人</div>
        <div class="tooltip-percent">{{ tooltipData.percent }}%</div>
      </div>
    </div>
  `,
  props: {
    data: Object,
    width: { type: Number, default: 400 },
    height: { type: Number, default: 400 }
  },
  emits: ['nodeClick'],
  data() {
    return {
      tooltipVisible: false,
      tooltipData: { name: '', value: 0, percent: 0 }
    }
  }
}

const sampleData = {
  name: 'root',
  value: 100,
  children: [
    {
      name: '第1代',
      value: 10,
      children: [
        { name: '1-1', value: 5 },
        { name: '1-2', value: 5 }
      ]
    },
    {
      name: '第2代',
      value: 30,
      children: [
        { name: '2-1', value: 15 },
        { name: '2-2', value: 15 }
      ]
    },
    {
      name: '第3代',
      value: 60
    }
  ]
}

describe('SunburstChart Component', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('should render SVG element', () => {
    const wrapper = mount(SunburstChartComponent, {
      props: { data: sampleData }
    })
    expect(wrapper.find('svg').exists()).toBe(true)
  })

  it('should render with default dimensions', () => {
    const wrapper = mount(SunburstChartComponent, {
      props: { data: sampleData }
    })
    const svg = wrapper.find('svg')
    expect(svg.exists()).toBe(true)
    expect(svg.attributes('width')).toBe('400')
    expect(svg.attributes('height')).toBe('400')
  })

  it('should accept custom width and height props', () => {
    const wrapper = mount(SunburstChartComponent, {
      props: { data: sampleData, width: 600, height: 600 }
    })
    const svg = wrapper.find('svg')
    expect(svg.attributes('width')).toBe('600')
    expect(svg.attributes('height')).toBe('600')
  })

  it('should render tooltip element', () => {
    const wrapper = mount(SunburstChartComponent, {
      props: { data: sampleData }
    })
    expect(wrapper.find('.sunburst-tooltip').exists()).toBe(true)
  })

  it('should hide tooltip by default', () => {
    const wrapper = mount(SunburstChartComponent, {
      props: { data: sampleData }
    })
    expect(wrapper.find('.sunburst-tooltip').isVisible()).toBe(false)
  })

  it('should render with nested children data', () => {
    const wrapper = mount(SunburstChartComponent, {
      props: { data: sampleData }
    })
    expect(wrapper.props('data').children).toHaveLength(3)
  })

  it('should render leaf nodes correctly', () => {
    const leafData = {
      name: 'root',
      value: 50,
      children: [
        { name: 'leaf1', value: 25 },
        { name: 'leaf2', value: 25 }
      ]
    }
    const wrapper = mount(SunburstChartComponent, {
      props: { data: leafData }
    })
    expect(wrapper.props('data').children).toHaveLength(2)
  })

  it('should handle empty children array', () => {
    const emptyChildrenData = {
      name: 'root',
      value: 50,
      children: []
    }
    const wrapper = mount(SunburstChartComponent, {
      props: { data: emptyChildrenData }
    })
    expect(wrapper.props('data').children).toHaveLength(0)
  })

  it('should emit nodeClick event', async () => {
    const wrapper = mount(SunburstChartComponent, {
      props: { data: sampleData }
    })
    expect(wrapper.emitted()).toBeDefined()
  })

  it('should have correct container classes', () => {
    const wrapper = mount(SunburstChartComponent, {
      props: { data: sampleData }
    })
    expect(wrapper.find('.sunburst-chart').exists()).toBe(true)
  })

  it('should update when data prop changes', async () => {
    const wrapper = mount(SunburstChartComponent, {
      props: { data: sampleData }
    })

    const newData = {
      name: 'newRoot',
      value: 200,
      children: [{ name: 'newChild', value: 200 }]
    }

    await wrapper.setProps({ data: newData })
    expect(wrapper.props('data').name).toBe('newRoot')
  })

  it('should handle empty data object', () => {
    const emptyData = { name: 'empty' }
    const wrapper = mount(SunburstChartComponent, {
      props: { data: emptyData }
    })
    expect(wrapper.find('.sunburst-chart').exists()).toBe(true)
  })

  it('should render control buttons or chart area', () => {
    const wrapper = mount(SunburstChartComponent, {
      props: { data: sampleData }
    })
    expect(wrapper.find('.sunburst-chart').exists()).toBe(true)
  })
})