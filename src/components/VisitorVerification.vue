<template>
  <div class="verification-container">
    <div class="verification-card">
      <div class="card-header">
        <h2>欢迎访问个人简历</h2>
        <p>请使用以下方式之一进行验证以继续访问</p>
      </div>
      
      <div v-if="error" class="error-message">
        <i class="fas fa-exclamation-circle"></i> {{ error }}
      </div>
      
      <form @submit.prevent="verifyAccess" class="verification-form">
        <div class="tabs">
          <div 
            v-for="tab in tabs" 
            :key="tab.type" 
            class="tab" 
            :class="{ active: activeTab === tab.type }"
            @click="activeTab = tab.type"
          >
            <i :class="tab.icon"></i> {{ tab.label }}
          </div>
        </div>
        
        <div class="input-group">
          <label :for="activeTab">{{ getInputLabel() }}</label>
          <input 
            :id="activeTab" 
            :type="getInputType()" 
            v-model="verificationValue" 
            :placeholder="getInputPlaceholder()"
            required
          >
        </div>
        
        <button type="submit" class="verify-btn" :disabled="verifying">
          <i v-if="verifying" class="fas fa-spinner fa-spin"></i>
          <span v-else>验证</span>
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import axios from 'axios';
import { useRouter } from 'vue-router';
import { API_URL } from '../config';

const router = useRouter();

// 验证状态
const verifying = ref(false);
const error = ref('');
const verificationValue = ref('');

// 标签页配置
const tabs = [
  { type: 'password', label: '密码验证', icon: 'fas fa-key' },
  { type: 'name', label: '姓名验证', icon: 'fas fa-user' },
  { type: 'email', label: '邮箱验证', icon: 'fas fa-envelope' },
  { type: 'phone', label: '电话验证', icon: 'fas fa-phone' }
];

// 当前选中标签页
const activeTab = ref('password');

// 获取输入框标签
const getInputLabel = () => {
  switch (activeTab.value) {
    case 'password': return '访问密码';
    case 'name': return '求职者姓名';
    case 'email': return '求职者邮箱';
    case 'phone': return '求职者电话';
    default: return '';
  }
};

// 获取输入框类型
const getInputType = () => {
  switch (activeTab.value) {
    case 'password': return 'password';
    case 'email': return 'email';
    case 'phone': return 'tel';
    default: return 'text';
  }
};

// 获取输入框提示文本
const getInputPlaceholder = () => {
  switch (activeTab.value) {
    case 'password': return '请输入访问密码';
    case 'name': return '请输入求职者的完整姓名';
    case 'email': return '请输入求职者的邮箱地址';
    case 'phone': return '请输入求职者的电话号码';
    default: return '';
  }
};

// 验证访问
const verifyAccess = async () => {
  verifying.value = true;
  error.value = '';
  
  try {
    const response = await axios.post(`${API_URL}/verify`, {
      verification_type: activeTab.value,
      value: verificationValue.value
    });
    
    if (response.data.success) {
      // 存储访客令牌
      localStorage.setItem('visitorToken', response.data.data.token);
      
      // 跳转到简历页面
      router.push('/');
    } else {
      error.value = response.data.message || '验证失败';
    }
  } catch (err) {
    console.error('验证失败:', err);
    if (err.response && err.response.status === 401) {
      error.value = '验证信息不正确，请重试';
    } else {
      error.value = '验证过程中发生错误，请稍后再试';
    }
  } finally {
    verifying.value = false;
  }
};
</script>

<style scoped>
.verification-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f0f2f5;
  padding: 20px;
}

.verification-card {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 450px;
  overflow: hidden;
}

.card-header {
  padding: 20px;
  background-color: var(--primary-color);
  color: white;
  text-align: center;
}

.card-header h2 {
  margin: 0 0 10px 0;
  font-size: 1.5rem;
}

.card-header p {
  margin: 0;
  opacity: 0.9;
}

.verification-form {
  padding: 20px;
}

.error-message {
  background-color: #fde8e8;
  color: #e02424;
  padding: 12px 15px;
  margin-bottom: 15px;
  border-radius: 5px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.tabs {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
  border-bottom: 1px solid #e2e8f0;
}

.tab {
  padding: 12px;
  text-align: center;
  cursor: pointer;
  color: #4b5563;
  font-weight: 500;
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 5px;
  font-size: 0.9rem;
  transition: all 0.2s;
}

.tab i {
  font-size: 1.1rem;
}

.tab.active {
  color: var(--primary-color);
  border-bottom: 2px solid var(--primary-color);
}

.tab:hover:not(.active) {
  background-color: #f8fafc;
}

.input-group {
  margin-bottom: 20px;
}

.input-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: #374151;
}

.input-group input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #d1d5db;
  border-radius: 5px;
  font-size: 1rem;
}

.input-group input:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.3);
}

.verify-btn {
  width: 100%;
  padding: 12px;
  background-color: var(--primary-color);
  color: white;
  border: none;
  border-radius: 5px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
}

.verify-btn:hover:not(:disabled) {
  background-color: var(--primary-dark);
}

.verify-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

@media (max-width: 480px) {
  .tabs {
    flex-wrap: wrap;
  }
  
  .tab {
    min-width: 50%;
    border-bottom: 1px solid #e2e8f0;
  }
  
  .tab.active {
    border-bottom: 2px solid var(--primary-color);
  }
}
</style> 