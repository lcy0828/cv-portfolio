<template>
  <section id="skills" class="skills-section">
    <div class="container">
      <div class="section-title" data-aos="fade-up">
        <h2>专业技能</h2>
      </div>
      
      <div v-if="skillCategories.length" class="skills-container">
        <div class="skills-tabs" data-aos="fade-up">
          <div 
            v-for="(category, index) in skillCategories" 
            :key="category.id"
            class="skill-tab"
            :class="{ 'active': activeCategory === index }"
            @click="setActiveCategory(index)"
          >
            <i :class="category.icon"></i>
            <span>{{ category.name }}</span>
          </div>
        </div>
        
        <div class="skills-content" data-aos="fade-up" data-aos-delay="100">
          <div class="skill-category-title">
            <h3>{{ skillCategories[activeCategory].name }}</h3>
            <p>{{ skillCategories[activeCategory].description }}</p>
          </div>
          
          <div class="skills-grid">
            <div 
              v-for="(skill, index) in skillCategories[activeCategory].skills" 
              :key="skill.id"
              class="skill-item"
              data-aos="zoom-in"
              :data-aos-delay="100 + (index * 50)"
            >
              <div class="skill-info">
                <h4>{{ skill.name }}</h4>
                <div class="skill-level">
                  <div 
                    class="skill-progress"
                    :style="{ width: skill.level + '%' }"
                  ></div>
                </div>
                <span class="skill-percentage">{{ skill.level }}%</span>
              </div>
              <div class="skill-detail">
                <p>{{ skill.description }}</p>
                <div v-if="skill.tags && skill.tags.length" class="skill-tags">
                  <span v-for="(tag, tagIndex) in skill.tags" :key="tagIndex" class="skill-tag">
                    {{ tag }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <div class="skills-diagram" data-aos="fade-left" data-aos-delay="200">
          <div class="diagram-title">技能雷达图</div>
          <div class="radar-chart">
            <!-- 这里可以使用canvas或其他图表库绘制雷达图 -->
            <div class="radar-placeholder">
              <div class="radar-point p1" style="top: 30%; left: 50%;">K8s</div>
              <div class="radar-point p2" style="top: 50%; left: 85%;">Docker</div>
              <div class="radar-point p3" style="top: 80%; left: 70%;">Linux</div>
              <div class="radar-point p4" style="top: 80%; left: 30%;">CI/CD</div>
              <div class="radar-point p5" style="top: 50%; left: 15%;">云平台</div>
              <div class="radar-circle c1"></div>
              <div class="radar-circle c2"></div>
              <div class="radar-circle c3"></div>
              <div class="radar-line l1"></div>
              <div class="radar-line l2"></div>
              <div class="radar-line l3"></div>
              <div class="radar-line l4"></div>
              <div class="radar-line l5"></div>
              <div class="radar-area"></div>
            </div>
          </div>
        </div>
      </div>
      
      <div v-else class="loading-container">
        <div class="loading-spinner"></div>
        <p>加载技能数据中...</p>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import apiService from '../services/api';

const activeCategory = ref(0);
const skillCategories = ref([]);

const setActiveCategory = (index) => {
  activeCategory.value = index;
};

// 从API获取技能数据
const fetchSkills = async () => {
  try {
    const response = await apiService.getSkills();
    if (response.success) {
      skillCategories.value = response.data;
    } else {
      console.error('获取技能数据失败:', response.message);
    }
  } catch (error) {
    console.error('获取技能数据出错:', error);
  }
};

onMounted(() => {
  fetchSkills();
});
</script>

<style scoped>
.skills-section {
  background-color: white;
  padding: 100px 0;
}

.skills-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 40px;
}

.skills-tabs {
  grid-column: 1 / 3;
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: 15px;
  margin-bottom: 40px;
}

.skill-tab {
  padding: 15px 25px;
  background-color: #f1f5f9;
  border-radius: 50px;
  cursor: pointer;
  transition: 0.3s;
  display: flex;
  align-items: center;
  gap: 8px;
}

.skill-tab i {
  font-size: 1.1rem;
  color: var(--primary-color);
}

.skill-tab.active {
  background-color: var(--primary-color);
  color: white;
}

.skill-tab.active i {
  color: white;
}

