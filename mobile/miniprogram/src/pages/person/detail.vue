<template>
  <view class="page">
    <view v-if="loading" class="loading">加载中...</view>
    <template v-else-if="person">
      <!-- 基本信息卡片 -->
      <view class="info-card">
        <view class="person-header">
          <view class="avatar">
            <text class="avatar-text">{{ person.name.slice(0, 1) }}</text>
          </view>
          <view class="basic-info">
            <view class="name-row">
              <text class="name">{{ person.name }}</text>
              <text v-if="person.styleName" class="style-name">{{ person.styleName }}</text>
            </view>
            <view class="gender-badge" :class="person.gender">
              {{ person.gender === '男' ? '♂ 男' : '♀ 女' }}
            </view>
          </view>
          <view :class="['alive-badge', person.isAlive ? 'alive' : 'deceased']">
            {{ person.isAlive ? '在世' : '已故' }}
          </view>
        </view>
      </view>

      <!-- 详细信息 -->
      <view class="detail-card">
        <view class="detail-title">基本信息</view>
        <view class="detail-grid">
          <view class="detail-item">
            <text class="detail-label">出生日期</text>
            <text class="detail-value">{{ person.birthDate || '未知' }}</text>
          </view>
          <view class="detail-item">
            <text class="detail-label">逝世日期</text>
            <text class="detail-value">{{ person.deathDate || (person.isAlive ? '-' : '未知') }}</text>
          </view>
          <view class="detail-item">
            <text class="detail-label">所属世代</text>
            <text class="detail-value">{{ person.generation ? `第 ${person.generation} 代` : '未知' }}</text>
          </view>
          <view class="detail-item">
            <text class="detail-label">出生地</text>
            <text class="detail-value">{{ person.birthPlace || '未知' }}</text>
          </view>
        </view>
      </view>

      <!-- 父亲信息 -->
      <view v-if="father" class="relation-card" @click="goPerson(father.id)">
        <view class="relation-icon father">父</view>
        <view class="relation-info">
          <text class="relation-label">父亲</text>
          <text class="relation-name">{{ father.name }}</text>
        </view>
        <text class="arrow">›</text>
      </view>

      <!-- 配偶信息 -->
      <view v-if="spouses.length > 0" class="section">
        <view class="section-title">配偶</view>
        <view
          v-for="spouse in spouses"
          :key="spouse.id"
          class="relation-card"
          @click="goPerson(spouse.spouseId)"
        >
          <view class="relation-icon spouse">配</view>
          <view class="relation-info">
            <text class="relation-label">{{ spouse.spouseGender === '男' ? '丈夫' : '妻子' }}</text>
            <text class="relation-name">{{ spouse.spouseName }}</text>
          </view>
          <text class="arrow">›</text>
        </view>
      </view>

      <!-- 子女信息 -->
      <view v-if="children.length > 0" class="section">
        <view class="section-title">子女 ({{ children.length }})</view>
        <view
          v-for="child in children"
          :key="child.id"
          class="relation-card"
          @click="goPerson(child.childId)"
        >
          <view class="relation-icon" :class="child.childGender === '男' ? 'male' : 'female'">
            {{ child.childGender === '男' ? '♂' : '♀' }}
          </view>
          <view class="relation-info">
            <text class="relation-name">{{ child.childName }}</text>
            <text class="relation-meta">
              {{ child.childGender }} · {{ child.birthDate || '未知出生日期' }}
            </text>
          </view>
          <text class="arrow">›</text>
        </view>
      </view>

      <!-- 生平介绍 -->
      <view v-if="person.biography" class="biography-card">
        <view class="detail-title">生平简介</view>
        <text class="biography-text">{{ person.biography }}</text>
      </view>

      <!-- 操作按钮 -->
      <view class="actions">
        <button class="action-btn primary" @click="viewTree">查看族谱</button>
      </view>
    </template>
    <view v-else class="error">加载失败</view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { getPerson } from '../../api/person';
import type { PersonDTO, SpouseDTO, ChildDTO } from '../../types/person';

const loading = ref(true);
const person = ref<PersonDTO | null>(null);
const father = ref<{ id: number; name: string } | null>(null);
const spouses = ref<SpouseDTO[]>([]);
const children = ref<ChildDTO[]>([]);

