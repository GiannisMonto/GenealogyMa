import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import Statistics from '@/pages/Statistics.vue'
import { createRouter, createWebHistory } from 'vue-router'
import { statisticsApi } from '@/api/client'

vi.mock('@/api/client', () => ({
  statisticsApi: {
    get: vi.fn()
  }
}))

const mockStatistics = {
  total_persons: 100,
  total_generations: 5,
  male_count: 55,
  female_count: 45,
  generation_distribution: [
    { generation: 1, count: 2 },
    { generation: 2, count: 4 },
    { generation: 3, count: 10 },
    { generation: 4, count: 30 },
    { generation: 5, count: 54 }
  ]
}

describe('Statistics Page', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('should render page title', async () => {
    (statisticsApi.get as any).mockResolvedValue({ data: { data: mockStatistics } })

    const wrapper = mount(Statistics, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/statistics', name: 'Statistics', component: Statistics }]
          })
        ]
      },
      stubs: {
        'router-link': { template: '<a class="router-link-stub"><slot /></a>' }
      }
    })

    await wrapper.vm.$nextTick()
    expect(wrapper.find('h1').text()).toBe('世代分布')
  })

  it('should render back link', async () => {
    (statisticsApi.get as any).mockResolvedValue({ data: { data: mockStatistics } })

    const wrapper = mount(Statistics, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/statistics', name: 'Statistics', component: Statistics }]
          })
        ]
      },
      stubs: {
        'router-link': { template: '<a class="router-link-stub"><slot /></a>' }
      }
    })

    await wrapper.vm.$nextTick()
    expect(wrapper.find('.back-link').text()).toBe('← 返回')
  })

  it('should render loading state initially', async () => {
    let resolveGet: any
    const pendingPromise = new Promise((resolve) => {
      resolveGet = resolve
    })
    ;(statisticsApi.get as any).mockReturnValue(pendingPromise)

    const wrapper = mount(Statistics, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/statistics', name: 'Statistics', component: Statistics }]
          })
        ]
      }
    })

    await wrapper.vm.$nextTick()
    expect(wrapper.find('.loading').text()).toBe('加载中...')

    resolveGet({ data: { data: mockStatistics } })
  })

  it('should render statistics cards after loading', async () => {
    (statisticsApi.get as any).mockResolvedValue({ data: { data: mockStatistics } })

    const wrapper = mount(Statistics, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/statistics', name: 'Statistics', component: Statistics }]
          })
        ]
      },
      stubs: {
        'router-link': { template: '<a class="router-link-stub"><slot /></a>' }
      }
    })

    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()

    const statCards = wrapper.findAll('.stat-card')
    expect(statCards).toHaveLength(4)
  })

  it('should render correct total persons count', async () => {
    (statisticsApi.get as any).mockResolvedValue({ data: { data: mockStatistics } })

    const wrapper = mount(Statistics, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/statistics', name: 'Statistics', component: Statistics }]
          })
        ]
      },
      stubs: {
        'router-link': { template: '<a class="router-link-stub"><slot /></a>' }
      }
    })

    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()

    const statValues = wrapper.findAll('.stat-value')
    expect(statValues[0].text()).toBe('100')
  })

  it('should render correct total generations count', async () => {
    (statisticsApi.get as any).mockResolvedValue({ data: { data: mockStatistics } })

    const wrapper = mount(Statistics, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/statistics', name: 'Statistics', component: Statistics }]
          })
        ]
      },
      stubs: {
        'router-link': { template: '<a class="router-link-stub"><slot /></a>' }
      }
    })

    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()

    const statValues = wrapper.findAll('.stat-value')
    expect(statValues[1].text()).toBe('5')
  })

  it('should render correct male and female counts', async () => {
    (statisticsApi.get as any).mockResolvedValue({ data: { data: mockStatistics } })

    const wrapper = mount(Statistics, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/statistics', name: 'Statistics', component: Statistics }]
          })
        ]
      },
      stubs: {
        'router-link': { template: '<a class="router-link-stub"><slot /></a>' }
      }
    })

    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()

    const statValues = wrapper.findAll('.stat-value')
    expect(statValues[2].text()).toBe('55')
    expect(statValues[3].text()).toBe('45')
  })

  it('should render generation distribution bars', async () => {
    (statisticsApi.get as any).mockResolvedValue({ data: { data: mockStatistics } })

    const wrapper = mount(Statistics, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/statistics', name: 'Statistics', component: Statistics }]
          })
        ]
      },
      stubs: {
        'router-link': { template: '<a class="router-link-stub"><slot /></a>' }
      }
    })

    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()

    const barItems = wrapper.findAll('.bar-item')
    expect(barItems).toHaveLength(5)
  })

  it('should render correct generation labels', async () => {
    (statisticsApi.get as any).mockResolvedValue({ data: { data: mockStatistics } })

    const wrapper = mount(Statistics, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/statistics', name: 'Statistics', component: Statistics }]
          })
        ]
      },
      stubs: {
        'router-link': { template: '<a class="router-link-stub"><slot /></a>' }
      }
    })

    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()

    const barLabels = wrapper.findAll('.bar-label')
    expect(barLabels[0].text()).toBe('第1代')
    expect(barLabels[4].text()).toBe('第5代')
  })

  it('should render correct bar values', async () => {
    (statisticsApi.get as any).mockResolvedValue({ data: { data: mockStatistics } })

    const wrapper = mount(Statistics, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/statistics', name: 'Statistics', component: Statistics }]
          })
        ]
      },
      stubs: {
        'router-link': { template: '<a class="router-link-stub"><slot /></a>' }
      }
    })

    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()

    const barValues = wrapper.findAll('.bar-value')
    expect(barValues[0].text()).toBe('2人')
    expect(barValues[4].text()).toBe('54人')
  })

  it('should render section title', async () => {
    (statisticsApi.get as any).mockResolvedValue({ data: { data: mockStatistics } })

    const wrapper = mount(Statistics, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/statistics', name: 'Statistics', component: Statistics }]
          })
        ]
      },
      stubs: {
        'router-link': { template: '<a class="router-link-stub"><slot /></a>' }
      }
    })

    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.distribution-section h2').text()).toBe('各世代人数分布')
  })

  it('should render error state on API failure', async () => {
    (statisticsApi.get as any).mockRejectedValue(new Error('Network error'))

    const wrapper = mount(Statistics, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/statistics', name: 'Statistics', component: Statistics }]
          })
        ]
      },
      stubs: {
        'router-link': { template: '<a class="router-link-stub"><slot /></a>' }
      }
    })

    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.error').text()).toBe('Network error')
  })

  it('should calculate max bar width correctly', async () => {
    (statisticsApi.get as any).mockResolvedValue({ data: { data: mockStatistics } })

    const wrapper = mount(Statistics, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/statistics', name: 'Statistics', component: Statistics }]
          })
        ]
      },
      stubs: {
        'router-link': { template: '<a class="router-link-stub"><slot /></a>' }
      }
    })

    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()

    const bars = wrapper.findAll('.bar-fill')
    const lastBar = bars[4]
    expect(lastBar.attributes('style')).toContain('width: 100%')
  })

  it('should handle empty statistics gracefully', async () => {
    const emptyStats = {
      total_persons: 0,
      total_generations: 0,
      male_count: 0,
      female_count: 0,
      generation_distribution: []
    }
    ;(statisticsApi.get as any).mockResolvedValue({ data: { data: emptyStats } })

    const wrapper = mount(Statistics, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/statistics', name: 'Statistics', component: Statistics }]
          })
        ]
      },
      stubs: {
        'router-link': { template: '<a class="router-link-stub"><slot /></a>' }
      }
    })

    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()

    expect(wrapper.findAll('.bar-item')).toHaveLength(0)
  })
})