.skill-tab:hover:not(.active) {
  background-color: #e2e8f0;
  transform: translateY(-3px);
}

.skills-content {
  grid-column: 1 / 2;
}

.skill-category-title {
  margin-bottom: 30px;
}

.skill-category-title h3 {
  font-size: 1.8rem;
  color: var(--primary-color);
  margin-bottom: 10px;
}

.skill-category-title p {
  color: var(--text-color);
  font-size: 1rem;
}

.skills-grid {
  display: grid;
  gap: 25px;
}

.skill-item {
  background-color: #f8f9fa;
  border-radius: 10px;
  padding: 20px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.05);
  transition: 0.3s;
}

.skill-item:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
}

.skill-info {
  margin-bottom: 15px;
}

.skill-info h4 {
  font-size: 1.1rem;
  margin-bottom: 10px;
  color: var(--primary-color);
}

.skill-level {
  height: 8px;
  background-color: #e2e8f0;
  border-radius: 10px;
  position: relative;
  margin-bottom: 5px;
  overflow: hidden;
}

.skill-progress {
  position: absolute;
  top: 0;
  left: 0;
  height: 100%;
  border-radius: 10px;
  background: linear-gradient(90deg, var(--secondary-color) 0%, var(--accent-color) 100%);
  transition: width 1s ease-in-out;
  width: 0;
}

.skill-percentage {
  font-size: 0.8rem;
  color: var(--primary-color);
  font-weight: 600;
}

.skill-detail p {
  margin-bottom: 10px;
  font-size: 0.9rem;
  line-height: 1.6;
}

.skill-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.skill-tag {
  font-size: 0.75rem;
  padding: 4px 10px;
  border-radius: 50px;
  background-color: rgba(59, 130, 246, 0.1);
  color: var(--secondary-color);
}

.skills-diagram {
  grid-column: 2 / 3;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.diagram-title {
  font-size: 1.4rem;
  color: var(--primary-color);
  margin-bottom: 20px;
  text-align: center;
}

.radar-chart {
  width: 100%;
  max-width: 400px;
  height: 400px;
  position: relative;
}

.radar-placeholder {
  width: 100%;
  height: 100%;
  position: relative;
  background-color: rgba(59, 130, 246, 0.05);
  border-radius: 50%;
}

.radar-circle {
  position: absolute;
  border-radius: 50%;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  border: 1px dashed var(--secondary-color);
}

.radar-circle.c1 {
  width: 90%;
  height: 90%;
  opacity: 0.2;
}

.radar-circle.c2 {
  width: 60%;
  height: 60%;
  opacity: 0.4;
}

.radar-circle.c3 {
  width: 30%;
  height: 30%;
  opacity: 0.6;
}

.radar-point {
  position: absolute;
  transform: translate(-50%, -50%);
  font-size: 0.9rem;
  color: var(--primary-color);
  font-weight: 600;
  z-index: 2;
}

.radar-line {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 50%;
  height: 1px;
  background-color: rgba(30, 58, 138, 0.1);
  transform-origin: left center;
}

.radar-line.l1 {
  transform: rotate(0deg);
}

.radar-line.l2 {
  transform: rotate(72deg);
}

.radar-line.l3 {
  transform: rotate(144deg);
}

.radar-line.l4 {
  transform: rotate(216deg);
}

.radar-line.l5 {
  transform: rotate(288deg);
}

.radar-area {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 90%;
  height: 90%;
  transform: translate(-50%, -50%);
  clip-path: polygon(
    50% 10%, 
    85% 35%, 
    75% 80%, 
    25% 80%, 
    15% 35%
  );
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.2) 0%, rgba(6, 182, 212, 0.2) 100%);
  border: 2px solid var(--accent-color);
  border-radius: 2px;
  z-index: 1;
}

@media (max-width: 992px) {
  .skills-container {
    grid-template-columns: 1fr;
  }
  
  .skills-tabs,
  .skills-content,
  .skills-diagram {
    grid-column: 1;
  }
  
  .skills-diagram {
    margin-top: 40px;
  }
}

@media (max-width: 576px) {
  .skill-tab {
    padding: 10px 20px;
    font-size: 0.9rem;
  }
  
  .skills-grid {
    gap: 15px;
  }
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
</style> 