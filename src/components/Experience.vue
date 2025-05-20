<template>
  <section id="experience" class="experience-section">
    <div class="container">
      <div class="section-title" data-aos="fade-up">
        <h2>工作经历</h2>
      </div>
      
      <div v-if="experiences.length" class="experience-timeline">
        <div
          v-for="(exp, index) in experiences"
          :key="exp.id"
          class="timeline-item"
          :class="{ 'left': index % 2 === 0, 'right': index % 2 === 1 }"
          data-aos="fade-up"
          :data-aos-delay="index * 100"
        >
          <div class="timeline-badge" :style="{ backgroundColor: exp.color }">
            <i :class="exp.icon"></i>
          </div>
          
          <div class="timeline-panel">
            <div class="timeline-heading">
              <span class="timeline-date">{{ exp.period }}</span>
              <h3 class="timeline-title">{{ exp.title }}</h3>
              <h4 class="timeline-company">{{ exp.company }}</h4>
              <span class="timeline-location">
                <i class="fas fa-map-marker-alt"></i> {{ exp.location }}
              </span>
            </div>
            
            <div class="timeline-body">
              <div class="timeline-section">
                <h5><i class="fas fa-tasks"></i> 工作职责</h5>
                <ul>
                  <li v-for="(item, i) in exp.responsibilities" :key="i">{{ item }}</li>
                </ul>
              </div>
              
              <div class="timeline-section">
                <h5><i class="fas fa-trophy"></i> 主要成就</h5>
                <ul>
                  <li v-for="(item, i) in exp.achievements" :key="i">{{ item }}</li>
                </ul>
              </div>
              
              <div class="timeline-tech">
                <h5><i class="fas fa-code"></i> 使用技术</h5>
                <div class="tech-tags">
                  <span v-for="(tech, i) in exp.technologies" :key="i" class="tech-tag">
                    {{ tech }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <div v-else class="loading-container">
        <div class="loading-spinner"></div>
        <p>加载工作经历数据中...</p>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import apiService from '../services/api';

// 工作经历数据
const experiences = ref([]);

// 从API获取工作经历
const fetchExperiences = async () => {
  try {
    const response = await apiService.getExperiences();
    if (response.success) {
      experiences.value = response.data;
    } else {
      console.error('获取工作经历失败:', response.message);
    }
  } catch (error) {
    console.error('获取工作经历出错:', error);
  }
};

onMounted(() => {
  fetchExperiences();
});
</script>

<style scoped>
.experience-section {
  padding: 100px 0;
  background-color: #f8f9fa;
  position: relative;
}

.experience-timeline {
  position: relative;
  padding: 40px 0;
}

.experience-timeline:before {
  content: '';
  position: absolute;
  top: 0;
  bottom: 0;
  width: 4px;
  background: var(--primary-color);
  left: 50%;
  margin-left: -2px;
  z-index: 1;
}

.timeline-item {
  position: relative;
  margin-bottom: 50px;
  width: 50%;
  box-sizing: border-box;
}

.timeline-item.left {
  left: 0;
  padding-right: 50px;
}

.timeline-item.right {
  left: 50%;
  padding-left: 50px;
}

.timeline-badge {
  width: 50px;
  height: 50px;
  line-height: 50px;
  font-size: 1.4em;
  text-align: center;
  position: absolute;
  top: 16px;
  border-radius: 50%;
  background: var(--primary-color);
  color: #fff;
  z-index: 10;
}

.timeline-item.left .timeline-badge {
  right: -25px;
}

.timeline-item.right .timeline-badge {
  left: -25px;
}

.timeline-panel {
  position: relative;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.08);
  padding: 25px;
}

.timeline-panel:before {
  content: '';
  position: absolute;
  top: 26px;
  border: 15px solid transparent;
}

.timeline-item.left .timeline-panel:before {
  right: -30px;
  border-left-color: #fff;
}

.timeline-item.right .timeline-panel:before {
  left: -30px;
  border-right-color: #fff;
}

.timeline-heading {
  margin-bottom: 20px;
}

.timeline-date {
  display: inline-block;
  padding: 6px 12px;
  background: var(--primary-color);
  color: white;
  border-radius: 20px;
  font-size: 0.85rem;
  margin-bottom: 15px;
}

.timeline-title {
  color: var(--primary-color);
  margin-bottom: 10px;
  font-size: 1.4rem;
}

.timeline-company {
  color: var(--secondary-color);
  margin-bottom: 8px;
  font-size: 1.1rem;
}

.timeline-location {
  display: block;
  color: #6c757d;
  font-size: 0.9rem;
  margin-bottom: 20px;
}

.timeline-location i {
  margin-right: 5px;
}

.timeline-section {
  margin-bottom: 20px;
}

.timeline-section h5 {
  color: var(--primary-color);
  font-size: 1rem;
  margin-bottom: 10px;
  border-bottom: 1px solid #eee;
  padding-bottom: 8px;
}

.timeline-section h5 i {
  margin-right: 8px;
}

.timeline-section ul {
  padding-left: 20px;
}

.timeline-section ul li {
  margin-bottom: 8px;
  line-height: 1.5;
}

.timeline-tech {
  margin-top: 20px;
}

.tech-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 10px;
}

.tech-tag {
  padding: 5px 10px;
  background-color: #e9f5ff;
  color: var(--secondary-color);
  border-radius: 20px;
  font-size: 0.85rem;
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
  .experience-timeline:before {
    left: 40px;
  }
  
  .timeline-item {
    width: 100%;
    padding-left: 90px;
    padding-right: 0;
  }
  
  .timeline-item.left,
  .timeline-item.right {
    left: 0;
    padding-left: 90px;
    padding-right: 0;
  }
  
  .timeline-item.left .timeline-badge,
  .timeline-item.right .timeline-badge {
    left: 15px;
  }
  
  .timeline-item.left .timeline-panel:before,
  .timeline-item.right .timeline-panel:before {
    left: -30px;
    border-right-color: #fff;
    border-left-color: transparent;
  }
}
</style> 