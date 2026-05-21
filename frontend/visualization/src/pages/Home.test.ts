import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import Home from '@/pages/Home.vue'
import { createRouter, createWebHistory } from 'vue-router'

describe('Home Page', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('should render page title', () => {
    const wrapper = mount(Home, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/', name: 'Home', component: Home }]
          })
        ]
      }
    })

    expect(wrapper.find('h1').text()).toBe('族谱可视化')
  })

  it('should render navigation cards', () => {
    const wrapper = mount(Home, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/', name: 'Home', component: Home }]
          })
        ]
      }
    })

    const navCards = wrapper.findAll('.nav-card')
    expect(navCards).toHaveLength(4)
  })

  it('should render genealogy tree nav card', () => {
    const wrapper = mount(Home, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/', name: 'Home', component: Home }]
          })
        ]
      }
    })

    const genealogyCard = wrapper.find('.nav-card:nth-child(1)')
    expect(genealogyCard.find('.nav-title').text()).toBe('族谱树')
    expect(genealogyCard.find('.nav-icon').text()).toBe('🌳')
  })
})
