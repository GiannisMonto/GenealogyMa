/**
 * 用户状态管理
 */
import { reactive, computed } from 'vue';
import { login as loginApi, logout as logoutApi, getCurrentUser } from '../api/auth';

interface UserState {
  token: string | null;
  userInfo: {
    id: number;
    username: string;
    name: string;
    role: string;
    avatar?: string;
  } | null;
  isLoggedIn: boolean;
}

const state = reactive<UserState>({
  token: uni.getStorageSync('token') || null,
  userInfo: uni.getStorageSync('userInfo') || null,
  isLoggedIn: !!uni.getStorageSync('token'),
});

export function useUserStore() {
  const isLoggedIn = computed(() => state.isLoggedIn);
  const userInfo = computed(() => state.userInfo);

  async function login(username: string, password: string) {
    const res = await loginApi({ username, password });
    state.token = res.token;
    state.userInfo = res.user;
    state.isLoggedIn = true;
    uni.setStorageSync('token', res.token);
    uni.setStorageSync('userInfo', res.user);
    return res;
  }

  async function logout() {
    try {
      await logoutApi();
    } catch {
      // ignore
    }
    state.token = null;
    state.userInfo = null;
    state.isLoggedIn = false;
    uni.removeStorageSync('token');
    uni.removeStorageSync('userInfo');
  }

  async function fetchUserInfo() {
    if (!state.token) return;
    try {
      const user = await getCurrentUser();
      state.userInfo = user;
      uni.setStorageSync('userInfo', user);
    } catch {
      await logout();
    }
  }

  return {
    state,
    isLoggedIn,
    userInfo,
    login,
    logout,
    fetchUserInfo,
  };
}