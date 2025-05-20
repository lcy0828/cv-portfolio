<template>
  <section id="contact" class="contact-section">
    <div class="container">
      <div class="section-title" data-aos="fade-up">
        <h2>联系方式</h2>
      </div>
      
      <div class="contact-container">
        <div class="contact-info" data-aos="fade-right">
          <div class="info-title">
            <h3>与我联系</h3>
            <p>如果您对我的经历和技能感兴趣，或者有任何问题，欢迎随时与我联系。</p>
          </div>
          
          <div class="info-items">
            <div class="info-item">
              <div class="info-icon">
                <i class="fas fa-map-marker-alt"></i>
              </div>
              <div class="info-text">
                <h4>地址</h4>
                <p>北京市朝阳区科技园区88号</p>
              </div>
            </div>
            
            <div class="info-item">
              <div class="info-icon">
                <i class="fas fa-envelope"></i>
              </div>
              <div class="info-text">
                <h4>电子邮箱</h4>
                <p><a href="mailto:zhangsan@example.com">zhangsan@example.com</a></p>
              </div>
            </div>
            
            <div class="info-item">
              <div class="info-icon">
                <i class="fas fa-phone-alt"></i>
              </div>
              <div class="info-text">
                <h4>电话</h4>
                <p><a href="tel:+8613812345678">138-1234-5678</a></p>
              </div>
            </div>
          </div>
          
          <div class="social-links">
            <a href="#" class="social-link">
              <i class="fab fa-github"></i>
            </a>
            <a href="#" class="social-link">
              <i class="fab fa-linkedin-in"></i>
            </a>
            <a href="#" class="social-link">
              <i class="fab fa-weixin"></i>
            </a>
            <a href="#" class="social-link">
              <i class="fab fa-zhihu"></i>
            </a>
          </div>
        </div>
        
        <div class="contact-form" data-aos="fade-left">
          <form @submit.prevent="submitForm">
            <div class="form-group">
              <label for="name">姓名</label>
              <input type="text" id="name" v-model="formData.name" required>
            </div>
            
            <div class="form-group">
              <label for="email">电子邮箱</label>
              <input type="email" id="email" v-model="formData.email" required>
            </div>
            
            <div class="form-group">
              <label for="subject">主题</label>
              <input type="text" id="subject" v-model="formData.subject" required>
            </div>
            
            <div class="form-group">
              <label for="message">留言</label>
              <textarea id="message" v-model="formData.message" rows="5" required></textarea>
            </div>
            
            <button type="submit" class="submit-btn" :disabled="formSubmitting">
              <span v-if="!formSubmitting">发送消息</span>
              <span v-else><i class="fas fa-spinner fa-spin"></i> 发送中...</span>
            </button>
            
            <div v-if="formSubmitted" class="form-success">
              <i class="fas fa-check-circle"></i> 消息已发送，感谢您的留言！
            </div>
          </form>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref, reactive } from 'vue';

const formData = reactive({
  name: '',
  email: '',
  subject: '',
  message: ''
});

const formSubmitting = ref(false);
const formSubmitted = ref(false);

const submitForm = () => {
  formSubmitting.value = true;
  
  // 模拟表单提交
  setTimeout(() => {
    formSubmitting.value = false;
    formSubmitted.value = true;
    
    // 重置表单
    formData.name = '';
    formData.email = '';
    formData.subject = '';
    formData.message = '';
    
    // 3秒后隐藏成功消息
    setTimeout(() => {
      formSubmitted.value = false;
    }, 3000);
  }, 1500);
};
</script>

<style scoped>
.contact-section {
  background-color: white;
  padding: 100px 0;
}

.contact-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 50px;
}

.contact-info {
  padding: 30px;
  background-color: #f8f9fa;
  border-radius: 10px;
  box-shadow: 0 5px 20px rgba(0, 0, 0, 0.05);
  height: 100%;
}

.info-title {
  margin-bottom: 30px;
}

.info-title h3 {
  font-size: 1.8rem;
  margin-bottom: 15px;
  color: var(--primary-color);
}

.info-title p {
  line-height: 1.6;
  font-size: 1rem;
}

.info-items {
  margin-bottom: 30px;
}

.info-item {
  display: flex;
  align-items: flex-start;
  margin-bottom: 25px;
}

.info-icon {
  flex-shrink: 0;
  width: 50px;
  height: 50px;
  background-color: rgba(59, 130, 246, 0.1);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20px;
  color: var(--secondary-color);
  font-size: 1.2rem;
  transition: all 0.3s;
}

.info-item:hover .info-icon {
  background-color: var(--secondary-color);
  color: white;
  transform: scale(1.1);
}

.info-text h4 {
  font-size: 1.1rem;
  margin-bottom: 5px;
  color: var(--primary-color);
}

.info-text p {
  font-size: 0.95rem;
  color: var(--text-color);
}

.info-text a {
  color: var(--secondary-color);
  transition: color 0.3s;
}

.info-text a:hover {
  color: var(--accent-color);
}

.social-links {
  display: flex;
  gap: 15px;
}

.social-link {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: var(--primary-color);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s;
}

.social-link:hover {
  background-color: var(--accent-color);
  transform: translateY(-3px);
}

.contact-form {
  padding: 30px;
  background-color: white;
  border-radius: 10px;
  box-shadow: 0 5px 20px rgba(0, 0, 0, 0.05);
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: var(--primary-color);
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 12px 15px;
  border: 1px solid #e2e8f0;
  border-radius: 5px;
  background-color: #f8f9fa;
  font-size: 0.95rem;
  transition: all 0.3s;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: var(--secondary-color);
  background-color: white;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.submit-btn {
  display: inline-block;
  padding: 12px 30px;
  background-color: var(--accent-color);
  color: white;
  border: none;
  border-radius: 5px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
}

.submit-btn:hover:not(:disabled) {
  background-color: var(--secondary-color);
  transform: translateY(-3px);
}

.submit-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.form-success {
  margin-top: 20px;
  padding: 15px;
  background-color: rgba(16, 185, 129, 0.1);
  border-radius: 5px;
  color: #10b981;
  display: flex;
  align-items: center;
  gap: 10px;
}

.form-success i {
  font-size: 1.2rem;
}

@media (max-width: 992px) {
  .contact-container {
    grid-template-columns: 1fr;
  }
  
  .contact-info {
    margin-bottom: 30px;
  }
}

@media (max-width: 576px) {
  .info-item {
    flex-direction: column;
  }
  
  .info-icon {
    margin-bottom: 15px;
    margin-right: 0;
  }
}
</style> 