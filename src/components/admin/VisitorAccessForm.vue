<template>
  <div class="visitor-access-form-container">
    <div v-if="loading" class="loading">
      <i class="fas fa-spinner fa-spin"></i> 加载中...
    </div>
    <div v-else-if="error" class="error-message">
      <i class="fas fa-exclamation-circle"></i> {{ error }}
    </div>
    <div v-else class="visitor-content">
      <div class="section-intro">
        <h3>访客权限管理</h3>
        <p>您可以设置多个访问密码，供不同的访客使用。所有访客必须通过验证才能查看您的简历。</p>
      </div>
      
      <div class="actions-bar">
        <button class="add-btn" @click="showAddAccess = true">
          <i class="fas fa-plus"></i> 添加访问密码
        </button>
      </div>
      
      <div class="access-list">
        <div class="table-responsive">
          <table class="access-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>类型</th>
                <th>密码值</th>
                <th>访问标识</th>
                <th>创建时间</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="access in accessList" :key="access.id">
                <td>{{ access.id }}</td>
                <td>{{ formatAccessType(access.access_type) }}</td>
                <td class="password-value">
                  <div class="masked-password" v-if="!visiblePasswords[access.id]">
                    •••••••••••
                  </div>
                  <div class="visible-password" v-else>
                    {{ access.value }}
                  </div>
                  <button class="toggle-password" @click="togglePasswordVisibility(access.id)">
                    <i :class="visiblePasswords[access.id] ? 'fas fa-eye-slash' : 'fas fa-eye'"></i>
                  </button>
                </td>
                <td>{{ access.access_key }}</td>
                <td>{{ formatDate(access.created_at) }}</td>
                <td>
                  <button class="delete-btn" @click="confirmDeleteAccess(access)">
                    <i class="fas fa-trash"></i>
                  </button>
                </td>
              </tr>
              <tr v-if="accessList.length === 0">
                <td colspan="6" class="no-data">暂无访客密码数据，请添加访客密码。</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
    
    <!-- 添加访客密码模态框 -->
    <div v-if="showAddAccess" class="modal-overlay">
      <div class="modal-container">
        <div class="modal-header">
          <h3>添加访客密码</h3>
          <button class="close-btn" @click="closeAccessForm">
            <i class="fas fa-times"></i>
          </button>
        </div>
        
        <div class="modal-body">
          <form class="access-form" @submit.prevent="saveAccess">
            <div class="form-group">
              <label for="accessType">访问类型</label>
              <select id="accessType" v-model="accessForm.accessType" disabled>
                <option value="password">密码</option>
              </select>
              <div class="form-hint">目前仅支持密码类型</div>
            </div>
            
            <div class="form-group">
              <label for="value">密码值 <span class="required">*</span></label>
              <input type="text" id="value" v-model="accessForm.value" required minlength="4">
              <div class="form-hint">访客需要输入此密码才能访问您的简历</div>
            </div>
            
            <div class="form-group">
              <label for="accessKey">访问标识 <span class="required">*</span></label>
              <input type="text" id="accessKey" v-model="accessForm.accessKey" required>
              <div class="form-hint">用于标识不同的访客群体，例如：HR、朋友、同事等</div>
            </div>
            
            <div class="form-actions">
              <button type="button" class="cancel-btn" @click="closeAccessForm">取消</button>
              <button type="submit" class="save-btn" :disabled="saving">
                <i v-if="saving" class="fas fa-spinner fa-spin"></i>
                <span v-else>保存</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
    
    <!-- 删除确认对话框 -->
    <div v-if="showDeleteConfirm" class="modal-overlay">
      <div class="confirm-dialog">
        <div class="dialog-header">
          <h3>确认删除</h3>
          <button class="close-btn" @click="showDeleteConfirm = false">
            <i class="fas fa-times"></i>
          </button>
        </div>
        
        <div class="dialog-body">
          <p>您确定要删除以下访客密码吗？</p>
          <div class="confirm-item" v-if="deletingAccess">
            <div><strong>访问类型:</strong> {{ formatAccessType(deletingAccess.access_type) }}</div>
            <div><strong>密码值:</strong> {{ deletingAccess.value }}</div>
            <div><strong>访问标识:</strong> {{ deletingAccess.access_key }}</div>
          </div>
          <p class="warning-text">删除后，使用此密码的访客将无法访问您的简历！</p>
        </div>
        
        <div class="dialog-actions">
          <button class="cancel-btn" @click="showDeleteConfirm = false">取消</button>
          <button class="delete-btn" @click="deleteAccess" :disabled="saving">
            <i v-if="saving" class="fas fa-spinner fa-spin"></i>
            <span v-else>确认删除</span>
          </button>
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
import { ref, reactive, onMounted } from 'vue';
import axios from 'axios';
import { API_URL } from '../../config';

