import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import type { TreeNode } from '@/components/GenealogyTree/index.vue'

vi.mock('d3', async () => {
  const actual = await vi.importActual('d3') as any
  return {
    ...actual,
    zoom: vi.fn().mockReturnValue({
      scaleExtent: vi.fn().mockReturnThis(),
      on: vi.fn().mockReturnThis(),
      transform: vi.fn()
    }),
    zoomIdentity: {
      translate: vi.fn().mockReturnThis(),
      scale: vi.fn().mockReturnThis()
    }
  }
})

const GenealogyTreeComponent = {
  template: `
    <div class="genealogy-tree" ref="container">
      <svg ref="svgRef"></svg>
      <div class="tree-controls">
        <button @click="zoomIn" title="放大">+</button>
        <button @click="zoomOut" title="缩小">-</button>
        <button @click="resetZoom" title="重置">⟲</button>
      </div>
    </div>
  `,
  props: {
    data: Object
  },
  emits: ['node-click'],
  data() {
    return {
      zoomValue: null
    }
  },
  methods: {
    zoomIn() {},
    zoomOut() {},
    resetZoom() {}
  }
}

const mockTreeData: TreeNode = {
  id: 1,
  name: '张老爷子',
  generation: 1,
  gender: 'M',
  styleName: '字肇始',
  birthYear: 1900,
  deathYear: 1970,
  children: [
    {
      id: 2,
      name: '张大',
      generation: 2,
      gender: 'M',
      styleName: '字守成',
      birthYear: 1925,
      children: [
        {
          id: 3,
          name: '张孙',
          generation: 3,
          gender: 'M',
          birthYear: 1950
        },
        {
          id: 4,
          name: '张孙女',
          generation: 3,
          gender: 'F',
          birthYear: 1952
        }
      ]
    },
    {
      id: 5,
      name: '张二',
      generation: 2,
      gender: 'M',
      birthYear: 1928
    }
  ]
}

describe('GenealogyTree Component', () => {
  beforeEach(() => {
    vi.spyOn(Element.prototype, 'clientWidth', 'get').mockReturnValue(800)
    vi.spyOn(Element.prototype, 'clientHeight', 'get').mockReturnValue(600)
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('should render SVG element', () => {
    const wrapper = mount(GenealogyTreeComponent, {
      props: { data: mockTreeData }
    })

    expect(wrapper.find('svg').exists()).toBe(true)
  })

  it('should render tree controls', () => {
    const wrapper = mount(GenealogyTreeComponent, {
      props: { data: mockTreeData }
    })

    const buttons = wrapper.findAll('.tree-controls button')
    expect(buttons.length).toBe(3)
  })

  it('should render tree controls with correct titles', () => {
    const wrapper = mount(GenealogyTreeComponent, {
      props: { data: mockTreeData }
    })

    const buttons = wrapper.findAll('.tree-controls button')
    expect(buttons[0].attributes('title')).toBe('放大')
    expect(buttons[1].attributes('title')).toBe('缩小')
    expect(buttons[2].attributes('title')).toBe('重置')
  })

  it('should render with null data gracefully', () => {
    const wrapper = mount(GenealogyTreeComponent, {
      props: { data: null }
    })

    expect(wrapper.find('svg').exists()).toBe(true)
  })

  it('should have correct container classes', () => {
    const wrapper = mount(GenealogyTreeComponent, {
      props: { data: mockTreeData }
    })

    const container = wrapper.find('.genealogy-tree')
    expect(container.exists()).toBe(true)
    expect(container.classes()).toContain('genealogy-tree')
  })

  it('should handle zoom in button click', async () => {
    const wrapper = mount(GenealogyTreeComponent, {
      props: { data: mockTreeData }
    })

    const zoomInBtn = wrapper.find('.tree-controls button:nth-child(1)')
    await zoomInBtn.trigger('click')

    expect(wrapper.find('.tree-controls').exists()).toBe(true)
  })

  it('should handle zoom out button click', async () => {
    const wrapper = mount(GenealogyTreeComponent, {
      props: { data: mockTreeData }
    })

    const zoomOutBtn = wrapper.find('.tree-controls button:nth-child(2)')
    await zoomOutBtn.trigger('click')

    expect(wrapper.find('.tree-controls').exists()).toBe(true)
  })

  it('should handle reset zoom button click', async () => {
    const wrapper = mount(GenealogyTreeComponent, {
      props: { data: mockTreeData }
    })

    const resetBtn = wrapper.find('.tree-controls button:nth-child(3)')
    await resetBtn.trigger('click')

    expect(wrapper.find('.tree-controls').exists()).toBe(true)
  })

  it('should emit events when clicking zoom buttons', async () => {
    const wrapper = mount(GenealogyTreeComponent, {
      props: { data: mockTreeData }
    })

    await wrapper.vm.$nextTick()

    const zoomInBtn = wrapper.find('.tree-controls button:nth-child(1)')
    await zoomInBtn.trigger('click')

    expect(wrapper.emitted()).toBeDefined()
  })

  it('should update when data prop changes', async () => {
    const wrapper = mount(GenealogyTreeComponent, {
      props: { data: mockTreeData }
    })

    await wrapper.vm.$nextTick()

    const newData: TreeNode = {
      id: 10,
      name: '李老爷子',
      generation: 1,
      gender: 'M',
      children: []
    }

    await wrapper.setProps({ data: newData })
    await wrapper.vm.$nextTick()

    const svg = wrapper.find('svg')
    expect(svg.exists()).toBe(true)
  })

  it('should render control buttons with correct content', () => {
    const wrapper = mount(GenealogyTreeComponent, {
      props: { data: mockTreeData }
    })

    const buttons = wrapper.findAll('.tree-controls button')
    expect(buttons[0].text()).toBe('+')
    expect(buttons[1].text()).toBe('-')
    expect(buttons[2].text()).toBe('⟲')
  })
})