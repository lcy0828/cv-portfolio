<template>
  <div class="certificates-form-container">
    <div v-if="loading" class="loading">
      <i class="fas fa-spinner fa-spin"></i> 加载中...
    </div>
    <div v-else-if="error" class="error-message">
      <i class="fas fa-exclamation-circle"></i> {{ error }}
    </div>
    <div v-else class="certificates-content">
      <div class="actions-bar">
        <button class="add-btn" @click="showAddCertificate = true">
          <i class="fas fa-plus"></i> 添加证书
        </button>
      </div>
      
      <div class="certificates-list">
        <div v-for="certificate in certificates" :key="certificate.id" class="certificate-card">
          <div class="certificate-icon" v-if="certificate.icon">
            <i :class="certificate.icon"></i>
          </div>
          
          <div class="certificate-content">
            <h3 class="certificate-name">{{ certificate.name }}</h3>
            <div class="certificate-org">{{ certificate.organization }}</div>
            <div class="certificate-date">{{ certificate.date }}</div>
            
            <p v-if="certificate.description" class="certificate-description">
              {{ certificate.description }}
            </p>
            
            <a v-if="certificate.link" :href="certificate.link" target="_blank" class="certificate-link">
              <i class="fas fa-external-link-alt"></i> 查看证书
            </a>
          </div>
          
          <div class="certificate-actions">
            <button class="edit-btn" @click="editCertificate(certificate)">
              <i class="fas fa-edit"></i>
            </button>
            <button class="delete-btn" @click="confirmDeleteCertificate(certificate)">
              <i class="fas fa-trash"></i>
            </button>
          </div>
        </div>
        
        <div v-if="certificates.length === 0" class="no-data">
          暂无证书数据，请添加您的证书信息。
        </div>
      </div>
    </div>
    
    <!-- 证书表单模态框 -->
    <div v-if="showAddCertificate || editingCertificate" class="modal-overlay">
      <div class="modal-container">
        <div class="modal-header">
          <h3>{{ editingCertificate ? '编辑证书' : '添加新证书' }}</h3>
          <button class="close-btn" @click="closeCertificateForm">
            <i class="fas fa-times"></i>
          </button>
        </div>
        
        <div class="modal-body">
          <form class="certificate-form" @submit.prevent="saveCertificate">
            <div class="form-row">
              <div class="form-group">
                <label for="name">证书名称 <span class="required">*</span></label>
                <input type="text" id="name" v-model="certificateForm.name" required>
              </div>
            </div>
            
            <div class="form-row">
              <div class="form-group">
                <label for="organization">颁发机构 <span class="required">*</span></label>
                <input type="text" id="organization" v-model="certificateForm.organization" required>
              </div>
              
              <div class="form-group">
                <label for="date">获得日期</label>
                <input type="text" id="date" v-model="certificateForm.date" placeholder="例如: 2022年6月">
              </div>
            </div>
            
            <div class="form-group">
              <label for="description">证书描述</label>
              <textarea id="description" v-model="certificateForm.description" rows="3"></textarea>
            </div>
            
            <div class="form-row">
              <div class="form-group">
                <label for="icon">图标类名</label>
                <input type="text" id="icon" v-model="certificateForm.icon" placeholder="例如: fas fa-award">
                <div class="icon-preview" v-if="certificateForm.icon">
                  <i :class="certificateForm.icon"></i>
                  <span>图标预览</span>
                </div>
              </div>
              
              <div class="form-group">
                <label for="link">证书链接</label>
                <input type="url" id="link" v-model="certificateForm.link" placeholder="https://example.com">
              </div>
            </div>
            
            <div class="form-group">
              <label for="sortOrder">排序顺序</label>
              <input type="number" id="sortOrder" v-model.number="certificateForm.sortOrder" min="0">
            </div>
            
            <div class="form-actions">
              <button type="button" class="cancel-btn" @click="closeCertificateForm">取消</button>
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
          <p>您确定要删除以下证书吗？</p>
          <div class="confirm-item" v-if="deletingCertificate">
            <strong>{{ deletingCertificate.name }}</strong>
            <div>{{ deletingCertificate.organization }}</div>
          </div>
          <p class="warning-text">此操作无法撤销！</p>
        </div>
        
        <div class="dialog-actions">
          <button class="cancel-btn" @click="showDeleteConfirm = false">取消</button>
          <button class="delete-btn" @click="deleteCertificate" :disabled="saving">
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

// 证书数据
const certificates = ref([]);

// 模态框状态
const showAddCertificate = ref(false);
const editingCertificate = ref(null);
const showDeleteConfirm = ref(false);
const deletingCertificate = ref(null);

