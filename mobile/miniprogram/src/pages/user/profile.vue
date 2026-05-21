<template>
  <view class="page">
    <!-- 用户信息卡片 -->
    <view class="profile-card">
      <view class="avatar-section">
        <view class="avatar">
          <text class="avatar-text">{{ userStore.userInfo?.name?.slice(0, 1) || '游' }}</text>
        </view>
        <view class="user-info">
          <text class="username">{{ userStore.userInfo?.name || '游客' }}</text>
          <text class="role-badge">{{ userStore.userInfo?.role || 'guest' }}</text>
        </view>
      </view>
    </view>

    <!-- 菜单列表 -->
    <view class="menu-section">
      <view class="menu-group">
        <view class="menu-item" @click="onMenuClick('favorites')">
          <text class="menu-icon">★</text>
          <text class="menu-text">我的收藏</text>
          <text class="arrow">›</text>
        </view>
        <view class="menu-item" @click="onMenuClick('settings')">
          <text class="menu-icon">⚙</text>
          <text class="menu-text">设置</text>
          <text class="arrow">›</text>
        </view>
      </view>

      <view class="menu-group">
        <view class="menu-item" @click="onMenuClick('about')">
          <text class="menu-icon">ℹ</text>
          <text class="menu-text">关于族谱</text>
          <text class="arrow">›</text>
        </view>
        <view class="menu-item" @click="onMenuClick('feedback')">
          <text class="menu-icon">✉</text>
          <text class="menu-text">意见反馈</text>
          <text class="arrow">›</text>
        </view>
      </view>
    </view>

    <!-- 退出登录 -->
    <view class="logout-section">
      <button v-if="userStore.isLoggedIn" class="logout-btn" @click="onLogout">
        退出登录
      </button>
      <button v-else class="login-btn" @click="goLogin">
        立即登录
      </button>
    </view>

    <!-- 版本信息 -->
    <view class="version-info">
      <text>族谱数字化管理平台 v1.0.0</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { useUserStore } from '../../store/user';

const userStore = useUserStore();

function onMenuClick(path: string) {
  uni.showToast({ title: `${path} 开发中`, icon: 'none' });
}

async function onLogout() {
  await userStore.logout();
  uni.showToast({ title: '已退出', icon: 'success' });
  setTimeout(() => {
    uni.redirectTo({ url: '/pages/login/index' });
  }, 1500);
}

function goLogin() {
  uni.navigateTo({ url: '/pages/login/index' });
}
</script>

<style scoped>
.page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 20rpx 30rpx;
}

.profile-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 24rpx;
  padding: 50rpx 40rpx;
  margin-bottom: 30rpx;
}

.avatar-section {
  display: flex;
  align-items: center;
}

.avatar {
  width: 140rpx;
  height: 140rpx;
  background-color: rgba(255, 255, 255, 0.2);
  border-radius: 70rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 30rpx;
}

.avatar-text {
  font-size: 56rpx;
  font-weight: 600;
  color: #ffffff;
}

.user-info {
  flex: 1;
}

.username {
  display: block;
  font-size: 40rpx;
  font-weight: 600;
  color: #ffffff;
  margin-bottom: 12rpx;
}

.role-badge {
  display: inline-block;
  padding: 8rpx 20rpx;
  background-color: rgba(255, 255, 255, 0.2);
  border-radius: 20rpx;
  font-size: 24rpx;
  color: #ffffff;
}

.menu-section {
  margin-bottom: 30rpx;
}

.menu-group {
  background-color: #ffffff;
  border-radius: 16rpx;
  margin-bottom: 20rpx;
  overflow: hidden;
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 32rpx 30rpx;
  border-bottom: 1rpx solid #f5f5f5;
}

.menu-item:last-child {
  border-bottom: none;
}

.menu-icon {
  font-size: 36rpx;
  margin-right: 20rpx;
  color: #667eea;
}

.menu-text {
  flex: 1;
  font-size: 30rpx;
  color: #333333;
}

.arrow {
  font-size: 36rpx;
  color: #cccccc;
}

.logout-section {
  margin-bottom: 30rpx;
}

.logout-btn,
.login-btn {
  width: 100%;
  height: 88rpx;
  border-radius: 44rpx;
  font-size: 32rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logout-btn {
  background-color: #ffffff;
  color: #ff4d4f;
}

.login-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #ffffff;
}

.version-info {
  text-align: center;
  padding: 30rpx;
}

.version-info text {
  font-size: 24rpx;
  color: #999999;
}
</style>