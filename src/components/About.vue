<template>
  <section id="about" class="about-section">
    <div class="container">
      <div class="section-title" data-aos="fade-up">
        <h2>关于我</h2>
      </div>
      
      <div class="about-content" v-if="profile">
        <div class="about-image" data-aos="fade-right" data-aos-delay="100">
          <div class="image-container">
            <img :src="profile.avatar || '../assets/img/profile.jpg'" alt="个人照片" />
          </div>
          <div class="tech-stack">
            <div class="tech-item"><i class="fab fa-linux"></i></div>
            <div class="tech-item"><i class="fab fa-docker"></i></div>
            <div class="tech-item"><i class="fas fa-dharmachakra"></i></div>
            <div class="tech-item"><i class="fab fa-aws"></i></div>
            <div class="tech-item"><i class="fab fa-github"></i></div>
          </div>
        </div>
        
        <div class="about-info" data-aos="fade-left" data-aos-delay="200">
          <h3>{{ profile.title }}</h3>
          <p class="intro">
            {{ profile.introduction }}
          </p>
          
          <div class="info-row">
            <div class="info-col">
              <ul class="info-list">
                <li><strong>姓名:</strong> <span>{{ profile.name }}</span></li>
                <li><strong>工作经验:</strong> <span>{{ profile.years_of_exp }}年</span></li>
                <li><strong>学历:</strong> <span>{{ profile.education }}</span></li>
              </ul>
            </div>
            <div class="info-col">
              <ul class="info-list">
                <li><strong>电子邮箱:</strong> <span>{{ profile.email }}</span></li>
                <li><strong>现居城市:</strong> <span>{{ profile.location }}</span></li>
                <li><strong>求职状态:</strong> <span class="highlight">{{ profile.job_status }}</span></li>
              </ul>
            </div>
          </div>
          
          <p class="philosophy">
            我的技术理念：<span class="highlight">"{{ profile.philosophy }}"</span>。我注重系统的稳定性和可靠性，同时不断学习和应用新技术来提高团队效率。
          </p>
          
          <div class="about-btns">
            <a :href="profile.resume_file_url" class="btn btn-primary" v-if="profile.resume_file_url">
              <i class="fas fa-download"></i> 下载简历
            </a>
          </div>
        </div>
      </div>
      
      <div v-else class="loading-container">
        <div class="loading-spinner"></div>
        <p>加载中...</p>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import apiService from '../services/api';

const profile = ref(null);

// 获取个人信息
const fetchProfile = async () => {
  try {
    const response = await apiService.getProfile();
    if (response.success) {
      profile.value = response.data;
    } else {
      console.error('获取个人信息失败:', response.message);
    }
  } catch (error) {
    console.error('获取个人信息出错:', error);
  }
};

onMounted(() => {
  fetchProfile();
});
</script>

<style scoped>
.about-section {
  background-color: #f8f9fa;
  padding: 100px 0;
}

.about-content {
  display: flex;
  align-items: center;
  gap: 50px;
}

.about-image {
  flex: 1;
  position: relative;
}

.image-container {
  border-radius: 10px;
  overflow: hidden;
  box-shadow: 0 15px 30px rgba(0, 0, 0, 0.1);
  position: relative;
}

.image-container img {
  width: 100%;
  display: block;
  transition: transform 0.5s;
}

.image-container:hover img {
  transform: scale(1.03);
}

.tech-stack {
  display: flex;
  justify-content: center;
  gap: 15px;
  margin-top: 20px;
}

.tech-item {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background-color: var(--primary-color);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
  transition: 0.3s;
}

.tech-item:hover {
  background-color: var(--accent-color);
  transform: translateY(-5px);
}

.about-info {
  flex: 1;
}

.about-info h3 {
  font-size: 1.8rem;
  margin-bottom: 20px;
  position: relative;
  padding-bottom: 15px;
  color: var(--primary-color);
}

.about-info h3:after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 50px;
  height: 3px;
  background: var(--accent-color);
}

.intro {
  margin-bottom: 25px;
  line-height: 1.8;
  font-size: 1.05rem;
}

.info-row {
  display: flex;
  gap: 30px;
  margin-bottom: 25px;
}

.info-col {
  flex: 1;
}

.info-list {
  list-style: none;
}

.info-list li {
  padding: 10px 0;
  border-bottom: 1px solid #eee;
}

.info-list li strong {
  color: var(--primary-color);
  font-weight: 600;
  margin-right: 10px;
}

.philosophy {
  line-height: 1.8;
  margin-bottom: 30px;
  font-style: italic;
}

.highlight {
  color: var(--accent-color);
  font-weight: 600;
}

.about-btns {
  margin-top: 30px;
}

.btn {
  padding: 12px 30px;
  border-radius: 50px;
  font-weight: 500;
  text-transform: uppercase;
  font-size: 0.9rem;
  letter-spacing: 1px;
  transition: 0.3s;
  cursor: pointer;
  display: inline-block;
}

.btn-primary {
  background-color: var(--accent-color);
  color: var(--light-color);
}

.btn-primary:hover {
  background-color: var(--secondary-color);
  transform: translateY(-3px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
}

.btn i {
  margin-right: 8px;
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
  .about-content {
    flex-direction: column;
  }
  
  .about-image, .about-info {
    width: 100%;
  }
  
  .about-image {
    margin-bottom: 40px;
  }
  
  .info-row {
    flex-direction: column;
    gap: 0;
  }
}
</style> 