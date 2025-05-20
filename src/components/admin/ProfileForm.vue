<template>
  <div class="profile-form-container">
    <div v-if="loading" class="loading">
      <i class="fas fa-spinner fa-spin"></i> 加载中...
    </div>
    <div v-else-if="error" class="error-message">
      <i class="fas fa-exclamation-circle"></i> {{ error }}
    </div>
    <form v-else @submit.prevent="saveProfile" class="profile-form">
      <div class="form-group">
        <label for="name">姓名</label>
        <input type="text" id="name" v-model="profile.name" required>
      </div>
      
      <div class="form-group">
        <label for="title">职位</label>
        <input type="text" id="title" v-model="profile.title" required>
      </div>
      
      <div class="form-group">
        <label for="avatar">头像URL</label>
        <input type="text" id="avatar" v-model="profile.avatar">
      </div>
      
      <div class="form-group">
        <label for="email">电子邮箱</label>
        <input type="email" id="email" v-model="profile.email">
      </div>
      
      <div class="form-group">
        <label for="phone">联系电话</label>
        <input type="text" id="phone" v-model="profile.phone">
      </div>
      
      <div class="form-group">
        <label for="location">所在地</label>
        <input type="text" id="location" v-model="profile.location">
      </div>
      
      <div class="form-group">
        <label for="years-exp">工作年限</label>
        <input type="number" id="years-exp" v-model.number="profile.yearsOfExp" min="0">
      </div>
      
      <div class="form-group">
        <label for="education">教育背景</label>
        <input type="text" id="education" v-model="profile.education">
      </div>
      
      <div class="form-group">
        <label for="job-status">求职状态</label>
        <select id="job-status" v-model="profile.jobStatus">
          <option value="可入职">可入职</option>
          <option value="在职-考虑机会">在职-考虑机会</option>
          <option value="暂不找工作">暂不找工作</option>
        </select>
      </div>
      
      <div class="form-group full-width">
        <label for="introduction">个人简介</label>
        <textarea id="introduction" v-model="profile.introduction" rows="5"></textarea>
      </div>
      
      <div class="form-group full-width">
        <label for="philosophy">个人理念</label>
        <textarea id="philosophy" v-model="profile.philosophy" rows="3"></textarea>
      </div>
      
      <div class="form-group">
        <label for="resume-url">简历文件URL</label>
        <input type="text" id="resume-url" v-model="profile.resumeFileURL">
      </div>
      
      <div class="form-actions">
        <button type="submit" class="save-btn" :disabled="saving">
          <i class="fas fa-save"></i> {{ saving ? '保存中...' : '保存更改' }}
        </button>
      </div>
    </form>
    
    <div v-if="saveSuccess" class="success-message">
      <i class="fas fa-check-circle"></i> 保存成功
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue';
import axios from 'axios';

const API_URL = 'http://localhost:8080/api';
const loading = ref(false);
const saving = ref(false);
const error = ref(null);
const saveSuccess = ref(false);

// 初始化空的个人信息对象
const profile = reactive({
  id: null,
  name: '',
  title: '',
  avatar: '',
  email: '',
  phone: '',
  location: '',
  introduction: '',
  yearsOfExp: 0,
  education: '',
  jobStatus: '可入职',
  philosophy: '',
  resumeFileURL: ''
});

// 获取个人信息
const fetchProfile = async () => {
  loading.value = true;
  error.value = null;
  
  try {
    const token = localStorage.getItem('token');
    const response = await axios.get(`${API_URL}/admin/profile`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    if (response.data.success) {
      // 更新个人信息
      const data = response.data.data;
      Object.keys(profile).forEach(key => {
        if (data[key] !== undefined) {
          profile[key] = data[key];
        }
      });
    } else {
      error.value = response.data.message || '获取个人信息失败';
    }
  } catch (err) {
    console.error('获取个人信息出错:', err);
    error.value = '获取个人信息时发生错误，请稍后再试';
  } finally {
    loading.value = false;
  }
};

// 保存个人信息
const saveProfile = async () => {
  saving.value = true;
  error.value = null;
  saveSuccess.value = false;
  
  try {
    const token = localStorage.getItem('token');
    const response = await axios.put(`${API_URL}/admin/profile`, profile, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });
    
    if (response.data.success) {
      saveSuccess.value = true;
      // 3秒后隐藏成功消息
      setTimeout(() => {
        saveSuccess.value = false;
      }, 3000);
    } else {
      error.value = response.data.message || '保存个人信息失败';
    }
  } catch (err) {
    console.error('保存个人信息出错:', err);
    error.value = '保存个人信息时发生错误，请稍后再试';
  } finally {
    saving.value = false;
  }
};

// 页面加载时获取个人信息
onMounted(() => {
  fetchProfile();
});
</script>

<style scoped>
.profile-form-container {
  padding: 20px 0;
  position: relative;
}

.loading, .error-message, .success-message {
  padding: 15px;
  margin-bottom: 20px;
  border-radius: 5px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.loading {
  background-color: #e9f0fd;
  color: #1a56db;
}

.error-message {
  background-color: #fde8e8;
  color: #e02424;
}

.success-message {
  background-color: #def7ec;
  color: #03543e;
  position: fixed;
  bottom: 20px;
  right: 20px;
  z-index: 100;
  padding: 12px 20px;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  animation: fadeInOut 3s ease-in-out;
}

@keyframes fadeInOut {
  0% { opacity: 0; transform: translateY(20px); }
  10% { opacity: 1; transform: translateY(0); }
  90% { opacity: 1; transform: translateY(0); }
  100% { opacity: 0; transform: translateY(-20px); }
}

.profile-form {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
}

.form-group {
  flex: 1 0 calc(50% - 20px);
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group.full-width {
  flex: 1 0 100%;
}

.form-group label {
  font-weight: 600;
  color: #374151;
}

.form-group input,
.form-group select,
.form-group textarea {
  padding: 10px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 1rem;
}

.form-group textarea {
  resize: vertical;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
}

.form-actions {
  flex: 1 0 100%;
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.save-btn {
  background-color: var(--primary-color);
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 6px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: all 0.3s;
}

.save-btn:hover {
  background-color: var(--primary-dark);
}

.save-btn:disabled {
  background-color: #9ca3af;
  cursor: not-allowed;
}

@media (max-width: 768px) {
  .form-group {
    flex: 1 0 100%;
  }
}
</style> 