const loading = ref(false);
const saving = ref(false);
const error = ref(null);
const saveSuccess = ref(false);
const successMessage = ref('');

// 访客密码数据
const accessList = ref([]);
const visiblePasswords = reactive({});

// 模态框状态
const showAddAccess = ref(false);
const showDeleteConfirm = ref(false);
const deletingAccess = ref(null);

// 访客密码表单
const accessForm = reactive({
  accessType: 'password',
  value: '',
  accessKey: ''
});

// 重置表单
const resetForm = () => {
  accessForm.accessType = 'password';
  accessForm.value = '';
  accessForm.accessKey = '';
};

// 格式化访问类型
const formatAccessType = (type) => {
  switch (type) {
    case 'password': return '密码';
    default: return type;
  }
};

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return '';
  const date = new Date(dateStr);
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  });
};

// 切换密码可见性
const togglePasswordVisibility = (id) => {
  visiblePasswords[id] = !visiblePasswords[id];
};

// 获取访客密码列表
const fetchAccessList = async () => {
  loading.value = true;
  error.value = null;
  
    try {    const token = localStorage.getItem('token');    console.log("获取访客密码列表，使用令牌类型: 管理员令牌");        const response = await axios.get(`${API_URL}/admin/visitor/access`, {      headers: {        'Authorization': `Bearer ${token}`      }    });        console.log("访客密码API响应:", response.data);        if (response.data.success) {      accessList.value = response.data.data || [];    } else {      error.value = response.data.message || '获取访客密码列表失败';    }
  } catch (err) {
    console.error('获取访客密码列表出错:', err);
    error.value = '获取访客密码列表时发生错误，请稍后再试';
  } finally {
    loading.value = false;
  }
};

// 关闭访客密码表单
const closeAccessForm = () => {
  showAddAccess.value = false;
  resetForm();
};

// 确认删除访客密码
const confirmDeleteAccess = (access) => {
  showDeleteConfirm.value = true;
  deletingAccess.value = access;
};

