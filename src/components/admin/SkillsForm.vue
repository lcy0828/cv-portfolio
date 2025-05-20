<template>
  <div class="skills-form-container">
    <div v-if="loading" class="loading">
      <i class="fas fa-spinner fa-spin"></i> 加载中...
    </div>
    <div v-else-if="error" class="error-message">
      <i class="fas fa-exclamation-circle"></i> {{ error }}
    </div>
    <div v-else class="skills-content">
      <!-- 分类列表 -->
      <div class="categories-section">
        <div class="section-header">
          <h3>技能分类</h3>
          <button class="add-btn" @click="showAddCategory = true">
            <i class="fas fa-plus"></i> 添加分类
          </button>
        </div>
        
        <div class="categories-list">
          <div 
            v-for="category in categories" 
            :key="category.id" 
            class="category-item"
            :class="{ active: selectedCategory && selectedCategory.id === category.id }"
            @click="selectCategory(category)"
          >
            <div class="category-info">
              <i :class="category.icon"></i>
              <span>{{ category.name }}</span>
            </div>
            <div class="category-actions">
              <button class="edit-btn" @click.stop="editCategory(category)">
                <i class="fas fa-edit"></i>
              </button>
              <button class="delete-btn" @click.stop="confirmDeleteCategory(category)">
                <i class="fas fa-trash"></i>
              </button>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 技能列表 -->
      <div class="skills-section" v-if="selectedCategory">
        <div class="section-header">
          <h3>{{ selectedCategory.name }} 的技能列表</h3>
          <button class="add-btn" @click="showAddSkill = true">
            <i class="fas fa-plus"></i> 添加技能
          </button>
        </div>
        
        <div class="skills-list">
          <div v-for="skill in selectedCategory.skills" :key="skill.id" class="skill-item">
            <div class="skill-header">
              <h4>{{ skill.name }}</h4>
              <div class="skill-actions">
                <button class="edit-btn" @click="editSkill(skill)">
                  <i class="fas fa-edit"></i>
                </button>
                <button class="delete-btn" @click="confirmDeleteSkill(skill)">
                  <i class="fas fa-trash"></i>
                </button>
              </div>
            </div>
            
            <div class="skill-level">
              <div class="level-bar">
                <div class="level-fill" :style="{ width: `${skill.level}%` }"></div>
              </div>
              <span>{{ skill.level }}%</span>
            </div>
            
            <p class="skill-description">{{ skill.description }}</p>
            
            <div class="skill-tags">
              <span v-for="(tag, index) in skill.tags" :key="index" class="tag">
                {{ tag }}
              </span>
            </div>
          </div>
        </div>
      </div>
      
      <div class="no-selection" v-else>
        <p>请选择一个技能分类或添加新分类</p>
      </div>
    </div>
    
    <!-- 添加/编辑分类的模态框 -->
    <div class="modal" v-if="showAddCategory || editingCategory">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ editingCategory ? '编辑分类' : '添加分类' }}</h3>
          <button class="close-btn" @click="closeModals">
            <i class="fas fa-times"></i>
          </button>
        </div>
        
        <form @submit.prevent="saveCategory" class="modal-form">
          <div class="form-group">
            <label for="category-name">分类名称</label>
            <input 
              type="text" 
              id="category-name" 
              v-model="categoryForm.name" 
              required
            >
          </div>
          
          <div class="form-group">
            <label for="category-description">分类描述</label>
            <textarea 
              id="category-description" 
              v-model="categoryForm.description" 
              rows="3"
            ></textarea>
          </div>
          
          <div class="form-group">
            <label for="category-icon">图标类名</label>
            <input 
              type="text" 
              id="category-icon" 
              v-model="categoryForm.icon" 
              placeholder="例如: fas fa-code"
            >
            <div class="icon-preview" v-if="categoryForm.icon">
              <i :class="categoryForm.icon"></i> 预览
            </div>
          </div>
          
          <div class="form-actions">
            <button type="button" class="cancel-btn" @click="closeModals">取消</button>
            <button type="submit" class="save-btn" :disabled="saving">
              {{ saving ? '保存中...' : '保存' }}
            </button>
          </div>
        </form>
      </div>
    </div>
    
    <!-- 添加/编辑技能的模态框 -->
    <div class="modal" v-if="showAddSkill || editingSkill">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ editingSkill ? '编辑技能' : '添加技能' }}</h3>
          <button class="close-btn" @click="closeModals">
            <i class="fas fa-times"></i>
          </button>
        </div>
        
        <form @submit.prevent="saveSkill" class="modal-form">
          <div class="form-group">
            <label for="skill-name">技能名称</label>
            <input 
              type="text" 
              id="skill-name" 
              v-model="skillForm.name" 
              required
            >
          </div>
          
          <div class="form-group">
            <label for="skill-level">掌握程度 ({{ skillForm.level }}%)</label>
            <input 
              type="range" 
              id="skill-level" 
              v-model.number="skillForm.level" 
              min="0" 
              max="100" 
              step="5"
            >
          </div>
          
          <div class="form-group">
            <label for="skill-description">技能描述</label>
            <textarea 
              id="skill-description" 
              v-model="skillForm.description" 
              rows="3"
            ></textarea>
          </div>
          
          <div class="form-group">
            <label for="skill-tags">技能标签 (用逗号分隔)</label>
            <input 
              type="text" 
              id="skill-tags" 
              v-model="tagsInput" 
              placeholder="例如: 前端,React,JavaScript"
            >
            <div class="tags-preview">
              <span v-for="(tag, index) in parsedTags" :key="index" class="tag">
                {{ tag }}
              </span>
            </div>
          </div>
          
          <div class="form-actions">
            <button type="button" class="cancel-btn" @click="closeModals">取消</button>
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
          <button class="close-btn" @click="closeModals">
            <i class="fas fa-times"></i>
          </button>
        </div>
        
        <div class="confirm-message">
          <i class="fas fa-exclamation-triangle"></i>
          <p>{{ deleteConfirmMessage }}</p>
        </div>
        
        <div class="form-actions">
          <button type="button" class="cancel-btn" @click="closeModals">取消</button>
          <button type="button" class="delete-confirm-btn" @click="confirmDelete" :disabled="saving">
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

