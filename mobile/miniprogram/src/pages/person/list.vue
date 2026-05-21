<template>
  <view class="page">
    <!-- 搜索栏 -->
    <view class="search-bar">
      <input
        v-model="keyword"
        class="search-input"
        placeholder="搜索姓名..."
        @confirm="onSearch"
      />
      <button class="search-btn" @click="onSearch">搜索</button>
    </view>

    <!-- 筛选标签 -->
    <scroll-view class="filter-tags" scroll-x>
      <view class="tag-list">
        <view
          v-for="item in genderOptions"
          :key="item.value"
          :class="['tag', { active: selectedGender === item.value }]"
          @click="onGenderFilter(item.value)"
        >
          {{ item.label }}
        </view>
      </view>
    </scroll-view>

    <!-- 统计卡片 -->
    <view class="stats-cards">
      <view class="stat-card">
        <text class="stat-value">{{ statistics.totalMembers }}</text>
        <text class="stat-label">总人数</text>
      </view>
      <view class="stat-card">
        <text class="stat-value">{{ statistics.aliveCount }}</text>
        <text class="stat-label">在世</text>
      </view>
      <view class="stat-card">
        <text class="stat-value">{{ statistics.deceasedCount }}</text>
        <text class="stat-label">已故</text>
      </view>
      <view class="stat-card">
        <text class="stat-value">{{ statistics.totalGenerations }}</text>
        <text class="stat-label">世代</text>
      </view>
    </view>

    <!-- 人物列表 -->
    <scroll-view class="person-list" scroll-y @scrolltolower="onLoadMore">
      <view
        v-for="person in persons"
        :key="person.id"
        class="person-card"
        @click="goDetail(person.id)"
      >
        <view class="person-avatar">
          <text class="avatar-text">{{ person.name.slice(0, 1) }}</text>
          <view :class="['gender-icon', person.gender === '男' ? 'male' : 'female']">
            {{ person.gender === '男' ? '♂' : '♀' }}
          </view>
        </view>
        <view class="person-info">
          <view class="person-header">
            <text class="person-name">{{ person.name }}</text>
            <text v-if="person.styleName" class="style-name">{{ person.styleName }}</text>
          </view>
          <view class="person-meta">
            <text v-if="person.generation">第{{ person.generation }}代</text>
            <text v-if="person.birthDate">· {{ person.birthDate }}</text>
            <text v-if="person.birthPlace">· {{ person.birthPlace }}</text>
          </view>
          <view class="person-status">
            <view :class="['status-dot', person.isAlive ? 'alive' : 'deceased']" />
            <text class="status-text">{{ person.isAlive ? '在世' : '已故' }}</text>
          </view>
        </view>
        <view class="arrow">›</view>
      </view>

      <!-- 加载更多 -->
      <view v-if="loading" class="loading">加载中...</view>
      <view v-else-if="!hasMore && persons.length > 0" class="no-more">没有更多了</view>
      <view v-else-if="persons.length === 0" class="empty">暂无数据</view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue';
import { searchPersons, getStatistics } from '../../api/person';
import type { PersonDTO } from '../../types/person';

const keyword = ref('');
const selectedGender = ref<string | null>(null);
const loading = ref(false);
const hasMore = ref(true);
const page = ref(1);
const pageSize = 20;

const persons = ref<PersonDTO[]>([]);
const statistics = reactive({
  totalMembers: 0,
  totalGenerations: 0,
  malesCount: 0,
  femalesCount: 0,
  aliveCount: 0,
  deceasedCount: 0,
});

const genderOptions = [
  { label: '全部', value: null },
  { label: '男', value: '男' },
  { label: '女', value: '女' },
];

async function fetchStatistics() {
  try {
    const res = await getStatistics();
    Object.assign(statistics, res);
  } catch {
    // ignore
  }
}

async function fetchPersons(reset = false) {
  if (loading.value) return;
  if (!reset && !hasMore.value) return;

  loading.value = true;
  try {
    const currentPage = reset ? 1 : page.value;
    const res = await searchPersons({
      keyword: keyword.value || undefined,
      gender: (selectedGender.value as '男' | '女') || undefined,
      page: currentPage,
      pageSize,
    });

    if (reset) {
      persons.value = res.data;
      page.value = 2;
    } else {
      persons.value.push(...res.data);
      page.value++;
    }

    hasMore.value = res.data.length === pageSize;
  } catch {
    uni.showToast({ title: '加载失败', icon: 'none' });
  } finally {
    loading.value = false;
  }
}

