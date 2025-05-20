<template>
  <section id="hero" class="hero-section">
    <div class="container hero-container">
      <div class="hero-content">
        <div class="avatar animate__animated animate__fadeIn">
          <img src="../assets/img/avatar.jpg" alt="个人头像" />
        </div>
        <h1 class="animate__animated animate__fadeInUp">张三</h1>
        <h2 class="animate__animated animate__fadeInUp animate__delay-1s">
          云原生运维工程师
        </h2>
        <div class="typed-text animate__animated animate__fadeInUp animate__delay-2s">
          <span>K8s专家</span> | <span>系统架构师</span> | <span>自动化部署专家</span>
        </div>
        <div class="hero-btns animate__animated animate__fadeInUp animate__delay-3s">
          <a href="#contact" class="btn btn-primary">联系我</a>
          <a href="#about" class="btn btn-outline">了解更多</a>
        </div>
      </div>
      <div class="terminal-container animate__animated animate__fadeIn animate__delay-4s">
        <div class="terminal-header">
          <div class="terminal-buttons">
            <span class="terminal-btn red"></span>
            <span class="terminal-btn yellow"></span>
            <span class="terminal-btn green"></span>
          </div>
          <div class="terminal-title">terminal</div>
        </div>
        <div class="terminal-body">
          <div class="terminal-line">$ whoami</div>
          <div class="terminal-line">张三 - 云原生运维工程师</div>
          <div class="terminal-line">$ ls skills</div>
          <div class="terminal-line">Kubernetes Docker Linux AWS 自动化部署 CI/CD</div>
          <div class="terminal-line">$ cat experience.txt</div>
          <div class="terminal-line">5年运维经验 | 3年K8s实践 | 10+大型项目</div>
          <div class="terminal-line terminal-cursor">$</div>
        </div>
      </div>
    </div>
    <div class="scroll-down animate__animated animate__bounce animate__infinite">
      <a href="#about"><i class="fas fa-chevron-down"></i></a>
    </div>
  </section>
</template>

<script setup>
import { onMounted } from 'vue';
import { gsap } from 'gsap';

onMounted(() => {
  const lines = document.querySelectorAll('.terminal-line');
  
  gsap.from(lines, {
    duration: 0.5,
    opacity: 0,
    y: 20,
    stagger: 0.3,
    delay: 1
  });
});
</script>

<style scoped>
.hero-section {
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--dark-color) 100%);
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  color: var(--light-color);
}

.hero-section::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: url('../assets/img/pattern.png');
  opacity: 0.1;
}

.hero-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  position: relative;
  z-index: 1;
}

.hero-content {
  flex: 1;
  text-align: center;
}

.avatar {
  width: 150px;
  height: 150px;
  border-radius: 50%;
  margin: 0 auto 20px;
  overflow: hidden;
  border: 5px solid rgba(255, 255, 255, 0.2);
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.hero-content h1 {
  font-size: 3.5rem;
  margin-bottom: 10px;
  color: var(--light-color);
}

.hero-content h2 {
  font-size: 2rem;
  margin-bottom: 20px;
  color: var(--accent-color);
}

.typed-text {
  font-size: 1.2rem;
  margin-bottom: 30px;
}

.hero-btns {
  display: flex;
  justify-content: center;
  gap: 20px;
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
}

.btn-primary {
  background-color: var(--accent-color);
  color: var(--light-color);
}

.btn-primary:hover {
  background-color: var(--secondary-color);
  transform: translateY(-3px);
}

.btn-outline {
  border: 2px solid var(--light-color);
  color: var(--light-color);
}

.btn-outline:hover {
  background-color: var(--light-color);
  color: var(--primary-color);
  transform: translateY(-3px);
}

.terminal-container {
  flex: 1;
  max-width: 500px;
  background-color: #1E1E1E;
  border-radius: 6px;
  overflow: hidden;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
  margin-left: 30px;
}

.terminal-header {
  background-color: #323232;
  padding: 10px 15px;
  display: flex;
  align-items: center;
}

.terminal-buttons {
  display: flex;
  gap: 6px;
}

.terminal-btn {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.terminal-btn.red {
  background-color: #FF5F56;
}

.terminal-btn.yellow {
  background-color: #FFBD2E;
}

.terminal-btn.green {
  background-color: #27C93F;
}

.terminal-title {
  margin-left: 30px;
  font-size: 0.8rem;
  color: #CCCCCC;
}

.terminal-body {
  padding: 15px;
  font-family: 'Courier New', monospace;
  font-size: 0.9rem;
  line-height: 1.5;
  color: #FFFFFF;
}

.terminal-line {
  margin-bottom: 8px;
}

.terminal-line:nth-child(even) {
  color: #63E2B7;
}

.terminal-cursor::after {
  content: '▋';
  animation: blink 1s step-end infinite;
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0; }
}

.scroll-down {
  position: absolute;
  bottom: 30px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 1.5rem;
}

.scroll-down a {
  color: var(--light-color);
  opacity: 0.7;
  transition: 0.3s;
}

.scroll-down a:hover {
  opacity: 1;
}

@media (max-width: 992px) {
  .hero-container {
    flex-direction: column;
  }
  
  .hero-content {
    margin-bottom: 40px;
  }
  
  .terminal-container {
    margin-left: 0;
    width: 100%;
  }
}

@media (max-width: 768px) {
  .hero-section {
    padding-top: 80px;
    height: auto;
    min-height: 100vh;
  }
  
  .hero-content h1 {
    font-size: 2.5rem;
  }
  
  .hero-content h2 {
    font-size: 1.5rem;
  }
}
</style> 