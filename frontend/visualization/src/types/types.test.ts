import { describe, it, expect } from 'vitest'
import type { Person, GenealogyNode, Statistics } from '@/types'

describe('Types', () => {
  describe('Person', () => {
    it('should accept valid person data', () => {
      const person: Person = {
        member_id: 1,
        name: '张三',
        style_name: '正一',
        gender: 'M',
        generation: 5,
        lineage_path: '1.2.3.4.5',
        father_id: 10,
        birth_year: 1900,
        death_year: 1970,
        birth_place: '浙江',
        burial_place: '浙江',
        biography: '生平简介',
        avatar_url: 'http://example.com/avatar.jpg'
      }

      expect(person.member_id).toBe(1)
      expect(person.name).toBe('张三')
      expect(person.gender).toBe('M')
      expect(person.generation).toBe(5)
    })

    it('should accept minimal person data', () => {
      const person: Person = {
        member_id: 1,
        name: '李四',
        gender: 'F',
        generation: 1,
        lineage_path: '1'
      }

      expect(person.member_id).toBe(1)
      expect(person.birth_year).toBeUndefined()
      expect(person.biography).toBeUndefined()
    })
  })

  describe('GenealogyNode', () => {
    it('should accept genealogy node with children', () => {
      const node: GenealogyNode = {
        id: 1,
        name: '王五',
        generation: 3,
        gender: 'M',
        children: [
          {
            id: 2,
            name: '王小一',
            generation: 4,
            gender: 'M',
            children: []
          }
        ],
        spouse: {
          id: 3,
          name: '王氏',
          generation: 3,
          gender: 'F'
        }
      }

      expect(node.children).toHaveLength(1)
      expect(node.spouse).toBeDefined()
    })
  })

  describe('Statistics', () => {
    it('should accept valid statistics data', () => {
      const stats: Statistics = {
        total_persons: 100,
        total_generations: 10,
        male_count: 55,
        female_count: 45,
        generation_distribution: [
          { generation: 1, count: 2 },
          { generation: 2, count: 4 },
          { generation: 3, count: 8 }
        ]
      }

      expect(stats.total_persons).toBe(100)
      expect(stats.generation_distribution).toHaveLength(3)
      expect(stats.male_count + stats.female_count).toBe(stats.total_persons)
    })
  })
})
