<template>
  <div class="flex flex-col h-screen">
    <!-- 导航栏 -->
    <div class="bg-indigo-600 text-white p-4 flex justify-between items-center">
      <button @click="router.push('/conversations')" class="text-white">
        &larr; 返回对话列表
      </button>
      <h1 class="text-xl font-bold">{{ conversation.Title || '对话详情' }}</h1>
      <div></div> <!-- 占位，保持标题居中 -->
    </div>
    
    <!-- 消息列表 -->
    <div class="flex-1 overflow-y-auto p-4 space-y-4">
      <div v-if="loading" class="text-center py-4">
        <p>加载中...</p>
      </div>
      
      <div v-else-if="messages.length === 0" class="text-center py-4">
        <p>暂无消息</p>
      </div>
      
      <div v-else>
        <div v-for="message in messages" :key="message.ID" class="flex">
          <div
            :class="[
              'max-w-3xl p-4 rounded-lg',
              message.Role === 'user'
                ? 'bg-indigo-100 ml-auto'
                : 'bg-white mr-auto border'
            ]"
          >
            <div class="text-sm font-medium text-gray-900">
              {{ message.Role === 'user' ? '你' : 'AI' }}
            </div>
            <div class="mt-1 text-sm text-gray-700">{{ message.Content }}</div>
            <div v-if="message.TokenCount > 0" class="mt-1 text-xs text-gray-500">
              Token数: {{ message.TokenCount }}
            </div>
          </div>
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
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import axios from 'axios';
import { useChatStore } from '../stores/chat';

const route = useRoute();
const router = useRouter();
const chatStore = useChatStore();
const conversationId = route.params.id as string;
const conversation = ref<any>({});
const messages = ref<any[]>([]);
const loading = ref(true);
const newMessage = ref('');
const selectedModel = ref('deepseek');
const isLoading = ref(false);

onMounted(async () => {
  try {
    // 获取对话详情
    const messagesResponse = await axios.get(`/api/chat/conversations/${conversationId}/messages`);
    messages.value = messagesResponse.data;
    
    // 如果有消息，使用第一条消息的会话信息
    if (messages.value.length > 0) {
      const conversationResponse = await axios.get(`/api/chat/conversations/${conversationId}`);
      conversation.value = conversationResponse.data;
      // 设置默认模型为最后一条消息使用的模型
      const lastMessage = messages.value[messages.value.length - 1];
      if (lastMessage && lastMessage.ModelName) {
        selectedModel.value = lastMessage.ModelName;
      }
    }
  } catch (error) {
    console.error('获取对话详情失败:', error);
  } finally {
    loading.value = false;
  }
});

const sendMessage = async () => {
  if (!newMessage.value.trim() || isLoading.value) return;

  // 添加用户消息到界面
  const userMessage = {
    ID: Date.now().toString(),
    Role: 'user',
    Content: newMessage.value,
    ConversationID: conversationId,
    TokenCount: 0,
    CreatedAt: new Date().toISOString()
  };

  messages.value.push(userMessage);
  isLoading.value = true;
  const messageText = newMessage.value;

  try {
    // 发送消息到服务器
    const response = await chatStore.sendMessage({
      model: selectedModel.value,
      message: messageText,
      conversation_id: parseInt(conversationId)
    });

    // 添加AI回复到界面
    const aiMessage = {
      ID: Date.now().toString() + '-ai',
      Role: 'assistant',
      Content: response.content,
      ConversationID: conversationId,
      TokenCount: response.token_count || 0,
      CreatedAt: new Date().toISOString()
    };

    messages.value.push(aiMessage);

    // 刷新对话信息
    const conversationResponse = await axios.get(`/api/chat/conversations/${conversationId}`);
    conversation.value = conversationResponse.data;
  } catch (error) {
    console.error('发送消息失败:', error);
    alert('发送消息失败，请重试');
  } finally {
    isLoading.value = false;
    newMessage.value = '';
  }
};
</script>