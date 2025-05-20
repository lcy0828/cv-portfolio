<template>
  <footer class="footer-section">
    <div class="container">
      <div class="footer-content">
        <div class="footer-top">
          <div class="footer-logo">
            <h2>张三</h2>
            <p>云原生运维工程师 / K8s专家 / 系统架构师</p>
          </div>
          
          <div class="footer-nav">
            <ul>
              <li><a href="#about">关于我</a></li>
              <li><a href="#skills">专业技能</a></li>
              <li><a href="#experience">工作经历</a></li>
              <li><a href="#projects">项目经验</a></li>
              <li><a href="#certificates">证书认证</a></li>
              <li><a href="#contact">联系方式</a></li>
            </ul>
          </div>
        </div>
        
        <div class="footer-divider"></div>
        
        <div class="footer-bottom">
          <div class="copyright">
            &copy; {{ currentYear }} 张三的个人简历. 保留所有权利.
          </div>
          
          <div class="back-to-top">
            <a href="#hero" @click="scrollToTop">
              <i class="fas fa-arrow-up"></i> 回到顶部
            </a>
          </div>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup>
import { ref, onMounted } from 'vue';

const currentYear = ref(new Date().getFullYear());

const scrollToTop = (e) => {
  e.preventDefault();
  window.scrollTo({
    top: 0,
    behavior: 'smooth'
  });
};

onMounted(() => {
  // 检查URL中的锚点，确保页面加载完成后能正确滚动
  const hash = window.location.hash;
  if (hash) {
    setTimeout(() => {
      const element = document.querySelector(hash);
      if (element) {
        element.scrollIntoView({ behavior: 'smooth' });
      }
    }, 500);
  }
});
</script>

<style scoped>
.footer-section {
  background-color: var(--primary-color);
  color: white;
  padding: 60px 0 30px;
}

.footer-content {
  display: flex;
  flex-direction: column;
}

.footer-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.footer-logo h2 {
  font-size: 2rem;
  color: white;
  margin-bottom: 10px;
}

.footer-logo p {
  font-size: 0.9rem;
  opacity: 0.8;
}

.footer-nav ul {
  display: flex;
  list-style: none;
  flex-wrap: wrap;
  gap: 20px;
}

.footer-nav a {
  color: white;
  font-size: 0.9rem;
  opacity: 0.8;
  transition: all 0.3s;
}

.footer-nav a:hover {
  opacity: 1;
  color: var(--accent-color);
}

.footer-divider {
  height: 1px;
  background-color: rgba(255, 255, 255, 0.1);
  margin-bottom: 20px;
}

.footer-bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.copyright {
  font-size: 0.85rem;
  opacity: 0.7;
}

.back-to-top a {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  color: var(--accent-color);
  font-size: 0.9rem;
  transition: all 0.3s;
}

.back-to-top a:hover {
  transform: translateY(-3px);
}

@media (max-width: 768px) {
  .footer-top {
    flex-direction: column;
    align-items: flex-start;
    gap: 20px;
  }
  
  .footer-nav ul {
    margin-top: 10px;
  }
  
  .footer-bottom {
    flex-direction: column;
    align-items: flex-start;
    gap: 15px;
  }
  
  .back-to-top {
    margin-top: 10px;
  }
}
</style> 