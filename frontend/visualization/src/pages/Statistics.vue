<template>
  <div class="statistics-page">
    <header class="page-header">
      <router-link to="/" class="back-link">← 返回</router-link>
      <h1>世代分布</h1>
    </header>

    <div class="statistics-content">
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="error" class="error">{{ error }}</div>
      <template v-else-if="statistics">
        <div class="stats-overview">
          <div class="stat-card">
            <span class="stat-value">{{ statistics.total_persons }}</span>
            <span class="stat-label">总人数</span>
          </div>
          <div class="stat-card">
            <span class="stat-value">{{ statistics.total_generations }}</span>
            <span class="stat-label">总代数</span>
          </div>
          <div class="stat-card male">
            <span class="stat-value">{{ statistics.male_count }}</span>
            <span class="stat-label">男性</span>
          </div>
          <div class="stat-card female">
            <span class="stat-value">{{ statistics.female_count }}</span>
            <span class="stat-label">女性</span>
          </div>
        </div>

        <div class="distribution-section">
          <h2>各世代人数分布</h2>
          <div class="distribution-chart">
            <div
              v-for="item in statistics.generation_distribution"
              :key="item.generation"
              class="bar-item"
            >
              <div class="bar-label">第{{ item.generation }}代</div>
              <div class="bar-container">
                <div
                  class="bar-fill"
                  :style="{ width: getBarWidth(item.count) + '%' }"
                ></div>
              </div>
              <div class="bar-value">{{ item.count }}人</div>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { statisticsApi } from '@/api/client'
import type { Statistics } from '@/types'

const statistics = ref<Statistics | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)

const maxCount = computed(() => {
  if (!statistics.value?.generation_distribution.length) return 0
  return Math.max(...statistics.value.generation_distribution.map((item) => item.count))
})

function getBarWidth(count: number): number {
  if (maxCount.value === 0) return 0
  return (count / maxCount.value) * 100
}

async function loadStatistics() {
  loading.value = true
  error.value = null
  try {
    const response = await statisticsApi.get()
    statistics.value = response.data.data
  } catch (e) {
    error.value = e instanceof Error ? e.message : '加载失败'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadStatistics()
})
</script>

<style scoped>
.statistics-page {
  width: 100%;
  min-height: 100vh;
  background: var(--color-background);
}

.page-header {
  background: var(--color-primary);
  color: white;
  padding: 16px 24px;
  display: flex;
  align-items: center;
  gap: 16px;
}

.back-link {
  color: white;
  text-decoration: none;
  font-size: 0.9rem;
}

.page-header h1 {
  font-size: 1.5rem;
  flex: 1;
}

.statistics-content {
  max-width: 800px;
  margin: 0 auto;
  padding: 24px;
}

.loading,
.error {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 200px;
  font-size: 1.1rem;
}

.error {
  color: #dc3545;
}

.stats-overview {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 32px;
}

.stat-card {
  background: white;
  border-radius: 12px;
  padding: 24px 16px;
  text-align: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.stat-value {
  display: block;
  font-size: 2rem;
  font-weight: 700;
  color: var(--color-primary);
  margin-bottom: 8px;
}

.stat-label {
  font-size: 0.9rem;
  color: #666;
}

.stat-card.male .stat-value {
  color: #1976d2;
}

.stat-card.female .stat-value {
  color: #e91e63;
}

.distribution-section {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.distribution-section h2 {
  font-size: 1.2rem;
  margin-bottom: 24px;
  color: #333;
}

.distribution-chart {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.bar-item {
  display: grid;
  grid-template-columns: 60px 1fr 60px;
  align-items: center;
  gap: 12px;
}

.bar-label {
  font-size: 0.9rem;
  color: #666;
  text-align: right;
}

.bar-container {
  height: 24px;
  background: #f0f0f0;
  border-radius: 4px;
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  background: var(--color-primary);
  border-radius: 4px;
  transition: width 0.3s ease;
  min-width: 2px;
}

.bar-value {
  font-size: 0.9rem;
  color: #333;
  font-weight: 500;
}
</style>
