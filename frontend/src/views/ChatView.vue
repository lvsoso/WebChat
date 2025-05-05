<template>
  <div class="flex flex-col h-screen">
    <!-- 导航栏 -->
    <div class="bg-indigo-600 text-white p-4 flex justify-between items-center">
      <h1 class="text-xl font-bold">Web Chat</h1>
      <button 
        @click="router.push('/conversations')"
        class="px-3 py-1 bg-white text-indigo-600 rounded hover:bg-gray-100"
      >
        查看历史对话
      </button>
    </div>
    <!-- 聊天历史区域 -->
    <div class="flex-1 overflow-y-auto p-4 space-y-4">
      <div v-for="message in messages" :key="message.id" class="flex">
        <div
          :class="[
            'max-w-3xl p-4 rounded-lg',
            message.role === 'user'
              ? 'bg-indigo-100 ml-auto'
              : 'bg-white mr-auto'
          ]"
        >
          <div class="text-sm font-medium text-gray-900">
            {{ message.role === 'user' ? '你' : 'AI' }}
          </div>
          <div class="mt-1 text-sm text-gray-700">{{ message.content }}</div>
        </div>
      </div>
    </div>

    <!-- 输入区域 -->
    <div class="border-t p-4 bg-white">
      <div class="flex space-x-4">
        <select
          v-model="selectedModel"
          class="block rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
        >
          <option value="gpt-4">GPT-4</option>
          <option value="gpt-3.5-turbo">GPT-3.5</option>
          <option value="deepseek">Deepseek</option>
        </select>
        <div class="flex-1">
          <textarea
            v-model="newMessage"
            rows="1"
            class="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
            placeholder="输入消息..."
            @keydown.enter.prevent="sendMessage"
          />
        </div>
        <button
          @click="sendMessage"
          class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700"
          :disabled="isLoading"
        >
          {{ isLoading ? '发送中...' : '发送' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useChatStore } from '../stores/chat';

const router = useRouter();

interface Message {
  id: string;
  role: 'user' | 'assistant';
  content: string;
}

const chatStore = useChatStore();
const messages = ref<Message[]>([]);
const newMessage = ref('');
const selectedModel = ref('deepseek');
const isLoading = ref(false);

const sendMessage = async () => {
  if (!newMessage.value.trim() || isLoading.value) return;

  const userMessage: Message = {
    id: Date.now().toString(),
    role: 'user',
    content: newMessage.value
  };

  messages.value.push(userMessage);
  isLoading.value = true;

  try {
    const response = await chatStore.sendMessage({
      model: selectedModel.value,
      message: newMessage.value
    });

    messages.value.push({
      id: Date.now().toString(),
      role: 'assistant',
      content: response.content
    });
  } catch (error) {
    console.error('发送消息失败:', error);
  } finally {
    isLoading.value = false;
    newMessage.value = '';
  }
};
</script>