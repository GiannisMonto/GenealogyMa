import { config } from '@vue/test-utils'
import { vi } from 'vitest'

export default config({
  global: {
    mocks: {
      console: {
        error: vi.fn(),
        warn: vi.fn()
      }
    }
  }
})
