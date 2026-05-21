<script setup lang="ts">
import { ref } from 'vue';

const hasLogin = ref(false);
const userInfo = ref<{ id: number; name: string; avatar?: string } | null>(null);

// #ifdef MP-WEIXIN
uni.getStorage({
  key: 'token',
  success: (res) => {
    if (res.data) {
      hasLogin.value = true;
      uni.getStorage({
        key: 'userInfo',
        success: (info) => {
          userInfo.value = info.data;
        },
      });
    }
  },
});
// #endif
</script>

<script lang="ts">
export default {
  onLaunch() {
    console.log('App Launch');
  },
  onShow() {
    console.log('App Show');
  },
  onHide() {
    console.log('App Hide');
  },
};
</script>

<style>
page {
  background-color: #f5f5f5;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}

.container {
  padding: 20rpx;
}

.card {
  background-color: #ffffff;
  border-radius: 16rpx;
  padding: 30rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.05);
}
</style>