// 证书表单数据
const certificateForm = reactive({
  name: '',
  organization: '',
  date: '',
  description: '',
  icon: '',
  link: '',
  sortOrder: 0
});

// 重置表单
const resetForm = () => {
  certificateForm.name = '';
  certificateForm.organization = '';
  certificateForm.date = '';
  certificateForm.description = '';
  certificateForm.icon = '';
  certificateForm.link = '';
  certificateForm.sortOrder = certificates.value.length > 0 
    ? Math.max(...certificates.value.map(c => c.sortOrder)) + 1 
    : 0;
};

// 获取证书列表
const fetchCertificates = async () => {
  loading.value = true;
  error.value = null;
  
  try {
    const token = localStorage.getItem('token');
    const response = await axios.get(`${API_URL}/admin/certificates`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    if (response.data.success) {
      certificates.value = response.data.data;
    } else {
      error.value = response.data.message || '获取证书列表失败';
    }
  } catch (err) {
    console.error('获取证书列表出错:', err);
    error.value = '获取证书列表时发生错误，请稍后再试';
  } finally {
    loading.value = false;
  }
};

// 编辑证书
const editCertificate = (certificate) => {
  editingCertificate.value = { ...certificate };
  
  // 填充表单数据
  certificateForm.name = certificate.name;
  certificateForm.organization = certificate.organization;
  certificateForm.date = certificate.date;
  certificateForm.description = certificate.description;
  certificateForm.icon = certificate.icon;
  certificateForm.link = certificate.link;
  certificateForm.sortOrder = certificate.sortOrder;
  
  showAddCertificate.value = true;
};

// 关闭证书表单
const closeCertificateForm = () => {
  showAddCertificate.value = false;
  editingCertificate.value = null;
  resetForm();
};

// 确认删除证书
const confirmDeleteCertificate = (certificate) => {
  showDeleteConfirm.value = true;
  deletingCertificate.value = certificate;
};

// 保存证书
const saveCertificate = async () => {
  saving.value = true;
  error.value = null;
  
  try {
    const token = localStorage.getItem('token');
    const certificateData = {
      name: certificateForm.name,
      organization: certificateForm.organization,
      date: certificateForm.date,
      description: certificateForm.description,
      icon: certificateForm.icon,
      link: certificateForm.link,
      sortOrder: certificateForm.sortOrder
    };
    
    let response;
    if (editingCertificate.value) {
      // 更新证书
      response = await axios.put(`${API_URL}/admin/certificates/${editingCertificate.value.id}`, certificateData, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      
      if (response.data.success) {
        // 更新本地数据
        const index = certificates.value.findIndex(c => c.id === editingCertificate.value.id);
        if (index !== -1) {
          certificates.value[index] = { ...certificates.value[index], ...response.data.data };
        }
        showSuccessMessage('证书更新成功');
      }
    } else {
      // 创建新证书
      response = await axios.post(`${API_URL}/admin/certificates`, certificateData, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      
      if (response.data.success) {
        // 添加到本地数据
        certificates.value.push(response.data.data);
        showSuccessMessage('证书添加成功');
      }
    }
    
    // 关闭表单
    if (response.data.success) {
      closeCertificateForm();
    } else {
      error.value = response.data.message || '保存证书失败';
    }
  } catch (err) {
    console.error('保存证书出错:', err);
    error.value = '保存证书时发生错误，请稍后再试';
  } finally {
    saving.value = false;
  }
};

// 删除证书
const deleteCertificate = async () => {
  if (!deletingCertificate.value) return;
  
  saving.value = true;
  
  try {
    const token = localStorage.getItem('token');
    const response = await axios.delete(`${API_URL}/admin/certificates/${deletingCertificate.value.id}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    if (response.data.success) {
      // 从本地数据中删除
      const index = certificates.value.findIndex(c => c.id === deletingCertificate.value.id);
      if (index !== -1) {
        certificates.value.splice(index, 1);
      }
      
      showSuccessMessage('证书删除成功');
      showDeleteConfirm.value = false;
      deletingCertificate.value = null;
    } else {
      error.value = response.data.message || '删除证书失败';
    }
  } catch (err) {
    console.error('删除证书出错:', err);
    error.value = '删除证书时发生错误，请稍后再试';
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

// 页面加载时获取证书数据
onMounted(() => {
  fetchCertificates();
  resetForm();
});
</script>

<style scoped>
.certificates-form-container {
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

.certificates-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.certificate-card {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  padding: 20px;
  display: flex;
  position: relative;
  transition: all 0.3s;
}

.certificate-card:hover {
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.certificate-icon {
  width: 50px;
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--primary-color);
  color: white;
  border-radius: 50%;
  margin-right: 15px;
}

.certificate-icon i {
  font-size: 1.5rem;
}

.certificate-content {
  flex: 1;
}

.certificate-name {
  margin: 0 0 5px 0;
  font-size: 1.2rem;
  color: #1f2937;
}

.certificate-org {
  color: #4b5563;
  font-weight: 600;
  margin-bottom: 5px;
}

.certificate-date {
  color: #6b7280;
  font-size: 0.9rem;
  margin-bottom: 10px;
}

.certificate-description {
  color: #4b5563;
  font-size: 0.95rem;
  margin: 10px 0;
  line-height: 1.5;
}

.certificate-link {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  color: var(--primary-color);
  text-decoration: none;
  font-size: 0.9rem;
  margin-top: 10px;
}

.certificate-link:hover {
  text-decoration: underline;
}

.certificate-actions {
  position: absolute;
  top: 10px;
  right: 10px;
  display: flex;
  gap: 5px;
}

.edit-btn, .delete-btn {
  background: none;
  border: none;
  cursor: pointer;
  width: 32px;
  height: 32px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.edit-btn:hover {
  color: var(--primary-color);
  background-color: rgba(59, 130, 246, 0.1);
}

.delete-btn:hover {
  color: #e02424;
  background-color: rgba(224, 36, 36, 0.1);
}

.no-data {
  grid-column: 1 / -1;
  background-color: white;
  border-radius: 8px;
  padding: 40px;
  text-align: center;
  color: #6b7280;
}

/* 模态框样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-container {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.modal-header {
  padding: 20px;
  background-color: #f8fafc;
  border-bottom: 1px solid #e2e8f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.5rem;
  color: #1f2937;
}

.close-btn {
  background: none;
  border: none;
  cursor: pointer;
  color: #6b7280;
  width: 32px;
  height: 32px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-btn:hover {
  background-color: #f1f5f9;
  color: #1f2937;
}

.modal-body {
  padding: 20px;
  overflow-y: auto;
}

.certificate-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-row {
  display: flex;
  gap: 15px;
}

.form-group {
  flex: 1;
  margin-bottom: 5px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: 600;
  color: #374151;
  font-size: 0.95rem;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  background-color: white;
  font-size: 0.95rem;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.3);
}

.required {
  color: #e02424;
}

.icon-preview {
  margin-top: 10px;
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 0.9rem;
  color: #4b5563;
}

.icon-preview i {
  font-size: 1.5rem;
  color: var(--primary-color);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 15px;
  margin-top: 20px;
}

.cancel-btn, .save-btn {
  padding: 10px 20px;
  border-radius: 4px;
  font-weight: 600;
  cursor: pointer;
}

.cancel-btn {
  background-color: #f1f5f9;
  color: #4b5563;
  border: 1px solid #d1d5db;
}

.cancel-btn:hover {
  background-color: #e5e7eb;
}

.save-btn {
  background-color: var(--primary-color);
  color: white;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  min-width: 100px;
}

.save-btn:hover:not(:disabled) {
  background-color: var(--primary-dark);
}

.save-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

/* 确认对话框样式 */
.confirm-dialog {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
  width: 90%;
  max-width: 450px;
  overflow: hidden;
}

.dialog-header {
  padding: 15px 20px;
  background-color: #f8fafc;
  border-bottom: 1px solid #e2e8f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.dialog-header h3 {
  margin: 0;
  font-size: 1.3rem;
  color: #1f2937;
}

.dialog-body {
  padding: 20px;
}

.confirm-item {
  margin: 15px 0;
  padding: 10px 15px;
  background-color: #f8fafc;
  border-radius: 4px;
  border-left: 3px solid var(--primary-color);
}

.confirm-item div {
  color: #6b7280;
  font-size: 0.9rem;
  margin-top: 5px;
}

.warning-text {
  color: #e02424;
  font-size: 0.9rem;
  margin-top: 15px;
}

.dialog-actions {
  padding: 15px 20px;
  background-color: #f8fafc;
  border-top: 1px solid #e2e8f0;
  display: flex;
  justify-content: flex-end;
  gap: 15px;
}

.dialog-actions .delete-btn {
  background-color: #e02424;
  color: white;
  padding: 8px 15px;
  border-radius: 4px;
  border: none;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 8px;
  width: auto;
  height: auto;
}

.dialog-actions .delete-btn:hover:not(:disabled) {
  background-color: #b91c1c;
}

.dialog-actions .delete-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

@media (max-width: 768px) {
  .certificates-list {
    grid-template-columns: 1fr;
  }
  
  .form-row {
    flex-direction: column;
    gap: 0;
  }
}
</style> 