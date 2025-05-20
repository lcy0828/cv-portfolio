<template>
  <div class="settings-form-container">
    <div v-if="loading" class="loading">
      <i class="fas fa-spinner fa-spin"></i> 加载中...
    </div>
    <div v-else-if="error" class="error-message">
      <i class="fas fa-exclamation-circle"></i> {{ error }}
    </div>
    
    <!-- 密码修改表单 -->
    <div class="settings-card">
      <div class="card-header">
        <h3>修改密码</h3>
        <p>定期更新您的密码可以增强账户安全性</p>
      </div>
      
      <div class="card-body">
        <form class="password-form" @submit.prevent="changePassword">
          <div class="form-group">
            <label for="currentPassword">当前密码 <span class="required">*</span></label>
            <div class="password-input">
              <input 
                :type="showCurrentPassword ? 'text' : 'password'" 
                id="currentPassword" 
                v-model="passwordForm.currentPassword" 
                required
              >
              <button type="button" class="toggle-password" @click="showCurrentPassword = !showCurrentPassword">
                <i :class="showCurrentPassword ? 'fas fa-eye-slash' : 'fas fa-eye'"></i>
              </button>
            </div>
          </div>
          
          <div class="form-group">
            <label for="newPassword">新密码 <span class="required">*</span></label>
            <div class="password-input">
              <input 
                :type="showNewPassword ? 'text' : 'password'" 
                id="newPassword" 
                v-model="passwordForm.newPassword" 
                required
                pattern=".{6,}"
                title="密码长度至少为6个字符"
              >
              <button type="button" class="toggle-password" @click="showNewPassword = !showNewPassword">
                <i :class="showNewPassword ? 'fas fa-eye-slash' : 'fas fa-eye'"></i>
              </button>
            </div>
            <div class="password-hint">密码长度至少为6个字符</div>
          </div>
          
          <div class="form-group">
            <label for="confirmPassword">确认新密码 <span class="required">*</span></label>
            <div class="password-input">
              <input 
                :type="showConfirmPassword ? 'text' : 'password'" 
                id="confirmPassword" 
                v-model="passwordForm.confirmPassword" 
                required
              >
              <button type="button" class="toggle-password" @click="showConfirmPassword = !showConfirmPassword">
                <i :class="showConfirmPassword ? 'fas fa-eye-slash' : 'fas fa-eye'"></i>
              </button>
            </div>
            <div class="password-error" v-if="!passwordsMatch && passwordForm.confirmPassword">
              两次输入的密码不一致
            </div>
          </div>
          
          <div class="form-actions">
            <button type="submit" class="save-btn" :disabled="saving || !passwordsMatch">
              <i v-if="saving" class="fas fa-spinner fa-spin"></i>
              <span v-else>更新密码</span>
            </button>
          </div>
        </form>
      </div>
    </div>
    
    <!-- 系统信息卡片 -->
    <div class="settings-card">
      <div class="card-header">
        <h3>系统信息</h3>
        <p>网站的基本信息和版本</p>
      </div>
      
      <div class="card-body">
        <div class="info-row">
          <div class="info-label">系统名称</div>
          <div class="info-value">简历管理系统</div>
        </div>
        
        <div class="info-row">
          <div class="info-label">版本</div>
          <div class="info-value">1.0.0</div>
        </div>
        
        <div class="info-row">
          <div class="info-label">前端技术</div>
          <div class="info-value">Vue 3 + Vite</div>
        </div>
        
        <div class="info-row">
          <div class="info-label">后端技术</div>
          <div class="info-value">Go + Gin</div>
        </div>
        
        <div class="info-row">
          <div class="info-label">数据库</div>
          <div class="info-value">SQLite</div>
        </div>
      </div>
    </div>
    
    <!-- 成功提示 -->
    <div v-if="saveSuccess" class="success-message">
      <i class="fas fa-check-circle"></i> {{ successMessage }}
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue';
import axios from 'axios';
import { API_URL } from '../../config';

