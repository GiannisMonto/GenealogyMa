import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import SearchPanel from '@/components/SearchPanel/index.vue'
import type { Person } from '@/types'

const mockPersons: Person[] = [
  {
    member_id: 1,
    name: '张三',
    style_name: '字云长',
    gender: 'M',
    generation: 5,
    lineage_path: '1.2.3.4.5'
  },
  {
    member_id: 2,
    name: '李四',
    style_name: '字翼德',
    gender: 'M',
    generation: 6,
    lineage_path: '1.2.3.4.5.6'
  },
  {
    member_id: 3,
    name: '王芳',
    style_name: '字貂蝉',
    gender: 'F',
    generation: 5,
    lineage_path: '1.2.3.4.6'
  },
  {
    member_id: 4,
    name: '赵六',
    gender: 'M',
    generation: 7,
    lineage_path: '1.2.3.4.5.6.7'
  },
  {
    member_id: 5,
    name: '张七',
    style_name: '字孟德',
    gender: 'M',
    generation: 6,
    lineage_path: '1.2.3.4.5.8'
  }
]

describe('SearchPanel Component', () => {
  it('should render search input', () => {
    const wrapper = mount(SearchPanel, {
      props: { persons: mockPersons }
    })

    expect(wrapper.find('.search-input').exists()).toBe(true)
  })

  it('should render search button', () => {
    const wrapper = mount(SearchPanel, {
      props: { persons: mockPersons }
    })

    expect(wrapper.find('.search-btn').exists()).toBe(true)
  })

  it('should use default placeholder when not provided', () => {
    const wrapper = mount(SearchPanel, {
      props: { persons: mockPersons }
    })

    const input = wrapper.find('.search-input')
    expect(input.attributes('placeholder')).toBeUndefined()
  })

  it('should use custom placeholder when provided', () => {
    const wrapper = mount(SearchPanel, {
      props: { persons: mockPersons, placeholder: '搜索人物姓名...' }
    })

    const input = wrapper.find('.search-input')
    expect(input.attributes('placeholder')).toBe('搜索人物姓名...')
  })

  it('should filter persons by name', async () => {
    const wrapper = mount(SearchPanel, {
      props: { persons: mockPersons }
    })

    const input = wrapper.find('.search-input')
    await input.setValue('张三')

    await wrapper.find('.search-input').trigger('input')

    expect(wrapper.findAll('.search-result-item')).toHaveLength(1)
    expect(wrapper.find('.person-name').text()).toBe('张三')
  })

  it('should filter persons by style name', async () => {
    const wrapper = mount(SearchPanel, {
      props: { persons: mockPersons }
    })

    const input = wrapper.find('.search-input')
    await input.setValue('云长')

    await wrapper.find('.search-input').trigger('input')

    expect(wrapper.findAll('.search-result-item')).toHaveLength(1)
    expect(wrapper.find('.person-name').text()).toBe('张三')
  })

  it('should show multiple results for partial match', async () => {
    const wrapper = mount(SearchPanel, {
      props: { persons: mockPersons }
    })

    const input = wrapper.find('.search-input')
    await input.setValue('张')

    await wrapper.find('.search-input').trigger('input')

    expect(wrapper.findAll('.search-result-item')).toHaveLength(2)
  })

  it('should render person gender correctly', async () => {
    const wrapper = mount(SearchPanel, {
      props: { persons: mockPersons }
    })

    const input = wrapper.find('.search-input')
    await input.setValue('王芳')

    await wrapper.find('.search-input').trigger('input')

    expect(wrapper.find('.person-gender').text()).toBe('女')
  })

  it('should render person generation correctly', async () => {
    const wrapper = mount(SearchPanel, {
      props: { persons: mockPersons }
    })

    const input = wrapper.find('.search-input')
    await input.setValue('李四')

    await wrapper.find('.search-input').trigger('input')

    expect(wrapper.find('.person-generation').text()).toBe('第6代')
  })

  it('should not show results when input is empty', async () => {
    const wrapper = mount(SearchPanel, {
      props: { persons: mockPersons }
    })

    const input = wrapper.find('.search-input')
    await input.setValue('')

    await wrapper.find('.search-input').trigger('input')

    expect(wrapper.find('.search-results').exists()).toBe(false)
  })

  it('should show no results message when no match', async () => {
    const wrapper = mount(SearchPanel, {
      props: { persons: mockPersons }
    })

    const input = wrapper.find('.search-input')
    await input.setValue('不存在的名字')

    await wrapper.find('.search-input').trigger('input')

    expect(wrapper.find('.no-results').exists()).toBe(true)
    expect(wrapper.find('.no-results').text()).toBe('未找到匹配的人物')
  })

  it('should emit select event when person clicked', async () => {
    const wrapper = mount(SearchPanel, {
      props: { persons: mockPersons }
    })

    const input = wrapper.find('.search-input')
    await input.setValue('张三')

    await wrapper.find('.search-input').trigger('input')

    await wrapper.find('.search-result-item').trigger('click')

    expect(wrapper.emitted('select')).toBeTruthy()
    const selectEvent = wrapper.emitted('select')![0]
    expect(selectEvent[0]).toEqual(mockPersons[0])
  })

  it('should search by name case insensitively', async () => {
    const wrapper = mount(SearchPanel, {
      props: { persons: mockPersons }
    })

    const input = wrapper.find('.search-input')
    await input.setValue('ZHANGSAN')

    await wrapper.find('.search-input').trigger('input')

    expect(wrapper.findAll('.search-result-item')).toHaveLength(0)
  })

  it('should limit results to maxResults', async () => {
    const manyPersons: Person[] = Array.from({ length: 20 }, (_, i) => ({
      member_id: i + 1,
      name: `测试用户${i + 1}`,
      gender: 'M',
      generation: 1,
      lineage_path: `${i + 1}`
    }))

    const wrapper = mount(SearchPanel, {
      props: { persons: manyPersons, maxResults: 5 }
    })

    const input = wrapper.find('.search-input')
    await input.setValue('测试')

    await wrapper.find('.search-input').trigger('input')

    expect(wrapper.findAll('.search-result-item')).toHaveLength(5)
  })

  it('should not render style name if not provided', async () => {
    const wrapper = mount(SearchPanel, {
      props: { persons: mockPersons }
    })

    const input = wrapper.find('.search-input')
    await input.setValue('赵六')

    await wrapper.find('.search-input').trigger('input')

    expect(wrapper.find('.person-style-name').exists()).toBe(false)
  })

  it('should navigate results with keyboard', async () => {
    const wrapper = mount(SearchPanel, {
      props: { persons: mockPersons }
    })

    const input = wrapper.find('.search-input')
    await input.setValue('张')

    await wrapper.find('.search-input').trigger('input')

    await wrapper.find('.search-input').trigger('keydown.down')

    expect(wrapper.find('.search-result-item.highlighted').exists()).toBe(true)
  })

  it('should close results on escape', async () => {
    const wrapper = mount(SearchPanel, {
      props: { persons: mockPersons }
    })

    const input = wrapper.find('.search-input')
    await input.setValue('张')

    await wrapper.find('.search-input').trigger('input')

    await wrapper.find('.search-input').trigger('keydown.escape')

    expect(wrapper.find('.search-results').exists()).toBe(false)
  })
})
