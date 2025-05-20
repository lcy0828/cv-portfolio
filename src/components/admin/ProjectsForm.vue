<template>
  <div class="projects-form-container">
    <div v-if="loading" class="loading">
      <i class="fas fa-spinner fa-spin"></i> 加载中...
    </div>
    <div v-else-if="error" class="error-message">
      <i class="fas fa-exclamation-circle"></i> {{ error }}
    </div>
    <div v-else class="projects-content">
      <div class="actions-bar">
        <button class="add-btn" @click="showAddProject = true">
          <i class="fas fa-plus"></i> 添加项目经验
        </button>
      </div>
      
      <div class="projects-list">
        <div v-for="project in projects" :key="project.id" class="project-card">
          <div class="card-image" v-if="project.image">
            <img :src="project.image" :alt="project.title">
          </div>
          <div class="card-header">
            <div class="project-category" v-if="project.category">{{ project.category }}</div>
            <div class="card-actions">
              <button class="edit-btn" @click="editProject(project)">
                <i class="fas fa-edit"></i>
              </button>
              <button class="delete-btn" @click="confirmDeleteProject(project)">
                <i class="fas fa-trash"></i>
              </button>
            </div>
          </div>
          
          <div class="card-body">
            <h3 class="project-title">{{ project.title }}</h3>
            <p class="project-description">{{ project.description }}</p>
            
            <div class="project-links" v-if="project.demoLink || project.repoLink">
              <a v-if="project.demoLink" :href="project.demoLink" target="_blank" class="demo-link">
                <i class="fas fa-globe"></i> 演示链接
              </a>
              <a v-if="project.repoLink" :href="project.repoLink" target="_blank" class="repo-link">
                <i class="fab fa-github"></i> 代码仓库
              </a>
            </div>
            
            <div class="project-metrics" v-if="project.metrics && project.metrics.length">
              <div v-for="(metric, idx) in project.metrics" :key="idx" class="metric-item">
                <div class="metric-value">{{ metric.value }}</div>
                <div class="metric-label">{{ metric.label }}</div>
              </div>
            </div>
            
            <div class="project-details">
              <div class="details-section" v-if="project.keyPoints && project.keyPoints.length">
                <h4>核心要点</h4>
                <ul>
                  <li v-for="(point, idx) in project.keyPoints" :key="idx">{{ point }}</li>
                </ul>
              </div>
              
              <div class="tech-stack" v-if="project.techStack && project.techStack.length">
                <h4>技术栈</h4>
                <div class="tech-list">
                  <span v-for="(tech, idx) in project.techStack" :key="idx" class="tech-tag">
                    {{ tech }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <div v-if="projects.length === 0" class="no-data">
          暂无项目经验，请添加您的项目经验。
        </div>
      </div>
    </div>
    
    <!-- 项目表单模态框 -->
    <div v-if="showAddProject || editingProject" class="modal-overlay">
      <div class="modal-container">
        <div class="modal-header">
          <h3>{{ editingProject ? '编辑项目' : '添加新项目' }}</h3>
          <button class="close-btn" @click="closeProjectForm">
            <i class="fas fa-times"></i>
          </button>
        </div>
        
        <div class="modal-body">
          <form class="project-form" @submit.prevent="saveProject">
            <div class="form-section">
              <h4>基本信息</h4>
              <div class="form-row">
                <div class="form-group">
                  <label for="title">项目名称 <span class="required">*</span></label>
                  <input type="text" id="title" v-model="projectForm.title" required>
                </div>
                <div class="form-group">
                  <label for="category">项目分类</label>
                  <input type="text" id="category" v-model="projectForm.category" placeholder="例如: 前端开发、后端开发等">
                </div>
              </div>
              
              <div class="form-group">
                <label for="description">项目描述 <span class="required">*</span></label>
                <textarea id="description" v-model="projectForm.description" rows="3" required></textarea>
              </div>
              
              <div class="form-group">
                <label for="image">项目图片URL</label>
                <input type="text" id="image" v-model="projectForm.image" placeholder="图片链接">
              </div>
              
              <div class="form-row">
                <div class="form-group">
                  <label for="demoLink">演示链接</label>
                  <input type="text" id="demoLink" v-model="projectForm.demoLink" placeholder="项目演示地址">
                </div>
                <div class="form-group">
                  <label for="repoLink">代码仓库</label>
                  <input type="text" id="repoLink" v-model="projectForm.repoLink" placeholder="例如: GitHub链接">
                </div>
              </div>
              
              <div class="form-group checkbox-group">
                <label class="checkbox-label">
                  <input type="checkbox" v-model="projectForm.showArchitecture">
                  <span>显示项目架构</span>
                </label>
              </div>
            </div>
            
            <div class="form-section">
              <h4>项目成果指标</h4>
              <div class="dynamic-list-container">
                <div v-for="(metric, index) in projectForm.metrics" :key="index" class="dynamic-item">
                  <div class="dynamic-item-content">
                    <div class="form-row">
                      <div class="form-group">
                        <label :for="'metricValue' + index">指标值</label>
                        <input type="text" :id="'metricValue' + index" v-model="metric.value" placeholder="例如: 200%、5000+">
                      </div>
                      <div class="form-group">
                        <label :for="'metricLabel' + index">指标名称</label>
                        <input type="text" :id="'metricLabel' + index" v-model="metric.label" placeholder="例如: 性能提升、用户增长">
                      </div>
                    </div>
                  </div>
                  <button type="button" class="remove-btn" @click="removeMetric(index)">
                    <i class="fas fa-trash"></i>
                  </button>
                </div>
                
                <button type="button" class="add-item-btn" @click="addMetric">
                  <i class="fas fa-plus"></i> 添加指标
                </button>
              </div>
            </div>
            
            <div class="form-section">
              <h4>项目核心要点</h4>
              <div class="dynamic-list-container">
                <div v-for="(point, index) in projectForm.keyPoints" :key="index" class="dynamic-item">
                  <div class="dynamic-item-content">
                    <div class="form-group">
                      <input type="text" v-model="projectForm.keyPoints[index]" placeholder="项目核心要点">
                    </div>
                  </div>
                  <button type="button" class="remove-btn" @click="removeKeyPoint(index)">
                    <i class="fas fa-trash"></i>
                  </button>
                </div>
                
                <button type="button" class="add-item-btn" @click="addKeyPoint">
                  <i class="fas fa-plus"></i> 添加要点
                </button>
              </div>
            </div>
            
            <div class="form-section">
              <h4>技术栈</h4>
              <div class="dynamic-list-container">
                <div v-for="(tech, index) in projectForm.techStack" :key="index" class="dynamic-item">
                  <div class="dynamic-item-content">
                    <div class="form-group">
                      <input type="text" v-model="projectForm.techStack[index]" placeholder="技术名称">
                    </div>
                  </div>
                  <button type="button" class="remove-btn" @click="removeTech(index)">
                    <i class="fas fa-trash"></i>
                  </button>
                </div>
                
                <button type="button" class="add-item-btn" @click="addTech">
                  <i class="fas fa-plus"></i> 添加技术
                </button>
              </div>
            </div>
            
            <div class="form-group">
              <label for="sortOrder">排序顺序</label>
              <input type="number" id="sortOrder" v-model.number="projectForm.sortOrder" min="0">
            </div>
            
            <div class="form-actions">
              <button type="button" class="cancel-btn" @click="closeProjectForm">取消</button>
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
          <p>您确定要删除以下项目吗？</p>
          <div class="confirm-item" v-if="deletingProject">
            <strong>{{ deletingProject.title }}</strong>
          </div>
          <p class="warning-text">此操作无法撤销！</p>
        </div>
        
        <div class="dialog-actions">
          <button class="cancel-btn" @click="showDeleteConfirm = false">取消</button>
          <button class="delete-btn" @click="deleteProject" :disabled="saving">
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
import { ref, reactive, computed, onMounted } from 'vue';
import axios from 'axios';

const API_URL = 'http://localhost:8080/api';
const loading = ref(false);
const saving = ref(false);
const error = ref(null);
const saveSuccess = ref(false);
const successMessage = ref('');

// 项目数据
const projects = ref([]);

// 模态框状态
const showAddProject = ref(false);
const editingProject = ref(null);
const showDeleteConfirm = ref(false);
const deletingProject = ref(null);

// 项目表单数据
const projectForm = reactive({
  title: '',
  category: '',
  description: '',
  image: '',
  demoLink: '',
  repoLink: '',
  showArchitecture: false,
  metrics: [],
  keyPoints: [],
  techStack: [],
  sortOrder: 0
});

// 重置表单
const resetForm = () => {
  projectForm.title = '';
  projectForm.category = '';
  projectForm.description = '';
  projectForm.image = '';
  projectForm.demoLink = '';
  projectForm.repoLink = '';
  projectForm.showArchitecture = false;
  projectForm.metrics = [];
  projectForm.keyPoints = [];
  projectForm.techStack = [];
  projectForm.sortOrder = projects.value.length > 0 ? Math.max(...projects.value.map(p => p.sortOrder)) + 1 : 0;
};

// 获取项目列表
const fetchProjects = async () => {
  loading.value = true;
  error.value = null;
  
  try {
    const token = localStorage.getItem('token');
    const response = await axios.get(`${API_URL}/admin/projects`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    if (response.data.success) {
      projects.value = response.data.data;
    } else {
      error.value = response.data.message || '获取项目经验失败';
    }
  } catch (err) {
    console.error('获取项目经验出错:', err);
    error.value = '获取项目经验时发生错误，请稍后再试';
  } finally {
    loading.value = false;
  }
};

// 编辑项目
const editProject = (project) => {
  editingProject.value = { ...project };
  
  // 填充表单数据
  projectForm.title = project.title;
  projectForm.category = project.category;
  projectForm.description = project.description;
  projectForm.image = project.image;
  projectForm.demoLink = project.demoLink;
  projectForm.repoLink = project.repoLink;
  projectForm.showArchitecture = project.showArchitecture;
  projectForm.metrics = project.metrics ? [...project.metrics] : [];
  projectForm.keyPoints = project.keyPoints ? [...project.keyPoints] : [];
  projectForm.techStack = project.techStack ? [...project.techStack] : [];
  projectForm.sortOrder = project.sortOrder;
  
  showAddProject.value = true;
};

// 关闭项目表单
const closeProjectForm = () => {
  showAddProject.value = false;
  editingProject.value = null;
  resetForm();
};

// 确认删除项目
const confirmDeleteProject = (project) => {
  showDeleteConfirm.value = true;
  deletingProject.value = project;
};

// 动态表单项操作 - 指标
const addMetric = () => {
  projectForm.metrics.push({ value: '', label: '' });
};

const removeMetric = (index) => {
  projectForm.metrics.splice(index, 1);
};

// 动态表单项操作 - 核心要点
const addKeyPoint = () => {
  projectForm.keyPoints.push('');
};

const removeKeyPoint = (index) => {
  projectForm.keyPoints.splice(index, 1);
};

// 动态表单项操作 - 技术栈
const addTech = () => {
  projectForm.techStack.push('');
};

const removeTech = (index) => {
  projectForm.techStack.splice(index, 1);
};

// 保存项目
const saveProject = async () => {
  saving.value = true;
  error.value = null;
  
  try {
    const token = localStorage.getItem('token');
    const projectData = {
      title: projectForm.title,
      category: projectForm.category,
      description: projectForm.description,
      image: projectForm.image,
      demoLink: projectForm.demoLink,
      repoLink: projectForm.repoLink,
      showArchitecture: projectForm.showArchitecture,
      metrics: projectForm.metrics.filter(m => m.value.trim() || m.label.trim()),
      keyPoints: projectForm.keyPoints.filter(p => p.trim()),
      techStack: projectForm.techStack.filter(t => t.trim()),
      sortOrder: projectForm.sortOrder
    };
    
    let response;
    if (editingProject.value) {
      // 更新项目
      response = await axios.put(`${API_URL}/admin/projects/${editingProject.value.id}`, projectData, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      
      if (response.data.success) {
        // 更新本地数据
        const index = projects.value.findIndex(p => p.id === editingProject.value.id);
        if (index !== -1) {
          projects.value[index] = { ...projects.value[index], ...response.data.data };
        }
        showSuccessMessage('项目更新成功');
      }
    } else {
      // 创建新项目
      response = await axios.post(`${API_URL}/admin/projects`, projectData, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      
      if (response.data.success) {
        // 添加到本地数据
        projects.value.push(response.data.data);
        showSuccessMessage('项目添加成功');
      }
    }
    
    // 关闭表单
    if (response.data.success) {
      closeProjectForm();
    } else {
      error.value = response.data.message || '保存项目失败';
    }
  } catch (err) {
    console.error('保存项目出错:', err);
    error.value = '保存项目时发生错误，请稍后再试';
  } finally {
    saving.value = false;
  }
};

// 删除项目
const deleteProject = async () => {
  if (!deletingProject.value) return;
  
  saving.value = true;
  
  try {
    const token = localStorage.getItem('token');
    const response = await axios.delete(`${API_URL}/admin/projects/${deletingProject.value.id}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    if (response.data.success) {
      // 从本地数据中删除
      const index = projects.value.findIndex(p => p.id === deletingProject.value.id);
      if (index !== -1) {
        projects.value.splice(index, 1);
      }
      
      showSuccessMessage('项目删除成功');
      showDeleteConfirm.value = false;
      deletingProject.value = null;
    } else {
      error.value = response.data.message || '删除项目失败';
    }
  } catch (err) {
    console.error('删除项目出错:', err);
    error.value = '删除项目时发生错误，请稍后再试';
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

// 页面加载时获取项目数据
onMounted(() => {
  fetchProjects();
  resetForm();
});
</script>

<style scoped>
.projects-form-container {
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

.projects-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 30px;
}

.project-card {
  background-color: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
}

.card-image {
  height: 180px;
  overflow: hidden;
}

.card-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}

.project-card:hover .card-image img {
  transform: scale(1.05);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 15px;
  background-color: #f8fafc;
  border-bottom: 1px solid #e2e8f0;
}

.project-category {
  font-size: 0.9rem;
  color: #6b7280;
  background-color: #f1f5f9;
  padding: 4px 10px;
  border-radius: 12px;
}

.card-actions {
  display: flex;
  gap: 8px;
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

.card-body {
  padding: 20px;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.project-title {
  margin: 0 0 10px 0;
  font-size: 1.3rem;
  color: #1f2937;
}

.project-description {
  color: #4b5563;
  margin-bottom: 15px;
  font-size: 0.95rem;
  line-height: 1.5;
}

.project-links {
  display: flex;
  gap: 15px;
  margin-bottom: 15px;
}

.demo-link, .repo-link {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.9rem;
  color: var(--primary-color);
  text-decoration: none;
}

.demo-link:hover, .repo-link:hover {
  text-decoration: underline;
}

.project-metrics {
  display: flex;
  justify-content: space-between;
  margin-bottom: 15px;
  padding: 15px 0;
  border-top: 1px solid #e2e8f0;
  border-bottom: 1px solid #e2e8f0;
}

.metric-item {
  text-align: center;
  flex: 1;
}

.metric-value {
  font-size: 1.4rem;
  font-weight: 600;
  color: var(--primary-color);
}

.metric-label {
  font-size: 0.8rem;
  color: #6b7280;
}

.project-details {
  margin-top: auto;
}

.details-section {
  margin-bottom: 15px;
}

.details-section h4, .tech-stack h4 {
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
  font-size: 0.9rem;
  color: #4b5563;
}

.tech-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tech-tag {
  background-color: #f1f5f9;
  color: #4b5563;
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 0.85rem;
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
  max-width: 800px;
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

.project-form {
  display: flex;
  flex-direction: column;
  gap: 25px;
}

.form-section {
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 20px;
  background-color: #f8fafc;
}

.form-section h4 {
  margin: 0 0 15px 0;
  font-size: 1.1rem;
  color: #1f2937;
  padding-bottom: 10px;
  border-bottom: 1px solid #e2e8f0;
}

.form-row {
  display: flex;
  gap: 15px;
  margin-bottom: 15px;
}

.form-group {
  flex: 1;
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: 600;
  color: #374151;
  font-size: 0.95rem;
}

.form-group input,
.form-group textarea,
.form-group select {
  width: 100%;
  padding: 10px;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  background-color: white;
  font-size: 0.95rem;
}

.form-group input:focus,
.form-group textarea:focus,
.form-group select:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.3);
}

.required {
  color: #e02424;
}

.checkbox-group {
  display: flex;
  align-items: center;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.checkbox-label input[type="checkbox"] {
  width: auto;
}

.dynamic-list-container {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.dynamic-item {
  display: flex;
  gap: 10px;
  align-items: center;
  background-color: white;
  padding: 12px;
  border-radius: 4px;
  border: 1px solid #e2e8f0;
}

.dynamic-item-content {
  flex: 1;
}

.remove-btn {
  background-color: #f1f5f9;
  border: none;
  width: 32px;
  height: 32px;
  border-radius: 4px;
  cursor: pointer;
  color: #6b7280;
  display: flex;
  align-items: center;
  justify-content: center;
}

.remove-btn:hover {
  background-color: #fee2e2;
  color: #e02424;
}

.add-item-btn {
  background-color: #f1f5f9;
  border: 1px dashed #d1d5db;
  padding: 8px 12px;
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: #4b5563;
  font-size: 0.9rem;
}

.add-item-btn:hover {
  background-color: #e5e7eb;
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
  .projects-list {
    grid-template-columns: 1fr;
  }
  
  .form-row {
    flex-direction: column;
    gap: 0;
  }
}
</style> 