<template>
  <view class="login-container">
    <view class="logo-section">
      <image class="logo" src="/static/logo.png" mode="aspectFit" />
      <text class="title">族谱数字化管理平台</text>
      <text class="subtitle">传承家族文化 · 铭记先人恩德</text>
    </view>

    <view class="form-section">
      <view class="form-item">
        <input
          v-model="form.username"
          type="text"
          placeholder="请输入用户名"
          class="input"
        />
      </view>
      <view class="form-item">
        <input
          v-model="form.password"
          type="password"
          placeholder="请输入密码"
          class="input"
        />
      </view>

      <button class="login-btn" :loading="loading" @click="handleLogin">登录</button>

      <view class="extra-links">
        <text class="link" @click="onForgotPassword">忘记密码？</text>
      </view>
    </view>

    <view class="wechat-login">
      <button class="wechat-btn" open-type="getPhoneNumber" @getphonenumber="onWechatLogin">
        <text class="wechat-icon">微</text>
        <text>微信授权登录</text>
      </button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import { useUserStore } from '../../store/user';

const userStore = useUserStore();
const loading = ref(false);

const form = reactive({
  username: '',
  password: '',
});

async function handleLogin() {
  if (!form.username || !form.password) {
    uni.showToast({ title: '请输入用户名和密码', icon: 'none' });
    return;
  }

  loading.value = true;
  try {
    await userStore.login(form.username, form.password);
    uni.showToast({ title: '登录成功', icon: 'success' });
    setTimeout(() => {
      uni.switchTab({ url: '/pages/person/list' });
    }, 1500);
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : '登录失败';
    uni.showToast({ title: msg, icon: 'none' });
  } finally {
    loading.value = false;
  }
}

function onForgotPassword() {
  uni.showToast({ title: '请联系管理员重置密码', icon: 'none' });
}

function onWechatLogin(e: { detail?: { code?: string } }) {
  if (e.detail?.code) {
    uni.showToast({ title: '微信登录开发中', icon: 'none' });
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 120rpx 60rpx 60rpx;
  display: flex;
  flex-direction: column;
}

.logo-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 80rpx;
}

.logo {
  width: 160rpx;
  height: 160rpx;
  border-radius: 32rpx;
  background-color: rgba(255, 255, 255, 0.2);
  margin-bottom: 30rpx;
}

.title {
  font-size: 44rpx;
  font-weight: 600;
  color: #ffffff;
  margin-bottom: 16rpx;
}

.subtitle {
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.8);
}

.form-section {
  background-color: #ffffff;
  border-radius: 24rpx;
  padding: 60rpx 40rpx;
  margin-bottom: 40rpx;
}

.form-item {
  margin-bottom: 32rpx;
}

.input {
  height: 88rpx;
  background-color: #f5f5f5;
  border-radius: 16rpx;
  padding: 0 30rpx;
  font-size: 30rpx;
}

.login-btn {
  width: 100%;
  height: 88rpx;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 44rpx;
  color: #ffffff;
  font-size: 32rpx;
  font-weight: 500;
  display: flex;
  align-items: center;
  justify-content: center;
}

.extra-links {
  display: flex;
  justify-content: center;
  margin-top: 32rpx;
}

.link {
  color: #667eea;
  font-size: 28rpx;
}

.wechat-login {
  flex: 1;
  display: flex;
  align-items: flex-end;
  justify-content: center;
  padding-bottom: 60rpx;
}

.wechat-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 88rpx;
  background-color: #07c160;
  border-radius: 44rpx;
  color: #ffffff;
  font-size: 30rpx;
}

.wechat-icon {
  width: 44rpx;
  height: 44rpx;
  background-color: #ffffff;
  border-radius: 50%;
  color: #07c160;
  font-size: 24rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16rpx;
  font-weight: bold;
}
</style>