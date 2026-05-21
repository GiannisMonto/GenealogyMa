<template>
  <view class="page">
    <view class="search-bar">
      <input
        v-model="searchKeyword"
        class="search-input"
        placeholder="搜索人物..."
        @confirm="onSearch"
      />
    </view>

    <view v-if="loading" class="loading">加载中...</view>
    <template v-else-if="treeData">
      <view class="tree-header">
        <text class="center-name">{{ treeData.centerPerson.name }}</text>
        <text class="tree-info">
          上查{{ treeData.ancestors.length }}代 · 下查{{ treeData.descendants.length }}代
        </text>
      </view>

      <!-- 简化的树形展示 -->
      <scroll-view class="tree-scroll" scroll-x>
        <view class="tree-container">
          <!-- 祖先区域 -->
          <view v-if="treeData.ancestors.length > 0" class="ancestors-section">
            <view class="section-label">祖先</view>
            <view class="ancestors-list">
              <view
                v-for="ancestor in treeData.ancestors"
                :key="ancestor.id"
                :class="['tree-node', ancestor.gender]"
                @click="onNodeClick(ancestor)"
              >
                <text class="node-name">{{ ancestor.name }}</text>
                <text class="node-gender">{{ ancestor.gender === '男' ? '♂' : '♀' }}</text>
              </view>
            </view>
          </view>

          <!-- 中心人物 -->
          <view class="center-node" @click="onNodeClick(treeData.centerPerson)">
            <view :class="['node-main', treeData.centerPerson.gender]">
              <text class="node-name">{{ treeData.centerPerson.name }}</text>
              <view class="node-icons">
                <text class="gender-icon">{{ treeData.centerPerson.gender === '男' ? '♂' : '♀' }}</text>
                <view :class="['alive-dot', treeData.centerPerson.isAlive ? 'alive' : 'deceased']" />
              </view>
            </view>
          </view>

          <!-- 后代区域 -->
          <view v-if="treeData.descendants.length > 0" class="descendants-section">
            <view class="section-label">后代 ({{ treeData.descendants.length }})</view>
            <view class="descendants-grid">
              <view
                v-for="desc in treeData.descendants"
                :key="desc.id"
                :class="['tree-node', desc.gender]"
                @click="onNodeClick(desc)"
              >
                <text class="node-name">{{ desc.name }}</text>
                <text class="node-gender">{{ desc.gender === '男' ? '♂' : '♀' }}</text>
              </view>
            </view>
          </view>
        </view>
      </scroll-view>
    </template>
    <view v-else class="empty">
      <text>暂无族谱数据</text>
    </view>

    <!-- 底部操作 -->
    <view class="actions">
      <button class="action-btn" @click="onExpand">展开全部</button>
      <button class="action-btn primary" @click="onShare">分享族谱</button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { getFamilyTree } from '../../api/person';
import type { FamilyTreeResponse, FamilyTreeNode } from '../../types/person';

const loading = ref(false);
const searchKeyword = ref('');
const treeData = ref<FamilyTreeResponse | null>(null);

async function loadTree(centerId: number) {
  loading.value = true;
  try {
    treeData.value = await getFamilyTree(centerId, 5, 5);
  } catch {
    uni.showToast({ title: '加载族谱失败', icon: 'none' });
  } finally {
    loading.value = false;
  }
}

function onSearch() {
  if (searchKeyword.value) {
    uni.showToast({ title: '搜索开发中', icon: 'none' });
  }
}

function onNodeClick(node: FamilyTreeNode) {
  uni.navigateTo({ url: `/pages/person/detail?id=${node.id}` });
}

function onExpand() {
  if (treeData.value) {
    loadTree(treeData.value.centerPerson.id);
  }
}

function onShare() {
  uni.showShareMenu({ withShareTicket: true });
}

onMounted(() => {
  const pages = getCurrentPages();
  const currentPage = pages[pages.length - 1] as { options?: { centerId?: string } };
  const centerId = Number(currentPage.options?.centerId || 1);
  loadTree(centerId);
});
</script>

<style scoped>
.page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 120rpx;
}

.search-bar {
  padding: 20rpx 30rpx;
  background-color: #ffffff;
}

.search-input {
  height: 72rpx;
  background-color: #f5f5f5;
  border-radius: 36rpx;
  padding: 0 30rpx;
  font-size: 28rpx;
}

.loading,
.empty {
  text-align: center;
  padding: 100rpx;
  color: #999999;
  font-size: 28rpx;
}

.tree-header {
  background-color: #ffffff;
  padding: 30rpx;
  text-align: center;
  border-bottom: 1rpx solid #f0f0f0;
}

.center-name {
  display: block;
  font-size: 36rpx;
  font-weight: 600;
  color: #333333;
  margin-bottom: 12rpx;
}

.tree-info {
  font-size: 24rpx;
  color: #999999;
}

.tree-scroll {
  padding: 20rpx;
}

.tree-container {
  display: flex;
  flex-direction: column;
  gap: 30rpx;
  min-width: 750rpx;
}

.ancestors-section,
.descendants-section {
  background-color: #ffffff;
  border-radius: 16rpx;
  padding: 24rpx;
}

.section-label {
  font-size: 26rpx;
  color: #999999;
  margin-bottom: 20rpx;
}

.ancestors-list {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.descendants-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150rpx, 1fr));
  gap: 16rpx;
}

.tree-node {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20rpx;
  background-color: #f5f5f5;
  border-radius: 12rpx;
  min-width: 120rpx;
}

.tree-node.男 {
  background-color: #e6f7ff;
}

.tree-node.女 {
  background-color: #fff1f0;
}

.node-name {
  font-size: 28rpx;
  color: #333333;
  font-weight: 500;
}

.node-gender {
  font-size: 22rpx;
  color: #999999;
  margin-top: 4rpx;
}

.center-node {
  display: flex;
  justify-content: center;
  padding: 40rpx;
}

.node-main {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 30rpx 60rpx;
  border-radius: 20rpx;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  box-shadow: 0 8rpx 24rpx rgba(102, 126, 234, 0.3);
}

.node-main.男 {
  background: linear-gradient(135deg, #1890ff 0%, #69c0ff 100%);
  box-shadow: 0 8rpx 24rpx rgba(24, 144, 255, 0.3);
}

.node-main.女 {
  background: linear-gradient(135deg, #ff4d4f 0%, #ff7875 100%);
  box-shadow: 0 8rpx 24rpx rgba(255, 77, 79, 0.3);
}

.node-main .node-name {
  font-size: 36rpx;
  font-weight: 600;
  color: #ffffff;
}

.node-icons {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin-top: 8rpx;
}

.gender-icon {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.8);
}

.alive-dot {
  width: 12rpx;
  height: 12rpx;
  border-radius: 50%;
}

.alive-dot.alive {
  background-color: #52c41a;
}

.alive-dot.deceased {
  background-color: rgba(255, 255, 255, 0.5);
}

.actions {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  gap: 20rpx;
  padding: 20rpx 30rpx;
  background-color: #ffffff;
  box-shadow: 0 -2rpx 10rpx rgba(0, 0, 0, 0.05);
}

.action-btn {
  flex: 1;
  height: 88rpx;
  border-radius: 44rpx;
  font-size: 32rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f5f5;
  color: #666666;
}

.action-btn.primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #ffffff;
}
</style>