<template>
  <div class="search-panel">
    <div class="search-input-wrapper">
      <input
        v-model="searchQuery"
        type="text"
        :placeholder="placeholder"
        class="search-input"
        @input="handleInput"
        @focus="showResults = true"
        @keydown.down.prevent="navigateResults(1)"
        @keydown.up.prevent="navigateResults(-1)"
        @keydown.enter.prevent="selectHighlighted"
        @keydown.escape="showResults = false"
      />
      <button @click="performSearch" class="search-btn" :disabled="!searchQuery.trim()">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="11" cy="11" r="8"/>
          <path d="m21 21-4.35-4.35"/>
        </svg>
      </button>
    </div>

    <div v-if="showResults && filteredPersons.length > 0" class="search-results">
      <div
        v-for="(person, index) in filteredPersons"
        :key="person.member_id"
        :class="['search-result-item', { highlighted: index === highlightedIndex }]"
        @click="selectPerson(person)"
        @mouseenter="highlightedIndex = index"
      >
        <div class="person-info">
          <span class="person-name">{{ person.name }}</span>
          <span v-if="person.style_name" class="person-style-name">({{ person.style_name }})</span>
        </div>
        <div class="person-meta">
          <span class="person-gender">{{ person.gender === 'M' ? '男' : '女' }}</span>
          <span class="person-generation">第{{ person.generation }}代</span>
        </div>
      </div>
    </div>

    <div v-if="showResults && searchQuery.trim() && filteredPersons.length === 0" class="no-results">
      未找到匹配的人物
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Person } from '@/types'

const props = defineProps<{
  persons: Person[]
  placeholder?: string
  maxResults?: number
}>()

const emit = defineEmits<{
  (e: 'select', person: Person): void
}>()

const searchQuery = ref('')
const showResults = ref(false)
const highlightedIndex = ref(-1)

const maxResults = props.maxResults ?? 10

const filteredPersons = computed(() => {
  if (!searchQuery.value.trim()) {
    return []
  }
  const query = searchQuery.value.trim().toLowerCase()
  return props.persons
    .filter(p => {
      const nameMatch = p.name.toLowerCase().includes(query)
      const styleNameMatch = p.style_name?.toLowerCase().includes(query)
      return nameMatch || styleNameMatch
    })
    .slice(0, maxResults)
})

function handleInput() {
  highlightedIndex.value = -1
  showResults.value = true
}

function navigateResults(direction: number) {
  if (filteredPersons.value.length === 0) return

  if (highlightedIndex.value === -1) {
    highlightedIndex.value = direction === 1 ? 0 : filteredPersons.value.length - 1
  } else {
    highlightedIndex.value = (highlightedIndex.value + direction + filteredPersons.value.length) % filteredPersons.value.length
  }
}

function selectHighlighted() {
  if (highlightedIndex.value >= 0 && highlightedIndex.value < filteredPersons.value.length) {
    selectPerson(filteredPersons.value[highlightedIndex.value])
  } else if (filteredPersons.value.length > 0) {
    selectPerson(filteredPersons.value[0])
  }
}

function selectPerson(person: Person) {
  emit('select', person)
  showResults.value = false
  searchQuery.value = ''
  highlightedIndex.value = -1
}

function performSearch() {
  if (filteredPersons.value.length > 0) {
    selectPerson(filteredPersons.value[0])
  }
}
</script>

<style scoped>
.search-panel {
  position: relative;
  width: 280px;
}

.search-input-wrapper {
  display: flex;
  gap: 8px;
}

.search-input {
  flex: 1;
  padding: 10px 14px;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  font-size: 14px;
  transition: border-color 0.2s;
}

.search-input:focus {
  outline: none;
  border-color: var(--color-primary, #1890ff);
}

.search-btn {
  padding: 10px 14px;
  border: none;
  border-radius: 8px;
  background: var(--color-primary, #1890ff);
  color: white;
  cursor: pointer;
  transition: background-color 0.2s;
}

.search-btn:hover:not(:disabled) {
  background: var(--color-primary-dark, #096dd9);
}

.search-btn:disabled {
  background: #d9d9d9;
  cursor: not-allowed;
}

.search-results {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  margin-top: 4px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  max-height: 320px;
  overflow-y: auto;
  z-index: 1000;
}

.search-result-item {
  padding: 12px 14px;
  cursor: pointer;
  border-bottom: 1px solid #f0f0f0;
  transition: background-color 0.15s;
}

.search-result-item:last-child {
  border-bottom: none;
}

.search-result-item:hover,
.search-result-item.highlighted {
  background-color: #f5f5f5;
}

.person-info {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 4px;
}

.person-name {
  font-weight: 600;
  color: #333;
}

.person-style-name {
  font-size: 12px;
  color: #888;
}

.person-meta {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: #666;
}

.no-results {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  margin-top: 4px;
  padding: 16px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  text-align: center;
  color: #888;
  font-size: 14px;
}
</style>
