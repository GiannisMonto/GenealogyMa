<template>
  <div class="person-detail-page">
    <header class="page-header">
      <router-link to="/" class="back-link">← 返回</router-link>
      <h1>{{ person?.name || '人物详情' }}</h1>
    </header>

    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="error" class="error">{{ error }}</div>

    <div v-else-if="person" class="detail-content">
      <div class="info-card">
        <div class="avatar">{{ person.name?.charAt(0) }}</div>
        <div class="basic-info">
          <h2>{{ person.name }}</h2>
          <p v-if="person.style_name">字：{{ person.style_name }}</p>
          <span class="badge">第{{ person.generation }}代</span>
        </div>
      </div>

      <div class="info-section">
        <h3>基本信息</h3>
        <div class="info-grid">
          <div class="info-item">
            <label>性别</label>
            <span>{{ person.gender === 'M' ? '男' : '女' }}</span>
          </div>
          <div class="info-item">
            <label>出生年份</label>
            <span>{{ person.birth_year || '未知' }}</span>
          </div>
          <div class="info-item">
            <label>逝世年份</label>
            <span>{{ person.death_year || '未知' }}</span>
          </div>
          <div class="info-item">
            <label>出生地</label>
            <span>{{ person.birth_place || '未知' }}</span>
          </div>
          <div class="info-item">
            <label>安葬地</label>
            <span>{{ person.burial_place || '未知' }}</span>
          </div>
        </div>
      </div>

      <div v-if="person.biography" class="info-section">
        <h3>生平简介</h3>
        <p>{{ person.biography }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { usePersonStore } from '@/stores/person'

const route = useRoute()
const personStore = usePersonStore()

const loading = ref(false)
const error = ref<string | null>(null)
const person = ref(personStore.currentPerson)

async function loadPerson() {
  const id = Number(route.params.id)
  if (!id) return

  loading.value = true
  error.value = null
  try {
    await personStore.fetchPersonById(id)
    person.value = personStore.currentPerson
  } catch (e) {
    error.value = e instanceof Error ? e.message : '加载失败'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadPerson()
})
</script>

<style scoped>
.person-detail-page {
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
  font-size: 1rem;
}

.detail-content {
  max-width: 800px;
  margin: 0 auto;
  padding: 24px;
}

.info-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 24px;
  margin-bottom: 24px;
}

.avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: var(--color-primary);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2rem;
  font-weight: 600;
}

.badge {
  display: inline-block;
  background: var(--color-secondary);
  color: white;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 0.9rem;
}

.info-section {
  background: white;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 24px;
}

.info-section h3 {
  color: var(--color-primary);
  margin-bottom: 16px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-item label {
  font-size: 0.85rem;
  color: #888;
}

.loading,
.error {
  text-align: center;
  padding: 60px 20px;
  font-size: 1.1rem;
}

.error {
  color: #dc3545;
}
</style>