// 保存访客密码
const saveAccess = async () => {
  saving.value = true;
  error.value = null;
  
  try {
    const token = localStorage.getItem('token');
    const accessData = {
      access_type: accessForm.accessType,
      value: accessForm.value,
      access_key: accessForm.accessKey
    };
    
    const response = await axios.post(`${API_URL}/admin/visitor/access`, accessData, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    if (response.data.success) {
      await fetchAccessList(); // 重新获取列表
      showSuccessMessage('访客密码添加成功');
      closeAccessForm();
    } else {
      error.value = response.data.message || '保存访客密码失败';
    }
  } catch (err) {
    console.error('保存访客密码出错:', err);
    if (err.response && err.response.status === 400) {
      error.value = '密码值不能重复，请使用其他密码值';
    } else {
      error.value = '保存访客密码时发生错误，请稍后再试';
    }
  } finally {
    saving.value = false;
  }
};

// 删除访客密码
const deleteAccess = async () => {
  if (!deletingAccess.value) return;
  
  saving.value = true;
  
  try {
    const token = localStorage.getItem('token');
    const response = await axios.delete(`${API_URL}/admin/visitor/access/${deletingAccess.value.id}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    if (response.data.success) {
      await fetchAccessList(); // 重新获取列表
      showSuccessMessage('访客密码删除成功');
      showDeleteConfirm.value = false;
      deletingAccess.value = null;
    } else {
      error.value = response.data.message || '删除访客密码失败';
    }
  } catch (err) {
    console.error('删除访客密码出错:', err);
    error.value = '删除访客密码时发生错误，请稍后再试';
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

// 页面加载时获取访客密码数据
onMounted(() => {
  fetchAccessList();
});
</script>

<style scoped>
.visitor-access-form-container {
  padding: 20px 0;
  position: relative;
}

.section-intro {
  margin-bottom: 20px;
}

.section-intro h3 {
  margin: 0 0 10px 0;
  font-size: 1.5rem;
  color: #1f2937;
}

.section-intro p {
  color: #6b7280;
  margin: 0;
  line-height: 1.5;
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

.actions-bar {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 20px;
}

.add-btn {
  background-color: var(--primary-color);
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.add-btn:hover {
  background-color: var(--primary-dark);
}

.table-responsive {
  overflow-x: auto;
}

.access-table {
  width: 100%;
  border-collapse: collapse;
  background-color: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.access-table th, .access-table td {
  padding: 12px 15px;
  text-align: left;
  border-bottom: 1px solid #e2e8f0;
}

.access-table th {
  background-color: #f8fafc;
  font-weight: 600;
  color: #4b5563;
}

.access-table tr:last-child td {
  border-bottom: none;
}

.password-value {
  position: relative;
  padding-right: 40px;
}

.masked-password, .visible-password {
  font-family: monospace;
}

.toggle-password {
  position: absolute;
  right: 10px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: #6b7280;
  cursor: pointer;
  padding: 5px;
}

.toggle-password:hover {
  color: var(--primary-color);
}

.delete-btn {
  background-color: #fee2e2;
  color: #dc2626;
  border: none;
  border-radius: 4px;
  padding: 6px 10px;
  cursor: pointer;
}

.delete-btn:hover {
  background-color: #fecaca;
}

.no-data {
  text-align: center;
  color: #6b7280;
  padding: 30px 0;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 100;
  padding: 20px;
}

.modal-container {
  background-color: white;
  border-radius: 8px;
  width: 100%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.confirm-dialog {
  background-color: white;
  border-radius: 8px;
  width: 100%;
  max-width: 450px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.modal-header, .dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  border-bottom: 1px solid #e2e8f0;
}

.modal-header h3, .dialog-header h3 {
  margin: 0;
  font-size: 1.25rem;
  color: #1f2937;
}

.close-btn {
  background: none;
  border: none;
  color: #6b7280;
  font-size: 1.2rem;
  cursor: pointer;
}

.close-btn:hover {
  color: #1f2937;
}

.modal-body, .dialog-body {
  padding: 20px;
}

.access-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
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

.form-group input, .form-group select {
  width: 100%;
  padding: 10px;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  background-color: white;
  font-size: 0.95rem;
}

.form-group input:focus, .form-group select:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.3);
}

.form-hint {
  font-size: 0.85rem;
  color: #6b7280;
  margin-top: 5px;
}

.required {
  color: #e02424;
}

.form-actions, .dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
}

.save-btn, .cancel-btn, .delete-btn {
  padding: 8px 16px;
  border-radius: 4px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.save-btn {
  background-color: var(--primary-color);
  color: white;
  border: none;
}

.save-btn:hover:not(:disabled) {
  background-color: var(--primary-dark);
}

.cancel-btn {
  background-color: #f3f4f6;
  color: #4b5563;
  border: 1px solid #d1d5db;
}

.cancel-btn:hover {
  background-color: #e5e7eb;
}

.delete-btn {
  background-color: #ef4444;
  color: white;
  border: none;
}

.delete-btn:hover:not(:disabled) {
  background-color: #dc2626;
}

.save-btn:disabled, .delete-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.confirm-item {
  background-color: #f8fafc;
  padding: 15px;
  border-radius: 5px;
  margin: 15px 0;
  font-size: 0.95rem;
}

.confirm-item div {
  margin-bottom: 5px;
}

.confirm-item div:last-child {
  margin-bottom: 0;
}

.warning-text {
  color: #dc2626;
  font-weight: 600;
}

@media (max-width: 768px) {
  .access-table {
    font-size: 0.85rem;
  }
  
  .access-table th, .access-table td {
    padding: 10px;
  }
}
</style> 