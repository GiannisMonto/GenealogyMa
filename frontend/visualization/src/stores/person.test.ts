import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { usePersonStore } from '@/stores/person'
import type { Person } from '@/types'

describe('Person Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('should have empty initial state', () => {
    const store = usePersonStore()
    expect(store.persons).toEqual([])
    expect(store.currentPerson).toBeNull()
    expect(store.loading).toBe(false)
    expect(store.error).toBeNull()
  })

  it('should set currentPerson when fetchPersonById is called', async () => {
    const mockPerson: Person = {
      member_id: 1,
      name: '测试人物',
      gender: 'M',
      generation: 1,
      lineage_path: '1'
    }

    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ code: 0, data: mockPerson })
    }))

    const store = usePersonStore()
    await store.fetchPersonById(1)

    expect(store.currentPerson).toEqual(mockPerson)
    expect(store.loading).toBe(false)
    expect(store.error).toBeNull()
  })

  it('should set error when fetchPersonById fails', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: false
    }))

    const store = usePersonStore()
    await store.fetchPersonById(999)

    expect(store.error).toBe('Failed to fetch person')
    expect(store.loading).toBe(false)
  })

  it('should set persons when fetchGenealogyTree is called', async () => {
    const mockPersons: Person[] = [
      { member_id: 1, name: '祖先', gender: 'M', generation: 1, lineage_path: '1' },
      { member_id: 2, name: '父亲', gender: 'M', generation: 2, lineage_path: '1.2', father_id: 1 },
      { member_id: 3, name: '儿子', gender: 'M', generation: 3, lineage_path: '1.2.3', father_id: 2 }
    ]

    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ code: 0, data: mockPersons })
    }))

    const store = usePersonStore()
    await store.fetchGenealogyTree()

    expect(store.persons).toEqual(mockPersons)
  })
})