// 分类和技能数据
const categories = ref([]);
const selectedCategory = ref(null);

// 模态框状态
const showAddCategory = ref(false);
const editingCategory = ref(null);
const showAddSkill = ref(false);
const editingSkill = ref(null);
const showDeleteConfirm = ref(false);
const deleteConfirmMessage = ref('');
const deleteCallback = ref(null);

// 表单数据
const categoryForm = reactive({
  id: null,
  name: '',
  description: '',
  icon: ''
});

const skillForm = reactive({
  id: null,
  categoryId: null, // 前端使用驼峰命名
  name: '',
  level: 80,
  description: '',
  tags: []
});

const tagsInput = ref('');

const parsedTags = computed(() => {
  return tagsInput.value
    .split(',')
    .map(tag => tag.trim())
    .filter(tag => tag !== '');
});

// 获取所有技能分类和技能
const fetchSkills = async () => {
  loading.value = true;
  error.value = null;
  
  try {
    const token = localStorage.getItem('token');
    const response = await axios.get(`${API_URL}/admin/skills`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    if (response.data.success) {
      categories.value = response.data.data;
      // 如果之前选择了分类，重新选择相同ID的分类
      if (selectedCategory.value) {
        const prevId = selectedCategory.value.id;
        const category = categories.value.find(c => c.id === prevId);
        if (category) {
          selectedCategory.value = category;
        } else {
          selectedCategory.value = null;
        }
      }
    } else {
      error.value = response.data.message || '获取技能数据失败';
    }
  } catch (err) {
    console.error('获取技能数据出错:', err);
    error.value = '获取技能数据时发生错误，请稍后再试';
  } finally {
    loading.value = false;
  }
};

// 选择分类
const selectCategory = (category) => {
  selectedCategory.value = category;
};

// 编辑分类
const editCategory = (category) => {
  editingCategory.value = category;
  categoryForm.id = category.id;
  categoryForm.name = category.name;
  categoryForm.description = category.description;
  categoryForm.icon = category.icon;
};

// 确认删除分类
const confirmDeleteCategory = (category) => {
  showDeleteConfirm.value = true;
  deleteConfirmMessage.value = `确定要删除"${category.name}"分类吗？这将同时删除该分类下的所有技能！`;
  deleteCallback.value = () => deleteCategory(category.id);
};

// 删除分类
const deleteCategory = async (categoryId) => {
  saving.value = true;
  
  try {
    const token = localStorage.getItem('token');
    const response = await axios.delete(`${API_URL}/admin/skill-categories/${categoryId}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    if (response.data.success) {
      if (selectedCategory.value && selectedCategory.value.id === categoryId) {
        selectedCategory.value = null;
      }
      await fetchSkills();
      showSuccessMessage('分类删除成功');
    } else {
      error.value = response.data.message || '删除分类失败';
    }
  } catch (err) {
    console.error('删除分类出错:', err);
    error.value = '删除分类时发生错误，请稍后再试';
  } finally {
    saving.value = false;
    closeModals();
  }
};

// 保存分类
const saveCategory = async () => {
  saving.value = true;
  
  try {
    const token = localStorage.getItem('token');
    const isEditing = !!editingCategory.value;
    const url = isEditing 
      ? `${API_URL}/admin/skill-categories/${categoryForm.id}` 
      : `${API_URL}/admin/skill-categories`;
    const method = isEditing ? 'put' : 'post';
    
    const response = await axios({
      method,
      url,
      data: categoryForm,
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });
    
    if (response.data.success) {
      await fetchSkills();
      showSuccessMessage(isEditing ? '分类更新成功' : '分类添加成功');
    } else {
      error.value = response.data.message || (isEditing ? '更新分类失败' : '添加分类失败');
    }
  } catch (err) {
    console.error('保存分类出错:', err);
    error.value = '保存分类时发生错误，请稍后再试';
  } finally {
    saving.value = false;
    closeModals();
  }
};

// 编辑技能
const editSkill = (skill) => {
  editingSkill.value = skill;
  skillForm.id = skill.id;
  skillForm.categoryId = skill.category_id || selectedCategory.value.id;
  skillForm.name = skill.name;
  skillForm.level = skill.level;
  skillForm.description = skill.description;
  skillForm.tags = [...skill.tags];
  tagsInput.value = skill.tags.join(', ');
};

// 确认删除技能
const confirmDeleteSkill = (skill) => {
  showDeleteConfirm.value = true;
  deleteConfirmMessage.value = `确定要删除"${skill.name}"技能吗？`;
  deleteCallback.value = () => deleteSkill(skill.id);
};

// 删除技能
const deleteSkill = async (skillId) => {
  saving.value = true;
  
  try {
    const token = localStorage.getItem('token');
    const response = await axios.delete(`${API_URL}/admin/skills/${skillId}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    if (response.data.success) {
      await fetchSkills();
      showSuccessMessage('技能删除成功');
    } else {
      error.value = response.data.message || '删除技能失败';
    }
  } catch (err) {
    console.error('删除技能出错:', err);
    error.value = '删除技能时发生错误，请稍后再试';
  } finally {
    saving.value = false;
    closeModals();
  }
};

// 保存技能
const saveSkill = async () => {
  saving.value = true;
  
  // 更新技能标签
  skillForm.tags = parsedTags.value;
  
  try {
    const token = localStorage.getItem('token');
    const isEditing = !!editingSkill.value;
    const url = isEditing 
      ? `${API_URL}/admin/skills/${skillForm.id}` 
      : `${API_URL}/admin/skills`;
    const method = isEditing ? 'put' : 'post';
    
    // 确保设置了分类ID
    if (!isEditing) {
      skillForm.categoryId = selectedCategory.value.id;
    } else if (!skillForm.categoryId) {
      // 如果编辑现有技能但缺少分类ID，使用原始技能的分类ID
      skillForm.categoryId = editingSkill.value.category_id || selectedCategory.value.id;
    }
    
    // 日志输出便于调试
    console.log('保存技能数据:', {
      isEditing,
      skillId: skillForm.id,
      categoryId: skillForm.categoryId
    });
    
    // 创建要发送的数据对象并映射字段名称
    const requestData = {
      id: skillForm.id,
      category_id: skillForm.categoryId, // 确保使用正确的字段名称
      name: skillForm.name,
      level: skillForm.level,
      description: skillForm.description,
      tags: skillForm.tags
    };
    
    console.log('发送到后端的数据:', requestData);
    
    const response = await axios({
      method,
      url,
      data: requestData,
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });
    
    if (response.data.success) {
      await fetchSkills();
      showSuccessMessage(isEditing ? '技能更新成功' : '技能添加成功');
    } else {
      error.value = response.data.message || (isEditing ? '更新技能失败' : '添加技能失败');
    }
  } catch (err) {
    console.error('保存技能出错:', err);
    error.value = '保存技能时发生错误，请稍后再试';
  } finally {
    saving.value = false;
    closeModals();
  }
};

// 确认删除
const confirmDelete = () => {
  if (deleteCallback.value) {
    deleteCallback.value();
  }
};

// 关闭所有模态框
const closeModals = () => {
  showAddCategory.value = false;
  editingCategory.value = null;
  showAddSkill.value = false;
  editingSkill.value = null;
  showDeleteConfirm.value = false;
  
  // 重置表单
  categoryForm.id = null;
  categoryForm.name = '';
  categoryForm.description = '';
  categoryForm.icon = '';
  
  skillForm.id = null;
  skillForm.categoryId = null;
  skillForm.name = '';
  skillForm.level = 80;
  skillForm.description = '';
  skillForm.tags = [];
  tagsInput.value = '';
};

// 显示成功消息
const showSuccessMessage = (message) => {
  successMessage.value = message;
  saveSuccess.value = true;
  setTimeout(() => {
    saveSuccess.value = false;
  }, 3000);
};

// 页面加载时获取技能数据
onMounted(() => {
  fetchSkills();
});
</script>

<style scoped>
.skills-form-container {
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

.skills-content {
  display: flex;
  gap: 30px;
}

.categories-section {
  flex: 1;
  max-width: 300px;
}

.skills-section {
  flex: 2;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.section-header h3 {
  margin: 0;
  font-size: 1.2rem;
}

.add-btn {
  background-color: var(--primary-color);
  color: white;
  border: none;
  padding: 8px 12px;
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 0.9rem;
}

.add-btn:hover {
  background-color: var(--primary-dark);
}

.categories-list {
  background-color: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.category-item {
  padding: 12px 15px;
  border-bottom: 1px solid #e2e8f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
  transition: all 0.2s;
}

.category-item:last-child {
  border-bottom: none;
}

.category-item:hover {
  background-color: #f1f5f9;
}

.category-item.active {
  background-color: #e2e8f0;
  border-left: 3px solid var(--accent-color);
}

.category-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.category-info i {
  font-size: 1.1rem;
  width: 20px;
  text-align: center;
  color: var(--primary-color);
}

.category-actions {
  display: flex;
  gap: 5px;
  opacity: 0.5;
  transition: opacity 0.2s;
}

.category-item:hover .category-actions {
  opacity: 1;
}

.edit-btn, .delete-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 5px;
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

.skills-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.skill-item {
  background-color: white;
  border-radius: 8px;
  padding: 15px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.skill-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.skill-header h4 {
  margin: 0;
  font-size: 1.1rem;
}

.skill-level {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.level-bar {
  flex: 1;
  height: 8px;
  background-color: #e2e8f0;
  border-radius: 4px;
  overflow: hidden;
}

.level-fill {
  height: 100%;
  background-color: var(--accent-color);
}

.skill-description {
  color: #4b5563;
  font-size: 0.95rem;
  margin-bottom: 10px;
}

.skill-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag {
  background-color: #f1f5f9;
  color: #4b5563;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.85rem;
}

.no-selection {
  flex: 2;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: white;
  border-radius: 8px;
  padding: 40px;
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
  max-width: 500px;
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
  padding: 10px;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  font-size: 1rem;
}

.form-group input[type="range"] {
  padding: 0;
}

.icon-preview {
  margin-top: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
  color: #4b5563;
}

.tags-preview {
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
}

.cancel-btn {
  background-color: #f3f4f6;
  color: #374151;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 600;
}

.save-btn {
  background-color: var(--primary-color);
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 600;
}

.delete-confirm-btn {
  background-color: #e02424;
  color: white;
  border: none;
  padding: 8px 16px;
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
  .skills-content {
    flex-direction: column;
  }
  
  .categories-section {
    max-width: none;
  }
  
  .skills-list {
    grid-template-columns: 1fr;
  }
}
</style> 