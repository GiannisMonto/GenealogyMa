<template>
  <div class="person-detail-overlay" @click.self="$emit('close')">
    <div class="person-detail">
      <button class="close-btn" @click="$emit('close')">×</button>

      <div class="person-header">
        <div class="avatar">{{ person.name.charAt(0) }}</div>
        <div class="basic-info">
          <h2>{{ person.name }}</h2>
          <p v-if="person.style_name">字：{{ person.style_name }}</p>
          <span class="generation-badge">第{{ person.generation }}代</span>
        </div>
      </div>

      <div class="info-section">
        <h3>基本信息</h3>
        <div class="info-grid">
          <div class="info-item">
            <label>性别</label>
            <span>{{ person.gender === 'M' ? '男' : '女' }}</span>
          </div>
          <div class="info-item" v-if="person.birth_year">
            <label>出生年份</label>
            <span>{{ person.birth_year }}</span>
          </div>
          <div class="info-item" v-if="person.death_year">
            <label>逝世年份</label>
            <span>{{ person.death_year }}</span>
          </div>
          <div class="info-item" v-if="person.birth_place">
            <label>出生地</label>
            <span>{{ person.birth_place }}</span>
          </div>
          <div class="info-item" v-if="person.burial_place">
            <label>安葬地</label>
            <span>{{ person.burial_place }}</span>
          </div>
        </div>
      </div>

      <div class="info-section" v-if="person.biography">
        <h3>生平简介</h3>
        <p class="biography">{{ person.biography }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Person } from '@/types'

defineProps<{
  person: Person
}>()

defineEmits<{
  (e: 'close'): void
}>()
</script>

<style scoped>
.person-detail-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.person-detail {
  background: white;
  border-radius: 16px;
  width: 480px;
  max-height: 80vh;
  overflow-y: auto;
  padding: 24px;
  position: relative;
}

.close-btn {
  position: absolute;
  top: 16px;
  right: 16px;
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #666;
}

.person-header {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 24px;
  padding-bottom: 24px;
  border-bottom: 1px solid var(--color-border);
}

.avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: var(--color-primary);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  font-weight: 600;
}

.basic-info h2 {
  margin: 0 0 4px;
  color: var(--color-primary);
}

.generation-badge {
  display: inline-block;
  background: var(--color-secondary);
  color: white;
  padding: 2px 10px;
  border-radius: 12px;
  font-size: 0.85rem;
  margin-top: 8px;
}

.info-section {
  margin-bottom: 20px;
}

.info-section h3 {
  font-size: 1rem;
  color: var(--color-primary);
  margin-bottom: 12px;
}

.info-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
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

.info-item span {
  font-size: 1rem;
}

.biography {
  line-height: 1.8;
  color: #444;
}
</style>
