<template>
  <div class="p-4">
    <div class="flex justify-between items-center mb-4">
      <h1 class="text-2xl font-bold">历史对话</h1>
      <button 
        @click="createNewConversation"
        class="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
      >
        新建对话
      </button>
    </div>
    
    <div v-if="loading" class="text-center py-4">
      <p>加载中...</p>
    </div>
    
    <div v-else-if="conversations.length === 0" class="text-center py-4">
      <p>暂无历史对话</p>
    </div>
    
    <div v-else class="space-y-2">
      <div 
        v-for="conversation in conversations" 
        :key="conversation.ID"
        class="border rounded-lg p-4 hover:bg-gray-50 flex justify-between items-center"
      >
        <div>
          <h2 class="font-medium">{{ conversation.Title }}</h2>
          <p class="text-sm text-gray-500">{{ formatDate(conversation.CreatedAt) }}</p>
        </div>
        
        <div class="flex space-x-2">
          <button 
            @click="viewConversation(conversation.ID)"
            class="px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            查看
          </button>
          <button 
            @click="deleteConversation(conversation.ID)"
            class="px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600"
          >
            删除
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useChatStore } from '../stores/chat';

const router = useRouter();
const chatStore = useChatStore();
const conversations = ref<any[]>([]);
const loading = ref(true);

onMounted(async () => {
  try {
    await chatStore.getConversations();
    conversations.value = chatStore.conversations;
  } catch (error) {
    console.error('获取对话失败:', error);
  } finally {
    loading.value = false;
  }
});

const viewConversation = (id: string) => {
  router.push(`/conversation/${id}`);
};

const createNewConversation = () => {
  router.push('/');
};

const deleteConversation = async (id: string) => {
  if (confirm('确定要删除这个对话吗？')) {
    try {
      await chatStore.deleteConversation(id);
      conversations.value = chatStore.conversations;
    } catch (error) {
      console.error('删除对话失败:', error);
    }
  }
};

const formatDate = (dateString: string) => {
  const date = new Date(dateString);
  return date.toLocaleString();
};
</script>