function onSearch() {
  hasMore.value = true;
  fetchPersons(true);
}

function onGenderFilter(value: string | null) {
  selectedGender.value = value;
  hasMore.value = true;
  fetchPersons(true);
}

function onLoadMore() {
  if (!loading.value && hasMore.value) {
    fetchPersons(false);
  }
}

function goDetail(id: number) {
  uni.navigateTo({ url: `/pages/person/detail?id=${id}` });
}

onMounted(() => {
  fetchStatistics();
  fetchPersons(true);
});
</script>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #f5f5f5;
}

.search-bar {
  display: flex;
  padding: 20rpx 30rpx;
  background-color: #ffffff;
  gap: 20rpx;
}

.search-input {
  flex: 1;
  height: 72rpx;
  background-color: #f5f5f5;
  border-radius: 36rpx;
  padding: 0 30rpx;
  font-size: 28rpx;
}

.search-btn {
  width: 120rpx;
  height: 72rpx;
  background-color: #667eea;
  color: #ffffff;
  border-radius: 36rpx;
  font-size: 28rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.filter-tags {
  background-color: #ffffff;
  padding: 0 30rpx 20rpx;
  white-space: nowrap;
}

.tag-list {
  display: flex;
  gap: 20rpx;
}

.tag {
  display: inline-block;
  padding: 10rpx 30rpx;
  background-color: #f5f5f5;
  border-radius: 30rpx;
  font-size: 26rpx;
  color: #666666;
}

.tag.active {
  background-color: #667eea;
  color: #ffffff;
}

.stats-cards {
  display: flex;
  padding: 20rpx 30rpx;
  gap: 20rpx;
}

.stat-card {
  flex: 1;
  background-color: #ffffff;
  border-radius: 16rpx;
  padding: 24rpx 0;
  text-align: center;
}

.stat-value {
  display: block;
  font-size: 40rpx;
  font-weight: 600;
  color: #667eea;
}

.stat-label {
  display: block;
  font-size: 24rpx;
  color: #999999;
  margin-top: 8rpx;
}

.person-list {
  flex: 1;
  padding: 20rpx 30rpx;
}

.person-card {
  display: flex;
  align-items: center;
  background-color: #ffffff;
  border-radius: 16rpx;
  padding: 30rpx;
  margin-bottom: 20rpx;
}

.person-avatar {
  position: relative;
  width: 100rpx;
  height: 100rpx;
  background-color: #667eea;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 24rpx;
}

.avatar-text {
  font-size: 40rpx;
  font-weight: 600;
  color: #ffffff;
}

.gender-icon {
  position: absolute;
  bottom: -4rpx;
  right: -4rpx;
  width: 36rpx;
  height: 36rpx;
  border-radius: 50%;
  font-size: 20rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2rpx solid #ffffff;
}

.gender-icon.male {
  background-color: #e6f7ff;
  color: #1890ff;
}

.gender-icon.female {
  background-color: #fff1f0;
  color: #ff4d4f;
}

.person-info {
  flex: 1;
}

.person-header {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.person-name {
  font-size: 32rpx;
  font-weight: 500;
  color: #333333;
}

.style-name {
  font-size: 26rpx;
  color: #999999;
}

.person-meta {
  font-size: 24rpx;
  color: #999999;
  margin-top: 8rpx;
}

.person-status {
  display: flex;
  align-items: center;
  gap: 8rpx;
  margin-top: 8rpx;
}

.status-dot {
  width: 12rpx;
  height: 12rpx;
  border-radius: 50%;
}

.status-dot.alive {
  background-color: #52c41a;
}

.status-dot.deceased {
  background-color: #8c8c8c;
}

.status-text {
  font-size: 24rpx;
  color: #999999;
}

.arrow {
  font-size: 40rpx;
  color: #cccccc;
}

.loading,
.no-more,
.empty {
  text-align: center;
  padding: 40rpx;
  color: #999999;
  font-size: 28rpx;
}
</style>