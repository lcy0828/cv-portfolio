<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <h1>简历管理系统</h1>
        <p>请登录以管理您的个人简历内容</p>
      </div>
      
      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label for="username">用户名</label>
          <div class="input-group">
            <i class="fas fa-user"></i>
            <input 
              type="text"
              id="username"
              v-model="username"
              placeholder="请输入用户名"
              required
            />
          </div>
        </div>
        
        <div class="form-group">
          <label for="password">密码</label>
          <div class="input-group">
            <i class="fas fa-lock"></i>
            <input 
              type="password"
              id="password"
              v-model="password"
              placeholder="请输入密码"
              required
            />
          </div>
        </div>
        
        <div v-if="error" class="error-message">
          <i class="fas fa-exclamation-triangle"></i>
          {{ error }}
        </div>
        
        <button type="submit" :disabled="isLoading" class="login-btn">
          <span v-if="!isLoading">登 录</span>
          <span v-else class="loading-spinner"></span>
        </button>
      </form>
      
      <div class="login-footer">
        <p>返回 <a href="/">个人主页</a></p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import apiService from '../services/api';

const router = useRouter();
const username = ref('');
const password = ref('');
const error = ref('');
const isLoading = ref(false);

const handleLogin = async () => {
  error.value = '';
  isLoading.value = true;
  
  try {
    const response = await apiService.login({
      username: username.value,
      password: password.value
    });
    
    if (response.success) {
      // 存储token
      localStorage.setItem('token', response.data.token);
      // 重定向到管理界面
      router.push('/admin');
    } else {
      error.value = response.message || '登录失败，请检查用户名和密码';
    }
  } catch (err) {
    error.value = '登录请求失败，请稍后重试';
    console.error('登录错误:', err);
  } finally {
    isLoading.value = false;
  }
};
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f8f9fa;
  padding: 20px;
}

.login-card {
  width: 100%;
  max-width: 450px;
  background-color: white;
  border-radius: 10px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.login-header {
  padding: 30px;
  background-color: var(--primary-color);
  color: white;
  text-align: center;
}

.login-header h1 {
  font-size: 1.8rem;
  margin-bottom: 10px;
}

.login-header p {
  opacity: 0.8;
  font-size: 0.95rem;
}

.login-form {
  padding: 30px;
}

.form-group {
  margin-bottom: 25px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: var(--text-color);
  font-size: 0.9rem;
}

.input-group {
  position: relative;
}

.input-group i {
  position: absolute;
  left: 15px;
  top: 50%;
  transform: translateY(-50%);
  color: #6c757d;
}

.input-group input {
  width: 100%;
  padding: 12px 15px 12px 45px;
  border: 1px solid #dee2e6;
  border-radius: 5px;
  font-size: 1rem;
  transition: border-color 0.3s;
}

.input-group input:focus {
  outline: none;
  border-color: var(--accent-color);
}

.error-message {
  background-color: #fee2e2;
  color: #ef4444;
  padding: 12px;
  border-radius: 5px;
  margin-bottom: 20px;
  font-size: 0.9rem;
  display: flex;
  align-items: center;
}

.error-message i {
  margin-right: 10px;
  font-size: 1rem;
}

.login-btn {
  width: 100%;
  padding: 12px;
  background-color: var(--primary-color);
  color: white;
  border: none;
  border-radius: 5px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.3s;
}

.login-btn:hover {
  background-color: var(--secondary-color);
}

.login-btn:disabled {
  background-color: #6c757d;
  cursor: not-allowed;
}

.loading-spinner {
  display: inline-block;
  width: 20px;
  height: 20px;
  border: 3px solid rgba(255,255,255,0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.login-footer {
  padding: 20px;
  text-align: center;
  border-top: 1px solid #dee2e6;
  font-size: 0.9rem;
}

.login-footer a {
  color: var(--accent-color);
  font-weight: 600;
}
</style> 