const loading = ref(false);
const saving = ref(false);
const error = ref(null);
const saveSuccess = ref(false);
const successMessage = ref('');

// 密码表单
const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
});

// 密码显示控制
const showCurrentPassword = ref(false);
const showNewPassword = ref(false);
const showConfirmPassword = ref(false);

// 密码匹配检查
const passwordsMatch = computed(() => {
  if (!passwordForm.confirmPassword) return true;
  return passwordForm.newPassword === passwordForm.confirmPassword;
});

// 修改密码
const changePassword = async () => {
  if (!passwordsMatch.value) return;
  
  saving.value = true;
  error.value = null;
  
  try {
    const token = localStorage.getItem('token');
    const response = await axios.put(`${API_URL}/admin/settings/password`, {
      current_password: passwordForm.currentPassword,
      new_password: passwordForm.newPassword
    }, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    if (response.data.success) {
      showSuccessMessage('密码修改成功');
      
      // 重置表单
      passwordForm.currentPassword = '';
      passwordForm.newPassword = '';
      passwordForm.confirmPassword = '';
    } else {
      error.value = response.data.message || '密码修改失败';
    }
  } catch (err) {
    console.error('修改密码出错:', err);
    if (err.response && err.response.status === 401) {
      error.value = '当前密码不正确';
    } else {
      error.value = '修改密码时发生错误，请稍后再试';
    }
  } finally {
    saving.value = false;
  }
};

// 显示成功消息
const showSuccessMessage = (message) => {
  successMessage.value = message;
  saveSuccess.value = true;
  setTimeout(() => {
    saveSuccess.value = false;
  }, 3000);
};

// 页面加载时初始化
onMounted(() => {
  // 可以在这里加载其他设置
});
</script>

<style scoped>
.settings-form-container {
  padding: 20px 0;
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 30px;
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

.settings-card {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.card-header {
  padding: 20px;
  background-color: #f8fafc;
  border-bottom: 1px solid #e2e8f0;
}

.card-header h3 {
  margin: 0 0 5px 0;
  font-size: 1.3rem;
  color: #1f2937;
}

.card-header p {
  margin: 0;
  font-size: 0.95rem;
  color: #6b7280;
}

.card-body {
  padding: 20px;
}

.password-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
  max-width: 500px;
}

.form-group {
  margin-bottom: 5px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: 600;
  color: #374151;
  font-size: 0.95rem;
}

.password-input {
  position: relative;
  display: flex;
  align-items: center;
}

.password-input input {
  width: 100%;
  padding: 10px;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  background-color: white;
  font-size: 0.95rem;
  padding-right: 40px;
}

.password-input input:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.3);
}

.toggle-password {
  position: absolute;
  right: 10px;
  background: none;
  border: none;
  color: #6b7280;
  cursor: pointer;
  padding: 5px;
}

.toggle-password:hover {
  color: var(--primary-color);
}

.password-hint {
  font-size: 0.85rem;
  color: #6b7280;
  margin-top: 5px;
}

.password-error {
  font-size: 0.85rem;
  color: #e02424;
  margin-top: 5px;
}

.required {
  color: #e02424;
}

.form-actions {
  display: flex;
  justify-content: flex-start;
  margin-top: 10px;
}

.save-btn {
  background-color: var(--primary-color);
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  min-width: 120px;
}

.save-btn:hover:not(:disabled) {
  background-color: var(--primary-dark);
}

.save-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.info-row {
  display: flex;
  margin-bottom: 15px;
  padding-bottom: 15px;
  border-bottom: 1px solid #e2e8f0;
}

.info-row:last-child {
  margin-bottom: 0;
  padding-bottom: 0;
  border-bottom: none;
}

.info-label {
  width: 100px;
  color: #4b5563;
  font-weight: 600;
}

.info-value {
  flex: 1;
  color: #1f2937;
}

@media (max-width: 768px) {
  .info-row {
    flex-direction: column;
  }
  
  .info-label {
    width: 100%;
    margin-bottom: 5px;
  }
}
</style> 