<template>
  <section id="projects" class="projects-section">
    <div class="container">
      <div class="section-title" data-aos="fade-up">
        <h2>项目经验</h2>
      </div>
      
      <div v-if="projects.length">
        <div class="projects-filter" data-aos="fade-up">
          <button 
            class="filter-btn" 
            :class="{ 'active': activeFilter === '*' }"
            @click="setFilter('*')"
          >
            全部
          </button>
          <button 
            v-for="category in uniqueCategories" 
            :key="category"
            class="filter-btn" 
            :class="{ 'active': activeFilter === category }"
            @click="setFilter(category)"
          >
            {{ category }}
          </button>
        </div>
        
        <div class="projects-grid" data-aos="fade-up">
          <div 
            v-for="project in filteredProjects" 
            :key="project.id"
            class="project-card"
            data-aos="zoom-in"
            :data-aos-delay="100"
          >
            <div class="project-image">
              <img :src="project.image" :alt="project.title" />
              <div class="project-links">
                <a v-if="project.demo_link" :href="project.demo_link" class="project-link" target="_blank">
                  <i class="fas fa-external-link-alt"></i>
                </a>
                <a v-if="project.repo_link" :href="project.repo_link" class="project-link" target="_blank">
                  <i class="fab fa-github"></i>
                </a>
              </div>
            </div>
            <div class="project-info">
              <div class="project-category">{{ project.category }}</div>
              <h3 class="project-title">{{ project.title }}</h3>
              <p class="project-description">{{ project.description }}</p>
              
              <div v-if="project.metrics && project.metrics.length" class="project-metrics">
                <div v-for="(metric, index) in project.metrics" :key="index" class="metric">
                  <div class="metric-value">{{ metric.value }}</div>
                  <div class="metric-label">{{ metric.label }}</div>
                </div>
              </div>
              
              <div class="project-keypoints">
                <h4>关键点</h4>
                <ul>
                  <li v-for="(point, index) in project.key_points" :key="index">
                    {{ point }}
                  </li>
                </ul>
              </div>
              
              <div class="project-tech">
                <div v-for="(tech, index) in project.tech_stack" :key="index" class="tech-badge">
                  {{ tech }}
                </div>
              </div>
              
              <button 
                v-if="project.show_architecture" 
                class="show-architecture-btn"
                @click="showArchitecture(project)"
              >
                <i class="fas fa-project-diagram"></i> 查看架构图
              </button>
            </div>
          </div>
        </div>
      </div>
      
      <div v-else class="loading-container">
        <div class="loading-spinner"></div>
        <p>加载项目数据中...</p>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import apiService from '../services/api';

const projects = ref([]);
const activeFilter = ref('*');

// 从API获取项目数据
const fetchProjects = async () => {
  try {
    const response = await apiService.getProjects();
    if (response.success) {
      projects.value = response.data;
    } else {
      console.error('获取项目数据失败:', response.message);
    }
  } catch (error) {
    console.error('获取项目数据出错:', error);
  }
};

// 设置过滤器
const setFilter = (filter) => {
  activeFilter.value = filter;
};

// 获取所有唯一的项目类别
const uniqueCategories = computed(() => {
  const categories = new Set(projects.value.map(project => project.category));
  return Array.from(categories);
});

// 过滤项目
const filteredProjects = computed(() => {
  if (activeFilter.value === '*') {
    return projects.value;
  }
  return projects.value.filter(project => project.category === activeFilter.value);
});

// 显示架构图
const showArchitecture = (project) => {
  // 这里可以实现查看架构图的逻辑，例如弹出模态框等
  console.log('显示项目架构图:', project.title);
};

onMounted(() => {
  fetchProjects();
});
</script>

<style scoped>
.projects-section {
  padding: 100px 0;
  background-color: white;
}

.projects-filter {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  margin-bottom: 40px;
  gap: 10px;
}

.filter-btn {
  background: none;
  border: 2px solid var(--primary-color);
  color: var(--primary-color);
  padding: 8px 20px;
  border-radius: 30px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.3s;
}

.filter-btn:hover, .filter-btn.active {
  background-color: var(--primary-color);
  color: white;
}

.projects-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 30px;
}

.project-card {
  border-radius: 10px;
  overflow: hidden;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s;
  background-color: white;
}

.project-card:hover {
  transform: translateY(-10px);
}

.project-image {
  height: 200px;
  overflow: hidden;
  position: relative;
}

.project-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s;
}

.project-card:hover .project-image img {
  transform: scale(1.1);
}

.project-links {
  position: absolute;
  top: 15px;
  right: 15px;
  display: flex;
  gap: 10px;
}

.project-link {
  width: 40px;
  height: 40px;
  background-color: rgba(255, 255, 255, 0.9);
  color: var(--primary-color);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: 0.3s;
}

.project-link:hover {
  background-color: var(--primary-color);
  color: white;
}

.project-info {
  padding: 25px;
}

.project-category {
  display: inline-block;
  background-color: var(--accent-color);
  color: white;
  padding: 5px 15px;
  border-radius: 20px;
  font-size: 0.8rem;
  margin-bottom: 15px;
}

.project-title {
  font-size: 1.5rem;
  margin-bottom: 15px;
  color: var(--primary-color);
}

.project-description {
  color: var(--text-color);
  margin-bottom: 20px;
  line-height: 1.6;
}

.project-metrics {
  display: flex;
  justify-content: space-between;
  margin: 25px 0;
  text-align: center;
}

.metric {
  flex: 1;
  padding: 0 10px;
}

.metric-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--primary-color);
}

.metric-label {
  font-size: 0.9rem;
  color: var(--text-color);
  margin-top: 5px;
}

.project-keypoints {
  margin: 20px 0;
}

.project-keypoints h4 {
  font-size: 1.1rem;
  margin-bottom: 15px;
  color: var(--primary-color);
}

.project-keypoints ul {
  padding-left: 20px;
}

.project-keypoints li {
  margin-bottom: 8px;
  line-height: 1.5;
}

.project-tech {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 20px;
}

.tech-badge {
  background-color: #f1f5f9;
  color: var(--secondary-color);
  padding: 5px 12px;
  border-radius: 5px;
  font-size: 0.85rem;
}

.show-architecture-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  background-color: var(--primary-color);
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 5px;
  margin-top: 20px;
  cursor: pointer;
  transition: 0.3s;
}

.show-architecture-btn:hover {
  background-color: var(--secondary-color);
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 300px;
}

.loading-spinner {
  border: 4px solid #f3f3f3;
  border-top: 4px solid var(--accent-color);
  border-radius: 50%;
  width: 50px;
  height: 50px;
  animation: spin 2s linear infinite;
  margin-bottom: 20px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

@media (max-width: 992px) {
  .projects-grid {
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  }
}

@media (max-width: 576px) {
  .projects-grid {
    grid-template-columns: 1fr;
  }
}
</style> 