import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Person } from '@/types'

export const usePersonStore = defineStore('person', () => {
  const persons = ref<Person[]>([])
  const currentPerson = ref<Person | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchPersons() {
    loading.value = true
    error.value = null
    try {
      const response = await fetch('/api/v1/persons?page=1&page_size=100')
      if (!response.ok) throw new Error('Failed to fetch persons')
      const data = await response.json()
      persons.value = data.data || []
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Unknown error'
    } finally {
      loading.value = false
    }
  }

  async function fetchPersonById(id: number) {
    loading.value = true
    error.value = null
    try {
      const response = await fetch(`/api/v1/persons/${id}`)
      if (!response.ok) throw new Error('Failed to fetch person')
      const data = await response.json()
      currentPerson.value = data.data
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Unknown error'
    } finally {
      loading.value = false
    }
  }

  async function fetchGenealogyTree(rootId?: number) {
    loading.value = true
    error.value = null
    try {
      const url = rootId
        ? `/api/v1/persons/${rootId}/genealogy-tree`
        : '/api/v1/persons/genealogy-tree'
      const response = await fetch(url)
      if (!response.ok) throw new Error('Failed to fetch genealogy tree')
      const data = await response.json()
      persons.value = data.data || []
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Unknown error'
    } finally {
      loading.value = false
    }
  }

  return {
    persons,
    currentPerson,
    loading,
    error,
    fetchPersons,
    fetchPersonById,
    fetchGenealogyTree
  }
})
