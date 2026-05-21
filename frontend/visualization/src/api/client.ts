import axios from 'axios'
import type { Person, Statistics } from '@/types'

const client = axios.create({
  baseURL: '/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

client.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

client.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
    }
    return Promise.reject(error)
  }
)

export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

export const personApi = {
  list(params: { page?: number; page_size?: number; search?: string }) {
    return client.get<ApiResponse<Person[]>>('/persons', { params })
  },

  get(id: number) {
    return client.get<ApiResponse<Person>>(`/persons/${id}`)
  },

  getGenealogyTree(rootId?: number) {
    const url = rootId ? `/persons/${rootId}/genealogy-tree` : '/persons/genealogy-tree'
    return client.get<ApiResponse<Person[]>>(url)
  },

  search(keyword: string) {
    return client.get<ApiResponse<Person[]>>('/persons/search', {
      params: { keyword }
    })
  }
}

export const statisticsApi = {
  get() {
    return client.get<ApiResponse<Statistics>>('/persons/statistics')
  }
}

export default client
