<template>
  <div class="experience-form-container">
    <div v-if="loading" class="loading">
      <i class="fas fa-spinner fa-spin"></i> 加载中...
    </div>
    <div v-else-if="error" class="error-message">
      <i class="fas fa-exclamation-circle"></i> {{ error }}
    </div>
    <div v-else class="experience-content">
      <div class="actions-bar">
        <button class="add-btn" @click="showAddExperience = true">
          <i class="fas fa-plus"></i> 添加工作经历
        </button>
      </div>
      
      <div class="experience-list">
        <div v-for="exp in experiences" :key="exp.id" class="experience-card">
          <div class="card-header" :style="{ backgroundColor: exp.color || '#f1f5f9' }">
            <div class="exp-icon">
              <i :class="exp.icon || 'fas fa-briefcase'"></i>
            </div>
            <div class="exp-period">{{ exp.period }}</div>
            <div class="card-actions">
              <button class="edit-btn" @click="editExperience(exp)">
                <i class="fas fa-edit"></i>
              </button>
              <button class="delete-btn" @click="confirmDeleteExperience(exp)">
                <i class="fas fa-trash"></i>
              </button>
            </div>
          </div>
          
          <div class="card-body">
            <h3 class="exp-title">{{ exp.title }}</h3>
            <div class="exp-company">
              <span>{{ exp.company }}</span>
              <span class="exp-location">{{ exp.location }}</span>
            </div>
            
            <div class="exp-details">
              <div class="details-section" v-if="exp.responsibilities && exp.responsibilities.length">
                <h4>工作职责</h4>
                <ul>
                  <li v-for="(item, idx) in exp.responsibilities" :key="idx">{{ item }}</li>
                </ul>
              </div>
              
              <div class="details-section" v-if="exp.achievements && exp.achievements.length">
                <h4>工作成就</h4>
                <ul>
                  <li v-for="(item, idx) in exp.achievements" :key="idx">{{ item }}</li>
                </ul>
              </div>
              
              <div class="details-section" v-if="exp.technologies && exp.technologies.length">
                <h4>使用技术</h4>
                <div class="technologies-list">
                  <span v-for="(tech, idx) in exp.technologies" :key="idx" class="tech-tag">
                    {{ tech }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <div v-if="experiences.length === 0" class="no-data">
          暂无工作经历，请添加您的工作经历。
        </div>
      </div>
    </div>
    
    <!-- 添加/编辑工作经历的模态框 -->
    <div class="modal" v-if="showAddExperience || editingExperience">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ editingExperience ? '编辑工作经历' : '添加工作经历' }}</h3>
          <button class="close-btn" @click="closeModal">
            <i class="fas fa-times"></i>
          </button>
        </div>
        
        <form @submit.prevent="saveExperience" class="modal-form">
          <div class="form-row">
            <div class="form-group">
              <label for="period">工作时间段</label>
              <input 
                type="text" 
                id="period" 
                v-model="expForm.period" 
                placeholder="例如：2020 - 至今"
                required
              >
            </div>
            
            <div class="form-group">
              <label for="sort-order">排序 (数字越小越靠前)</label>
              <input 
                type="number" 
                id="sort-order" 
                v-model.number="expForm.sortOrder" 
                min="1"
              >
            </div>
          </div>
          
          <div class="form-row">
            <div class="form-group">
              <label for="title">职位</label>
              <input 
                type="text" 
                id="title" 
                v-model="expForm.title" 
                required
              >
            </div>
            
            <div class="form-group">
              <label for="company">公司名称</label>
              <input 
                type="text" 
                id="company" 
                v-model="expForm.company" 
                required
              >
            </div>
          </div>
          
          <div class="form-row">
            <div class="form-group">
              <label for="location">工作地点</label>
              <input 
                type="text" 
                id="location" 
                v-model="expForm.location"
              >
            </div>
            
            <div class="form-group">
              <label for="color">卡片颜色</label>
              <input 
                type="color" 
                id="color" 
                v-model="expForm.color"
              >
            </div>
          </div>
          
          <div class="form-group">
            <label for="icon">图标类名</label>
            <input 
              type="text" 
              id="icon" 
              v-model="expForm.icon" 
              placeholder="例如：fas fa-briefcase"
            >
            <div class="icon-preview" v-if="expForm.icon">
              <i :class="expForm.icon"></i> 预览
            </div>
          </div>
          
          <div class="form-group">
            <label>工作职责 (每行一项)</label>
            <div class="list-editor">
              <div v-for="(item, idx) in expForm.responsibilities" :key="idx" class="list-item">
                <input 
                  type="text" 
                  v-model="expForm.responsibilities[idx]" 
                  :placeholder="`职责 ${idx + 1}`"
                >
                <button type="button" class="remove-btn" @click="removeListItem('responsibilities', idx)">
                  <i class="fas fa-times"></i>
                </button>
              </div>
              <button type="button" class="add-list-item" @click="addListItem('responsibilities')">
                <i class="fas fa-plus"></i> 添加职责
              </button>
            </div>
          </div>
          
          <div class="form-group">
            <label>工作成就 (每行一项)</label>
            <div class="list-editor">
              <div v-for="(item, idx) in expForm.achievements" :key="idx" class="list-item">
                <input 
                  type="text" 
                  v-model="expForm.achievements[idx]" 
                  :placeholder="`成就 ${idx + 1}`"
                >
                <button type="button" class="remove-btn" @click="removeListItem('achievements', idx)">
                  <i class="fas fa-times"></i>
                </button>
              </div>
              <button type="button" class="add-list-item" @click="addListItem('achievements')">
                <i class="fas fa-plus"></i> 添加成就
              </button>
            </div>
          </div>
          
          <div class="form-group">
            <label for="technologies">使用技术 (用逗号分隔)</label>
            <input 
              type="text" 
              id="technologies" 
              v-model="techsInput" 
              placeholder="例如：React, JavaScript, CSS"
            >
            <div class="techs-preview" v-if="parsedTechs.length">
              <span v-for="(tech, idx) in parsedTechs" :key="idx" class="tech-tag">
                {{ tech }}
              </span>
            </div>
          </div>
          
          <div class="form-actions">
            <button type="button" class="cancel-btn" @click="closeModal">取消</button>
            <button type="submit" class="save-btn" :disabled="saving">
              {{ saving ? '保存中...' : '保存' }}
            </button>
          </div>
        </form>
      </div>
    </div>
    
    <!-- 确认删除的模态框 -->
    <div class="modal" v-if="showDeleteConfirm">
      <div class="modal-content confirm-modal">
        <div class="modal-header">
          <h3>确认删除</h3>
          <button class="close-btn" @click="closeModal">
            <i class="fas fa-times"></i>
          </button>
        </div>
        
        <div class="confirm-message">
          <i class="fas fa-exclamation-triangle"></i>
          <p>确定要删除 {{ deletingExperience?.company }} 的 {{ deletingExperience?.title }} 职位记录吗？</p>
        </div>
        
        <div class="form-actions">
          <button type="button" class="cancel-btn" @click="closeModal">取消</button>
          <button type="button" class="delete-confirm-btn" @click="deleteExperience" :disabled="saving">
            {{ saving ? '删除中...' : '确认删除' }}
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
import { ref, reactive, computed, onMounted } from 'vue';
import axios from 'axios';

