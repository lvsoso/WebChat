import { defineStore } from 'pinia';
import axios from 'axios';

interface SendMessagePayload {
  model: string;
  message: string;
  conversation_id?: number;
}

interface MessageResponse {
  content: string;
  token_count: number;
}

export const useChatStore = defineStore('chat', {
  state: () => ({
    conversations: [] as any[],
  }),

  actions: {
    async sendMessage(payload: SendMessagePayload): Promise<MessageResponse> {
      try {
        const response = await axios.post('/api/chat/send', payload);
        return response.data;
      } catch (error) {
        console.error('发送消息失败:', error);
        throw error;
      }
    },

    async getConversations() {
      try {
        const response = await axios.get('/api/chat/conversations');
        this.conversations = response.data;
      } catch (error) {
        console.error('获取对话历史失败:', error);
        throw error;
      }
    },

    async deleteConversation(conversationId: string) {
      try {
        await axios.delete(`/api/chat/conversations/${conversationId}`);
        this.conversations = this.conversations.filter(
          (conv) => conv.id !== conversationId
        );
      } catch (error) {
        console.error('删除对话失败:', error);
        throw error;
      }
    },
    
    async getConversationMessages(conversationId: string) {
      try {
        const response = await axios.get(`/api/chat/conversations/${conversationId}/messages`);
        return response.data;
      } catch (error) {
        console.error('获取对话消息失败:', error);
        throw error;
      }
    },

    async getConversation(conversationId: string) {
      try {
        const response = await axios.get(`/api/chat/conversations/${conversationId}`);
        return response.data;
      } catch (error) {
        console.error('获取对话详情失败:', error);
        throw error;
      }
    },
  },
});