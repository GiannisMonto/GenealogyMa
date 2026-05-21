<template>
  <div class="home-page">
    <header class="home-header">
      <h1>族谱可视化</h1>
      <p class="subtitle">GenealogyMa - 数字化族谱管理平台</p>
    </header>

    <div class="home-content">
      <nav class="home-nav">
        <router-link to="/genealogy" class="nav-card">
          <span class="nav-icon">🌳</span>
          <span class="nav-title">族谱树</span>
          <span class="nav-desc">交互式族谱树可视化</span>
        </router-link>

        <div class="nav-card" @click="goToStatistics">
          <span class="nav-icon">📊</span>
          <span class="nav-title">世代分布</span>
          <span class="nav-desc">各世代人数统计</span>
        </div>

        <div class="nav-card" @click="goToMigration">
          <span class="nav-icon">🗺️</span>
          <span class="nav-title">迁徙地图</span>
          <span class="nav-desc">族人迁徙路径</span>
        </div>

        <div class="nav-card" @click="goToHall">
          <span class="nav-icon">🏛️</span>
          <span class="nav-title">3D 宗祠</span>
          <span class="nav-desc">沉浸式3D体验</span>
        </div>
      </nav>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

function goToStatistics() {
  router.push('/statistics')
}

function goToMigration() {
  router.push('/migration')
}

function goToHall() {
  router.push('/hall')
}
</script>

<style scoped>
.home-page {
  width: 100%;
  height: 100vh;
  background: linear-gradient(135deg, var(--color-background) 0%, #fff 100%);
}

.home-header {
  text-align: center;
  padding: 60px 20px 40px;
  background: var(--color-primary);
  color: white;
}

.home-header h1 {
  font-size: 2.5rem;
  margin-bottom: 10px;
}

.subtitle {
  font-size: 1.1rem;
  opacity: 0.9;
}

.home-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 20px;
}

.home-nav {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 24px;
}

.nav-card {
  background: white;
  border-radius: 12px;
  padding: 32px 24px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  text-decoration: none;
  color: inherit;
  display: block;
}

.nav-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.nav-icon {
  font-size: 3rem;
  display: block;
  margin-bottom: 16px;
}

.nav-title {
  font-size: 1.3rem;
  font-weight: 600;
  color: var(--color-primary);
  display: block;
  margin-bottom: 8px;
}

.nav-desc {
  font-size: 0.9rem;
  color: #666;
}
</style>