const API_URL = 'http://localhost:8080/api';
const loading = ref(false);
const saving = ref(false);
const error = ref(null);
const saveSuccess = ref(false);
const successMessage = ref('');

// 工作经历数据
const experiences = ref([]);

// 模态框状态
const showAddExperience = ref(false);
const editingExperience = ref(null);
const showDeleteConfirm = ref(false);
const deletingExperience = ref(null);

// 表单数据
const expForm = reactive({
  id: null,
  period: '',
  title: '',
  company: '',
  location: '',
  color: '#1E40AF',
  icon: 'fas fa-briefcase',
  responsibilities: [],
  achievements: [],
  technologies: [],
  sortOrder: 1
});

const techsInput = ref('');

const parsedTechs = computed(() => {
  return techsInput.value
    .split(',')
    .map(tech => tech.trim())
    .filter(tech => tech !== '');
});

// 获取工作经历列表
const fetchExperiences = async () => {
  loading.value = true;
  error.value = null;
  
  try {
    const token = localStorage.getItem('token');
    const response = await axios.get(`${API_URL}/admin/experiences`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    if (response.data.success) {
      experiences.value = response.data.data;
    } else {
      error.value = response.data.message || '获取工作经历失败';
    }
  } catch (err) {
    console.error('获取工作经历出错:', err);
    error.value = '获取工作经历时发生错误，请稍后再试';
  } finally {
    loading.value = false;
  }
};

// 编辑工作经历
const editExperience = (exp) => {
  editingExperience.value = exp;
  
  // 填充表单数据
  expForm.id = exp.id;
  expForm.period = exp.period;
  expForm.title = exp.title;
  expForm.company = exp.company;
  expForm.location = exp.location || '';
  expForm.color = exp.color || '#1E40AF';
  expForm.icon = exp.icon || 'fas fa-briefcase';
  expForm.responsibilities = [...(exp.responsibilities || [])];
  expForm.achievements = [...(exp.achievements || [])];
  expForm.technologies = [...(exp.technologies || [])];
  expForm.sortOrder = exp.sortOrder || 1;
  
  // 更新技术输入框
  techsInput.value = (exp.technologies || []).join(', ');
};

// 确认删除工作经历
const confirmDeleteExperience = (exp) => {
  showDeleteConfirm.value = true;
  deletingExperience.value = exp;
};

// 删除工作经历
const deleteExperience = async () => {
  if (!deletingExperience.value) return;
  
  saving.value = true;
  
  try {
    const token = localStorage.getItem('token');
    const response = await axios.delete(`${API_URL}/admin/experiences/${deletingExperience.value.id}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    if (response.data.success) {
      await fetchExperiences();
      showSuccessMessage('工作经历删除成功');
    } else {
      error.value = response.data.message || '删除工作经历失败';
    }
  } catch (err) {
    console.error('删除工作经历出错:', err);
    error.value = '删除工作经历时发生错误，请稍后再试';
  } finally {
    saving.value = false;
    closeModal();
  }
};

// 保存工作经历
const saveExperience = async () => {
  saving.value = true;
  
  // 更新技术标签
  expForm.technologies = parsedTechs.value;
  
  try {
    const token = localStorage.getItem('token');
    const isEditing = !!editingExperience.value;
    const url = isEditing 
      ? `${API_URL}/admin/experiences/${expForm.id}` 
      : `${API_URL}/admin/experiences`;
    const method = isEditing ? 'put' : 'post';
    
    const response = await axios({
      method,
      url,
      data: expForm,
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });
    
    if (response.data.success) {
      await fetchExperiences();
      showSuccessMessage(isEditing ? '工作经历更新成功' : '工作经历添加成功');
    } else {
      error.value = response.data.message || (isEditing ? '更新工作经历失败' : '添加工作经历失败');
    }
  } catch (err) {
    console.error('保存工作经历出错:', err);
    error.value = '保存工作经历时发生错误，请稍后再试';
  } finally {
    saving.value = false;
    closeModal();
  }
};

// 添加列表项
const addListItem = (field) => {
  expForm[field].push('');
};

// 移除列表项
const removeListItem = (field, index) => {
  expForm[field].splice(index, 1);
};

// 关闭模态框
const closeModal = () => {
  showAddExperience.value = false;
  editingExperience.value = null;
  showDeleteConfirm.value = false;
  deletingExperience.value = null;
  
  // 重置表单
  expForm.id = null;
  expForm.period = '';
  expForm.title = '';
  expForm.company = '';
  expForm.location = '';
  expForm.color = '#1E40AF';
  expForm.icon = 'fas fa-briefcase';
  expForm.responsibilities = [];
  expForm.achievements = [];
  expForm.technologies = [];
  expForm.sortOrder = 1;
  
  techsInput.value = '';
};

// 显示成功消息
const showSuccessMessage = (message) => {
  successMessage.value = message;
  saveSuccess.value = true;
  setTimeout(() => {
    saveSuccess.value = false;
  }, 3000);
};

// 页面加载时获取工作经历数据
onMounted(() => {
  fetchExperiences();
});
</script>

<style scoped>
.experience-form-container {
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

.experience-list {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.experience-card {
  background-color: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.card-header {
  position: relative;
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 15px;
  color: white;
}

.exp-icon {
  background-color: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2rem;
}

.exp-period {
  font-weight: 600;
  font-size: 1.1rem;
}

.card-actions {
  position: absolute;
  right: 15px;
  top: 15px;
  display: flex;
  gap: 10px;
  z-index: 2;
}

.edit-btn, .delete-btn {
  background-color: rgba(255, 255, 255, 0.2);
  border: none;
  color: white;
  width: 32px;
  height: 32px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: background-color 0.2s;
}

.edit-btn:hover {
  background-color: rgba(255, 255, 255, 0.3);
}

.delete-btn:hover {
  background-color: rgba(220, 38, 38, 0.8);
}

.card-body {
  padding: 20px;
}

.exp-title {
  margin: 0 0 5px 0;
  font-size: 1.3rem;
  color: #1f2937;
}

.exp-company {
  font-size: 1.1rem;
  color: #4b5563;
  margin-bottom: 15px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.exp-location {
  font-size: 0.9rem;
  color: #6b7280;
  display: flex;
  align-items: center;
  gap: 5px;
}

.exp-location::before {
  content: '•';
  color: #9ca3af;
}

.exp-details {
  display: flex;
  flex-direction: column;
  gap: 15px;
  color: #4b5563;
}

.details-section h4 {
  margin: 0 0 8px 0;
  font-size: 1rem;
  color: #374151;
}

.details-section ul {
  margin: 0;
  padding-left: 20px;
}

.details-section li {
  margin-bottom: 5px;
}

.details-section li:last-child {
  margin-bottom: 0;
}

.technologies-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tech-tag {
  background-color: #f1f5f9;
  color: #4b5563;
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 0.9rem;
}

.no-data {
  background-color: white;
  border-radius: 8px;
  padding: 40px;
  text-align: center;
  color: #6b7280;
}

.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background-color: white;
  border-radius: 8px;
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.confirm-modal {
  max-width: 400px;
}

.modal-header {
  padding: 15px 20px;
  border-bottom: 1px solid #e2e8f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: sticky;
  top: 0;
  background-color: white;
  z-index: 10;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.2rem;
}

.close-btn {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 1.2rem;
  color: #6b7280;
}

.close-btn:hover {
  color: #111827;
}

.modal-form {
  padding: 20px;
}

.form-row {
  display: flex;
  gap: 15px;
  margin-bottom: 15px;
}

.form-row .form-group {
  flex: 1;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: 600;
  color: #374151;
}

.form-group input,
.form-group textarea,
.form-group select {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  font-size: 1rem;
}

.form-group input[type="color"] {
  padding: 2px;
  height: 40px;
}

.icon-preview {
  margin-top: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
  color: #4b5563;
}

.list-editor {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.list-item {
  display: flex;
  align-items: center;
  gap: 10px;
}

.remove-btn {
  background-color: #f3f4f6;
  border: none;
  width: 32px;
  height: 32px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  flex-shrink: 0;
}

.remove-btn:hover {
  background-color: #e5e7eb;
  color: #dc2626;
}

.add-list-item {
  background: none;
  border: 1px dashed #d1d5db;
  padding: 8px 12px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  cursor: pointer;
  color: #6b7280;
  width: 100%;
}

.add-list-item:hover {
  background-color: #f3f4f6;
  color: var(--primary-color);
}

.techs-preview {
  margin-top: 8px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
  position: sticky;
  bottom: 0;
  background-color: white;
  padding-top: 15px;
  border-top: 1px solid #e2e8f0;
}

.cancel-btn {
  background-color: #f3f4f6;
  color: #374151;
  border: none;
  padding: 10px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 600;
}

.save-btn {
  background-color: var(--primary-color);
  color: white;
  border: none;
  padding: 10px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 600;
}

.delete-confirm-btn {
  background-color: #e02424;
  color: white;
  border: none;
  padding: 10px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 600;
}

.confirm-message {
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 15px;
}

.confirm-message i {
  font-size: 2rem;
  color: #e02424;
}

@media (max-width: 768px) {
  .form-row {
    flex-direction: column;
    gap: 0;
  }
}
</style> 