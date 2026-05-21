<template>
  <div class="genealogy-page">
    <header class="page-header">
      <h1>族谱树</h1>
      <div class="header-controls">
        <input
          v-model="searchKeyword"
          type="text"
          placeholder="搜索人物..."
          class="search-input"
          @keyup.enter="handleSearch"
        />
        <button @click="handleSearch" class="search-btn">搜索</button>
        <button @click="resetView" class="reset-btn">重置</button>
      </div>
    </header>

    <div class="tree-container" ref="treeContainer">
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="error" class="error">{{ error }}</div>
      <GenealogyTree v-else :data="treeData" @node-click="handleNodeClick" />
    </div>

    <PersonDetail
      v-if="selectedPerson"
      :person="selectedPerson"
      @close="selectedPerson = null"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { usePersonStore } from '@/stores/person'
import GenealogyTree from '@/components/GenealogyTree/index.vue'
import PersonDetail from '@/components/PersonDetail/index.vue'
import type { Person } from '@/types'

const personStore = usePersonStore()

const searchKeyword = ref('')
const selectedPerson = ref<Person | null>(null)
const treeContainer = ref<HTMLElement | null>(null)

const loading = ref(false)
const error = ref<string | null>(null)
const treeData = ref<any>(null)

async function loadTreeData() {
  loading.value = true
  error.value = null
  try {
    await personStore.fetchGenealogyTree()
    treeData.value = transformToTreeData(personStore.persons)
  } catch (e) {
    error.value = e instanceof Error ? e.message : '加载失败'
  } finally {
    loading.value = false
  }
}

function transformToTreeData(persons: Person[]): any {
  if (!persons.length) return null

  const personMap = new Map<number, any>()
  persons.forEach((p) => {
    personMap.set(p.member_id, {
      id: p.member_id,
      name: p.name,
      generation: p.generation,
      gender: p.gender,
      styleName: p.style_name,
      birthYear: p.birth_year,
      deathYear: p.death_year,
      children: []
    })
  })

  let root: any = null
  personMap.forEach((node) => {
    if (node.fatherId && personMap.has(node.fatherId)) {
      personMap.get(node.fatherId).children.push(node)
    } else {
      root = node
    }
  })

  return root
}

function handleSearch() {
  // Search implementation
}

function resetView() {
  loadTreeData()
}

function handleNodeClick(person: Person) {
  selectedPerson.value = person
}

onMounted(() => {
  loadTreeData()
})
</script>

<style scoped>
.genealogy-page {
  width: 100%;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--color-background);
}

.page-header {
  background: var(--color-primary);
  color: white;
  padding: 16px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-shrink: 0;
}

.page-header h1 {
  font-size: 1.5rem;
}

.header-controls {
  display: flex;
  gap: 12px;
}

.search-input {
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  width: 200px;
}

.search-btn,
.reset-btn {
  padding: 8px 20px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}

.search-btn {
  background: var(--color-secondary);
  color: white;
}

.reset-btn {
  background: rgba(255, 255, 255, 0.2);
  color: white;
}

.tree-container {
  flex: 1;
  overflow: hidden;
  position: relative;
}

.loading,
.error {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  font-size: 1.2rem;
}

.error {
  color: #dc3545;
}
</style>
