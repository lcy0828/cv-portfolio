<template>
  <section id="certificates" class="certificates-section">
    <div class="container">
      <div class="section-title" data-aos="fade-up">
        <h2>专业证书</h2>
      </div>
      
      <div v-if="certificates.length" class="certificates-container">
        <div 
          v-for="(cert, index) in certificates" 
          :key="cert.id"
          class="certificate-card"
          data-aos="fade-up"
          :data-aos-delay="index * 100"
        >
          <div class="certificate-icon">
            <i :class="cert.icon"></i>
          </div>
          <div class="certificate-body">
            <h3>{{ cert.name }}</h3>
            <div class="certificate-issuer">
              <span>{{ cert.organization }}</span>
              <div class="certificate-date">{{ cert.date }}</div>
            </div>
            <p>{{ cert.description }}</p>
            <a v-if="cert.link" :href="cert.link" class="view-certificate" target="_blank">
              <i class="fas fa-external-link-alt"></i> 查看证书
            </a>
          </div>
        </div>
      </div>
      
      <div v-else class="loading-container">
        <div class="loading-spinner"></div>
        <p>加载证书数据中...</p>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import apiService from '../services/api';

const certificates = ref([]);

// 从API获取证书数据
const fetchCertificates = async () => {
  try {
    const response = await apiService.getCertificates();
    if (response.success) {
      certificates.value = response.data;
    } else {
      console.error('获取证书数据失败:', response.message);
    }
  } catch (error) {
    console.error('获取证书数据出错:', error);
  }
};

onMounted(() => {
  fetchCertificates();
});
</script>

<style scoped>
.certificates-section {
  padding: 100px 0;
  background-color: white;
}

.certificates-container {
  display: flex;
  flex-direction: column;
  gap: 30px;
  max-width: 800px;
  margin: 0 auto;
}

.certificate-card {
  display: flex;
  background-color: white;
  border-radius: 10px;
  box-shadow: 0 5px 25px rgba(0, 0, 0, 0.08);
  overflow: hidden;
  transition: transform 0.3s;
}

.certificate-card:hover {
  transform: translateY(-5px);
}

.certificate-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100px;
  background-color: var(--primary-color);
  color: white;
  font-size: 2rem;
  transition: all 0.3s;
}

.certificate-card:hover .certificate-icon {
  background-color: var(--accent-color);
}

.certificate-body {
  padding: 25px;
  flex: 1;
}

.certificate-body h3 {
  font-size: 1.4rem;
  margin-bottom: 10px;
  color: var(--primary-color);
}

.certificate-issuer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  font-size: 0.95rem;
  color: var(--secondary-color);
  border-bottom: 1px solid #eee;
  padding-bottom: 10px;
}

.certificate-date {
  background-color: #f1f5f9;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 0.85rem;
  color: var(--text-color);
}

.certificate-body p {
  margin-bottom: 20px;
  line-height: 1.7;
  color: var(--text-color);
}

.view-certificate {
  display: inline-flex;
  align-items: center;
  font-size: 0.9rem;
  color: var(--accent-color);
  border: 1px solid var(--accent-color);
  padding: 6px 16px;
  border-radius: 50px;
  transition: all 0.3s;
}

.view-certificate i {
  margin-right: 8px;
}

.view-certificate:hover {
  background-color: var(--accent-color);
  color: white;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 200px;
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

@media (max-width: 768px) {
  .certificate-card {
    flex-direction: column;
  }
  
  .certificate-icon {
    width: 100%;
    height: 80px;
  }
  
  .certificate-issuer {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
}
</style> 