async function loadData() {
  const pages = getCurrentPages();
  const currentPage = pages[pages.length - 1] as { options?: { id?: string } };
  const id = Number(currentPage.options?.id || 0);

  if (!id) {
    uni.showToast({ title: '参数错误', icon: 'none' });
    return;
  }

  loading.value = true;
  try {
    const res = await getPerson(id, true);
    person.value = res;
    // 简化处理，实际应从API获取关联数据
  } catch {
    uni.showToast({ title: '加载失败', icon: 'none' });
  } finally {
    loading.value = false;
  }
}

function goPerson(id: number) {
  uni.navigateTo({ url: `/pages/person/detail?id=${id}` });
}

function viewTree() {
  if (person.value) {
    uni.navigateTo({ url: `/pages/genealogy/tree?centerId=${person.value.id}` });
  }
}

onMounted(loadData);
</script>

<style scoped>
.page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 20rpx 30rpx 120rpx;
}

.loading,
.error {
  text-align: center;
  padding: 100rpx;
  color: #999999;
}

.info-card,
.detail-card,
.biography-card {
  background-color: #ffffff;
  border-radius: 16rpx;
  padding: 30rpx;
  margin-bottom: 20rpx;
}

.person-header {
  display: flex;
  align-items: center;
}

.avatar {
  width: 120rpx;
  height: 120rpx;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 60rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 30rpx;
}

.avatar-text {
  font-size: 48rpx;
  font-weight: 600;
  color: #ffffff;
}

.basic-info {
  flex: 1;
}

.name-row {
  display: flex;
  align-items: center;
  gap: 16rpx;
  margin-bottom: 12rpx;
}

.name {
  font-size: 40rpx;
  font-weight: 600;
  color: #333333;
}

.style-name {
  font-size: 28rpx;
  color: #999999;
}

.gender-badge {
  display: inline-block;
  padding: 6rpx 20rpx;
  border-radius: 20rpx;
  font-size: 24rpx;
}

.gender-badge.男 {
  background-color: #e6f7ff;
  color: #1890ff;
}

.gender-badge.女 {
  background-color: #fff1f0;
  color: #ff4d4f;
}

.alive-badge {
  padding: 8rpx 20rpx;
  border-radius: 20rpx;
  font-size: 24rpx;
}

.alive-badge.alive {
  background-color: #f6ffed;
  color: #52c41a;
}

.alive-badge.deceased {
  background-color: #f5f5f5;
  color: #8c8c8c;
}

.detail-title {
  font-size: 30rpx;
  font-weight: 500;
  color: #333333;
  margin-bottom: 24rpx;
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24rpx;
}

.detail-item {
  display: flex;
  flex-direction: column;
}

.detail-label {
  font-size: 24rpx;
  color: #999999;
  margin-bottom: 8rpx;
}

.detail-value {
  font-size: 28rpx;
  color: #333333;
}

.section {
  margin-bottom: 20rpx;
}

.section-title {
  font-size: 30rpx;
  font-weight: 500;
  color: #333333;
  margin-bottom: 16rpx;
}

.relation-card {
  display: flex;
  align-items: center;
  background-color: #ffffff;
  border-radius: 16rpx;
  padding: 24rpx 30rpx;
  margin-bottom: 16rpx;
}

.relation-icon {
  width: 80rpx;
  height: 80rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32rpx;
  margin-right: 24rpx;
}

.relation-icon.father {
  background-color: #e6f7ff;
  color: #1890ff;
}

.relation-icon.spouse {
  background-color: #fff1f0;
  color: #ff4d4f;
}

.relation-icon.male {
  background-color: #e6f7ff;
  color: #1890ff;
}

.relation-icon.female {
  background-color: #fff1f0;
  color: #ff4d4f;
}

.relation-info {
  flex: 1;
}

.relation-label {
  font-size: 24rpx;
  color: #999999;
  display: block;
  margin-bottom: 4rpx;
}

.relation-name {
  font-size: 30rpx;
  color: #333333;
  font-weight: 500;
}

.relation-meta {
  font-size: 24rpx;
  color: #999999;
}

.arrow {
  font-size: 40rpx;
  color: #cccccc;
}

.biography-text {
  font-size: 28rpx;
  color: #666666;
  line-height: 1.6;
}

.actions {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 20rpx 30rpx;
  background-color: #ffffff;
  box-shadow: 0 -2rpx 10rpx rgba(0, 0, 0, 0.05);
}

.action-btn {
  height: 88rpx;
  border-radius: 44rpx;
  font-size: 32rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.action-btn.primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #ffffff;
}
</style>