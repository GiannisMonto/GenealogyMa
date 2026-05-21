import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import PersonDetail from '@/pages/PersonDetail.vue'
import { createRouter, createWebHistory } from 'vue-router'
import { usePersonStore } from '@/stores/person'

vi.mock('@/stores/person', () => ({
  usePersonStore: vi.fn()
}))

const mockPerson = {
  member_id: 1,
  name: '张三',
  style_name: '字伯约',
  gender: 'M',
  generation: 5,
  lineage_path: '1.2.3.4.5',
  father_id: 2,
  birth_year: 1950,
  death_year: 2020,
  birth_place: '江西赣州',
  burial_place: '江西赣州',
  biography: '张三的生平简介'
}

describe('PersonDetail Page', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('should render page title with person name', async () => {
    const mockStore = {
      currentPerson: mockPerson,
      fetchPersonById: vi.fn()
    };
    (usePersonStore as any).mockReturnValue(mockStore)

    const wrapper = mount(PersonDetail, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/person/:id', name: 'PersonDetail', component: PersonDetail }]
          })
        ]
      },
      props: {},
      stubs: {
        'router-link': { template: '<a class="router-link-stub"><slot /></a>' }
      }
    })

    await wrapper.vm.$nextTick()

    expect(wrapper.find('h1').text()).toBe('张三')
  })

  it('should render back link', () => {
    const mockStore = {
      currentPerson: mockPerson,
      fetchPersonById: vi.fn()
    };
    (usePersonStore as any).mockReturnValue(mockStore)

    const wrapper = mount(PersonDetail, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/person/:id', name: 'PersonDetail', component: PersonDetail }]
          })
        ]
      },
      stubs: {
        'router-link': { template: '<a class="router-link-stub"><slot /></a>' }
      }
    })

    expect(wrapper.find('.back-link').text()).toBe('← 返回')
  })

  it('should render avatar with first character of name', () => {
    const mockStore = {
      currentPerson: mockPerson,
      fetchPersonById: vi.fn()
    };
    (usePersonStore as any).mockReturnValue(mockStore)

    const wrapper = mount(PersonDetail, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/person/:id', name: 'PersonDetail', component: PersonDetail }]
          })
        ]
      },
      stubs: {
        'router-link': { template: '<a class="router-link-stub"><slot /></a>' }
      }
    })

    expect(wrapper.find('.avatar').text()).toBe('张')
  })

  it('should render generation badge', () => {
    const mockStore = {
      currentPerson: mockPerson,
      fetchPersonById: vi.fn()
    };
    (usePersonStore as any).mockReturnValue(mockStore)

    const wrapper = mount(PersonDetail, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/person/:id', name: 'PersonDetail', component: PersonDetail }]
          })
        ]
      },
      stubs: {
        'router-link': { template: '<a class="router-link-stub"><slot /></a>' }
      }
    })

    expect(wrapper.find('.badge').text()).toBe('第5代')
  })

  it('should render style name if available', () => {
    const mockStore = {
      currentPerson: mockPerson,
      fetchPersonById: vi.fn()
    };
    (usePersonStore as any).mockReturnValue(mockStore)

    const wrapper = mount(PersonDetail, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/person/:id', name: 'PersonDetail', component: PersonDetail }]
          })
        ]
      },
      stubs: {
        'router-link': { template: '<a class="router-link-stub"><slot /></a>' }
      }
    })

    expect(wrapper.find('.basic-info p').text()).toBe('字：字伯约')
  })

  it('should render gender correctly', () => {
    const malePerson = { ...mockPerson }
    const mockStore = {
      currentPerson: malePerson,
      fetchPersonById: vi.fn()
    };
    (usePersonStore as any).mockReturnValue(mockStore)

    const wrapper = mount(PersonDetail, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/person/:id', name: 'PersonDetail', component: PersonDetail }]
          })
        ]
      },
      stubs: {
        'router-link': { template: '<a class="router-link-stub"><slot /></a>' }
      }
    })

    const infoItems = wrapper.findAll('.info-item')
    const genderItem = infoItems.find(wr => wr.find('label').text() === '性别')
    expect(genderItem?.find('span').text()).toBe('男')
  })

  it('should render female gender correctly', () => {
    const femalePerson = { ...mockPerson, gender: 'F' }
    const mockStore = {
      currentPerson: femalePerson,
      fetchPersonById: vi.fn()
    };
    (usePersonStore as any).mockReturnValue(mockStore)

    const wrapper = mount(PersonDetail, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/person/:id', name: 'PersonDetail', component: PersonDetail }]
          })
        ]
      },
      stubs: {
        'router-link': { template: '<a class="router-link-stub"><slot /></a>' }
      }
    })

    const infoItems = wrapper.findAll('.info-item')
    const genderItem = infoItems.find(wr => wr.find('label').text() === '性别')
    expect(genderItem?.find('span').text()).toBe('女')
  })

  it('should render biography section if available', () => {
    const mockStore = {
      currentPerson: mockPerson,
      fetchPersonById: vi.fn()
    };
    (usePersonStore as any).mockReturnValue(mockStore)

    const wrapper = mount(PersonDetail, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/person/:id', name: 'PersonDetail', component: PersonDetail }]
          })
        ]
      },
      stubs: {
        'router-link': { template: '<a class="router-link-stub"><slot /></a>' }
      }
    })

    expect(wrapper.find('.info-section h3').text()).toBe('基本信息')
    expect(wrapper.find('.info-section:last-child h3').text()).toBe('生平简介')
  })

  it('should render info grid with all fields', () => {
    const mockStore = {
      currentPerson: mockPerson,
      fetchPersonById: vi.fn()
    };
    (usePersonStore as any).mockReturnValue(mockStore)

    const wrapper = mount(PersonDetail, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/person/:id', name: 'PersonDetail', component: PersonDetail }]
          })
        ]
      },
      stubs: {
        'router-link': { template: '<a class="router-link-stub"><slot /></a>' }
      }
    })

    const infoItems = wrapper.findAll('.info-item')
    expect(infoItems).toHaveLength(5)

    const labels = infoItems.map(el => el.find('label').text())
    expect(labels).toContain('性别')
    expect(labels).toContain('出生年份')
    expect(labels).toContain('逝世年份')
    expect(labels).toContain('出生地')
    expect(labels).toContain('安葬地')
  })

  it('should show default value for missing birth year', () => {
    const personWithoutBirthYear = { ...mockPerson, birth_year: undefined }
    const mockStore = {
      currentPerson: personWithoutBirthYear,
      fetchPersonById: vi.fn()
    };
    (usePersonStore as any).mockReturnValue(mockStore)

    const wrapper = mount(PersonDetail, {
      global: {
        plugins: [
          createRouter({
            history: createWebHistory(),
            routes: [{ path: '/person/:id', name: 'PersonDetail', component: PersonDetail }]
          })
        ]
      },
      stubs: {
        'router-link': { template: '<a class="router-link-stub"><slot /></a>' }
      }
    })

    const birthYearItem = wrapper.findAll('.info-item').find(el => el.find('label').text() === '出生年份')
    expect(birthYearItem?.find('span').text()).toBe('未知')
  })
})
