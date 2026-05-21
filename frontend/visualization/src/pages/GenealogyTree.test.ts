import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import GenealogyTreePage from '@/pages/GenealogyTree.vue'
import { createRouter, createWebHistory } from 'vue-router'

describe('GenealogyTree Page', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    setActivePinia(createPinia())
  })

  it('should render page header with title', () => {
    const wrapper = mount(GenealogyTreePage, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/genealogy', name: 'Genealogy', component: GenealogyTreePage }]
          })
        ],
        stubs: {
          GenealogyTree: true,
          PersonDetail: true
        }
      }
    })

    expect(wrapper.find('h1').text()).toBe('族谱树')
  })

  it('should render search input and buttons', () => {
    const wrapper = mount(GenealogyTreePage, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/genealogy', name: 'Genealogy', component: GenealogyTreePage }]
          })
        ],
        stubs: {
          GenealogyTree: true,
          PersonDetail: true
        }
      }
    })

    expect(wrapper.find('.search-input').exists()).toBe(true)
    expect(wrapper.find('.search-btn').exists()).toBe(true)
    expect(wrapper.find('.reset-btn').exists()).toBe(true)
  })

  it('should render search input with placeholder', () => {
    const wrapper = mount(GenealogyTreePage, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/genealogy', name: 'Genealogy', component: GenealogyTreePage }]
          })
        ],
        stubs: {
          GenealogyTree: true,
          PersonDetail: true
        }
      }
    })

    const searchInput = wrapper.find('.search-input')
    expect(searchInput.attributes('placeholder')).toBe('搜索人物...')
  })

  it('should have header controls with correct classes', () => {
    const wrapper = mount(GenealogyTreePage, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/genealogy', name: 'Genealogy', component: GenealogyTreePage }]
          })
        ],
        stubs: {
          GenealogyTree: true,
          PersonDetail: true
        }
      }
    })

    expect(wrapper.find('.header-controls').exists()).toBe(true)
  })

  it('should have tree container', () => {
    const wrapper = mount(GenealogyTreePage, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/genealogy', name: 'Genealogy', component: GenealogyTreePage }]
          })
        ],
        stubs: {
          GenealogyTree: true,
          PersonDetail: true
        }
      }
    })

    expect(wrapper.find('.tree-container').exists()).toBe(true)
  })
})