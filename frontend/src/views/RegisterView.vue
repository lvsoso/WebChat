<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
          注册Web Chat
        </h2>
      </div>
      <form class="mt-8 space-y-6" @submit.prevent="handleRegister">
        <div class="rounded-md shadow-sm -space-y-px">
          <div>
            <label for="email-address" class="sr-only">邮箱地址</label>
            <input
              v-model="email"
              id="email-address"
              name="email"
              type="email"
              required
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
              placeholder="邮箱地址"
            />
          </div>
          <div>
            <label for="password" class="sr-only">密码</label>
            <input
              v-model="password"
              id="password"
              name="password"
              type="password"
              required
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
              placeholder="密码（至少6位）"
            />
          </div>
          <div>
            <label for="confirm-password" class="sr-only">确认密码</label>
            <input
              v-model="confirmPassword"
              id="confirm-password"
              name="confirm-password"
              type="password"
              required
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
              placeholder="确认密码"
            />
          </div>
        </div>

        <div class="flex items-center justify-between">
          <div class="text-sm">
            <router-link to="/login" class="font-medium text-indigo-600 hover:text-indigo-500">
              已有账号？去登录
            </router-link>
          </div>
        </div>

        <div>
          <button
            type="submit"
            :disabled="isLoading || password !== confirmPassword"
            class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:bg-indigo-400"
          >
            {{ isLoading ? '注册中...' : '注册' }}
          </button>
        </div>
        <div v-if="password !== confirmPassword && confirmPassword" class="text-red-500 text-sm text-center">
          两次输入的密码不一致
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';

const router = useRouter();
const authStore = useAuthStore();

const email = ref('');
const password = ref('');
const confirmPassword = ref('');
const isLoading = ref(false);

const handleRegister = async () => {
  if (isLoading.value || password.value !== confirmPassword.value) return;
  
  isLoading.value = true;
  try {
    await authStore.register({
      email: email.value,
      password: password.value
    });
    router.push('/login');
  } catch (error) {
    console.error('注册失败:', error);
  } finally {
    isLoading.value = false;
  }
};
</script>