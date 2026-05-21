import { describe, it, expect } from 'vitest'
import { personApi, statisticsApi } from '@/api/client'

describe('API Client', () => {
  describe('personApi', () => {
    it('should have list method', () => {
      expect(typeof personApi.list).toBe('function')
    })

    it('should have get method', () => {
      expect(typeof personApi.get).toBe('function')
    })

    it('should have getGenealogyTree method', () => {
      expect(typeof personApi.getGenealogyTree).toBe('function')
    })

    it('should have search method', () => {
      expect(typeof personApi.search).toBe('function')
    })
  })

  describe('statisticsApi', () => {
    it('should have get method', () => {
      expect(typeof statisticsApi.get).toBe('function')
    })
  })
})
