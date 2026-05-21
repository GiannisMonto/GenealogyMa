import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import PersonDetailComponent from '@/components/PersonDetail/index.vue'
import type { Person } from '@/types'

const mockPerson: Person = {
  member_id: 1,
  name: '张三',
  style_name: '字云长',
  gender: 'M',
  generation: 5,
  lineage_path: '1.2.3.4.5',
  birth_year: 1950,
  death_year: 2020,
  birth_place: '湖北武汉',
  burial_place: '湖北武汉祖山',
  biography: '张三的生平简介...'
}

describe('PersonDetail Component', () => {
  it('should render person name', () => {
    const wrapper = mount(PersonDetailComponent, {
      props: { person: mockPerson }
    })

    expect(wrapper.find('h2').text()).toBe('张三')
  })

  it('should render style name when provided', () => {
    const wrapper = mount(PersonDetailComponent, {
      props: { person: mockPerson }
    })

    expect(wrapper.text()).toContain('字：字云长')
  })

  it('should render generation badge', () => {
    const wrapper = mount(PersonDetailComponent, {
      props: { person: mockPerson }
    })

    expect(wrapper.find('.generation-badge').text()).toBe('第5代')
  })

  it('should render gender correctly', () => {
    const wrapper = mount(PersonDetailComponent, {
      props: { person: mockPerson }
    })

    expect(wrapper.text()).toContain('男')
  })

  it('should render birth year when provided', () => {
    const wrapper = mount(PersonDetailComponent, {
      props: { person: mockPerson }
    })

    expect(wrapper.text()).toContain('1950')
  })

  it('should render death year when provided', () => {
    const wrapper = mount(PersonDetailComponent, {
      props: { person: mockPerson }
    })

    expect(wrapper.text()).toContain('2020')
  })

  it('should render biography when provided', () => {
    const wrapper = mount(PersonDetailComponent, {
      props: { person: mockPerson }
    })

    expect(wrapper.find('.biography').text()).toContain('张三的生平简介...')
  })

  it('should emit close event when close button clicked', async () => {
    const wrapper = mount(PersonDetailComponent, {
      props: { person: mockPerson }
    })

    await wrapper.find('.close-btn').trigger('click')

    expect(wrapper.emitted('close')).toBeTruthy()
  })

  it('should not render biography section when not provided', () => {
    const personWithoutBio: Person = { ...mockPerson, biography: undefined }
    const wrapper = mount(PersonDetailComponent, {
      props: { person: personWithoutBio }
    })

    expect(wrapper.find('.biography').exists()).toBe(false)
  })

  it('should show avatar with first character of name', () => {
    const wrapper = mount(PersonDetailComponent, {
      props: { person: mockPerson }
    })

    expect(wrapper.find('.avatar').text()).toBe('张')
  })
})